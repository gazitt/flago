package flago

import (
	"flag"
	"testing"
)

var (
	testArgs = []string{
		"--A",
		"--B=100",
		"--C",
		"10.5",
		"--D=hello",
		"--E",
		"100",
		"argument-1",
		"argument-2",
		"argument-3",
		"argument-4",
		"argument-5",
		"argument-6",
	}
)

func BenchmarkFlag(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		fs := flag.NewFlagSet("test", flag.ContinueOnError)
		fs.Bool("A", false, "description")
		fs.Int("B", 0, "description")
		fs.Float64("C", 0, "description")
		fs.String("D", "", "description")
		fs.Uint64("E", 0, "description")
		fs.Parse(testArgs)
	}
}

func BenchmarkAliasNotSpecified(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		fs := NewFlagSet("test", ContinueOnError)
		fs.Bool("A", -1, false, "description", nil)
		fs.Int("B", -1, 0, "description", nil)
		fs.Float64("C", -1, 0, "description", nil)
		fs.String("D", -1, "", "description", nil)
		fs.Uint64("E", -1, 0, "description", nil)
		fs.Parse(testArgs)
	}
}

func BenchmarkAliasSpecified(b *testing.B) {
	args := []string{
		"argument-1",
		"--A",
		"argument-2",
		"-b=100",
		"argument-3",
		"--C",
		"10.5",
		"argument-4",
		"-d",
		"hello",
		"argument-5",
		"--E=100",
		"argument-6",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		fs := NewFlagSet("test", ContinueOnError)
		fs.Bool("A", 'a', false, "description", nil)
		fs.Int("B", 'b', 0, "description", nil)
		fs.Float64("C", 'c', 0, "description", nil)
		fs.String("D", 'd', "", "description", nil)
		fs.Uint64("E", 'e', 0, "description", nil)
		fs.Parse(args)
	}
}
func BenchmarkAliasContinuous(b *testing.B) {
	args := []string{
		"argument-1",
		"-abc",
		"argument-2",
		"-de",
		"argument-3",
		"argument-4",
		"argument-5",
		"argument-6",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		fs := NewFlagSet("test", ContinueOnError)
		fs.Bool("A", 'a', false, "description", nil)
		fs.Bool("B", 'b', false, "description", nil)
		fs.Bool("C", 'c', false, "description", nil)
		fs.Bool("D", 'd', false, "description", nil)
		fs.Bool("E", 'e', false, "description", nil)
		fs.Parse(args)
	}
}

func BenchmarkSubCommand(b *testing.B) {
	args := []string{
		"A",
		"argument-1",
		"--B",
		"argument-2",
		"-c",
		"argument-3",
		"--D",
		"argument-4",
		"-e",
		"argument-5",
		"argument-6",
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		fs := NewFlagSet("test", ContinueOnError)
		fs.BoolSubCommand("A", 'a', "description",
			fs.BoolSubFlag("B", 'b', false, "description", nil),
			fs.BoolSubFlag("C", 'c', false, "description", nil),
			fs.BoolSubFlag("D", 'd', false, "description", nil),
			fs.BoolSubFlag("E", 'e', false, "description", nil),
		)
		fs.Parse(args)
	}
}
