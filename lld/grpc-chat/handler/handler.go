package handler

import (
	"context"
	"fmt"
	"sync"
	"system-design/lld/grpc-chat/gen"
)

type Connection struct {
	gen.UnimplementedBroadcastServer
	stream gen.Broadcast_CreateStreamServer
	id     string
	active bool
	error  chan error
}

type Pool struct {
	gen.UnimplementedBroadcastServer
	Connection []*Connection
}

func (p *Pool) CreateStream(pconn *gen.Connect, stream gen.Broadcast_CreateStreamServer) error {
	conn := &Connection{
		stream: stream,
		id:     pconn.User.ID,
		active: true,
		error:  make(chan error),
	}

	p.Connection = append(p.Connection, conn)

	return <-conn.error
}

func (p *Pool) BroadcastMessage(ctx context.Context, msg *gen.Message) (*gen.Close, error) {
	wait := sync.WaitGroup{}
	done := make(chan int)

	for _, conn := range p.Connection {
		wait.Add(1)

		go func(msg *gen.Message, conn *Connection) {
			defer wait.Done()
			if conn.active {
				err := conn.stream.Send(msg)
				fmt.Printf("Sending message to: %v from %v\n", conn.id, msg.ID)

				if err != nil {
					fmt.Printf("Error with Stream: %v - Error: %v\n", conn.stream, err)
					conn.active = false
					conn.error <- err
				}
			}
		}(msg, conn)
	}
	go func() {
		wait.Wait()
		close(done)
	}()
	<-done
	return &gen.Close{}, nil
}
