package main

import (
	"context"
	"dictionary-app/server/graph/model"
	"github.com/machinebox/graphql"
)

type GraphQLClient interface {
	Run(ctx context.Context, req *graphql.Request, res interface{}) error
}

type QueriesHandler interface {
	SendCreateMutation(ctx context.Context, input model.FullRecordInput) (bool, error)
	SendReceiveMutation(ctx context.Context, input string) (string, error)
	SendRemoveMutation(ctx context.Context, input string) (bool, error)
	SendRemoveTranslationMutation(ctx context.Context, word string, input string) (bool, error)
	SendRemoveExampleMutation(ctx context.Context, input model.FullRecordInput) (bool, error)
	SendAddTranslationMutation(ctx context.Context, input model.FullRecordInput) (bool, error)
	SendAddExampleMutation(ctx context.Context, input model.FullRecordInput) (bool, error)
}
