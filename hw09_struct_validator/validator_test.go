package hw09structvalidator

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

type UserRole string

// Test the function on different structures and other types.
type (
	User struct {
		ID     string `json:"id" validate:"len:36"`
		Name   string
		Age    int             `validate:"min:18|max:50"`
		Email  string          `validate:"regexp:^\\w+@\\w+\\.\\w+$"`
		Role   UserRole        `validate:"in:admin,stuff"`
		Phones []string        `validate:"len:11"`
		meta   json.RawMessage //nolint
	}

	App struct {
		Version string `validate:"len:5"`
	}

	Token struct {
		Header    []byte
		Payload   []byte
		Signature []byte
	}

	Response struct {
		Code int    `validate:"in:200,404,500"`
		Body string `json:"omitempty"`
	}

	SliceOfInts struct {
		Value []int `validate:"min:1|max:100"`
	}

	BrokenSyntax struct {
		Value string `validate:"len:"`
	}

	BrokenSyntax2 struct {
		Value string `validate:":"`
	}

	BrokenSyntax3 struct {
		Value string `validate:"len"`
	}

	BrokenSyntax4 struct {
		Value int `validate:"min:"`
	}

	BrokenSyntax5 struct {
		Value int `validate:":"`
	}
)

func TestValidate(t *testing.T) {
	tests := []struct {
		in          interface{}
		expectedErr error
	}{
		{
			in:          123,
			expectedErr: ErrUnsupportedType,
		},
		{
			in:          BrokenSyntax{},
			expectedErr: ErrConstraintIsInvalid,
		},
		{
			in:          BrokenSyntax2{},
			expectedErr: ErrConstraintIsInvalid,
		},
		{
			in:          BrokenSyntax3{},
			expectedErr: ErrConstraintIsInvalid,
		},
		{
			in:          BrokenSyntax4{},
			expectedErr: ErrConstraintIsInvalid,
		},
		{
			in:          BrokenSyntax5{},
			expectedErr: ErrConstraintIsInvalid,
		},
		{
			in: SliceOfInts{
				Value: []int{-10, 0, 150},
			},
			expectedErr: ValidationErrors{
				{
					Field: "Value[0]",
					Err:   ErrIntLessThen(1),
				},
				{
					Field: "Value[1]",
					Err:   ErrIntLessThen(1),
				},
				{
					Field: "Value[2]",
					Err:   ErrIntGreaterThen(100),
				},
			},
		},
		{
			in: Response{
				Code: 0,
				Body: "",
			},
			expectedErr: ValidationErrors{
				{
					Field: "Code",
					Err:   ErrValueNotIn{"200", "404", "500"},
				},
			},
		},
		{
			in: User{
				ID:    "363636363636363636363636363636363636",
				Name:  "",
				Age:   60,
				Email: "test@test.test",
				Role:  "admin",
				Phones: []string{
					"91123456789",
				},
			},
			expectedErr: ValidationErrors{
				{
					Field: "Age",
					Err:   ErrIntGreaterThen(50),
				},
			},
		},
		{
			in: App{
				Version: "a",
			},
			expectedErr: ValidationErrors{
				{
					Field: "Version",
					Err:   ErrStringLengthIsInvalid,
				},
			},
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			tt := tt
			t.Parallel()

			err := Validate(tt.in)
			require.Equal(t, tt.expectedErr, err)
		})
	}
}
