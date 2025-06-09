package option

import (
	"fmt"
	motmedelErrors "github.com/Motmedel/utils_go/pkg/errors"
	argumentParserErrors "github.com/vphpersson/argument_parser/pkg/errors"
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

func (bv *base) GetShortName() string {
	return string(bv.ShortName)
}

func (bv *base) GetLongName() string {
	return bv.LongName
}

func (bv *base) GetUsage() string {
	return bv.Usage
}

func (bv *base) GetRequired() bool {
	return bv.Required
}

func (bv *base) GetNargs() string {
	return bv.Nargs
}

type IntVar struct {
	base
	Value *int
}

func NewIntVar(shortName rune, longName string, usage string, required bool, value *int) *IntVar {
	return &IntVar{
		base: base{
			ShortName: shortName,
			LongName:  longName,
			Usage:     usage,
			Required:  required,
		},
		Value: value,
	}
}

func (intVar *IntVar) Set(in string) error {
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

type StringVar struct {
	base
	Value *string
}

func NewStringVar(shortName rune, longName string, usage string, required bool, value *string) *StringVar {
	return &StringVar{
		base: base{
			ShortName: shortName,
			LongName:  longName,
			Usage:     usage,
			Required:  required,
		},
		Value: value,
	}
}

type BoolVar struct {
	base
	Value *bool
}

func (boolVar *BoolVar) Set(in string) error {
	if boolVar.Value == nil {
		return motmedelErrors.NewWithTrace(argumentParserErrors.ErrNilValue)
	}

	*boolVar.Value = true

	return nil
}

func NewBoolVar(shortName rune, longName string, usage string, required bool, value *bool) *BoolVar {
	return &BoolVar{
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

func (stringVar *StringVar) Set(in string) error {
	if stringVar.Value == nil {
		return motmedelErrors.NewWithTrace(argumentParserErrors.ErrNilValue)
	}

	*stringVar.Value = in

	return nil
}

type CountedVar struct {
	base
	Count *int
}

func (countedVar *CountedVar) Set(in string) error {
	if countedVar.Count == nil {
		return motmedelErrors.NewWithTrace(argumentParserErrors.ErrNilValue)
	}

	*countedVar.Count++

	return nil
}
