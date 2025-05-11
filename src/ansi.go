package main

import "fmt"

func highlightRed(s string) string {
	return fmt.Sprintf("\x1b[%dm\x1b[%dm%s\x1b[0m", 30, 101, s)
}

func highlightYellow(s string) string {
	return fmt.Sprintf("\x1b[%dm\x1b[%dm%s\x1b[0m", 30, 103, s)
}

func highlightGreen(s string) string {
	return fmt.Sprintf("\x1b[%dm\x1b[%dm%s\x1b[0m", 30, 102, s)
}

func highlightCyan(s string) string {
	return fmt.Sprintf("\x1b[%dm\x1b[%dm%s\x1b[0m", 30, 106, s)
}

func highlightBlue(s string) string {
	return fmt.Sprintf("\x1b[%dm\x1b[%dm%s\x1b[0m", 97, 104, s)
}

func highlightMagenta(s string) string {
	return fmt.Sprintf("\x1b[%dm\x1b[%dm%s\x1b[0m", 97, 105, s)
}

func highlightGrey(s string) string {
	return fmt.Sprintf("\x1b[%dm\x1b[%dm%s\x1b[0m", 30, 47, s)
}

func highlightWhite(s string) string {
	return fmt.Sprintf("\x1b[%dm\x1b[%dm%s\x1b[0m", 30, 107, s)
}

// yellow_line := fmt.Sprintf("\x1b[%dm\x1b[%dm%s\x1b[0m", 30, 103, "Yellow")
// green_line := fmt.Sprintf("\x1b[%dm\x1b[%dm%s\x1b[0m", 30, 102, "Green")
// cyan_line := fmt.Sprintf("\x1b[%dm\x1b[%dm%s\x1b[0m", 30, 106, "Cyan")
// blue_line := fmt.Sprintf("\x1b[%dm\x1b[%dm%s\x1b[0m", 97, 104, "Blue")
// magenta_line := fmt.Sprintf("\x1b[%dm\x1b[%dm%s\x1b[0m", 97, 105, "Magenta")
// grey_line := fmt.Sprintf("\x1b[%dm\x1b[%dm%s\x1b[0m", 30, 47, "Grey")
// white_line := fmt.Sprintf("\x1b[%dm\x1b[%dm%s\x1b[0m", 30, 107, "White")
