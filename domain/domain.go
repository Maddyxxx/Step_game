package domain

// Entity представляет базовый интерфейс для всех сущностей
type Entity interface {
	TableName() string
	Validate() error
}

// UserState представляет состояние пользователя
type UserState struct {
	ChatID       int64  `db:"chat_id"`
	UserName     string `db:"user_name"`
	ScenarioName string `db:"scenario_name"`
	StepName     int    `db:"step_name"`
}

func (UserState) TableName() string { return "userstate" }

func (u UserState) Validate() error {
	if u.ChatID <= 0 {
		return ErrInvalidChatID
	}
	return nil
}

// Request представляет запрос пользователя
type Request struct {
	Date      string `db:"date"`
	UserName  string `db:"user_name"`
	Operation string `db:"operation"`
	Result    string `db:"result"`
}

func (Request) TableName() string { return "requests" }

func (r Request) Validate() error {
	if r.UserName == "" || r.Operation == "" {
		return ErrInvalidRequest
	}
	return nil
}
