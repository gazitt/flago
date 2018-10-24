package flago

const (
	NAME = "test-command"
)

var DefaultUsage = Usage

func ResetForTesting(usage func()) {
	CommandLine = NewFlagSet(NAME, ContinueOnError)
	CommandLine.Usage = commandLineUsage
	Usage = usage
}
