package entity

import "time"

type Word struct {
	EN               string
	RU               string
	NextRetry        time.Time
	Interval         int32
	CorrectAnswers   int32
	IncorrectAnswers int32
}
