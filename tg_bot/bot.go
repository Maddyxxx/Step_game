package tg_bot

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	conf "Step_game/config"
	mod "Step_game/models"
)

type Bot struct {
	db    *sqlx.DB
	bot   *tgbotapi.BotAPI
	state mod.UserState
}

// Run - запуск бота
func (b *Bot) Run() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60 // todo REST

	for update := range b.bot.GetUpdatesChan(u) {
		if msg := update.Message; msg != nil {
			b.handleMessage(msg)
		}
	}
}

// handleMessage - обработка входящих сообщений
func (b *Bot) handleMessage(msg *tgbotapi.Message) {
	b.state.ChatID = msg.Chat.ID
	intent := conf.GetIntent(strings.ToLower(msg.Text))
	err := b.state.GetData(b.db) // попытка получить состояние пользователя

	if intent != nil {
		b.processIntent(intent, msg)
	} else if err == nil {
		b.state.Context["textMsg"] = strings.ToLower(msg.Text)
		scenario := mapScenarios[b.state.ScenarioName]
		scenario.ContinueScenario(b)
	} else {
		log.Printf("No message in intents for textMsg: %s", msg.Text)
		b.sendMsg(conf.DefaultAnswer, nil)
	}
}

// processIntent - обработка намерения
func (b *Bot) processIntent(intent *conf.Intent, msg *tgbotapi.Message) {
	scenario := intent.Scenario
	if scenario != "" {
		scenario := mapScenarios[scenario]
		b.state.UserName = msg.From.UserName
		//b.state.Context["operationType"] = conf.OperTypes[msg.Text]
		scenario.StartScenario(b)
	} else {
		buttons := b.makeButtons(intent.Buttons)
		b.sendMsg(intent.Answer, buttons)
	}
}

// MakeButtons - создание кнопок
func (b *Bot) makeButtons(buttons []string) *tgbotapi.ReplyKeyboardMarkup {
	if buttons == nil {
		return nil
	}

	var keyboardButtons [][]tgbotapi.KeyboardButton
	for _, button := range buttons {
		keyboardButtons = append(keyboardButtons, tgbotapi.NewKeyboardButtonRow(tgbotapi.KeyboardButton{Text: button}))
	}

	return &tgbotapi.ReplyKeyboardMarkup{
		Keyboard:        keyboardButtons,
		ResizeKeyboard:  true,
		OneTimeKeyboard: true,
	}
}

func (b *Bot) sendMsg(msg interface{}, buttons interface{}) {
	Msg := tgbotapi.NewMessage(b.state.ChatID, fmt.Sprintf("%v", msg))
	if buttons != nil {
		Msg.ReplyMarkup = buttons
	}
	b.bot.Send(Msg)
}

func (b *Bot) saveRequest(result string) {
	date := time.Now()
	req := mod.Request{
		Date:      date.Format("02.01.2006 15:04:05"),
		UserName:  b.state.UserName,
		Operation: b.state.ScenarioName,
		Result:    result,
	}
	req.InsertData(b.db)
}

func InitBot(key string, db *sqlx.DB, debug bool) *Bot {
	bot, err := tgbotapi.NewBotAPI(key)
	if err != nil {
		log.Printf("Error connecting to Telegram: %v", err)
		return nil
	}
	bot.Debug = debug
	log.Printf("Authorized on account %s", bot.Self.UserName)

	return &Bot{
		db:    db,
		bot:   bot,
		state: mod.UserState{Context: map[string]interface{}{}},
	}
}
