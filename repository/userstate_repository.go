package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"Step_game/domain"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

type UserStateRepo struct {
	db     *sqlx.DB
	repo   Repository
	logger *zap.Logger
}

func NewUserStateRepository(db *sqlx.DB, logger *zap.Logger) *UserStateRepo {
	return &UserStateRepo{
		db:     db,
		logger: logger.Named("userstate_repo"),
	}
}

// GetByChatID возвращает UserState по chat_id
func (r *UserStateRepo) GetByChatID(ctx context.Context, chatID int64) (*domain.UserState, error) {
	const op = "UserStateRepo.GetByChatID"

	if chatID <= 0 {
		return nil, fmt.Errorf("%s: %w", op, domain.ErrInvalidChatID)
	}

	var state domain.UserState
	query := "SELECT chat_id, user_name, scenario_name, step_name FROM userstate WHERE chat_id = ?"

	if err := r.db.GetContext(ctx, &state, query, chatID); err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	// Загружаем контекст
	if err := r.loadContext(ctx, &state); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &state, nil
}

// UpdateStepAndContext обновляет только step_name и context
func (r *UserStateRepo) UpdateStepAndContext(
	ctx context.Context, chatID int64, stepName int, context map[string]interface{},
) error {

	const op = "UserStateRepo.UpdateStepAndContext"

	contextJSON, err := json.Marshal(context)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	query := "UPDATE userstate SET step_name = ?, context = ? WHERE chat_id = ?"
	result, err := r.db.ExecContext(ctx, query, stepName, contextJSON, chatID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if rows == 0 {
		return domain.ErrNotFound
	}

	return nil
}

func (r *UserStateRepo) loadContext(ctx context.Context, state *domain.UserState) error {
	var contextJSON []byte
	query := "SELECT context FROM userstate WHERE chat_id = ?"

	if err := r.db.GetContext(ctx, &contextJSON, query, state.ChatID); err != nil {
		return fmt.Errorf("load context: %w", err)
	}

	if len(contextJSON) > 0 {
		var contextMap map[string]interface{}
		if err := json.Unmarshal(contextJSON, &contextMap); err != nil {
			return fmt.Errorf("unmarshal context: %w", err)
		}
		state.Context = contextMap
	}

	return nil
}
