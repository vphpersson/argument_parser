package argument_parser

import (
	"errors"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	argumentParserErrors "github.com/vphpersson/argument_parser/pkg/errors"
	"github.com/vphpersson/argument_parser/pkg/types/option"
)

var diffOpts = []cmp.Option{cmpopts.EquateEmpty()}

func TestArgumentParserParseArgs(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name      string
		options   []option.Option
		args      []string
		intValue  int
		intsValue    []int
		strValue  string
		stringsValue []string
		boolValue bool
		wantErr   error
	}{
		{
			name: "full",
			options: []option.Option{
				option.NewIntOption('i', "int", "An int option", false, nil),
				option.NewStringOption('s', "str", "A string option", false, nil),
				option.NewBoolOption('b', "bool", "A bool option", false, nil),
				option.NewStringsOption('a', "array", "An array option", false, nil),
				option.NewIntsOption('n', "numbers", "An array of ints option", false, nil),
			},
			args:      []string{"-i", "42", "--str", "abc", "--bool", "-a", "a", "b", "-n", "1", "2"},
			intValue:  42,
			intsValue: []int{1, 2},
			strValue:  "abc",
			stringsValue: []string{"a", "b"},
			boolValue: true,
			wantErr:   nil,
		},
		{
			name: "multiple same",
			options: []option.Option{
				option.NewIntOption('i', "int", "An int option", false, nil),
			},
			args:     []string{"-i", "1", "-i", "2"},
			intValue: 2,
		},
		{
			name: "unset option 1",
			options: []option.Option{
				option.NewIntOption('i', "int", "An int option", false, nil),
			},
			args:    []string{"-i"},
			wantErr: argumentParserErrors.ErrUnsetOption,
		},
		{
			name: "unset option 2",
			options: []option.Option{
				option.NewIntOption('i', "int", "An int option", false, nil),
				option.NewBoolOption('b', "bool", "A bool option", false, nil),
			},
			args:    []string{"-i", "42", "--bool", "--int"},
			wantErr: argumentParserErrors.ErrUnsetOption,
		},
		{
			name: "name not found",
			options: []option.Option{
				option.NewIntOption('i', "int", "An int option", false, nil),
			},
			args:    []string{"--unknown", "1"},
			wantErr: argumentParserErrors.ErrNameNotFound,
		},
		{
			name: "name not found",
			options: []option.Option{
				option.NewIntOption('i', "ii", "An int option", false, nil),
				option.NewIntOption('p', "pp", "An int option", false, nil),
			},
			args:    []string{"-ip"},
			wantErr: argumentParserErrors.ErrNonZeroArgCombinedOption,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			var intValue int
			var stringValue string
			var boolValue bool
			stringsValue := make([]string, 0)
			intsValue := make([]int, 0)

			for _, opt := range testCase.options {
				switch typedOpt := opt.(type) {
				case *option.IntOption:
					typedOpt.Value = &intValue
				case *option.StringOption:
					typedOpt.Value = &stringValue
				case *option.BoolOption:
					typedOpt.Value = &boolValue
				case *option.StringsOption:
					typedOpt.Value = &stringsValue
				case *option.IntsOption:
					typedOpt.Value = &intsValue
				default:
					t.Fatalf("Unexpected option type: %T", opt)
					return
				}
			}

			parser := &ArgumentParser{Options: testCase.options}

			err := parser.ParseArgs(testCase.args)
			if testCase.wantErr != nil {
				if !errors.Is(err, testCase.wantErr) {
					t.Errorf("unexpected error = %v, want %v", err, testCase.wantErr)
					return
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error = %v", err)
				return
			}

			if intValue != testCase.intValue {
				t.Errorf("Expected int value = %v, got %v", testCase.intValue, intValue)
			}

			if stringValue != testCase.strValue {
				t.Errorf("Expected string value = %v, got %v", testCase.strValue, stringValue)
			}

			if boolValue != testCase.boolValue {
				t.Errorf("Expected bool value = %v, got %v", testCase.boolValue, boolValue)
			}

			if diff := cmp.Diff(testCase.stringsValue, stringsValue, diffOpts...); diff != "" {
				t.Errorf("strings value mismatch (-expected +got):\n%s", diff)
			}

			if diff := cmp.Diff(testCase.intsValue, intsValue, diffOpts...); diff != "" {
				t.Errorf("ints value mismatch (-expected +got):\n%s", diff)
			}
		})
	}
}

func TestGetArgumentName(t *testing.T) {
	t.Parallel()

	cases := []struct {
		in   string
		want []string
	}{
		{"", nil},
		{"-n", []string{"n"}},
		{"--filename", []string{"filename"}},
		{"-yes", []string{"y", "e", "s"}},
		{"no", nil},
		{"nope--nope", nil},
		{"nope-nope", nil},
	}

	for _, tc := range cases {
		t.Run(tc.in, func(t *testing.T) {
			t.Parallel()
			got := getArgumentNames(tc.in)

			if diff := cmp.Diff(tc.want, got, diffOpts...); diff != "" {
				t.Errorf("mismatch (-expected +got):\n%s", diff)
			}
		})
	}
}

func TestMakeNameToOption(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name     string
		options  []option.Option
		wantErr  error
		contains map[string]bool // map of option names expected in the result
	}{
		{
			name:    "No options",
			options: nil,
			wantErr: nil,
		},
		{
			name: "Duplicate long name",
			options: []option.Option{
				option.NewIntOption('a', "same", "usage", false, nil),
				option.NewIntOption('b', "same", "usage2", false, nil),
			},
			wantErr: argumentParserErrors.ErrMultipleOptionsWithSameName,
		},
		{
			name: "OK names",
			options: []option.Option{
				option.NewIntOption('a', "as", "usage", false, nil),
				option.NewIntOption('b', "bs", "usage2", false, nil),
			},
			wantErr: nil,
			contains: map[string]bool{
				"a":  true,
				"as": true,
				"b":  true,
				"bs": true,
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()
			got, err := makeNameToOption(testCase.options)
			if testCase.wantErr != nil {
				if !errors.Is(err, testCase.wantErr) {
					t.Errorf("Expected error %v, got %v", testCase.wantErr, err)
				}
				return
			}
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}
			if testCase.contains != nil {
				for name, wantPresent := range testCase.contains {
					_, found := got[name]
					if found != wantPresent {
						t.Errorf("Name %q present in map: %v, want %v", name, found, wantPresent)
					}
				}
			}
		})
	}
}
