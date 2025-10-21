package tg_bot

import (
	conf "Step_game/config"
	"Step_game/domain"
	hand "Step_game/handlers"
	"Step_game/repository"
	"context"
	"fmt"
	"go.uber.org/zap"
)

func (s *Scenarios) StartScenario(b *Bot) {
	b.repo.Delete(b.ctx, b.state.ChatID, b.state)
	b.logger.Info("Starting scenario", zap.String("scenarioName", s.scenarioName))
	b.state.ScenarioName = s.scenarioName
	b.state.StepName = s.firstStep
	b.context["textToSend"] = nil

	err := b.repo.Create(b.ctx, b.state)
	if err == nil {
		b.logger.Info("Success inserting data for starting scenario")
		s.sendStep(b, s.firstStep)
	} else {
		b.logger.Error("Error inserting data ", zap.Error(err))
		b.sendMsg("Произошла ошибка начала сценария", nil)
	}
}

// ContinueScenario - продолжение сценария при наличии state
func (s *Scenarios) ContinueScenario(b *Bot) {
	b.logger.Info("Continuing scenario %s, step N %d",
		zap.String("scenarioName", s.scenarioName), zap.Int("StepName", b.state.StepName))
	step := s.steps[b.state.StepName]
	resp := hand.HandlerDo(step.handler, b.context)

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
		step.action = fmt.Sprintf("%v", b.context["action"])
		s.makeAction(b, step)
	default:
		b.sendMsg(b.context["error"], nil)
		s.finishScenario(b, "response error")
	}
}

// finishScenario - удаляет состояние пользователя после прохождения или прерывания сценария,
// сохраняет результат прохождения сценария
func (s *Scenarios) finishScenario(b *Bot, result string) {
	// b.saveRequest(result)
	err := b.repo.Delete(b.ctx, b.state.ChatID, b.state)
	if err != nil {
		b.logger.Error("Failed to delete UserState", zap.Error(err))
	}
	b.sendMsg(conf.FinishScenarioAnswer, nil)
	b.logger.Info("Scenario for ChatID", zap.String("result", result),
		zap.String("scenarioName", s.scenarioName), zap.Int64("ChatID", b.state.ChatID))
}

// sendStep - отправка следующего шага по сценарию
func (s *Scenarios) sendStep(b *Bot, stepNum int) {
	b.logger.Info("Sending step for scenario", zap.Int("stepNum", stepNum),
		zap.String("scenarioName", s.scenarioName))
	defer func(usRepo *repository.UserStateRepo, ctx context.Context, state *domain.UserState) {
		err := b.repo.Update(ctx, state)
		if err != nil {
			b.logger.Error("Failed to update UserState", zap.Error(err))
		}
	}(b.usRepo, b.ctx, b.state)
	step := s.steps[stepNum]
	buttons := b.makeButtons(step.buttons)

	if textToSend, ok := b.context["textToSend"]; ok && textToSend != nil {
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
		resp := hand.HandlerDo(step.handler, b.context)
		b.sendMsg(resp, nil)
		return
	}

	hand.HandlerDo(step.action, b.context)
	defer s.finishScenario(b, "interrupted")
	b.sendMsg(b.context["error"], nil)
}
