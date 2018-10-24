package flago

import (
	"strconv"
)

type float64Value float64

func newFloat64Value(p *float64, value float64) *float64Value {
	*p = value
	return (*float64Value)(p)
}

func Float64(name string, alias rune, value float64, usage string, callback Callback) *float64 {
	p := new(float64)
	CommandLine.Var(newFloat64Value(p, value), name, alias, usage, 0, callback)
	return p
}

func Float64Var(p *float64, name string, alias rune, value float64, usage string, callback Callback) {
	CommandLine.Var(newFloat64Value(p, value), name, alias, usage, 0, callback)
}

func Float64SubFlag(name string, alias rune, value float64, usage string, callback Callback) *Flag {
	return CommandLine.Float64VarSubFlag(nil, name, alias, value, usage, callback)
}

func Float64VarSubFlag(p *float64, name string, alias rune, value float64, usage string, callback Callback) *Flag {
	if p == nil {
		p = new(float64)
	}
	return CommandLine.Var(newFloat64Value(p, value), name, alias, usage, NESTED, callback)
}

func (f *FlagSet) Float64(name string, alias rune, value float64, usage string, callback Callback) *float64 {
	p := new(float64)
	f.Var(newFloat64Value(p, value), name, alias, usage, 0, callback)
	return p
}

func (f *FlagSet) Float64Var(p *float64, name string, alias rune, value float64, usage string, callback Callback) {
	f.Var(newFloat64Value(p, value), name, alias, usage, 0, callback)
}

func (f *FlagSet) Float64SubFlag(name string, alias rune, value float64, usage string, callback Callback) *Flag {
	return f.Float64VarSubFlag(nil, name, alias, value, usage, callback)
}

func (f *FlagSet) Float64VarSubFlag(p *float64, name string, alias rune, value float64, usage string, callback Callback) *Flag {
	if p == nil {
		p = new(float64)
	}
	return f.Var(newFloat64Value(p, value), name, alias, usage, NESTED, callback)
}

func (f *float64Value) Set(s string) error {
	v, err := strconv.ParseFloat(s, 64)
	if err != nil {
		err = numError(err)
	}
	*f = float64Value(v)
	return err
}

func (f *float64Value) Get() interface{} {
	return float64(*f)
}

func (f *float64Value) String() string {
	return strconv.FormatFloat(float64(*f), 'g', -1, 64)
}
