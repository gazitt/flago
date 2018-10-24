package flago

import (
	"fmt"
	"net/url"
)

type URLValue struct {
	URL *url.URL
}

func (v URLValue) String() string {
	if v.URL != nil {
		return v.URL.String()
	}
	return ""
}

func (v URLValue) Set(s string) error {
	if u, err := url.Parse(s); err != nil {
		return err
	} else {
		*v.URL = *u
	}
	return nil
}

func (v URLValue) Get() interface{} {
	return v
}

var u = &url.URL{}

func ExampleValue() {
	fs := NewFlagSet("example_value_test", ContinueOnError)
	fs.Var(&URLValue{u}, "url", 'u', "URL to parse", 0, nil)

	fs.Parse([]string{"--url", "https://golang.org/pkg/flag/"})
	fmt.Printf(`{scheme: %q, host: %q, path: %q}`, u.Scheme, u.Host, u.Path)

	// Output:
	// {scheme: "https", host: "golang.org", path: "/pkg/flag/"}
}
