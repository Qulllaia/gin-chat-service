package websockets

import (
	"encoding/json"
	. "main/types"
	"main/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
);

type WSConnection struct {
	Actor *ConnectorActor
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var user_id_conns_with_chat_ids = map[int]UserWSData{} 


func (ws *WSConnection) WebsocketsInit(context *gin.Context) {
	conn, err := upgrader.Upgrade(context.Writer, context.Request, nil);
	
	if err != nil {
		println("Error Connection", err.Error())
		return;
	}

	defer conn.Close();
	
	claims, err := utils.ExtractClaimsFromCookie(context);

	if err != nil {
		println("Error JWT", err.Error())
		return;
	}
	ws.Actor.AddClient(conn, claims.UserID)

	for {

		messageType, message_ws, err := conn.ReadMessage();

		if err != nil {
			println("Error Message", err.Error());
			break;
		}

		var message MessageWS;

		err = json.Unmarshal(message_ws, &message)

		ws.Actor.Send(message, messageType, conn);
	}
}


