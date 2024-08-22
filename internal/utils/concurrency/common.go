package concurrency

import "log"

func channelRelay[M any](inChan chan M, outChan chan M) {
	for {
		if out, ok := <-inChan; ok {
			outChan <- out
			log.Println("relaying", out)
		} else {
			break
		}
	}
}
