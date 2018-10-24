package flago

import (
	"strconv"
)

type int64Value int64

func newInt64Value(p *int64, value int64) *int64Value {
	*p = value
	return (*int64Value)(p)
}

func Int64(name string, alias rune, value int64, usage string, callback Callback) *int64 {
	p := new(int64)
	CommandLine.Var(newInt64Value(p, value), name, alias, usage, 0, callback)
	return p
}

func Int64Var(p *int64, name string, alias rune, value int64, usage string, callback Callback) {
	CommandLine.Var(newInt64Value(p, value), name, alias, usage, 0, callback)
}

func Int64SubFlag(name string, alias rune, value int64, usage string, callback Callback) *Flag {
	return CommandLine.Int64VarSubFlag(nil, name, alias, value, usage, callback)
}

func Int64VarSubFlag(p *int64, name string, alias rune, value int64, usage string, callback Callback) *Flag {
	if p == nil {
		p = new(int64)
	}
	return CommandLine.Var(newInt64Value(p, value), name, alias, usage, NESTED, callback)
}

func (f *FlagSet) Int64(name string, alias rune, value int64, usage string, callback Callback) *int64 {
	p := new(int64)
	f.Var(newInt64Value(p, value), name, alias, usage, 0, callback)
	return p
}

func (f *FlagSet) Int64Var(p *int64, name string, alias rune, value int64, usage string, callback Callback) {
	f.Var(newInt64Value(p, value), name, alias, usage, 0, callback)
}

func (f *FlagSet) Int64SubFlag(name string, alias rune, value int64, usage string, callback Callback) *Flag {
	return f.Int64VarSubFlag(nil, name, alias, value, usage, callback)
}

func (f *FlagSet) Int64VarSubFlag(p *int64, name string, alias rune, value int64, usage string, callback Callback) *Flag {
	if p == nil {
		p = new(int64)
	}
	return f.Var(newInt64Value(p, value), name, alias, usage, NESTED, callback)
}

func (i *int64Value) Set(s string) error {
	v, err := strconv.ParseInt(s, 0, 64)
	if err != nil {
		err = numError(err)
	}
	*i = int64Value(v)
	return err
}

func (i *int64Value) Get() interface{} {
	return int64(*i)
}

func (i *int64Value) String() string {
	return strconv.FormatInt(int64(*i), 10)
}
