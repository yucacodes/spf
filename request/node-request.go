package request

import "github.com/yucacodes/secure-port-forwarding/config"

type NodeRequest struct {
	Id                                config.NodeId                             `yaml:"id" json:"id"`
	PublishService                    *PublishServiceRequest                    `yaml:"publishService" json:"publishService"`
	StreamingToServiceClient          *StreamingToServiceClient                 `yaml:"streamingToServiceClient" json:"streamingToServiceClient"`
	ForeignServiceClientConectionPair *ForeignServiceClientConectionPairRequest `yaml:"foreignServiceClientConectionPair" json:"foreignServiceClientConectionPair"`
}

type StreamingToServiceClient struct {
	Service string `yaml:"service" json:"service"`
	Client  string `yaml:"client" json:"client"`
}

type PublishServiceRequest struct {
	Service string `yaml:"service" json:"service"`
}

type ForeignServiceClientConectionPairRequest struct {
	Service string `yaml:"service" json:"service"`
	Client  string `yaml:"client" json:"client"`
}
