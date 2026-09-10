package service

import (
	"context"
	"errors"
	repository "go-chatapp/repository/generated"
)

type ChatService struct {
	queries *repository.Queries
}

func NewChatService(queries *repository.Queries) *ChatService {
	return &ChatService{
		queries: queries,
	}
}

func (s *ChatService) SendMessage(ctx context.Context, userID int64, conversationID int64, content string) (repository.CreateMessageRow, error) {
	isParticipant, err := s.queries.IsUserInConversation(
		ctx,
		repository.IsUserInConversationParams{
			ConversationID: conversationID,
			UserID:         userID,
		},
	)
	if err != nil {
		return repository.CreateMessageRow{}, err
	}
	if !isParticipant {
		return repository.CreateMessageRow{}, errors.New("user is not a participant")
	}

	message, err := s.queries.CreateMessage(
		ctx,
		repository.CreateMessageParams{
			ConversationID: conversationID,
			SenderID:       userID,
			Content:        content,
			MessageType:    "text",
		},
	)
	if err != nil {
		return repository.CreateMessageRow{}, err
	}

	return message, nil
}
