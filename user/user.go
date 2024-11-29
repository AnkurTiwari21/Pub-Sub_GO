package user

import "net"

type User struct{
	UserConnections []*net.Conn
}