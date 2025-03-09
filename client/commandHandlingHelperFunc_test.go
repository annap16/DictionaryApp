package main

import(
	"testing"
	"github.com/stretchr/testify/assert"

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
			assert.Equal(t,tt.want, got)
		})
	}
}

func TestCheckAddExampleSyntax(t *testing.T){
	tests := []struct {
		name string
		command string
		want bool
	}{
		{
			name : "Success - One Example",
			command: "modify add example word translation [Example one]",
			want: true,
		},
		{
			name: "Success - Multiple Examples",
			command: "modify add example word translation [Example one] [Example two] [Example3]",
			want: true,
		},
		{
			name: "Success - Ignoring Capital Letters",
			command: "MODIfy ADd example word translation [Example one]",
			want: true,
		},
		{
			name: "Failure - No example",
			command: "modify add example word translation",
			want: false,
		},
		{
			name: "Failure - Modify Not Specified",
			command: "add example word translation [Example one]",
			want: false,
		},
		{
			name: "Failure - Example Not Specified",
			command: "modify translation word translation [Example one]",
			want: false,
		},
		{
			name: "Failure - Missing Word In Command",
			command: "modify example translation [Example one]",
			want: false,
		},
		{
			name: "Failure - Missing Translation In Command",
			command: "modify example word [Example one]",
			want: false,
		},
		{
			name: "Failure - Wrong Brackets",
			command: "modify add example word translation [whf[fwuef]sfhqwiu dfhweu ]][]",
			want: false,
		},
		{
			name: "Failure - Empty Brackets",
			command: "modify add example word translation []",
			want: false,
		},
		{
			name: "Failure - Correct And Wrong Brackets Combined",
			command: "modify add example word translation [Okey okey][fweufhy[fwygfuw]]",
			want: false,
		},
		{
			name: "Failure - Correct And Empty Brackets Combined",
			command: "modify add example word translation [Here] []",
			want: false,
		},

	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CheckAddExampleSyntax(tt.command)
			assert.Equal(t,tt.want, got)
		})
	}
}

func TestCheckAddTranslationSyntax(t *testing.T){
	tests := []struct {
		name string
		command string
		want bool
	}{
		{
			name: "Success - Translation Without Examples",
			command: "modify add translation word newTranslation",
			want: true,
		},
		{
			name: "Success - Translation With One Example",
			command: "modify add translation word newTrans [Example one]",
			want: true,
		},
		{
			name: "Success - Translation With Many Examples",
			command: "modify add translation word newTrans [Example one] [Example two] [Example3]",
			want: true,
		},
		{
			name: "Success - Ignoring Capital Letters",
			command: "MODify ADd TRANSlation word newTrans [Example one]",
			want: true,
		},
		{
			name: "Failure - Missing Modify In Command",
			command: "add translation word newTrans [Example one]",
			want: false,
		},
		{
			name: "Failure - Missing Add In Command",
			command: "modify translation word newTrans [Example one]",
			want: false,
		},
		{
			name: "Failure - Missing Translation In Command",
			command: "modify add word newTrans [Example one]",
			want: false,
		},
		{
			name: "Failure - Wrong Brackets",
			command: "modify add translation word newTrans [huwff[weiufh]]]",
			want: false,
		},
		{
			name: "Failure - Empty Brackets",
			command: "modify add translation word newTrans []",
			want: false,
		},
		{
			name: "Failure - Correct And Wrong Brackets Combined",
			command: "modify add translation word newTrans [cwyfgu[ew]]][ [Correct]",
			want: false,
		},
		{
			name: "Failure - Correct And Empty Brackets Combined",
			command: "modify add translation word newTrans [Correct] []",
			want: false,
		},

	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CheckAddTranslationSyntax(tt.command)
			assert.Equal(t,tt.want, got)
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