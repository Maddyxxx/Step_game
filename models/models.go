package models

// UserState Состояние пользователя внутри сценария
type UserState struct {
	ChatID       int64
	UserName     string
	ScenarioName string
	StepName     int
	Context      map[string]interface{}
}

// Request - Лог запросов
type Request struct {
	Date      string
	UserName  string
	Operation string
	Result    string // Result - результат прохождения сценария
}
