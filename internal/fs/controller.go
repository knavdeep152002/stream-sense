package fs

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"sort"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/knavdeep152002/stream-sense/internal/models"
	"github.com/knavdeep152002/stream-sense/internal/redis"
	"github.com/knavdeep152002/stream-sense/internal/utils"
)

const (
	maxChunkSize = int64(1 << 20) // 5MB

	uploadDir       = "uploads/chunks"
	uploadsFinalDir = "uploads"
)

// processChunk will parse the chunk data from the request and store in a file on disk.
func processChunk(ctx *gin.Context) error {
	chunk, err := ParseChunk(ctx.Request)
	if err != nil {
		return fmt.Errorf("failed to parse chunk %w", err)
	}

	// Let's create the dir to store the file chunks.
	dir := fmt.Sprintf("%s/%s", uploadDir, chunk.UploadID)
	if err := os.MkdirAll(dir, 02750); err != nil {
		return err
	}

	if err := StoreChunk(chunk); err != nil {
		return err
	}

	return nil
}

// completeChunk rebulds the file chunks into the original full file.
// It then stores the file on disk.
func (fs *FSHandler) completeChunk(ctx *gin.Context, uploadID, filename string, vidId uuid.UUID) error {
	uploadedChunkDir := fmt.Sprintf("%s/%s", uploadDir, uploadID)

	newFile, err := os.Create(fmt.Sprintf("%s/%s", uploadsFinalDir, fmt.Sprintf("%s.mp4", vidId.String())))
	if err != nil {
		return fmt.Errorf("failed creating file %w", err)
	}
	defer newFile.Close()

	err = RebuildFile(uploadedChunkDir, newFile)
	if err != nil {
		return fmt.Errorf("failed to rebuild file %w", err)
	}
	err = fs.saveUploadToDb(vidId, ctx.GetUint("userID"), filename)
	if err != nil {
		return fmt.Errorf("failed to save upload to db %w", err)
	}
	// send for audio processing
	contextData := &utils.ChunkReceiver{
		Filename:  fmt.Sprintf("%s.mp4", vidId.String()),
		UploadDir: uploadsFinalDir,
		VideoID:   vidId.String(),
	}
	data, err := json.Marshal(contextData)
	if err != nil {
		return fmt.Errorf("failed to marshal data %w", err)
	}

	redis.RedisClient.Publish(ctx, "ss-audio-preprocess", data)

	return nil
}

func (fs *FSHandler) saveUploadToDb(videoId uuid.UUID, userId uint, fileName string) (err error) {
	// save the video to the db
	video := &models.Uploads{
		VideoId:  videoId.String(),
		FileName: fileName,
		UserID:   userId,
	}
	if err := fs.DB.Create(video).Error; err != nil {
		return err
	}
	return
}

// ParseChunk parse the request body and creates our chunk struct. It expects the data to be sent in a
// specific order and handles validating the order.
func ParseChunk(r *http.Request) (*utils.Chunk, error) {
	var chunk utils.Chunk

	buf := new(bytes.Buffer)

	reader, err := r.MultipartReader()
	if err != nil {
		return nil, err
	}

	// start readings parts
	// 1. upload id
	// 2. chunk number
	// 3. total chunks
	// 4. total file size
	// 5. file name
	// 6. chunk data

	// 1
	if err := getPart("upload_id", reader, buf); err != nil {
		return nil, err
	}

	chunk.UploadID = buf.String()
	buf.Reset()

	// dir to where we store our chunk
	chunk.UploadDir = fmt.Sprintf("%s/%s", uploadDir, chunk.UploadID)

	// 2
	if err := getPart("chunk_number", reader, buf); err != nil {
		return nil, err
	}

	parsedChunkNumber, err := strconv.ParseInt(buf.String(), 10, 32)
	if err != nil {
		return nil, err
	}

	chunk.ChunkNumber = int32(parsedChunkNumber)
	buf.Reset()

	// 3
	if err := getPart("total_chunks", reader, buf); err != nil {
		return nil, err
	}

	parsedTotalChunksNumber, err := strconv.ParseInt(buf.String(), 10, 32)
	if err != nil {
		return nil, err
	}

	chunk.TotalChunks = int32(parsedTotalChunksNumber)
	buf.Reset()

	// 4
	if err := getPart("total_file_size", reader, buf); err != nil {
		return nil, err
	}

	parsedTotalFileSizeNumber, err := strconv.ParseInt(buf.String(), 10, 64)
	if err != nil {
		return nil, err
	}

	chunk.TotalFileSize = parsedTotalFileSizeNumber
	buf.Reset()

	// 5
	if err := getPart("file_name", reader, buf); err != nil {
		return nil, err
	}

	chunk.Filename = buf.String()
	buf.Reset()

	// 6
	part, err := reader.NextPart()
	if err != nil {
		return nil, fmt.Errorf("failed reading chunk part %w", err)
	}

	chunk.Data = part

	return &chunk, nil
}

// StoreChunk stores the chunk on disk for it to later be processed when all other file chunks have been uploaded.
func StoreChunk(chunk *utils.Chunk) (err error) {
	// create the dir to store the file chunks
	chunkFile, err := os.Create(fmt.Sprintf("%s/%d", chunk.UploadDir, chunk.ChunkNumber))
	if err != nil {
		return err
	}

	if _, err := io.CopyN(chunkFile, chunk.Data, maxChunkSize); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// ByChunk is a helper type to sort the files by name. Since the name of the file is it's chunk number
// it makes rebuilding the file a trivial task.
type ByChunk []os.DirEntry

func (a ByChunk) Len() int      { return len(a) }
func (a ByChunk) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a ByChunk) Less(i, j int) bool {
	ai, _ := strconv.Atoi(a[i].Name())
	aj, _ := strconv.Atoi(a[j].Name())
	return ai < aj
}

// RebuildFile grabs all the files from the directory passed on concantinates them to build the original file.
// It stores the file contents in a temp file and returns it.
func RebuildFile(dir string, fullFile *os.File) error {
	fileInfos, err := os.ReadDir(dir)
	if err != nil {
		return err
	}

	sort.Sort(ByChunk(fileInfos))
	for _, fs := range fileInfos {
		if err := appendChunk(dir, fs.Name(), fullFile); err != nil {
			return err
		}
	}

	if err := os.RemoveAll(dir); err != nil {
		return err
	}

	return nil
}

func appendChunk(uploadDir, fileName string, fullFile *os.File) error {
	// src, err := os.Open(uploadDir + "/" + fileName)
	// if err != nil {
	// 	return err
	// }
	fileContent, err := os.ReadFile(uploadDir + "/" + fileName)
	if err != nil {
		return err
	}
	// defer src.Close()
	if _, err := fullFile.Write(fileContent); err != nil {
		return err
	}
	// if _, err := io.Copy(fullFile, src); err != nil {
	// 	return err
	// }
	return nil
}

func getPart(expectedPart string, reader *multipart.Reader, buf *bytes.Buffer) error {
	part, err := reader.NextPart()
	if err != nil {
		return fmt.Errorf("failed reading %s part %w", expectedPart, err)
	}

	if part.FormName() != expectedPart {
		return fmt.Errorf("invalid form name for part. Expected %s got %s", expectedPart, part.FormName())
	}

	if _, err := io.Copy(buf, part); err != nil {
		return fmt.Errorf("failed copying %s part %w", expectedPart, err)
	}

	return nil
}
