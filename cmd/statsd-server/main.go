package statsd_server

import (
	"fmt"
	"github.com/caarlos0/env"
	_ "github.com/joho/godotenv/autoload" // by default we will read .env before the real env vars
	"github.com/mitchellh/cli"
	"github.com/yusufsyaifudin/statsd-server/internal/statsdserver"
	"log"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

const (
	TCP = "tcp"
	UDP = "udp"
)

type Config struct {
	Port int `env:"PORT" envDefault:"8125" validate:"required"`
}

type StatsdServer struct {
	Config Config
}

var _ cli.Command = (*StatsdServer)(nil)

func New(cfg Config) (*StatsdServer, error) {
	return &StatsdServer{
		Config: cfg,
	}, nil
}

func (s *StatsdServer) Help() string {
	return "running UDP server for statsd"
}

func (s *StatsdServer) Run(args []string) int {
	// *** Parse and validate config input
	cfg := Config{}
	err := env.Parse(&cfg)
	if err != nil {
		err = fmt.Errorf("cannot parse env var: %w", err)
		log.Println(err)
		return 1
	}

	metricsAddr := fmt.Sprintf(":%d", cfg.Port)
	resolveUDPAddr, err := net.ResolveUDPAddr(UDP, metricsAddr)
	if err != nil {
		err = fmt.Errorf("cannot parse udp address: %w", err)
		log.Println(err)
		return 1
	}

	if resolveUDPAddr == nil {
		resolveUDPAddr = &net.UDPAddr{
			IP:   net.IPv6unspecified,
			Port: cfg.Port,
			Zone: "",
		}
	}

	if resolveUDPAddr.IP.Equal(net.IP{}) {
		metricsAddr = fmt.Sprintf("[%s]%s", net.IPv6unspecified, metricsAddr)
	}

	var errChan = make(chan error)
	var packetConn net.PacketConn
	wgServer := sync.WaitGroup{}
	wgServer.Add(1)
	go func() {
		defer wgServer.Done()

		log.Printf("running %s server on address %s\n", UDP, metricsAddr)
		var _err error
		packetConn, _err = net.ListenPacket(UDP, metricsAddr)
		if err != nil {
			err = fmt.Errorf("cannot listen packet address: %w", err)
			errChan <- _err
			return
		}
	}()

	// Wait server to be really up
	wgServer.Wait()

	statsdServerMsgHandlerCfg := statsdserver.ListenerCfg{
		PacketConn:   packetConn,
		FallbackAddr: resolveUDPAddr,
	}
	statsdServerMsgHandler, err := statsdserver.NewListener(statsdServerMsgHandlerCfg)
	if err != nil {
		err = fmt.Errorf("failed to prepare message handler: %w", err)
		log.Println(err)
		return 1
	}

	statsdServerMsgHandler.Run()
	log.Println("server is up and running")

	var signalChan = make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)
	select {
	case s := <-signalChan:
		msg := fmt.Sprintf("got an interrupt: %+v", s)

		if packetConn != nil {
			err = packetConn.Close()
			if err != nil {
				msg = fmt.Sprintf("%s: connection close with error: %s", msg, err)
			}
		}

		log.Println(msg)

	case _err := <-errChan:
		if _err != nil {
			msg := fmt.Sprintf("error while running server: %s", _err)
			log.Println(msg)
		}
	}

	return 0
}

func (s *StatsdServer) Synopsis() string {
	return "running UDP server for statsd"
}
