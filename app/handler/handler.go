package handler

import (
	"go.uber.org/zap"
)

// Handler - struct used through program and passed to go-micro.
type Handler struct {
	zapLog *zap.Logger
}

// NewHandler returns a Handler object
func NewHandler(zapLog *zap.Logger) *Handler {
	return &Handler{zapLog}
}
