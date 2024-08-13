package ws

import (
	"encoding/json"
	"health-monitoring/types"
	"testing"
	"time"

	"github.com/gorilla/websocket"
)

func TestWsMachine(t *testing.T) {
	url := "ws://localhost:9521/websocket"
	c, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		t.Fatalf("dial %s: %v", url, err)
	}
	defer c.Close()

	if err := c.WriteControl(websocket.PingMessage, []byte{}, time.Now().Add(9*time.Second)); err != nil {
		t.Fatalf("ping websocket failed: %v", err)
	}

	done := make(chan struct{})
	go func() {
		defer close(done)
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				t.Log("read:", err)
				return
			}
			t.Logf("recv: %s", message)

			response := &types.WsResponse{}
			if err := json.Unmarshal(message, response); err != nil {
				println("parse websocket response failed:", err)
				break
			}
			if response.Type == uint32(types.WsMtMachineInfo) && response.Code == 0 {
				println("received machine info success")
				break
			}
		}
	}()

	var reqId uint64 = 0

	onlineReq := &types.WsOnlineRequest{
		NodeId: "123456789",
	}
	reqBody, err := json.Marshal(onlineReq)
	if err != nil {
		t.Fatalf("marshal online request body failed: %v", err)
	}
	req := &types.WsRequest{
		WsHeader: types.WsHeader{
			Version:   0,
			Timestamp: time.Now().Unix(),
			Id:        reqId,
			Type:      uint32(types.WsMtOnline),
			PubKey:    []byte(""),
			Sign:      []byte(""),
		},
		Body: reqBody,
	}
	reqBytes, err := json.Marshal(req)
	if err != nil {
		t.Fatalf("marshal online request failed: %v", err)
	}
	if err := c.WriteMessage(websocket.TextMessage, reqBytes); err != nil {
		t.Fatalf("send websocket message failed: %v", err)
	}

	reqId++

	machineInfo := &types.WsMachineInfoRequest{
		Project:        "DecentralGPT",
		Models:         make([]types.ModelInfo, 0),
		GPUName:        "NVIDIA RTX A5000",
		UtilizationGPU: 30,
		MemoryTotal:    24564,
		MemoryUsed:     22128,
	}
	machineInfo.Models = append(machineInfo.Models, types.ModelInfo{
		Model: "Codestral-22B-v0.1",
	})
	reqBody, err = json.Marshal(machineInfo)
	if err != nil {
		t.Fatalf("marshal machine info request body failed: %v", err)
	}
	req2 := &types.WsRequest{
		WsHeader: types.WsHeader{
			Version:   0,
			Timestamp: time.Now().Unix(),
			Id:        reqId,
			Type:      uint32(types.WsMtMachineInfo),
			PubKey:    []byte(""),
			Sign:      []byte(""),
		},
		Body: reqBody,
	}
	reqBytes, err = json.Marshal(req2)
	if err != nil {
		t.Fatalf("marshal machine info request failed: %v", err)
	}
	if err := c.WriteMessage(websocket.TextMessage, reqBytes); err != nil {
		t.Fatalf("send websocket message failed: %v", err)
	}

	<-done

	err = c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
	if err != nil {
		t.Log("write close:", err)
		return
	}
	<-time.After(time.Second)
}
