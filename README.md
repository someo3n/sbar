# sbar
simple status bar program

## example config
here is my config. put this in `~/.config/sbar/config.yml`

```yaml
delimiter: " ┇ "
delimiter-on-edges: true
tick-rate: 1000
blocks:
  - interval: 30
    command: [ "bar-disk" ]
  - interval: 5
    command: [ "bar-net" ]
  - interval: 5
    command: [ "bar-hw" ]
  - interval: 10
    command: [ "bar-bat" ]
  - interval: 0
    command: [ "bar-vol" ]
  - interval: 0
    command: [ "bar-kb" ]
  - interval: 10
    command: [ "bar-clock" ]
```

> `command` does not spawn a shell, but it does respect `PATH`
