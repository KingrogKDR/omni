package client

import (
	"context"
	"errors"
	"fmt"
	"io"
	"iter"

	"github.com/KingrogKDR/omni/internal/storage"
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

func (c *Client) Put(ctx context.Context, key, value []byte) (*kvpb.PutResponse, error) {
	resp, err := c.rpc.Put(ctx, &kvpb.PutRequest{
		Pair: &kvpb.KeyValue{
			Key: key,
		},
	})
	if err != nil {
		return nil, fmt.Errorf("put error: %w", err)
	}

	return resp, nil
}

func (c *Client) Get(ctx context.Context, key []byte) (*kvpb.GetResponse, error) {
	resp, err := c.rpc.Get(ctx, &kvpb.GetRequest{
		Key: key,
	})
	if err != nil {
		return nil, fmt.Errorf("get error: %w", err)
	}

	return resp, nil
}

func (c *Client) Delete(ctx context.Context, key []byte) (*kvpb.DeleteResponse, error) {
	resp, err := c.rpc.Delete(ctx, &kvpb.DeleteRequest{
		Key: key,
	})
	if err != nil {
		return nil, fmt.Errorf("delete error: %w", err)
	}

	return resp, nil
}

func (c *Client) list(ctx context.Context, req *kvpb.ListRequest) (iter.Seq2[[]byte, []byte], *storage.ReadError) {
	readErr := &storage.ReadError{}
	seq := func(yield func(k, v []byte) bool) {
		stream, err := c.rpc.List(ctx, req)
		if err != nil {
			readErr.SetErr(fmt.Errorf("list opening stream: %w", err))
			return
		}
		for {
			resp, err := stream.Recv()
			if errors.Is(err, io.EOF) {
				return
			}
			if err != nil {
				readErr.SetErr(fmt.Errorf("list: receiving: %w", err))
				return
			}
			if !yield(resp.Pair.Key, resp.Pair.Value) {
				return
			}
		}
	}

	return seq, readErr
}

func (c *Client) List(ctx context.Context, prefix []byte, limit uint32) (iter.Seq2[[]byte, []byte], *storage.ReadError) {
	req := &kvpb.ListRequest{Prefix: prefix}
	if limit > 0 {
		req.Limit = &limit
	}
	return c.list(ctx, req)
}

func (c *Client) Scan(ctx context.Context, start, end []byte) (iter.Seq2[[]byte, []byte], *storage.ReadError) {
	readErr := &storage.ReadError{}
	seq := func(yield func(k, v []byte) bool) {
		stream, err := c.rpc.Scan(ctx, &kvpb.ScanRequest{
			Start: start,
			End:   end,
		})
		if err != nil {
			readErr.SetErr(fmt.Errorf("scan error: %w", err))
			return
		}

		for {
			resp, err := stream.Recv()

			if errors.Is(err, io.EOF) {
				return
			}
			if err != nil {
				readErr.SetErr(fmt.Errorf("scan: receiving: %w", err))
				return
			}
			if !yield(resp.Pair.Key, resp.Pair.Value) {
				return
			}
		}
	}

	return seq, readErr
}

func (c *Client) Close() error {
	return c.conn.Close()
}
