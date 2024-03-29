package cmd

import (
	"log"

	"github.com/asavie/xdp"
	bpf "github.com/kerwenwwer/xdp-gossip/bpf"
	common "github.com/kerwenwwer/xdp-gossip/common"
	"github.com/vishvananda/netlink"
)

// New initializes the local node list
func ProgramHandler(LinkName string, obj *bpf.BpfObjects, debug bool) (*xdp.Program, *xdp.Socket) {
	// Get netlink by name
	link, err := netlink.LinkByName(LinkName)
	if err != nil {
		log.Fatalf("[Error]: Failed to get link by name %v", err)
	}

	// Attach Tc program
	if err := bpf.AttachTC(obj, link); err != nil {
		log.Fatalf("[Error]: Failed to attach TC: %v", err)
	}

	if debug {
		log.Printf("[Info]: TC attached. ")
	}
	//Attach XDP program
	program, err := bpf.AttachXDP(obj, link.Attrs().Index)
	if err != nil {
		log.Fatalf("[Error]: Failed to attach XDP: %v", err)
	}

	xsk, err := xdp.NewSocket(link.Attrs().Index, 0, &xdp.SocketOptions{
		NumFrames:              128,
		FrameSize:              2048,
		FillRingNumDescs:       64,
		CompletionRingNumDescs: 64,
		RxRingNumDescs:         64,
		TxRingNumDescs:         64,
	})
	if err != nil {
		log.Fatal("error: failed to create an XDP socket: ", err)
	}

	if err := program.Register(0, xsk.FD()); err != nil {
		log.Fatal("error: failed to register socket in BPF map: ", err)
	}
	defer program.Unregister(0)

	if debug {
		log.Printf("[Info]: AF_XDP registered.")
	}

	return program, xsk
	//return nil, nil
}

// func (nl *NodeList) storeWithCheck(node common.Node) {
// 	nl.nodes.Range(func(k, v interface{}) bool {
// 		existingNode, ok := k.(common.Node)
// 		if ok && existingNode.Addr == node.Addr {
// 			fmt.Printf("Node with Addr %s already exists!\n", node.Addr)
// 			return false
// 		}
// 		return true
// 	})

// 	nl.nodes.Store(node, time.Now().Unix())
// }

type MyPacket common.Packet
