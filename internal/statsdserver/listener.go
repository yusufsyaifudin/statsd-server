package statsdserver

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/yusufsyaifudin/statsd-server/pkg/parser"
	"log"
	"net"
)

type ListenerCfg struct {
	PacketConn   net.PacketConn
	FallbackAddr net.Addr
}

type Listener struct {
	Config ListenerCfg
}

func NewListener(cfg ListenerCfg) (*Listener, error) {
	return &Listener{Config: cfg}, nil
}

// Run TODO handle
//
//	panic: too many concurrent operations on a single file or socket (max 1048575)
func (l *Listener) Run() {
	listener := l.Config.PacketConn
	for {
		message := make([]byte, 512)
		n, addr, err := listener.ReadFrom(message)
		if err != nil {
			continue
		}

		if addr == nil {
			addr = l.Config.FallbackAddr
		}

		buf := bytes.NewBuffer(message[0:n])
		go l.handleMessage(addr, buf)
	}
}

// handleMessage will receive the message, the basic example is on here https://github.com/yusufsyaifudin/gographite
func (l *Listener) handleMessage(addr net.Addr, buf *bytes.Buffer) {
	// parse message
	msg := buf.String()

	// ** handle multi-packet, separated with newline
	metrics, err := parser.Parse(msg)
	if err != nil {
		fmt.Println(err)
		return
	}

	for idx, metric := range metrics {
		b, err := json.Marshal(metric)
		if err != nil {
			err = fmt.Errorf("cannot parse metric idx %d to json: %w", idx, err)
			log.Println(err)

			fmt.Printf("%+v\n", metric)
			continue
		}

		fmt.Println(string(b))
	}

}
