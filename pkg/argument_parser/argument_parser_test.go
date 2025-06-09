package argument_parser

import (
	"errors"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	argumentParserErrors "github.com/vphpersson/argument_parser/pkg/errors"
	"github.com/vphpersson/argument_parser/pkg/types/option"
	"testing"
)

var diffOpts = []cmp.Option{cmpopts.EquateEmpty()}

func TestArgumentParserParseArgs(t *testing.T) {
	t.Parallel()

	type testCase struct {
		name      string
		options   []option.Option
		args      []string
		intValue  int
		strValue  string
		boolValue bool
		wantErr   error
	}

	testCases := []testCase{
		{
			name: "full",
			options: []option.Option{
				option.NewIntVar('i', "int", "An int option", false, nil),
				option.NewStringVar('s', "str", "A string option", false, nil),
				option.NewBoolVar('b', "bool", "A bool option", false, nil),
			},
			args:      []string{"-i", "42", "--str", "abc", "--bool"},
			intValue:  42,
			strValue:  "abc",
			boolValue: true,
			wantErr:   nil,
		},
		{
			name: "multiple same",
			options: []option.Option{
				option.NewIntVar('i', "int", "An int option", false, nil),
			},
			args:     []string{"-i", "1", "-i", "2"},
			intValue: 2,
		},
		{
			name: "unset option 1",
			options: []option.Option{
				option.NewIntVar('i', "int", "An int option", false, nil),
			},
			args:    []string{"-i"},
			wantErr: argumentParserErrors.ErrUnsetOption,
		},
		{
			name: "unset option 2",
			options: []option.Option{
				option.NewIntVar('i', "int", "An int option", false, nil),
				option.NewBoolVar('b', "bool", "A bool option", false, nil),
			},
			args:    []string{"-i", "42", "--bool", "--int"},
			wantErr: argumentParserErrors.ErrUnsetOption,
		},
		{
			name: "name not found",
			options: []option.Option{
				option.NewIntVar('i', "int", "An int option", false, nil),
			},
			args:    []string{"--unknown", "1"},
			wantErr: argumentParserErrors.ErrNameNotFound,
		},
		{
			name: "name not found",
			options: []option.Option{
				option.NewIntVar('i', "ii", "An int option", false, nil),
				option.NewIntVar('p', "pp", "An int option", false, nil),
			},
			args:    []string{"-ip"},
			wantErr: argumentParserErrors.ErrNonZeroArgCombinedOption,
		},
		// Add more test cases as needed
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var intVal int
			var strVal string
			var boolVal bool
			for _, opt := range tt.options {
				switch v := opt.(type) {
				case *option.IntVar:
					v.Value = &intVal
				case *option.StringVar:
					v.Value = &strVal
				case *option.BoolVar:
					v.Value = &boolVal
				}
			}

			parser := &ArgumentParser{Options: tt.options}

			err := parser.ParseArgs(tt.args)
			if tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Errorf("unexpected error = %v, want %v", err, tt.wantErr)
					return
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error = %v", err)
				return
			}

			if tt.intValue != 0 && intVal != tt.intValue {
				t.Errorf("Expected int value = %v, got %v", tt.intValue, intVal)
			}

			if tt.strValue != "" && strVal != tt.strValue {
				t.Errorf("Expected string value = %v, got %v", tt.strValue, strVal)
			}
			// Similarly test for boolVal if used in cases
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
	tests := []struct {
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
				option.NewIntVar('a', "same", "usage", false, nil),
				option.NewIntVar('b', "same", "usage2", false, nil),
			},
			wantErr: argumentParserErrors.ErrMultipleOptionsWithSameName,
		},
		{
			name: "OK names",
			options: []option.Option{
				option.NewIntVar('a', "as", "usage", false, nil),
				option.NewIntVar('b', "bs", "usage2", false, nil),
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

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			got, err := makeNameToOption(tc.options)
			if tc.wantErr != nil {
				if !errors.Is(err, tc.wantErr) {
					t.Errorf("Expected error %v, got %v", tc.wantErr, err)
				}
				return
			}
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}
			if tc.contains != nil {
				for name, wantPresent := range tc.contains {
					_, found := got[name]
					if found != wantPresent {
						t.Errorf("Name %q present in map: %v, want %v", name, found, wantPresent)
					}
				}
			}
		})
	}
}
