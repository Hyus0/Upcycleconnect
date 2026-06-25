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
	UserID    int    `json:"user_id"`
	ForumID     int    `json:"forum_id"`
	TopicTitle  string `json:"topic_title"`
	Author   string `json:"author"`
	Role     string `json:"role"`
	Content  string `json:"content"`
	PostedAt string `json:"postedAt"`
	ImageProfil string `json:"image_profil"`
}

type CreateTopicRequest struct {
    UserID  int    `json:"user_id"`
    SalonID int    `json:"salon_id"`
    Titre   string `json:"titre"`
    Sujet   string `json:"sujet"`
}

type SendMessageRequest struct {
    UserID  int    `json:"user_id"`
    ForumID int    `json:"forum_id"`
    Contenu string `json:"contenu"`
}

type ReportDetail struct {
	ReporterName string `json:"reporter_name"`
	DateSignalement string `json:"date_signalement"`
	Motif        string `json:"motif"`
}

type ReportedMessage struct {
	MessageID      int    `json:"message_id"`
	Contenu        string `json:"contenu"`
	DateEnvoi      string `json:"date_envoi"`
	AuthorID       int    `json:"author_id"`
	AuthorName     string `json:"author_name"`
	AuthorImage    string `json:"author_image"`
	NbSignalements int    `json:"nb_signalements"`
	Details        []ReportDetail `json:"details"`
}

type BannedUser struct {
	ID       int    `json:"id"`
	FullName string `json:"full_name"`
	Role     string `json:"role"`
	ImageProfil string `json:"image_profil"`
}

type ModTopic struct {
	ID         int    `json:"id"`
	Titre      string `json:"titre"`
	Sujet      string `json:"sujet"`
	AuthorID   int    `json:"author_id"`
	AuthorName string `json:"author_name"`
	Date       string `json:"date"`
	MsgCount   int    `json:"msg_count"`
}