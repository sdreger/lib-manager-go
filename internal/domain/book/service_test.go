package book

import (
	"context"
	"github.com/stretchr/testify/assert"
	"log/slog"
	"os"
	"testing"
)

func TestService_GetById(t *testing.T) {
	ctx := context.Background()
	service := GetService()

	mockStore := NewMockStore(t)
	mockStore.EXPECT().GetByID(ctx, bookID).Return(getTestBook(), nil).Once()
	injectMocks(service, mockStore)

	book, err := service.GetBookByID(ctx, bookID)
	if assert.NoError(t, err, "should get book by id") {
		assert.Equal(t, getTestBook(), book, "books should be equal")
	}
}

func GetService() *Service {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	return NewService(logger, nil)
}

func injectMocks(service *Service, store *MockStore) {
	service.store = store
}
