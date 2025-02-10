package handlers

import (
	"context"
	"net/http"
)

type HTTPHandler func(ctx context.Context, w http.ResponseWriter, r *http.Request) error
