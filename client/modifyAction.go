package main

import(
	"context"
	"dictionary-app/server/graph/model"
)

type ModifyAction interface {
    Execute() (bool, error)
}


type ModifyAddCommand struct {
    handler QueriesHandler
    targetType string
    word string
    translation string
    examples []string
}

func (m *ModifyAddCommand) Execute() (bool, error) {
    var success bool
    var err error

	input := model.FullRecordInput{
		Word:        m.word,
		Translation: m.translation,
		Examples:    m.examples,
	}

    switch m.targetType {
    case "tłumaczenie":
        success, err = m.handler.SendAddTranslationMutation(context.Background(), input)
    case "przykład":
        success, err = m.handler.SendAddExampleMutation(context.Background(), input)
    default:
        return false, nil
    }
    if err != nil {
        return false, err
    }

    return success, nil
}

type ModifyDeleteCommand struct {
    handler QueriesHandler
    targetType string
    word string
	translation string
    examples []string
}

func (m *ModifyDeleteCommand) Execute() (bool, error) {
    var success bool
    var err error

    switch m.targetType {
    case "tłumaczenie":
        success, err = m.handler.SendRemoveTranslationMutation(context.Background(), m.word, m.translation)
    case "przykład":
        input := model.FullRecordInput{
            Word:        m.word,
            Translation: m.translation,
            Examples:    m.examples,
        }
        success, err = m.handler.SendRemoveExampleMutation(context.Background(), input)
    default:
        return false, nil
    }

    if err != nil {
        return false, err
    }
	
    return success, nil
}
