package mocks

import (
	"context"
)

type MockTxManager struct{}

func (f *MockTxManager) Transaction(ctx context.Context, fn func(ctx context.Context) error) error {
	return fn(ctx)
}
