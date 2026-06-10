$ErrorActionPreference = "Stop"
$BASE = "http://localhost:8080"
$script:passed = 0
$script:failed = 0
$script:token = $null
$script:uid = $null
$script:cid = $null
$script:mid = $null

function GET($url, $token) { $h = @{}; if ($token) { $h.Authorization = "Bearer $token" }; return Invoke-RestMethod -Uri "$BASE$url" -Method Get -Headers $h -ContentType "application/json" }
function POST($url, $body, $token) { $h = @{'Content-Type' = 'application/json'}; if ($token) { $h.Authorization = "Bearer $token" }; return Invoke-RestMethod -Uri "$BASE$url" -Method Post -Headers $h -Body ($body | ConvertTo-Json -Compress -Depth 10) }
function PUT($url, $body, $token) { $h = @{'Content-Type' = 'application/json'}; if ($token) { $h.Authorization = "Bearer $token" }; return Invoke-RestMethod -Uri "$BASE$url" -Method Put -Headers $h -Body ($body | ConvertTo-Json -Compress -Depth 10) }
function DEL($url, $token) { $h = @{'Content-Type' = 'application/json'}; if ($token) { $h.Authorization = "Bearer $token" }; return Invoke-RestMethod -Uri "$BASE$url" -Method Delete -Headers $h }
function CHECK($n, $r) { if ($r.success) { Write-Host "  OK: $n" -ForegroundColor Green; $script:passed++ } else { Write-Host "  FAIL: $n" -ForegroundColor Red; $script:failed++ } }

Write-Host "=== FULL TEST SUITE ===" -ForegroundColor Cyan

# 1
CHECK "Health" (GET "/health" $null)
# 2
CHECK "Captcha Generate" (GET "/api/captcha/generate" $null)

# 3
Write-Host "--- AUTH ---" -ForegroundColor Yellow
$reg = POST "/api/auth/register" @{username="user1"; email="u1@x.com"; password="pass123"; display_name="User One"} $null
CHECK "Register" $reg
$script:token = $reg.data.token
$script:uid = $reg.data.user.id
Write-Host "Token: $($script:token.Substring(0, [Math]::Min(20, $script:token.Length)))..." -ForegroundColor DarkGray

if ($script:token) {
    $tk = $script:token
    CHECK "Login" (POST "/api/auth/login" @{email="u1@x.com"; password="pass123"} $null)
    CHECK "Profile" (GET "/api/users/profile" $tk)
    CHECK "Update Profile" (PUT "/api/users/profile" @{display_name="Updated Name"} $tk)
    CHECK "Search Users" (GET "/api/users/search?query=user" $tk)
    CHECK "Get User By ID" (GET "/api/users/$($script:uid)" $tk)
    CHECK "Get By Username" (GET "/api/users/username/user1" $tk)
    CHECK "Settings GET" (GET "/api/account/settings" $tk)
    CHECK "Settings PUT" (PUT "/api/account/settings" @{language="ru"; theme="dark"; notifications=$true; sound_enabled=$true; last_seen_mode="contacts"} $tk)
    CHECK "Refresh Token" (GET "/api/auth/refresh" $tk)
    CHECK "Change Password" (PUT "/api/auth/change-password" @{old_password="pass123"; new_password="newpass456"} $tk)

    Write-Host "--- CHATS ---" -ForegroundColor Yellow
    $chat = POST "/api/chats" @{name="Test Group"; type="group"} $tk
    CHECK "Create Chat" $chat
    $script:cid = $chat.data.id
    if ($script:cid) {
        $cid = $script:cid
        CHECK "List Chats" (GET "/api/chats" $tk)
        CHECK "Get Chat" (GET "/api/chats/$cid" $tk)
        CHECK "Search Chats" (GET "/api/chats/search?query=Test" $tk)
        CHECK "Notification GET" (GET "/api/chats/$cid/notifications" $tk)
        CHECK "Notification PUT" (PUT "/api/chats/$cid/notifications" @{muted=$true} $tk)

        Write-Host "--- MESSAGES ---" -ForegroundColor Yellow
        $msg = POST "/api/chats/$cid/messages" @{content="Hello World!"; type="text"} $tk
        CHECK "Send Message" $msg
        $script:mid = $msg.data.id
        if ($script:mid) {
            $mid = $script:mid
            CHECK "List Messages" (GET "/api/chats/$cid/messages" $tk)
            CHECK "Get Message" (GET "/api/messages/$mid" $tk)
            CHECK "Edit Message" (PUT "/api/messages/$mid" @{content="Edited Content"} $tk)
            CHECK "Add Reaction" (POST "/api/messages/$mid/reactions" @{emoji="like"} $tk)
            CHECK "Pin Message" (PUT "/api/messages/$mid/pin" @{} $tk)
            CHECK "Star Message" (POST "/api/messages/$mid/star" @{} $tk)
            CHECK "Unstar Message" (DEL "/api/messages/$mid/star" $tk)
            CHECK "Starred List" (GET "/api/messages/starred" $tk)
            CHECK "Pinned List" (GET "/api/chats/$cid/pinned" $tk)
            CHECK "Chat Media" (GET "/api/chats/$cid/media" $tk)
            CHECK "Export Chat" (GET "/api/chats/$cid/export" $tk)
            CHECK "Search Msg" (GET "/api/chats/$cid/messages/search?query=Hello" $tk)
            CHECK "Search All" (GET "/api/messages/search?query=Hello" $tk)
            CHECK "Mark Read" (POST "/api/messages/$mid/read" @{} $tk)
            CHECK "Delete For Me" (DEL "/api/messages/$mid/for-me" $tk)
            CHECK "Delete Message" (DEL "/api/messages/$mid" $tk)
        }

        Write-Host "--- CHAT ACTIONS ---" -ForegroundColor Yellow
        CHECK "Pin Chat" (POST "/api/chats/$cid/pin" @{} $tk)
        CHECK "Unpin Chat" (DEL "/api/chats/$cid/pin" $tk)
        CHECK "Archive Chat" (POST "/api/chats/$cid/archive" @{} $tk)
        CHECK "Unarchive Chat" (POST "/api/chats/$cid/unarchive" @{} $tk)
        CHECK "List Archived" (GET "/api/chats/archived" $tk)
        CHECK "Hide Chat" (POST "/api/chats/$cid/hide" @{} $tk)
        CHECK "Mark Chat Read" (POST "/api/chats/$cid/read" @{} $tk)
    }

    Write-Host "--- FEATURES ---" -ForegroundColor Yellow
    CHECK "Drafts Save" (POST "/api/drafts" @{chat_id=$script:cid; content="test draft"} $tk)
    CHECK "Drafts Get" (GET "/api/drafts?chat_id=$($script:cid)" $tk)
    CHECK "Sessions" (GET "/api/sessions" $tk)
    CHECK "Blocked List" (GET "/api/users/blocked" $tk)
    CHECK "Email Verification" (POST "/api/verification/email/send" @{email="test@x.com"} $tk)
    CHECK "Phone Verification" (POST "/api/verification/phone/send" @{phone="+1234567890"} $tk)
    CHECK "Bookmarks List" (GET "/api/bookmarks" $tk)
    CHECK "E2E Register Key" (POST "/api/e2e/keys" @{public_key="test-pub"; private_key_encrypted="test-priv"} $tk)
    CHECK "E2E Get Key" (GET "/api/e2e/keys/$($script:uid)" $tk)
    CHECK "Reports List" (GET "/api/reports" $tk)
    CHECK "Stickers List" (GET "/api/stickers/packs" $tk)
    CHECK "GIFs List" (GET "/api/gifs" $tk)
    CHECK "Scheduled List" (GET "/api/messages/scheduled" $tk)
    CHECK "Bots List" (GET "/api/bots" $tk)
    CHECK "Admin Dashboard" (GET "/api/admin/dashboard" $tk)
    CHECK "Admin Users" (GET "/api/admin/users" $tk)
    CHECK "Admin Messages" (GET "/api/admin/messages" $tk)
    CHECK "Admin Settings" (GET "/api/admin/settings" $tk)
    CHECK "Admin Logs" (GET "/api/admin/logs" $tk)
    CHECK "Admin IP Blocks" (GET "/api/admin/ip-blocks" $tk)
    CHECK "Admin Setting Update" (PUT "/api/admin/settings" @{key="captcha_enabled"; value="true"} $tk)
    CHECK "Link Preview" (GET "/api/preview?url=https://example.com" $null)
}

Write-Host ""
Write-Host "=== RESULTS ===" -ForegroundColor Cyan
Write-Host " PASSED: $($script:passed)" -ForegroundColor Green
Write-Host " FAILED: $($script:failed)" -ForegroundColor Red
if ($script:failed -eq 0) { Write-Host " ALL TESTS PASSED!" -ForegroundColor Green } else { Write-Host " SOME TESTS FAILED!" -ForegroundColor Red }
