package main

import (
	"anki-bot/internal/config"
	"anki-bot/internal/adapter/telegram"
	"fmt"
	"log"
)

func main() {
	cfg, err := config.GetConfig(".env")
	if err != nil {
		panic(err)
	}
	client := telegram.New("api.telegram.org", cfg.TelegramToken)

	offset := 0
	fmt.Println("Bot has been started")

	for {
		updates, err := client.FetchUpdates(offset, 10)
		if err != nil {
			log.Printf("error fetching updates: %v", err)
			continue
		}

		for _, upd := range updates {
			if upd.Message != nil {
				fmt.Printf("Пользователь пишет: %s\n", upd.Message.Text)

				replyText := "Ты сказал: " + upd.Message.Text

				err := client.SendMessage(upd.Message.Chat.ID, replyText)
				if err != nil {
					fmt.Printf("Не удалось отправить ответ: %v\n", err)
				}
			}
			offset = upd.ID + 1
			if upd.Message != nil {
				fmt.Printf("[%d] Message by %d: %s\n", upd.ID, upd.Message.Chat.ID, upd.Message.Text)
			}
		}
	}

}
