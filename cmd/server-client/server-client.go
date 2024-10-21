package main

import (
	"fmt"

	"github.com/cego/nebula-provisioner/protocol"
	conf "github.com/slackhq/nebula/config"
	"google.golang.org/grpc"
)

type serverClient struct {
	conn   *grpc.ClientConn
	client protocol.ServerCommandClient
}

func NewClient(config *conf.C) (*serverClient, error) {

	var opts []grpc.DialOption

	opts = append(opts, grpc.WithInsecure())

	socketPath := config.GetString("command.socket", "/tmp/nebula-provisioner.socket") // TODO Change default path
	conn, err := grpc.Dial(fmt.Sprintf("unix://%s", socketPath), opts...)
	if err != nil {
		return nil, err
	}

	client := protocol.NewServerCommandClient(conn)

	return &serverClient{conn, client}, nil
}

func (c serverClient) Close() error {
	return c.conn.Close()
}
