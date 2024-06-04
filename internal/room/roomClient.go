package room

import (
	"log"
	"time"

	"github.com/knavdeep152002/stream-sense/internal/utils/concurrency"
)

type RoomClientHandler struct {
	ControlChan chan *RoomController
	OutChan     chan *RoomStreamResponse
	Publisher   *concurrency.Publisher[*RoomStreamResponse]
	Observer    *concurrency.Observer[*RoomController]
	AdminId     uint
	Ready       bool
	state       RoomState
	TimeStamp   float64
	VideoId     string
	errs        chan error
}

func (rch *RoomClientHandler) takeControl() {
	log.Println("Starting stream")
	defer func() {
		log.Println("Canceling stream")
		rch.Ready = false
	}()

	run := true

	for run {
		select {
		case out, ok := <-rch.ControlChan:
			if ok {
				log.Println("action based on control: ", ok, out)
				roomStreamResponse := &RoomStreamResponse{
					RoomID:  out.RoomID,
					AdminId: rch.AdminId,
				}
				switch out.MsgType {
				case "control":
					if out.Play {
						rch.state = PLAY
						roomStreamResponse.MsgType = "play"
					} else if out.Pause {
						rch.state = PAUSE
						roomStreamResponse.MsgType = "pause"
					} else if out.UpdateTimestamp > 0 {
						rch.TimeStamp = out.UpdateTimestamp
						roomStreamResponse.MsgType = "update"
					}
					roomStreamResponse.TimeStamp = rch.TimeStamp
					roomStreamResponse.Status = rch.state
					roomStreamResponse.UserID = out.UserID
					rch.OutChan <- roomStreamResponse
				case "current_status":
					log.Println("Sending current status")
					roomStreamResponse.MsgType = "current_status"
					roomStreamResponse.Status = rch.state
					roomStreamResponse.TimeStamp = rch.TimeStamp
					roomStreamResponse.UserID = out.UserID
					rch.OutChan <- roomStreamResponse
				default:
					log.Println("Invalid control message")
				}

				// client.Send(out)
			} else {
				log.Println("Command channel closed")
				run = false
				return
			}
		case <-time.After(time.Second):
		}
	}
	log.Println("Ending stream")
}

func NewRuntimeClientHandler(videoId string, adminId uint) (rch *RoomClientHandler) {
	rch = &RoomClientHandler{
		OutChan:     make(chan *RoomStreamResponse, 10),
		ControlChan: make(chan *RoomController, 10),
		errs:        make(chan error, 10),
		VideoId:     videoId,
		AdminId:     adminId,
		Ready:       true,
	}
	rch.Publisher = concurrency.NewPublisher(rch.OutChan)
	rch.Observer = concurrency.NewObserver(rch.ControlChan)
	go rch.takeControl()
	return
}
