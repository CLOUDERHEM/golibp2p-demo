package main

import (
	"context"
	"fmt"
	golog "github.com/ipfs/go-log/v2"
	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/core/peerstore"
	"github.com/libp2p/go-libp2p/p2p/protocol/circuitv2/client"
	ma "github.com/multiformats/go-multiaddr"
)

var (
	targetPeerIdString = "QmRGGu8681RyNRtFs1Z3EeXmkatwkyQwhDXfpLgnfXejta"

	relayPeerIdString = "QmfNspDuxDv6bnFbzntc33YEgy3sGM9NTg9gCeMfUqcoeY"
	relayAddr         = "/ip4/127.0.0.1/udp/64324/quic-v1"
)

func main() {
	golog.SetAllLoggers(golog.LevelDebug)

	node, err := createHost()
	if err != nil {
		panic(err)
	}
	fmt.Println("my peerId is: ", node.ID())

	relayInfo := connectToRelay(context.Background(), node)

	targetAddrInfo, err := getTargetAddrInfoWithRelay(relayInfo)
	if err != nil {
		return
	}
	node.Peerstore().AddAddrs(targetAddrInfo.ID, targetAddrInfo.Addrs, peerstore.PermanentAddrTTL)

	// node.Network().(*swarm.Swarm).Backoff().Clear(targetAddrInfo.ID)

	// connect to target
	err = node.Connect(context.Background(), *targetAddrInfo)
	if err != nil {
		panic(err)
	}

	// must use WithAllowLimitedConn
	ctx := network.WithAllowLimitedConn(context.Background(), "test")
	stream, err := node.NewStream(ctx, targetAddrInfo.ID, "/test")
	if err != nil {
		panic(err)
	}
	fmt.Println("created stream: ", stream.ID())

	for {
		fmt.Print("input msg: ")
		var msg string
		_, err := fmt.Scan(&msg)
		if err != nil {
			panic(err)
		}

		n, err := stream.Write([]byte(msg))
		if err != nil && err.Error() == "stream reset" {
			stream, err = node.NewStream(network.WithAllowLimitedConn(context.Background(), "test"),
				targetAddrInfo.ID, "/test")
			if err != nil {
				panic(err)
			}
		} else {
			fmt.Printf("\nwrite msg: %v, len: %v\n", msg, n)
		}
	}

}

func getTargetAddrInfoWithRelay(relayInfo peer.AddrInfo) (*peer.AddrInfo, error) {
	targetAddrWithRelay, err := ma.NewMultiaddr(
		fmt.Sprintf("/p2p/%s/p2p-circuit/p2p/%s", relayInfo.ID.String(), targetPeerIdString),
	)
	if err != nil {
		return nil, err
	}
	return peer.AddrInfoFromP2pAddr(targetAddrWithRelay)
}

func connectToRelay(ctx context.Context, node host.Host) peer.AddrInfo {
	relayAddr, err := ma.NewMultiaddr(relayAddr)
	if err != nil {
		panic(err)
	}
	relayPeerId, err := peer.Decode(relayPeerIdString)
	if err != nil {
		panic(err)
	}
	relayAddrInfo := peer.AddrInfo{
		ID:    relayPeerId,
		Addrs: []ma.Multiaddr{relayAddr},
	}
	err = node.Connect(ctx, relayAddrInfo)
	if err != nil {
		panic(err)
	} else {
		fmt.Println("connected to relay: ", relayPeerIdString)
	}
	_, err = client.Reserve(ctx, node, relayAddrInfo)
	if err != nil {
		panic(err)
	}
	return relayAddrInfo
}

func createHost() (host.Host, error) {
	key, err := loadOrCreatePrivateKey("private_key.pem")
	if err != nil {
		return nil, err
	}
	node, err := libp2p.New(
		libp2p.NoListenAddrs,
		libp2p.Identity(key),
		libp2p.EnableRelay(),
	)
	if err != nil {
		return nil, err
	}

	connLogger := &ConnLogger{}
	node.Network().Notify(connLogger)

	return node, err
}
