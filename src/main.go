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

const appName = "kidsh"

var (
	version        = "1.0.0"
	nonzeroExit    bool
	commandReader  io.Reader
	postionalArg0  string
	positionalArgs []string
)

func init() {
	log.SetOutput(os.Stderr)
	log.SetFlags(0)
	flag.BoolVar(&flags.PrintVersionAndExit, "version", false, "print version and exit")
	flag.BoolVar(&flags.Verbose, "v", false, "verbose")
	flag.BoolVar(&flags.ExitOnError, "e", false, "exit on error")
	flag.BoolVar(&flags.DryRun, "n", false, "dry-run")
	flag.Parse()
}

func onExecuteError(command []string, err error) {
	nonzeroExit = true
	log.Printf("execute %v: %v", command, err)
	if flags.ExitOnError {
		log.Fatalf("exiting on error")
	}
}

func execute(command []string) {
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
		if err := builtin.Func(args); err != nil {
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

func executeFromReader(r io.Reader) {
	os.Stdout.Write([]byte(GreenText + ">>> " + NormalText))
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		execute(strings.Fields(scanner.Text()))
		os.Stdout.Write([]byte(GreenText + ">>> " + NormalText))
	}
}

func init() {
	registerCommand(Command{
		Name:        "time",
		Aliases:     []string{},
		Description: "Display the current time",
		Func:        doTime,
	})
	registerCommand(Command{
		Name:        "date",
		Aliases:     []string{},
		Description: "Display the current date",
		Func:        doDate,
	})
	registerCommand(Command{
		Name:        "datetime",
		Aliases:     []string{"dt"},
		Description: "Display the current date and time",
		Func:        doDatetime,
	})
	registerCommand(Command{
		Name:        "colors",
		Aliases:     []string{"color"},
		Description: "Display colors",
		Func:        doColors,
	})
	registerCommand(Command{
		Name:        "days",
		Aliases:     []string{"day", "week"},
		Description: "Display days of the week",
		Func:        doDays,
	})
	registerCommand(Command{
		Name:        "months",
		Aliases:     []string{"month"},
		Description: "Display months of the year",
		Func:        doMonths,
	})
	registerCommand(Command{
		Name:        "calendar",
		Aliases:     []string{"cal"},
		Description: "Display the current month as a calendar",
		Func:        doCal,
	})
	registerCommand(Command{
		Name:        "news",
		Aliases:     []string{},
		Description: "Display the news",
		Func:        doNews,
	})
	registerCommand(Command{ // TODO: Reconsider the name of this command.
		Name:        "message",
		Aliases:     []string{"msg", "mesg", "announce"},
		Description: "Send a message",
		Func:        doMsg,
	})
	registerCommand(Command{
		Name:        "birthdays",
		Aliases:     []string{"birthday", "bday"},
		Description: "Display your birthday and those of your family members",
		Func:        doBday,
	})
	registerCommand(Command{ // TODO: Reconsider this command entirely: why not add / subtract, etc.?
		Name:        "calculator",
		Aliases:     []string{"calc"},
		Description: "Calculator",
		Func:        doCalc,
	})
	registerCommand(Command{
		Name:        "alphabet",
		Aliases:     []string{"abc"},
		Description: "Display the alphabet",
		Func:        doABC,
	})
	registerCommand(Command{
		Name:        "beep",
		Aliases:     []string{},
		Description: "Make a beep sound",
		Func:        doBeep,
	})
	registerCommand(Command{ // TODO: about calling 911
		Name:        "help",
		Aliases:     []string{"helpme", "cmds"},
		Description: "Display all commands, aliases, and descriptions",
		Func:        doHelp,
	})
	registerCommand(Command{
		Name:        "exit",
		Aliases:     []string{"quit"},
		Description: "Quit the Shell",
		Func:        doExit,
	})
	registerCommand(Command{
		Name:        "numbers",
		Aliases:     []string{"nums", "num"},
		Description: "Display Numbers",
		Func:        doNum,
	})
	registerCommand(Command{
		Name:        "compare",
		Aliases:     []string{"cmp"},
		Description: "Compare two or more numbers",
		Func:        doCompare,
	})
	registerCommand(Command{
		Name:        "count",
		Aliases:     []string{"cnt"},
		Description: "Count up to a number",
		Func:        doCount,
	})
	registerCommand(Command{
		Name:        "sort",
		Aliases:     []string{},
		Description: "Sort words or numbers",
		Func:        doSort,
	})
	registerCommand(Command{
		Name:        "unique",
		Aliases:     []string{"uniq", "distinct"},
		Description: "Remove duplicates from a list so that they are all unique / distinct",
		Func:        doUniq,
	})
	registerCommand(Command{
		Name:        "pwd",
		Aliases:     []string{"cwd"},
		Description: "Print the current working directory",
		Func:        doPwd,
	})
	registerCommand(Command{
		Name:        "cd",
		Aliases:     []string{},
		Description: "Change the current working directory",
		Func:        doCd,
	})
	registerCommand(Command{
		Name:        "list",
		Aliases:     []string{"ls"},
		Description: "List the files and folders in the current working directory",
		Func:        doLs,
	})
	registerCommand(Command{
		Name:        "first",
		Aliases:     []string{},
		Description: "Print the first item in a list",
		Func:        doFirst,
	})
	registerCommand(Command{
		Name:        "last",
		Aliases:     []string{},
		Description: "Print the last item in a list",
		Func:        doLast,
	})
	registerCommand(Command{
		Name:        "reverse",
		Aliases:     []string{"rev"},
		Description: "Print the arguments in reverse order",
		Func:        doRev,
	})
	registerCommand(Command{
		Name:        "add",
		Aliases:     []string{"sum", "total"},
		Description: "Print the sum of all arguments added together",
		Func:        doAdd,
	})
	registerCommand(Command{
		Name:        "multiply",
		Aliases:     []string{"mult", "mul"},
		Description: "Print the product of all arguments multiplied together",
		Func:        doMultiply,
	})
	registerCommand(Command{
		Name:        "weather",
		Aliases:     []string{"wtr"},
		Description: "Print the weather",
		Func:        doWeather,
	})
	registerCommand(Command{
		Name:        "lowercase",
		Aliases:     []string{"lower"},
		Description: "Lowercase the arguments",
		Func:        doLower,
	})
	registerCommand(Command{
		Name:        "uppercase",
		Aliases:     []string{"upper"},
		Description: "Uppercase the arguments",
		Func:        doUpper,
	})
	registerCommand(Command{
		Name:        "environment",
		Aliases:     []string{"env"},
		Description: "Print the environment variables",
		Func:        doEnv,
	})
	registerCommand(Command{
		Name:        "shuffle",
		Aliases:     []string{"shuf"},
		Description: "Randomly re-arrange the arguments",
		Func:        doShuffle,
	})
	registerCommand(Command{
		Name:        "random",
		Aliases:     []string{"rand"},
		Description: "Print a random number",
		Func:        doRandom,
	})
	registerCommand(Command{
		Name:        "cointoss",
		Aliases:     []string{"coin", "flip", "coinflip"},
		Description: "Flip a coin",
		Func:        doFlip,
	})
	registerCommand(Command{
		Name:        "sleep",
		Aliases:     []string{"wait"},
		Description: "Pause for some amount of time",
		Func:        doSleep,
	})
	registerCommand(Command{
		Name:        "compass",
		Aliases:     []string{},
		Description: "Print a compass",
		Func:        doCompass,
	})
	registerCommand(Command{
		Name:        "reset",
		Aliases:     []string{},
		Description: "Reset the terminal",
		Func:        doReset,
	})
	registerCommand(Command{
		Name:        "ipaddresses",
		Aliases:     []string{"ipaddress", "ip"},
		Description: "Display my IP address",
		Func:        doIp,
	})
	registerCommand(Command{
		Name:        "seasons",
		Aliases:     []string{"season"},
		Description: "Display the seasons of the year",
		Func:        doSeasons,
	})
	registerCommand(Command{
		Name:        "uptime",
		Aliases:     []string{},
		Description: "Display the uptime of the system",
		Func:        doUptime,
	})
	registerCommand(Command{
		Name:        "push",
		Aliases:     []string{},
		Description: "Push a string to a stack",
		Func:        doPush,
	})
	registerCommand(Command{
		Name:        "pop",
		Aliases:     []string{},
		Description: "Pop a string from the stack",
		Func:        doPop,
	})
	registerCommand(Command{
		Name:        "stack",
		Aliases:     []string{},
		Description: "Display the contents of the stack",
		Func:        doPrintStack,
	})
	registerCommand(Command{
		Name:        "queue",
		Aliases:     []string{},
		Description: "Display the contents of the queue",
		Func:        doPrintQueue,
	})
	registerCommand(Command{
		Name:        "enqueue",
		Aliases:     []string{},
		Description: "Add something to the queue",
		Func:        doEnqueue,
	})
	registerCommand(Command{
		Name:        "dequeue",
		Aliases:     []string{},
		Description: "Remove the next item from the queue",
		Func:        doDequeue,
	})
	registerCommand(Command{
		Name:        "todo",
		Aliases:     []string{},
		Description: "Display the todo list or add something to it",
		Func:        doTodo,
	})
	registerCommand(Command{
		Name:        "done",
		Aliases:     []string{},
		Description: "Mark a todo item as done either by name or index",
		Func:        doDone,
	})
	registerCommand(Command{
		Name:        "home",
		Aliases:     []string{},
		Description: "Display my home address",
		Func:        doHomeAddress,
	})
	registerCommand(Command{
		Name:        "birthday",
		Aliases:     []string{"bday"},
		Description: "Display my birthday",
		Func:        doBirthday,
	})
	registerCommand(Command{
		Name:        "age",
		Aliases:     []string{},
		Description: "Display my age",
		Func:        doAge,
	})
	registerCommand(Command{
		Name:        "countdown",
		Aliases:     []string{"tminus"},
		Description: "Display a countdown",
		Func:        doCountdown,
	})
	registerCommand(Command{
		Name:        "nock",
		Aliases:     []string{},
		Description: "Evaluate a Nock expression (prints 0 on error)",
		Func:        doNock,
	})
}

func main() {
	if flags.PrintVersionAndExit {
		fmt.Println(version)
		os.Exit(0)
	}
	commandReader = os.Stdin
	postionalArg0, _ = os.Executable()
	switch {
	case flag.NArg() > 0:
		path := flag.Arg(0)
		f, err := os.Open(path)
		if err != nil {
			log.Fatalf("open %q: %v", path, err)
		}
		if flag.NArg() > 1 {
			positionalArgs = flag.Args()[1:]
		}
		commandReader = f
	}
	executeFromReader(commandReader)
	if nonzeroExit {
		os.Exit(1)
	}
}
