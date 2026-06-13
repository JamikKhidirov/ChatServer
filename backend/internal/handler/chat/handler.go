package chathandler

import (
	"io"
	"os"
	"path/filepath"

	chatdomain "ChatServerGolang/backend/internal/domain/chat"
	"ChatServerGolang/backend/internal/service"
	"ChatServerGolang/backend/internal/ws"
	"ChatServerGolang/backend/pkg/response"

	"github.com/gin-gonic/gin"
)

type ChatHandler struct {
	chatService service.ChatService
	hub         *ws.Hub
}

func NewChatHandler(chatService service.ChatService, hub *ws.Hub) *ChatHandler {
	return &ChatHandler{chatService: chatService, hub: hub}
}

// CreateChat creates a new chat (private or group)
// @Tags Чаты
// @Summary Создать новый чат
// @Security BearerAuth
// @Accept json
// @Produce json
// @Description Создаёт новый личный или групповой чат. Для личного чата укажите тип private и список участников. Для группового — тип group, название и участников.
// @Param request body chatdomain.CreateChatRequest true "Данные для создания чата: type (тип: private/group/channel, обязательно), participant_ids (ID участников, обязательно), name (название для группы, опционально), description (описание группы, опционально)"
// @Success 201 {object} chatdomain.ChatResponse "Чат успешно создан"
// @Failure 400 {object} response.ErrorResponse "Неверные входные данные"
// @Router /chats [post]
func (h *ChatHandler) CreateChat(c *gin.Context) {
	userID, _ := c.Get("userID")

	var req chatdomain.CreateChatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	chat, err := h.chatService.CreateChat(userID.(string), &req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	c.Set("chatResponse", chat)
	response.JSON(c, 201, chat)
}

// GetChat returns chat details
// @Tags Чаты
// @Summary Получить детали чата
// @Security BearerAuth
// @Produce json
// @Description Возвращает полную информацию о чате, включая участников, настройки и последнее сообщение.
// @Param id path string true "ID чата"
// @Success 200 {object} chatdomain.ChatResponse "Детали чата"
// @Failure 404 {object} response.ErrorResponse "Чат не найден"
// @Router /chats/{id} [get]
func (h *ChatHandler) GetChat(c *gin.Context) {
	userID, _ := c.Get("userID")
	chatID := c.Param("id")

	chat, err := h.chatService.GetChat(chatID, userID.(string))
	if err != nil {
		response.NotFound(c, err.Error())
		return
	}

	response.JSON(c, 200, chat)
}

// SearchChats searches the user's chats by name
// @Tags Чаты
// @Summary Найти чаты по названию
// @Security BearerAuth
// @Produce json
// @Description Ищет чаты пользователя по названию или имени собеседника.
// @Param q query string true "Поисковый запрос (название чата или имя участника)"
// @Success 200 {array} chatdomain.ChatResponse "Найденные чаты"
// @Failure 400 {object} response.ErrorResponse "Отсутствует поисковый запрос"
// @Router /chats/search [get]
func (h *ChatHandler) SearchChats(c *gin.Context) {
	userID, _ := c.Get("userID")
	query := c.Query("q")
	if query == "" {
		response.BadRequest(c, "query parameter q is required")
		return
	}

	chats, err := h.chatService.SearchChats(userID.(string), query)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.JSON(c, 200, chats)
}

// ListChats returns all chats for the authenticated user
// @Tags Чаты
// @Summary Получить список чатов
// @Security BearerAuth
// @Produce json
// @Description Возвращает все чаты аутентифицированного пользователя, включая личные и групповые диалоги, отсортированные по времени последней активности.
// @Success 200 {array} chatdomain.ChatResponse "Список чатов пользователя"
// @Failure 400 {object} response.ErrorResponse "Ошибка получения списка чатов"
// @Router /chats [get]
func (h *ChatHandler) ListChats(c *gin.Context) {
	userID, _ := c.Get("userID")

	chats, err := h.chatService.ListChats(userID.(string))
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.JSON(c, 200, chats)
}

// DeleteChat deletes a chat (owner only)
// @Tags Чаты
// @Summary Удалить чат
// @Security BearerAuth
// @Produce json
// @Description Удаляет чат. Доступно только владельцу группы или канала. Личные чаты удаляются для всех участников.
// @Param id path string true "ID чата"
// @Success 200 {object} response.MessageResponse "Чат успешно удалён"
// @Failure 400 {object} response.ErrorResponse "Ошибка удаления чата"
// @Router /chats/{id} [delete]
func (h *ChatHandler) DeleteChat(c *gin.Context) {
	userID, _ := c.Get("userID")
	chatID := c.Param("id")

	if err := h.chatService.DeleteChat(chatID, userID.(string)); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.JSON(c, 200, gin.H{"message": "chat deleted"})
}

// AddParticipant adds a user to a group chat
// @Tags Чаты
// @Summary Добавить участника в группу
// @Security BearerAuth
// @Accept json
// @Produce json
// @Description Добавляет нового участника в групповой чат. Требуются права администратора.
// @Param id path string true "ID чата"
// @Param request body chatdomain.AddParticipantRequest true "Данные для добавления: user_id (ID пользователя, обязательно)"
// @Success 200 {object} response.MessageResponse "Участник добавлен"
// @Failure 400 {object} response.ErrorResponse "Ошибка добавления участника"
// @Router /chats/{id}/participants [post]
func (h *ChatHandler) AddParticipant(c *gin.Context) {
	userID, _ := c.Get("userID")
	chatID := c.Param("id")

	var req chatdomain.AddParticipantRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := h.chatService.AddParticipant(chatID, req.UserID, userID.(string)); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.JSON(c, 200, gin.H{"message": "participant added"})
}

// RemoveParticipant removes a user from a group chat
// @Tags Чаты
// @Summary Удалить участника из группы
// @Security BearerAuth
// @Produce json
// @Description Удаляет указанного участника из группового чата. Требуются права администратора.
// @Param id path string true "ID чата"
// @Param userId path string true "ID пользователя для удаления"
// @Success 200 {object} response.MessageResponse "Участник удалён"
// @Failure 400 {object} response.ErrorResponse "Ошибка удаления участника"
// @Router /chats/{id}/participants/{userId} [delete]
func (h *ChatHandler) RemoveParticipant(c *gin.Context) {
	userID, _ := c.Get("userID")
	chatID := c.Param("id")
	targetUserID := c.Param("userId")

	if err := h.chatService.RemoveParticipant(chatID, targetUserID, userID.(string)); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.JSON(c, 200, gin.H{"message": "participant removed"})
}

// MarkAsRead marks all messages in a chat as read
// @Tags Чаты
// @Summary Отметить чат как прочитанный
// @Security BearerAuth
// @Produce json
// @Description Помечает все непрочитанные сообщения в чате как прочитанные для текущего пользователя.
// @Param id path string true "ID чата"
// @Success 200 {object} response.MessageResponse "Чат отмечен как прочитанный"
// @Failure 400 {object} response.ErrorResponse "Ошибка отметки"
// @Router /chats/{id}/read [post]
func (h *ChatHandler) MarkAsRead(c *gin.Context) {
	userID, _ := c.Get("userID")
	chatID := c.Param("id")

	if err := h.chatService.MarkAsRead(chatID, userID.(string)); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.JSON(c, 200, gin.H{"message": "marked as read"})
}

// SetRole changes a participant's role (admin/member)
// @Tags Чаты
// @Summary Изменить роль участника
// @Security BearerAuth
// @Accept json
// @Produce json
// @Description Изменяет роль участника в групповом чате: admin (администратор) или member (участник). Требуются права владельца.
// @Param id path string true "ID чата"
// @Param userId path string true "ID целевого пользователя"
// @Param request body object{role=string} true "Новая роль: admin (администратор) или member (участник)"
// @Success 200 {object} response.MessageResponse "Роль участника обновлена"
// @Failure 400 {object} response.ErrorResponse "Ошибка изменения роли"
// @Router /chats/{id}/participants/{userId}/role [put]
func (h *ChatHandler) SetRole(c *gin.Context) {
	userID, _ := c.Get("userID")
	chatID := c.Param("id")
	targetUserID := c.Param("userId")

	var req struct {
		Role string `json:"role" binding:"required,oneof=admin member"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := h.chatService.SetRole(chatID, targetUserID, userID.(string), req.Role); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.JSON(c, 200, gin.H{"message": "role updated to " + req.Role})
}

// LeaveGroup removes the authenticated user from a group chat
// @Tags Чаты
// @Summary Покинуть групповой чат
// @Security BearerAuth
// @Produce json
// @Description Позволяет аутентифицированному пользователю выйти из группового чата.
// @Param id path string true "ID чата"
// @Success 200 {object} response.MessageResponse "Вы покинули группу"
// @Failure 400 {object} response.ErrorResponse "Ошибка выхода из группы"
// @Router /chats/{id}/leave [post]
func (h *ChatHandler) LeaveGroup(c *gin.Context) {
	userID, _ := c.Get("userID")
	chatID := c.Param("id")

	if err := h.chatService.LeaveGroup(chatID, userID.(string)); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.JSON(c, 200, gin.H{"message": "left the group"})
}

// HideChat hides a chat from the user's chat list
// @Tags Чаты
// @Summary Скрыть чат из списка
// @Security BearerAuth
// @Produce json
// @Description Скрывает указанный чат из основного списка чатов пользователя. Чат можно восстановить через поиск.
// @Param id path string true "ID чата"
// @Success 200 {object} response.MessageResponse "Чат скрыт"
// @Failure 400 {object} response.ErrorResponse "Ошибка скрытия чата"
// @Router /chats/{id}/hide [post]
func (h *ChatHandler) HideChat(c *gin.Context) {
	userID, _ := c.Get("userID")
	chatID := c.Param("id")

	if err := h.chatService.HideChat(chatID, userID.(string)); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.JSON(c, 200, gin.H{"message": "chat hidden"})
}

// UpdateGroup updates group chat name/description/photo
// @Tags Чаты
// @Summary Обновить групповой чат
// @Security BearerAuth
// @Accept json
// @Produce json
// @Description Обновляет название, описание или аватар группового чата. Требуются права администратора.
// @Param id path string true "ID чата"
// @Param request body chatdomain.UpdateGroupRequest true "Обновляемые поля: name (название, опционально), description (описание, опционально), avatar_url (URL аватара, опционально)"
// @Success 200 {object} response.MessageResponse "Группа обновлена"
// @Failure 400 {object} response.ErrorResponse "Ошибка обновления группы"
// @Router /chats/{id} [put]
func (h *ChatHandler) UpdateGroup(c *gin.Context) {
	userID, _ := c.Get("userID")
	chatID := c.Param("id")

	var req chatdomain.UpdateGroupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := h.chatService.UpdateGroup(chatID, userID.(string), &req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.JSON(c, 200, gin.H{"message": "group updated"})
}

// PinChat pins a chat to the top of the list
// @Tags Чаты
// @Summary Закрепить чат вверху списка
// @Security BearerAuth
// @Produce json
// @Description Закрепляет чат в верхней части списка диалогов для быстрого доступа.
// @Param id path string true "ID чата"
// @Success 200 {object} response.MessageResponse "Чат закреплён"
// @Failure 400 {object} response.ErrorResponse "Ошибка закрепления чата"
// @Router /chats/{id}/pin [post]
func (h *ChatHandler) PinChat(c *gin.Context) {
	userID, _ := c.Get("userID")
	chatID := c.Param("id")

	if err := h.chatService.PinChat(chatID, userID.(string)); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.JSON(c, 200, gin.H{"message": "chat pinned"})
}

// UnpinChat unpins a chat from the top of the list
// @Tags Чаты
// @Summary Открепить чат
// @Security BearerAuth
// @Produce json
// @Description Открепляет ранее закреплённый чат, возвращая его в обычный порядок списка.
// @Param id path string true "ID чата"
// @Success 200 {object} response.MessageResponse "Чат откреплён"
// @Failure 400 {object} response.ErrorResponse "Ошибка открепления чата"
// @Router /chats/{id}/pin [delete]
func (h *ChatHandler) UnpinChat(c *gin.Context) {
	userID, _ := c.Get("userID")
	chatID := c.Param("id")

	if err := h.chatService.UnpinChat(chatID, userID.(string)); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.JSON(c, 200, gin.H{"message": "chat unpinned"})
}

// ArchiveChat archives a chat
// @Tags Чаты
// @Summary Архивировать чат
// @Security BearerAuth
// @Produce json
// @Description Перемещает чат в архив. Архивированные чаты не отображаются в основном списке, но сохраняют все сообщения.
// @Param id path string true "ID чата"
// @Success 200 {object} response.MessageResponse "Чат архивирован"
// @Failure 400 {object} response.ErrorResponse "Ошибка архивации чата"
// @Router /chats/{id}/archive [post]
func (h *ChatHandler) ArchiveChat(c *gin.Context) {
	userID, _ := c.Get("userID")
	chatID := c.Param("id")

	if err := h.chatService.ArchiveChat(chatID, userID.(string)); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.JSON(c, 200, gin.H{"message": "chat archived"})
}

// UnarchiveChat unarchives a chat
// @Tags Чаты
// @Summary Разархивировать чат
// @Security BearerAuth
// @Produce json
// @Description Возвращает архивированный чат обратно в основной список диалогов.
// @Param id path string true "ID чата"
// @Success 200 {object} response.MessageResponse "Чат разархивирован"
// @Failure 400 {object} response.ErrorResponse "Ошибка разархивации чата"
// @Router /chats/{id}/unarchive [post]
func (h *ChatHandler) UnarchiveChat(c *gin.Context) {
	userID, _ := c.Get("userID")
	chatID := c.Param("id")

	if err := h.chatService.UnarchiveChat(chatID, userID.(string)); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.JSON(c, 200, gin.H{"message": "chat unarchived"})
}

// ListArchivedChats returns the user's archived chats
// @Tags Чаты
// @Summary Получить архивированные чаты
// @Security BearerAuth
// @Produce json
// @Description Возвращает список архивированных чатов пользователя.
// @Success 200 {array} chatdomain.ChatResponse "Список архивированных чатов"
// @Failure 400 {object} response.ErrorResponse "Ошибка получения архивированных чатов"
// @Router /chats/archived [get]
func (h *ChatHandler) ListArchivedChats(c *gin.Context) {
	userID, _ := c.Get("userID")

	chats, err := h.chatService.ListArchivedChats(userID.(string))
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.JSON(c, 200, chats)
}

// TransferOwnership transfers group ownership to another participant
// @Tags Чаты
// @Summary Передать права владельца группы
// @Security BearerAuth
// @Accept json
// @Produce json
// @Description Передаёт права владельца группового чата другому участнику. Доступно только текущему владельцу.
// @Param id path string true "ID чата"
// @Param request body object{userId=string} true "ID нового владельца: userId (строка, ID пользователя, обязательно)"
// @Success 200 {object} response.MessageResponse "Права владельца переданы"
// @Failure 400 {object} response.ErrorResponse "Ошибка передачи прав"
// @Router /chats/{id}/transfer-ownership [post]
func (h *ChatHandler) TransferOwnership(c *gin.Context) {
	userID, _ := c.Get("userID")
	chatID := c.Param("id")

	var req struct {
		UserID string `json:"userId" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := h.chatService.TransferOwnership(chatID, userID.(string), req.UserID); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.JSON(c, 200, gin.H{"message": "ownership transferred"})
}

// SetSlowMode sets slow mode interval for a group chat
// @Tags Чаты
// @Summary Установить медленный режим
// @Security BearerAuth
// @Accept json
// @Produce json
// @Description Устанавливает интервал медленного режима в групповом чате. Участники смогут отправлять сообщения не чаще указанного интервала. 0 — отключить.
// @Param id path string true "ID чата"
// @Param request body object{seconds=integer} true "Интервал медленного режима в секундах (0-3600, 0 = отключено): seconds (целое число, обязательно)"
// @Success 200 {object} response.MessageResponse "Медленный режим обновлён"
// @Failure 400 {object} response.ErrorResponse "Неверное значение интервала"
// @Router /chats/{id}/slow-mode [put]
func (h *ChatHandler) SetSlowMode(c *gin.Context) {
	userID, _ := c.Get("userID")
	chatID := c.Param("id")

	var req struct {
		Seconds int `json:"seconds" binding:"min=0,max=3600"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := h.chatService.SetSlowMode(chatID, userID.(string), req.Seconds); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.JSON(c, 200, gin.H{"message": "slow mode updated"})
}

// PromoteToAdmin promotes a participant to admin role
// @Tags Чаты
// @Summary Повысить до администратора
// @Security BearerAuth
// @Accept json
// @Produce json
// @Description Повышает участника группового чата до роли администратора. Доступно только владельцу группы.
// @Param id path string true "ID чата"
// @Param request body object{userId=string} true "ID целевого пользователя: userId (строка, обязательно)"
// @Success 200 {object} response.MessageResponse "Пользователь повышен до администратора"
// @Failure 400 {object} response.ErrorResponse "Ошибка повышения"
// @Router /chats/{id}/promote [post]
func (h *ChatHandler) PromoteToAdmin(c *gin.Context) {
	userID, _ := c.Get("userID")
	chatID := c.Param("id")

	var req struct {
		UserID string `json:"userId" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := h.chatService.SetRole(chatID, req.UserID, userID.(string), "admin"); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.JSON(c, 200, gin.H{"message": "user promoted to admin"})
}

// DemoteFromAdmin demotes a participant from admin to member
// @Tags Чаты
// @Summary Понизить с администратора до участника
// @Security BearerAuth
// @Accept json
// @Produce json
// @Description Понижает администратора группового чата до обычного участника. Доступно только владельцу группы.
// @Param id path string true "ID чата"
// @Param request body object{userId=string} true "ID целевого пользователя: userId (строка, обязательно)"
// @Success 200 {object} response.MessageResponse "Пользователь понижен до участника"
// @Failure 400 {object} response.ErrorResponse "Ошибка понижения"
// @Router /chats/{id}/demote [post]
func (h *ChatHandler) DemoteFromAdmin(c *gin.Context) {
	userID, _ := c.Get("userID")
	chatID := c.Param("id")

	var req struct {
		UserID string `json:"userId" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := h.chatService.SetRole(chatID, req.UserID, userID.(string), "member"); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.JSON(c, 200, gin.H{"message": "user demoted to member"})
}

// UploadChatPhoto uploads a photo for a group chat
// @Tags Чаты
// @Summary Загрузить фото чата
// @Security BearerAuth
// @Accept multipart/form-data
// @Produce json
// @Description Загружает новое изображение для аватара группового чата. Поддерживаются форматы JPG, PNG и WEBP.
// @Param id path string true "ID чата"
// @Param photo formData file true "Файл изображения для аватара чата"
// @Success 200 {object} response.MessageResponse "Фото чата обновлено"
// @Failure 400 {object} response.ErrorResponse "Ошибка загрузки фото"
// @Router /chats/{id}/photo [post]
func (h *ChatHandler) UploadChatPhoto(c *gin.Context) {
	userID, _ := c.Get("userID")
	chatID := c.Param("id")

	file, header, err := c.Request.FormFile("photo")
	if err != nil {
		response.BadRequest(c, "photo file required")
		return
	}
	defer file.Close()

	ext := filepath.Ext(header.Filename)
	if ext == "" {
		ext = ".jpg"
	}

	uploadDir := "uploads/chat_photos"
	os.MkdirAll(uploadDir, 0755)

	fileName := chatID + ext
	filePath := filepath.Join(uploadDir, fileName)

	out, err := os.Create(filePath)
	if err != nil {
		response.InternalError(c, "failed to save file")
		return
	}
	defer out.Close()

	if _, err := io.Copy(out, file); err != nil {
		response.InternalError(c, "failed to save file")
		return
	}

	photoURL := "/uploads/chat_photos/" + fileName
	if err := h.chatService.UpdateGroup(chatID, userID.(string), &chatdomain.UpdateGroupRequest{AvatarURL: photoURL}); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.JSON(c, 200, gin.H{"photoUrl": photoURL})
}

// GetOnlineMembers returns online members of a chat
// @Tags Чаты
// @Summary Получить онлайн-участников чата
// @Security BearerAuth
// @Produce json
// @Description Возвращает список ID участников чата, которые в данный момент находятся онлайн.
// @Param id path string true "ID чата"
// @Success 200 {object} object{userIds=[]string} "Список ID онлайн-участников: userIds (массив строк)"
// @Failure 400 {object} response.ErrorResponse "Ошибка получения списка"
// @Router /chats/{id}/online [get]
func (h *ChatHandler) GetOnlineMembers(c *gin.Context) {
	userID, _ := c.Get("userID")
	chatID := c.Param("id")

	chat, err := h.chatService.GetChat(chatID, userID.(string))
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	var onlineIDs []string
	for _, p := range chat.Participants {
		if h.hub.IsOnline(p.ID) {
			onlineIDs = append(onlineIDs, p.ID)
		}
	}

	response.JSON(c, 200, gin.H{"userIds": onlineIDs})
}

// SetChatPermissions sets group permissions (who can send messages, add members, etc.)
// @Tags Чаты
// @Summary Установить права группы
// @Security BearerAuth
// @Accept json
// @Produce json
// @Description Настраивает права участников группы: кто может отправлять сообщения и добавлять участников. Доступно администраторам.
// @Param id path string true "ID чата"
// @Param request body object{whoCanSend=string,whoCanAdd=string} true "Права доступа: who_can_send (кто может писать: everyone/admins, обязательно), who_can_add (кто может добавлять: everyone/admins, обязательно)"
// @Success 200 {object} response.MessageResponse "Права обновлены"
// @Failure 400 {object} response.ErrorResponse "Ошибка установки прав"
// @Router /chats/{id}/permissions [put]
func (h *ChatHandler) SetChatPermissions(c *gin.Context) {
	chatID := c.Param("id")

	var req struct {
		WhoCanSend string `json:"whoCanSend" binding:"required,oneof=everyone admins"`
		WhoCanAdd  string `json:"whoCanAdd" binding:"required,oneof=everyone admins"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.JSON(c, 200, gin.H{"chatId": chatID, "whoCanSend": req.WhoCanSend, "whoCanAdd": req.WhoCanAdd})
}

// SetChatWallpaper sets the wallpaper for a chat
// @Tags Чаты
// @Summary Установить обои чата
// @Security BearerAuth
// @Accept multipart/form-data
// @Produce json
// @Description Загружает и устанавливает фоновое изображение (обои) для указанного чата.
// @Param id path string true "ID чата"
// @Param wallpaper formData file true "Изображение для обоев чата"
// @Success 200 {object} response.MessageResponse "Обои чата установлены"
// @Failure 400 {object} response.ErrorResponse "Ошибка загрузки обоев"
// @Router /chats/{id}/wallpaper [post]
func (h *ChatHandler) SetChatWallpaper(c *gin.Context) {
	chatID := c.Param("id")

	file, header, err := c.Request.FormFile("wallpaper")
	if err != nil {
		response.BadRequest(c, "wallpaper file required")
		return
	}
	defer file.Close()

	ext := filepath.Ext(header.Filename)
	if ext == "" {
		ext = ".jpg"
	}

	uploadDir := "uploads/wallpapers"
	os.MkdirAll(uploadDir, 0755)

	fileName := chatID + ext
	filePath := filepath.Join(uploadDir, fileName)

	out, err := os.Create(filePath)
	if err != nil {
		response.InternalError(c, "failed to save wallpaper")
		return
	}
	defer out.Close()

	io.Copy(out, file)

	wallpaperURL := "/uploads/wallpapers/" + fileName
	response.JSON(c, 200, gin.H{"wallpaperUrl": wallpaperURL})
}

// StartPrivateChat finds or creates a private chat with another user
// @Tags Чаты
// @Summary Начать личный чат с пользователем
// @Security BearerAuth
// @Produce json
// @Description Находит существующий или создаёт новый личный чат с указанным пользователем. Нельзя создать чат с самим собой.
// @Param userId path string true "ID пользователя для начала диалога"
// @Success 200 {object} chatdomain.ChatResponse "Существующий личный чат"
// @Success 201 {object} chatdomain.ChatResponse "Новый личный чат создан"
// @Failure 400 {object} response.ErrorResponse "Неверный ID пользователя"
// @Failure 404 {object} response.ErrorResponse "Пользователь не найден"
// @Router /chats/start/{userId} [post]
func (h *ChatHandler) StartPrivateChat(c *gin.Context) {
	userID, _ := c.Get("userID")
	targetUserID := c.Param("userId")

	if targetUserID == userID.(string) {
		response.BadRequest(c, "cannot start chat with yourself")
		return
	}

	req := &chatdomain.CreateChatRequest{
		Type:           chatdomain.ChatPrivate,
		ParticipantIDs: []string{targetUserID},
	}

	chat, err := h.chatService.CreateChat(userID.(string), req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.JSON(c, 200, chat)
}
