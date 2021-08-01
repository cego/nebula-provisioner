package main

import (
	"github.com/slyngdk/nebula-provisioner/protocol"
	"google.golang.org/grpc"
)

type serverClient struct {
	conn   *grpc.ClientConn
	client protocol.ServerCommandClient
}

func NewClient() (*serverClient, error) {

	var opts []grpc.DialOption

	opts = append(opts, grpc.WithInsecure())

	conn, err := grpc.Dial("unix:///tmp/nebula-provisioner.socket", opts...) // TODO change socket path
	if err != nil {
		return nil, err
	}

	client := protocol.NewServerCommandClient(conn)

	return &serverClient{conn, client}, nil
}

func (c serverClient) Close() error {
	return c.conn.Close()
}
