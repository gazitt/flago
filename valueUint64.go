package flago

import (
	"strconv"
)

type uint64Value uint64

func newUint64Value(p *uint64, value uint64) *uint64Value {
	*p = value
	return (*uint64Value)(p)
}

func Uint64(name string, alias rune, value uint64, usage string, callback Callback) *uint64 {
	p := new(uint64)
	CommandLine.Var(newUint64Value(p, value), name, alias, usage, 0, callback)
	return p
}

func Uint64Var(p *uint64, name string, alias rune, value uint64, usage string, callback Callback) {
	CommandLine.Var(newUint64Value(p, value), name, alias, usage, 0, callback)
}

func Uint64SubFlag(name string, alias rune, value uint64, usage string, callback Callback) *Flag {
	return CommandLine.Uint64VarSubFlag(nil, name, alias, value, usage, callback)
}

func Uint64VarSubFlag(p *uint64, name string, alias rune, value uint64, usage string, callback Callback) *Flag {
	if p == nil {
		p = new(uint64)
	}
	return CommandLine.Var(newUint64Value(p, value), name, alias, usage, NESTED, callback)
}

func (f *FlagSet) Uint64(name string, alias rune, value uint64, usage string, callback Callback) *uint64 {
	p := new(uint64)
	f.Var(newUint64Value(p, value), name, alias, usage, 0, callback)
	return p
}

func (f *FlagSet) Uint64Var(p *uint64, name string, alias rune, value uint64, usage string, callback Callback) {
	f.Var(newUint64Value(p, value), name, alias, usage, 0, callback)
}

func (f *FlagSet) Uint64SubFlag(name string, alias rune, value uint64, usage string, callback Callback) *Flag {
	return f.Uint64VarSubFlag(nil, name, alias, value, usage, callback)
}

func (f *FlagSet) Uint64VarSubFlag(p *uint64, name string, alias rune, value uint64, usage string, callback Callback) *Flag {
	if p == nil {
		p = new(uint64)
	}
	return f.Var(newUint64Value(p, value), name, alias, usage, NESTED, callback)
}

func (i *uint64Value) Set(s string) error {
	v, err := strconv.ParseUint(s, 0, 64)
	if err != nil {
		err = numError(err)
	}
	*i = uint64Value(v)
	return err
}

func (i *uint64Value) Get() interface{} {
	return uint64(*i)
}

func (i *uint64Value) String() string {
	return strconv.FormatUint(uint64(*i), 10)
}
