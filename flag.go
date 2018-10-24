package flago

/*
	Package flago implements command-line flag parsing.
	Based on the built-in flag library
		flag: https://github.com/golang/go/tree/master/src/flag

	Command line flag syntax

		// long name
		--flag
		--flag value
		--flag=value

		// short name
		-f
		-f value
		-f=value

		// mixed
		-abc
		-abc value

		// subcommand. defined init sub-command
		$ command init -a -b -c

		// there are other values in before a flag, does not end
		$ command other -a other -b other -c other

*/

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
)

const (
	_ALIAS_PREFIX = "|ALIAS|"
	// _COMAND_PREFIX = "|COMMAND|"
)

// ErrHelp is the error returned if the -help or -h flag is invoked
// but no such flag is defined.
var ErrHelp = errors.New("flag: help requested")

// errParse is returned by Set if a flag's value fails to parse, such as with an invalid integer for Int.
// It then gets wrapped through failf to provide more information.
var errParse = errors.New("parse error")

// errRange is returned by Set if a flag's value is out of range.
// It then gets wrapped through failf to provide more information.
var errRange = errors.New("value out of range")

func numError(err error) error {
	ne, ok := err.(*strconv.NumError)
	if !ok {
		return err
	}
	if ne.Err == strconv.ErrSyntax {
		return errParse
	}
	if ne.Err == strconv.ErrRange {
		return errRange
	}
	return err
}

// ErrorHandling defines how FlagSet.Parse behaves if the parse fails.
type ErrorHandling int

// These constants cause FlagSet.Parse to behave as described if the parse fails.
const (
	ContinueOnError ErrorHandling = iota // Return a descriptive error.
	ExitOnError                          // Call os.Exit(2).
	PanicOnError                         // Call panic with a descriptive error.
)

// A FlagSet represents a set of defined flags. The zero value of a FlagSet
// has no name and has ContinueOnError error handling.
type FlagSet struct {
	// Usage is the function called when an error occurs while parsing flags.
	// The field is a function (not a method) that may be changed to point to
	// a custom error handler. What happens after Usage is called depends
	// on the ErrorHandling setting; for the command line, this defaults
	// to ExitOnError, which exits the program after calling Usage.
	Usage func()

	name          string
	args          []string // argument other than flag
	parsed        bool
	index         int
	flags         map[string]*Flag
	errorHandling ErrorHandling
	output        io.Writer
}

func NewFlagSet(name string, errorHandling ErrorHandling) *FlagSet {
	f := new(FlagSet)
	f.name = name
	f.errorHandling = errorHandling
	f.Usage = f.defaultUsage
	return f
}

// CommandLine is the default set of command-line flags, parsed from os.Args.
// The top-level functions such as BoolVar, Arg, and so on are wrappers for the
// methods of CommandLine.
var CommandLine = NewFlagSet(os.Args[0], ExitOnError)

// Init sets the name and error handling property for a flag set.
// By default, the zero FlagSet uses an empty name and the
// ContinueOnError error handling policy.
func (f *FlagSet) Init(name string, errorHandling ErrorHandling) {
	f.name = name
	f.errorHandling = errorHandling
}

// Name returns the name of the flag set.
func (f *FlagSet) Name() string {
	return f.name
}

// Output returns the destination for usage and error messages. os.Stderr is returned if
// output was not set or was set to nil.
func (f *FlagSet) Output() io.Writer {
	if f.output == nil {
		return os.Stderr
	}
	return f.output
}

// SetOutput sets the destination for usage and error messages.
// If output is nil, os.Stderr is used.
func (f *FlagSet) SetOutput(w io.Writer) {
	f.output = w
}

// Parse parses the command-line flags from os.Args[1:]. Must be called
// after all flags are defined and before flags are accessed by the program.
func Parse() error {
	return CommandLine.Parse(os.Args[1:])
}

// Parsed reports whether f.Parse has been called.
func (f *FlagSet) Parsed() bool { return f.parsed }

// Parsed reports whether the command-line flags have been parsed.
func Parsed() bool { return CommandLine.parsed }

// NFlag returns the number of flags that have been set.
func (f *FlagSet) NFlag() int { return len(f.flags) }

// NFlag returns the number of command-line flags that have been set.
func NFlag() int { return len(CommandLine.flags) }

// Args returns the non-flag arguments.
func (f *FlagSet) Args() []string { return f.args }

// Args returns the non-flag command-line arguments.
func Args() []string { return CommandLine.args }

// NArg is the number of arguments remaining after flags have been processed.
func (f *FlagSet) NArg() int { return len(f.args) }

// NArg is the number of arguments remaining after flags have been processed.
func NArg() int { return len(CommandLine.args) }

// Arg returns the i'th argument. Arg(0) is the first remaining argument
// after flags have been processed. Arg returns an empty string if the
// requested element does not exist.
func (f *FlagSet) Arg(i int) string {
	if i < 0 || i >= len(f.args) {
		return ""
	}
	return f.args[i]
}

// Arg returns the i'th command-line argument. Arg(0) is the first remaining argument
// after flags have been processed. Arg returns an empty string if the
// requested element does not exist.
func Arg(i int) string { return CommandLine.Arg(i) }

// PrintDefaults prints, to standard error unless configured otherwise,
// a usage message showing the default settings of all defined
// command-line flags.
func PrintDefaults() {
	CommandLine.PrintDefaults()
}

// ErrorHandling
func (f *FlagSet) ErrorHandling() ErrorHandling {
	return f.errorHandling
}

var Usage = func() {
	fmt.Fprintf(CommandLine.Output(), "\nUsage: %s\n\n", CommandLine.name)
	PrintDefaults()
}

func init() {
	// Override generic FlagSet default Usage with call to global Usage.
	// Note: This is not CommandLine.Usage = Usage,
	// because we want any eventual call to use any updated value of Usage,
	// not the value it has when this line is run.
	CommandLine.Usage = commandLineUsage
}

func commandLineUsage() {
	Usage()
}

func (f *FlagSet) addSubCommandName(name string) {
	f.name += " " + name
}

type Value interface {
	Set(string) error
	Get() interface{}
	String() string
}

// not define the Getter
// https://github.com/golang/go/blob/5ddb20912043ff7ad722a27cc93a7e68d1c5ec78/src/flag/flag.go#L296
// type Getter interface {
// 	Value
// 	Get() interface{}
// }

type Callback func(v Value) error

// A Flag represents the state of a flag.
type Flag struct {
	Name         string
	Usage        string
	Alias        rune
	Value        Value
	callback     Callback
	DefValue     string // default value (as text); for usage message
	flags        map[string]*Flag
	isSubCommand bool
}

func (f *Flag) IsSubCommand() bool {
	return f.isSubCommand
}

func isValidAlias(alias rune) bool {
	// alias is must be single alphabet letter
	if alias > 0 {
		// 48-57  0~9
		// 65-90  A~Z
		// 97-122 a~z
		// if (alias > 64 && alias < 91) || (alias > 96 && alias < 123) {
		if (alias > 47 && alias < 58) || (alias > 64 && alias < 91) || (alias > 96 && alias < 123) {
			return true
		}
		panic(fmt.Sprintf("`%c' is invalid as an alias", alias))
	}
	return false
}

func (f *FlagSet) flagAlreadyThere(name string) {
	if _, alreadythere := f.flags[name]; alreadythere {
		s := fmt.Sprintf("flag redefined: %s", name)
		fmt.Fprintln(f.Output(), s)
		panic(s)
	}
}

func (f *Flag) subflagAlreadyThere(w io.Writer, name string) {
	if _, alreadythere := f.flags[name]; alreadythere {
		s := fmt.Sprintf("flag redefined: %s", name)
		fmt.Fprintln(w, s)
		panic(s)
	}
}

func aliasToKey(alias rune) string {
	return _ALIAS_PREFIX + string(alias)
}

// func commandToKey(name string) string {
// 	return _COMAND_PREFIX + name
// }

// Var  defines a flag with the specified long short name, usage string, bit flags, callback, sub-flags.
// fifth argument:
//   0      : normal flag
//   COMMAND: sub-command
//   NESTED : nested sub-command or sub-flag.
//            so, for nested subcommands,
//            specify as follows COMMAND|NESTED
func (f *FlagSet) Var(value Value, name string, alias rune, usage string, u uint, callback Callback, subflags ...*Flag) *Flag {
	flag := &Flag{
		Name:         name,
		Usage:        usage,
		Value:        value,
		DefValue:     value.String(),
		callback:     callback,
		isSubCommand: u&COMMAND == COMMAND,
	}

	if len(subflags) > 0 {
		flag.flags = make(map[string]*Flag)

		for _, v := range subflags {
			flag.subflagAlreadyThere(f.Output(), v.Name)
			flag.flags[v.Name] = v
			if isValidAlias(v.Alias) {
				a := aliasToKey(v.Alias)
				flag.subflagAlreadyThere(f.Output(), a)
				flag.flags[a] = v
			}
		}
	}

	f.flagAlreadyThere(name)

	if isValidAlias(alias) {
		flag.Alias = alias
	}

	if u&NESTED != NESTED {
		if f.flags == nil {
			f.flags = make(map[string]*Flag)
		}
		f.flags[flag.Name] = flag

		if flag.Alias != 0 {
			a := aliasToKey(flag.Alias)
			f.flagAlreadyThere(a)
			f.flags[a] = flag
		}
	}

	return flag
}

// Var  defines a flag with the specified long short name, usage string, bit flags, callback, sub-flags.
func Var(value Value, name string, alias rune, usage string, u uint, callback Callback, subflags ...*Flag) *Flag {
	return CommandLine.Var(value, name, alias, usage, u, callback, subflags...)
}
