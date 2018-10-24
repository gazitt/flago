package flago

import (
	"strconv"
)

const (
	COMMAND uint = 1 << iota
	NESTED
)

type boolValue bool

func newBoolValue(p *bool, value bool) *boolValue {
	*p = value
	return (*boolValue)(p)
}

func Bool(name string, alias rune, value bool, usage string, callback Callback) *bool {
	p := new(bool)
	CommandLine.Var(newBoolValue(p, value), name, alias, usage, 0, callback)
	return p
}

func BoolVar(p *bool, name string, alias rune, value bool, usage string, callback Callback) {
	CommandLine.Var(newBoolValue(p, value), name, alias, usage, 0, callback)
}

// BoolVarSubFlag
func BoolSubFlag(name string, alias rune, value bool, usage string, callback Callback) *Flag {
	return CommandLine.BoolVarSubFlag(nil, name, alias, value, usage, callback)
}

// BoolVarSubFlag
func BoolVarSubFlag(p *bool, name string, alias rune, value bool, usage string, callback Callback) *Flag {
	if p == nil {
		p = new(bool)
	}
	return CommandLine.Var(newBoolValue(p, value), name, alias, usage, NESTED, callback)
}

func BoolSubCommand(name string, alias rune, usage string, subflags ...*Flag) *bool {
	p := new(bool)
	CommandLine.Var(newBoolValue(p, false), name, alias, usage, COMMAND, nil, subflags...)
	return p
}

func BoolVarSubCommand(p *bool, name string, alias rune, usage string, subflags ...*Flag) {
	CommandLine.Var(newBoolValue(p, false), name, alias, usage, COMMAND, nil, subflags...)
}

func BoolVarSubCommandNested(p *bool, name string, alias rune, usage string, subflags ...*Flag) *Flag {
	return CommandLine.Var(newBoolValue(p, false), name, alias, usage, COMMAND|NESTED, nil, subflags...)
}

func (f *FlagSet) Bool(name string, alias rune, value bool, usage string, callback Callback) *bool {
	p := new(bool)
	f.Var(newBoolValue(p, value), name, alias, usage, 0, callback)
	return p
}

func (f *FlagSet) BoolVar(p *bool, name string, alias rune, value bool, usage string, callback Callback) {
	f.Var(newBoolValue(p, value), name, alias, usage, 0, callback)
}

func (f *FlagSet) BoolSubFlag(name string, alias rune, value bool, usage string, callback Callback) *Flag {
	return f.BoolVarSubFlag(nil, name, alias, value, usage, callback)
}

func (f *FlagSet) BoolVarSubFlag(p *bool, name string, alias rune, value bool, usage string, callback Callback) *Flag {
	if p == nil {
		p = new(bool)
	}
	return f.Var(newBoolValue(p, value), name, alias, usage, NESTED, callback)
}

func (f *FlagSet) BoolSubCommand(name string, alias rune, usage string, subflags ...*Flag) *bool {
	p := new(bool)
	f.Var(newBoolValue(p, false), name, alias, usage, COMMAND, nil, subflags...)
	return p
}

func (f *FlagSet) BoolVarSubCommand(p *bool, name string, alias rune, usage string, subflags ...*Flag) {
	f.Var(newBoolValue(p, false), name, alias, usage, COMMAND, nil, subflags...)
}

func (f *FlagSet) BoolVarSubCommandNested(p *bool, name string, alias rune, usage string, subflags ...*Flag) *Flag {
	return f.Var(newBoolValue(p, false), name, alias, usage, COMMAND|NESTED, nil, subflags...)
}

func (b *boolValue) Set(s string) error {
	v, err := strconv.ParseBool(s)
	if err != nil {
		err = errParse
	}
	*b = boolValue(v)
	return err
}

func (b *boolValue) String() string {
	return strconv.FormatBool(bool(*b))
}

func (b *boolValue) Get() interface{} {
	return bool(*b)
}

func (b *boolValue) IsBoolFlag() bool {
	return true
}

type boolFlag interface {
	Value
	IsBoolFlag() bool
}
