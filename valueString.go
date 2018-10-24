package flago

type stringValue string

func newStringValue(p *string, value string) *stringValue {
	*p = value
	return (*stringValue)(p)
}

func String(name string, alias rune, value string, usage string, callback Callback) *string {
	p := new(string)
	CommandLine.Var(newStringValue(p, value), name, alias, usage, 0, callback)
	return p
}

func StringVar(p *string, name string, alias rune, value string, usage string, callback Callback) {
	CommandLine.Var(newStringValue(p, value), name, alias, usage, 0, callback)
}

func StringSubFlag(name string, alias rune, value string, usage string, callback Callback) *Flag {
	return CommandLine.StringVarSubFlag(nil, name, alias, value, usage, callback)
}

func StringVarSubFlag(p *string, name string, alias rune, value string, usage string, callback Callback) *Flag {
	if p == nil {
		p = new(string)
	}
	return CommandLine.Var(newStringValue(p, value), name, alias, usage, NESTED, callback)
}

func (f *FlagSet) String(name string, alias rune, value string, usage string, callback Callback) *string {
	p := new(string)
	f.Var(newStringValue(p, value), name, alias, usage, 0, callback)
	return p
}

func (f *FlagSet) StringVar(p *string, name string, alias rune, value string, usage string, callback Callback) {
	f.Var(newStringValue(p, value), name, alias, usage, 0, callback)
}

func (f *FlagSet) StringSubFlag(name string, alias rune, value string, usage string, callback Callback) *Flag {
	return f.StringVarSubFlag(nil, name, alias, value, usage, callback)
}

func (f *FlagSet) StringVarSubFlag(p *string, name string, alias rune, value string, usage string, callback Callback) *Flag {
	if p == nil {
		p = new(string)
	}
	return f.Var(newStringValue(p, value), name, alias, usage, NESTED, callback)
}

func (s *stringValue) Set(v string) error {
	*s = stringValue(v)
	return nil
}

func (s *stringValue) String() string {
	return string(*s)
}

func (s *stringValue) Get() interface{} {
	return string(*s)
}
