package model

type User struct {
	ID               uint     `json:"id"`
	UserName         string   `json:"username"`
	UserEmail        string   `json:"email"`
	UserPasswordHash string   `json:"password"`
	UserTweets       []Tweet  `json:"tweets"`
	UserTweetsFeed   []Tweet  `json:"feeds"`
	UserFollowers    []string `json:"followers"`
	UserFollowing    []string `json:"following"`
}
