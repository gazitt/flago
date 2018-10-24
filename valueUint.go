package flago

import (
	"strconv"
)

type uintValue uint

func newUintValue(p *uint, value uint) *uintValue {
	*p = value
	return (*uintValue)(p)
}

func Uint(name string, alias rune, value uint, usage string, callback Callback) *uint {
	p := new(uint)
	CommandLine.Var(newUintValue(p, value), name, alias, usage, 0, callback)
	return p
}

func UintVar(p *uint, name string, alias rune, value uint, usage string, callback Callback) {
	CommandLine.Var(newUintValue(p, value), name, alias, usage, 0, callback)
}

func UintSubFlag(name string, alias rune, value uint, usage string, callback Callback) *Flag {
	return CommandLine.UintVarSubFlag(nil, name, alias, value, usage, callback)
}

func UintVarSubFlag(p *uint, name string, alias rune, value uint, usage string, callback Callback) *Flag {
	if p == nil {
		p = new(uint)
	}
	return CommandLine.Var(newUintValue(p, value), name, alias, usage, NESTED, callback)
}

func (f *FlagSet) Uint(name string, alias rune, value uint, usage string, callback Callback) *uint {
	p := new(uint)
	f.Var(newUintValue(p, value), name, alias, usage, 0, callback)
	return p
}

func (f *FlagSet) UintVar(p *uint, name string, alias rune, value uint, usage string, callback Callback) {
	f.Var(newUintValue(p, value), name, alias, usage, 0, callback)
}

func (f *FlagSet) UintSubFlag(name string, alias rune, value uint, usage string, callback Callback) *Flag {
	return f.UintVarSubFlag(nil, name, alias, value, usage, callback)
}

func (f *FlagSet) UintVarSubFlag(p *uint, name string, alias rune, value uint, usage string, callback Callback) *Flag {
	if p == nil {
		p = new(uint)
	}
	return f.Var(newUintValue(p, value), name, alias, usage, NESTED, callback)
}

func (i *uintValue) Set(s string) error {
	v, err := strconv.ParseUint(s, 0, strconv.IntSize)
	if err != nil {
		err = numError(err)
	}
	*i = uintValue(v)
	return err
}

func (i *uintValue) Get() interface{} {
	return uint(*i)
}

func (i *uintValue) String() string {
	return strconv.FormatUint(uint64(*i), 10)
}
