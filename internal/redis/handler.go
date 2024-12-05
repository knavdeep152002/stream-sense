package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os/exec"

	"path/filepath"

	"github.com/knavdeep152002/stream-sense/internal/utils"
	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()
var RedisClient = redis.NewClient(&redis.Options{
	Addr: "localhost:6379",
})

func Publish(channel string, message []byte) error {
	return RedisClient.Publish(ctx, channel, message).Err()
}

func Subscribe(channel string) *redis.PubSub {
	pubSub := RedisClient.Subscribe(ctx, channel)
	for {
		msg, err := pubSub.ReceiveMessage(ctx)
		if err != nil {
			panic(err)
		}
		fmt.Println(msg.Channel, msg.Payload)
	}
}

func Observe() {
	log.Println("Observing")
	pubSub := RedisClient.PSubscribe(ctx, "ss-*")
	log.Println("Subscribed")
	for {
		msg, err := pubSub.ReceiveMessage(ctx)
		if err != nil {
			panic(err)
		}
		switch msg.Channel {
		case "ss-audio-preprocess":
			//Convert the chunk to wav using ffmpeg
			log.Println("Received ss-audio-preprocess : ", msg.Payload)
			chunkData := msg.Payload
			chunk := &utils.ChunkReceiver{}
			err := json.Unmarshal([]byte(chunkData), chunk)
			if err != nil {
				log.Println("Failed to unmarshal chunk")
				continue
			}
			// send file to generate segments
			go utils.GenerateSegments(filepath.Join(chunk.UploadDir, chunk.Filename), chunk.Filename, chunk.VideoID)
			// read the content and send to ffmpeg
			err = utils.ConvertMP4ToWav(filepath.Join(chunk.UploadDir, chunk.Filename), chunk.Filename)
			if err != nil {
				log.Println("Failed to convert mp4 to wav", err)
				continue
			}

			data := &utils.TranscribeArgs{
				AudioPath: filepath.Join(utils.WavsDir, chunk.Filename+".wav"),
				VideoId:   chunk.Filename,
			}

			dataBytes, err := json.Marshal(data)
			if err != nil {
				log.Println("Failed to marshal transcribe args")
				continue
			}

			err = Publish("ss-transcribe", dataBytes)
			if err != nil {
				log.Println("Failed to publish transcribe message")
				continue
			}

		case "ss-transcribe":
			args := &utils.TranscribeArgs{}
			err := json.Unmarshal([]byte(msg.Payload), args)
			if err != nil {
				log.Println("Failed to unmarshal transcribe args")
				continue
			}
			log.Println("Received transcribe : ", msg.Payload)
			command := "source venv/bin/activate && python3 transcribe.py --input_audio " + args.AudioPath + " --vid_id " + args.VideoId
			go func() {
				e := exec.Command("/bin/sh", "-c", command)
				x, err := e.CombinedOutput()
				if err != nil {
					log.Println("Failed to transcribe", err)
				}
				log.Println("Output: ", string(x))
			}()
		default:
			log.Println("Received unknown message")
		}
	}
}
