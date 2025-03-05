package main

import(
	"testing"

)


func TestCheckCreateSyntax(t *testing.T) {
	tests := []struct {
		name string
		command string
		want bool
	}{
		{
			name : "Success - Command Without Examples",
			command: "create xyz abc",
			want: true,
		},
		{
			name: "Success - Command With Capital Letters", 
			command: "CREATe something else", 
			want: true,
		},
		{
			name: "Success - Command With 2 Examples",
			command: "create word translation [Sentence one] [Sentence two]", 
			want: true,
		},
		{
			name: "Failure - Too Many Arguments Before Examples",
			command: "Create word translation something [else]", 
			want: false,
		},
		{
			name: "Failure - Not Enough Arguments Before Examples",
			command: "create missingTranslation [value1] [value2]", 
			want: false,
		},
		{
			name: "Failure - Empty Create Command With Whitespace Character",
			command: "create ", 
			want: false,
		},
		{
			name: "Failure - Command Different Than Create",
			command: "random command", 
			want: false,
		},
		{
			name: "Failure - Wrong Brackets",
			command: "create word translation [sentence[sentence2]]]", 
			want: false,
		},
		{
			name: "Failure - Empty Brackets",
			command: "create word translation []", 
			want: false,
		},
		{
			name: "Failure - Correct And Wrong Brackets Combined",
			command: "create word translation ][sentence]][]", 
			want: false,
		},
		{
			name: "Failure - Correct And Empty Brackets Combined",
			command: "create word [][ here] translation", 
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CheckCreateSyntax(tt.command)
			if got != tt.want {
				t.Errorf("CheckCreateSyntax(%q) = %v, want %v", tt.command, got, tt.want)
			}
		})
	}
}

func TestParseQuery(t *testing.T) {
	tests := []struct {
		name string
		query string
		want []string
	}{
		{
			name: "Success - Create Query With 2 Examples",
			query: "create word translation [value1] [value2 value3]", 
			want: []string{"value1", "value2 value3"},
		},
		{
			name: "Success - Create Query Without Any Examples",
			query: "create word translation", 
			want: []string{},
		},
		{
			name: "Success - Create Query With 1 Example",
			query: "create word translation [single]", 
			want: []string{"single"},
		},
		{
			name: "Success - Create Query 3 Examples",
			query: "create word translation [one] [two] [three]", 
			want: []string{"one", "two", "three"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ParseQuery(tt.query)
			if !equalSlices(got, tt.want) {
				t.Errorf("ParseQuery(%q) = %v, want %v", tt.query, got, tt.want)
			}
		})
	}
}