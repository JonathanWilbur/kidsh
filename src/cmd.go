package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net"
	"net/http"
	"os"
	"os/exec"
	"slices"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/emersion/go-vcard"
	"github.com/lukechampine/nock"
	"github.com/mmcdole/gofeed"
)

// BibleVerse represents a single verse from the Bible
type BibleVerse struct {
	BookID   string `json:"book_id"`
	BookName string `json:"book_name"`
	Chapter  int    `json:"chapter"`
	Verse    int    `json:"verse"`
	Text     string `json:"text"`
}

// BibleResponse represents the complete response from the Bible API
type BibleResponse struct {
	Reference       string       `json:"reference"`
	Verses          []BibleVerse `json:"verses"`
	Text            string       `json:"text"`
	TranslationID   string       `json:"translation_id"`
	TranslationName string       `json:"translation_name"`
	TranslationNote string       `json:"translation_note"`
}

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

const stackEnv = "KIDSH_STACK"
const queueEnv = "KIDSH_QUEUE"

// TODO: Make these files a full path given by an environment variable.
const todoFile = "todo.db"
const contactsFile = "contacts.vcf"
const myName = "John Doe"

const separator = "\x1E" // ASCII Record Separator (RS)

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
	rssURL := os.Getenv("KIDSH_RSS_URL")
	if rssURL == "" {
		return fmt.Errorf("KIDSH_RSS_URL environment variable not set")
	}

	fp := gofeed.NewParser()
	feed, err := fp.ParseURL(rssURL)
	if err != nil {
		return fmt.Errorf("error parsing RSS feed: %v", err)
	}

	fmt.Printf("%sLatest News from %s%s\n\n", BoldGreenText, feed.Title, NormalText)
	
	// Show last 5 items
	maxItems := 5
	if len(feed.Items) < maxItems {
		maxItems = len(feed.Items)
	}

	for i := 0; i < maxItems; i++ {
		item := feed.Items[i]
		fmt.Printf("%s%d. %s%s\n", BoldBlueText, i+1, item.Title, NormalText)
		if item.Description != "" {
			fmt.Printf("   %s\n", item.Description)
		}
		if item.Link != "" {
			fmt.Printf("   %sLink: %s%s\n", FaintText, item.Link, NormalText)
		}
		if item.Published != "" {
			fmt.Printf("   %sPublished: %s%s\n", FaintText, item.Published, NormalText)
		}
		fmt.Println()
	}

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
		os.Exit(0)
		return nil
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

func doCompare(args []string) error {
	if len(args) < 2 {
		fmt.Println("You have to type in more than one number, silly!")
		return nil
	}

	// Parse all arguments as integers
	nums := make([]int, 0, len(args))
	for _, arg := range args {
		num, err := strconv.Atoi(arg)
		if err != nil {
			return fmt.Errorf("'%s' is not a valid number", arg)
		}
		nums = append(nums, num)
	}

	// Make a copy for sorting
	sortedNums := make([]int, len(nums))
	copy(sortedNums, nums)
	sort.Ints(sortedNums)

	// Special case for exactly two numbers
	if len(nums) == 2 {
		if nums[0] == nums[1] {
			fmt.Printf("%d is equal to %d\n", nums[0], nums[1])
		} else if nums[0] > nums[1] {
			fmt.Printf("%d is larger than %d\n", nums[0], nums[1])
		} else {
			fmt.Printf("%d is larger than %d\n", nums[1], nums[0])
		}
		return nil
	}

	// Case for more than two numbers
	smallest := sortedNums[0]
	largest := sortedNums[len(sortedNums)-1]

	fmt.Printf("The smallest number is %d\n", smallest)
	fmt.Printf("The largest number is %d\n", largest)

	// Print the numbers in sorted order
	fmt.Print("Numbers in ascending order: ")
	for i, num := range sortedNums {
		fmt.Printf("%d", num)
		if i < len(sortedNums)-1 {
			fmt.Print(", ")
		}
	}
	fmt.Println()
	return nil
}

func doCount(args []string) error {
	if len(args) != 1 && (len(args) == 2 && args[0] != "to") {
		fmt.Println("You have to tell me what number to count to, Silly!")
		return nil
	}
	arg := args[len(args)-1]
	num, err := strconv.Atoi(arg)
	if err != nil {
		return fmt.Errorf("'%s' is not a valid number", arg)
	}
	if num > 100 {
		fmt.Println("That number is too big. Try a smaller number.")
		return nil
	}
	for i := 0; i <= num; i++ {
		fmt.Printf("%d ", i)
	}
	fmt.Println()
	return nil
}

func doSortLex(args []string) error {
	sort.Strings(args)
	for _, s := range args {
		fmt.Printf("%s ", s)
	}
	fmt.Println()
	return nil
}

func doSort(args []string) error {
	// Parse all arguments as integers
	nums := make([]int, 0, len(args))
	for _, arg := range args {
		num, err := strconv.Atoi(arg)
		if err != nil {
			// If any fail, assume we want to do a lexicographic sort.
			return doSortLex(args)
		}
		nums = append(nums, num)
	}

	sortedNums := make([]int, len(nums))
	copy(sortedNums, nums)
	sort.Ints(sortedNums)
	for i, num := range sortedNums {
		fmt.Printf("%d", num)
		if i < len(sortedNums)-1 {
			fmt.Print(", ")
		}
	}
	fmt.Println()
	return nil
}

func doUniq(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("no arguments provided to uniq command")
	}
	seen := make(map[string]bool)
	unique := []string{}
	for _, item := range args {
		if !seen[item] {
			seen[item] = true
			unique = append(unique, item)
		}
	}
	fmt.Println("Unique items:")
	for _, item := range unique {
		fmt.Println(item)
	}
	fmt.Printf("\nFound %d unique items from %d total items\n", len(unique), len(args))
	if len(unique) < len(args) {
		fmt.Printf("Removed %d duplicate(s)\n", len(args)-len(unique))
	}
	return nil
}

func doPwd(args []string) error {
	dir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current directory: %v", err)
	}
	fmt.Println(dir)
	return nil
}

func doCd(args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("exactly one argument required, got %d: %q", len(args), args)
	}
	err := os.Chdir(args[0])
	if err != nil {
		return fmt.Errorf("failed to change directory to '%s': %v", args[0], err)
	}
	return nil
}

func doLs(args []string) error {
	files, err := os.ReadDir(".")
	if err != nil {
		return fmt.Errorf("failed to list directory contents: %v", err)
	}

	for _, file := range files {
		if file.IsDir() {
			fmt.Printf("%s/\n", file.Name())
		} else {
			fmt.Println(file.Name())
		}
	}
	return nil
}

func doFirst(args []string) error {
	if len(args) == 0 {
		fmt.Println("You need to supply an argument, silly!")
		return nil
	}
	fmt.Println(args[0])
	return nil
}

func doLast(args []string) error {
	if len(args) == 0 {
		fmt.Println("you need to supply an argument, silly!")
		return nil
	}
	arg := args[len(args)-1]
	fmt.Println(arg)
	return nil
}

func doRev(args []string) error {
	slices.Reverse(args)
	for _, arg := range args {
		fmt.Printf("%s ", arg)
	}
	fmt.Println()
	return nil
}

func doAdd(args []string) error {
	sum := 0
	for _, arg := range args {
		num, err := strconv.Atoi(arg)
		if err != nil {
			return err
		}
		sum += num
	}
	fmt.Printf("The total is: %d\n", sum)
	return nil
}

func doMultiply(args []string) error {
	if len(args) < 1 {
		fmt.Println("you need to supply an argument, silly!")
		return nil
	}
	result, err := strconv.Atoi(args[0])
	if err != nil {
		return err
	}
	for _, arg := range args[1:] {
		num, err := strconv.Atoi(arg)
		if err != nil {
			return err
		}
		result *= num
	}
	fmt.Printf("The product is: %d\n", result)
	return nil
}

const DEFAULT_WEATHER_URL = "https://wttr.in/St.%20Johns,%20Florida?format=3&u"

// const DEFAULT_WEATHER_URL = "https://wttr.in/St.%20Johns,%20Florida?2Anu"

func doWeather(args []string) error {
	weatherUrl := os.Getenv("WEATHER_URL")
	if len(weatherUrl) < 8 { // Could not be valid if smaller.
		weatherUrl = DEFAULT_WEATHER_URL
	}
	resp, err := http.Get(weatherUrl)
	if err != nil {
		return fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error reading response: %v", err)
	}
	fmt.Println(string(body))
	return nil
}

func doUpper(args []string) error {
	for i, arg := range args {
		fmt.Print(strings.ToUpper(arg))
		if i < len(args)-1 {
			fmt.Print(" ")
		}
	}
	fmt.Println()
	return nil
}

func doLower(args []string) error {
	for i, arg := range args {
		fmt.Print(strings.ToLower(arg))
		if i < len(args)-1 {
			fmt.Print(" ")
		}
	}
	fmt.Println()
	return nil
}

func doEnv(args []string) error {
	envVars := os.Environ()
	sort.Strings(envVars)
	for _, env := range envVars {
		fmt.Println(env)
	}
	return nil
}

func doShuffle(args []string) error {
	shuffled := make([]string, len(args))
	copy(shuffled, args)

	// Fisher-Yates shuffle algorithm
	for i := len(shuffled) - 1; i > 0; i-- {
		j := time.Now().UnixNano() % int64(i+1)
		shuffled[i], shuffled[int(j)] = shuffled[int(j)], shuffled[i]
	}

	for i, arg := range shuffled {
		fmt.Print(arg)
		if i < len(shuffled)-1 {
			fmt.Print(" ")
		}
	}
	fmt.Println()
	return nil
}

func doRandom(args []string) error {
	max := 100 // Default max value

	// If an argument is provided, use it as the max value
	if len(args) > 0 {
		var err error
		max, err = strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("invalid number '%s': %v", args[0], err)
		}
		if max <= 0 {
			return fmt.Errorf("maximum value must be positive, got %d", max)
		}
	}

	// Generate random number between 0 and max
	source := time.Now().UnixNano()
	randomNum := int(source % int64(max+1))

	fmt.Printf("Random number (0-%d): %d\n", max, randomNum)
	return nil
}

func doFlip(args []string) error {
	source := time.Now().UnixNano()
	result := "Tails"
	if source%2 == 0 {
		result = "Heads"
	}
	fmt.Printf("The coin flip result is: %s\n", result)
	return nil
}

func doSleep(args []string) error {
	// Default sleep time is 1 second
	sleepTime := 1.0

	// Parse duration if provided as argument
	if len(args) > 0 {
		var err error
		sleepTime, err = strconv.ParseFloat(args[0], 64)
		if err != nil {
			return fmt.Errorf("invalid sleep duration '%s': %v", args[0], err)
		}
		if sleepTime < 0 {
			return fmt.Errorf("sleep duration cannot be negative, got %f", sleepTime)
		}
	}

	// Convert to duration and sleep
	duration := time.Duration(sleepTime * float64(time.Second))
	if duration > (10 * time.Second) {
		return fmt.Errorf("that's too long")
	}
	time.Sleep(duration)
	return nil
}

func doReset(args []string) error {
	// ANSI escape sequence to reset the terminal
	resetSequence := "\033c"
	fmt.Print(resetSequence)
	return nil
}

func doCompass(args []string) error {
	fmt.Println("   NW   N    NE  ")
	fmt.Println("        |        ")
	fmt.Println("   W ---+--- E   ")
	fmt.Println("        |        ")
	fmt.Println("   SW   S    SE  ")
	fmt.Println()
	fmt.Println("N = North")
	fmt.Println("E = East")
	fmt.Println("S = South")
	fmt.Println("W = West")
	fmt.Println()
	fmt.Println("NE = Northeast")
	fmt.Println("SE = Southeast")
	fmt.Println("SW = Southwest")
	fmt.Println("NW = Northwest")
	return nil
}

func doIp(args []string) error {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return fmt.Errorf("failed to get interface addresses: %v", err)
	}
	fmt.Println("My IP Addresses:")
	for _, addr := range addrs {
		// Check if this is an IP network address
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			// Show IPv4 addresses
			if ipnet.IP.To4() != nil {
				fmt.Printf("  %s\n", ipnet.IP.String())
			} else if ipnet.IP.To16() != nil {
				fmt.Printf("  %s\n", ipnet.IP.String())
			}
		}
	}
	return nil
}

func doSeasons(args []string) error {
	// Define emoji and colors for each season
	spring := fmt.Sprintf("\033[42;37m Spring \033[0m") // Green background, white text
	summer := fmt.Sprintf("\033[43;30m Summer \033[0m") // Yellow background, black text
	autumn := fmt.Sprintf("\033[41;37m Autumn \033[0m") // Red background, white text
	winter := fmt.Sprintf("\033[44;37m Winter \033[0m") // Blue background, white text

	// Get current season in Northern Hemisphere
	now := time.Now()
	month := now.Month()
	var currentSeason string

	switch {
	case month >= 3 && month <= 5:
		currentSeason = spring
	case month >= 6 && month <= 8:
		currentSeason = summer
	case month >= 9 && month <= 11:
		currentSeason = autumn
	default:
		currentSeason = winter
	}

	fmt.Printf("Current season (Northern Hemisphere): %s\n", currentSeason)
	fmt.Println("The Four Seasons:")
	fmt.Println(spring)
	fmt.Println("  - March, April, May (Northern Hemisphere)")
	fmt.Println("  - September, October, November (Southern Hemisphere)")
	fmt.Println(summer)
	fmt.Println("  - June, July, August (Northern Hemisphere)")
	fmt.Println("  - December, January, February (Southern Hemisphere)")
	fmt.Println(autumn)
	fmt.Println("  - September, October, November (Northern Hemisphere)")
	fmt.Println("  - March, April, May (Southern Hemisphere)")
	fmt.Println(winter)
	fmt.Println("  - December, January, February (Northern Hemisphere)")
	fmt.Println("  - June, July, August (Southern Hemisphere)")
	return nil
}

func doUptime(args []string) error {
	// Read uptime info from /proc/uptime on Linux
	data, err := os.ReadFile("/proc/uptime")
	if err != nil {
		return fmt.Errorf("failed to read uptime: %v", err)
	}

	// Parse uptime - first value is total uptime in seconds
	fields := strings.Fields(string(data))
	if len(fields) < 1 {
		return fmt.Errorf("invalid uptime data")
	}

	uptimeSeconds, err := strconv.ParseFloat(fields[0], 64)
	if err != nil {
		return fmt.Errorf("failed to parse uptime: %v", err)
	}

	// Calculate days, hours, minutes, seconds
	days := int(uptimeSeconds / 86400)
	hours := int((uptimeSeconds - float64(days)*86400) / 3600)
	minutes := int((uptimeSeconds - float64(days)*86400 - float64(hours)*3600) / 60)
	seconds := int(uptimeSeconds - float64(days)*86400 - float64(hours)*3600 - float64(minutes)*60)

	fmt.Print("System uptime: ")
	if days > 0 {
		fmt.Printf("%d day(s), ", days)
	}
	fmt.Printf("%02d:%02d:%02d\n", hours, minutes, seconds)

	// Get current time
	now := time.Now()
	fmt.Printf("Current time: %s\n", now.Format("15:04:05"))

	return nil
}

func doPush(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("no values provided to push")
	}

	current := os.Getenv(stackEnv)
	newEntries := strings.Join(args, separator)

	var newStack string
	if current == "" {
		newStack = newEntries
	} else {
		newStack = current + separator + newEntries
	}

	return os.Setenv(stackEnv, newStack)
}

func doPop(args []string) error {
	current := os.Getenv(stackEnv)
	if current == "" {
		return fmt.Errorf("stack is empty")
	}

	parts := strings.Split(current, separator)
	if len(parts) == 0 {
		return fmt.Errorf("stack is empty")
	}

	popped := parts[len(parts)-1]
	remaining := parts[:len(parts)-1]

	if len(remaining) == 0 {
		_ = os.Unsetenv(stackEnv)
	} else {
		_ = os.Setenv(stackEnv, strings.Join(remaining, separator))
	}

	fmt.Println(popped)
	return nil
}

func doPrintStack(args []string) error {
	current := os.Getenv(stackEnv)
	if current == "" {
		fmt.Println("Stack is empty.")
		return nil
	}

	parts := strings.Split(current, separator)

	// A few readable ANSI foreground colors (30–37, skipping black)
	colors := []int{31, 32, 33, 34, 35, 36, 37}

	for i, val := range parts {
		color := colors[i%len(colors)]
		fmt.Printf("\033[%dm[%d] %s\033[0m\n", color, i, val)
	}
	return nil
}

func doEnqueue(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("no values provided to enqueue")
	}

	current := os.Getenv(queueEnv)
	newEntries := strings.Join(args, separator)

	var newQueue string
	if current == "" {
		newQueue = newEntries
	} else {
		newQueue = current + separator + newEntries
	}

	return os.Setenv("KIDSH_QUEUE", newQueue)
}

func doDequeue(args []string) error {
	current := os.Getenv(queueEnv)
	if current == "" {
		return fmt.Errorf("queue is empty")
	}

	parts := strings.Split(current, separator)
	if len(parts) == 0 {
		return fmt.Errorf("queue is empty")
	}

	dequeued := parts[0]
	remaining := parts[1:]

	if len(remaining) == 0 {
		_ = os.Unsetenv(queueEnv)
	} else {
		_ = os.Setenv(queueEnv, strings.Join(remaining, separator))
	}

	fmt.Println(dequeued)
	return nil
}

func doPrintQueue(args []string) error {
	current := os.Getenv(queueEnv)
	if current == "" {
		fmt.Println("Queue is empty.")
		return nil
	}

	parts := strings.Split(current, separator)
	colors := []int{31, 32, 33, 34, 35, 36, 37}

	for i, val := range parts {
		color := colors[i%len(colors)]
		fmt.Printf("\033[%dm[%d] %s\033[0m\n", color, i, val)
	}
	return nil
}

func readTodos() ([]string, error) {
	data, err := os.ReadFile(todoFile)
	if err != nil && !os.IsNotExist(err) {
		return nil, err
	}
	if len(data) == 0 {
		return []string{}, nil
	}
	return strings.Split(string(data), separator), nil
}

func writeTodos(todos []string) error {
	if len(todos) == 0 {
		return os.Remove(todoFile)
	}
	return os.WriteFile(todoFile, []byte(strings.Join(todos, separator)), 0644)
}

func doTodo(args []string) error {
	todos, err := readTodos()
	if err != nil {
		return err
	}

	if len(args) == 0 {
		if len(todos) == 0 {
			fmt.Println("No todos.")
			return nil
		}
		colors := []int{31, 32, 33, 34, 35, 36, 37}
		for i, val := range todos {
			color := colors[i%len(colors)]
			fmt.Printf("\033[%dm[%d] %s\033[0m\n", color, i, val)
		}
		return nil
	}

	newItem := strings.Join(args, " ")
	todos = append(todos, newItem)
	return writeTodos(todos)
}

func doDone(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("specify index or prefix to mark done")
	}

	todos, err := readTodos()
	if err != nil {
		return err
	}

	if len(todos) == 0 {
		return fmt.Errorf("no todos to mark done")
	}

	target := strings.ToLower(strings.Join(args, " "))
	idx := -1

	// Try numeric index
	if n, err := strconv.Atoi(target); err == nil && n >= 0 && n < len(todos) {
		idx = n
	} else {
		// Try prefix match
		for i, todo := range todos {
			if strings.HasPrefix(strings.ToLower(todo), target) {
				idx = i
				break
			}
		}
	}

	if idx == -1 || idx >= len(todos) {
		return fmt.Errorf("todo not found: %s", target)
	}

	done := todos[idx]
	todos = append(todos[:idx], todos[idx+1:]...)
	if err := writeTodos(todos); err != nil {
		return err
	}

	fmt.Printf("Done: %s\n", done)
	return nil
}

func doHomeAddress(args []string) error {
	f, err := os.Open(contactsFile)
	if err != nil {
		return err
	}
	defer f.Close()

	dec := vcard.NewDecoder(f)
	for {
		card, err := dec.Decode()
		if err != nil {
			if err.Error() == "EOF" {
				break
			}
			return err
		}

		if fn := card.PreferredValue(vcard.FieldFormattedName); fn == myName {
			addresses := card[vcard.FieldAddress]
			for _, a := range addresses {
				if strings.Contains(strings.ToLower(a.Params.Get("TYPE")), "home") {
					parts := strings.Split(a.Value, ";")
					labels := []string{
						"PO Box",
						"Extended Address",
						"Street Address",
						"Locality",
						"Region",
						"Postal Code",
						"Country",
					}

					fmt.Printf("%sMy home address is:%s\n", BoldGreenText, NormalText)
					lines := []string{}
					if parts[0] != "" {
						lines = append(lines, parts[0]) // PO Box
					}
					if parts[1] != "" {
						lines = append(lines, parts[1]) // Extended
					}
					if parts[2] != "" {
						lines = append(lines, parts[2]) // Street
					}
					line2 := strings.TrimSpace(strings.Join([]string{parts[3], parts[4], parts[5]}, " "))
					if line2 != "" {
						lines = append(lines, line2) // City Region Postal
					}
					if parts[6] != "" {
						lines = append(lines, parts[6]) // Country
					}
					for _, l := range lines {
						fmt.Println(l)
					}

					fmt.Printf("\n%sFormatted Address Fields:%s\n", BoldGreenText, NormalText)
					maxLabelLen := 0
					for _, label := range labels {
						if len(label) > maxLabelLen {
							maxLabelLen = len(label)
						}
					}
					for i := 0; i < len(parts) && i < len(labels); i++ {
						if parts[i] != "" {
							fmt.Printf("%-*s: %s\n", maxLabelLen, labels[i], parts[i])
						}
					}

					return nil
				}
			}
			return fmt.Errorf("home address not found for %s", myName)
		}
	}

	return fmt.Errorf("contact not found: %s", myName)
}

func doBirthday(args []string) error {
	f, err := os.Open(contactsFile)
	if err != nil {
		return err
	}
	defer f.Close()

	dec := vcard.NewDecoder(f)
	for {
		card, err := dec.Decode()
		if err != nil {
			if err.Error() == "EOF" {
				break
			}
			return err
		}

		if fn := card.PreferredValue(vcard.FieldFormattedName); fn == myName {
			bdayRaw := card.PreferredValue(vcard.FieldBirthday)
			if bdayRaw == "" {
				return fmt.Errorf("birthday not found for %s", myName)
			}

			var bday time.Time
			var parsed bool

			// Try full date first
			if t, err := time.Parse("2006-01-02", bdayRaw); err == nil {
				bday = t
				parsed = true
			} else if t, err := time.Parse("--01-02", bdayRaw); err == nil {
				// vCard 4.0 allows partial date like --MM-DD
				bday = t
				parsed = true
			}

			if parsed {
				fmt.Printf("My birthday is: %s\n", bday.Format("January 2, 2006"))

				now := time.Now()
				if bday.Month() == now.Month() && bday.Day() == now.Day() {
					fmt.Println("Today is your birthday! Yay!")
				}
			} else {
				fmt.Printf("My birthday is: %s\n", bdayRaw)
			}

			return nil
		}
	}

	// TODO: Print out birthdays of family members

	return fmt.Errorf("contact not found: %s", myName)
}

func doAge(args []string) error {
	f, err := os.Open(contactsFile)
	if err != nil {
		return err
	}
	defer f.Close()

	dec := vcard.NewDecoder(f)
	for {
		card, err := dec.Decode()
		if err != nil {
			if err.Error() == "EOF" {
				break
			}
			return err
		}

		if fn := card.PreferredValue(vcard.FieldFormattedName); fn == myName {
			bdayRaw := card.PreferredValue(vcard.FieldBirthday)
			if bdayRaw == "" {
				return fmt.Errorf("birthday not found for %s", myName)
			}

			bday, err := time.Parse("2006-01-02", bdayRaw)
			if err != nil {
				return fmt.Errorf("could not parse birthday: %s", bdayRaw)
			}

			now := time.Now()
			age := now.Year() - bday.Year()
			if now.Month() < bday.Month() || (now.Month() == bday.Month() && now.Day() < bday.Day()) {
				age-- // hasn't had birthday yet this year
			}

			fmt.Printf("My age is: %d\n", age)
			return nil
		}
	}

	return fmt.Errorf("contact not found: %s", myName)
}

func doCountdown(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("please provide a starting number")
	}

	start, err := strconv.Atoi(args[0])
	if err != nil || start < 0 {
		return fmt.Errorf("invalid number: %s", args[0])
	}

	if start > 10 {
		// This is so the kid doesn't lock themselves out of the command line forever.
		return fmt.Errorf("that's too long")
	}

	for i := start; i >= 0; i-- {
		// ANSI: clear line (\033[2K) and return carriage (\r)
		fmt.Fprintf(os.Stdout, "\033[2K\r%d", i)
		time.Sleep(1 * time.Second)
	}
	fmt.Fprintf(os.Stdout, "\033[2K\rGo!\n")
	return nil
}

func doNock(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("please provide a Nock expression")
	}
	src := strings.Join(args, " ")
	expr := nock.Parse(src)
	fmt.Println(expr.String())
	return nil
}

func doRepeat(args []string) error {
	if len(args) < 2 {
		return fmt.Errorf("usage: repeat <count> <text...>")
	}

	count, err := strconv.Atoi(args[0])
	if err != nil || count < 0 {
		return fmt.Errorf("invalid repeat count: %s", args[0])
	}

	line := strings.Join(args[1:], " ")
	for i := 0; i < count; i++ {
		fmt.Println(line)
	}
	return nil
}

func doSubtract(args []string) error {
	if len(args) < 2 {
		return fmt.Errorf("please provide two numbers: minuend subtrahend")
	}

	a, err1 := strconv.Atoi(args[0])
	b, err2 := strconv.Atoi(args[1])
	if err1 != nil || err2 != nil {
		return fmt.Errorf("both arguments must be integers")
	}

	fmt.Println(a - b)
	return nil
}

func doCountGame(args []string) error {
	count := rand.Intn(9) + 1 // 1–9

	// TODO: Format to be a little more readable
	fmt.Println(strings.Repeat("O", count))

	fmt.Print("How many Os? ")
	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		return fmt.Errorf("failed to read input: %v", err)
	}

	input = strings.TrimSpace(input)
	n, err := strconv.Atoi(input)
	if err != nil {
		fmt.Println("Please enter a number.")
		return nil
	}

	// TODO: color
	if n == count {
		fmt.Println("That's correct!")
	} else {
		fmt.Println("That is incorrect.")
	}
	return nil
}

func isTextFile(data []byte) bool {
	// Check first 512 bytes or entire file if smaller
	checkLen := 512
	if len(data) < checkLen {
		checkLen = len(data)
	}
	
	// If file is empty, consider it text
	if checkLen == 0 {
		return true
	}

	// Count printable ASCII characters
	printableCount := 0
	for i := 0; i < checkLen; i++ {
		if data[i] >= 32 && data[i] <= 126 || data[i] == '\n' || data[i] == '\r' || data[i] == '\t' {
			printableCount++
		}
	}

	// If more than 80% of characters are printable, consider it text
	return float64(printableCount)/float64(checkLen) > 0.8
}

func doCat(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("please specify a file to display")
	}

	filename := args[0]
	data, err := os.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("error reading file: %v", err)
	}

	if !isTextFile(data) {
		return fmt.Errorf("warning: this appears to be a binary file. Use a different tool to view it.")
	}

	fmt.Print(string(data))
	return nil
}

func doAnd(args []string) error {
	if len(args) != 2 {
		return fmt.Errorf("and requires exactly 2 operands")
	}
	
	op1 := isTruthy(args[0])
	op2 := isTruthy(args[1])
	
	result := op1 && op2
	fmt.Printf("%s AND %s = %s\n", args[0], args[1], boolToStr(result))
	return nil
}

func doOr(args []string) error {
	if len(args) != 2 {
		return fmt.Errorf("or requires exactly 2 operands")
	}
	
	op1 := isTruthy(args[0])
	op2 := isTruthy(args[1])
	
	result := op1 || op2
	fmt.Printf("%s OR %s = %s\n", args[0], args[1], boolToStr(result))
	return nil
}

func doXor(args []string) error {
	if len(args) != 2 {
		return fmt.Errorf("xor requires exactly 2 operands")
	}
	
	op1 := isTruthy(args[0])
	op2 := isTruthy(args[1])
	
	result := op1 != op2
	fmt.Printf("%s XOR %s = %s\n", args[0], args[1], boolToStr(result))
	return nil
}

func doNot(args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("not requires exactly 1 operand")
	}
	
	op := isTruthy(args[0])
	result := !op
	
	fmt.Printf("NOT %s = %s\n", args[0], boolToStr(result))
	return nil
}

func isTruthy(s string) bool {
	s = strings.ToLower(strings.TrimSpace(s))
	return s == "true" || s == "1"
}

func boolToStr(b bool) string {
	if b {
		return "true"
	}
	return "false"
}

func doFamily(args []string) error {
	return doCat([]string{"family.txt"})
}

func doBedtime(args []string) error {
	now := time.Now()
	// FIXME: get bedtime from config
	bedtime := time.Date(now.Year(), now.Month(), now.Day(), 21, 0, 0, 0, now.Location())
	
	// If it's already past bedtime, show tomorrow's bedtime
	if now.After(bedtime) {
		bedtime = bedtime.Add(24 * time.Hour)
	}
	
	timeUntilBedtime := bedtime.Sub(now)
	hours := int(timeUntilBedtime.Hours())
	minutes := int(timeUntilBedtime.Minutes()) % 60
	
	fmt.Printf("Bedtime is at %s\n", bedtime.Format("3:04 PM"))
	fmt.Printf("Time until bedtime: %d hours and %d minutes\n", hours, minutes)
	return nil
}

func doPrintOut(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("printOut requires at least one argument")
	}
	
	// Concatenate all arguments with spaces
	text := strings.Join(args, " ")
	
	// Create command to pipe to lpr
	cmd := exec.Command("lpr")
	cmd.Stdin = strings.NewReader(text)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	
	return cmd.Run()
}

func doSpeak(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("speak requires at least one argument")
	}
	
	// Check if espeak is available
	_, err := exec.LookPath("espeak")
	if err != nil {
		return fmt.Errorf("espeak command not found: %v", err)
	}
	
	// Concatenate all arguments with spaces
	text := strings.Join(args, " ")
	
	// Create command to invoke espeak
	cmd := exec.Command("espeak", text)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	
	return cmd.Run()
}

func doBible(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("bible requires at least one argument")
	}
	
	query := strings.Join(args, "+")
	
	// Create HTTP client with timeout
	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	
	// Make the request
	resp, err := client.Get("https://bible-api.com/" + query)
	if err != nil {
		return fmt.Errorf("failed to make request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("failed to make request: %s", resp.Status)
	}
	
	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response: %v", err)
	}
	bibleResponse := BibleResponse{}
	err = json.Unmarshal(body, &bibleResponse)
	if err != nil {
		return fmt.Errorf("failed to unmarshal response: %v", err)
	}

	for _, verse := range bibleResponse.Verses {
		fmt.Printf("%s %d:%d\n", verse.BookName, verse.Chapter, verse.Verse)
		fmt.Println(verse.Text)
		fmt.Println()
	}
	
	return nil
}

var cmds = map[string]*Command{}

func registerCommand(cmd Command) {
	cmds[cmd.Name] = &cmd
	for _, alias := range cmd.Aliases {
		cmds[alias] = &cmd
	}
}
