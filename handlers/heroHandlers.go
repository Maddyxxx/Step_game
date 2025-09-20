package handlers

// handleOperationType - обрабатывает введенный тип операции и возвращает результат
func handleOperationType(context map[string]interface{}) string {
	context["checkType"] = context["textMsg"]
	return "next"
}

func handleChoseHero(context map[string]interface{}) string {
	hero := context["hero"].(string)
	context["equipment"] = hero
	context["equipment"] = hero
	return "next"
}

func handleHistoryHero(context map[string]interface{}) string {
	return "next"
}

func handleMakeAttack(context map[string]interface{}) string {
	// todo здесь будет логика боя
	equipment := context["equipment"].([]string)
	if equipment != nil {
		battle()
		return "next"
	}
	return ""
}

func battle() string {
	// todo во время боя здоровье героя отнимается, прибавляется опыт
	return "result"
}
