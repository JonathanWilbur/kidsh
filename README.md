# Kid Shell

**Work in progress. This whole readme is likely to be wrong.**

This is a lightweight / low-capability secure shell for young children to fool
around on a computer. It has simple commands for doing easy tasks in a simple
way that children can understand, such as pressing `c` to display colors, or
`d` to display days of the week, or `n` to take notes. It compiles as a static
binary that only supports built-ins, and (not implemented yet) it will be able
to ran with dropped capabilities, so your little one cannot accidentally do
anything to mess up your computer.

## Build

```bash
go build -o kidsh src/main.go src/cmd.go src/ansi.go src/config.go
```

## Usage

I expect users to full-screen the window where this shell is running so their
kids cannot easily escape it. Even better would be to define a boot menu
configuration where Linux defines `kidsh` as PID 0, so they cannot escape it.
Alternatively, you could just open up this shell in a real terminal that is not
running a GUI.
