package receiverbot

import (
	"film-adviser/repository"
	"film-adviser/settings"
	"fmt"
	"log"

	"github.com/mymmrac/telego"
	tu "github.com/mymmrac/telego/telegoutil"
)

type RecomendBot struct {
	bot   *telego.Bot
	token string
	repo  repository.Repository
}

func New() *RecomendBot {
	return &RecomendBot{}
}

func (rb *RecomendBot) MustInit(repo repository.Repository) {
	rb.token = settings.GetSettings().TgReceiverToken
	var err error
	rb.bot, err = telego.NewBot(rb.token)
	if err != nil {
		fmt.Println(err)
		log.Fatal("Could not start sender bot!")
	}
	rb.repo = repo
}
func (rb RecomendBot) PickFilm(chatid int64) string {

	if err, res := rb.repo.PickRandom(chatid); err == nil {
		return res
	} else {
		fmt.Println("Could not pick film !")
		return ""
	}
}

func (rb RecomendBot) SendAnswer() {
	inlineKeyboard := tu.InlineKeyboard(
		tu.InlineKeyboardRow( // Row 1
			tu.InlineKeyboardButton("Порекомендуй фильм").
				WithCallbackData("recomend_film"),
		),
	)
	updates, _ := rb.bot.UpdatesViaLongPolling(nil)
	defer rb.bot.StopLongPolling()

	for update := range updates {
		var chatID int64 // ID чата

		if update.Message != nil {
			chatID = update.Message.Chat.ID
		} else if update.CallbackQuery != nil {
			chatID = update.CallbackQuery.Message.GetChat().ID
		} else {
			continue
		}
		if update.CallbackQuery != nil {
			callbackData := update.CallbackQuery.Data
			if callbackData == "recomend_film" {
				message := tu.Message(
					tu.ID(chatID), // Используем правильный ID чата
					rb.PickFilm(chatID),
				).WithReplyMarkup(inlineKeyboard)

				// Отправка сообщения
				_, _ = rb.bot.SendMessage(message)
			}
		} else {
			message := tu.Message(
				tu.ID(chatID), // Используем правильный ID чата
				"Давайте порекомендую вам фильм",
			).WithReplyMarkup(inlineKeyboard)

			// Отправка сообщения
			_, _ = rb.bot.SendMessage(message)
		}

	}

}
