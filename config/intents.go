package config

import "go.uber.org/zap"

type Intent struct {
	Tokens   []string
	Scenario string
	Answer   string
	Buttons  []string
}

var helpIntent = Intent{
	Tokens: []string{"/help"},
	Answer: HelpAnswer,
}

var startIntent = Intent{
	Tokens:  []string{"/start", "старт", "главное меню"},
	Answer:  "Выберите раздел",
	Buttons: []string{"Продолжить", "Новая игра", "Правила"},
}

// game intents

var newGameIntent = Intent{
	Tokens:   []string{"новая игра"},
	Scenario: "newGame",
}

var continueGameIntent = Intent{
	Tokens:   []string{"продолжить"},
	Scenario: "continueGame",
}

var Intents = []Intent{
	helpIntent, startIntent, newGameIntent, continueGameIntent,
}

func GetIntent(token string, log *zap.Logger) *Intent {
	for _, intent := range Intents {
		for _, token_ := range intent.Tokens {
			if token == token_ {
				return &intent
			}
		}
	}
	return nil
}
