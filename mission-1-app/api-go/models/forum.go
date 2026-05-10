package models

type ForumCategory struct {
	ID          int          `json:"id"`
	Name        string       `json:"name"`
	Description string       `json:"description"`
	Topics      []ForumTopic `json:"topics"`
}

type ForumTopic struct {
	ID           int            `json:"id"`
	Title        string         `json:"title"`
	Preview      string         `json:"preview"`
	LastActivity string         `json:"lastActivity"`
	Messages     []ForumMessage `json:"messages"`
}

type ForumMessage struct {
	ID       int    `json:"id"`
	Author   string `json:"author"`
	Role     string `json:"role"`
	Content  string `json:"content"`
	PostedAt string `json:"postedAt"`
}