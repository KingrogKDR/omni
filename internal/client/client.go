package client

import (
	"context"
	"fmt"

	kvpb "github.com/KingrogKDR/omni/proto/gen/kv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	rpc  kvpb.OmniClient
	conn *grpc.ClientConn
}

func NewClient(serverAddr string) (*Client, error) {
	clientConn, err := grpc.NewClient(serverAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	omniClient := kvpb.NewOmniClient(clientConn)
	return &Client{
		rpc:  omniClient,
		conn: clientConn,
	}, nil
}

func (c *Client) Put(ctx context.Context, key, value string) (*kvpb.PutResponse, error) {
	resp, err := c.rpc.Put(ctx, &kvpb.PutRequest{
		Key:   []byte(key),
		Value: []byte(value),
	})
	if err != nil {
		return nil, fmt.Errorf("put error: %w", err)
	}

	return resp, nil
}

func (c *Client) Get(ctx context.Context, key string) (*kvpb.GetResponse, error) {
	resp, err := c.rpc.Get(ctx, &kvpb.GetRequest{
		Key: []byte(key),
	})
	if err != nil {
		return nil, fmt.Errorf("get error: %w", err)
	}

	return resp, nil
}

// func (c *Client) Delete(ctx context.Context, key string) (*kvpb.DeleteResponse, error) {
// 	resp, err := c.rpc.Delete(ctx, &kvpb.DeleteRequest{
// 		Key: []byte(key),
// 	})
// 	if err != nil {
// 		return nil, fmt.Errorf("delete error: %w", err)
// 	}

// 	return resp, nil
// }

func (c *Client) Close() error {
	return c.conn.Close()
}
