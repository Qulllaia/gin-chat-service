package handlers

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"

	"main/types"

	"github.com/gorilla/websocket"
	"github.com/hcl/audioduration"
)

type MediaHandler struct{}

type ChuncksCache struct {
	Chunks []Chunk
}

type Chunk struct {
	Bytes []byte
	Index int
}

func NewMediaHandler() types.Handler {
	return &MediaHandler{}
}

const STATIC_BUFFER_READ_LIMIT = 1024 * 100

func (mh *MediaHandler) Handle(message types.MessageWS, messageType int, conn *websocket.Conn, actor types.Actor) {
	index := message.Index

	filename := "example.mp3"
	filePath := "./" + filename

	// rdb := actor.GetRDB()

	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	defer file.Close()

	buffer := make([]byte, STATIC_BUFFER_READ_LIMIT)
	chunk, err := file.ReadAt(buffer, int64(STATIC_BUFFER_READ_LIMIT*index))
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	size, err := audioduration.Mp3(file)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fileInfo, err := file.Stat()
	if err != nil {
		fmt.Println("Ошибка получения информации:", err)
		return
	}

	responseData, err := json.Marshal(map[string]interface{}{
		"type":       "audio_chunk",
		"data":       base64.StdEncoding.EncodeToString(buffer[:chunk]),
		"size":       size,
		"size_bytes": fileInfo.Size(),
	})
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	conn.WriteMessage(messageType, []byte(responseData))
}
