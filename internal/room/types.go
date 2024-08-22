package room

type roomResult struct {
	RoomId   uint   `json:"room_id"`
	VideoID  string `json:"video_id"`
	RoomLink string `json:"room_link"`
}

type RoomControls struct {
	Play            bool
	Pause           bool
	UpdateTimestamp float64
}

type RoomController struct {
	MsgType string
	UserID  uint
	RoomID  uint
	RoomControls
}

type RoomStreamResponse struct {
	RoomID    uint
	AdminId   uint
	MsgType   string
	TimeStamp float64
	Message   string
	UserID    uint
	Status    RoomState
}

type RoomState int

const (
	ROOMSTATE RoomState = iota
	PLAY
	PAUSE
)
