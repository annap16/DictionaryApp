package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCheckCreateSyntax(t *testing.T) {
	tests := []struct {
		name    string
		command string
		want    bool
	}{
		{
			name:    "Success - Command Without Examples",
			command: "dodaj xyz abc",
			want:    true,
		},
		{
			name:    "Success - Command With Capital Letters",
			command: "DODaj something else",
			want:    true,
		},
		{
			name:    "Success - Command With 2 Examples",
			command: "dodaj word translation [Sentence one] [Sentence two]",
			want:    true,
		},
		{
			name:    "Failure - Too Many Arguments Before Examples",
			command: "Dodaj word translation something [else]",
			want:    false,
		},
		{
			name:    "Failure - Not Enough Arguments Before Examples",
			command: "dodaj missingTranslation [value1] [value2]",
			want:    false,
		},
		{
			name:    "Failure - Empty Create Command With Whitespace Character",
			command: "dodaj ",
			want:    false,
		},
		{
			name:    "Failure - Command Different Than Create",
			command: "random command",
			want:    false,
		},
		{
			name:    "Failure - Wrong Brackets",
			command: "dodaj word translation [sentence[sentence2]]]",
			want:    false,
		},
		{
			name:    "Failure - Empty Brackets",
			command: "dodaj word translation []",
			want:    false,
		},
		{
			name:    "Failure - Correct And Wrong Brackets Combined",
			command: "dodaj word translation ][sentence]][]",
			want:    false,
		},
		{
			name:    "Failure - Correct And Empty Brackets Combined",
			command: "dodaj word [][ here] translation",
			want:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CheckCreateSyntax(tt.command)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestCheckAddExampleSyntax(t *testing.T) {
	tests := []struct {
		name    string
		command string
		want    bool
	}{
		{
			name:    "Success - One Example",
			command: "modyfikuj dodaj przykład word translation [Example one]",
			want:    true,
		},
		{
			name:    "Success - Multiple Examples",
			command: "modyfikuj dodaj przykład word translation [Example one] [Example two] [Example3]",
			want:    true,
		},
		{
			name:    "Success - Ignoring Capital Letters",
			command: "MODYfikuj DODaj przykład word translation [Example one]",
			want:    true,
		},
		{
			name:    "Failure - No example",
			command: "modyfikuj dodaj przykład word translation",
			want:    false,
		},
		{
			name:    "Failure - Modify Not Specified",
			command: "add example word translation [Example one]",
			want:    false,
		},
		{
			name:    "Failure - Example Not Specified",
			command: "modify translation word translation [Example one]",
			want:    false,
		},
		{
			name:    "Failure - Missing Word In Command",
			command: "modify example translation [Example one]",
			want:    false,
		},
		{
			name:    "Failure - Missing Translation In Command",
			command: "modify example word [Example one]",
			want:    false,
		},
		{
			name:    "Failure - Wrong Brackets",
			command: "modyfikuj dodaj przykład word translation [whf[fwuef]sfhqwiu dfhweu ]][]",
			want:    false,
		},
		{
			name:    "Failure - Empty Brackets",
			command: "modyfikuj dodaj przykład word translation []",
			want:    false,
		},
		{
			name:    "Failure - Correct And Wrong Brackets Combined",
			command: "modyfikuj dodaj przykład word translation [Okey okey][fweufhy[fwygfuw]]",
			want:    false,
		},
		{
			name:    "Failure - Correct And Empty Brackets Combined",
			command: "modyfikuj dodaj przykład word translation [Here] []",
			want:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CheckAddExampleSyntax(tt.command)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestCheckAddTranslationSyntax(t *testing.T) {
	tests := []struct {
		name    string
		command string
		want    bool
	}{
		{
			name:    "Success - Translation Without Examples",
			command: "modyfikuj dodaj tłumaczenie word newTranslation",
			want:    true,
		},
		{
			name:    "Success - Translation With One Example",
			command: "modyfikuj dodaj tłumaczenie word newTrans [Example one]",
			want:    true,
		},
		{
			name:    "Success - Translation With Many Examples",
			command: "modyfikuj dodaj tłumaczenie word newTrans [Example one] [Example two] [Example3]",
			want:    true,
		},
		{
			name:    "Success - Ignoring Capital Letters",
			command: "MODYfikuj DODaj TŁUMaczenie word newTrans [Example one]",
			want:    true,
		},
		{
			name:    "Failure - Missing Modify In Command",
			command: "add translation word newTrans [Example one]",
			want:    false,
		},
		{
			name:    "Failure - Missing Add In Command",
			command: "modify translation word newTrans [Example one]",
			want:    false,
		},
		{
			name:    "Failure - Missing Translation In Command",
			command: "modify add word newTrans [Example one]",
			want:    false,
		},
		{
			name:    "Failure - Wrong Brackets",
			command: "modyfikuj dodaj tłumaczenie word newTrans [huwff[weiufh]]]",
			want:    false,
		},
		{
			name:    "Failure - Empty Brackets",
			command: "modyfikuj dodaj tłumaczenie word newTrans []",
			want:    false,
		},
		{
			name:    "Failure - Correct And Wrong Brackets Combined",
			command: "modyfikuj dodaj tłumaczenie word newTrans [cwyfgu[ew]]][ [Correct]",
			want:    false,
		},
		{
			name:    "Failure - Correct And Empty Brackets Combined",
			command: "modyfikuj dodaj tłumaczenie word newTrans [Correct] []",
			want:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CheckAddTranslationSyntax(tt.command)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestParseQuery(t *testing.T) {
	tests := []struct {
		name  string
		query string
		want  []string
	}{
		{
			name:  "Success - Create Query With 2 Examples",
			query: "dodaj word translation [value1] [value2 value3]",
			want:  []string{"value1", "value2 value3"},
		},
		{
			name:  "Success - Create Query Without Any Examples",
			query: "dodaj word translation",
			want:  []string{},
		},
		{
			name:  "Success - Create Query With 1 Example",
			query: "dodaj word translation [single]",
			want:  []string{"single"},
		},
		{
			name:  "Success - Create Query 3 Examples",
			query: "dodaj word translation [one] [two] [three]",
			want:  []string{"one", "two", "three"},
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
