package unpack

import (
	"testing"
)

func TestUnpack(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
		wantErr  bool
	}{
		{
			name:     "simple unpack",
			input:    "a4bc2d5e",
			expected: "aaaabccddddde",
			wantErr:  false,
		},
		{
			name:     "no digits",
			input:    "abcd",
			expected: "abcd",
			wantErr:  false,
		},
		{
			name:     "empty string",
			input:    "",
			expected: "",
			wantErr:  false,
		},
		{
			name:     "single character",
			input:    "a",
			expected: "a",
			wantErr:  false,
		},
		{
			name:     "single repetition",
			input:    "a5",
			expected: "aaaaa",
			wantErr:  false,
		},
		{
			name:     "multi-digit number",
			input:    "a10",
			expected: "aaaaaaaaaa",
			wantErr:  false,
		},
		{
			name:     "starts with digit",
			input:    "4abc",
			expected: "",
			wantErr:  true,
		},
		{
			name:     "only digits",
			input:    "45",
			expected: "",
			wantErr:  true,
		},
		{
			name:     "digit without preceding char in middle",
			input:    "ab2c3",
			expected: "abbccc",
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := Unpack(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("Unpack() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if result != tt.expected {
				t.Errorf("Unpack() = %q, want %q", result, tt.expected)
			}
		})
	}
}

func TestUnpackWithEscape(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
		wantErr  bool
	}{
		{
			name:     "escaped digits simple",
			input:    `qwe\4\5`,
			expected: "qwe45",
			wantErr:  false,
		},
		{
			name:     "escaped digit with repetition",
			input:    `qwe\45`,
			expected: "qwe44444",
			wantErr:  false,
		},
		{
			name:     "escape at end of string",
			input:    `abc\`,
			expected: "",
			wantErr:  true,
		},
		{
			name:     "complex escapes",
			input:    `a\2b\3\4c`,
			expected: "a2b34c",
			wantErr:  false,
		},
		{
			name:     "escaped letter with repetition",
			input:    `\a2`,
			expected: "aa",
			wantErr:  false,
		},
		{
			name:     "multiple escapes and digits",
			input:    `\4\5\6`,
			expected: "456",
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := Unpack(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("Unpack() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if result != tt.expected {
				t.Errorf("Unpack() = %q, want %q", result, tt.expected)
			}
		})
	}
}
