package validator

import (
	"strings"
	"testing"
)

func TestValidateURL(t *testing.T){
	tests := []struct{
		name string
		input string
		wantErr bool
	}{
		{
			name: "valid http URL",
			input: "http://github.com",
			wantErr: false,
		},
		{
			name: "valid https URL",
			input: "https://github.com",
			wantErr: false,
		},
		{
			name: "empty URL",
			input: "",
			wantErr: true,
		},
		{
			name: "URL with spaces",
			input: "   https://github.com   ",
			wantErr: false,
		},
		{
			name: "URL with no scheme",
			input: "github.com",
			wantErr: true,
		},
		{
			name: "URL with invalid scheme",
			input: "ftp://github.com",
			wantErr: true,
		},
		{
			name: "Too long URL",
			input: "https://"+strings.Repeat("a",2001),
			wantErr: true,
		},
	}

	for _, tc := range tests{
		t.Run(tc.name, func(t *testing.T) {
			_,err := ValidateURL(tc.input)
			if (err != nil) != tc.wantErr{
				t.Errorf("ValidateURL() error = %v, want error %v", err, tc.wantErr)
			}
		})
	}

}
