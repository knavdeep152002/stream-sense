package fs

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/knavdeep152002/stream-sense/internal/models"
	"gorm.io/gorm"
)

type FSHandler struct {
	DB *gorm.DB
}

func CreateFSHandler(db *gorm.DB) *FSHandler {
	return &FSHandler{
		DB: db,
	}
}

//	@Summary		Upload file
//	@Description	Upload file
//	@Tags			fs
//	@ID				file.upload
//	@Accept			multipart/form-data
//	@Produce		json
//	@Param			upload_id		formData	string	true	"file upload requirements"
//	@Param			chunk_number	formData	int		true	"file upload requirements"
//	@Param			total_chunks	formData	int		true	"file upload requirements"
//	@Param			total_file_size	formData	int		true	"file upload requirements"
//	@Param			file_name		formData	string	true	"file upload requirements"
//	@Param			file			formData	file	true	"file upload requirements"
//	@Success		201				{string}	string	"file saved"
//	@Success		200				{string}	string	"ok"
//	@Failure		400				{object}	object	"We need ID!!"
//	@Failure		404				{object}	object	"Can not find ID"
//	@Security		Bearer
//	@Router			/upload [post]
func (fs *FSHandler) UploadChunk(c *gin.Context) {
	err := processChunk(c)
	if err != nil {
		c.JSON(500, gin.H{"message": "error in processing file : " + err.Error()})
		return
	}
	c.JSON(201, gin.H{"message": "file saved"})
}

//	@Summary		Complete file upload
//	@Description	Complete file upload
//	@Tags			fs
//	@ID				file.complete
//	@Accept			json
//	@Produce		json
//	@Param			upload_id	query		string	true	"file upload requirements"
//	@Param			file_name	query		string	true	"file upload requirements"
//	@Success		201			{string}	string	"file saved"
//	@Success		200			{string}	string	"ok"
//	@Failure		400			{object}	object	"We need ID!!"
//	@Failure		404			{object}	object	"Can not find ID"
//	@Security		Bearer
//	@Router			/complete [post]
func (fs *FSHandler) CompleteUpload(c *gin.Context) {
	uploadID := c.Query("upload_id")
	filename := c.Query("file_name")
	vidId, err := uuid.NewRandom()
	if err != nil {
		c.JSON(500, gin.H{"message": "error in generating video id : " + err.Error()})
		return
	}
	err = fs.completeChunk(c, uploadID, filename, vidId)
	if err != nil {
		c.JSON(500, gin.H{"message": "error in completing file upload : " + err.Error()})
		return
	}
	c.JSON(201, gin.H{"message": "file saved"})
}

//	@Summary		Get user uploads
//	@Description	Get user uploads
//	@Tags			fs
//	@ID				user.uploads
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	[]models.Uploads
//	@Failure		400	{object}	object
//	@Security		Bearer
//	@Router			/uploads [get]
func (fs *FSHandler) GetUserUploads(c *gin.Context) {
	// get all uploads for the user
	userId := c.GetUint("userID")
	var uploads []models.Uploads
	err := fs.DB.Where("user_id=?", userId).Find(&uploads).Error
	if err != nil {
		c.JSON(500, gin.H{"message": "error in getting user uploads : " + err.Error()})
		return
	}
	var uploadResult []uploadsResult
	for _, upload := range uploads {
		uploadResult = append(uploadResult, uploadsResult{
			FileName: upload.FileName,
			VideoId:  upload.VideoId,
		})
	}
	c.JSON(200, gin.H{"data": uploadResult})
}
