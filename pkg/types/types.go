package types

type Parser interface {
	ParseArgs([]string) error
	GetCommand() string
}
