package flago

import (
	"fmt"
	"os"
	"strconv"
)

func (f *FlagSet) cut() string {
	v := f.args[f.index]
	f.args = append(f.args[:f.index], f.args[f.index+1:]...)
	return v
}

// failf prints to standard error a formatted error and usage message and
// returns the error.
func (f *FlagSet) failf(format string, a ...interface{}) error {
	err := fmt.Errorf(format, a...)
	fmt.Fprintln(f.Output(), err)
	return err
}

// setValue sets value, and execute if callback is not nil
func (f *FlagSet) setValue(flag *Flag, value string, hasValue bool) error {
	var err error
	// boolean value is inverted unless a value is explicitly specified with "="
	if b, ok := flag.Value.(boolFlag); ok && b.IsBoolFlag() {
		if hasValue {
			err = flag.Value.Set(value)
		} else if v, ok := flag.Value.Get().(bool); ok {
			err = flag.Value.Set(strconv.FormatBool(!v))
		} else {
			err = fmt.Errorf("option `--%s' type not a boolean", flag.Name)
		}

	} else if hasValue {
		err = flag.Value.Set(value)

	} else if f.index < len(f.args) {
		err = flag.Value.Set(f.cut())

	} else {
		err = fmt.Errorf("option `--%s' requires an argument", flag.Name)
	}

	// callback function
	if err == nil && flag.callback != nil {
		err = flag.callback(flag.Value)
	}

	if err != nil {
		return f.failf("%s", err)
	}
	return nil
}

// parseOne parses one flag.
func (f *FlagSet) parseOne() error {
	v := f.args[f.index]

	if len(v) > 1 && v[0] == '-' {
		f.cut()

		n := 1
		if v[1] == '-' {
			n++
		}

		// get string after the "-"
		name := v[n:]

		if len(name) == 0 || name[0] == '-' || name[0] == '=' {
			return f.failf("invalid syntax as an option `%s'", v)
		}

		// when value is specified by "="
		value := ""
		hasValue := false
		for j := 0; j < len(name); j++ {
			if name[j] == '=' {
				value = name[j+1:]
				hasValue = true
				name = name[:j]
				break
			}
		}

		switch n {
		case 2: // long option
			if flag, ok := f.flags[name]; ok && !flag.IsSubCommand() {
				return f.setValue(flag, value, hasValue)
			}
			if name == "help" {
				f.usage()
				return ErrHelp
			}
			return f.failf("unrecognized option `--%s'", name)

		case 1: // short option
			if hasValue {
				if len(name) == 1 {
					if flag, ok := f.flags[aliasToKey(rune(name[0]))]; ok && !flag.IsSubCommand() {
						return f.setValue(flag, value, hasValue)
					}
					return f.failf("unrecognized option `-%s'", name)
				}
				return f.failf("unrecognized option `-%s'", name)
			}

			for _, r := range name {
				if flag, ok := f.flags[aliasToKey(r)]; ok && !flag.IsSubCommand() {
					if err := f.setValue(flag, "", false); err != nil {
						return err
					}
					continue
				}

				// to output help message
				if r == 'h' {
					f.usage()
					return ErrHelp
				}

				return f.failf("unrecognized option `%c'", r)
			}

		} // end switch

	} else {
		f.index++
	}

	return nil
}

// Parse parses flag definitions from the argument list
func (f *FlagSet) Parse(arguments []string) error {
	f.parsed = true
	f.index = 0

	var index int
	// when a sub-command continues to be at the beginning of arguments,
	// sets a flags to the f.flags again
	for _, v := range arguments {
		flag, ok := f.flags[v]

		// long name also may be one letter
		// so, find as a short name in case of not OK
		if !ok && len(v) == 1 {
			flag, ok = f.flags[aliasToKey(rune(v[0]))]
		}

		if ok && flag.IsSubCommand() {
			f.addSubCommandName(flag.Name)
			if err := flag.Value.Set("true"); err != nil {
				switch f.errorHandling {
				case ContinueOnError:
					return err
				case ExitOnError:
					os.Exit(2)
				case PanicOnError:
					panic(err)
				}
			}

			f.flags = flag.flags
			index++
			continue
		}

		break
	}

	if index > 0 {
		f.args = arguments[index:]
	} else {
		f.args = arguments
	}

	// parses!
	for f.index < len(f.args) {
		err := f.parseOne()

		if err != nil {
			switch f.errorHandling {
			case ContinueOnError:
				return err
			case ExitOnError:
				os.Exit(2)
			case PanicOnError:
				panic(err)
			}
		}
	}

	return nil
}
