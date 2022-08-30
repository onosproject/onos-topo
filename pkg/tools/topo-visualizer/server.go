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
	"os"
	"sync"
)

const (
	devIndexPath = "pkg/tools/topo-visualizer/index.html"
	indexPath    = "/var/topo-visualizer/index.html"
)

var (
	log      = logging.GetLogger("visualizer")
	upgrader = websocket.Upgrader{}
	maxID    = uint32(0)
	devMode  = false
)

var homeTemplate *template.Template

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
	loadTemplate()
	http.HandleFunc("/watch", s.watchChanges)
	http.HandleFunc("/", s.home)
	return http.ListenAndServe(":5152", nil)
}

// Load index.html template from its expected production location; sets devMode if not found
func loadTemplate() {
	if _, err := os.Stat(indexPath); err == nil {
		if homeTemplate, err = template.New("index.html").ParseFiles(indexPath); err != nil {
			log.Errorf("Unable to parse template %s: %+v", indexPath, err)
		}
		log.Infof("Running with production template...")
		return
	}
	log.Infof("Running in dev mode...")
	devMode = true
}

// Sets up a new web-client and Watch topology changes and relay them onto the web-client
func (s *Server) watchChanges(w http.ResponseWriter, r *http.Request) {
	log.Infof("Received new web client connection")
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Error("Unable to upgrade HTTP connection:", err)
		return
	}
	defer ws.Close()

	s.lock.Lock()
	maxID++
	wc := &webClient{
		id: maxID,
		ch: make(chan *topo.WatchResponse),
	}
	s.clients[wc.id] = wc
	s.lock.Unlock()
	log.Infof("Web client %d connected", wc.id)

	go s.watchTopology(wc)

	for msg := range wc.ch {
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
	log.Infof("Web client %d disconnected", wc.id)
}

// Watch topology changes and relay them onto the web-client
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
	if devMode {
		var err error
		if homeTemplate, err = template.New("index.html").ParseFiles(devIndexPath); err != nil {
			log.Errorf("Unable to parse template %s: %+v", devIndexPath, err)
			return
		}
	}
	_ = homeTemplate.Execute(w, "ws://"+r.Host+"/watch")
}
