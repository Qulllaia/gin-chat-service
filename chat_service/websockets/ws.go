package websockets

import (
	"main/controller/utils"
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

type UserWSData struct {
	user_id int;
	ws *websocket.Conn
}

var chat_id_conns = map[string][]UserWSData{} 

func WebsocketsInit(context *gin.Context) {
	conn, err := upgrader.Upgrade(context.Writer, context.Request, nil);

	if err != nil {
		println("Error Connection", err.Error())
		return;
	}

	cookie := context.Request.Cookies();

	jwt_token := "";

	for _, val := range cookie {
		if val.Name == "session_token" {
			jwt_token = val.Value;
		}
	}

	claims, err := utils.DecodeJWT(jwt_token);

	if err != nil {
		println("Error JWT", err.Error())
		return;
	}


	chat_id, _ := context.GetQuery("chat_id")
	chat_id_conns[chat_id] = append(chat_id_conns[chat_id], UserWSData{user_id: claims.UserID, ws:conn});

	
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
				if getter_conn.ws != conn {
					err  = getter_conn.ws.WriteMessage(messageType, message);
					println(chat_id, getter_conn.user_id)
				}
		
				if err != nil {
					println("Send Error", err.Error(), len(getter_conns));
				}
			}
		}
	}
}