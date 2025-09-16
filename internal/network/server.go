package network

import (
	"bufio"
	"fmt"
	"go.uber.org/zap"
	"in-memory-kv/internal/compute"
	"in-memory-kv/internal/config"
	"net"
	"strings"
	"sync"
)

// Server – структура TCP-сервера базы данных.
// Она хранит:
// – конфигурацию,
// – ссылку на движок хранения,
// – семафор (канал) для ограничения числа клиентов,
// – waitGroup для корректного завершения всех горутин.
type Server struct {
	cfg     config.Config
	compute compute.Compute
	limiter chan struct{}
	wg      sync.WaitGroup
	log     *zap.Logger
}

func NewServer(cfg config.Config, compute compute.Compute, log *zap.Logger) *Server {
	return &Server{
		cfg:     cfg,
		compute: compute,
		limiter: make(chan struct{}, cfg.Network.MaxConnections),
		log:     log,
	}
}

// ListenAndServe запускает сервер и слушает входящие подключения.
func (s *Server) StartServer() error {
	ln, err := net.Listen("tcp", s.cfg.Network.Address)
	if err != nil {
		return err
	}
	fmt.Printf("Server listening on %s (max %d connections)\n",
		s.cfg.Network.Address, cap(s.limiter))

	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("Accept error:", err)
			continue
		}

		s.limiter <- struct{}{}

		s.wg.Add(1)
		go s.handleConnection(conn)
	}
}

// handleConnection обрабатывает каждое соединение в отдельной горутине.
func (s *Server) handleConnection(conn net.Conn) {
	defer func() {
		conn.Close()
		<-s.limiter
		s.wg.Done()
	}()

	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.TrimSpace(line) == "" {
			continue
		}
		if line == "EXIT" {
			conn.Write([]byte("BYE\n"))
			return
		}

		resp := s.compute.Handle(line)

		if !strings.HasSuffix(resp, "\n") {
			resp += "\n"
		}
		_, err := conn.Write([]byte(resp))
		if err != nil {
			s.log.Error("Write error", zap.Error(err))
			return
		}
	}

	if err := scanner.Err(); err != nil {
		s.log.Error("Scanner error", zap.Error(err))
	}
}

func (s *Server) Shutdown() {
	s.wg.Wait()
}
