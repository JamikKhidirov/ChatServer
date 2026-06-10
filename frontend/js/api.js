;(function(global) {
  'use strict';

  var authToken = localStorage.getItem('messenger_token') || '';
  var currentEmoji = '';

  // Derive API base from current page when served from same origin
  var baseURL = '';
  if (global.location && global.location.host) {
    baseURL = global.location.protocol + '//' + global.location.host;
  }

  async function apiCall(method, path, body, isFormData) {
    var fullURL = baseURL + path;
    var needsAuth = !path.startsWith('/api/auth/register') && !path.startsWith('/api/auth/login');
    if (needsAuth && !authToken) { toast('Please login first', 'error'); throw new Error('No auth'); }
    var opts = { method: method };
    var headers = {};
    if (authToken) headers['Authorization'] = 'Bearer ' + authToken;
    if (body && !isFormData) { headers['Content-Type'] = 'application/json'; opts.body = JSON.stringify(body); }
    else if (body && isFormData) { opts.body = body; }
    if (Object.keys(headers).length > 0) opts.headers = headers;
    try {
      var res = await fetch(fullURL, opts);
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

  function errorCheck(btn, r) {
    if (btn) setLoading(btn, false);
    if (r && r.error && r.status === 401) {
      authToken = ''; updateTokenDisplay(); updateConnectionStatus();
      localStorage.removeItem('messenger_token');
    }
  }

  function showRes(id, btn, r) {
    if (btn) setLoading(btn, false);
    showResult(id, (r && !r.error) ? r.data : (r || null));
  }

  // Auth
  global.apiRegister = async function(btn) {
    setLoading(btn, true);
    var body = { username: idVal('regUsername'), email: idVal('regEmail'), password: idVal('regPassword'), displayName: idVal('regDisplayName') };
    var r = await apiCall('POST', '/api/auth/register', body);
    if (!r.error && r.data && r.data.token) { authToken = r.data.token; updateTokenDisplay(); updateConnectionStatus(); }
    showRes('resultAuth', btn, r);
  };

  global.apiLogin = async function(btn) {
    setLoading(btn, true);
    var body = { email: idVal('loginEmail'), password: idVal('loginPassword') };
    var r = await apiCall('POST', '/api/auth/login', body);
    if (!r.error && r.data && r.data.token) { authToken = r.data.token; updateTokenDisplay(); updateConnectionStatus(); toast('Login successful', 'success'); }
    showRes('resultAuth', btn, r);
  };

  global.apiRefreshToken = async function(btn) {
    setLoading(btn, true);
    var r = await apiCall('GET', '/api/auth/refresh');
    if (!r.error && r.data && r.data.token) { authToken = r.data.token; updateTokenDisplay(); toast('Token refreshed', 'success'); }
    showRes('resultAuth', btn, r);
  };

  global.apiChangePassword = async function(btn) {
    setLoading(btn, true);
    var body = { oldPassword: idVal('oldPassword'), newPassword: idVal('newPassword') };
    var r = await apiCall('PUT', '/api/auth/change-password', body);
    showRes('resultAuth', btn, r);
    if (!r.error) toast('Password changed', 'success');
  };

  // Profile
  global.apiGetProfile = async function(btn) {
    setLoading(btn, true);
    var r = await apiCall('GET', '/api/users/profile');
    showRes('resultProfile', btn, r);
  };

  global.apiUpdateProfile = async function(btn) {
    setLoading(btn, true);
    var body = {};
    var map = { displayName:'updDisplayName', bio:'updBio', phone:'updPhone', gender:'updGender', dateOfBirth:'updDob' };
    for (var k in map) { if (!map.hasOwnProperty(k)) continue; var v = idVal(map[k]); if (v) body[k] = v; }
    if (Object.keys(body).length === 0) { toast('Fill at least one field', 'error'); setLoading(btn, false); return; }
    var r = await apiCall('PUT', '/api/users/profile', body);
    showRes('resultProfile', btn, r);
  };

  global.apiUploadAvatar = async function(btn) {
    setLoading(btn, true);
    var file = idFile('avatarFile');
    if (!file) { toast('Select a file', 'error'); setLoading(btn, false); return; }
    var fd = new FormData(); fd.append('avatar', file);
    var r = await apiCall('POST', '/api/users/avatar', fd, true);
    showRes('resultProfile', btn, r);
  };

  global.apiUpdateStatus = async function(btn) {
    setLoading(btn, true);
    var body = { text: idVal('statusText'), status: idVal('statusType') };
    var r = await apiCall('PUT', '/api/users/status', body);
    showRes('resultProfile', btn, r);
  };

  global.apiDeleteAccount = async function(btn) {
    if (!confirm('Delete your account permanently?')) return;
    setLoading(btn, true);
    var r = await apiCall('DELETE', '/api/users/account');
    if (!r.error) { authToken = ''; updateTokenDisplay(); updateConnectionStatus(); }
    showRes('resultProfile', btn, r);
  };

  global.apiGetUserById = async function(btn) {
    setLoading(btn, true);
    var id = idVal('userIdLookup');
    if (!id) { toast('Enter a user ID', 'error'); setLoading(btn, false); return; }
    var r = await apiCall('GET', '/api/users/' + encodeURIComponent(id));
    showRes('resultProfile', btn, r);
  };

  global.apiGetUserByUsername = async function(btn) {
    setLoading(btn, true);
    var u = idVal('usernameLookup');
    if (!u) { toast('Enter username', 'error'); setLoading(btn, false); return; }
    var r = await apiCall('GET', '/api/users/username/' + encodeURIComponent(u));
    showRes('resultProfile', btn, r);
  };

  global.apiSavePushToken = async function(btn) {
    setLoading(btn, true);
    var body = { token: idVal('pushToken'), platform: idVal('pushPlatform') };
    var r = await apiCall('POST', '/api/users/push-token', body);
    showRes('resultProfile', btn, r);
  };

  global.apiTestPush = async function(btn) {
    setLoading(btn, true);
    var r = await apiCall('POST', '/api/users/push-test');
    showRes('resultProfile', btn, r);
  };

  // Settings
  global.apiGetSettings = async function(btn) {
    setLoading(btn, true);
    var r = await apiCall('GET', '/api/account/settings');
    showRes('resultSettings', btn, r);
  };

  global.apiUpdateSettings = async function(btn) {
    setLoading(btn, true);
    var body = {};
    var lang = idVal('setLanguage'), theme = idVal('setTheme');
    if (lang) body.language = lang;
    if (theme) body.theme = theme;
    body.notifications = idChk('setNotifications');
    body.soundEnabled = idChk('setSoundEnabled');
    body.lastSeenMode = idVal('setLastSeenMode');
    var r = await apiCall('PUT', '/api/account/settings', body);
    showRes('resultSettings', btn, r);
  };

  // Contacts
  global.apiSyncContacts = async function(btn) {
    setLoading(btn, true);
    var raw = idVal('contactsSyncData');
    if (!raw) { toast('Enter JSON data', 'error'); setLoading(btn, false); return; }
    var data;
    try { data = JSON.parse(raw); } catch(e) { toast('Invalid JSON', 'error'); setLoading(btn, false); return; }
    var r = await apiCall('POST', '/api/contacts/sync', data);
    showRes('resultSyncContacts', btn, r);
  };

  global.apiGetContacts = async function(btn) {
    setLoading(btn, true);
    var r = await apiCall('GET', '/api/contacts');
    showRes('resultContacts', btn, r);
  };

  global.apiSearchContacts = async function(btn) {
    setLoading(btn, true);
    var q = idVal('contactSearchQuery');
    if (!q) { toast('Enter query', 'error'); setLoading(btn, false); return; }
    var r = await apiCall('GET', '/api/contacts/search?q=' + encodeURIComponent(q));
    showRes('resultSearchContacts', btn, r);
  };

  global.apiFindRegistered = async function(btn) {
    setLoading(btn, true);
    var r = await apiCall('GET', '/api/contacts/registered');
    showRes('resultRegistered', btn, r);
  };

  // Chats
  global.apiSearchChats = async function(btn) {
    setLoading(btn, true);
    var q = idVal('chatSearchQuery');
    if (!q) { toast('Enter search query', 'error'); setLoading(btn, false); return; }
    var r = await apiCall('GET', '/api/chats/search?q=' + encodeURIComponent(q));
    showRes('resultSearchChats', btn, r);
  };

  global.apiListChats = async function(btn) {
    setLoading(btn, true);
    var r = await apiCall('GET', '/api/chats');
    showRes('resultListChats', btn, r);
  };

  global.apiCreatePrivateChat = async function(btn) {
    setLoading(btn, true);
    var userId = idVal('privateChatUserId');
    if (!userId) { toast('Enter user ID', 'error'); setLoading(btn, false); return; }
    var body = { type: 'private', participantIds: [userId] };
    var r = await apiCall('POST', '/api/chats', body);
    showRes('resultListChats', btn, r);
  };

  global.apiCreateGroupChat = async function(btn) {
    setLoading(btn, true);
    var name = idVal('groupName'), raw = idVal('groupParticipants');
    if (!name) { toast('Enter group name', 'error'); setLoading(btn, false); return; }
    if (!raw) { toast('Enter participant JSON', 'error'); setLoading(btn, false); return; }
    var participants;
    try { participants = JSON.parse(raw); } catch(e) { toast('Invalid JSON', 'error'); setLoading(btn, false); return; }
    if (!Array.isArray(participants) || !participants.length) { toast('Non-empty array required', 'error'); setLoading(btn, false); return; }
    var r = await apiCall('POST', '/api/chats', { type: 'group', name: name, participantIds: participants });
    showRes('resultListChats', btn, r);
  };

  global.apiGetChat = async function(btn) {
    setLoading(btn, true);
    var id = idVal('chatIdGet');
    if (!id) { toast('Enter chat ID', 'error'); setLoading(btn, false); return; }
    var r = await apiCall('GET', '/api/chats/' + encodeURIComponent(id));
    showRes('resultGetChat', btn, r);
  };

  global.apiHideChat = async function(btn) {
    setLoading(btn, true);
    var id = idVal('chatIdAction');
    if (!id) { toast('Enter chat ID', 'error'); setLoading(btn, false); return; }
    var r = await apiCall('POST', '/api/chats/' + encodeURIComponent(id) + '/hide');
    showRes('resultGetChat', btn, r);
  };

  global.apiDeleteChat = async function(btn) {
    var id = idVal('chatIdAction');
    if (!id) { toast('Enter chat ID', 'error'); return; }
    if (!confirm('Delete chat?')) return;
    setLoading(btn, true);
    var r = await apiCall('DELETE', '/api/chats/' + encodeURIComponent(id));
    showRes('resultGetChat', btn, r);
  };

  global.apiLeaveGroup = async function(btn) {
    setLoading(btn, true);
    var id = idVal('chatIdAction');
    if (!id) { toast('Enter chat ID', 'error'); setLoading(btn, false); return; }
    var r = await apiCall('POST', '/api/chats/' + encodeURIComponent(id) + '/leave');
    showRes('resultGetChat', btn, r);
  };

  global.apiMarkChatRead = async function(btn) {
    setLoading(btn, true);
    var id = idVal('chatIdAction');
    if (!id) { toast('Enter chat ID', 'error'); setLoading(btn, false); return; }
    var r = await apiCall('POST', '/api/chats/' + encodeURIComponent(id) + '/read');
    showRes('resultGetChat', btn, r);
  };

  global.apiAddParticipant = async function(btn) {
    setLoading(btn, true);
    var chatId = idVal('partChatId'), userId = idVal('partUserId');
    if (!chatId || !userId) { toast('Enter chat ID and user ID', 'error'); setLoading(btn, false); return; }
    var r = await apiCall('POST', '/api/chats/' + encodeURIComponent(chatId) + '/participants', { userId: userId });
    showRes('resultGetChat', btn, r);
  };

  global.apiRemoveParticipant = async function(btn) {
    setLoading(btn, true);
    var chatId = idVal('partChatId'), userId = idVal('partUserId');
    if (!chatId || !userId) { toast('Enter chat ID and user ID', 'error'); setLoading(btn, false); return; }
    var r = await apiCall('DELETE', '/api/chats/' + encodeURIComponent(chatId) + '/participants/' + encodeURIComponent(userId));
    showRes('resultGetChat', btn, r);
  };

  global.apiUpdateRole = async function(btn) {
    setLoading(btn, true);
    var chatId = idVal('roleChatId'), userId = idVal('roleUserId'), role = idVal('roleType');
    if (!chatId || !userId) { toast('Enter chat ID and user ID', 'error'); setLoading(btn, false); return; }
    var r = await apiCall('PUT', '/api/chats/' + encodeURIComponent(chatId) + '/participants/' + encodeURIComponent(userId) + '/role', { role: role });
    showRes('resultGetChat', btn, r);
  };

  global.apiGetChatNotifications = async function(btn) {
    setLoading(btn, true);
    var id = idVal('notifChatId');
    if (!id) { toast('Enter chat ID', 'error'); setLoading(btn, false); return; }
    var r = await apiCall('GET', '/api/chats/' + encodeURIComponent(id) + '/notifications');
    showRes('resultGetChat', btn, r);
  };

  global.apiUpdateChatNotifications = async function(btn) {
    setLoading(btn, true);
    var id = idVal('notifChatId'), mode = idVal('notifMode');
    if (!id) { toast('Enter chat ID', 'error'); setLoading(btn, false); return; }
    var r = await apiCall('PUT', '/api/chats/' + encodeURIComponent(id) + '/notifications', { mode: mode });
    showRes('resultGetChat', btn, r);
  };

  // Messages
  global.apiListMessages = async function(btn) {
    setLoading(btn, true);
    var chatId = idVal('msgChatIdList');
    if (!chatId) { toast('Enter chat ID', 'error'); setLoading(btn, false); return; }
    var limit = idVal('msgLimit'), offset = idVal('msgOffset');
    var path = '/api/chats/' + encodeURIComponent(chatId) + '/messages';
    var params = [];
    if (limit) params.push('limit=' + encodeURIComponent(limit));
    if (offset) params.push('offset=' + encodeURIComponent(offset));
    if (params.length) path += '?' + params.join('&');
    var r = await apiCall('GET', path);
    showRes('resultListMessages', btn, r);
  };

  global.apiSendMessage = async function(btn) {
    setLoading(btn, true);
    var chatId = idVal('sendMsgChatId');
    if (!chatId) { toast('Enter chat ID', 'error'); setLoading(btn, false); return; }
    var body = { type: idVal('sendMsgType'), content: idVal('sendMsgContent') };
    if (!body.content) { toast('Enter content', 'error'); setLoading(btn, false); return; }
    var r = await apiCall('POST', '/api/chats/' + encodeURIComponent(chatId) + '/messages', body);
    showRes('resultListMessages', btn, r);
  };

  global.apiUploadFile = async function(btn) {
    setLoading(btn, true);
    var chatId = idVal('fileMsgChatId'), file = idFile('fileMsgFile');
    if (!chatId) { toast('Enter chat ID', 'error'); setLoading(btn, false); return; }
    if (!file) { toast('Select a file', 'error'); setLoading(btn, false); return; }
    var fd = new FormData(); fd.append('file', file);
    var r = await apiCall('POST', '/api/chats/' + encodeURIComponent(chatId) + '/messages/file', fd, true);
    showRes('resultListMessages', btn, r);
  };

  global.apiEditMessage = async function(btn) {
    setLoading(btn, true);
    var id = idVal('msgIdAction'), content = idVal('editMsgContent');
    if (!id) { toast('Enter message ID', 'error'); setLoading(btn, false); return; }
    if (!content) { toast('Enter new content', 'error'); setLoading(btn, false); return; }
    var r = await apiCall('PUT', '/api/messages/' + encodeURIComponent(id), { content: content });
    showRes('resultListMessages', btn, r);
  };

  global.apiDeleteMessage = async function(btn) {
    var id = idVal('msgIdAction');
    if (!id) { toast('Enter message ID', 'error'); return; }
    if (!confirm('Delete message?')) return;
    setLoading(btn, true);
    var r = await apiCall('DELETE', '/api/messages/' + encodeURIComponent(id));
    showRes('resultListMessages', btn, r);
  };

  global.apiResendMessage = async function(btn) {
    setLoading(btn, true);
    var chatId = idVal('resendChatId'), msgId = idVal('resendMsgId');
    if (!chatId || !msgId) { toast('Enter chat ID and message ID', 'error'); setLoading(btn, false); return; }
    var r = await apiCall('POST', '/api/chats/' + encodeURIComponent(chatId) + '/messages/' + encodeURIComponent(msgId) + '/resend');
    showRes('resultListMessages', btn, r);
  };

  global.apiSearchMessages = async function(btn) {
    setLoading(btn, true);
    var chatId = idVal('searchMsgChatId'), q = idVal('searchMsgQuery');
    if (!chatId) { toast('Enter chat ID', 'error'); setLoading(btn, false); return; }
    if (!q) { toast('Enter query', 'error'); setLoading(btn, false); return; }
    var r = await apiCall('GET', '/api/chats/' + encodeURIComponent(chatId) + '/messages/search?q=' + encodeURIComponent(q));
    showRes('resultSearchMessages', btn, r);
  };

  global.apiAddReaction = async function(btn) {
    setLoading(btn, true);
    var msgId = idVal('reactMsgId'), emoji = idVal('reactEmoji');
    if (!msgId) { toast('Enter message ID', 'error'); setLoading(btn, false); return; }
    if (!emoji) { toast('Select an emoji', 'error'); setLoading(btn, false); return; }
    var r = await apiCall('POST', '/api/messages/' + encodeURIComponent(msgId) + '/reactions', { emoji: emoji });
    showRes('resultListMessages', btn, r);
  };

  global.apiRemoveReaction = async function(btn) {
    setLoading(btn, true);
    var msgId = idVal('reactMsgId'), emoji = idVal('reactEmoji');
    if (!msgId) { toast('Enter message ID', 'error'); setLoading(btn, false); return; }
    if (!emoji) { toast('Select emoji', 'error'); setLoading(btn, false); return; }
    var r = await apiCall('DELETE', '/api/messages/' + encodeURIComponent(msgId) + '/reactions?emoji=' + encodeURIComponent(emoji));
    showRes('resultListMessages', btn, r);
  };

  global.apiPinMessage = async function(btn) {
    setLoading(btn, true);
    var id = idVal('pinMsgId');
    if (!id) { toast('Enter message ID', 'error'); setLoading(btn, false); return; }
    var r = await apiCall('PUT', '/api/messages/' + encodeURIComponent(id) + '/pin');
    showRes('resultListMessages', btn, r);
  };

  global.apiUnpinMessage = async function(btn) {
    setLoading(btn, true);
    var id = idVal('pinMsgId');
    if (!id) { toast('Enter message ID', 'error'); setLoading(btn, false); return; }
    var r = await apiCall('DELETE', '/api/messages/' + encodeURIComponent(id) + '/pin');
    showRes('resultListMessages', btn, r);
  };

  global.apiMarkMsgRead = async function(btn) {
    setLoading(btn, true);
    var id = idVal('pinMsgId');
    if (!id) { toast('Enter message ID', 'error'); setLoading(btn, false); return; }
    var r = await apiCall('POST', '/api/messages/' + encodeURIComponent(id) + '/read');
    showRes('resultListMessages', btn, r);
  };

  global.apiForwardMessage = async function(btn) {
    setLoading(btn, true);
    var msgId = idVal('forwardMsgId'), chatId = idVal('forwardChatId');
    if (!msgId || !chatId) { toast('Enter message ID and target chat ID', 'error'); setLoading(btn, false); return; }
    var r = await apiCall('POST', '/api/messages/' + encodeURIComponent(msgId) + '/forward', { chatId: chatId });
    showRes('resultListMessages', btn, r);
  };

  global.apiGetPinnedMessages = async function(btn) {
    setLoading(btn, true);
    var chatId = idVal('pinnedChatId');
    if (!chatId) { toast('Enter chat ID', 'error'); setLoading(btn, false); return; }
    var r = await apiCall('GET', '/api/chats/' + encodeURIComponent(chatId) + '/pinned');
    showRes('resultPinned', btn, r);
  };

  // Calls
  global.apiInitiateCall = async function(btn) {
    setLoading(btn, true);
    var body = { chatId: idVal('callChatId'), type: idVal('callType') };
    if (!body.chatId) { toast('Enter chat ID', 'error'); setLoading(btn, false); return; }
    var r = await apiCall('POST', '/api/calls/initiate', body);
    showRes('resultInitCall', btn, r);
  };

  global.apiRespondCall = async function(btn) {
    setLoading(btn, true);
    var id = idVal('callIdRespond'), action = idVal('callResponse');
    if (!id) { toast('Enter call ID', 'error'); setLoading(btn, false); return; }
    var r = await apiCall('POST', '/api/calls/' + encodeURIComponent(id) + '/respond', { action: action });
    showRes('resultInitCall', btn, r);
  };

  global.apiEndCall = async function(btn) {
    setLoading(btn, true);
    var id = idVal('callIdEnd');
    if (!id) { toast('Enter call ID', 'error'); setLoading(btn, false); return; }
    var r = await apiCall('POST', '/api/calls/' + encodeURIComponent(id) + '/end');
    showRes('resultInitCall', btn, r);
  };

  global.apiGetCall = async function(btn) {
    setLoading(btn, true);
    var id = idVal('callIdGet');
    if (!id) { toast('Enter call ID', 'error'); setLoading(btn, false); return; }
    var r = await apiCall('GET', '/api/calls/' + encodeURIComponent(id));
    showRes('resultGetCall', btn, r);
  };

  global.apiGetCallHistory = async function(btn) {
    setLoading(btn, true);
    var chatId = idVal('callHistoryChatId');
    if (!chatId) { toast('Enter chat ID', 'error'); setLoading(btn, false); return; }
    var r = await apiCall('GET', '/api/calls/history/' + encodeURIComponent(chatId));
    showRes('resultCallHistory', btn, r);
  };

  // Block
  global.apiBlockUser = async function(btn) {
    setLoading(btn, true);
    var userId = idVal('blockUserId');
    if (!userId) { toast('Enter user ID', 'error'); setLoading(btn, false); return; }
    var r = await apiCall('POST', '/api/users/block', { blockedId: userId });
    showRes('resultBlock', btn, r);
  };

  global.apiUnblockUser = async function(btn) {
    setLoading(btn, true);
    var userId = idVal('unblockUserId');
    if (!userId) { toast('Enter user ID', 'error'); setLoading(btn, false); return; }
    var r = await apiCall('DELETE', '/api/users/block/' + encodeURIComponent(userId));
    showRes('resultBlock', btn, r);
  };

  global.apiListBlocked = async function(btn) {
    setLoading(btn, true);
    var r = await apiCall('GET', '/api/users/blocked');
    showRes('resultBlocked', btn, r);
  };

  // Search
  global.apiSearchUsers = async function(btn) {
    setLoading(btn, true);
    var q = idVal('searchQuery');
    if (!q) { toast('Enter query', 'error'); setLoading(btn, false); return; }
    var r = await apiCall('GET', '/api/users/search?q=' + encodeURIComponent(q));
    showRes('resultSearchUsers', btn, r);
  };

  // Expose internals for app.js
  global._api = { getToken: function() { return authToken; }, setToken: function(t) { authToken = t; localStorage.setItem('messenger_token', t); }, clearToken: function() { authToken = ''; localStorage.removeItem('messenger_token'); } };

  // Helpers
  function idVal(id) { return (document.getElementById(id) || {}).value || ''; }
  function idChk(id) { return (document.getElementById(id) || {}).checked || false; }
  function idFile(id) { var el = document.getElementById(id); return (el && el.files && el.files[0]) || null; }
})(window);
