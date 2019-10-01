package server

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"regexp"
	"sync"

	"github.com/gps/gps-tracker/server/workers"

	"github.com/gps/gps-tracker/models"

	"github.com/pkg/errors"

	"github.com/gps/gps-tracker/conf"
)

type Server struct {
	config       conf.Config
	MaxReadBytes int64

	listener   net.Listener
	conns      map[net.Conn]struct{}
	mu         sync.Mutex
	inShutdown bool
}

func New(config conf.Config) *Server {
	return &Server{
		config: config,
	}
}

func (srv *Server) ListenAndServe() error {
	config := srv.config
	log := config.Log()

	log.Info("starting server\n")

	listener, err := net.Listen("tcp", config.Tracker().Info())
	if err != nil {
		return err
	}

	defer listener.Close()

	srv.listener = listener

	for {
		if srv.inShutdown {
			break
		}

		conn, err := listener.Accept()
		if err != nil {
			log.Printf("error accepting connection %v", err)
			continue
		}

		log.Infof("accepted connection from %v", conn.RemoteAddr())

		srv.trackConn(conn)

		go func() {
			err := srv.handle(conn)
			if err != nil {
				log.Printf("error handling request %v", err)
			}
		}()
	}

	return nil
}

func (srv *Server) trackConn(c net.Conn) {
	defer srv.mu.Unlock()

	srv.mu.Lock()

	if srv.conns == nil {
		srv.conns = make(map[net.Conn]struct{})
	}

	srv.conns[c] = struct{}{}
}

func (srv *Server) handle(conn net.Conn) error {
	config := srv.config
	log := config.Log()
	db := config.DB()

	defer func() {
		log.Printf("closing connection from %v", conn.RemoteAddr())
		_ = conn.Close()
		srv.deleteConn(conn)
	}()

	r := bufio.NewReader(conn)
	w := bufio.NewWriter(conn)

	scanr := bufio.NewScanner(r)

	for {
		scanned := scanr.Scan()
		if !scanned {
			if err := scanr.Err(); err != nil {
				return err
			}
			break
		}

		fileWrite(scanr.Text())

		results, err := db.FindAllSupportedTrackers()
		if err != nil {
			return err
		}

		if results == nil {
			return errors.New("Please add supported trackers")
		}

		input := scanr.Text()

		for _, tracker := range results {
			ok, err := regexp.Match(tracker.Pattern, []byte(input))
			if err != nil {
				//need to log
				continue
			}

			if !ok {
				//need to log
				continue
			}

			worker := srv.getWorker(tracker)
			if worker == nil {
				//need to log
				continue
			}

			worker.Work(input, w)
		}
	}

	return nil
}

func (srv *Server) getWorker(tracker *models.SupportedTracker) workers.WorkerInterface {
	if tracker.Type == "chinese" {
		return workers.NewChinese(srv.config, tracker)
	}

	return nil
}

func (srv *Server) deleteConn(conn net.Conn) {
	defer srv.mu.Unlock()
	srv.mu.Lock()
	delete(srv.conns, conn)
}

func fileWrite(input string) {
	f, err := os.OpenFile("logs/tmp.txt", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		panic(err)
	}

	defer f.Close()

	if _, err = f.WriteString(fmt.Sprintf("%s\n", input)); err != nil {
		panic(err)
	}
}
