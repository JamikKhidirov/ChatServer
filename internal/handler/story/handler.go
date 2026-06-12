package storyhandler

import (
	"io"
	"os"
	"path/filepath"

	storydomain "ChatServerGolang/internal/domain/story"
	"ChatServerGolang/internal/service"
	"ChatServerGolang/pkg/response"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type StoryHandler struct {
	storyService service.StoryService
}

func NewStoryHandler(storyService service.StoryService) *StoryHandler {
	return &StoryHandler{storyService: storyService}
}

// CreateStory creates a new story (photo/video that disappears after 24h)
// @Tags Stories
// @Security BearerAuth
// @Accept multipart/form-data
// @Produce json
// @Param type formData string true "Story type (photo/video)"
// @Param caption formData string false "Optional caption"
// @Param file formData file true "Story media file"
// @Success 201 {object} storydomain.StoryResponse
// @Failure 400 {object} response.ErrorResponse
// @Router /stories [post]
func (h *StoryHandler) CreateStory(c *gin.Context) {
	userID, _ := c.Get("userID")

	storyType := c.PostForm("type")
	if storyType == "" {
		storyType = "photo"
	}

	file, header, err := c.Request.FormFile("file")
	if err != nil {
		response.BadRequest(c, "file required")
		return
	}
	defer file.Close()

	uploadDir := "uploads/stories"
	os.MkdirAll(uploadDir, 0755)

	ext := filepath.Ext(header.Filename)
	if ext == "" {
		ext = ".jpg"
	}

	fileName := uuid.New().String() + ext
	filePath := filepath.Join(uploadDir, fileName)

	out, err := os.Create(filePath)
	if err != nil {
		response.InternalError(c, "failed to save file")
		return
	}
	defer out.Close()

	io.Copy(out, file)

	req := &storydomain.CreateStoryRequest{
		Type:    storydomain.StoryType(storyType),
		Caption: c.PostForm("caption"),
	}

	story, err := h.storyService.CreateStory(userID.(string), req, fileName, "/uploads/stories/"+fileName)
	if err != nil {
		os.Remove(filePath)
		response.BadRequest(c, err.Error())
		return
	}

	response.JSON(c, 201, story)
}

// GetMyStories returns current user's active stories
// @Tags Stories
// @Security BearerAuth
// @Produce json
// @Success 200 {array} storydomain.StoryResponse
// @Router /stories/my [get]
func (h *StoryHandler) GetMyStories(c *gin.Context) {
	userID, _ := c.Get("userID")

	stories, err := h.storyService.GetMyStories(userID.(string))
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.JSON(c, 200, stories)
}

// GetFollowingStories returns stories from contacts and channels
// @Tags Stories
// @Security BearerAuth
// @Produce json
// @Success 200 {array} storydomain.StoryResponse
// @Router /stories [get]
func (h *StoryHandler) GetFollowingStories(c *gin.Context) {
	userID, _ := c.Get("userID")

	stories, err := h.storyService.GetFollowingStories(userID.(string))
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.JSON(c, 200, stories)
}

// GetStoryByID returns a story by ID and marks it as viewed
// @Tags Stories
// @Security BearerAuth
// @Produce json
// @Param id path string true "Story ID"
// @Success 200 {object} storydomain.StoryResponse
// @Failure 404 {object} response.ErrorResponse
// @Router /stories/{id} [get]
func (h *StoryHandler) GetStoryByID(c *gin.Context) {
	userID, _ := c.Get("userID")
	storyID := c.Param("id")

	story, err := h.storyService.GetStoryByID(storyID, userID.(string))
	if err != nil {
		response.NotFound(c, err.Error())
		return
	}

	response.JSON(c, 200, story)
}

// DeleteStory deletes a story (owner only)
// @Tags Stories
// @Security BearerAuth
// @Produce json
// @Param id path string true "Story ID"
// @Success 200 {object} response.MessageResponse
// @Failure 400 {object} response.ErrorResponse
// @Router /stories/{id} [delete]
func (h *StoryHandler) DeleteStory(c *gin.Context) {
	userID, _ := c.Get("userID")
	storyID := c.Param("id")

	if err := h.storyService.DeleteStory(storyID, userID.(string)); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.JSON(c, 200, gin.H{"message": "story deleted"})
}

// GetStoryViews returns viewers of a story (owner only)
// @Tags Stories
// @Security BearerAuth
// @Produce json
// @Param id path string true "Story ID"
// @Success 200 {array} storydomain.StoryView
// @Failure 400 {object} response.ErrorResponse
// @Router /stories/{id}/views [get]
func (h *StoryHandler) GetStoryViews(c *gin.Context) {
	userID, _ := c.Get("userID")
	storyID := c.Param("id")

	views, err := h.storyService.GetStoryViews(storyID, userID.(string))
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.JSON(c, 200, views)
}
