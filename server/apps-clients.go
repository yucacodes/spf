package server

import (
	"fmt"
	"net"
	"strconv"

	"github.com/google/uuid"
	transfers "github.com/yucacodes/secure-port-forwarding/messages"
)

type ConnectClientRequest struct {
	ClientId string
}

func listenAppClients(connApp *ConnectedApp) {
	defer connApp.Conn.Close()
	clientsListener, err := net.Listen("tcp", "0.0.0.0:"+strconv.Itoa(connApp.Config.Port))
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer clientsListener.Close()
	fmt.Println("App is listening on port " + strconv.Itoa(connApp.Config.Port))

	for {
		appClientConn, err := clientsListener.Accept()
		if err != nil {
			fmt.Println("Error:", err)
			continue
		}

		fmt.Println("Sub Client connected", appClientConn)
		appClientId := uuid.New()
		req := ConnectClientRequest{ClientId: appClientId.String()}

		err = transfers.Send(connApp.Conn, req)
		if err != nil {
			fmt.Println("Error on request connection to App")
		}
		connApp.Clients[appClientId.String()] = appClientConn
	}
}
