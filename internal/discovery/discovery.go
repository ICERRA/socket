package discovery

import (
	"net"

	client "github.com/etcd-io/etcd/clientv3"
)

type Identity = string

const (
	Prefix = "/service/access_layer/"
)

type Server struct {
	endpoints map[Identity]Endpoint

	etcdClient client.Config
}

type Endpoint struct {
	Identity string
	GRPCAddr net.Addr
}

func New() *Server {
	client.New(client.Config{
		Endpoints:            nil,
		AutoSyncInterval:     0,
		DialTimeout:          0,
		DialKeepAliveTime:    0,
		DialKeepAliveTimeout: 0,
		MaxCallSendMsgSize:   0,
		MaxCallRecvMsgSize:   0,
		TLS:                  nil,
		Username:             "",
		Password:             "",
		RejectOldCluster:     false,
		DialOptions:          nil,
		LogConfig:            nil,
		Context:              nil,
		PermitWithoutStream:  false,
	})
}

func (s Server) Register(endpoint Endpoint) error {
	return nil
}

func (s Server) Discovery(prefix string) (Endpoint, error) {
	return Endpoint{}, nil
}
