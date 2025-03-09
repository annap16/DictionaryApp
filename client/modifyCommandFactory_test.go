package main

import(
	"testing"
	"github.com/stretchr/testify/assert"

)

func TestCreateAction(t *testing.T) {
    mockHandler := new(MockQueriesHandler)

    tests := []struct {
        name string
        command string
        expectedAction ModifyAction
		expectedOutput string
    }{
        {
            name: "Success - Create Add Command",
            command: "modify add translation word translation",
            expectedAction: &ModifyAddCommand{
                handler: mockHandler,
                targetType: "translation",
                word: "word",
                translation: "translation",
                examples: nil,
            },
			expectedOutput: "",
        },
        {
            name: "Success - Create Delete Command",
            command: "modify delete translation word translation",
            expectedAction: &ModifyDeleteCommand{
                handler: mockHandler,
                targetType: "translation",
                word: "word",
                translation: "translation",
                examples: nil,
            },
			expectedOutput: "",
        },
		{
            name: "Failure - Wrong Key word",
            command: "modify create translation word translation",
            expectedAction: nil,
			expectedOutput: "Invalid modify action\n",
        },
		
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            factory := &ModifyCommandFactory{}

			var action ModifyAction
			output := captureStdout(func() {
				action = factory.CreateAction(mockHandler, tt.command)
			})
            

            assert.Equal(t, tt.expectedAction, action)
			assert.Equal(t, tt.expectedOutput, output)
        })
    }
}
