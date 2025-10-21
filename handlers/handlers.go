package handlers

import (
	"log"

	"github.com/jmoiron/sqlx"
)

type HandlerFunc func(map[string]interface{}) string

type HandlerFuncDb func(*sqlx.DB, map[string]interface{}) string

var handlers = map[string]HandlerFunc{
	"handleOperationType": handleOperationType,
	"handleChooseHero":    handleChooseHero,
	"handleHistoryHero":   handleHistoryHero,
	"handleMakeAttack":    handleMakeAttack,
}

// checkStep - проверяет сообщение на конкретные команды
// "stop" - завершение сценария; "back" - шаг назад
func checkStep(msg interface{}) string {
	stopIntents := []string{"/stop", "stop", "стоп"}
	if msg == "назад" {
		return "back"
	}

	for _, stopIntent := range stopIntents {
		if msg == stopIntent {
			return "stop"
		}
	}
	return ""
}

// HandlerDo - выполняет определенное действие в зависимости от handler
// при условии что шаг пройден, возвращает "next"
// возвращает результат функции checkStep, если пришла конкретная команда от пользователя
func HandlerDo(handler interface{}, context map[string]interface{}) string {
	log.Printf("HandlerDo called with handler: %v, context: %v", handler, context)

	if check := checkStep(context["textMsg"]); check != "" {
		return check
	}

	if handlerFunc, exists := handlers[handler.(string)]; exists {
		return handlerFunc(context)
	}

	return ""
}
