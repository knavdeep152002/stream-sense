package utils

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"
)

const (
	WavsDir    = "wavs/"
	SegmentDir = "segments/"
)

func createIfNotExists(dir string) (err error) {
	if _, err = os.Stat(dir); os.IsNotExist(err) {
		err = os.Mkdir(dir, 0755)
		if err != nil {
			return
		}
	}
	return
}

func ConvertMP4ToWav(filePath string, fileName string) (err error) {
	// Check if the wavs directory exists
	err = createIfNotExists(WavsDir)
	if err != nil {
		return
	}

	// Convert the mp4 file to wav
	command := "ffmpeg" + " -i " + filePath + " -y " + WavsDir + fileName + ".wav"
	log.Println("Running command: ", command)
	e := exec.Command("/bin/sh", "-c", command)
	x, err := e.CombinedOutput()
	log.Println("Output: ", string(x))
	if err != nil {
		return
	}
	return
}

func GenerateSegments(filePath string, fileName string, videoID string) {
	err := createIfNotExists(SegmentDir)
	if err != nil {
		log.Println("Failed to create segment directory", err)
		return
	}
	currentSegmentsDir := path.Join(SegmentDir, videoID)
	err = createIfNotExists(currentSegmentsDir)
	if err != nil {
		log.Println("Failed to create current segment directory", err)
		return
	}
	command := fmt.Sprintf("ffmpeg -i %s -c:v libx264 -preset veryfast -g 48 -keyint_min 48 -sc_threshold 0 -b:v 2500k -maxrate 2500k -bufsize 5000k -c:a aac -b:a 128k -hls_time 10 -hls_playlist_type vod -hls_segment_filename \"%s\" %s", filePath, path.Join(currentSegmentsDir, "output%03d.ts"), path.Join(currentSegmentsDir, "playlist.m3u8"))
	log.Println("Running command: ", command)
	e := exec.Command("/bin/sh", "-c", command)
	x, err := e.CombinedOutput()
	log.Println("Output: ", string(x))
	if err != nil {
		log.Println("Failed to generate segments", err)
	}
}
