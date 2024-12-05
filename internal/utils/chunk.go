package utils

import "io"

// Chunk is a chunk of a file.
// It contains information to be able to put the full file back together
// when all file chunks have been uploaded.
type Chunk struct {
	UploadID      string // unique id for the current upload.
	ChunkNumber   int32
	TotalChunks   int32
	TotalFileSize int64 // in bytes
	Filename      string
	Data          io.Reader
	UploadDir     string
}

type ChunkReceiver struct {
	Filename  string
	UploadDir string
	VideoID  string
}

type TranscribeArgs struct {
	AudioPath string
	VideoId   string
}
