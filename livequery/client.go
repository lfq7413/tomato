package livequery

import "encoding/json"

type client struct {
	id int
	ws *webSocket
}

func newClient(id int, ws *webSocket) *client {
	return &client{
		id: id,
		ws: ws,
	}
}

func pushResponse(ws *webSocket, msg string) {
	ws.send(msg)
}

func pushError(ws *webSocket, code int, errMsg string, reconnect bool) {
	errResp := M{
		"op":        "error",
		"error":     errMsg,
		"code":      code,
		"reconnect": reconnect,
	}
	data, err := json.Marshal(errResp)
	if err != nil {
		return
	}
	pushResponse(ws, string(data))
}
