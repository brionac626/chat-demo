package main

import (
	"log"
	"net"

	"github.com/gobwas/ws"
	"github.com/rs/xid"
)

func main() {
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalln(err)
	}

	userPool := NewUserPool()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println(err)
			continue
		}

		upgradeAndStoreConn(conn, userPool)
	}
}

func upgradeAndStoreConn(conn net.Conn, userPool *UserPool) {
	_, err := ws.Upgrade(conn)
	if err != nil {
		log.Println(err)
		return
	}

	userID := xid.New().String()
	user, _ := userPool.Users.LoadOrStore(userID, &User{Conn: conn, ID: userID})
	u := user.(*User)
	userPool.GetAllUser()
	go func() {
		for {
			if err := u.Receive(); err != nil {
				log.Println(err)
				u.Conn.Close()
				userPool.DeleteCloseUser(u)
				break
			}
		}
	}()
}
