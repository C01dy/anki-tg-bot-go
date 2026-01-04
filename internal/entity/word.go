package entity

import "time"

type Word struct {
	EN               string
	RU               string
	NextRetry        time.Time
	Interval         int
	EaseFactor   	 float64
	Repetitions 	 int
}
