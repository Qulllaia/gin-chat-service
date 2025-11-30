package handlers

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"os"

	"main/types"

	"github.com/gorilla/websocket"
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

func (mh *MediaHandler) Handle(message types.MessageWS, messageType int, conn *websocket.Conn, actor types.Actor) {
	index := message.Index

	filename := "example.mp3"
	filePath := "./" + filename

	rdb := actor.GetRDB()

	check, err := rdb.DoesDataExists("1")
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	if *check == 1 {
		var chunks ChuncksCache
		redisData, _ := rdb.GetData("1")
		err := json.Unmarshal([]byte(redisData), &chunks)
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		for _, cacheChunk := range chunks.Chunks {
			if cacheChunk.Index == index {
				responseData, err := json.Marshal(map[string]interface{}{
					"type":  "audio_chunk",
					"chunk": 0,
					"data":  base64.StdEncoding.EncodeToString(cacheChunk.Bytes),
					"total": 1,
				})
				if err != nil {
					fmt.Println(err.Error())
					return
				}
				conn.WriteMessage(messageType, []byte(responseData))
			}
		}

	} else {

		file, err := os.Open(filePath)
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		defer file.Close()

		buffer := make([]byte, 1024*1024)
		chunkNumber := 0
		chunk := 0
		for {
			n, err := file.Read(buffer)
			if err == io.EOF {
				break
			}

			chunkNumber++
			if chunkNumber > 0 {
				chunk = n
				partedChunk := chunk / 3

				chucnkCache := ChuncksCache{
					Chunks: []Chunk{
						{
							Bytes: buffer[partedChunk : partedChunk*2],
							Index: 1,
						},
						{
							Bytes: buffer[partedChunk*2 : partedChunk*3],
							Index: 2,
						},
					},
				}
				chunkMarshal, _ := json.Marshal(chucnkCache)
				err = rdb.SetData("1", string(chunkMarshal))
				if err != nil {
					fmt.Println(err.Error())
					return
				}
				break
			}
		}

		responseData, err := json.Marshal(map[string]interface{}{
			"type":  "audio_chunk",
			"chunk": 0,
			"data":  base64.StdEncoding.EncodeToString(buffer[:chunk/3]),
			"total": chunkNumber,
		})
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		conn.WriteMessage(messageType, []byte(responseData))
	}
}
