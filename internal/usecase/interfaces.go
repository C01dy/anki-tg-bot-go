package usecase


import "anki-bot/internal/entity"

type WordRepository interface {
	Save(w entity.Word, userID int64) error 
	GetForReview(userID int64) ([]entity.Word, error)  
}

type BotSender interface {
	SendMessage(chatID int64, text string) error
}

