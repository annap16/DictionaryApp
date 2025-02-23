package main

import(
	"fmt"
)
func main() {
	// srv := handler.NewDefaultServer(graphql.NewExecutableSchema(graphql.Config{Resolvers: &graphql.Resolver{}}))

	//http.Handle("/graphql", srv)
	//http.Handle("/playground", playground.Handler("GraphQL playground", "/graphql"))
	fmt.Println("Hello World!")
	WaitForUserInput()
}


