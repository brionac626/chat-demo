package main

type Chatroom struct {
	ID        string
	Name      string
	Members   []*User
	CreatedAt int64
}
