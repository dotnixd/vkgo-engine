package vkapi

import "../utils"

// Method это сборник методов
type Method struct {
	VK *VK
}

// SendMessage отправляет сообщение
func (m *Method) SendMessage(msg string, peerID int64, enableMentions ...bool) []byte {
	args := map[string]string{
		"message":   msg,
		"peer_id":   utils.Int64ToString(peerID),
		"random_id": "0"}
	if len(enableMentions) == 0 {
		args["disable_mentions"] = "1"
	}

	return m.VK.Request("messages.send", args)
}

// GetConversationMembers получает всех участников беседы
func (m *Method) GetConversationMembers(peerID int64) []byte {
	return m.VK.Request("messages.getConversationMembers", map[string]string{
		"peer_id": utils.Int64ToString(peerID)})
}

// UserGet получает информацию об пользователе
func (m *Method) UserGet(userID int64) []byte {
	return m.VK.Request("users.get", map[string]string{
		"user_ids": utils.Int64ToString(userID)})
}
