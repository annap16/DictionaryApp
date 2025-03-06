package main

import(
	"regexp"

)

func CheckCreateSyntax(command string) bool{
	pattern := `(?i)^create\s+([^\[\]\s]+)\s+([^\[\]\s]+)(\s+\[[^\[\]]+\])*?$`
	re := regexp.MustCompile(pattern)

	return re.MatchString(command)
}

func CheckAddExampleSyntax(command string) bool{
	pattern := `(?i)^modify\s+add\s+example\s+(\S+)\s+(\S+)(\s+\[.*?\])+`

	re := regexp.MustCompile(pattern)

	return re.MatchString(command)
}

func CheckAddTranslationSyntax(command string) bool{
	pattern := `(?i)^modify\s+add\s+translation\s+(\S+)\s+(\S+)(\s+\[.*?\])*$`
	re := regexp.MustCompile(pattern)

	return re.MatchString(command)
}

func ParseQuery(query string) ([]string) {
	re := regexp.MustCompile(`\[(.*?)\]`)
	matches := re.FindAllStringSubmatch(query, -1)

	var sentences []string
	for _, match := range matches {
		sentences = append(sentences, match[1]) 
	}

	return sentences
}
