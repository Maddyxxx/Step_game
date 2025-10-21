package repository

import (
	"Step_game/domain"
	"context"
)

// Repository определяет контракт для работы с данными
type Repository interface {
	// GetByID возвращает сущность по ID
	GetByID(ctx context.Context, id int64, entity domain.Entity) error

	// GetAll возвращает все сущности указанного типа
	GetAll(ctx context.Context, entityType domain.Entity) ([]domain.Entity, error)

	// Create создает новую сущность
	Create(ctx context.Context, entity domain.Entity) error

	// Update обновляет существующую сущность
	Update(ctx context.Context, entity domain.Entity) error

	// Delete удаляет сущность
	Delete(ctx context.Context, id int64, entityType domain.Entity) error
}

// UserStateRepository специализированный репозиторий для UserState
type UserStateRepository interface {
	Repository
	GetByChatID(ctx context.Context, chatID int64) (*domain.UserState, error)
	UpdateStepAndContext(
		ctx context.Context, chatID int64, stepName int, context map[string]interface{}) error
}
