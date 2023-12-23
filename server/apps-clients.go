package server

import (
	"fmt"
	"net"

	"github.com/google/uuid"
	"github.com/yucacodes/secure-port-forwarding/transfers"
)

type ConnectClientRequest struct {
	clientId string
}

func listenAppClients(connApp *ConnectedApp) {
	defer connApp.Conn.Close()
	clientsListener, err := net.Listen("tcp", "0.0.0.0:"+string(connApp.Config.Port))
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer clientsListener.Close()
	fmt.Println("App is listening on port " + string(connApp.Config.Port))

	for {
		appClientConn, err := clientsListener.Accept()
		if err != nil {
			fmt.Println("Error:", err)
			continue
		}

		fmt.Println("Sub Client connected", appClientConn)
		appClientId := uuid.New()
		req := ConnectClientRequest{clientId: appClientId.String()}

		err = transfers.Write(connApp.Conn, req)
		if err != nil {
			fmt.Println("Error on request connection to App")
		}
		connApp.Clients[appClientId.String()] = appClientConn
	}
}
