package main

import (
	"encoding/json"
	"fmt"
	"net"
	"sync"

	"github.com/gobwas/ws"

	"github.com/brionac626/chat-demo/models"
	"github.com/pkg/errors"

	"github.com/rs/xid"

	"github.com/gobwas/ws/wsutil"
)

type UserPool struct {
	Users sync.Map
}

type User struct {
	ID     string
	Name   string
	Conn   net.Conn
	Locker sync.Mutex
}

func NewUserPool() *UserPool {
	return &UserPool{}
}

func (up *UserPool) DeleteCloseUser(u *User) {
	up.Users.Delete(u.ID)
}

func (up *UserPool) GetAllUser() {
	up.Users.Range(func(key, value interface{}) bool {
		fmt.Println(key, value)
		return true
	})
}

func NewUser(conn net.Conn) *User {
	return &User{
		ID:     xid.New().String(),
		Conn:   conn,
		Locker: sync.Mutex{},
	}
}

func (u *User) Receive() error {
	req, err := u.readRequest()
	if err != nil {
		return err
	}

	switch req.Method {
	case "echo":
		fmt.Printf("payload: \n%+v\n", req)
		for k, v := range req.Payload {
			fmt.Println(k, v)
		}
	case "publish":
		resp := models.UserResponse{UserID: u.ID, Status: true, Message: "get message"}
		payload, err := json.Marshal(resp)
		if err != nil {
			return err
		}

		u.Conn.Write(payload)
	default:
		return errors.Errorf("unknow method: %s", req.Method)
	}

	return nil
}

func (u *User) readRequest() (*models.UserRequest, error) {
	u.Locker.Lock()
	defer u.Locker.Unlock()

	b, s, err := wsutil.ReadData(u.Conn, ws.StateServerSide)
	if err != nil {
		fmt.Println("456")
		return nil, err
	}

	if s.IsData() {
		var req models.UserRequest
		if err := json.Unmarshal(b, &req); err != nil {
			return nil, err
		}

		return &req, nil
	}

	return nil, errors.Errorf("status not data: %v", s)
}

func (u *User) WriteMessage(msg []byte) error {
	return nil
}
