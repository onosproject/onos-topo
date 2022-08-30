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
	id     uint32
	ch     chan *topo.WatchResponse
	ws     *websocket.Conn
	client topo.TopoClient
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
	return http.ListenAndServe(":12345", nil)
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
		ws: ws,
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
	wc.client = topo.NewTopoClient(s.topoConn)
	ctx := context.Background()
	stream, err := wc.client.Watch(ctx, &topo.WatchRequest{})
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
	_ = homeTemplate.Execute(w, "ws://"+r.Host+"/watch")
}

var homeTemplate = template.Must(template.New("index").Parse(`
<!DOCTYPE html>
<html>
<head>
<meta charset="utf-8">
    <style> body { margin: 0; } </style>
    <script src="//unpkg.com/force-graph"></script>
</head>
<body>
  	<div id="graph"></div>
    <script>  
		const canvas = document.getElementById("graph");
		var nodes = [], links = [];

		const Graph = ForceGraph()(canvas)
			//.onNodeClick(inspectNode)
			.graphData({nodes: nodes, links: links})
 			.nodeCanvasObject((node, ctx, globalScale) => {
				  const label = node.id;
				  const fontSize = 12/globalScale;
				  ctx.font = '' + fontSize + 'px Sans-Serif';
				  const textWidth = ctx.measureText(label).width;
				  const bckgDimensions = [textWidth, fontSize].map(n => n + fontSize * 0.2); // some padding
	
				  ctx.fillStyle = 'rgba(255, 255, 255, 0.8)';
				  ctx.fillRect(node.x - bckgDimensions[0] / 2, node.y - bckgDimensions[1] / 2, ...bckgDimensions);
	
				  ctx.textAlign = 'center';
				  ctx.textBaseline = 'middle';
				  ctx.fillStyle = '#000';
				  ctx.fillText(label, node.x, node.y);
	
				  node.__bckgDimensions = bckgDimensions; // to re-use in nodePointerAreaPaint
			})
     		.linkDirectionalArrowLength(5)
	        .linkDirectionalArrowRelPos(1)
			.linkCanvasObjectMode(() => 'after')
    		.linkCanvasObject((link, ctx, globalScale) => {
    			const MAX_FONT_SIZE = 12/globalScale;
    			const LABEL_NODE_MARGIN = Graph.nodeRelSize() * 1.5;
    			
    			const start = link.source;
    			const end = link.target;
    			
    			// ignore unbound links
    			if (typeof start !== 'object' || typeof end !== 'object') return;
    			
    			// calculate label positioning
    			const textPos = Object.assign(...['x', 'y'].map(c => ({
    			  [c]: start[c] + (end[c] - start[c]) / 2 // calc middle point
    			})));
    			
    			const relLink = { x: end.x - start.x, y: end.y - start.y };
    			
    			const maxTextLength = Math.sqrt(Math.pow(relLink.x, 2) + Math.pow(relLink.y, 2)) - LABEL_NODE_MARGIN * 2;
    			
    			let textAngle = Math.atan2(relLink.y, relLink.x);
    			// maintain label vertical orientation for legibility
    			if (textAngle > Math.PI / 2) textAngle = -(Math.PI - textAngle);
    			if (textAngle < -Math.PI / 2) textAngle = -(-Math.PI - textAngle);
    			
    			const label = link.data.id;

    			// estimate fontSize to fit in link length
    			ctx.font = '1px Sans-Serif';
    			const fontSize = Math.min(MAX_FONT_SIZE, maxTextLength / ctx.measureText(label).width);
    			ctx.font = '' + fontSize + 'px Sans-Serif';
    			const textWidth = ctx.measureText(label).width;
    			const bckgDimensions = [textWidth, fontSize].map(n => n + fontSize * 0.2); // some padding

    			// draw text label (with background rect)
    			ctx.save();
    			ctx.translate(textPos.x, textPos.y);
    			ctx.rotate(textAngle);

    			ctx.fillStyle = 'rgba(255, 255, 255, 0.8)';
    			ctx.fillRect(- bckgDimensions[0] / 2, - bckgDimensions[1] / 2, ...bckgDimensions);

    			ctx.textAlign = 'center';
    			ctx.textBaseline = 'middle';
    			ctx.fillStyle = 'darkgrey';
    			ctx.fillText(label, 0, 0);
    			ctx.restore();
    		});

		var ws = new WebSocket("{{.}}");
		ws.onopen = function(evt) {
			console.log("Connected");
		}
		ws.onclose = function(evt) {
			console.log("Disconnected");
			ws = null;
		}
		ws.onmessage = function(evt) {
			//print("Event: " + evt.data);
			data = JSON.parse(evt.data);
			console.log(data);
			if (data.event === "replay" || data.event === "added") {
				if (data.entity) {
					nodes = [...nodes, { id: data.id, data: data }];
				} else if (data.relation) {
					links = [...links, { source: data.relation.src, target: data.relation.tgt, data: data }];
				}
				Graph.graphData({nodes: nodes, links: links})
			}
		}
		ws.onerror = function(evt) {
			console.error("ERROR: " + evt.data);
		}
    </script>
</body>
</html>
`))
