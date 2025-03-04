package main

import(
	"github.com/machinebox/graphql"
	)

func main() {
	client := graphql.NewClient("http://localhost:8080/query")
	queriesHandlerQL := &QueriesHandlerQL {client: client} 
	WaitForUserInput(queriesHandlerQL)
}



