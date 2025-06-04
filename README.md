# sbar
simple status bar program

## example config
here is my config. put this in `~/.config/sbar/config.yml`

```
delimiter: " ┇ "
delimiter-on-edges: true
tick-rate: 1000
toggle-alt-signal: 3
blocks:
  - interval: 5
    command: [ "bar-net" ]
  - interval: 30
    command: [ "bar-disk" ]
  - interval: 5
    command: [ "bar-hw" ]
  - interval: 10
    command: [ "bar-bat" ]
    in-alt: true
  - interval: 0
    command: [ "bar-vol" ]
    signal: 1
    in-alt: true
  - interval: 0
    command: [ "bar-kb" ]
    signal: 2
  - interval: 10
    command: [ "bar-clock" ]
    in-alt: true
```

> `command` does not spawn a shell, but it does respect `PATH`
