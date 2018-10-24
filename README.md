
## flago

command-line flag parser

<br/>

# Features

* Options with long and short names
* Supports callback function
* Supports sub command

<br/>

# Installation

```
$ go get github.com/gazitt/flago
```

<br/>

# Command line flag syntax

```

--flag
--flag value
--flag=value

-f
-f value
-f=value

// mixed
-abc
-abc value

// subcommand
$ command init -a -b -c

// there are other values in before a flag, does not end
$ command other -a other -b other -c other

```

<br/>

# Usage

```go

package main

import (
	"fmt"
	"os"
	"time"

	"github.com/gazitt/flago"
)

type Options struct {
	Bool     bool
	Duration time.Duration
	Float64  float64
	Int      int
	Int64    int64
	String   string
	Uint     uint
	Uint64   uint64
}

var (
	o Options
)

func init() {
	// bool
	flago.BoolVar(&o.Bool, "bool-flag", 'b', false, "description", nil)
	// time.Duration
	flago.DurationVar(&o.Duration, "duration-flag", 'd', 10*time.Second, "description", nil)
	// float64
	flago.Float64Var(&o.Float64, "float64-flag", 'f', 10.5, "description", nil)
	// int
	flago.IntVar(&o.Int, "int-flag", 'i', 100, "description", nil)
	// int64
	flago.Int64Var(&o.Int64, "int64-flag", 'I', 200, "description", nil)
	// string
	flago.StringVar(&o.String, "string-flag", 's', "", "description", nil)
	// uint
	flago.UintVar(&o.Uint, "uint-flag", 'u', 100, "description", nil)
	// uint64
	flago.Uint64Var(&o.Uint64, "uint64-flag", 'U', 200, "description", nil)
}

func main() {
	flago.Parse()

	fmt.Println("bool:", o.Bool)
	fmt.Println("duration:", o.Duration)
	fmt.Println("float64:", o.Float64)
	fmt.Println("int:", o.Int)
	fmt.Println("int64:", o.Int64)
	fmt.Println("string:", o.String)
	fmt.Println("uint:", o.Uint)
	fmt.Println("uint64:", o.Uint64)

	fmt.Println("argv:", flago.Args())
}

```

<br/>

* Not using short flag

```go
// Specify less than to 1
flago.BoolVar(&o.Bool, "bool-flag", 0, false, "description", nil)
```

<br/>

* Callback function

```go
flago.Bool("version", 'v', false, "Output version information and exit", func(_ flago.Value) error {
	return fmt.Errorf("%s v0.0.1", flago.CommandLine.Name())
})

flago.String("string-number", 'n', "", "description", func(v flago.Value) error {
	for _, vv := range v.String() {
		if vv < 48 || vv > 57 {
			return fmt.Errorf("`%s' is not a number", string(vv))
		}
	}
	return nil
})
```

<br/>

* SubCommand

```go
flago.BoolVarSubCommand(&o.Init.Flag, "init", -1, "description",
	flago.BoolVarSubFlag(&o.Init.Bare, "bare", -1, false, "descriptin", nil),
	flago.BoolVarSubFlag(&o.Init.Quiet, "guiet", 'q', false, "descriptin", nil),
)

flago.BoolVarSubCommand(&o.Log.Flag, "log", -1, "description",
	flago.BoolVarSubFlag(&o.Log.Oneline, "oneline", -1, false, "descriptin", nil),
	flago.BoolVarSubFlag(&o.Log.Graph, "graph", -1, false, "descriptin", nil),

	// Further an additional subcommand
	flago.BoolVarSubCommandNested(&o...
		flago.BoolVarSubFlag(&o...
		flago.BoolVarSubCommandNested(&o...

)
```
