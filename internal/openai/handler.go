package openai

import (
	"log"

	"github.com/gin-gonic/gin"
)

// @Summary		Answer a question based on the video context
// @Description	Answer a question based on the video context
// @Tags			OpenAI
// @Accept			json
// @Produce		json
// @Param			videoId		path		string	true	"Video ID"
// @Param			question	body		query	true	"Question"
// @Success		200			{object}	string
// @Failure		400			{object}	string
// @Failure		404			{object}	string
// @Failure		500			{object}	string
// @Security		Bearer
// @Router			/qa/{videoId} [POST]
func VideoIntelligence(c *gin.Context) {
	vidID := c.Param("videoId")
	question := &query{}
	if err := c.BindJSON(question); err != nil {
		c.JSON(400, gin.H{"message": "error in binding json : " + err.Error()})
		return
	}
	log.Println("Question: ", question.Question)
	log.Println("Video ID: ", vidID)
	ans, err, status := qa(vidID, question.Question)
	if err != nil {
		c.JSON(status, gin.H{"message": "error in answering question : " + err.Error()})
		return
	}
	c.JSON(200, gin.H{"message": ans})
}
