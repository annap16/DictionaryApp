package main

import (
	"fmt"

)

func ModifyCommandFactory(handler QueriesHandler, command string) ModifyAction {
    params, err := ParseModifyCommand(command)
    if err != nil {
        fmt.Println(err)
        return nil
    }

    switch params.Action {
    case "add":
        return &ModifyAddCommand{
            handler: handler,
            targetType: params.TargetType,
            word: params.Word,
            translation: params.Translation,
            examples: params.Examples,
        }
    case "delete":
        return &ModifyDeleteCommand{
            handler: handler,
            targetType: params.TargetType,
            word: params.Word,
			translation: params.Translation,
            examples: params.Examples, 
        }
    }

    return nil
}
