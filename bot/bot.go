package bot

import (
	"fmt"
	"os"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	database "github.com/obbap1/ctbot/db"
	"gorm.io/gorm"
)

func Init(db *gorm.DB) {
	bot, err := tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_APITOKEN"))
	if err != nil {
		panic(err)
	}
	bot.Debug = true

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil {
			fmt.Printf("from %s, text %s", update.Message.From.UserName, update.Message.Text)
			if update.Message.From.IsBot {
				// sorry, we dont support bots
				return
			}
			// eg. save paschal 05/12
			// delete paschal
			// get 05/12
			// update paschal 04/12
			text := strings.Split(update.Message.Text, " ")
			if len(text) == 0 {
				// invalid text
				return
			}
			operation := strings.ToLower(strings.TrimSpace(text[0]))
			if !IsValidOp(operation) {
				// invalid operation, please use either a "save", "get", "update" or "delete"
				return
			}
			op := ToOp(operation)
			if r := op.Requirements(); r.numberOfCommands != len(text)-1 {
				// invalid text, please use this format
				return
			}

			name, day, month, err := op.Convert(text)
			if err != nil {
				// invalid text
				return
			}

			if err := op.Do(&database.Reminder{
				UserID:    int(update.Message.From.ID),
				Month:     month,
				Day:       day,
				Celebrant: name,
			}); err != nil {
				// couldnt perform operation
				return
			}

			switch op {
			case Save:
				if result := db.Create(&database.Reminder{
					UserID:    int(update.Message.From.ID),
					Month:     month,
					Day:       day,
					Celebrant: name,
				}); result.Error != nil {
					// invalid message, try again later
					return
				}
			case Update:
				reminder := &database.Reminder{
					UserID:    int(update.Message.From.ID),
					Month:     month,
					Day:       day,
					Celebrant: name}
				// fix
				if result := db.Save(reminder); result.Error != nil {
					// invalid message, try again later
					return
				}
			case Delete:
				// db.Delete()
			}

		}
	}
}

func sendMessage(bot *tgbotapi.BotAPI, chatID int64, text string) error {
	_, err := bot.Send(tgbotapi.NewMessage(chatID, text))
	return err
}
