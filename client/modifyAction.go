package main

import(
	"fmt"
	"context"
	"log"
	"github.com/annap16/DictionaryApp/graph/model"
)

type ModifyAction interface {
    Execute() bool
}


type ModifyAddCommand struct {
    handler QueriesHandler
    targetType string
    word string
    translation string
    examples []string
}

// TODO better error handling and returning err not only bool

func (m *ModifyAddCommand) Execute() bool {
    var success bool
    var err error

	input := model.FullRecordInput{
		Word:        m.word,
		Translation: m.translation,
		Examples:    m.examples,
	}

    switch m.targetType {
    case "translation":
        success, err = m.handler.SendAddTranslationMutation(context.Background(), input)
    case "example":
        success, err = m.handler.SendAddExampleMutation(context.Background(), input)
    default:
        fmt.Println("Invalid modify add command")
        return false
    }

    if err != nil {
        log.Println("Error:", err)
        return false
    }

    return success
}

type ModifyDeleteCommand struct {
    handler QueriesHandler
    targetType string
    word string
	translation string
    examples []string
}

func (m *ModifyDeleteCommand) Execute() bool {
    var success bool
    var err error

    switch m.targetType {
    case "translation":
        success, err = m.handler.SendRemoveTranslationMutation(context.Background(), m.word, m.translation)
    case "example":
        input := model.FullRecordInput{
            Word:        m.word,
            Translation: m.translation,
            Examples:    m.examples,
        }
        success, err = m.handler.SendRemoveExampleMutation(context.Background(), input)
    default:
        fmt.Println("Invalid modify delete command")
        return false
    }

    if err != nil {
        log.Println("Error:", err)
        return false
    }
	
    return success
}
