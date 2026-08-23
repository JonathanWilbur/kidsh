# Kid Shell

This is a lightweight / low-capability secure shell for young children to fool
around on a computer. It has simple commands for doing easy tasks in a simple
way that children can understand, such as pressing `c` to display colors, or
`d` to display days of the week, or `n` to take notes. It only supports
built-ins, and (not implemented yet) it will be able to ran with dropped
capabilities, so your little one cannot accidentally do anything to mess up
your computer.

## Snapshot of Commands Supported

```
NAME                ALIASES             DESCRIPTION
====                =======             ===========
add                 sum,total           Print the sum of all arguments added together
age                                     Display my age
alphabet            abc                 Display the alphabet
and                                     Logical AND
bedtime             bed                 Display the bedtime
beep                                    Make a beep sound
bible                                   Display a Bible verse
birthday            bday                Display my birthday
birthdays           birthday,bday       Display your birthday and those of your family members
calendar            cal                 Display the current month as a calendar
cd                                      Change the current working directory
cointoss            coin,flip,coinflip    Flip a coin
colors              color               Display colors
compare             cmp                 Compare two or more numbers
compass                                 Print a compass
count               cnt                 Count up to a number
countdown           tminus              Display a countdown
countgame                               Guess the number of Os
date                                    Display the current date
datetime            dt                  Display the current date and time
days                day,week            Display days of the week
dequeue                                 Remove the next item from the queue
done                                    Mark a todo item as done either by name or index
enqueue                                 Add something to the queue
environment         env                 Print the environment variables
exit                quit                Quit the Shell
family              fam                 Display information about your family
first                                   Print the first item in a list
help                helpme,cmds         Display all commands, aliases, and descriptions
home                                    Display my home address
ipaddresses         ipaddress,ip        Display my IP address
last                                    Print the last item in a list
list                ls                  List the files and folders in the current working directory
lowercase           lower               Lowercase the arguments
months              month               Display months of the year
multiply            mult,mul            Print the product of all arguments multiplied together
news                                    Show the news
nock                                    Evaluate a Nock expression (prints 0 on error)
not                                     Logical NOT
numbers             nums,num            Display Numbers
or                                      Logical OR
pop                                     Pop a string from the stack
printout            printer             Print out a string to the printer
push                                    Push a string to a stack
pwd                 cwd                 Print the current working directory
queue                                   Display the contents of the queue
random              rand                Print a random number
read                cat                 Read a file
repeat                                  Repeat the line the specified number of times
reset                                   Reset the terminal
reverse             rev                 Print the arguments in reverse order
seasons             season              Display the seasons of the year
shuffle             shuf                Randomly re-arrange the arguments
sleep               wait                Pause for some amount of time
sort                                    Sort words or numbers
speak                                   Speak a string
stack                                   Display the contents of the stack
subtract            sub                 Subtract one number from another
time                                    Display the current time
todo                                    Display the todo list or add something to it
unique              uniq,distinct       Remove duplicates from a list so that they are all unique / distinct
uppercase           upper               Uppercase the arguments
uptime                                  Display the uptime of the system
weather             wtr                 Print the weather
xor                                     Logical XOR
```

See the [man page](./man/kidsh.1) for further documentation.

## Build

```bash
go build -o kidsh ./src
```

## Install

Publishing a GitHub Release (from a tag such as `v1.0.0`) builds packages and
attaches them to that release as assets. Linux packages are produced for
amd64, ARM64, and RISC-V. macOS archives cover Apple Silicon, Intel, and a
universal binary. Snap and Flatpak are built for amd64 and ARM64. The Homebrew
formula compiles for whatever Mac or Linux machine you install on.

Before tagging a release, update `debian/changelog` and the `%changelog`
section in `kidsh.spec` by hand. The Debian package version is taken from
`debian/changelog`, so that entry must match the git tag.

```bash
# Debian / Ubuntu (also _arm64.deb, _riscv64.deb)
sudo apt install ./kidsh_*_amd64.deb

# Fedora / RHEL (also .aarch64.rpm, .riscv64.rpm)
sudo rpm -i kidsh-*.x86_64.rpm

# Alpine (install the matching *.rsa.pub from the release first)
sudo cp kidsh@wilbur.space.rsa.pub /etc/apk/keys/
sudo apk add ./kidsh-*.apk

# Arch Linux (also aarch64, riscv64)
sudo pacman -U ./kidsh-*-x86_64.pkg.tar.zst

# Snap
sudo snap install --dangerous ./kidsh_*.snap

# Flatpak (also kidsh-aarch64.flatpak)
flatpak install --user ./kidsh-x86_64.flatpak

# macOS tarball (Apple Silicon, Intel, or universal)
tar -xf kidsh-*-darwin-universal.tar.gz

# Homebrew (Linux or macOS)
brew install --formula ./kidsh.rb
```

Alpine packages are signed. Generate a key once with
`packaging/scripts/gen-alpine-key.sh`, store the private key as the
`ALPINE_ABUILD_PRIVKEY` repository secret, and commit the public
`*.rsa.pub` file it prints.

Optional configuration is read from `$KIDSH_CONFIG`, then
`~/.config/kidsh/config.json`, then `/etc/kidsh.json`. An example file ships as
`kidsh.json.example` in the package docs.

I don't know if I will ever publish these. I expect this to be a low-effort,
low maintenance command line tool.

## Usage

I expect users to full-screen the window where this shell is running so their
kids cannot easily escape it. Even better would be to define a boot menu
configuration where Linux defines `kidsh` as PID 0, so they cannot escape it.
Alternatively, you could just open up this shell in a real terminal that is not
running a GUI.

## AI Usage Statement

Basically, most of the code in this project was written by AI, except for the
basic core of the shell.
