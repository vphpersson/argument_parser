package option

import (
	"fmt"
	motmedelErrors "github.com/Motmedel/utils_go/pkg/errors"
	argumentParserErrors "github.com/vphpersson/argument_parser/pkg/errors"
	"os"
	"strconv"
)

type Option interface {
	Set(string) error
	GetShortName() string
	GetLongName() string
	GetUsage() string
	GetRequired() bool
	GetNargs() string
}

type base struct {
	ShortName rune
	LongName  string
	Usage     string
	Required  bool
	Nargs     string
}

func (b *base) GetShortName() string {
	return string(b.ShortName)
}

func (b *base) GetLongName() string {
	return b.LongName
}

func (b *base) GetUsage() string {
	return b.Usage
}

func (b *base) GetRequired() bool {
	return b.Required
}

func (b *base) GetNargs() string {
	return b.Nargs
}

type IntOption struct {
	base
	Value *int
}

func (intVar *IntOption) Set(in string) error {
	value, err := strconv.Atoi(in)
	if err != nil {
		return motmedelErrors.NewWithTrace(fmt.Errorf("strvconv atoi: %w", err))
	}

	if intVar.Value == nil {
		return motmedelErrors.NewWithTrace(argumentParserErrors.ErrNilValue)
	}

	*intVar.Value = value

	return nil
}

func NewIntOption(shortName rune, longName string, usage string, required bool, value *int) *IntOption {
	return &IntOption{
		base: base{
			ShortName: shortName,
			LongName:  longName,
			Usage:     usage,
			Required:  required,
		},
		Value: value,
	}
}

type IntsOption struct {
	base
	Value *[]int
}

func (intsOption *IntsOption) Set(in string) error {
	if intsOption.Value == nil {
		return motmedelErrors.NewWithTrace(argumentParserErrors.ErrNilValue)
	}

	value, err := strconv.Atoi(in)
	if err != nil {
		return motmedelErrors.NewWithTrace(fmt.Errorf("strvconv atoi: %w", err))
	}

	*intsOption.Value = append(*intsOption.Value, value)
	return nil
}

func NewIntsOption(shortName rune, longName string, usage string, required bool, value *[]int) *IntsOption {
	return &IntsOption{
		base: base{
			ShortName: shortName,
			LongName:  longName,
			Usage:     usage,
			Required:  required,
			Nargs:     "+",
		},
		Value: value,
	}
}

type StringOption struct {
	base
	Value *string
}

func (stringOption *StringOption) Set(in string) error {
	if stringOption.Value == nil {
		return motmedelErrors.NewWithTrace(argumentParserErrors.ErrNilValue)
	}

	*stringOption.Value = in

	return nil
}

func NewStringOption(shortName rune, longName string, usage string, required bool, value *string) *StringOption {
	return &StringOption{
		base: base{
			ShortName: shortName,
			LongName:  longName,
			Usage:     usage,
			Required:  required,
		},
		Value: value,
	}
}

type StringsOption struct {
	base
	Value *[]string
}

func (stringsOption *StringsOption) Set(in string) error {
	if stringsOption.Value == nil {
		return motmedelErrors.NewWithTrace(argumentParserErrors.ErrNilValue)
	}

	*stringsOption.Value = append(*stringsOption.Value, in)
	return nil
}

func NewStringsOption(shortName rune, longName string, usage string, required bool, value *[]string) *StringsOption {
	return &StringsOption{
		base: base{
			ShortName: shortName,
			LongName:  longName,
			Usage:     usage,
			Required:  required,
			Nargs:     "+",
		},
		Value: value,
	}
}

type BoolOption struct {
	base
	Value *bool
}

func (boolOption *BoolOption) Set(in string) error {
	if boolOption.Value == nil {
		return motmedelErrors.NewWithTrace(argumentParserErrors.ErrNilValue)
	}

	*boolOption.Value = true

	return nil
}

func NewBoolOption(shortName rune, longName string, usage string, required bool, value *bool) *BoolOption {
	return &BoolOption{
		base: base{
			ShortName: shortName,
			LongName:  longName,
			Usage:     usage,
			Required:  required,
			Nargs:     "0",
		},
		Value: value,
	}
}

type CountedOption struct {
	base
	Count *int
}

func (countedOption *CountedOption) Set(in string) error {
	if countedOption.Count == nil {
		return motmedelErrors.NewWithTrace(argumentParserErrors.ErrNilValue)
	}

	*countedOption.Count++

	return nil
}

func NewFileOption(shortName rune, longName string, usage string, required bool, file *os.File) *FileOption {
	return &FileOption{
		base: base{
			ShortName: shortName,
			LongName:  longName,
			Usage:     usage,
			Required:  required,
		},
		File: file,
	}
}

func NewFileOptionExtra(
	shortName rune,
	longName string,
	usage string,
	required bool,
	flag int,
	mode os.FileMode,
	file *os.File,
) *FileOption {
	return &FileOption{
		base: base{
			ShortName: shortName,
			LongName:  longName,
			Usage:     usage,
			Required:  required,
		},
		File: file,
		Flag: flag,
		Mode: mode,
	}
}

type FileOption struct {
	base
	File *os.File
	Flag int
	Mode os.FileMode
}

func (fileOption *FileOption) Set(in string) error {
	if fileOption.File == nil {
		return motmedelErrors.NewWithTrace(argumentParserErrors.ErrNilValue)
	}

	value, err := os.OpenFile(in, fileOption.Flag, fileOption.Mode)
	if err != nil {
		return motmedelErrors.NewWithTrace(fmt.Errorf("os open file: %w", err))
	}

	*fileOption.File = *value

	return nil
}
