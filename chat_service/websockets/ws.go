package websockets

import (
	"encoding/json"
	"main/controller/utils"
	. "main/database/models"
	"main/database/queries"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/lib/pq"
);

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type UserWSData struct {
	chat_ids pq.Int64Array;
	ws *websocket.Conn
}

type WSConnection struct {
	WSQ *queries.WSQueries
}

type MessageWS struct {
	Chat_id string `json:"chat_id"`;
	Message string `json:"messages"`;
}

var user_id_conns_with_chat_ids = map[int]UserWSData{} 


// TODO: Перенести нахуй sql запросы в query и нормально поделить код, типы и прочее
func (ws *WSConnection) WebsocketsInit(context *gin.Context) {
	conn, err := upgrader.Upgrade(context.Writer, context.Request, nil);

	if err != nil {
		println("Error Connection", err.Error())
		return;
	}

	claims, err := utils.ExtractClaimsFromCookie(context);

	if err != nil {
		println("Error JWT", err.Error())
		return;
	}
    var user User

	if err := ws.WSQ.GetUserChatList(claims, &user); err != nil {
		println("GetUserChatListError", err.Error())
	}
    
	user_id_conns_with_chat_ids[claims.UserID] = 
			UserWSData{chat_ids: user.Chat_list, ws:conn};

	
	defer conn.Close();

	for {

		messageType, message_ws, err := conn.ReadMessage();

		if err != nil {
			println("Error Message", err.Error());
			break;
		}

		var message MessageWS;

		err = json.Unmarshal(message_ws, &message)

		var getter_conn *websocket.Conn;
		for user_id, value := range user_id_conns_with_chat_ids {
			chat_id, _ := strconv.Atoi(message.Chat_id); 

			if contains(value.chat_ids, int64(chat_id)) {
				getter_conn = value.ws;
				if getter_conn != conn {
					err  = getter_conn.WriteMessage(messageType, []byte(message.Message));
				} else {
			
					if err = ws.WSQ.InsertMessageIntoChatHistory(chat_id, user_id, string(message.Message)); err != nil {
						println(err.Error());	
					}
				}
			}
		}

	}
}


func contains(slice pq.Int64Array, item int64) bool {
    for _, element := range slice {
        if element == item {
            return true
        }
    }
    return false
}