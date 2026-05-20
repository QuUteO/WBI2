package main

import (
	"reflect"
	"testing"
)

func TestGrep(t *testing.T) {
	tests := []struct {
		name   string
		lines  []string
		config Config
		want   []string
	}{
		{
			name:  "фиксированная строка",
			lines: []string{"test", "testing", "tested"},
			config: Config{
				pattern: "test",
				fixed:   true,
			},
			want: []string{"test"},
		},
		{
			name:  "контекст вокруг",
			lines: []string{"line1", "line2", "match", "line4", "line5"},
			config: Config{
				pattern: "match",
				context: 1,
			},
			want: []string{"line2", "match", "line4"},
		},
		{
			name:  "контекст с разделителями",
			lines: []string{"match1", "line2", "line3", "match4", "line5"},
			config: Config{
				pattern: "match",
				after:   1,
				before:  1,
			},
			want: []string{"match1", "line2", "--", "line3", "match4", "line5"},
		},
		{
			name:  "перекрывающийся контекст",
			lines: []string{"line1", "match1", "line3", "match2", "line5"},
			config: Config{
				pattern: "match",
				context: 1,
			},
			want: []string{"line1", "match1", "line3", "match2", "line5"},
		},
		{
			name:  "несколько совпадений с разделителями",
			lines: []string{"a1", "a2", "match1", "a4", "a5", "match2", "a7", "a8"},
			config: Config{
				pattern: "match",
				context: 1,
			},
			want: []string{"a2", "match1", "a4", "--", "a5", "match2", "a7"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := grep(tt.lines, tt.config)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("grep() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGrepAdditional(t *testing.T) {
	// Тест на игнорирование регистра
	t.Run("игнорирование регистра", func(t *testing.T) {
		lines := []string{"Apple", "BANANA", "Cherry"}
		config := Config{
			pattern:    "apple",
			ignoreCase: true,
		}
		got := grep(lines, config)
		want := []string{"Apple"}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %v, want %v", got, want)
		}
	})

	// Тест с номерами строк
	t.Run("номера строк", func(t *testing.T) {
		lines := []string{"apple", "banana", "cherry"}
		config := Config{
			pattern: "a",
			lineNum: true,
		}
		got := grep(lines, config)
		want := []string{"1:apple", "2:banana"}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %v, want %v", got, want)
		}
	})

	// Тест на количество
	t.Run("количество", func(t *testing.T) {
		lines := []string{"apple", "banana", "cherry", "apricot"}
		config := Config{
			pattern: "a",
			count:   true,
		}
		got := grep(lines, config)
		want := []string{"3"}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %v, want %v", got, want)
		}
	})
}
