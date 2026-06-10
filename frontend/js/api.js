var authToken = localStorage.getItem('messenger_token') || '';
var currentEmoji = '';

async function apiCall(method, path, body, isFormData) {
  var needsAuth = !path.startsWith('/api/auth/register') && !path.startsWith('/api/auth/login');
  if (needsAuth && !authToken) { toast('Please login first', 'error'); throw new Error('No auth'); }
  var opts = { method: method };
  var headers = {};
  if (authToken) headers['Authorization'] = 'Bearer ' + authToken;
  if (body && !isFormData) { headers['Content-Type'] = 'application/json'; opts.body = JSON.stringify(body); }
  else if (body && isFormData) { opts.body = body; }
  if (Object.keys(headers).length > 0) opts.headers = headers;
  try {
    var res = await fetch('http://localhost:8080' + path, opts);
    var ct = res.headers.get('content-type') || '';
    var textBody = await res.text();
    var data;
    if (ct.includes('application/json')) { try { data = JSON.parse(textBody); } catch(e) { data = textBody; } }
    else { data = textBody; }
    if (!res.ok) {
      var msg = 'Request failed';
      if (typeof data === 'object' && data !== null) msg = data.error || data.message || data.detail || JSON.stringify(data);
      else if (typeof data === 'string' && data) msg = data;
      toast('HTTP ' + res.status + ': ' + (typeof msg === 'string' && msg.length > 80 ? msg.substring(0,80)+'...' : msg), 'error');
      return { error: true, status: res.status, data: data };
    }
    return { error: false, status: res.status, data: data };
  } catch (e) {
    toast('Network error: ' + e.message, 'error');
    return { error: true, data: e.message };
  }
}

window.apiRegister = async function(btn) {
  setLoading(btn, true);
  var body = { username: document.getElementById('regUsername').value, email: document.getElementById('regEmail').value, password: document.getElementById('regPassword').value, displayName: document.getElementById('regDisplayName').value };
  var r = await apiCall('POST', '/api/auth/register', body);
  setLoading(btn, false);
  if (!r.error && r.data && r.data.token) { authToken = r.data.token; updateTokenDisplay(); updateConnectionStatus(); }
  showResult('resultAuth', r.data || r);
};

window.apiLogin = async function(btn) {
  setLoading(btn, true);
  var body = { email: document.getElementById('loginEmail').value, password: document.getElementById('loginPassword').value };
  var r = await apiCall('POST', '/api/auth/login', body);
  setLoading(btn, false);
  if (!r.error && r.data && r.data.token) { authToken = r.data.token; updateTokenDisplay(); updateConnectionStatus(); toast('Login successful', 'success'); }
  showResult('resultAuth', r.data || r);
};

window.apiRefreshToken = async function(btn) {
  setLoading(btn, true);
  var r = await apiCall('GET', '/api/auth/refresh');
  setLoading(btn, false);
  if (!r.error && r.data && r.data.token) { authToken = r.data.token; updateTokenDisplay(); toast('Token refreshed', 'success'); }
  showResult('resultAuth', r.data || r);
};

window.apiChangePassword = async function(btn) {
  setLoading(btn, true);
  var body = { oldPassword: document.getElementById('oldPassword').value, newPassword: document.getElementById('newPassword').value };
  var r = await apiCall('PUT', '/api/auth/change-password', body);
  setLoading(btn, false);
  showResult('resultAuth', r.data || r);
  if (!r.error) toast('Password changed', 'success');
};

window.apiGetProfile = async function(btn) {
  setLoading(btn, true);
  var r = await apiCall('GET', '/api/users/profile');
  setLoading(btn, false);
  showResult('resultProfile', r.data || r);
};

window.apiUpdateProfile = async function(btn) {
  setLoading(btn, true);
  var body = {};
  var map = { displayName:'updDisplayName', bio:'updBio', phone:'updPhone', gender:'updGender', dateOfBirth:'updDob' };
  for (var k in map) { if (!map.hasOwnProperty(k)) continue; var v = document.getElementById(map[k]).value; if (v) body[k] = v; }
  if (Object.keys(body).length === 0) { toast('Fill at least one field', 'error'); setLoading(btn, false); return; }
  var r = await apiCall('PUT', '/api/users/profile', body);
  setLoading(btn, false);
  showResult('resultProfile', r.data || r);
};

window.apiUploadAvatar = async function(btn) {
  setLoading(btn, true);
  var file = document.getElementById('avatarFile').files[0];
  if (!file) { toast('Select a file', 'error'); setLoading(btn, false); return; }
  var fd = new FormData(); fd.append('avatar', file);
  var r = await apiCall('POST', '/api/users/avatar', fd, true);
  setLoading(btn, false);
  showResult('resultProfile', r.data || r);
};

window.apiUpdateStatus = async function(btn) {
  setLoading(btn, true);
  var body = { text: document.getElementById('statusText').value, type: document.getElementById('statusType').value };
  var r = await apiCall('PUT', '/api/users/status', body);
  setLoading(btn, false);
  showResult('resultProfile', r.data || r);
};

window.apiDeleteAccount = async function(btn) {
  if (!confirm('Delete your account permanently?')) return;
  setLoading(btn, true);
  var r = await apiCall('DELETE', '/api/users/account');
  setLoading(btn, false);
  showResult('resultProfile', r.data || r);
  if (!r.error) { authToken = ''; updateTokenDisplay(); updateConnectionStatus(); }
};

window.apiGetUserById = async function(btn) {
  setLoading(btn, true);
  var id = document.getElementById('userIdLookup').value;
  if (!id) { toast('Enter a user ID', 'error'); setLoading(btn, false); return; }
  var r = await apiCall('GET', '/api/users/' + encodeURIComponent(id));
  setLoading(btn, false);
  showResult('resultProfile', r.data || r);
};

window.apiGetUserByUsername = async function(btn) {
  setLoading(btn, true);
  var u = document.getElementById('usernameLookup').value;
  if (!u) { toast('Enter username', 'error'); setLoading(btn, false); return; }
  var r = await apiCall('GET', '/api/users/username/' + encodeURIComponent(u));
  setLoading(btn, false);
  showResult('resultProfile', r.data || r);
};

window.apiSavePushToken = async function(btn) {
  setLoading(btn, true);
  var body = { token: document.getElementById('pushToken').value, platform: document.getElementById('pushPlatform').value };
  var r = await apiCall('POST', '/api/users/push-token', body);
  setLoading(btn, false);
  showResult('resultProfile', r.data || r);
};

window.apiTestPush = async function(btn) {
  setLoading(btn, true);
  var r = await apiCall('POST', '/api/users/push-test');
  setLoading(btn, false);
  showResult('resultProfile', r.data || r);
};

window.apiGetSettings = async function(btn) {
  setLoading(btn, true);
  var r = await apiCall('GET', '/api/account/settings');
  setLoading(btn, false);
  showResult('resultSettings', r.data || r);
};

window.apiUpdateSettings = async function(btn) {
  setLoading(btn, true);
  var body = {};
  var lang = document.getElementById('setLanguage').value;
  var theme = document.getElementById('setTheme').value;
  if (lang) body.language = lang;
  if (theme) body.theme = theme;
  body.notifications = document.getElementById('setNotifications').checked;
  body.soundEnabled = document.getElementById('setSoundEnabled').checked;
  body.lastSeenMode = document.getElementById('setLastSeenMode').value;
  var r = await apiCall('PUT', '/api/account/settings', body);
  setLoading(btn, false);
  showResult('resultSettings', r.data || r);
};

window.apiSyncContacts = async function(btn) {
  setLoading(btn, true);
  var raw = document.getElementById('contactsSyncData').value;
  if (!raw) { toast('Enter JSON data', 'error'); setLoading(btn, false); return; }
  var data;
  try { data = JSON.parse(raw); } catch(e) { toast('Invalid JSON', 'error'); setLoading(btn, false); return; }
  var r = await apiCall('POST', '/api/contacts/sync', data);
  setLoading(btn, false);
  showResult('resultSyncContacts', r.data || r);
};

window.apiGetContacts = async function(btn) {
  setLoading(btn, true);
  var r = await apiCall('GET', '/api/contacts');
  setLoading(btn, false);
  showResult('resultContacts', r.data || r);
};

window.apiSearchContacts = async function(btn) {
  setLoading(btn, true);
  var q = document.getElementById('contactSearchQuery').value;
  if (!q) { toast('Enter query', 'error'); setLoading(btn, false); return; }
  var r = await apiCall('GET', '/api/contacts/search?q=' + encodeURIComponent(q));
  setLoading(btn, false);
  showResult('resultSearchContacts', r.data || r);
};

window.apiFindRegistered = async function(btn) {
  setLoading(btn, true);
  var r = await apiCall('GET', '/api/contacts/registered');
  setLoading(btn, false);
  showResult('resultRegistered', r.data || r);
};

window.apiSearchChats = async function(btn) {
  setLoading(btn, true);
  var q = document.getElementById('chatSearchQuery').value;
  if (!q) { toast('Enter search query', 'error'); setLoading(btn, false); return; }
  var r = await apiCall('GET', '/api/chats/search?q=' + encodeURIComponent(q));
  setLoading(btn, false);
  showResult('resultSearchChats', r.data || r);
};

window.apiListChats = async function(btn) {
  setLoading(btn, true);
  var r = await apiCall('GET', '/api/chats');
  setLoading(btn, false);
  showResult('resultListChats', r.data || r);
};

window.apiCreatePrivateChat = async function(btn) {
  setLoading(btn, true);
  var userId = document.getElementById('privateChatUserId').value;
  if (!userId) { toast('Enter user ID', 'error'); setLoading(btn, false); return; }
  var body = { type: 'private', participantIds: [userId] };
  var r = await apiCall('POST', '/api/chats', body);
  setLoading(btn, false);
  showResult('resultListChats', r.data || r);
};

window.apiCreateGroupChat = async function(btn) {
  setLoading(btn, true);
  var name = document.getElementById('groupName').value;
  var raw = document.getElementById('groupParticipants').value;
  if (!name) { toast('Enter group name', 'error'); setLoading(btn, false); return; }
  if (!raw) { toast('Enter participant JSON', 'error'); setLoading(btn, false); return; }
  var participants;
  try { participants = JSON.parse(raw); } catch(e) { toast('Invalid JSON', 'error'); setLoading(btn, false); return; }
  if (!Array.isArray(participants) || !participants.length) { toast('Non-empty array required', 'error'); setLoading(btn, false); return; }
  var body = { type: 'group', name: name, participants: participants };
  var r = await apiCall('POST', '/api/chats', body);
  setLoading(btn, false);
  showResult('resultListChats', r.data || r);
};

window.apiGetChat = async function(btn) {
  setLoading(btn, true);
  var id = document.getElementById('chatIdGet').value;
  if (!id) { toast('Enter chat ID', 'error'); setLoading(btn, false); return; }
  var r = await apiCall('GET', '/api/chats/' + encodeURIComponent(id));
  setLoading(btn, false);
  showResult('resultGetChat', r.data || r);
};

window.apiHideChat = async function(btn) {
  setLoading(btn, true);
  var id = document.getElementById('chatIdAction').value;
  if (!id) { toast('Enter chat ID', 'error'); setLoading(btn, false); return; }
  var r = await apiCall('POST', '/api/chats/' + encodeURIComponent(id) + '/hide');
  setLoading(btn, false);
  showResult('resultGetChat', r.data || r);
};

window.apiDeleteChat = async function(btn) {
  var id = document.getElementById('chatIdAction').value;
  if (!id) { toast('Enter chat ID', 'error'); return; }
  if (!confirm('Delete chat?')) return;
  setLoading(btn, true);
  var r = await apiCall('DELETE', '/api/chats/' + encodeURIComponent(id));
  setLoading(btn, false);
  showResult('resultGetChat', r.data || r);
};

window.apiLeaveGroup = async function(btn) {
  setLoading(btn, true);
  var id = document.getElementById('chatIdAction').value;
  if (!id) { toast('Enter chat ID', 'error'); setLoading(btn, false); return; }
  var r = await apiCall('POST', '/api/chats/' + encodeURIComponent(id) + '/leave');
  setLoading(btn, false);
  showResult('resultGetChat', r.data || r);
};

window.apiMarkChatRead = async function(btn) {
  setLoading(btn, true);
  var id = document.getElementById('chatIdAction').value;
  if (!id) { toast('Enter chat ID', 'error'); setLoading(btn, false); return; }
  var r = await apiCall('POST', '/api/chats/' + encodeURIComponent(id) + '/read');
  setLoading(btn, false);
  showResult('resultGetChat', r.data || r);
};

window.apiAddParticipant = async function(btn) {
  setLoading(btn, true);
  var chatId = document.getElementById('partChatId').value, userId = document.getElementById('partUserId').value;
  if (!chatId || !userId) { toast('Enter chat ID and user ID', 'error'); setLoading(btn, false); return; }
  var r = await apiCall('POST', '/api/chats/' + encodeURIComponent(chatId) + '/participants', { userId: userId });
  setLoading(btn, false);
  showResult('resultGetChat', r.data || r);
};

window.apiRemoveParticipant = async function(btn) {
  setLoading(btn, true);
  var chatId = document.getElementById('partChatId').value, userId = document.getElementById('partUserId').value;
  if (!chatId || !userId) { toast('Enter chat ID and user ID', 'error'); setLoading(btn, false); return; }
  var r = await apiCall('DELETE', '/api/chats/' + encodeURIComponent(chatId) + '/participants/' + encodeURIComponent(userId));
  setLoading(btn, false);
  showResult('resultGetChat', r.data || r);
};

window.apiUpdateRole = async function(btn) {
  setLoading(btn, true);
  var chatId = document.getElementById('roleChatId').value, userId = document.getElementById('roleUserId').value, role = document.getElementById('roleType').value;
  if (!chatId || !userId) { toast('Enter chat ID and user ID', 'error'); setLoading(btn, false); return; }
  var r = await apiCall('PUT', '/api/chats/' + encodeURIComponent(chatId) + '/participants/' + encodeURIComponent(userId) + '/role', { role: role });
  setLoading(btn, false);
  showResult('resultGetChat', r.data || r);
};

window.apiGetChatNotifications = async function(btn) {
  setLoading(btn, true);
  var id = document.getElementById('notifChatId').value;
  if (!id) { toast('Enter chat ID', 'error'); setLoading(btn, false); return; }
  var r = await apiCall('GET', '/api/chats/' + encodeURIComponent(id) + '/notifications');
  setLoading(btn, false);
  showResult('resultGetChat', r.data || r);
};

window.apiUpdateChatNotifications = async function(btn) {
  setLoading(btn, true);
  var id = document.getElementById('notifChatId').value, mode = document.getElementById('notifMode').value;
  if (!id) { toast('Enter chat ID', 'error'); setLoading(btn, false); return; }
  var r = await apiCall('PUT', '/api/chats/' + encodeURIComponent(id) + '/notifications', { mode: mode });
  setLoading(btn, false);
  showResult('resultGetChat', r.data || r);
};

window.apiListMessages = async function(btn) {
  setLoading(btn, true);
  var chatId = document.getElementById('msgChatIdList').value;
  if (!chatId) { toast('Enter chat ID', 'error'); setLoading(btn, false); return; }
  var limit = document.getElementById('msgLimit').value, offset = document.getElementById('msgOffset').value;
  var path = '/api/chats/' + encodeURIComponent(chatId) + '/messages';
  var params = [];
  if (limit) params.push('limit=' + encodeURIComponent(limit));
  if (offset) params.push('offset=' + encodeURIComponent(offset));
  if (params.length) path += '?' + params.join('&');
  var r = await apiCall('GET', path);
  setLoading(btn, false);
  showResult('resultListMessages', r.data || r);
};

window.apiSendMessage = async function(btn) {
  setLoading(btn, true);
  var chatId = document.getElementById('sendMsgChatId').value;
  if (!chatId) { toast('Enter chat ID', 'error'); setLoading(btn, false); return; }
  var body = { type: document.getElementById('sendMsgType').value, content: document.getElementById('sendMsgContent').value };
  if (!body.content) { toast('Enter content', 'error'); setLoading(btn, false); return; }
  var r = await apiCall('POST', '/api/chats/' + encodeURIComponent(chatId) + '/messages', body);
  setLoading(btn, false);
  showResult('resultListMessages', r.data || r);
};

window.apiUploadFile = async function(btn) {
  setLoading(btn, true);
  var chatId = document.getElementById('fileMsgChatId').value, file = document.getElementById('fileMsgFile').files[0];
  if (!chatId) { toast('Enter chat ID', 'error'); setLoading(btn, false); return; }
  if (!file) { toast('Select a file', 'error'); setLoading(btn, false); return; }
  var fd = new FormData(); fd.append('file', file);
  var r = await apiCall('POST', '/api/chats/' + encodeURIComponent(chatId) + '/messages/file', fd, true);
  setLoading(btn, false);
  showResult('resultListMessages', r.data || r);
};

window.apiEditMessage = async function(btn) {
  setLoading(btn, true);
  var id = document.getElementById('msgIdAction').value, content = document.getElementById('editMsgContent').value;
  if (!id) { toast('Enter message ID', 'error'); setLoading(btn, false); return; }
  if (!content) { toast('Enter new content', 'error'); setLoading(btn, false); return; }
  var r = await apiCall('PUT', '/api/messages/' + encodeURIComponent(id), { content: content });
  setLoading(btn, false);
  showResult('resultListMessages', r.data || r);
};

window.apiDeleteMessage = async function(btn) {
  var id = document.getElementById('msgIdAction').value;
  if (!id) { toast('Enter message ID', 'error'); return; }
  if (!confirm('Delete message?')) return;
  setLoading(btn, true);
  var r = await apiCall('DELETE', '/api/messages/' + encodeURIComponent(id));
  setLoading(btn, false);
  showResult('resultListMessages', r.data || r);
};

window.apiResendMessage = async function(btn) {
  setLoading(btn, true);
  var chatId = document.getElementById('resendChatId').value, msgId = document.getElementById('resendMsgId').value;
  if (!chatId || !msgId) { toast('Enter chat ID and message ID', 'error'); setLoading(btn, false); return; }
  var r = await apiCall('POST', '/api/chats/' + encodeURIComponent(chatId) + '/messages/' + encodeURIComponent(msgId) + '/resend');
  setLoading(btn, false);
  showResult('resultListMessages', r.data || r);
};

window.apiSearchMessages = async function(btn) {
  setLoading(btn, true);
  var chatId = document.getElementById('searchMsgChatId').value, q = document.getElementById('searchMsgQuery').value;
  if (!chatId) { toast('Enter chat ID', 'error'); setLoading(btn, false); return; }
  if (!q) { toast('Enter query', 'error'); setLoading(btn, false); return; }
  var r = await apiCall('GET', '/api/chats/' + encodeURIComponent(chatId) + '/messages/search?q=' + encodeURIComponent(q));
  setLoading(btn, false);
  showResult('resultSearchMessages', r.data || r);
};

window.apiAddReaction = async function(btn) {
  setLoading(btn, true);
  var msgId = document.getElementById('reactMsgId').value, emoji = document.getElementById('reactEmoji').value;
  if (!msgId) { toast('Enter message ID', 'error'); setLoading(btn, false); return; }
  if (!emoji) { toast('Select an emoji', 'error'); setLoading(btn, false); return; }
  var r = await apiCall('POST', '/api/messages/' + encodeURIComponent(msgId) + '/reactions', { emoji: emoji });
  setLoading(btn, false);
  showResult('resultListMessages', r.data || r);
};

window.apiRemoveReaction = async function(btn) {
  setLoading(btn, true);
  var msgId = document.getElementById('reactMsgId').value, emoji = document.getElementById('reactEmoji').value;
  if (!msgId) { toast('Enter message ID', 'error'); setLoading(btn, false); return; }
  if (!emoji) { toast('Select emoji', 'error'); setLoading(btn, false); return; }
  var r = await apiCall('DELETE', '/api/messages/' + encodeURIComponent(msgId) + '/reactions?emoji=' + encodeURIComponent(emoji));
  setLoading(btn, false);
  showResult('resultListMessages', r.data || r);
};

window.apiPinMessage = async function(btn) {
  setLoading(btn, true);
  var id = document.getElementById('pinMsgId').value;
  if (!id) { toast('Enter message ID', 'error'); setLoading(btn, false); return; }
  var r = await apiCall('PUT', '/api/messages/' + encodeURIComponent(id) + '/pin');
  setLoading(btn, false);
  showResult('resultListMessages', r.data || r);
};

window.apiUnpinMessage = async function(btn) {
  setLoading(btn, true);
  var id = document.getElementById('pinMsgId').value;
  if (!id) { toast('Enter message ID', 'error'); setLoading(btn, false); return; }
  var r = await apiCall('DELETE', '/api/messages/' + encodeURIComponent(id) + '/pin');
  setLoading(btn, false);
  showResult('resultListMessages', r.data || r);
};

window.apiMarkMsgRead = async function(btn) {
  setLoading(btn, true);
  var id = document.getElementById('pinMsgId').value;
  if (!id) { toast('Enter message ID', 'error'); setLoading(btn, false); return; }
  var r = await apiCall('POST', '/api/messages/' + encodeURIComponent(id) + '/read');
  setLoading(btn, false);
  showResult('resultListMessages', r.data || r);
};

window.apiForwardMessage = async function(btn) {
  setLoading(btn, true);
  var msgId = document.getElementById('forwardMsgId').value, chatId = document.getElementById('forwardChatId').value;
  if (!msgId || !chatId) { toast('Enter message ID and target chat ID', 'error'); setLoading(btn, false); return; }
  var r = await apiCall('POST', '/api/messages/' + encodeURIComponent(msgId) + '/forward', { chatId: chatId });
  setLoading(btn, false);
  showResult('resultListMessages', r.data || r);
};

window.apiGetPinnedMessages = async function(btn) {
  setLoading(btn, true);
  var chatId = document.getElementById('pinnedChatId').value;
  if (!chatId) { toast('Enter chat ID', 'error'); setLoading(btn, false); return; }
  var r = await apiCall('GET', '/api/chats/' + encodeURIComponent(chatId) + '/pinned');
  setLoading(btn, false);
  showResult('resultPinned', r.data || r);
};

window.apiInitiateCall = async function(btn) {
  setLoading(btn, true);
  var body = { chatId: document.getElementById('callChatId').value, type: document.getElementById('callType').value };
  if (!body.chatId) { toast('Enter chat ID', 'error'); setLoading(btn, false); return; }
  var r = await apiCall('POST', '/api/calls/initiate', body);
  setLoading(btn, false);
  showResult('resultInitCall', r.data || r);
};

window.apiRespondCall = async function(btn) {
  setLoading(btn, true);
  var id = document.getElementById('callIdRespond').value, action = document.getElementById('callResponse').value;
  if (!id) { toast('Enter call ID', 'error'); setLoading(btn, false); return; }
  var r = await apiCall('POST', '/api/calls/' + encodeURIComponent(id) + '/respond', { action: action });
  setLoading(btn, false);
  showResult('resultInitCall', r.data || r);
};

window.apiEndCall = async function(btn) {
  setLoading(btn, true);
  var id = document.getElementById('callIdEnd').value;
  if (!id) { toast('Enter call ID', 'error'); setLoading(btn, false); return; }
  var r = await apiCall('POST', '/api/calls/' + encodeURIComponent(id) + '/end');
  setLoading(btn, false);
  showResult('resultInitCall', r.data || r);
};

window.apiGetCall = async function(btn) {
  setLoading(btn, true);
  var id = document.getElementById('callIdGet').value;
  if (!id) { toast('Enter call ID', 'error'); setLoading(btn, false); return; }
  var r = await apiCall('GET', '/api/calls/' + encodeURIComponent(id));
  setLoading(btn, false);
  showResult('resultGetCall', r.data || r);
};

window.apiGetCallHistory = async function(btn) {
  setLoading(btn, true);
  var chatId = document.getElementById('callHistoryChatId').value;
  if (!chatId) { toast('Enter chat ID', 'error'); setLoading(btn, false); return; }
  var r = await apiCall('GET', '/api/calls/history/' + encodeURIComponent(chatId));
  setLoading(btn, false);
  showResult('resultCallHistory', r.data || r);
};

window.apiBlockUser = async function(btn) {
  setLoading(btn, true);
  var userId = document.getElementById('blockUserId').value;
  if (!userId) { toast('Enter user ID', 'error'); setLoading(btn, false); return; }
  var r = await apiCall('POST', '/api/users/block', { blockedId: userId });
  setLoading(btn, false);
  showResult('resultBlock', r.data || r);
};

window.apiUnblockUser = async function(btn) {
  setLoading(btn, true);
  var userId = document.getElementById('unblockUserId').value;
  if (!userId) { toast('Enter user ID', 'error'); setLoading(btn, false); return; }
  var r = await apiCall('DELETE', '/api/users/block/' + encodeURIComponent(userId));
  setLoading(btn, false);
  showResult('resultBlock', r.data || r);
};

window.apiListBlocked = async function(btn) {
  setLoading(btn, true);
  var r = await apiCall('GET', '/api/users/blocked');
  setLoading(btn, false);
  showResult('resultBlocked', r.data || r);
};

window.apiSearchUsers = async function(btn) {
  setLoading(btn, true);
  var q = document.getElementById('searchQuery').value;
  if (!q) { toast('Enter query', 'error'); setLoading(btn, false); return; }
  var r = await apiCall('GET', '/api/users/search?q=' + encodeURIComponent(q));
  setLoading(btn, false);
  showResult('resultSearchUsers', r.data || r);
};
