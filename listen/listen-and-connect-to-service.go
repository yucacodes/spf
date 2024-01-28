package listen

import "github.com/yucacodes/secure-port-forwarding/config"

type ListenAndConnectToService struct {
	listen *config.Listen
	node   *config.Node
}
