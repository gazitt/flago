package flago

import (
	"strconv"
)

type intValue int

func newIntValue(p *int, value int) *intValue {
	*p = value
	return (*intValue)(p)
}

func Int(name string, alias rune, value int, usage string, callback Callback) *int {
	p := new(int)
	CommandLine.Var(newIntValue(p, value), name, alias, usage, 0, callback)
	return p
}

func IntVar(p *int, name string, alias rune, value int, usage string, callback Callback) {
	CommandLine.Var(newIntValue(p, value), name, alias, usage, 0, callback)
}

func IntSubFlag(name string, alias rune, value int, usage string, callback Callback) *Flag {
	return CommandLine.IntVarSubFlag(nil, name, alias, value, usage, callback)
}

func IntVarSubFlag(p *int, name string, alias rune, value int, usage string, callback Callback) *Flag {
	if p == nil {
		p = new(int)
	}
	return CommandLine.Var(newIntValue(p, value), name, alias, usage, NESTED, callback)
}

func (f *FlagSet) Int(name string, alias rune, value int, usage string, callback Callback) *int {
	p := new(int)
	f.Var(newIntValue(p, value), name, alias, usage, 0, callback)
	return p
}

func (f *FlagSet) IntVar(p *int, name string, alias rune, value int, usage string, callback Callback) {
	f.Var(newIntValue(p, value), name, alias, usage, 0, callback)
}

func (f *FlagSet) IntSubFlag(name string, alias rune, value int, usage string, callback Callback) *Flag {
	return f.IntVarSubFlag(nil, name, alias, value, usage, callback)
}

func (f *FlagSet) IntVarSubFlag(p *int, name string, alias rune, value int, usage string, callback Callback) *Flag {
	if p == nil {
		p = new(int)
	}
	return f.Var(newIntValue(p, value), name, alias, usage, NESTED, callback)
}

func (i *intValue) Set(s string) error {
	v, err := strconv.ParseInt(s, 0, strconv.IntSize)
	if err != nil {
		err = numError(err)
	}
	*i = intValue(v)
	return err
}

func (i *intValue) String() string {
	return strconv.Itoa(int(*i))
}

func (i *intValue) Get() interface{} {
	return int(*i)
}
