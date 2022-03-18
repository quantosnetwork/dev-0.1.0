package main

import (
	"context"
	"github.com/davecgh/go-spew/spew"
	"quantos/0.1.0/config"
)

func main() {
	c := config.NewNodeContext(context.Background())
	spew.Dump(c.Value("nctx"))
	spew.Dump(c.Value("node_id"))
}
