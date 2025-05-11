package main

import (
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

const (
	Normal = iota
	Bold   // bold or increased intensity
	Faint  // faint, decreased intensity or second colour
	Italics
	Underline
	Blink
	FastBlink
	Inverse
	Hidden
	Strikeout
	PrimaryFont
	AltFont1
	AltFont2
	AltFont3
	AltFont4
	AltFont5
	AltFont6
	AltFont7
	AltFont8
	AltFont9
	Gothic // fraktur
	DoubleUnderline
	NormalColor // normal colour or normal intensity (neither bold nor faint)
	NotItalics  // not italicized, not fraktur
	NotUnderlined
	Steady     // not Blink or FastBlink
	Reserved26 // reserved for proportional spacing as specified in CCITT Recommendation T.61
	NotInverse // Positive
	NotHidden  // Revealed
	NotStrikeout
	Black
	Red
	Green
	Yellow
	Blue
	Magenta
	Cyan
	White
	Reserved38 // intended for setting character foreground colour as specified in ISO 8613-6 [CCITT Recommendation T.416]
	Default    // default display colour (implementation-defined)
	BlackBackground
	RedBackground
	GreenBackground
	YellowBackground
	BlueBackground
	MagentaBackground
	CyanBackground
	WhiteBackground
	Reserved48        // reserved for future standardization; intended for setting character background colour as specified in ISO 8613-6 [CCITT Recommendation T.416]
	DefaultBackground // default background colour (implementation-defined)
	Reserved50        // reserved for cancelling the effect of the rendering aspect established by parameter value 26
	Framed
	Encircled
	Overlined
	NotFramed // NotEncircled
	NotOverlined
	Reserved56
	Reserved57
	Reserved58
	Reserved59
	IdeogramUnderline       // ideogram underline or right side line
	IdeogramDoubleUnderline // ideogram double underline or double line on the right side
	IdeogramOverline        // ideogram overline or left side line
	IdeogramDoubleOverline  // ideogram double overline or double line on the left side
	IdeogramStress          // ideogram stress marking
	IdeogramCancel          // cancels the effect of the rendition aspects established by parameter values IdeogramUnderline to IdeogramStress
	reserved66              // This should be 66
)

const (
	NormalText       = "\033[0m" // Turn off all attributes
	BlackText        = "\033[30m"
	RedText          = "\033[31m"
	GreenText        = "\033[32m"
	YellowText       = "\033[33m"
	BlueText         = "\033[34m"
	MagentaText      = "\033[35m"
	CyanText         = "\033[36m"
	WhiteText        = "\033[37m"
	DefaultColorText = "\033[39m" // Normal text color
	BoldText         = "\033[1m"
	BoldBlackText    = "\033[1;30m"
	BoldRedText      = "\033[1;31m"
	BoldGreenText    = "\033[1;32m"
	BoldYellowText   = "\033[1;33m"
	BoldBlueText     = "\033[1;34m"
	BoldMagentaText  = "\033[1;35m"
	BoldCyanText     = "\033[1;36m"
	FaintText        = "\033[2m"
	FaintBlackText   = "\033[2;30m"
	FaintRedText     = "\033[2;31m"
	FaintGreenText   = "\033[2;32m"
	FaintYellowText  = "\033[2;33m"
	FaintBlueText    = "\033[2;34m"
	FaintMagentaText = "\033[2;35m"
	FaintCyanText    = "\033[2;36m"
	FaintWhiteText   = "\033[2;37m"
	DefaultText      = "\033[22;39m" // Normal text color and intensity
)

type Command struct {
	Name        string
	Aliases     []string
	Description string
	Func        func([]string) error
}

func nextCurrAndPrev(i int, options []string) (string, string, string) {
	if i == 0 {
		return options[len(options)-1], options[0], options[1]
	}
	if i == len(options)-1 {
		return options[len(options)-2], options[len(options)-1], options[0]
	}
	return options[i-1], options[i], options[i+1]
}

func doDatetime(args []string) error {
	now := time.Now()
	formattedDateTime := now.Format("Monday, January 2, 2006 at 15:04:05 MST")
	fmt.Printf("The date and time is now %s\n", formattedDateTime)
	return nil
}

func doTime(args []string) error {
	fmt.Print("The time is now ")
	fmt.Println(time.Now().Format("15:04:05"))
	return nil
}

func doDays(args []string) error {
	days := []string{
		highlightRed("Sunday"),
		highlightYellow("Monday"),
		highlightGreen("Tuesday"),
		highlightCyan("Wednesday"),
		highlightBlue("Thursday"),
		highlightMagenta("Friday"),
		highlightWhite("Saturday"),
	}
	for _, day := range days {
		fmt.Println(day)
	}
	fmt.Println()
	yesterday, today, tomorrow := nextCurrAndPrev(int(time.Now().Weekday()), days)
	fmt.Printf("Today is %s\n", today)
	fmt.Printf("Yesterday was %s\n", yesterday)
	fmt.Printf("Tomorrow is %s\n", tomorrow)
	return nil
}

func doMonths(args []string) error {
	months := []string{
		highlightRed("January"),
		highlightYellow("February"),
		highlightGreen("March"),
		highlightCyan("April"),
		highlightBlue("May"),
		highlightMagenta("June"),
		highlightRed("July"),
		highlightYellow("August"),
		highlightGreen("September"),
		highlightCyan("October"),
		highlightBlue("November"),
		highlightMagenta("December"),
	}
	for i, month := range months {
		fmt.Printf("%d. %s\n", i+1, month)
	}
	fmt.Println()
	last, curr, prev := nextCurrAndPrev(int(time.Now().Month())-1, months)
	fmt.Printf("This month is %s\n", curr)
	fmt.Printf("Last month was %s\n", last)
	fmt.Printf("Next month is %s\n", prev)
	return nil
}

func doCal(args []string) error {
	now := time.Now()
	year, month, day := now.Date()
	firstDay := time.Date(year, month, 1, 0, 0, 0, 0, now.Location())
	firstDayWeekday := int(firstDay.Weekday())
	daysInMonth := 32 - time.Date(year, month, 32, 0, 0, 0, 0, now.Location()).Day()
	fmt.Printf("\n%s %d\n", month.String(), year)
	fmt.Println("Sun Mon Tue Wed Thu Fri Sat")
	for i := 0; i < firstDayWeekday; i++ {
		fmt.Print("    ")
	}
	for d := 1; d <= daysInMonth; d++ {
		if d == day {
			// Highlight current day with cyan
			fmt.Printf("%s%2d%s  ", CyanText, d, NormalText)
		} else {
			fmt.Printf("%2d  ", d)
		}

		// Start a new line after Saturday
		if (firstDayWeekday+d)%7 == 0 {
			fmt.Println()
		}
	}
	// Add a final newline if the last day wasn't a Saturday
	if (firstDayWeekday+daysInMonth)%7 != 0 {
		fmt.Println()
	}
	fmt.Println()
	return nil
}

func doNews(args []string) error {
	return nil
}

func doMsg(args []string) error {
	return nil
}

func doBday(args []string) error {
	return nil
}

func doCalc(args []string) error {
	return nil
}

func doBeep(args []string) error {
	os.Stdout.Write([]byte("Beep!\007\n"))
	return nil
}

func doABC(args []string) error {
	fmt.Println("ABCDEFGHIJKLMNOPQRSTUVWXYZ")
	fmt.Println("abcdefghijklmnopqrstuvwxyz")
	return nil
}

func doNum(args []string) error {
	fmt.Println("0123456789")
	fmt.Println()
	fmt.Println("0 = Zero")
	fmt.Println("1 = One")
	fmt.Println("2 = Two")
	fmt.Println("3 = Three")
	fmt.Println("4 = Four")
	fmt.Println("5 = Five")
	fmt.Println("6 = Six")
	fmt.Println("7 = Seven")
	fmt.Println("8 = Eight")
	fmt.Println("9 = Nine")
	fmt.Println("10 = Ten")
	fmt.Println("11 = Eleven")
	fmt.Println("12 = Twelve")
	fmt.Println("20 = Twenty")
	fmt.Println("30 = Thirty")
	fmt.Println("40 = Forty")
	fmt.Println("50 = Fifty")
	fmt.Println("60 = Sixty")
	fmt.Println("70 = Seventy")
	fmt.Println("80 = Eighty")
	fmt.Println("90 = Ninety")
	fmt.Println("100 = One Hundred")
	return nil
}

func doHelp(args []string) error {
	sorted := make([]*Command, 0, len(cmds))
	for name, cmd := range cmds {
		if name != cmd.Name {
			continue
		}
		sorted = append(sorted, cmd)
	}
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Name < sorted[j].Name
	})
	fmt.Printf("%-16s    %-16s    %s\n", "NAME", "ALIASES", "DESCRIPTION")
	fmt.Printf("%-16s    %-16s    %s\n", "====", "=======", "===========")
	for _, cmd := range sorted {
		fmt.Printf("%-16s    %-16s    %s\n", cmd.Name, strings.Join(cmd.Aliases, ","), cmd.Description)
	}
	return nil
}

func doColors(args []string) error {
	red_line := highlightRed("Red")
	yellow_line := highlightYellow("Yellow")
	green_line := highlightGreen("Green")
	cyan_line := highlightCyan("Cyan")
	blue_line := highlightBlue("Blue")
	magenta_line := highlightMagenta("Magenta")
	grey_line := highlightGrey("Grey")
	white_line := highlightWhite("White")
	fmt.Println(red_line)
	fmt.Println(yellow_line)
	fmt.Println(green_line)
	fmt.Println(cyan_line)
	fmt.Println(blue_line)
	fmt.Println(magenta_line)
	fmt.Println(grey_line)
	fmt.Println(white_line)
	return nil
}

func doExit(args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("exactly one argument required, got %d: %q", len(args), args)
	}
	code, err := strconv.ParseInt(args[0], 32, 10)
	if err != nil {
		return err
	}
	os.Exit(int(code))
	return nil
}

func doDate(args []string) error {
	fmt.Print("Today's date is ")
	fmt.Println(time.Now().Format("Monday, January 2, 2006"))
	return nil
}

// func doCompare(args []string) error {
// 	if len(args) == 0 {
// 		fmt.Println("You need to provide arguments to compare")
// 		fmt.Println("For example: 'compare 5 3'")
// 		return nil
// 	}
// 	if len(args) == 1 {
// 		fmt.Println("You need to provide two or more arguments to compare")
// 		fmt.Println("For example: 'compare 5 3'")
// 	}

// 	return nil
// }

var cmds = map[string]*Command{}

func registerCommand(cmd Command) {
	cmds[cmd.Name] = &cmd
	for _, alias := range cmd.Aliases {
		cmds[alias] = &cmd
	}
}
