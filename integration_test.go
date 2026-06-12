package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"testing"
	"time"
)

const baseURL = "http://localhost:8080/api"

type testUser struct {
	email    string
	password string
	username string
	id       string
	token    string
}

func TestIntegration(t *testing.T) {
	if os.Getenv("INTEGRATION") != "1" {
		t.Skip("Set INTEGRATION=1 to run")
	}

	user1 := &testUser{
		email:    fmt.Sprintf("test%d@mail.com", time.Now().UnixNano()),
		password: "testpass123",
		username: fmt.Sprintf("user%d", time.Now().UnixNano()),
	}
	user2 := &testUser{
		email:    fmt.Sprintf("test2%d@mail.com", time.Now().UnixNano()),
		password: "testpass456",
		username: fmt.Sprintf("user2%d", time.Now().UnixNano()),
	}

	t.Run("Health", func(t *testing.T) {
		resp := request(t, "GET", "/health", nil, "")
		assertStatus(t, resp, 200)
	})

	t.Run("Register User1", func(t *testing.T) {
		body := map[string]string{
			"username":    user1.username,
			"email":       user1.email,
			"password":    user1.password,
			"displayName": "Test User 1",
		}
		resp := request(t, "POST", "/api/auth/register", body, "")
		assertStatus(t, resp, 200)
		var result map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&result)
		if data, ok := result["data"].(map[string]interface{}); ok {
			user1.id = data["id"].(string)
			user1.token = data["token"].(string)
		}
		t.Logf("User1: id=%s token=%s", user1.id, user1.token)
	})

	if user1.token == "" {
		t.Fatal("Failed to register user1")
	}

	t.Run("Register User2", func(t *testing.T) {
		body := map[string]string{
			"username":    user2.username,
			"email":       user2.email,
			"password":    user2.password,
			"displayName": "Test User 2",
		}
		resp := request(t, "POST", "/api/auth/register", body, "")
		assertStatus(t, resp, 200)
		var result map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&result)
		if data, ok := result["data"].(map[string]interface{}); ok {
			user2.id = data["id"].(string)
			user2.token = data["token"].(string)
		}
		t.Logf("User2: id=%s", user2.id)
	})

	if user2.token == "" {
		t.Fatal("Failed to register user2")
	}

	var chatID string
	var msgID string

	t.Run("Create Private Chat", func(t *testing.T) {
		body := map[string]interface{}{
			"type":           "private",
			"participantIds": []string{user2.id},
		}
		resp := request(t, "POST", "/api/chats", body, user1.token)
		assertStatus(t, resp, 200)
		var result map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&result)
		if data, ok := result["data"].(map[string]interface{}); ok {
			chatID = data["id"].(string)
		}
		t.Logf("ChatID: %s", chatID)
	})

	if chatID == "" {
		t.Fatal("Failed to create chat")
	}

	t.Run("Get Chats", func(t *testing.T) {
		resp := request(t, "GET", "/api/chats", nil, user1.token)
		assertStatus(t, resp, 200)
	})

	t.Run("Get Chat", func(t *testing.T) {
		resp := request(t, "GET", "/api/chats/"+chatID, nil, user1.token)
		assertStatus(t, resp, 200)
	})

	t.Run("Send Message", func(t *testing.T) {
		body := map[string]string{
			"content": "Hello from test!",
			"type":    "text",
		}
		resp := request(t, "POST", "/api/chats/"+chatID+"/messages", body, user1.token)
		assertStatus(t, resp, 201)
		var result map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&result)
		if data, ok := result["data"].(map[string]interface{}); ok {
			msgID = data["id"].(string)
		}
		t.Logf("MsgID: %s", msgID)
	})

	if msgID == "" {
		t.Fatal("Failed to send message")
	}

	t.Run("Get Messages", func(t *testing.T) {
		resp := request(t, "GET", "/api/chats/"+chatID+"/messages", nil, user1.token)
		assertStatus(t, resp, 200)
	})

	t.Run("Edit Message", func(t *testing.T) {
		body := map[string]string{"content": "Edited content"}
		resp := request(t, "PUT", "/api/messages/"+msgID, body, user1.token)
		assertStatus(t, resp, 200)
	})

	t.Run("Add Reaction", func(t *testing.T) {
		body := map[string]string{"emoji": "👍"}
		resp := request(t, "POST", "/api/messages/"+msgID+"/reactions", body, user1.token)
		assertStatus(t, resp, 200)
	})

	t.Run("Star Message", func(t *testing.T) {
		resp := request(t, "POST", "/api/messages/"+msgID+"/star", nil, user1.token)
		assertStatus(t, resp, 200)
	})

	t.Run("Mark Message Read", func(t *testing.T) {
		resp := request(t, "POST", "/api/messages/"+msgID+"/read", nil, user2.token)
		assertStatus(t, resp, 200)
	})

	t.Run("Get User Profile", func(t *testing.T) {
		resp := request(t, "GET", "/api/users/profile", nil, user1.token)
		assertStatus(t, resp, 200)
	})

	t.Run("Search Users", func(t *testing.T) {
		resp := request(t, "GET", "/api/users/search?q="+user2.username, nil, user1.token)
		assertStatus(t, resp, 200)
	})

	t.Run("Get User By ID", func(t *testing.T) {
		resp := request(t, "GET", "/api/users/"+user2.id, nil, user1.token)
		assertStatus(t, resp, 200)
	})

	t.Run("Last Seen", func(t *testing.T) {
		resp := request(t, "GET", "/api/users/"+user2.id+"/last-seen", nil, user1.token)
		assertStatus(t, resp, 200)
	})

	t.Run("Block User", func(t *testing.T) {
		body := map[string]string{"blockedId": user2.id}
		resp := request(t, "POST", "/api/users/block", body, user1.token)
		assertStatus(t, resp, 200)
	})

	t.Run("Get Blocked Users", func(t *testing.T) {
		resp := request(t, "GET", "/api/users/blocked", nil, user1.token)
		assertStatus(t, resp, 200)
	})

	t.Run("Unblock User", func(t *testing.T) {
		resp := request(t, "DELETE", "/api/users/block/"+user2.id, nil, user1.token)
		assertStatus(t, resp, 200)
	})

	t.Run("Create Group Chat", func(t *testing.T) {
		body := map[string]interface{}{
			"name":           "Test Group",
			"description":    "A test group chat",
			"type":           "group",
			"participantIds": []string{user2.id},
		}
		resp := request(t, "POST", "/api/chats", body, user1.token)
		assertStatus(t, resp, 200)
	})

	t.Run("Pin Chat", func(t *testing.T) {
		resp := request(t, "POST", "/api/chats/"+chatID+"/pin", nil, user1.token)
		assertStatus(t, resp, 200)
	})

	t.Run("Unpin Chat", func(t *testing.T) {
		resp := request(t, "DELETE", "/api/chats/"+chatID+"/pin", nil, user1.token)
		assertStatus(t, resp, 200)
	})

	t.Run("Archive Chat", func(t *testing.T) {
		resp := request(t, "POST", "/api/chats/"+chatID+"/archive", nil, user1.token)
		assertStatus(t, resp, 200)
	})

	t.Run("Unarchive Chat", func(t *testing.T) {
		resp := request(t, "POST", "/api/chats/"+chatID+"/unarchive", nil, user1.token)
		assertStatus(t, resp, 200)
	})

	t.Run("Search Chats", func(t *testing.T) {
		resp := request(t, "GET", "/api/chats/search?q=Test", nil, user1.token)
		assertStatus(t, resp, 200)
	})

	t.Run("Account Settings", func(t *testing.T) {
		resp := request(t, "GET", "/api/account/settings", nil, user1.token)
		assertStatus(t, resp, 200)
	})

	t.Run("Update Account Settings", func(t *testing.T) {
		body := map[string]interface{}{
			"language":    "en",
			"theme":       "dark",
			"soundEnabled": true,
		}
		resp := request(t, "PUT", "/api/account/settings", body, user1.token)
		assertStatus(t, resp, 200)
	})

	t.Run("Change Username", func(t *testing.T) {
		newUsername := fmt.Sprintf("new%d", time.Now().UnixNano())
		body := map[string]string{"username": newUsername}
		resp := request(t, "PUT", "/api/users/username", body, user1.token)
		assertStatus(t, resp, 200)
	})

	t.Run("Starred Messages", func(t *testing.T) {
		resp := request(t, "GET", "/api/messages/starred", nil, user1.token)
		assertStatus(t, resp, 200)
	})

	t.Run("Bulk Mark Read", func(t *testing.T) {
		body := map[string]interface{}{
			"messageIds": []string{msgID},
			"chatId":     chatID,
		}
		resp := request(t, "POST", "/api/messages/read/bulk", body, user2.token)
		assertStatus(t, resp, 200)
	})

	t.Run("Online Members", func(t *testing.T) {
		resp := request(t, "GET", "/api/chats/"+chatID+"/online", nil, user1.token)
		assertStatus(t, resp, 200)
	})

	t.Run("Promote/Demote Admin", func(t *testing.T) {
		if chatID == "" {
			t.Skip("No group chat")
		}
		body := map[string]string{"userId": user2.id}
		resp := request(t, "POST", "/api/chats/"+chatID+"/promote", body, user1.token)
		assertStatus(t, resp, 200)

		resp = request(t, "POST", "/api/chats/"+chatID+"/demote", body, user1.token)
		assertStatus(t, resp, 200)
	})

	t.Run("Transfer Ownership", func(t *testing.T) {
		if chatID == "" {
			t.Skip("No chat")
		}
		body := map[string]string{"userId": user2.id}
		resp := request(t, "POST", "/api/chats/"+chatID+"/transfer-ownership", body, user1.token)
		assertStatus(t, resp, 200)
	})

	t.Run("Message History", func(t *testing.T) {
		resp := request(t, "GET", "/api/messages/"+msgID+"/history", nil, user1.token)
		assertStatus(t, resp, 200)
	})

	t.Run("Report Message", func(t *testing.T) {
		body := map[string]string{"reason": "inappropriate content"}
		resp := request(t, "POST", "/api/messages/"+msgID+"/report", body, user1.token)
		assertStatus(t, resp, 200)
	})

	t.Run("Delete Message", func(t *testing.T) {
		resp := request(t, "DELETE", "/api/messages/"+msgID, nil, user1.token)
		assertStatus(t, resp, 200)
	})

	t.Run("Delete Account", func(t *testing.T) {
		resp := request(t, "DELETE", "/api/users/account", nil, user2.token)
		assertStatus(t, resp, 200)
	})

	t.Run("Notification Settings", func(t *testing.T) {
		body := map[string]bool{"muted": true}
		resp := request(t, "PUT", "/api/chats/"+chatID+"/notifications", body, user1.token)
		assertStatus(t, resp, 200)

		resp = request(t, "GET", "/api/chats/"+chatID+"/notifications", nil, user1.token)
		assertStatus(t, resp, 200)
	})

	t.Run("Slow Mode", func(t *testing.T) {
		body := map[string]int{"seconds": 30}
		resp := request(t, "PUT", "/api/chats/"+chatID+"/slow-mode", body, user1.token)
		assertStatus(t, resp, 200)
	})

	t.Run("Chat Permissions", func(t *testing.T) {
		body := map[string]string{
			"whoCanSend": "everyone",
			"whoCanAdd":  "admins",
		}
		resp := request(t, "PUT", "/api/chats/"+chatID+"/permissions", body, user1.token)
		assertStatus(t, resp, 200)
	})

	t.Run("Create Chat Folder", func(t *testing.T) {
		body := map[string]interface{}{
			"name":   "Work",
			"chatIds": []string{chatID},
		}
		resp := request(t, "POST", "/api/folders", body, user1.token)
		assertStatus(t, resp, 200)
	})

	t.Run("Get Folders", func(t *testing.T) {
		resp := request(t, "GET", "/api/folders", nil, user1.token)
		assertStatus(t, resp, 200)
	})
}

func request(t *testing.T, method, url string, body interface{}, token string) *http.Response {
	t.Helper()
	var reqBody []byte
	if body != nil {
		reqBody, _ = json.Marshal(body)
	}
	req, err := http.NewRequest(method, url, bytes.NewReader(reqBody))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("Request failed: %v (url: %s)", err, url)
	}
	return resp
}

func assertStatus(t *testing.T, resp *http.Response, expected int) {
	t.Helper()
	if resp.StatusCode != expected {
		var result map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&result)
		t.Errorf("Expected status %d, got %d for %s %s. Body: %v", expected, resp.StatusCode, resp.Request.Method, resp.Request.URL, result)
	}
}
