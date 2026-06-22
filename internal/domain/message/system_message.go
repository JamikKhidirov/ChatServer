package messagedomain

type SystemAction string

const (
	SysUserJoined       SystemAction = "user_joined"
	SysUserLeft         SystemAction = "user_left"
	SysUserRemoved      SystemAction = "user_removed"
	SysUserAdded        SystemAction = "user_added"
	SysRoleChanged      SystemAction = "role_changed"
	SysChatCreated      SystemAction = "chat_created"
	SysChatRenamed      SystemAction = "chat_renamed"
	SysChatPhotoChanged SystemAction = "chat_photo_changed"
	SysMessagePinned    SystemAction = "message_pinned"
	SysMessageUnpinned  SystemAction = "message_unpinned"
	SysGroupMigrated    SystemAction = "group_migrated"
)

type SystemMessageContent struct {
	Action     SystemAction `json:"action"`
	UserID     string       `json:"userId"`
	UserName   string       `json:"userName"`
	TargetID   string       `json:"targetId,omitempty"`
	TargetName string       `json:"targetName,omitempty"`
	Extra      string       `json:"extra,omitempty"`
}
