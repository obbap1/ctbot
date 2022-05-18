package bot

import (
	"fmt"
	"os"
	"strconv"
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
			text := strings.Split(update.Message.Text, " ")
			if len(text) != 3 {
				// invalid text, please use this format
				return
			}
			operation := strings.ToLower(strings.TrimSpace(text[0]))
			if !IsValidOp(operation) {
				// invalid operation, please use either a "save", "update" or "delete"
				return
			}
			op := ToOp(operation)
			name := strings.ToLower(strings.TrimSpace(text[1]))
			date := strings.Split(strings.TrimSpace(text[2]), "/")
			if len(date) != 2 {
				// invalid date format, please use DD/MM
				return
			}
			day, err := strconv.Atoi(date[0])
			if err != nil {
				return
			}
			if day <= 0 || day > 31 {
				// invalid day, must be between 1 and 31
				return
			}
			month, err := strconv.Atoi(date[1])
			if err != nil {
				return
			}
			if month <= 0 || month > 12 {
				// invalid month, month must be between 1 and 12
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
