package client

import (
	"context"
	"errors"
	"fmt"

	kvpb "github.com/KingrogKDR/omni/proto/gen/kv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var ErrKeyNotFound = errors.New("key not found")

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

func NewClientFromConn(conn *grpc.ClientConn) *Client {
	return &Client{rpc: kvpb.NewOmniClient(conn), conn: conn}
}

func (c *Client) Put(ctx context.Context, cf, key, value []byte) error {
	resp, err := c.rpc.Put(ctx, &kvpb.PutRequest{
		Cf:  cf,
		Key: key,
		Val: value,
	})
	if err != nil {
		return fmt.Errorf("put error: %w", err)
	}
	if !resp.Success {
		return fmt.Errorf("put failed: server reported unsuccessful write")
	}

	return nil
}

func (c *Client) Get(ctx context.Context, cf, key []byte) ([]byte, error) {
	resp, err := c.rpc.Get(ctx, &kvpb.GetRequest{
		Cf:  cf,
		Key: key,
	})
	if err != nil {
		return nil, fmt.Errorf("get error: %w", err)
	}

	if resp.NotFound {
		return nil, ErrKeyNotFound
	}

	return resp.Val, nil
}

func (c *Client) Delete(ctx context.Context, cf, key []byte) error {
	resp, err := c.rpc.Delete(ctx, &kvpb.DeleteRequest{
		Cf:  cf,
		Key: key,
	})
	if !resp.Success {
		return fmt.Errorf("delete error: %w", err)
	}

	return nil
}

func (c *Client) Scan(ctx context.Context, cf []byte, scanOptions *kvpb.ScanOptions) ([]*kvpb.KeyValue, error) {
	resp, err := c.rpc.Scan(ctx, &kvpb.ScanRequest{
		Cf:          cf,
		ScanOptions: scanOptions,
	})
	if err != nil {
		return nil, fmt.Errorf("scan error: %w", err)
	}
	return resp.Pairs, nil
}

func (c *Client) Close() error {
	return c.conn.Close()
}
