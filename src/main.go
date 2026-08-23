package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"
)

var flags struct {
	DryRun              bool
	ExitOnError         bool
	Verbose             bool
	PrintVersionAndExit bool
}

var (
	version       = "1.0.0"
	nonzeroExit   bool
	commandReader io.Reader
)

func init() {
	log.SetOutput(os.Stderr)
	log.SetFlags(0)
	flag.BoolVar(&flags.PrintVersionAndExit, "version", false, "print version and exit")
	flag.BoolVar(&flags.Verbose, "v", false, "verbose")
	flag.BoolVar(&flags.ExitOnError, "e", false, "exit on error")
	flag.BoolVar(&flags.DryRun, "n", false, "dry-run")
}

func onExecuteError(command []string, err error) {
	nonzeroExit = true
	log.Printf("execute %v: %v", command, err)
	if flags.ExitOnError {
		log.Fatalf("exiting on error")
	}
}

func execute(config *Config, command []string) {
	if len(command) == 0 {
		return
	}
	name := command[0]
	var args []string
	if len(command) > 1 {
		args = command[1:]
	}
	if flags.DryRun {
		return
	}
	if builtin, ok := cmds[name]; ok {
		if err := builtin.Func(config, args); err != nil {
			onExecuteError(command, fmt.Errorf("builtin %q: %v", name, err))
		}
		return
	}
	cmd := exec.Command(name, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		onExecuteError(command, err)
	}
}

func executeFromReader(config *Config, r io.Reader) {
	os.Stdout.Write([]byte(GreenText + ">>> " + NormalText))
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		execute(config, strings.Fields(scanner.Text()))
		os.Stdout.Write([]byte(GreenText + ">>> " + NormalText))
	}
}

func registerCommands(config *Config) {
	registerCommand(Command{
		Name:        "time",
		Aliases:     []string{},
		Description: "Display the current time",
		Func:        doTime,
	}, config)
	registerCommand(Command{
		Name:        "date",
		Aliases:     []string{},
		Description: "Display the current date",
		Func:        doDate,
	}, config)
	registerCommand(Command{
		Name:        "datetime",
		Aliases:     []string{"dt"},
		Description: "Display the current date and time",
		Func:        doDatetime,
	}, config)
	registerCommand(Command{
		Name:        "colors",
		Aliases:     []string{"color"},
		Description: "Display colors",
		Func:        doColors,
	}, config)
	registerCommand(Command{
		Name:        "days",
		Aliases:     []string{"day", "week"},
		Description: "Display days of the week",
		Func:        doDays,
	}, config)
	registerCommand(Command{
		Name:        "months",
		Aliases:     []string{"month"},
		Description: "Display months of the year",
		Func:        doMonths,
	}, config)
	registerCommand(Command{
		Name:        "calendar",
		Aliases:     []string{"cal"},
		Description: "Display the current month as a calendar",
		Func:        doCal,
	}, config)
	registerCommand(Command{
		Name:        "news",
		Aliases:     []string{},
		Description: "Display the news",
		Func:        doNews,
	}, config)
	registerCommand(Command{
		Name:        "birthdays",
		Aliases:     []string{"birthday", "bday"},
		Description: "Display your birthday and those of your family members",
		Func:        doBirthday,
	}, config)
	registerCommand(Command{
		Name:        "alphabet",
		Aliases:     []string{"abc"},
		Description: "Display the alphabet",
		Func:        doABC,
	}, config)
	registerCommand(Command{
		Name:        "beep",
		Aliases:     []string{},
		Description: "Make a beep sound",
		Func:        doBeep,
	}, config)
	registerCommand(Command{
		Name:        "help",
		Aliases:     []string{"helpme", "cmds"},
		Description: "Display all commands, aliases, and descriptions",
		Func:        doHelp,
	}, config)
	registerCommand(Command{
		Name:        "exit",
		Aliases:     []string{"quit"},
		Description: "Quit the Shell",
		Func:        doExit,
	}, config)
	registerCommand(Command{
		Name:        "numbers",
		Aliases:     []string{"nums", "num"},
		Description: "Display Numbers",
		Func:        doNum,
	}, config)
	registerCommand(Command{
		Name:        "compare",
		Aliases:     []string{"cmp"},
		Description: "Compare two or more numbers",
		Func:        doCompare,
	}, config)
	registerCommand(Command{
		Name:        "count",
		Aliases:     []string{"cnt"},
		Description: "Count up to a number",
		Func:        doCount,
	}, config)
	registerCommand(Command{
		Name:        "sort",
		Aliases:     []string{},
		Description: "Sort words or numbers",
		Func:        doSort,
	}, config)
	registerCommand(Command{
		Name:        "unique",
		Aliases:     []string{"uniq", "distinct"},
		Description: "Remove duplicates from a list so that they are all unique / distinct",
		Func:        doUniq,
	}, config)
	registerCommand(Command{
		Name:        "pwd",
		Aliases:     []string{"cwd"},
		Description: "Print the current working directory",
		Func:        doPwd,
	}, config)
	registerCommand(Command{
		Name:        "cd",
		Aliases:     []string{},
		Description: "Change the current working directory",
		Func:        doCd,
	}, config)
	registerCommand(Command{
		Name:        "list",
		Aliases:     []string{"ls"},
		Description: "List the files and folders in the current working directory",
		Func:        doLs,
	}, config)
	registerCommand(Command{
		Name:        "first",
		Aliases:     []string{},
		Description: "Print the first item in a list",
		Func:        doFirst,
	}, config)
	registerCommand(Command{
		Name:        "last",
		Aliases:     []string{},
		Description: "Print the last item in a list",
		Func:        doLast,
	}, config)
	registerCommand(Command{
		Name:        "reverse",
		Aliases:     []string{"rev"},
		Description: "Print the arguments in reverse order",
		Func:        doRev,
	}, config)
	registerCommand(Command{
		Name:        "add",
		Aliases:     []string{"sum", "total"},
		Description: "Print the sum of all arguments added together",
		Func:        doAdd,
	}, config)
	registerCommand(Command{
		Name:        "multiply",
		Aliases:     []string{"mult", "mul"},
		Description: "Print the product of all arguments multiplied together",
		Func:        doMultiply,
	}, config)
	registerCommand(Command{
		Name:        "weather",
		Aliases:     []string{"wtr"},
		Description: "Print the weather",
		Func:        doWeather,
	}, config)
	registerCommand(Command{
		Name:        "lowercase",
		Aliases:     []string{"lower"},
		Description: "Lowercase the arguments",
		Func:        doLower,
	}, config)
	registerCommand(Command{
		Name:        "uppercase",
		Aliases:     []string{"upper"},
		Description: "Uppercase the arguments",
		Func:        doUpper,
	}, config)
	registerCommand(Command{
		Name:        "environment",
		Aliases:     []string{"env"},
		Description: "Print the environment variables",
		Func:        doEnv,
	}, config)
	registerCommand(Command{
		Name:        "shuffle",
		Aliases:     []string{"shuf"},
		Description: "Randomly re-arrange the arguments",
		Func:        doShuffle,
	}, config)
	registerCommand(Command{
		Name:        "random",
		Aliases:     []string{"rand"},
		Description: "Print a random number",
		Func:        doRandom,
	}, config)
	registerCommand(Command{
		Name:        "cointoss",
		Aliases:     []string{"coin", "flip", "coinflip"},
		Description: "Flip a coin",
		Func:        doFlip,
	}, config)
	registerCommand(Command{
		Name:        "sleep",
		Aliases:     []string{"wait"},
		Description: "Pause for some amount of time",
		Func:        doSleep,
	}, config)
	registerCommand(Command{
		Name:        "compass",
		Aliases:     []string{},
		Description: "Print a compass",
		Func:        doCompass,
	}, config)
	registerCommand(Command{
		Name:        "reset",
		Aliases:     []string{},
		Description: "Reset the terminal",
		Func:        doReset,
	}, config)
	registerCommand(Command{
		Name:        "ipaddresses",
		Aliases:     []string{"ipaddress", "ip"},
		Description: "Display my IP address",
		Func:        doIp,
	}, config)
	registerCommand(Command{
		Name:        "seasons",
		Aliases:     []string{"season"},
		Description: "Display the seasons of the year",
		Func:        doSeasons,
	}, config)
	registerCommand(Command{
		Name:        "uptime",
		Aliases:     []string{},
		Description: "Display the uptime of the system",
		Func:        doUptime,
	}, config)
	registerCommand(Command{
		Name:        "push",
		Aliases:     []string{},
		Description: "Push a string to a stack",
		Func:        doPush,
	}, config)
	registerCommand(Command{
		Name:        "pop",
		Aliases:     []string{},
		Description: "Pop a string from the stack",
		Func:        doPop,
	}, config)
	registerCommand(Command{
		Name:        "stack",
		Aliases:     []string{},
		Description: "Display the contents of the stack",
		Func:        doPrintStack,
	}, config)
	registerCommand(Command{
		Name:        "queue",
		Aliases:     []string{},
		Description: "Display the contents of the queue",
		Func:        doPrintQueue,
	}, config)
	registerCommand(Command{
		Name:        "enqueue",
		Aliases:     []string{},
		Description: "Add something to the queue",
		Func:        doEnqueue,
	}, config)
	registerCommand(Command{
		Name:        "dequeue",
		Aliases:     []string{},
		Description: "Remove the next item from the queue",
		Func:        doDequeue,
	}, config)
	registerCommand(Command{
		Name:        "todo",
		Aliases:     []string{},
		Description: "Display the todo list or add something to it",
		Func:        doTodo,
	}, config)
	registerCommand(Command{
		Name:        "done",
		Aliases:     []string{},
		Description: "Mark a todo item as done either by name or index",
		Func:        doDone,
	}, config)
	registerCommand(Command{
		Name:        "home",
		Aliases:     []string{},
		Description: "Display my home address",
		Func:        doHomeAddress,
	}, config)
	registerCommand(Command{
		Name:        "birthday",
		Aliases:     []string{"bday"},
		Description: "Display my birthday",
		Func:        doBirthday,
	}, config)
	registerCommand(Command{
		Name:        "age",
		Aliases:     []string{},
		Description: "Display my age",
		Func:        doAge,
	}, config)
	registerCommand(Command{
		Name:        "countdown",
		Aliases:     []string{"tminus"},
		Description: "Display a countdown",
		Func:        doCountdown,
	}, config)
	registerCommand(Command{
		Name:        "nock",
		Aliases:     []string{},
		Description: "Evaluate a Nock expression (prints 0 on error)",
		Func:        doNock,
	}, config)
	registerCommand(Command{
		Name:        "repeat",
		Aliases:     []string{},
		Description: "Repeat the line the specified number of times",
		Func:        doRepeat,
	}, config)
	registerCommand(Command{
		Name:        "subtract",
		Aliases:     []string{"sub"},
		Description: "Subtract one number from another",
		Func:        doSubtract,
	}, config)
	registerCommand(Command{
		Name:        "countgame",
		Aliases:     []string{},
		Description: "Guess the number of Os",
		Func:        doCountGame,
	}, config)
	registerCommand(Command{
		Name:        "news",
		Aliases:     []string{},
		Description: "Show the news",
		Func:        doNews,
	}, config)
	registerCommand(Command{
		Name:        "read",
		Aliases:     []string{"cat"},
		Description: "Read a file",
		Func:        doCat,
	}, config)
	registerCommand(Command{
		Name:        "and",
		Aliases:     []string{},
		Description: "Logical AND",
		Func:        doAnd,
	}, config)
	registerCommand(Command{
		Name:        "or",
		Aliases:     []string{},
		Description: "Logical OR",
		Func:        doOr,
	}, config)
	registerCommand(Command{
		Name:        "xor",
		Aliases:     []string{},
		Description: "Logical XOR",
		Func:        doXor,
	}, config)
	registerCommand(Command{
		Name:        "not",
		Aliases:     []string{},
		Description: "Logical NOT",
		Func:        doNot,
	}, config)
	registerCommand(Command{
		Name:        "family",
		Aliases:     []string{"fam"},
		Description: "Display information about your family",
		Func:        doFamily,
	}, config)
	registerCommand(Command{
		Name:        "bedtime",
		Aliases:     []string{"bed"},
		Description: "Display the bedtime",
		Func:        doBedtime,
	}, config)
	registerCommand(Command{
		Name:        "printout",
		Aliases:     []string{"printer"},
		Description: "Print out a string to the printer",
		Func:        doPrintOut,
	}, config)
	registerCommand(Command{
		Name:        "speak",
		Aliases:     []string{},
		Description: "Speak a string",
		Func:        doSpeak,
	}, config)
	registerCommand(Command{
		Name:        "bible",
		Aliases:     []string{},
		Description: "Display a Bible verse",
		Func:        doBible,
	}, config)
}

func main() {
	flag.Parse()
	if flags.PrintVersionAndExit {
		fmt.Println(version)
		os.Exit(0)
	}
	commandReader = os.Stdin
	config := getConfig()
	registerCommands(config)
	switch {
	case flag.NArg() > 0:
		path := flag.Arg(0)
		f, err := os.Open(path)
		if err != nil {
			log.Fatalf("open %q: %v", path, err)
		}
		commandReader = f
	}
	executeFromReader(config, commandReader)
	if nonzeroExit {
		os.Exit(1)
	}
}
