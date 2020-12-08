package models

type UserRequest struct {
	UserID  string
	Method  string
	Payload map[string]interface{}
}

type UserResponse struct {
	UserID  string
	Status  bool
	Message string
}

type Message struct {
	ChatroomID  string
	AutherID    string
	Message     string
	PublishedAt int64
}
