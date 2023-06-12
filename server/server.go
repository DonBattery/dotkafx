package server

import (
	"bufio"
	"fmt"
	"net"
	"strings"
	"time"

	"dotkafx/log"
	"dotkafx/model"
	"dotkafx/scheduler"
	"dotkafx/sound"
	"dotkafx/tools"
)

type Server struct {
	fx  *sound.Player
	sch *scheduler.Scheduler
	cmd model.RootCommand
}

func NewServer(fx *sound.Player, sch *scheduler.Scheduler, cmd model.RootCommand) *Server {
	return &Server{
		fx:  fx,
		sch: sch,
		cmd: cmd,
	}
}

func (srv *Server) handleConnection(conn net.Conn) {
	log.Debug("Received TCP connection. Local address: %s Remote address: %s", conn.LocalAddr(), conn.RemoteAddr())
	defer func() {
		if err := conn.Close(); err != nil {
			log.Error("Failed to close TCP connection properly: %s Local address: %s Remote address: %s", err, conn.LocalAddr(), conn.RemoteAddr())
		}
	}()

	request, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		log.Error("Failed to read client request: %s", err)
		return
	}
	request = strings.TrimSpace(request)
	log.Info("Request received: %s", request)

	response := ""

	switch {

	case request == "test":
		log.Debug("Testing Sound Output")
		srv.fx.Play(sound.ChaosDunk)
		response = "Test succeeded"

	case request == "start":
		response = srv.sch.Start()

	case request == "stop":
		response = srv.sch.Stop()

	case request == "pause":
		response = srv.sch.Pause()

	case request == "shutdown":
		srv.fx.Play(sound.DotkaFXServerIsShuttingDown)
		time.Sleep(time.Second * 3)
		log.Shutdown("gg wp")

	case strings.HasPrefix(request, "back"):
		amount, err := tools.ParseSuffixAmount(request, "back")
		if err != nil {
			response = fmt.Sprintf("Incorrect input value for backward seconds: %s", err)
		} else {
			response = srv.sch.Back(amount)
		}

	case strings.HasPrefix(request, "forward"):
		amount, err := tools.ParseSuffixAmount(request, "forward")
		if err != nil {
			response = fmt.Sprintf("Incorrect input value for forward seconds: %s", err)
		} else {
			response = srv.sch.Forward(amount)
		}

	default:
		response = fmt.Sprintf("Unknown command: %s Allowed commands: start, stop, pause, back[seconds], forward[seconds], shutdown", request)
	}

	log.Debug("Sending response: %s Local address: %s Remote address: %s", response, conn.LocalAddr(), conn.RemoteAddr())
	_, err = conn.Write([]byte(response + "\n"))
	if err != nil {
		log.Error("Failed to write response: %s Local address: %s Remote address: %s", err, conn.LocalAddr(), conn.RemoteAddr())
	}
}

func (srv *Server) soundPlayer() error {
	for {
		sound := <-srv.sch.EventChan
		srv.fx.Play(sound)
	}
}

func (srv *Server) Run() error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", srv.cmd.Port))
	if err != nil {
		return err
	}

	log.Info("DotkaFX server listening on TCP Port %d", srv.cmd.Port)

	go srv.soundPlayer()

	srv.fx.Play(sound.DotkaFXSercerIsOnline)

	for {
		conn, err := lis.Accept()
		if err != nil {
			return err
		}
		go srv.handleConnection(conn)
	}
}
