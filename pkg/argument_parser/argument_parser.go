package argument_parser

import (
	"fmt"
	motmedelErrors "github.com/Motmedel/utils_go/pkg/errors"
	argumentParserErrors "github.com/vphpersson/argument_parser/pkg/errors"
	"github.com/vphpersson/argument_parser/pkg/types"
	"github.com/vphpersson/argument_parser/pkg/types/option"
	"os"
	"strings"
)

type ArgumentParser struct {
	Options []option.Option
	Parsers []types.Parser
}

func (argumentParser *ArgumentParser) Parse() error {
	osArgs := os.Args
	if len(osArgs) <= 1 {
		return nil
	}

	arguments := os.Args[1:]
	if err := argumentParser.ParseArgs(arguments); err != nil {
		return motmedelErrors.New(fmt.Errorf("parse args: %w", err), arguments)
	}

	return nil
}

func getArgumentNames(argument string) []string {
	if argument == "" {
		return nil
	}

	if before, longName, found := strings.Cut(argument, "--"); found && before == "" {
		return []string{longName}
	} else if before, shortName, found := strings.Cut(argument, "-"); found && before == "" {
		return strings.Split(shortName, "")
	} else {
		return nil
	}
}

func makeNameToOption(options []option.Option) (map[string]option.Option, error) {
	if len(options) == 0 {
		return nil, nil
	}

	nameToOption := make(map[string]option.Option)
	for _, opt := range options {
		if opt == nil {
			continue
		}

		for _, name := range []string{opt.GetShortName(), opt.GetLongName()} {
			if name == "" {
				continue
			}

			if _, ok := nameToOption[name]; ok {
				return nil, motmedelErrors.NewWithTrace(
					fmt.Errorf(
						"%w: %s",
						argumentParserErrors.ErrMultipleOptionsWithSameName,
						name,
					),
				)
			}

			nameToOption[name] = opt
		}
	}

	return nameToOption, nil
}

func (argumentParser *ArgumentParser) ParseArgs(arguments []string) error {
	if len(arguments) == 0 {
		return nil
	}

	if parsers := argumentParser.Parsers; len(parsers) > 0 {
		firstArg := arguments[0]

		for _, parser := range parsers {
			if parser == nil {
				continue
			}

			if parser.GetCommand() == firstArg {
				subCommandArguments := arguments[1:]
				if err := parser.ParseArgs(subCommandArguments); err != nil {
					return motmedelErrors.New(
						fmt.Errorf("subcommand parse args: %w", err),
						parser,
						subCommandArguments,
					)
				}

				return nil
			}
		}
	}

	options := argumentParser.Options
	nameToOption, err := makeNameToOption(options)
	if err != nil {
		return motmedelErrors.New(fmt.Errorf("make name to options: %w", err), options)
	}
	if len(nameToOption) == 0 {
		return nil
	}

	// TODO: Do something with positionals/rest?

	var currentOptionPtr *option.Option
	var currentName string
	var optionSet bool

	for _, argument := range arguments {
		switch argument {
		case "-h", "--help":
			// TODO: Do something helpful...
		case "--":
			// TODO: Do something with positionals/rest?
			return nil
		default:
			if names := getArgumentNames(argument); len(names) > 0 {
				for _, name := range names {
					opt, ok := nameToOption[name]
					if !ok {
						return motmedelErrors.NewWithTrace(
							fmt.Errorf("%w: %s", argumentParserErrors.ErrNameNotFound, name),
						)
					}
					if options == nil {
						return motmedelErrors.NewWithTrace(argumentParserErrors.ErrNilOption)
					}

					if currentOptionPtr != nil && !optionSet {
						return motmedelErrors.NewWithTrace(
							fmt.Errorf("%w: %s", argumentParserErrors.ErrUnsetOption, currentName),
						)
					}

					// Check if the argument is a combined option that takes a value.
					if len(names) > 0 && opt.GetNargs() != "0" {
						return motmedelErrors.NewWithTrace(
							fmt.Errorf("%w: %s", argumentParserErrors.ErrNonZeroArgCombinedOption, name),
						)
					}

					if opt.GetNargs() == "0" {
						if err := opt.Set(""); err != nil {
							return motmedelErrors.New(fmt.Errorf("option set: %w", err), opt)
						}
					} else {
						currentOptionPtr = &opt
						currentName = name
						optionSet = false
					}
				}
			} else {
				// TODO: Support positional arguments?
				if currentOptionPtr == nil {
					return motmedelErrors.NewWithTrace(argumentParserErrors.ErrNoCurrentOption)
				}

				currentOption := *currentOptionPtr
				if currentOption == nil {
					return motmedelErrors.NewWithTrace(argumentParserErrors.ErrNilOption)
				}

				if err := currentOption.Set(argument); err != nil {
					return motmedelErrors.New(fmt.Errorf("option set: %w", err), currentOption, argument)
				}

				optionSet = true
			}
		}
	}

	if currentOptionPtr != nil && !optionSet {
		return motmedelErrors.NewWithTrace(
			fmt.Errorf("%w: %s", argumentParserErrors.ErrUnsetOption, currentName),
		)
	}

	return nil
}
