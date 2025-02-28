package main

import(
	"fmt"
	"os"
	"bufio"
	"strings"
)

func WaitForUserInput() {
	// Initializing variables and responsibility chain
	reader := bufio.NewReader(os.Stdin)
	exit := false

	createHandler := &CreateCommandHandler{}
	reciveHandler := &ReceiveCommandHandler{}
	modifyHandler := &ModifyCommandHandler{}
	removeHandler := &RemoveCommandHandler{}

	createHandler.SetNext(reciveHandler)
	reciveHandler.SetNext(modifyHandler)
	modifyHandler.SetNext(removeHandler)

	// Input handling
	for !exit {
		fmt.Println("Enter command:")
		input, err := reader.ReadString('\n')

		if err != nil {
			fmt.Println("Error while reading input data", err)
			return
		}

		// Formatting input string
		input = strings.TrimRight(input, " \t\r\n")

		// Checking for exit condition
		if(strings.ToLower(input)=="exit"){
			exit = true
			break
		}

		// Handling command based on its type
		if !createHandler.HandleCommand(input){
			fmt.Println("Unexpected Command")
		}
	}	
	
}