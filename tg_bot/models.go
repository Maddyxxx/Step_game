package tg_bot

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
	}
)
