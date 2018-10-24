package flago

import (
	"fmt"
	"reflect"
	"sort"
	"strings"
)

var (
	Indent = 2
	Pad    = 20
)

func (f *FlagSet) defaultUsage() {
	fmt.Fprintf(f.Output(), "\nUsage: %s\n\n", f.name)
	f.PrintDefaults()
}

func (f *FlagSet) usage() {
	if f.Usage == nil {
		f.defaultUsage()
	} else {
		f.Usage()
	}
}

// sortFlags returns the flags as a slice in lexicographical sorted order.
func sortFlags(flags map[string]*Flag) []*Flag {
	list := make(sort.StringSlice, 0, len(flags))
	i := 0
	for k, v := range flags {
		// Excluding a mapping by the short name of the flag, because of duplication
		if strings.HasPrefix(k, _ALIAS_PREFIX) {
			continue
		}
		list = append(list, v.Name)
		i++
	}
	list.Sort()
	result := make([]*Flag, len(list))
	for i, v := range list {
		result[i] = flags[v]
	}
	return result
}

func visitAll(depth int, flags map[string]*Flag, fn func(int, *Flag)) {
	for _, v := range sortFlags(flags) {
		fn(depth, v)
	}
}

func (f *FlagSet) VisitAll(fn func(int, *Flag)) {
	visitAll(1, f.flags, fn)
}

func VisitAll(fn func(int, *Flag)) {
	CommandLine.VisitAll(fn)
}

// Lookup returns the Flag structure of the named flag, returning nil if none exists.
func (f *FlagSet) Lookup(name string) *Flag {
	return f.flags[name]
}

// Lookup returns the Flag structure of the named command-line flag,
// returning nil if none exists.
func Lookup(name string) *Flag {
	return CommandLine.flags[name]
}

// UnquoteUsage extracts a back-quoted name from the usage
// string for a flag and returns it and the un-quoted usage.
// Given "a `name` to show" it returns ("name", "a name to show").
// If there are no back quotes, the name is an educated guess of the
// type of the flag's value, or the empty string if the flag is boolean.
func UnquoteUsage(flag *Flag) (name string, usage string) {
	// Look for a back-quoted name, but avoid the strings package.
	usage = flag.Usage
	for i := 0; i < len(usage); i++ {
		if usage[i] == '`' {
			for j := i + 1; j < len(usage); j++ {
				if usage[j] == '`' {
					name = usage[i+1 : j]
					usage = usage[:i] + name + usage[j+1:]
					return name, usage
				}
			}
			break // Only one back quote; use type name.
		}
	}
	// No explicit name, so use type if we can find one.
	name = "value"
	switch flag.Value.(type) {
	case boolFlag:
		name = ""
	case *durationValue:
		name = "duration"
	case *float64Value:
		name = "float"
	case *intValue, *int64Value:
		name = "int"
	case *stringValue:
		name = "string"
	case *uintValue, *uint64Value:
		name = "uint"
	}
	return
}

// ValueType
func ValueType(f *Flag) string {
	switch f.Value.(type) {
	case *boolValue:
		return "bool"
	case *stringValue:
		return "string"
	case *intValue, *int64Value:
		return "int"
	case *durationValue:
		return "duration"
	case *float64Value:
		return "float"
	case *uintValue, *uint64Value:
		return "uint"
	default:
		return "value"
	}
}

// IsZeroValue determines whether the string represents the zero
// value for a flag.
func IsZeroValue(flag *Flag, value string) bool {
	// Build a zero value of the flag's Value type, and see if the
	// result of calling its String method equals the value passed in.
	// This works unless the Value type is itself an interface type.
	typ := reflect.TypeOf(flag.Value)
	var z reflect.Value
	if typ.Kind() == reflect.Ptr {
		z = reflect.New(typ.Elem())
	} else {
		z = reflect.Zero(typ)
	}
	return value == z.Interface().(Value).String()
}

// GetFlagName
func (f *Flag) GetFlagName() string {
	if f.Alias > 0 {
		if f.IsSubCommand() {
			return fmt.Sprintf(" %c,   %s  ", f.Alias, f.Name)
		}
		return fmt.Sprintf("-%c, --%s", f.Alias, f.Name)
	}

	if f.IsSubCommand() {
		return fmt.Sprintf("%s      ", f.Name)
	}
	return fmt.Sprintf("    --%s", f.Name)
}

// PrintDefaults default value and value type is not include in the output string
// if you need them you can define custom functions
func (f *FlagSet) PrintDefaults() {
	var options, command string
	indent, pad := Indent, Pad
	f.VisitAll(func(depth int, flag *Flag) {
		name := flag.GetFlagName()

		n := depth*indent + 2
		if len(name) > pad {
			n += len(name)
		} else {
			n += pad
		}

		s := fmt.Sprintf("%s%-20s", strings.Repeat(" ", depth*indent), name)
		for i, u := range strings.Split(flag.Usage, "\n") {
			if i > 0 {
				s += strings.Repeat(" ", n)
			} else {
				s += "  "
			}
			s += u + "\n"
		}

		if flag.IsSubCommand() {
			command += s
		} else {
			options += s
		}
	})

	s := ""
	if len(options) > 0 {
		s += fmt.Sprintf("Options:\n%s", options)
	}
	if len(command) > 0 {
		s += fmt.Sprintf("\nCommand:\n%s", command)
	}

	fmt.Fprintln(f.Output(), s)
}
