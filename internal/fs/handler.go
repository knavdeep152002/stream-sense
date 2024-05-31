package fs

import "github.com/gin-gonic/gin"

// @Summary		Upload file
// @Description	Upload file
// @Tags			fs
// @ID				file.upload
// @Accept			multipart/form-data
// @Produce		json
// @Param			upload_id		formData	string	true	"file upload requirements"
// @Param			chunk_number	formData	int		true	"file upload requirements"
// @Param			total_chunks	formData	int		true	"file upload requirements"
// @Param			total_file_size	formData	int		true	"file upload requirements"
// @Param			file_name		formData	string	true	"file upload requirements"
// @Param			file			formData	file	true	"file upload requirements"
// @Success		201				{string}	string	"file saved"
// @Success		200				{string}	string	"ok"
// @Failure		400				{object}	object	"We need ID!!"
// @Failure		404				{object}	object	"Can not find ID"
// @Router			/upload [post]
func UploadChunk(c *gin.Context) {
	err := processChunk(c)
	if err != nil {
		c.JSON(500, gin.H{"message": "error in processing file : " + err.Error()})
		return
	}
	c.JSON(201, gin.H{"message": "file saved"})
}

// @Summary		Complete file upload
// @Description	Complete file upload
// @Tags			fs
// @ID				file.complete
// @Accept			json
// @Produce		json
// @Param			upload_id	query		string	true	"file upload requirements"
// @Param			file_name	query		string	true	"file upload requirements"
// @Success		201			{string}	string	"file saved"
// @Success		200			{string}	string	"ok"
// @Failure		400			{object}	object	"We need ID!!"
// @Failure		404			{object}	object	"Can not find ID"
// @Router			/complete [post]
func CompleteUpload(c *gin.Context) {
	uploadID := c.Query("upload_id")
	filename := c.Query("file_name")
	err := completeChunk(c, uploadID, filename)
	if err != nil {
		c.JSON(500, gin.H{"message": "error in completing file upload : " + err.Error()})
		return
	}
	c.JSON(201, gin.H{"message": "file saved"})
}
