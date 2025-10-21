package repository

import (
	"Step_game/domain"
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

type SQLXRepository struct {
	db     *sqlx.DB
	logger *zap.Logger
}

func NewSQLXRepository(db *sqlx.DB, logger *zap.Logger) *SQLXRepository {
	return &SQLXRepository{
		db:     db,
		logger: logger.Named("repository"),
	}
}

// GetByID возвращает сущность по ID
func (r *SQLXRepository) GetByID(ctx context.Context, id int64, entity domain.Entity) error {
	const op = "SQLXRepository.GetByID"

	if err := entity.Validate(); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	query := fmt.Sprintf("SELECT * FROM %s WHERE chat_id = ?", entity.TableName())

	// sqlx для маппинга полей
	if err := r.db.GetContext(ctx, entity, query, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.ErrNotFound
		}
		return fmt.Errorf("%s: %w", op, err)
	}

	// Обработка JSON полей для UserState
	if userState, ok := entity.(*domain.UserState); ok {
		return r.scanUserStateContext(ctx, userState)
	}

	return nil
}

// GetAll - возвращает все сущности указанного типа
func (r *SQLXRepository) GetAll(ctx context.Context, entityType domain.Entity) ([]domain.Entity, error) {
	const op = "SQLXRepository.GetAll"

	query := fmt.Sprintf("SELECT * FROM %s", entityType.TableName())

	switch entityType.(type) {
	case domain.UserState:
		var states []domain.UserState
		if err := r.db.SelectContext(ctx, &states, query); err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}

		entities := make([]domain.Entity, len(states))
		for i := range states {
			if err := r.scanUserStateContext(ctx, &states[i]); err != nil {
				return nil, fmt.Errorf("%s: %w", op, err)
			}
			entities[i] = &states[i]
		}
		return entities, nil

	case domain.Request:
		var requests []domain.Request
		if err := r.db.SelectContext(ctx, &requests, query); err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}

		entities := make([]domain.Entity, len(requests))
		for i := range requests {
			entities[i] = &requests[i]
		}
		return entities, nil

	default:
		return nil, fmt.Errorf("%s: %w: %T", op, domain.ErrInvalidEntity, entityType)
	}
}

// Create - создает новую сущность
func (r *SQLXRepository) Create(ctx context.Context, entity domain.Entity) error {
	const op = "SQLXRepository.Create"

	if err := entity.Validate(); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	defer tx.Rollback()

	switch e := entity.(type) {
	case domain.UserState:
		if err := r.createUserState(ctx, tx, &e); err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}
	case domain.Request:
		if err := r.createRequest(ctx, tx, &e); err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}
	default:
		return fmt.Errorf("%s: %w: %T", op, domain.ErrInvalidEntity, entity)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

// Update - обновляет существующую сущность
func (r *SQLXRepository) Update(ctx context.Context, entity domain.Entity) error {
	const op = "SQLXRepository.Update"

	if err := entity.Validate(); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	defer tx.Rollback()

	switch e := entity.(type) {
	case domain.UserState:
		if err := r.updateUserState(ctx, tx, &e); err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}
	default:
		return fmt.Errorf("%s: %w: %T", op, domain.ErrInvalidEntity, entity)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

// Delete - удаляет сущность
func (r *SQLXRepository) Delete(ctx context.Context, id int64, entityType domain.Entity) error {
	const op = "SQLXRepository.Delete"

	if id <= 0 {
		return fmt.Errorf("%s: %w", op, domain.ErrInvalidChatID)
	}

	query := fmt.Sprintf("DELETE FROM %s WHERE chat_id = ?", entityType.TableName())

	result, err := r.db.ExecContext(ctx, query, id)
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

// Вспомогательные методы для UserState
func (r *SQLXRepository) createUserState(ctx context.Context, tx *sqlx.Tx, state *domain.UserState) error {
	contextJSON, err := json.Marshal(state.Context)
	if err != nil {
		return fmt.Errorf("marshal context: %w", err)
	}

	query := `INSERT INTO userstate (chat_id, user_name, scenario_name, step_name, context) 
              VALUES (?, ?, ?, ?, ?)`

	_, err = tx.ExecContext(ctx, query,
		state.ChatID, state.UserName, state.ScenarioName, state.StepName, contextJSON)
	return err
}

func (r *SQLXRepository) updateUserState(ctx context.Context, tx *sqlx.Tx, state *domain.UserState) error {
	contextJSON, err := json.Marshal(state.Context)
	if err != nil {
		return fmt.Errorf("marshal context: %w", err)
	}

	query := `UPDATE userstate SET user_name = ?, scenario_name = ?, step_name = ?, context = ? 
              WHERE chat_id = ?`

	result, err := tx.ExecContext(ctx, query,
		state.UserName, state.ScenarioName, state.StepName, contextJSON, state.ChatID)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return domain.ErrNotFound
	}

	return nil
}

func (r *SQLXRepository) createRequest(ctx context.Context, tx *sqlx.Tx, req *domain.Request) error {
	query := `INSERT INTO requests (date, user_name, operation, result) VALUES (?, ?, ?, ?)`
	_, err := tx.ExecContext(ctx, query, req.Date, req.UserName, req.Operation, req.Result)
	return err
}

func (r *SQLXRepository) scanUserStateContext(ctx context.Context, state *domain.UserState) error {
	var contextJSON []byte
	var contextMap map[string]interface{}

	// Получаем JSON контекст отдельным запросом
	query := "SELECT context FROM userstate WHERE chat_id = ?"
	if err := r.db.GetContext(ctx, &contextJSON, query, state.ChatID); err != nil {
		return fmt.Errorf("get context: %w", err)
	}

	if len(contextJSON) > 0 {
		if err := json.Unmarshal(contextJSON, &contextMap); err != nil {
			return fmt.Errorf("unmarshal context: %w", err)
		}
		state.Context = contextMap
	}

	return nil
}
