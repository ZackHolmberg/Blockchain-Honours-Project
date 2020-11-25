package blockchain

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/grandcat/zeroconf"
)

// ============================ Communication ============================

// Communicator implements CommunicationsComponent and facilities Blockchain communication
type Communicator struct {
	peer          *zeroconf.Server
	peerNodes     []string
	peerMessage   chan string
	clientMessage chan string
	//Not sure if clientMessage will be necessary, not sure how to best handle receiving new transactions from client/messageQueue
}

// GetPeerChains is the interface method that retrieves the copy of the blockchain from every peer on the network
func (c Communicator) GetPeerChains() {
	//TODO: Implement
}

// RecieveFromClient is the interface method that receives messages
// from the client (MIGHT NOT BE NECESSARY)
func (c Communicator) RecieveFromClient() {
	//TODO: Implement
}

// SendToClient is the interface method that sends messages
// to the client (MIGHT NOT BE NECESSARY)
func (c Communicator) SendToClient() {
	//TODO: Implement
}

// RecieveFromNetwork is the interface method that
// receives UDP message from peers on the network
func (c Communicator) RecieveFromNetwork() {
	//TODO: Implement
}

// BroadcastToNetwork is the interface method that uses
// UDP to broadcast a message to all the peers on the network
func (c Communicator) BroadcastToNetwork() {
	//TODO: Implement
}

// == Non-interface methods ==

// PingNetwork is the interface method that
func (c Communicator) PingNetwork() {
	//Will use SendToFromNetwork
	//TODO: Implement
}

// HandlePingFromNetwork is the interface method that
func (c Communicator) HandlePingFromNetwork() {
	//Will use RecieveFromNetwork
	//TODO: Implement
}

// NewCommunicator returns a new initiliazed Communicator
func NewCommunicator(name string, service string, domain string, port int) Communicator {

	newCommunicator := Communicator{}
	newCommunicator.initializeCommunicator(name, service, domain, port)
	discoverServices()
	return newCommunicator
}

// TerminateCommunicator cleans up and terminates the service
func (c Communicator) TerminateCommunicator() {
	fmt.Println("Terminating service...")
	c.peer.Shutdown()
}

// initializeCommunicator initializes a new communicator by initializing
// a ZeroConf service and discovering other services
func (c *Communicator) initializeCommunicator(name string, service string, domain string, port int) {

	fmt.Println("Starting service...")
	//Register the service
	peer, err := zeroconf.Register(name, service, domain, port, []string{"txtv=0", "lo=1", "la=2"}, nil)
	if err != nil {
		panic(err)
	}

	fmt.Println("Published service:")
	fmt.Println("- Name:", name)
	fmt.Println("- Type:", service)
	fmt.Println("- Domain:", domain)
	fmt.Println("- Port:", port)

	//Set this Communicator's service reference to the newly created service
	c.peer = peer

	// ============== TODO ==============
	// When a service/peer connects to the network (for the first time I'd think), it needs
	// to run consensus so that it only creates a genesis node if its the first node on the network
	// Otherwise, we dont want to generate a genesis node but instead "download" a copy of the current chain
}

func discoverServices() {
	// Discover all services on the '_blockchain-P2P-Network._udp' blockchain network
	resolver, err := zeroconf.NewResolver(nil)
	if err != nil {
		log.Fatalln("Failed to initialize resolver:", err.Error())
	}

	entries := make(chan *zeroconf.ServiceEntry)
	go func(results <-chan *zeroconf.ServiceEntry) {
		for entry := range results {
			log.Println("\n\n", entry, "\n")
		}
		log.Println("No more entries.")
	}(entries)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(5))
	defer cancel()
	err = resolver.Browse(ctx, "_blockchain-P2P-Network._udp", "local.", entries)
	if err != nil {
		log.Fatalln("Failed to browse:", err.Error())
	}

	<-ctx.Done()
	// Wait some additional time to see debug messages on go routine shutdown.
	time.Sleep(1 * time.Second)
}
