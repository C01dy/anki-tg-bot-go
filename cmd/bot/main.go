package main

import (
	"anki-bot/internal/adapter/repository"
	"anki-bot/internal/adapter/telegram"
	"anki-bot/internal/config"
	"anki-bot/internal/usecase"
	"fmt"
	"log"
	"strings"
)

func main() {
	cfg, err := config.GetConfig(".env")
	if err != nil {
		log.Fatalf("Config error: %v", err)
	}

	sqliteRepo, err := repository.NewSQLiteRepo(cfg.DB)
	if err != nil {
		log.Fatalf("Failed to init db: %v", err)
	}

	client := telegram.New("api.telegram.org", cfg.TelegramToken)

	vocabularyService := usecase.NewVocabularyService(sqliteRepo, client)

	offset := 0
	fmt.Println("Bot has been started")

	for {
		updates, err := client.FetchUpdates(offset, 10)
		if err != nil {
			log.Printf("error fetching updates: %v", err)
			continue
		}

		for _, upd := range updates {
			offset = upd.ID + 1 

			if upd.Message == nil {
				continue
			}

			chatID := upd.Message.Chat.ID
			text := upd.Message.Text

			log.Printf("[%d] Message from %d: %s", upd.ID, chatID, text)

			if strings.HasPrefix(text, "/add") {
				parts := strings.Fields(text)
				
				if len(parts) < 3 {
					client.SendMessage(chatID, "Используй формат: /add <en> <ru>")
					continue
				}

				en := parts[1]
				ru := parts[2]

				err := vocabularyService.AddWord(chatID, en, ru)
				if err != nil {
					log.Printf("Failed to add word: %v", err)
					client.SendMessage(chatID, "Ошибка при сохранении слова ❌")
				} else {
					client.SendMessage(chatID, fmt.Sprintf("Слово '%s' сохранено! ✅", en))
				}
				continue
			}

			client.SendMessage(chatID, "Я пока умею только добавлять слова через /add")
		}
	}
}