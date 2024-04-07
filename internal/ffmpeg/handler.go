package ffmpeg

import (
	"log"
	"os"
	"os/exec"
)

const (
	WavsDir = "wavs/"
)

func ConvertMP4ToWav(filePath string, fileName string) (err error) {
	// Convert the mp4 file to wav

	// Check if the wavs directory exists
	if _, err = os.Stat(WavsDir); os.IsNotExist(err) {
		err = os.Mkdir(WavsDir, 0755)
		if err != nil {
			return
		}
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
