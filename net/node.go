package net

import (
	_ "github.com/libp2p/go-libp2p"
	_ "github.com/libp2p/go-libp2p-core/peer"
	_ "github.com/multiformats/go-multiaddr"
)

type Node struct {
	host host.Host
}
