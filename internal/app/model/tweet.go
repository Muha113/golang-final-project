package model

import "time"

type Tweet struct {
	ID           uint      `json:"id"`
	UserName     string    `json:"userId"`
	TweetMessage string    `jsong:"message"`
	DateTime     time.Time `json:"datetime"`
}
