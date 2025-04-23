package main

import (
	"bytes"
	"context"
	"dictionary-app/server/graph/model"
	"github.com/machinebox/graphql"
	"github.com/stretchr/testify/mock"
	"os"
)

type MockGraphQLClient struct {
	mock.Mock
}

func (m *MockGraphQLClient) Run(ctx context.Context, req *graphql.Request, res interface{}) error {
	args := m.Called(ctx, req, res)
	return args.Error(0)
}

type MockQueriesHandler struct {
	client *MockGraphQLClient
	mock.Mock
}

func (m *MockQueriesHandler) SendCreateMutation(ctx context.Context, input model.FullRecordInput) (bool, error) {
	args := m.Called(ctx, input)
	return args.Bool(0), args.Error(1)
}

func (m *MockQueriesHandler) SendReceiveMutation(ctx context.Context, input string) (string, error) {
	args := m.Called(ctx, input)
	return args.String(0), args.Error(1)
}

func (m *MockQueriesHandler) SendRemoveMutation(ctx context.Context, input string) (bool, error) {
	args := m.Called(ctx, input)
	return args.Bool(0), args.Error(1)
}

func (m *MockQueriesHandler) SendRemoveTranslationMutation(ctx context.Context, word string, input string) (bool, error) {
	args := m.Called(ctx, word, input)
	return args.Bool(0), args.Error(1)
}

func (m *MockQueriesHandler) SendRemoveExampleMutation(ctx context.Context, input model.FullRecordInput) (bool, error) {
	args := m.Called(ctx, input)
	return args.Bool(0), args.Error(1)
}

func (m *MockQueriesHandler) SendAddTranslationMutation(ctx context.Context, input model.FullRecordInput) (bool, error) {
	args := m.Called(ctx, input)
	return args.Bool(0), args.Error(1)
}

func (m *MockQueriesHandler) SendAddExampleMutation(ctx context.Context, input model.FullRecordInput) (bool, error) {
	args := m.Called(ctx, input)
	return args.Bool(0), args.Error(1)
}

func captureStdout(f func()) string {
	originalStdout := os.Stdout
	readPipe, writePipe, _ := os.Pipe()
	defer readPipe.Close()

	os.Stdout = writePipe
	defer func() { os.Stdout = originalStdout }()

	f()

	writePipe.Close()
	var buf bytes.Buffer
	_, _ = buf.ReadFrom(readPipe)
	return buf.String()
}
