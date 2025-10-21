package tg_bot

import "Step_game/heroes"

type Scenarios struct {
	scenarioName string
	firstStep    int
	steps        map[int]ScenarioStep
}

type ScenarioStep struct {
	text     string
	handler  string
	action   string
	nextStep int
	prevStep int
	buttons  []string
}

// Инициализация сценариев
func newScenario(name string, firstStep int, steps map[int]ScenarioStep) Scenarios {
	return Scenarios{
		scenarioName: name,
		firstStep:    firstStep,
		steps:        steps,
	}
}

// Создание сценариев
var (
	newGame = newScenario("newGame", 1, map[int]ScenarioStep{
		1: {
			text:     "Выберите персонажа",
			handler:  "handleChooseHero",
			nextStep: 2,
			prevStep: 1,
			buttons:  heroes.HeroList,
		},
		2: {
			text:     "",
			action:   "handleHistoryHero",
			nextStep: 3,
			prevStep: 1,
			buttons:  []string{"Далее"},
		},
		3: {
			text:     "step 3 coming soon",
			nextStep: 4,
			prevStep: 2,
		},
		4: {
			text: "step 4 coming soon",
		},
	})
	unitAttackScenario = newScenario("unitAttackScenario", 1, map[int]ScenarioStep{
		1: {
			text:     "",
			handler:  "handleAttack",
			nextStep: 2,
			prevStep: 1,
			buttons:  []string{"Нанести удар", "Уклониться"},
		},
		2: {
			text:     "",
			handler:  "",
			nextStep: 3,
			prevStep: 1,
		},
		3: {
			action: "",
		},
	})

	mapScenarios = map[string]Scenarios{
		"attackUnitScenario": unitAttackScenario,
		"newGame":            newGame,
	}
)
