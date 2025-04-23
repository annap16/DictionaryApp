package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func WaitForUserInput(handler QueriesHandler) {
	reader := bufio.NewReader(os.Stdin)
	exit := false

	createHandler := &CreateCommandHandler{handler: handler}
	reciveHandler := &ReceiveCommandHandler{handler: handler}
	modifyHandler := &ModifyCommandHandler{handler: handler, modifyFactory: &ModifyCommandFactory{}}
	removeHandler := &RemoveCommandHandler{handler: handler}

	createHandler.SetNext(reciveHandler)
	reciveHandler.SetNext(modifyHandler)
	modifyHandler.SetNext(removeHandler)

	for !exit {
		green := "\033[32m"
		reset := "\033[0m"

		fmt.Println(green + "Wprowadź polecenie:" + reset)
		input, err := reader.ReadString('\n')

		if err != nil {
			fmt.Println("Wystąpił błąd podczas wczytywania danych", err)
			return
		}

		input = strings.TrimRight(input, " \t\r\n")

		if strings.ToLower(input) == "exit" {
			exit = true
			break
		}

		if !createHandler.HandleCommand(input) {
			fmt.Println("Niepoprawne polecenie")
		}
	}

}
