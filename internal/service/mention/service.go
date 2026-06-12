package mentionservice

import (
	"regexp"
	"strings"

	messagedomain "ChatServerGolang/internal/domain/message"
	"ChatServerGolang/internal/repository"
)

var mentionRegex = regexp.MustCompile(`@(\w{3,32})`)

type MentionService interface {
	ExtractAndSaveMentions(messageID, content string) []*messagedomain.MentionedUser
	GetMentionsByMessageID(messageID string) ([]*messagedomain.Mention, error)
}

type mentionService struct {
	userRepo    repository.UserRepository
	messageRepo repository.MessageRepository
}

func NewMentionService(userRepo repository.UserRepository, messageRepo repository.MessageRepository) MentionService {
	return &mentionService{
		userRepo:    userRepo,
		messageRepo: messageRepo,
	}
}

func (s *mentionService) ExtractAndSaveMentions(messageID, content string) []*messagedomain.MentionedUser {
	matches := mentionRegex.FindAllStringSubmatchIndex(content, -1)
	if len(matches) == 0 {
		return nil
	}

	var mentioned []*messagedomain.MentionedUser
	seen := make(map[string]bool)

	for _, m := range matches {
		if len(m) < 4 {
			continue
		}
		username := content[m[2]:m[3]]
		if seen[username] {
			continue
		}
		seen[username] = true

		user, err := s.userRepo.FindByUsername(username)
		if err != nil || user == nil {
			continue
		}

		if err := s.messageRepo.SaveMention(messageID, user.ID, username); err != nil {
			continue
		}

		mentioned = append(mentioned, &messagedomain.MentionedUser{
			UserID:   user.ID,
			Username: username,
			Offset:   m[2],
			Length:   m[3] - m[2],
		})
	}

	return mentioned
}

func (s *mentionService) GetMentionsByMessageID(messageID string) ([]*messagedomain.Mention, error) {
	return s.messageRepo.GetMentionsByMessageID(messageID)
}

// Helper to find all mentioned user IDs in a message
func FindMentionedUsernames(content string) []string {
	matches := mentionRegex.FindAllStringSubmatch(content, -1)
	usernames := make([]string, 0, len(matches))
	seen := make(map[string]bool)
	for _, m := range matches {
		if !seen[m[1]] {
			seen[m[1]] = true
			usernames = append(usernames, strings.ToLower(m[1]))
		}
	}
	return usernames
}





