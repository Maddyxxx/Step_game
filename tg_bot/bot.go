package tg_bot

import (
	"Step_game/repository"
	"context"
	"fmt"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	conf "Step_game/config"
	"Step_game/domain"
)

// Config конфигурация бота
type Config struct {
	Token string
	Debug bool
}

type Bot struct {
	api     *tgbotapi.BotAPI
	state   *domain.UserState
	repo    *repository.SQLXRepository
	usRepo  *repository.UserStateRepo
	logger  *zap.Logger
	ctx     context.Context
	cancel  context.CancelFunc
	context map[string]interface{}
}

// Run - запуск бота
func (b *Bot) Run() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60 // todo REST
	for update := range b.api.GetUpdatesChan(u) {
		if msg := update.Message; msg != nil {
			b.handleMessage(msg)
		}
	}
}

// handleMessage - обработка входящих сообщений
func (b *Bot) handleMessage(msg *tgbotapi.Message) {
	b.state.ChatID = msg.Chat.ID
	b.logger.Info("Получено новое сообщение: ", zap.String("Message", msg.Text))
	intent := conf.GetIntent(strings.ToLower(msg.Text), b.logger)
	err := b.repo.GetByID(b.ctx, b.state.ChatID, b.state) // попытка получить состояние пользователя
	fmt.Println(err)
	if intent != nil {
		b.processIntent(intent, msg, b.ctx)
	} else if err == nil {
		// если пользователь уже есть - продолжение сценария
		b.context["textMsg"] = strings.ToLower(msg.Text)
		scenario := mapScenarios[b.state.ScenarioName]
		scenario.ContinueScenario(b)
	} else {
		b.logger.Info("There is no message in intents for textMsg: %s", zap.String(
			"Message", msg.Text), zap.Error(err))
		b.sendMsg(conf.DefaultAnswer, nil)
	}
}

// processIntent - обработка намерения
func (b *Bot) processIntent(intent *conf.Intent, msg *tgbotapi.Message, ctx context.Context) {
	scenario := intent.Scenario
	if scenario != "" {
		scenario := mapScenarios[scenario]
		b.state.UserName = msg.From.UserName
		//todo b.context["operationType"] = conf.OperTypes[msg.Text]
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
	b.api.Send(Msg)
	b.logger.Info("Sending message: ", zap.String("Message", Msg.Text))
}

// todo доработать
//func (b *Bot) saveRequest(result string) {
//	date := time.Now()
//	req := mod.Request{
//		Date:      date.Format("02.01.2006 15:04:05"),
//		UserName:  b.state.UserName,
//		Operation: b.state.ScenarioName,
//		Result:    result,
//	}
//	req.InsertData(b.db)
//}

// Dependencies зависимости для создания бота
type Dependencies struct {
	DB     *sqlx.DB
	Logger *zap.Logger
}

// NewBot - создает новый экземпляр бота
func NewBot(cfg Config, deps Dependencies) (*Bot, error) {
	const op = "bot.NewBot"

	// Инициализация Telegram API
	botAPI, err := tgbotapi.NewBotAPI(cfg.Token)
	if err != nil {
		return nil, fmt.Errorf("%s: create bot API: %w", op, err)
	}

	botAPI.Debug = cfg.Debug

	// Создание репозиториев
	repo := repository.NewSQLXRepository(deps.DB, deps.Logger)

	// Создание контекста с возможностью отмены
	ctx, cancel := context.WithCancel(context.Background())

	bot := &Bot{
		api:     botAPI,
		state:   &domain.UserState{},
		repo:    repo,
		logger:  deps.Logger,
		ctx:     ctx,
		cancel:  cancel,
		context: make(map[string]interface{}),
	}

	deps.Logger.Info("Bot authorized successfully",
		zap.String("username", botAPI.Self.UserName),
		zap.Bool("debug", cfg.Debug),
	)

	return bot, nil
}

// Close - освобождает ресурсы бота
func (b *Bot) Close() error {
	b.logger.Info("Shutting down bot...")

	// Отмена контекста
	if b.cancel != nil {
		b.cancel()
	}

	return nil
}
