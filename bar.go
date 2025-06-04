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

const SIGRTMIN = 34 // Standard SIGRTMIN on Linux

type Block struct {
	Signal   *uint8   `yaml:"signal"` 
	Interval uint     `yaml:"interval"`
	Command  []string `yaml:"command"`
	InAlt    bool     `yaml:"in-alt"`
	current  uint
	text     string
}

func (b *Block) Run() {
	cmd := exec.Command(b.Command[0], b.Command[1:]...)
	output, err := cmd.Output()
	if err != nil {
		b.text = "<err>"
		fmt.Fprintf(os.Stderr, "error running %s: %v\n", b.Command[0], err)
		return
	}
	b.text = strings.TrimSpace(string(output))
}

func (b *Block) Tick() bool {
	if b.Interval == 0 {
		return false
	}
	b.current = (b.current + 1) % b.Interval
	return b.current == 0
}

type Bar struct {
	Blocks              []*Block `yaml:"blocks"`
	Delimiter           string   `yaml:"delimiter"`
	AddDelimiterOnEdges bool     `yaml:"delimiter-on-edges"`
	TickRate            uint     `yaml:"tick-rate"`
	ToggleAltSignal     *uint8   `yaml:"toggle-alt-signal"`
	signalMap           map[os.Signal]*Block
	showingAlt          bool
}

func (b *Bar) getActiveBlocks() []*Block {
	if !b.showingAlt {
		return b.Blocks
	}

	var altBlocks []*Block
	for _, block := range b.Blocks {
		if block.InAlt {
			altBlocks = append(altBlocks, block)
		}
	}
	return altBlocks
}

func (b Bar) String() string {
	var builder strings.Builder
	activeBlocks := b.getActiveBlocks()

	if b.AddDelimiterOnEdges {
		builder.WriteString(b.Delimiter)
	}
	for i, block := range activeBlocks {
		builder.WriteString(block.text)
		if i < len(activeBlocks)-1 || b.AddDelimiterOnEdges {
			builder.WriteString(b.Delimiter)
		}
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

func (b *Bar) toggleAlt() {
	b.showingAlt = !b.showingAlt
	b.Update()
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
	var signals []os.Signal

	for _, block := range b.Blocks {
		if block.Signal != nil {
			sig := syscall.Signal(SIGRTMIN + int(*block.Signal))
			b.signalMap[sig] = block
			signals = append(signals, sig)
		}
	}

	if b.ToggleAltSignal != nil && *b.ToggleAltSignal != 0 {
		toggleSig := syscall.Signal(SIGRTMIN + int(*b.ToggleAltSignal))
		signals = append(signals, toggleSig)
	}

	if len(signals) == 0 {
		return
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, signals...)
	go func() {
		for sig := range sigChan {
			if b.ToggleAltSignal != nil && *b.ToggleAltSignal != 0 &&
				sig == syscall.Signal(SIGRTMIN+int(*b.ToggleAltSignal)) {
				b.toggleAlt()
				continue
			}

			if block, ok := b.signalMap[sig]; ok {
				block.Run()
				b.Update()
			}
		}
	}()
}
