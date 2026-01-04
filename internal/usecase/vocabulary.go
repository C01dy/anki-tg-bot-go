package usecase

import (
	"anki-bot/internal/entity"
	"fmt"
	"math"
	"time"
)

type VocabularyService struct {
	repo   WordRepository
	sender BotSender
}

func NewVocabularyService(r WordRepository, s BotSender) *VocabularyService {
	return &VocabularyService{
		repo:   r,
		sender: s,
	}
}

func (vs *VocabularyService) AddWord(userID int64, en, ru string) error {
	word := entity.Word{
		EN:        en,
		RU:        ru,
		NextRetry: time.Now(),
	}

	if err := vs.repo.Save(word, userID); err != nil {
		return fmt.Errorf("error when try to save %s/%s word for user %d: %w", en, ru, userID, err)
	}

	return nil
}

func (vs *VocabularyService) GetWord(userID int64, en string) (entity.Word, error) {
	// TODO: implement
	return entity.Word{}, nil
}

func (vs *VocabularyService) calculateSM2(w entity.Word, quality int) entity.Word {
	if quality < 0 { quality = 0 }
	if quality > 5 { quality = 5 }

	w.EaseFactor = w.EaseFactor + (0.1 - float64(5-quality)*(0.08+float64(5-quality)*0.02))
	if w.EaseFactor < 1.3 {
		w.EaseFactor = 1.3
	}

	if quality < 3 {
		w.Repetitions = 0
		w.Interval = 1
	} else {
		w.Repetitions++
		
		switch w.Repetitions {
		case 1:
			w.Interval = 1
		case 2:
			w.Interval = 6
		default:
			w.Interval = int(math.Round(float64(w.Interval) * w.EaseFactor))
		}
	}

	w.NextRetry = time.Now().AddDate(0, 0, w.Interval)
	
	return w
}
