package server

import (
	"errors"
	"fmt"
	"net"

	transfers "github.com/yucacodes/secure-port-forwarding/messages"
	"github.com/yucacodes/secure-port-forwarding/stream"
)

type AppRequest struct {
	Key              string
	SetAppConnection bool
	CallbackToClient string
}

type ConnectedApp struct {
	Conn    net.Conn
	Config  AppConfig
	Clients map[string]net.Conn
}

func listenApps(allowedApps []AppConfig, appsListener net.Listener) {
	appsPool := make(map[string]*ConnectedApp)

	for {
		conn, err := appsListener.Accept()
		if err != nil {
			fmt.Println("Error:", err)
			continue
		}
		fmt.Println("App connected")
		go handleAppConnection(conn, allowedApps, appsPool)
	}
}

func handleAppConnection(conn net.Conn, allowedApps []AppConfig, appsPool map[string]*ConnectedApp) {
	defer conn.Close()

	app, req, err := ValidateAppConnectionRequest(conn, allowedApps)
	if err != nil {
		return
	}

	if req.SetAppConnection {
		prevConn := appsPool[app.Key]
		if prevConn != nil {
			prevConn.Conn.Close()
		}
		connectedApp := &ConnectedApp{Conn: conn, Config: *app, Clients: make(map[string]net.Conn)}
		appsPool[app.Key] = connectedApp
		listenAppClients(connectedApp)
	} else if req.CallbackToClient != "" {
		connectedApp, exist := appsPool[app.Key]
		if !exist {
			return
		}
		client, exist := connectedApp.Clients[req.CallbackToClient]
		if !exist {
			return
		}
		stream.HandlePairStream(conn, client)
	}
}

func ValidateAppConnectionRequest(conn net.Conn, allowedApps []AppConfig) (*AppConfig, *AppRequest, error) {
	req := AppRequest{}
	err := transfers.Receive(conn, &req)
	if err != nil {
		return nil, nil, err
	}
	for i := range allowedApps {
		if allowedApps[i].Key == req.Key {
			return &allowedApps[i], &req, nil
		}
	}
	return nil, nil, errors.New("not found config for requested app")
}
