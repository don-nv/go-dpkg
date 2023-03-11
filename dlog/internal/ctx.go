package internal

import "context"

type ReadContext func(ctx context.Context) []KV

type KV struct {
	Key   string
	Value interface{}
}
