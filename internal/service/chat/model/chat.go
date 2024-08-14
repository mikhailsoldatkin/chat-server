package model

// Chat represents a business logic chat model.
type Chat struct {
	ID int64
}

// ChatUser represents the business logic association between a chat and its users.
type ChatUser struct {
	ChatID int64
	UserID int64
}
