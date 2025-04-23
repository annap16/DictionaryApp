package graph

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.
import "dictionary-app/server/database"

type Resolver struct {
	DBInterface *database.DBInterface
}
