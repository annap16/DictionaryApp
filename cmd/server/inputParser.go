package main

import(
	"fmt"
	"os"
	"bufio"
	"strings"
)

func WaitForUserInput() {
	reader := bufio.NewReader(os.Stdin)
	exit := false

	createHandler := &CreateCommandHandler{}
	reciveHandler := &ReciveCommandHandler{}
	modifyHandler := &ModifyCommandHandler{}
	removeHandler := &RemoveCommandHandler{}

	createHandler.SetNext(reciveHandler)
	reciveHandler.SetNext(modifyHandler)
	modifyHandler.SetNext(removeHandler)

	for !exit {
		fmt.Println("Enter command:")
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error while reading input data", err)
			return
		}
		input = input[:len(input)-1]
		input = strings.ToLower(input)

		if(input=="exit"){
			exit = true
			break
		}
		if !createHandler.HandleCommand(input){
			fmt.Println("Unexpected Command")
		}
	}	
	
}