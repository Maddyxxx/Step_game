package tg_bot

import (
	"fmt"
	"log"

	conf "Step_game/config"
	hand "Step_game/handlers"
)

func (s *Scenarios) StartScenario(b *Bot) {
	b.state.DeleteData(b.db)
	log.Printf("Starting scenario %s", s.scenarioName)
	b.state.ScenarioName = s.scenarioName
	b.state.StepName = s.firstStep
	b.state.Context["textToSend"] = nil
	err := b.state.InsertData(b.db)
	if err == nil {
		log.Printf("Success inserting data for starting scenario")
		s.sendStep(b, s.firstStep)
	} else {
		log.Printf("Error inserting data: %v", err)
		b.sendMsg("Произошла ошибка начала сценария", nil)
	}

}

// ContinueScenario - продолжение сценария при наличии state
func (s *Scenarios) ContinueScenario(b *Bot) {
	log.Printf("Continuing scenario %s, step N %d", s.scenarioName, b.state.StepName)
	step := s.steps[b.state.StepName]
	resp := hand.HandlerDo(step.handler, b.state.Context)

	switch resp {
	case "stop":
		s.finishScenario(b, "interrupted")
	case "next":
		b.state.StepName = step.nextStep
		s.sendStep(b, step.nextStep)
	case "back":
		b.state.StepName = step.prevStep
		s.sendStep(b, step.prevStep)
	case "action":
		step.action = fmt.Sprintf("%v", b.state.Context["action"])
		s.makeAction(b, step)
	default:
		b.sendMsg(b.state.Context["error"], nil)
		s.finishScenario(b, "response error")
	}
}

// finishScenario - удаляет состояние пользователя после прохождения или прерывания сценария,
// сохраняет результат прохождения сценария
func (s *Scenarios) finishScenario(b *Bot, result string) {
	defer log.Printf("%s scenario %v for ChatID %d", result, s.scenarioName, b.state.ChatID)
	// b.saveRequest(result)
	b.state.DeleteData(b.db)
	b.sendMsg(conf.FinishScenarioAnswer, nil)
}

// sendStep - отправка следующего шага по сценарию
func (s *Scenarios) sendStep(b *Bot, stepNum int) {
	log.Printf("Sending step N %d for scenario %s", stepNum, s.scenarioName)
	defer b.state.UpdateData(b.db)
	step := s.steps[stepNum]
	buttons := b.makeButtons(step.buttons)

	if textToSend, ok := b.state.Context["textToSend"]; ok && textToSend != nil {
		b.sendMsg(textToSend, buttons)
	} else {
		if step.text != "" {
			b.sendMsg(step.text, buttons)
		} else {
			b.sendMsg("Ошибка: текст ответа не сформирован", buttons)
		}
	}

	if step.action != "" {
		s.makeAction(b, step)
	}
}

// makeAction - выполнение действия в step.Action
func (s *Scenarios) makeAction(b *Bot, step ScenarioStep) {
	if s.scenarioName == "dealUpdate" {
		defer s.finishScenario(b, "finished")
		resp := hand.HandlerDo(step.handler, b.state.Context)
		b.sendMsg(resp, nil)
		return
	}

	hand.HandlerDo(step.action, b.state.Context)
	defer s.finishScenario(b, "interrupted")
	b.sendMsg(b.state.Context["error"], nil)

}
