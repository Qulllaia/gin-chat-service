package websockets

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
);

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var chat_id_conns = map[string][]*websocket.Conn{} 

func WebsocketsInit(context *gin.Context) {
	conn, err := upgrader.Upgrade(context.Writer, context.Request, nil);

	chat_id, _ := context.GetQuery("chat_id")
	chat_id_conns[chat_id] = append(chat_id_conns[chat_id], conn);
	if err != nil {
		println("Error Connection", err.Error())
		return;
	}
	
	defer conn.Close();

	for {

		messageType, message, err := conn.ReadMessage();

		if err != nil {
			println("Error Message", err.Error());
			break;
		}

		chat_id, chat_id_get_successfull := context.GetQuery("chat_id")
		if(chat_id_get_successfull) {

			getter_conns := chat_id_conns[chat_id]
			
			for _, getter_conn := range getter_conns {
				if getter_conn != conn {
					err  = getter_conn.WriteMessage(messageType, message);
				}
		
				if err != nil {
					println("Send Error", err.Error(), len(getter_conns));
				}
			}
		}
	}
}