package flago

import (
	"time"
)

type durationValue time.Duration

func newDurationValue(p *time.Duration, value time.Duration) *durationValue {
	*p = value
	return (*durationValue)(p)
}

func Duration(name string, alias rune, value time.Duration, usage string, callback Callback) *time.Duration {
	p := new(time.Duration)
	CommandLine.Var(newDurationValue(p, value), name, alias, usage, 0, callback)
	return p
}

func DurationVar(p *time.Duration, name string, alias rune, value time.Duration, usage string, callback Callback) {
	CommandLine.Var(newDurationValue(p, value), name, alias, usage, 0, callback)
}

func DurationSubFlag(name string, alias rune, value time.Duration, usage string, callback Callback) *Flag {
	return CommandLine.DurationVarSubFlag(nil, name, alias, value, usage, callback)
}

func DurationVarSubFlag(p *time.Duration, name string, alias rune, value time.Duration, usage string, callback Callback) *Flag {
	if p == nil {
		p = new(time.Duration)
	}
	return CommandLine.Var(newDurationValue(p, value), name, alias, usage, NESTED, callback)
}

func (f *FlagSet) Duration(name string, alias rune, value time.Duration, usage string, callback Callback) *time.Duration {
	p := new(time.Duration)
	f.Var(newDurationValue(p, value), name, alias, usage, 0, callback)
	return p
}

func (f *FlagSet) DurationVar(p *time.Duration, name string, alias rune, value time.Duration, usage string, callback Callback) {
	f.Var(newDurationValue(p, value), name, alias, usage, 0, callback)
}

func (f *FlagSet) DurationSubFlag(name string, alias rune, value time.Duration, usage string, callback Callback) *Flag {
	return f.DurationVarSubFlag(nil, name, alias, value, usage, callback)
}

func (f *FlagSet) DurationVarSubFlag(p *time.Duration, name string, alias rune, value time.Duration, usage string, callback Callback) *Flag {
	if p == nil {
		p = new(time.Duration)
	}
	return f.Var(newDurationValue(p, value), name, alias, usage, NESTED, callback)
}

func (d *durationValue) Set(s string) error {
	v, err := time.ParseDuration(s)
	if err != nil {
		err = errParse
	}
	*d = durationValue(v)
	return err
}

func (d *durationValue) Get() interface{} {
	return time.Duration(*d)
}

func (d *durationValue) String() string {
	return (*time.Duration)(d).String()
}
