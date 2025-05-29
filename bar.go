package main

import (
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

const SIGRTMIN = 34

type Block struct {
	signal   os.Signal
	current  uint
	text     string
	Interval uint     `yaml:"interval"`
	Command  []string `yaml:"command"`
}

func (b *Block) Run() {
	command := exec.Command(b.Command[0], b.Command[1:]...)
	output, err := command.Output()
	if err != nil {
		return
	}
	b.text = strings.TrimSpace(string(output))
}

func (b *Block) Tick() bool {
	if b.Interval == 0 {
		return false
	}
	b.current++
	if b.current >= b.Interval {
		b.current = 0
		return true
	}
	return false
}

type Bar struct {
	Blocks              []*Block `yaml:"blocks"`
	Delimiter           string   `yaml:"delimiter"`
	AddDelimiterOnEdges bool     `yaml:"delimiter-on-edges"`
	TickRate            uint     `yaml:"tick-rate"`

	signalMap map[os.Signal]*Block
}

func (b Bar) String() string {
	var builder strings.Builder

	if b.AddDelimiterOnEdges {
		builder.WriteString(b.Delimiter)
	}

	for _, block := range b.Blocks {
		builder.WriteString(block.text)
		builder.WriteString(b.Delimiter)
	}

	return builder.String()
}

func (b *Bar) Tick() {
	changed := false
	for _, block := range b.Blocks {
		if block.Tick() {
			block.Run()
			changed = true
		}
	}
	if changed {
		b.Update()
	}
}

func (b *Bar) Update() {
	fmt.Println(b.String())
}

func (b *Bar) Loop() {
	b.setupSignals()
	b.setup()

	tickDuration := time.Duration(b.TickRate) * time.Millisecond
	for {
		b.Tick()
		time.Sleep(tickDuration)
	}
}

func (b *Bar) setup() {
	for _, block := range b.Blocks {
		block.Run()
	}

	b.Update()
}

func (b *Bar) setupSignals() {
	b.signalMap = make(map[os.Signal]*Block)
	signals := []os.Signal{}

	for idx, block := range b.Blocks {
		sig := syscall.Signal(SIGRTMIN + idx)
		block.signal = sig
		b.signalMap[sig] = block
		signals = append(signals, sig)
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, signals...)

	go func() {
		for sig := range sigChan {
			if block, ok := b.signalMap[sig]; ok {
				block.Run()
				b.Update()
			}
		}
	}()
}
