package main

import (
	"context"
	"flag"
	"log"
	"math/rand"
	"mutual-exclusion/proto"
	"net"
	"strconv"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Peer struct {
	proto.UnimplementedPeerServer
	port     int
	nextPort int
}

var (
	defaultPort = 8430
	port        = flag.Int("port", defaultPort, "port number")
	nextPort    = flag.Int("next", defaultPort+1, "peer port number")
)

var nextPeer proto.PeerClient

func (peer *Peer) SendToken(ctx context.Context, token *proto.Token) (*proto.Acknowledgement, error) {
	// wait until connected
	for nextPeer == nil {
		time.Sleep(time.Second)
	}

	// randomly select if critical section is wanted or not (~10% chance)
	if rand.Int()%10 == 0 {
		// critical section in use!
		log.Println(token.Timestamp)
		time.Sleep(time.Second * 2)
		// critical section is not used after this line
	}

	// Increment timestamp to keep track of critical section usage, even when not used
	go nextPeer.SendToken(context.Background(), &proto.Token{Timestamp: token.Timestamp + 1})

	return &proto.Acknowledgement{}, nil
}

func main() {
	flag.Parse()

	peer := &Peer{
		port:     *port,
		nextPort: *nextPort,
	}

	go startListener(peer)

	conn, _ := grpc.Dial(":"+strconv.Itoa(peer.nextPort), grpc.WithBlock(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	log.Println("Connected")
	defer conn.Close()

	nextPeer = proto.NewPeerClient(conn)

	establishRingIfDefault(peer)

	for {
		time.Sleep(time.Hour)
	}
}

func establishRingIfDefault(peer *Peer) {
	if peer.port == defaultPort {
		nextPeer.SendToken(context.Background(), &proto.Token{Timestamp: 0})
	}
}

// code adapted from TAs
// https://github.com/Mai-Sigurd/grpcTimeRequestExample?tab=readme-ov-file#setting-up-the-files
func startListener(peer *Peer) {
	// Create a new grpc server
	grpcServer := grpc.NewServer()

	// Make the server listen at the given port (convert int port to string)
	listener, err := net.Listen("tcp", ":"+strconv.Itoa(peer.port))

	if err != nil {
		log.Fatalf("Could not create the listener %v", err)
	}
	log.Printf("Started listening at port: %d\n", peer.port)

	// Register the grpc server and serve its listener
	proto.RegisterPeerServer(grpcServer, peer)
	serveError := grpcServer.Serve(listener)
	if serveError != nil {
		log.Fatalf("Could not serve listener")
	}
}
