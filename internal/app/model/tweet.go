package model

import "time"

type Tweet struct {
	ID           uint      `json:"id"`
	UserID       uint      `json:"userId"`
	TweetMessage string    `jsong:"message"`
	DateTime     time.Time `json:"datetime"`
}
