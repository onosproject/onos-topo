// SPDX-FileCopyrightText: 2022-present Intel Corporation
//
// SPDX-License-Identifier: Apache-2.0

package visualizer

import (
	"context"
	"github.com/gorilla/websocket"
	"github.com/onosproject/onos-api/go/onos/topo"
	"github.com/onosproject/onos-lib-go/pkg/logging"
	"google.golang.org/grpc"
	"html/template"
	"io"
	"net/http"
	"sync"
)

var (
	log      = logging.GetLogger("visualizer")
	upgrader = websocket.Upgrader{}
	maxID    = uint32(0)
)

type webClient struct {
	id uint32
	ch chan *topo.WatchResponse
}

// Server is an HTTP/WS server for the web-based visualizer client
type Server struct {
	topoConn *grpc.ClientConn
	clients  map[uint32]*webClient
	lock     sync.RWMutex
}

// NewServer creates a new HTTP/WS server for the web-based visualizer client
func NewServer(conn *grpc.ClientConn) *Server {
	return &Server{
		topoConn: conn,
		clients:  make(map[uint32]*webClient),
	}
}

// Serve starts the HTTP/WS server
func (s *Server) Serve() error {
	http.HandleFunc("/watch", s.watchChanges)
	http.HandleFunc("/", s.home)
	return http.ListenAndServe(":5152", nil)
}

func (s *Server) watchChanges(w http.ResponseWriter, r *http.Request) {
	log.Infof("Received new web client connection")
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Error("Unable to upgrade HTTP connection:", err)
		return
	}
	defer ws.Close()

	s.lock.Lock()
	defer s.lock.Unlock()
	maxID++
	newClient := &webClient{
		id: maxID,
		ch: make(chan *topo.WatchResponse),
	}
	s.clients[newClient.id] = newClient
	s.lock.Unlock()
	log.Infof("Web client %d connected", newClient.id)

	go s.watchTopology(newClient)

	for msg := range newClient.ch {
		b, err := EncodeTopoEvent(msg)
		if err != nil {
			log.Warn("Unable to encode:", err)
		}
		err = ws.WriteMessage(websocket.TextMessage, b)
		if err != nil {
			log.Errorf("Unable to write:", err)
			break
		}
	}
	log.Infof("Web client %d disconnected")
}

func (s *Server) watchTopology(wc *webClient) {
	client := topo.NewTopoClient(s.topoConn)
	ctx := context.Background()
	stream, err := client.Watch(ctx, &topo.WatchRequest{})
	if err != nil {
		log.Errorf("Unable to connect to onos-topo: %+v", err)
		return
	}

	for {
		msg, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Errorf("Unable read from onos-topo stream: %+v", err)
			break
		}
		wc.ch <- msg
	}

	s.lock.Lock()
	defer s.lock.Unlock()
	delete(s.clients, wc.id)
	log.Infof("Web client %d disconnected")
}

func (s *Server) home(w http.ResponseWriter, r *http.Request) {
	homeTemplate, err := template.New("index.html").ParseFiles("pkg/tools/topo-visualizer/index.html")
	if err != nil {
		log.Errorf("Unable to parse template: %+v", err)
		return
	}
	_ = homeTemplate.Execute(w, "ws://"+r.Host+"/watch")
}
