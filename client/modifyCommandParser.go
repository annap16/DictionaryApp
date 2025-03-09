package main

import(
	"strings"
    "errors"
)

type ModifyCommandParams struct {
    Action     string // "add", "delete"
    TargetType string // "translation", "example"
    Word       string
    Translation string 
    Examples   []string
}

func ParseModifyCommand(command string) (*ModifyCommandParams, error) {
    parts := strings.Split(command, " ")
    if len(parts) < 3 {
        return nil, errors.New("Wprowadzono niepoprawne polecenie")
    }

    params := &ModifyCommandParams{
        Action:     strings.ToLower(parts[1]),
        TargetType: strings.ToLower(parts[2]),
    }

    switch params.Action {
    case "add":
        if params.TargetType == "translation" && CheckAddTranslationSyntax(command) {
            params.Word = parts[3]
            params.Translation = parts[4]
            params.Examples = ParseQuery(command)
        } else if params.TargetType == "example" && CheckAddExampleSyntax(command) {
            params.Word = parts[3]
			params.Translation = parts[4]
            params.Examples = ParseQuery(command)
        } else {
            return nil, errors.New("Niepoprawna składnia dla polecenia dodawania")
        }
    case "delete":
        if params.TargetType == "translation" && len(parts) == 5 {
            params.Word = parts[3]
			params.Translation = parts[4]
        } else if params.TargetType == "example" && len(parts) > 5 {
            params.Word = parts[3]
			params.Translation = parts[4]
            params.Examples = ParseQuery(command)
        } else {
            return nil, errors.New("Niepoprawna składnia dla polecenia usuwania")
        }
    default:
        return nil, errors.New("Niepoprawna składnia dla polecenia modyfikacji słowa")
    }

    return params, nil
}
