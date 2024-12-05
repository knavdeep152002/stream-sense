package utils

import "os"

var (
	OPENAI_API_KEY      = os.Getenv("OPENAI_API_KEY")
	ROOT_TRANSCRIPT_DIR = os.Getenv("ROOT_TRANSCRIPT_DIR")
)
