package handler

import (
	repository "github.com/softcorp-io/hqs_department_service/repository"
	"go.uber.org/zap"
)

// Handler - struct used through program and passed to go-micro.
type Handler struct {
	zapLog     *zap.Logger
	repository repository.Repository
}

// NewHandler returns a Handler object
func NewHandler(zapLog *zap.Logger, repository repository.Repository) *Handler {
	return &Handler{zapLog, repository}
}
