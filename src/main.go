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
