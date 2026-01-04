package usecase

import (
	"anki-bot/internal/entity"
	"fmt"
	"time"
)

type VocabularyService struct {
	repo WordRepository
	sender BotSender
}

func NewVocabularyService(r WordRepository, s BotSender) *VocabularyService {
	return &VocabularyService{
		repo: r,
		sender: s,
	}
}

func (vs *VocabularyService) AddWord(userID int64, en, ru string) error {
	word := entity.Word{
		EN: en,
		RU: ru,
		NextRetry: time.Now(),
	}

	if err := vs.repo.Save(word, userID); err != nil {
		return fmt.Errorf("error when try to save %s/%s word for user %d: %w", en, ru, userID, err)
	}

	return nil
}
