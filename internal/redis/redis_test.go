package redis_test

import (
	"context"
	"log"
	"os/exec"
	"testing"

	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()
var rdb = redis.NewClient(&redis.Options{
	Addr:     "localhost:6379",
	Password: "", // no password set
	DB:       0,  // use default DB
})

func Test(t *testing.T) {
	command := "source venv/bin/activate && python3 transcribe.py --input_audio " + "wavs/test.wav" + " --vid_id " + "1"
	e := exec.Command("/bin/sh", "-c", command)
	x, err := e.CombinedOutput()
	if err != nil {
		log.Println("Failed to transcribe", err)
	}
	log.Println("Output: ", string(x))

}

func Test2(t *testing.T) {

}
