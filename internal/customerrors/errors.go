package customerrors

import "fmt"

// NotFoundError represents an error for a missing entity with additional context.
type NotFoundError struct {
	Entity string
	ID     int64
}

// Error implements the error interface for NotFoundError.
func (e *NotFoundError) Error() string {
	return fmt.Sprintf("%s %d not found", e.Entity, e.ID)
}

// NewNotFoundError creates a new NotFoundError for a given entity and ID.
func NewNotFoundError(entity string, id int64) error {
	return &NotFoundError{
		Entity: entity,
		ID:     id,
	}
}

// UserNotInChatError represents an error indicating that a user can't be found in a chat.
type UserNotInChatError struct {
	UserID int64
	ChatID int64
}

// Error implements the error interface for UserNotInChatError.
func (e *UserNotInChatError) Error() string {
	return fmt.Sprintf("user %d not found in chat %d", e.UserID, e.ChatID)
}

// NewUserNotInChatError creates a new UserNotInChatError.
func NewUserNotInChatError(userID, chatID int64) error {
	return &UserNotInChatError{
		UserID: userID,
		ChatID: chatID,
	}
}
