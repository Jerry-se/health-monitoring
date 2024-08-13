package ws

import (
	"encoding/json"
	"net/http"
	"time"

	"health-monitoring/log"
	"health-monitoring/types"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		if r.Method != "GET" {
			return false
		}
		if r.URL.Path != "/echo" && r.URL.Path != "/websocket" {
			return false
		}
		return true
	},
} // use default options

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 30 * time.Second // 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10
)

func Ws(ctx *gin.Context) {
	w, r := ctx.Writer, ctx.Request
	var nodeId string
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, "Upgrade to websocket failed", http.StatusUpgradeRequired)
		log.Log.Error("Upgrade to websocket failed:", err)
		return
	}
	defer func() {
		log.Log.WithFields(logrus.Fields{
			"node_id": nodeId,
		}).Info("connection stopped")
		c.Close()
	}()

	c.SetReadDeadline(time.Now().Add(pongWait))
	c.SetPingHandler(func(appData string) error {
		c.SetReadDeadline(time.Now().Add(pongWait))
		log.Log.WithFields(logrus.Fields{
			"node_id": nodeId,
		}).Info("ping handler")
		return nil
	})

	for {
		mt, message, err := c.ReadMessage()
		if err != nil {
			log.Log.WithFields(logrus.Fields{
				"node_id": nodeId,
			}).Info("read:", err)
			break
		}
		log.Log.WithFields(logrus.Fields{
			"node_id": nodeId,
		}).Infof("recv message: %v %s", mt, message)

		req := &types.WsRequest{}
		if err := json.Unmarshal(message, req); err != nil {
			log.Log.WithFields(logrus.Fields{
				"node_id": nodeId,
			}).Error("parse request failed", err)
			writeWsResponse(c, nodeId, &types.WsResponse{
				WsHeader: types.WsHeader{
					Version:   0,
					Timestamp: time.Now().Unix(),
					Id:        0,
					Type:      0,
					PubKey:    []byte(""),
					Sign:      []byte(""),
				},
				Code:    uint32(types.ErrCodeParam),
				Message: "parse request failed",
				Body:    []byte(""),
			})
			continue
		}

		switch req.Type {
		case uint32(types.WsMtOnline):
			onlineReq := &types.WsOnlineRequest{}
			if err := json.Unmarshal(req.Body, onlineReq); err != nil {
				log.Log.WithFields(logrus.Fields{
					"node_id": nodeId,
				}).Error("parse online request failed", err)
				writeWsResponse(c, nodeId, &types.WsResponse{
					WsHeader: types.WsHeader{
						Version:   0,
						Timestamp: time.Now().Unix(),
						Id:        req.Id,
						Type:      req.Type,
						PubKey:    []byte(""),
						Sign:      []byte(""),
					},
					Code:    uint32(types.ErrCodeParam),
					Message: "parse online request failed",
					Body:    []byte(""),
				})
				continue
			}
			nodeId = onlineReq.NodeId
			handleWsOnlineRequest(c, nodeId, req, onlineReq)
		case uint32(types.WsMtMachineInfo):
			handleWsMachineInfoRequest(c, nodeId, req)
		default:
			log.Log.WithFields(logrus.Fields{
				"node_id": nodeId,
			}).Error("unknowned request message type")
			writeWsResponse(c, nodeId, &types.WsResponse{
				WsHeader: types.WsHeader{
					Version:   0,
					Timestamp: time.Now().Unix(),
					Id:        req.Id,
					Type:      req.Type,
					PubKey:    []byte(""),
					Sign:      []byte(""),
				},
				Code:    uint32(types.ErrCodeParam),
				Message: "unknowned request message type",
				Body:    []byte(""),
			})
		}

	}
}

func writeWsResponse(c *websocket.Conn, nodeId string, res *types.WsResponse) error {
	resBytes, err := json.Marshal(res)
	if err != nil {
		log.Log.WithFields(logrus.Fields{
			"node_id": nodeId,
		}).Error("marshal reponse failed", err)
		return err
	}
	err = c.WriteMessage(websocket.TextMessage, resBytes)
	if err != nil {
		log.Log.WithFields(logrus.Fields{
			"node_id": nodeId,
		}).Error("write response message failed", err)
		return err
	}
	return nil
}

func handleWsOnlineRequest(c *websocket.Conn, nodeId string, req *types.WsRequest, online *types.WsOnlineRequest) error {
	writeWsResponse(c, nodeId, &types.WsResponse{
		WsHeader: types.WsHeader{
			Version:   0,
			Timestamp: time.Now().Unix(),
			Id:        req.Id,
			Type:      req.Type,
			PubKey:    []byte(""),
			Sign:      []byte(""),
		},
		Code:    0,
		Message: "ok",
		Body:    []byte(""),
	})
	return nil
}

func handleWsMachineInfoRequest(c *websocket.Conn, nodeId string, req *types.WsRequest) error {
	miReq := &types.WsMachineInfoRequest{}
	if err := json.Unmarshal(req.Body, miReq); err != nil {
		log.Log.WithFields(logrus.Fields{
			"node_id": nodeId,
		}).Error("parse machine info request failed", err)
		writeWsResponse(c, nodeId, &types.WsResponse{
			WsHeader: types.WsHeader{
				Version:   0,
				Timestamp: time.Now().Unix(),
				Id:        req.Id,
				Type:      req.Type,
				PubKey:    []byte(""),
				Sign:      []byte(""),
			},
			Code:    uint32(types.ErrCodeParam),
			Message: "parse machine info request failed",
			Body:    []byte(""),
		})
		return err
	}
	log.Log.WithFields(logrus.Fields{
		"node_id": nodeId,
	}).WithField("machine info", miReq).Info("update machine info")
	writeWsResponse(c, nodeId, &types.WsResponse{
		WsHeader: types.WsHeader{
			Version:   0,
			Timestamp: time.Now().Unix(),
			Id:        req.Id,
			Type:      req.Type,
			PubKey:    []byte(""),
			Sign:      []byte(""),
		},
		Code:    0,
		Message: "ok",
		Body:    []byte(""),
	})
	return nil
}
