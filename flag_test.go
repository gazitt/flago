package flago

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"sort"
	"strconv"
	"strings"
	"testing"
	"time"
)

func boolString(s string) string {
	if s == "0" {
		return "false"
	}
	return "true"
}

func TestEverything(t *testing.T) {
	ResetForTesting(nil)
	Bool("test_bool", 'b', false, "bool value", nil)
	Int("test_int", 'i', 0, "int value", nil)
	Int64("test_int64", 'I', 0, "int64 value", nil)
	Uint("test_uint", 'u', 0, "uint value", nil)
	Uint64("test_uint64", 'U', 0, "uint64 value", nil)
	String("test_string", 's', "0", "string value", nil)
	Float64("test_float64", 'f', 0, "float64 value", nil)
	Duration("test_duration", 'd', 0, "time.Duration value", nil)

	m := make(map[string]*Flag)
	desired := "0"
	visitor := func(_ int, f *Flag) {
		if len(f.Name) > 5 && f.Name[0:5] == "test_" {
			m[f.Name] = f
			ok := false
			switch {
			case f.Value.String() == desired:
				ok = true
			case f.Name == "test_bool" && f.Value.String() == boolString(desired):
				ok = true
			case f.Name == "test_duration" && f.Value.String() == desired+"s":
				ok = true
			}
			if !ok {
				t.Error("Visit: bad value", f.Value.String(), "for", f.Name)
			}
		}
	}
	CommandLine.VisitAll(visitor)
	if len(m) != 8 {
		t.Error("VisitAll misses some flags")
		for k, v := range m {
			t.Log(k, *v)
		}
	}

	mm := map[string]string{
		"test_bool":     "true",
		"test_int":      "1",
		"test_int64":    "1",
		"test_uint":     "1",
		"test_uint64":   "1",
		"test_string":   "1",
		"test_float64":  "1",
		"test_duration": "1s",
	}
	// Now set all flags
	for k, v := range m {
		v.Value.Set(mm[k])
	}
	desired = "1"
	CommandLine.VisitAll(visitor)
	if len(m) != 8 {
		t.Error("Visit fails after set")
		for k, v := range m {
			t.Log(k, *v)
		}
	}

	// Now test they're visited in sort order.
	var flagNames []string
	CommandLine.VisitAll(func(_ int, f *Flag) { flagNames = append(flagNames, f.Name) })
	if !sort.StringsAreSorted(flagNames) {
		t.Errorf("flag names not sorted: %v", flagNames)
	}
}

func TestGet(t *testing.T) {
	ResetForTesting(nil)
	Bool("test_bool", 'b', true, "bool value", nil)
	Int("test_int", 'i', 1, "int value", nil)
	Int64("test_int64", 'I', 2, "int64 value", nil)
	Uint("test_uint", 'u', 3, "uint value", nil)
	Uint64("test_uint64", 'U', 4, "uint64 value", nil)
	String("test_string", 's', "5", "string value", nil)
	Float64("test_float64", 'f', 6, "float64 value", nil)
	Duration("test_duration", 'd', 7, "time.Duration value", nil)

	visitor := func(_ int, f *Flag) {
		if len(f.Name) > 5 && f.Name[0:5] == "test_" {
			var ok bool
			switch f.Name {
			case "test_bool":
				ok = f.Value.Get() == true
			case "test_int":
				ok = f.Value.Get() == int(1)
			case "test_int64":
				ok = f.Value.Get() == int64(2)
			case "test_uint":
				ok = f.Value.Get() == uint(3)
			case "test_uint64":
				ok = f.Value.Get() == uint64(4)
			case "test_string":
				ok = f.Value.Get() == "5"
			case "test_float64":
				ok = f.Value.Get() == float64(6)
			case "test_duration":
				ok = f.Value.Get() == time.Duration(7)
			}
			if !ok {
				t.Errorf("Visit: bad value %T(%v) for %s", f.Value.Get(), f.Value.Get(), f.Name)
			}
		}
	}
	CommandLine.VisitAll(visitor)
}

func TestUsage(t *testing.T) {
	for _, v := range []string{"--help", "-h"} {
		called := false
		ResetForTesting(func() { called = true })

		CommandLine.Parse([]string{v})

		if !called {
			t.Error("did not call Usage for unknown flag")
		}
	}
}

func testParse(f *FlagSet, t *testing.T) {
	if f.Parsed() {
		t.Error("f.Parse() = true before Parse")
	}
	boolFlag := f.Bool("bool", 'b', false, "bool value", nil)
	bool2Flag := f.Bool("bool2", 'B', false, "bool2 value", nil)
	intFlag := f.Int("int", 'i', 0, "int value", nil)
	int64Flag := f.Int64("int64", 'I', 0, "int64 value", nil)
	uintFlag := f.Uint("uint", 'u', 0, "uint value", nil)
	uint64Flag := f.Uint64("uint64", 'U', 0, "uint64 value", nil)
	stringFlag := f.String("string", 's', "0", "string value", nil)
	float64Flag := f.Float64("float64", 'f', 0, "float64 value", nil)
	durationFlag := f.Duration("duration", 'd', 5*time.Second, "time.Duration value", nil)
	extra := "one-extra-argument"
	args := []string{
		"--bool",
		extra,
		"--bool2=true",
		"not-flag",
		"--int", "22",
		"not-flag",
		"--int64", "0x23",
		"--uint", "24",
		"not-flag",
		"--uint64", "25",
		"--string", "hello",
		"--float64", "2718e28",
		"not-flag",
		"--duration", "2m",
		"not-flag",
	}
	if err := f.Parse(args); err != nil {
		t.Fatal(err)
	}
	if !f.Parsed() {
		t.Error("f.Parse() = false after Parse")
	}
	if *boolFlag != true {
		t.Error("bool flag should be true, is ", *boolFlag)
	}
	if *bool2Flag != true {
		t.Error("bool2 flag should be true, is ", *bool2Flag)
	}
	if *intFlag != 22 {
		t.Error("int flag should be 22, is ", *intFlag)
	}
	if *int64Flag != 0x23 {
		t.Error("int64 flag should be 0x23, is ", *int64Flag)
	}
	if *uintFlag != 24 {
		t.Error("uint flag should be 24, is ", *uintFlag)
	}
	if *uint64Flag != 25 {
		t.Error("uint64 flag should be 25, is ", *uint64Flag)
	}
	if *stringFlag != "hello" {
		t.Error("string flag should be `hello`, is ", *stringFlag)
	}
	if *float64Flag != 2718e28 {
		t.Error("float64 flag should be 2718e28, is ", *float64Flag)
	}
	if *durationFlag != 2*time.Minute {
		t.Error("duration flag should be 2m, is ", *durationFlag)
	}
	if len(f.Args()) != 6 {
		t.Error("expected one argument, got", len(f.Args()))
	} else if f.Args()[0] != extra {
		t.Errorf("expected argument %q got %q", extra, f.Args()[0])
	}
}

func TestParse(t *testing.T) {
	ResetForTesting(func() { t.Error("bad parse") })
	testParse(CommandLine, t)
}

func TestFlagSetParse(t *testing.T) {
	testParse(NewFlagSet("test", ContinueOnError), t)
}

// Declare a user-defined flag type.
type flagVar []string

func (f *flagVar) String() string {
	return fmt.Sprint([]string(*f))
}

func (f *flagVar) Set(value string) error {
	*f = append(*f, value)
	return nil
}

func (f *flagVar) Get() interface{} {
	return f
}

func TestUserDefined(t *testing.T) {
	var flags FlagSet
	flags.Init("test", ContinueOnError)
	var v flagVar
	flags.Var(&v, "v", 'v', "usage", 0, nil)
	if err := flags.Parse([]string{"--v", "1", "-v", "2", "--v=3"}); err != nil {
		t.Error(err)
	}
	if len(v) != 3 {
		t.Fatal("expected 3 args; got ", len(v))
	}
	expect := "[1 2 3]"
	if v.String() != expect {
		t.Errorf("expected value %q got %q", expect, v.String())
	}
}

func TestUserDefinedForCommandLine(t *testing.T) {
	const help = "HELP"
	var result string
	ResetForTesting(func() { result = help })
	Usage()
	if result != help {
		t.Fatalf("got %q; expected %q", result, help)
	}
}

// Declare a user-defined boolean flag type.
type boolFlagVar struct {
	count int
}

func (b *boolFlagVar) String() string {
	return fmt.Sprintf("%d", b.count)
}

func (b *boolFlagVar) Set(value string) error {
	if value == "true" {
		b.count++
	}
	return nil
}

func (b *boolFlagVar) Get() interface{} {
	return *b
}

func (b *boolFlagVar) IsBoolFlag() bool {
	return b.count < 2
}

func TestUserDefinedBool(t *testing.T) {
	var flags FlagSet
	flags.Init("test", ContinueOnError)
	flags.SetOutput(ioutil.Discard)

	var b boolFlagVar
	var err error
	flags.Var(&b, "bool", 'b', "usage", 0, nil)
	if err = flags.Parse([]string{"-b=true", "--bool=true", "-b=false", "--bool", "-b", "barg", "--bool"}); err != nil {
		if b.count < 2 {
			t.Error(err)
		}
	}

	if b.count != 2 {
		t.Errorf("want: %d; got: %d", 2, b.count)
	}

	if err == nil {
		t.Error("expected error; got none")
	}
}

func TestSetOutput(t *testing.T) {
	var flags FlagSet
	var buf bytes.Buffer
	flags.SetOutput(&buf)
	flags.Init("test", ContinueOnError)
	flags.Parse([]string{"--unknown"})
	if out := buf.String(); !strings.Contains(out, "--unknown") {
		t.Logf("expected output mentioning unknown; got %q", out)
	}
}

// Test that -help invokes the usage message and returns ErrHelp.
func TestHelp(t *testing.T) {
	var helpCalled = false
	fs := NewFlagSet("help test", ContinueOnError)
	fs.Usage = func() { helpCalled = true }
	var flag bool
	fs.BoolVar(&flag, "flag", -1, false, "regular flag", nil)
	// Regular flag invocation should work
	err := fs.Parse([]string{"--flag=true"})
	if err != nil {
		t.Fatal("expected no error; got ", err)
	}
	if !flag {
		t.Error("flag was not set by -flag")
	}
	if helpCalled {
		t.Error("help called for regular flag")
		helpCalled = false // reset for next test
	}

	for _, v := range []string{"--help", "-h"} {
		// Help flag should work as expected.
		err = fs.Parse([]string{v})
		if err == nil {
			t.Fatal("error expected")
		}
		if err != ErrHelp {
			t.Fatal("expected ErrHelp; got ", err)
		}
		if !helpCalled {
			t.Fatal("help was not called")
		}
	}

	// If we define a help flag, that should override.
	var help bool

	fs.BoolVar(&help, "help", 'h', false, "help flag", nil)
	helpCalled = false
	for _, v := range []string{"--help", "-h"} {
		err = fs.Parse([]string{v})
		if err != nil {
			t.Fatal("expected no error for defined --help; got ", err)
		}
		if helpCalled {
			t.Fatal("help was called; should not have been for defined help flag")
		}
	}
}

const defaultOutput = `Options:
  -a, --A               for bootstrapping, allow 'any' type
      --Alongflagname   disable bounds checking
  -c, --C               a boolean defaulting to true
  -d, --D               set relative path for local imports
      --E               issue 23543
      --F               a non-zero number
  -g, --G               a float that defaults to zero
      --M               a multiline
                        help
                        string
      --N               a non-zero int
      --O               a flag
                        multiline help string
      --Z               an int that defaults to zero
      --maxT            set timeout for dial

Commands:
   A,   subcmdA         A subcommand

`

const subCmdAOutput = `Options:
  -a, --subflagA        a subflag
  -b, --subflagB        b subflag
      --subflagC        c subflag

Commands:
   B,   subcmdB         B subcommand nested

`

const subCmdBOutput = `Options:
      --subflagA        a subflag
  -b, --subflagB        b subflag
  -c, --subflagC        c subflag

`

func TestPrintDefaults(t *testing.T) {
	fs := NewFlagSet("print defaults test", ContinueOnError)
	var buf bytes.Buffer
	var asubcmdnested bool
	fs.SetOutput(&buf)
	fs.Bool("A", 'a', false, "for bootstrapping, allow 'any' type", nil)
	fs.Bool("Alongflagname", -1, false, "disable bounds checking", nil)
	fs.Bool("C", 'c', true, "a boolean defaulting to true", nil)
	fs.String("D", 'd', "", "set relative path for local imports", nil)
	fs.String("E", -1, "0", "issue 23543", nil)
	fs.Float64("F", -1, 2.7, "a non-zero number", nil)
	fs.Float64("G", 'g', 0, "a float that defaults to zero", nil)
	fs.String("M", -1, "", "a multiline\nhelp\nstring", nil)
	fs.Int("N", -1, 27, "a non-zero int", nil)
	fs.Bool("O", -1, true, "a flag\nmultiline help string", nil)
	fs.Int("Z", -1, 0, "an int that defaults to zero", nil)
	fs.Duration("maxT", -1, 0, "set timeout for dial", nil)
	fs.BoolSubCommand("subcmdA", 'A', "A subcommand",
		fs.BoolSubFlag("subflagA", 'a', false, "a subflag", nil),
		fs.BoolSubFlag("subflagB", 'b', false, "b subflag", nil),
		fs.BoolVarSubCommandNested(&asubcmdnested, "subcmdB", 'B', "B subcommand nested",
			fs.BoolSubFlag("subflagB", 'b', false, "b subflag", nil),
			fs.BoolSubFlag("subflagA", -1, false, "a subflag", nil),
			fs.BoolSubFlag("subflagC", 'c', false, "c subflag", nil),
		),
		fs.BoolSubFlag("subflagC", -1, false, "c subflag", nil),
	)

	fs.PrintDefaults()
	got := buf.String()
	if got != defaultOutput {
		index := -1
		for i := 0; i < len(defaultOutput) && i < len(got); i++ {
			if got[i] != defaultOutput[i] {
				index = i
			}
		}

		t.Errorf("difference position: %d/%d\n\n- want:\n\n'%s'\n\n- got:\n\n'%s'\n",
			index, len(got), defaultOutput, got)
	}

	fs.Parse([]string{"subcmdA"})
	buf.Reset()

	fs.PrintDefaults()
	got = buf.String()
	if got != subCmdAOutput {
		index := -1
		for i := 0; i < len(subCmdAOutput) && i < len(got); i++ {
			if got[i] != subCmdAOutput[i] {
				index = i
			}
		}

		t.Errorf("difference position: %d/%d\n\n- want:\n\n'%s'\n\n- got:\n\n'%s'\n",
			index, len(got), subCmdAOutput, got)
	}

	fs.Parse([]string{"subcmdB"})
	buf.Reset()

	fs.PrintDefaults()
	got = buf.String()
	if got != subCmdBOutput {
		index := -1
		for i := 0; i < len(subCmdBOutput) && i < len(got); i++ {
			if got[i] != subCmdBOutput[i] {
				index = i
			}
		}

		t.Errorf("difference position: %d/%d\n\n- want:\n\n'%s'\n\n- got:\n\n'%s'\n",
			index, len(got), subCmdBOutput, got)
	}

}

// Issue 19230: validate range of Int and Uint flag values.
func TestIntFlagOverflow(t *testing.T) {
	if strconv.IntSize != 32 {
		return
	}
	ResetForTesting(nil)
	Int("i", -1, 0, "", nil)
	Uint("u", -1, 0, "", nil)
	if err := CommandLine.flags["i"].Value.Set("2147483648"); err == nil {
		t.Error("unexpected success setting Int")
	}

	if err := CommandLine.flags["u"].Value.Set("4294967296"); err == nil {
		t.Error("unexpected success setting Uint")
	}
}

func TestGetters(t *testing.T) {
	expectedName := "flag set"
	expectedErrorHandling := ContinueOnError
	expectedOutput := io.Writer(os.Stderr)
	fs := NewFlagSet(expectedName, expectedErrorHandling)

	if fs.Name() != expectedName {
		t.Errorf("unexpected name: got %s, expected %s", fs.Name(), expectedName)
	}
	if fs.ErrorHandling() != expectedErrorHandling {
		t.Errorf("unexpected ErrorHandling: got %d, expected %d", fs.ErrorHandling(), expectedErrorHandling)
	}
	if fs.Output() != expectedOutput {
		t.Errorf("unexpected output: got %#v, expected %#v", fs.Output(), expectedOutput)
	}

	expectedName = "gopher"
	expectedErrorHandling = ExitOnError
	expectedOutput = os.Stdout
	fs.Init(expectedName, expectedErrorHandling)
	fs.SetOutput(expectedOutput)

	if fs.Name() != expectedName {
		t.Errorf("unexpected name: got %s, expected %s", fs.Name(), expectedName)
	}
	if fs.ErrorHandling() != expectedErrorHandling {
		t.Errorf("unexpected ErrorHandling: got %d, expected %d", fs.ErrorHandling(), expectedErrorHandling)
	}
	if fs.Output() != expectedOutput {
		t.Errorf("unexpected output: got %v, expected %v", fs.Output(), expectedOutput)
	}
}

func TestParseError(t *testing.T) {
	for _, typ := range []string{"bool", "int", "int64", "uint", "uint64", "float64", "duration"} {
		fs := NewFlagSet("parse error test", ContinueOnError)
		fs.SetOutput(ioutil.Discard)
		_ = fs.Bool("bool", -1, false, "", nil)
		_ = fs.Int("int", -1, 0, "", nil)
		_ = fs.Int64("int64", -1, 0, "", nil)
		_ = fs.Uint("uint", -1, 0, "", nil)
		_ = fs.Uint64("uint64", -1, 0, "", nil)
		_ = fs.Float64("float64", -1, 0, "", nil)
		_ = fs.Duration("duration", -1, 0, "", nil)
		// Strings cannot give errors.
		args := []string{"--" + typ + "=x"}
		err := fs.Parse(args) // x is not a valid setting for any flag.
		if err == nil {
			t.Errorf("Parse(%q)=%v; expected parse error", args, err)
			continue
		}
		if !strings.Contains(err.Error(), "invalid") && !strings.Contains(err.Error(), "parse error") {
			t.Errorf("Parse(%q)=%v; expected parse error", args, err)
		}
	}
}

func TestRangeError(t *testing.T) {
	bad := []string{
		"--int=123456789012345678901",
		"--int64=123456789012345678901",
		"--uint=123456789012345678901",
		"--uint64=123456789012345678901",
		"--float64=1e1000",
	}
	for _, arg := range bad {
		fs := NewFlagSet("parse error test", ContinueOnError)
		fs.SetOutput(ioutil.Discard)
		_ = fs.Int("int", -1, 0, "", nil)
		_ = fs.Int64("int64", -1, 0, "", nil)
		_ = fs.Uint("uint", -1, 0, "", nil)
		_ = fs.Uint64("uint64", -1, 0, "", nil)
		_ = fs.Float64("float64", -1, 0, "", nil)
		// Strings cannot give errors, and bools and durations do not return strconv.NumError.
		err := fs.Parse([]string{arg})
		if err == nil {
			t.Errorf("Parse(%q)=%v; expected range error", arg, err)
			continue
		}
		if !strings.Contains(err.Error(), "invalid") && !strings.Contains(err.Error(), "value out of range") {
			t.Errorf("Parse(%q)=%v; expected range error", arg, err)
		}
	}
}

func TestAddSubCommandName(t *testing.T) {
	var B, C, D bool
	const Name = "Test Add Sub Command Name"
	args := []string{"A", "B", "C", "D"}
	expect := fmt.Sprintf("%s %s", Name, strings.Join(args, " "))

	ResetForTesting(nil)
	CommandLine.Init(Name, ContinueOnError)

	BoolSubCommand(args[0], -1, "description",
		BoolVarSubCommandNested(&B, args[1], -1, "description",
			BoolVarSubCommandNested(&C, args[2], -1, "description",
				BoolVarSubCommandNested(&D, args[3], -1, "description"),
			),
		),
	)

	CommandLine.Parse(args)
	actual := CommandLine.Name()
	if actual != expect {
		t.Errorf("\ngot : %s, want: %s\n", actual, expect)
	}
}

func TestSubCommand(t *testing.T) {
	data := []struct {
		args   []string
		expect map[string]bool
	}{
		{
			args: []string{"A", "-b"},
			expect: map[string]bool{
				"A": true, "B": true, "C": false, "D": false, "E": false, "F": false,
			},
		},
		{
			args: []string{"A", "C", "--D"},
			expect: map[string]bool{
				"A": true, "B": false, "C": true, "D": true, "E": false, "F": false,
			},
		},
		{
			args: []string{"A", "C", "E"},
			expect: map[string]bool{
				"A": true, "B": false, "C": true, "D": false, "E": true, "F": false,
			},
		},
		{
			args: []string{"A", "C", "E", "-f"},
			expect: map[string]bool{
				"A": true, "B": false, "C": true, "D": false, "E": true, "F": true,
			},
		},
	}

	var A, B, C, D, E, F bool
	for i, v := range data {
		ResetForTesting(nil)
		BoolVarSubCommand(&A, "A", 'a', "description",
			BoolVarSubFlag(&B, "B", 'b', false, "description", nil),
			BoolVarSubCommandNested(&C, "C", 'c', "description",
				BoolVarSubFlag(&D, "D", 'd', false, "description", nil),
				BoolVarSubCommandNested(&E, "E", 'e', "description",
					BoolVarSubFlag(&F, "F", 'f', false, "description", nil),
				),
			),
		)

		if err := CommandLine.Parse(v.args); err != nil {
			t.Errorf(" %d: error: %s\n", i, err)
			continue
		}
		actual := map[string]bool{
			"A": A, "B": B, "C": C, "D": D, "E": E, "F": F,
		}

		for k, v := range v.expect {
			if actual[k] != v {
				t.Errorf(" %d: flag name `%s' - got: %t, want: %t\n", i, k, actual[k], v)
			}
		}
	}
}

func TestShortFlag(t *testing.T) {
	data := []struct {
		name            string
		args            []string
		shouldBeAnError bool
	}{
		{
			name:            "A consecutive short flags",
			args:            []string{"-abc"},
			shouldBeAnError: false,
		},
		{
			name:            "Equals can not be used with consecutive short flags",
			args:            []string{"-bcd=value"},
			shouldBeAnError: true,
		},
		{
			name:            "Specify the value by a equal",
			args:            []string{"-d=value"},
			shouldBeAnError: false,
		},
		{
			name:            "Unrecognized short flags",
			args:            []string{"-xyz"},
			shouldBeAnError: true,
		},
		{
			name:            "Unrecognized short flags",
			args:            []string{"-z=value"},
			shouldBeAnError: true,
		},
	}

	for _, v := range data {
		ResetForTesting(nil)
		CommandLine.SetOutput(ioutil.Discard)
		Bool("A", 'a', false, "", nil)
		Bool("B", 'b', false, "", nil)
		Bool("C", 'c', false, "", nil)
		String("D", 'd', "", "", nil)

		err := CommandLine.Parse(v.args)

		if v.shouldBeAnError {
			if err == nil {
				t.Errorf("\n%s: should be an error\n", v.name)
			}
			continue
		}

		if err != nil {
			t.Errorf("\n%s: error: %s\n", v.name, err)
		}
	}
}
