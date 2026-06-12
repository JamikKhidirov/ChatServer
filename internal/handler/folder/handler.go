package folderhandler

import (
	"ChatServerGolang/internal/domain/chat"
	"ChatServerGolang/internal/service"
	"ChatServerGolang/pkg/response"

	"github.com/gin-gonic/gin"
)

type ChatFolderHandler struct {
	folderService service.ChatFolderService
}

func NewChatFolderHandler(folderService service.ChatFolderService) *ChatFolderHandler {
	return &ChatFolderHandler{folderService: folderService}
}

// CreateFolder creates a new chat folder
// @Tags Folders
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body chatdomain.CreateChatFolderRequest true "Folder name and chat IDs"
// @Success 201 {object} chatdomain.ChatFolder
// @Failure 400 {object} response.ErrorResponse
// @Router /folders [post]
func (h *ChatFolderHandler) CreateFolder(c *gin.Context) {
	userID, _ := c.Get("userID")

	var req chatdomain.CreateChatFolderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	folder, err := h.folderService.Create(userID.(string), &req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.JSON(c, 201, folder)
}

// ListFolders returns all chat folders for the user
// @Tags Folders
// @Security BearerAuth
// @Produce json
// @Success 200 {array} chatdomain.ChatFolder
// @Failure 400 {object} response.ErrorResponse
// @Router /folders [get]
func (h *ChatFolderHandler) ListFolders(c *gin.Context) {
	userID, _ := c.Get("userID")

	folders, err := h.folderService.List(userID.(string))
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.JSON(c, 200, folders)
}

// UpdateFolder updates a chat folder's name or chat list
// @Tags Folders
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "Folder ID"
// @Param request body chatdomain.UpdateChatFolderRequest true "Fields to update"
// @Success 200 {object} chatdomain.ChatFolder
// @Failure 400 {object} response.ErrorResponse
// @Router /folders/{id} [put]
func (h *ChatFolderHandler) UpdateFolder(c *gin.Context) {
	userID, _ := c.Get("userID")
	folderID := c.Param("id")

	var req chatdomain.UpdateChatFolderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	folder, err := h.folderService.Update(folderID, userID.(string), &req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.JSON(c, 200, folder)
}

// DeleteFolder deletes a chat folder
// @Tags Folders
// @Security BearerAuth
// @Produce json
// @Param id path string true "Folder ID"
// @Success 200 {object} response.MessageResponse
// @Failure 400 {object} response.ErrorResponse
// @Router /folders/{id} [delete]
func (h *ChatFolderHandler) DeleteFolder(c *gin.Context) {
	userID, _ := c.Get("userID")
	folderID := c.Param("id")

	if err := h.folderService.Delete(folderID, userID.(string)); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.JSON(c, 200, gin.H{"message": "folder deleted"})
}
