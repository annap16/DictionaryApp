package main

import ()

type CommandFactory interface {
	CreateAction(handler QueriesHandler, command string) (ModifyAction, error)
}

type ModifyCommandFactory struct {
}

func (m *ModifyCommandFactory) CreateAction(handler QueriesHandler, command string) (ModifyAction, error) {
	params, err := ParseModifyCommand(command)
	if err != nil {
		return nil, err
	}

	switch params.Action {
	case "dodaj":
		return &ModifyAddCommand{
			handler:     handler,
			targetType:  params.TargetType,
			word:        params.Word,
			translation: params.Translation,
			examples:    params.Examples,
		}, nil
	case "usuń":
		return &ModifyDeleteCommand{
			handler:     handler,
			targetType:  params.TargetType,
			word:        params.Word,
			translation: params.Translation,
			examples:    params.Examples,
		}, nil
	}

	return nil, nil
}
