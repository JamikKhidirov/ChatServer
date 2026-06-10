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
    var u = idVal('regUsername'), e = idVal('regEmail'), p = idVal('regPassword'), dn = idVal('regDisplayName');
    if (!u || !e || !p) { toast('Fill all required fields', 'error'); setLoading(btn, false); return; }
    var body = { username: u, email: e, password: p, display_name: dn };
    var r = await apiCall('POST', '/api/auth/register', body);
    if (!r.error && r.data && r.data.token) { _api.setToken(r.data.token); updateTokenDisplay(); updateConnectionStatus(); }
    showRes('resultAuth', btn, r);
  };

  global.apiLogin = async function(btn) {
    setLoading(btn, true);
    var body = { email: idVal('loginEmail'), password: idVal('loginPassword') };
    var r = await apiCall('POST', '/api/auth/login', body);
    if (!r.error && r.data && r.data.token) { authToken = r.data.token; updateTokenDisplay(); updateConnectionStatus(); toast('Login successful', 'success'); }
    showRes('resultAuth', btn, r);
  };

  global.apiSendEmailLoginCode = async function(btn) {
    setLoading(btn, true);
    var email = idVal('loginEmailCodeAddr');
    if (!email) { toast('Enter email', 'error'); setLoading(btn, false); return; }
    var r = await apiCall('POST', '/api/auth/login/email', { email: email });
    showRes('resultAuth', btn, r);
  };

  global.apiVerifyEmailLoginCode = async function(btn) {
    setLoading(btn, true);
    var email = idVal('loginEmailCodeAddr');
    var code = idVal('loginEmailCodeVerify');
    if (!email || !code) { toast('Enter email and code', 'error'); setLoading(btn, false); return; }
    var r = await apiCall('POST', '/api/auth/login/email/verify', { email: email, code: code });
    if (!r.error && r.data && r.data.token) { authToken = r.data.token; updateTokenDisplay(); updateConnectionStatus(); toast('Email login successful', 'success'); }
    showRes('resultAuth', btn, r);
  };

  global.apiSendPhoneLoginCode = async function(btn) {
    setLoading(btn, true);
    var phone = idVal('loginPhoneCodeNum');
    if (!phone) { toast('Enter phone number', 'error'); setLoading(btn, false); return; }
    var r = await apiCall('POST', '/api/auth/login/phone', { phone: phone });
    showRes('resultAuth', btn, r);
  };

  global.apiVerifyPhoneLoginCode = async function(btn) {
    setLoading(btn, true);
    var phone = idVal('loginPhoneCodeNum');
    var code = idVal('loginPhoneCodeVerify');
    if (!phone || !code) { toast('Enter phone and code', 'error'); setLoading(btn, false); return; }
    var r = await apiCall('POST', '/api/auth/login/phone/verify', { phone: phone, code: code });
    if (!r.error && r.data && r.data.token) { authToken = r.data.token; updateTokenDisplay(); updateConnectionStatus(); toast('Phone login successful', 'success'); }
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
    var token = idVal('pushToken'), provider = idVal('pushPlatform');
    if (!token || !provider) { toast('Fill token and platform', 'error'); setLoading(btn, false); return; }
    var body = { token: token, provider: provider };
    var r = await apiCall('PUT', '/api/users/push-token', body);
    showRes('resultProfile', btn, r);
  };

  global.apiTestPush = async function(btn) {
    setLoading(btn, true);
    var r = await apiCall('POST', '/api/users/push-test', { title: 'Test Notification', body: 'This is a test push from API tester' });
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
    var id = idVal('notifChatId'), muted = idVal('notifMode') === 'none';
    if (!id) { toast('Enter chat ID', 'error'); setLoading(btn, false); return; }
    var r = await apiCall('PUT', '/api/chats/' + encodeURIComponent(id) + '/notifications', { muted: muted });
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

  // ====== NEW FEATURES ======

  // Pinned chats
  global.apiPinChat = async function(btn) {
    setLoading(btn, true);
    var id = idVal('pinChatId');
    if (!id) { toast('Enter chat ID', 'error'); setLoading(btn, false); return; }
    var r = await apiCall('POST', '/api/chats/' + encodeURIComponent(id) + '/pin');
    showRes('resultChatActions', btn, r);
  };
  global.apiUnpinChat = async function(btn) {
    setLoading(btn, true);
    var id = idVal('pinChatId');
    if (!id) { toast('Enter chat ID', 'error'); setLoading(btn, false); return; }
    var r = await apiCall('DELETE', '/api/chats/' + encodeURIComponent(id) + '/pin');
    showRes('resultChatActions', btn, r);
  };

  // Archive
  global.apiArchiveChat = async function(btn) {
    setLoading(btn, true);
    var id = idVal('archiveChatId');
    if (!id) { toast('Enter chat ID', 'error'); setLoading(btn, false); return; }
    var r = await apiCall('POST', '/api/chats/' + encodeURIComponent(id) + '/archive');
    showRes('resultChatActions', btn, r);
  };
  global.apiUnarchiveChat = async function(btn) {
    setLoading(btn, true);
    var id = idVal('archiveChatId');
    if (!id) { toast('Enter chat ID', 'error'); setLoading(btn, false); return; }
    var r = await apiCall('POST', '/api/chats/' + encodeURIComponent(id) + '/unarchive');
    showRes('resultChatActions', btn, r);
  };
  global.apiListArchived = async function(btn) {
    setLoading(btn, true);
    var r = await apiCall('GET', '/api/chats/archived');
    showRes('resultChatActions', btn, r);
  };

  // Transfer ownership
  global.apiTransferOwnership = async function(btn) {
    setLoading(btn, true);
    var chatId = idVal('transferChatId'), userId = idVal('transferUserId');
    if (!chatId || !userId) { toast('Enter chat ID and user ID', 'error'); setLoading(btn, false); return; }
    var r = await apiCall('POST', '/api/chats/' + encodeURIComponent(chatId) + '/transfer-ownership', { userId: userId });
    showRes('resultChatActions', btn, r);
  };

  // Star messages
  global.apiStarMessage = async function(btn) {
    setLoading(btn, true);
    var id = idVal('starMsgId');
    if (!id) { toast('Enter message ID', 'error'); setLoading(btn, false); return; }
    var r = await apiCall('POST', '/api/messages/' + encodeURIComponent(id) + '/star');
    showRes('resultStarred', btn, r);
  };
  global.apiUnstarMessage = async function(btn) {
    setLoading(btn, true);
    var id = idVal('starMsgId');
    if (!id) { toast('Enter message ID', 'error'); setLoading(btn, false); return; }
    var r = await apiCall('DELETE', '/api/messages/' + encodeURIComponent(id) + '/star');
    showRes('resultStarred', btn, r);
  };
  global.apiGetStarred = async function(btn) {
    setLoading(btn, true);
    var r = await apiCall('GET', '/api/messages/starred');
    showRes('resultStarred', btn, r);
  };

  // Delete for me
  global.apiDeleteForMe = async function(btn) {
    setLoading(btn, true);
    var id = idVal('deleteForMeMsgId');
    if (!id) { toast('Enter message ID', 'error'); setLoading(btn, false); return; }
    var r = await apiCall('DELETE', '/api/messages/' + encodeURIComponent(id) + '/for-me');
    showRes('resultListMessages', btn, r);
  };

  // Search all messages
  global.apiSearchAllMessages = async function(btn) {
    setLoading(btn, true);
    var q = idVal('searchAllQ');
    if (!q) { toast('Enter query', 'error'); setLoading(btn, false); return; }
    var r = await apiCall('GET', '/api/messages/search?q=' + encodeURIComponent(q));
    showRes('resultListMessages', btn, r);
  };

  // Export chat
  global.apiExportChat = async function(btn) {
    setLoading(btn, true);
    var id = idVal('exportChatId');
    if (!id) { toast('Enter chat ID', 'error'); setLoading(btn, false); return; }
    var r = await apiCall('GET', '/api/chats/' + encodeURIComponent(id) + '/export');
    showRes('resultExportChat', btn, r);
  };

  // Chat media
  global.apiGetChatMedia = async function(btn) {
    setLoading(btn, true);
    var chatId = idVal('mediaChatId'), type = idVal('mediaType');
    if (!chatId) { toast('Enter chat ID', 'error'); setLoading(btn, false); return; }
    var path = '/api/chats/' + encodeURIComponent(chatId) + '/media';
    if (type) path += '?type=' + encodeURIComponent(type);
    var r = await apiCall('GET', path);
    showRes('resultMedia', btn, r);
  };

  // Polls
  global.apiCreatePoll = async function(btn) {
    setLoading(btn, true);
    var chatId = idVal('pollChatId');
    if (!chatId) { toast('Enter chat ID', 'error'); setLoading(btn, false); return; }
    var raw = idVal('pollOptions');
    if (!raw) { toast('Enter options JSON array', 'error'); setLoading(btn, false); return; }
    var opts;
    try { opts = JSON.parse(raw); } catch(e) { toast('Invalid JSON', 'error'); setLoading(btn, false); return; }
    var body = { chatId: chatId, question: idVal('pollQuestion'), options: opts, isAnonymous: idChk('pollAnonymous'), multipleChoice: idChk('pollMultiple') };
    var r = await apiCall('POST', '/api/chats/' + encodeURIComponent(chatId) + '/polls', body);
    showRes('resultPolls', btn, r);
  };
  global.apiGetPolls = async function(btn) {
    setLoading(btn, true);
    var chatId = idVal('pollChatIdGet');
    if (!chatId) { toast('Enter chat ID', 'error'); setLoading(btn, false); return; }
    var r = await apiCall('GET', '/api/chats/' + encodeURIComponent(chatId) + '/polls');
    showRes('resultPolls', btn, r);
  };
  global.apiVotePoll = async function(btn) {
    setLoading(btn, true);
    var pollId = idVal('votePollId'), idx = idVal('voteOptionIdx');
    if (!pollId || idx === '') { toast('Enter poll ID and option index', 'error'); setLoading(btn, false); return; }
    var r = await apiCall('POST', '/api/polls/' + encodeURIComponent(pollId) + '/vote', { optionIndex: parseInt(idx) });
    showRes('resultPolls', btn, r);
  };
  global.apiClosePoll = async function(btn) {
    setLoading(btn, true);
    var pollId = idVal('closePollId');
    if (!pollId) { toast('Enter poll ID', 'error'); setLoading(btn, false); return; }
    var r = await apiCall('POST', '/api/polls/' + encodeURIComponent(pollId) + '/close');
    showRes('resultPolls', btn, r);
  };

  // Stickers
  global.apiCreateStickerPack = async function(btn) {
    setLoading(btn, true);
    var r = await apiCall('POST', '/api/stickers/packs', { name: idVal('stickPackName'), animated: idChk('stickPackAnimated') });
    showRes('resultStickers', btn, r);
  };
  global.apiListStickerPacks = async function(btn) {
    setLoading(btn, true);
    var r = await apiCall('GET', '/api/stickers/packs');
    showRes('resultStickers', btn, r);
  };
  global.apiGetMyStickerPacks = async function(btn) {
    setLoading(btn, true);
    var r = await apiCall('GET', '/api/stickers/packs/my');
    showRes('resultStickers', btn, r);
  };
  global.apiGetStickerPack = async function(btn) {
    setLoading(btn, true);
    var id = idVal('stickPackIdGet');
    if (!id) { toast('Enter pack ID', 'error'); setLoading(btn, false); return; }
    var r = await apiCall('GET', '/api/stickers/packs/' + encodeURIComponent(id));
    showRes('resultStickers', btn, r);
  };
  global.apiAddSticker = async function(btn) {
    setLoading(btn, true);
    var packId = idVal('stickPackIdAdd');
    if (!packId) { toast('Enter pack ID', 'error'); setLoading(btn, false); return; }
    var r = await apiCall('POST', '/api/stickers/packs/' + encodeURIComponent(packId) + '/stickers', { emoji: idVal('stickEmoji'), imageUrl: idVal('stickImageUrl') });
    showRes('resultStickers', btn, r);
  };
  global.apiDeleteStickerPack = async function(btn) {
    if (!confirm('Delete pack?')) return;
    setLoading(btn, true);
    var id = idVal('stickPackIdDel');
    if (!id) { toast('Enter pack ID', 'error'); setLoading(btn, false); return; }
    var r = await apiCall('DELETE', '/api/stickers/packs/' + encodeURIComponent(id));
    showRes('resultStickers', btn, r);
  };
  global.apiGetStickerLibrary = async function(btn) {
    setLoading(btn, true);
    var r = await apiCall('GET', '/api/stickers/library');
    showRes('resultStickers', btn, r);
  };
  global.apiAddStickerToLibrary = async function(btn) {
    setLoading(btn, true);
    var id = idVal('stickLibAddId');
    if (!id) { toast('Enter sticker ID', 'error'); setLoading(btn, false); return; }
    var r = await apiCall('POST', '/api/stickers/library', { stickerId: id });
    showRes('resultStickers', btn, r);
  };

  // Drafts
  global.apiSaveDraft = async function(btn) {
    setLoading(btn, true);
    var r = await apiCall('POST', '/api/drafts', { chatId: idVal('draftChatId'), content: idVal('draftContent') });
    showRes('resultDrafts', btn, r);
  };
  global.apiGetDraft = async function(btn) {
    setLoading(btn, true);
    var chatId = idVal('draftGetChatId');
    if (!chatId) { toast('Enter chat ID', 'error'); setLoading(btn, false); return; }
    var r = await apiCall('GET', '/api/drafts?chatId=' + encodeURIComponent(chatId));
    showRes('resultDrafts', btn, r);
  };

  // Scheduled messages
  global.apiScheduleMessage = async function(btn) {
    setLoading(btn, true);
    var body = { chatId: idVal('schedChatId'), content: idVal('schedContent'), type: idVal('schedType'), scheduledAt: idVal('schedAt') };
    if (!body.chatId || !body.content || !body.scheduledAt) { toast('Fill required fields', 'error'); setLoading(btn, false); return; }
    var r = await apiCall('POST', '/api/messages/schedule', body);
    showRes('resultScheduled', btn, r);
  };
  global.apiGetScheduled = async function(btn) {
    setLoading(btn, true);
    var r = await apiCall('GET', '/api/messages/scheduled');
    showRes('resultScheduled', btn, r);
  };
  global.apiCancelScheduled = async function(btn) {
    setLoading(btn, true);
    var id = idVal('cancelSchedId');
    if (!id) { toast('Enter scheduled message ID', 'error'); setLoading(btn, false); return; }
    var r = await apiCall('DELETE', '/api/messages/scheduled/' + encodeURIComponent(id));
    showRes('resultScheduled', btn, r);
  };

  // Sessions
  global.apiGetSessions = async function(btn) {
    setLoading(btn, true);
    var r = await apiCall('GET', '/api/sessions');
    showRes('resultSessions', btn, r);
  };
  global.apiDeleteSession = async function(btn) {
    setLoading(btn, true);
    var id = idVal('delSessionId');
    if (!id) { toast('Enter session ID', 'error'); setLoading(btn, false); return; }
    var r = await apiCall('DELETE', '/api/sessions/' + encodeURIComponent(id));
    showRes('resultSessions', btn, r);
  };
  global.apiDeleteAllSessions = async function(btn) {
    if (!confirm('Delete all sessions?')) return;
    setLoading(btn, true);
    var r = await apiCall('DELETE', '/api/sessions');
    showRes('resultSessions', btn, r);
  };

  // Bots
  global.apiCreateBot = async function(btn) {
    setLoading(btn, true);
    var r = await apiCall('POST', '/api/bots', { name: idVal('botName'), webhookUrl: idVal('botWebhook') });
    showRes('resultBots', btn, r);
  };
  global.apiGetMyBots = async function(btn) {
    setLoading(btn, true);
    var r = await apiCall('GET', '/api/bots');
    showRes('resultBots', btn, r);
  };
  global.apiUpdateBot = async function(btn) {
    setLoading(btn, true);
    var id = idVal('botIdUpdate');
    if (!id) { toast('Enter bot ID', 'error'); setLoading(btn, false); return; }
    var body = {};
    var n = idVal('botNameUpdate'); if (n) body.name = n;
    var w = idVal('botWebhookUpdate'); if (w) body.webhookUrl = w;
    var r = await apiCall('PUT', '/api/bots/' + encodeURIComponent(id), body);
    showRes('resultBots', btn, r);
  };
  global.apiDeleteBot = async function(btn) {
    if (!confirm('Delete bot?')) return;
    setLoading(btn, true);
    var id = idVal('botIdDel');
    if (!id) { toast('Enter bot ID', 'error'); setLoading(btn, false); return; }
    var r = await apiCall('DELETE', '/api/bots/' + encodeURIComponent(id));
    showRes('resultBots', btn, r);
  };
  global.apiRegenerateBotToken = async function(btn) {
    setLoading(btn, true);
    var id = idVal('botIdRegen');
    if (!id) { toast('Enter bot ID', 'error'); setLoading(btn, false); return; }
    var r = await apiCall('POST', '/api/bots/' + encodeURIComponent(id) + '/regenerate-token');
    showRes('resultBots', btn, r);
  };

  // Saved GIFs
  global.apiSaveGif = async function(btn) {
    setLoading(btn, true);
    var url = idVal('gifUrl');
    if (!url) { toast('Enter GIF URL', 'error'); setLoading(btn, false); return; }
    var r = await apiCall('POST', '/api/gifs', { url: url });
    showRes('resultGifs', btn, r);
  };
  global.apiGetSavedGifs = async function(btn) {
    setLoading(btn, true);
    var r = await apiCall('GET', '/api/gifs');
    showRes('resultGifs', btn, r);
  };
  global.apiDeleteGif = async function(btn) {
    setLoading(btn, true);
    var url = idVal('gifDelUrl');
    if (!url) { toast('Enter GIF URL', 'error'); setLoading(btn, false); return; }
    var r = await apiCall('DELETE', '/api/gifs', { url: url });
    showRes('resultGifs', btn, r);
  };

  // === New Feature API Functions ===

  // Captcha
  global.apiGenerateCaptcha = async function(btn) {
    setLoading(btn, true);
    var r = await apiCall('GET', '/api/captcha/generate');
    showRes('resultSecurity', btn, r);
  };
  global.apiVerifyCaptcha = async function(btn) {
    setLoading(btn, true);
    var token = idVal('captchaToken');
    var sol = idVal('captchaSolution');
    if (!token || !sol) { toast('Fill token and solution', 'error'); setLoading(btn, false); return; }
    var r = await apiCall('POST', '/api/captcha/verify', { token: token, solution: sol });
    showRes('resultSecurity', btn, r);
  };

  // E2E Encryption
  global.apiRegisterE2EKey = async function(btn) {
    setLoading(btn, true);
    var pubKey = idVal('e2ePubKey');
    var privKey = idVal('e2ePrivKey');
    if (!pubKey || !privKey) { toast('Fill public and private key', 'error'); setLoading(btn, false); return; }
    var r = await apiCall('POST', '/api/e2e/keys', { public_key: pubKey, private_key_encrypted: privKey });
    showRes('resultSecurity', btn, r);
  };
  global.apiGetE2EPublicKey = async function(btn) {
    setLoading(btn, true);
    var userId = idVal('e2eUserId');
    if (!userId) { toast('Enter user ID', 'error'); setLoading(btn, false); return; }
    var r = await apiCall('GET', '/api/e2e/keys/' + encodeURIComponent(userId));
    showRes('resultSecurity', btn, r);
  };

  // Email Verification
  global.apiSendEmailVerification = async function(btn) {
    setLoading(btn, true);
    var email = idVal('verEmail');
    if (!email) { toast('Enter email', 'error'); setLoading(btn, false); return; }
    var r = await apiCall('POST', '/api/verification/email/send', { email: email });
    showRes('resultSecurity', btn, r);
  };
  global.apiVerifyEmail = async function(btn) {
    setLoading(btn, true);
    var code = idVal('verEmailCode');
    if (!code) { toast('Enter code', 'error'); setLoading(btn, false); return; }
    var r = await apiCall('POST', '/api/verification/email/verify', { code: code });
    showRes('resultSecurity', btn, r);
  };

  // Phone Verification
  global.apiSendPhoneVerification = async function(btn) {
    setLoading(btn, true);
    var phone = idVal('verPhone');
    if (!phone) { toast('Enter phone', 'error'); setLoading(btn, false); return; }
    var r = await apiCall('POST', '/api/verification/phone/send', { phone: phone });
    showRes('resultSecurity', btn, r);
  };
  global.apiVerifyPhone = async function(btn) {
    setLoading(btn, true);
    var code = idVal('verPhoneCode');
    if (!code) { toast('Enter code', 'error'); setLoading(btn, false); return; }
    var r = await apiCall('POST', '/api/verification/phone/verify', { code: code });
    showRes('resultSecurity', btn, r);
  };

  // Bookmarks
  global.apiBookmarkMessage = async function(btn) {
    setLoading(btn, true);
    var msgId = idVal('bmMsgId');
    var chatId = idVal('bmChatId');
    if (!msgId || !chatId) { toast('Fill message and chat ID', 'error'); setLoading(btn, false); return; }
    var r = await apiCall('POST', '/api/bookmarks', { message_id: msgId, chat_id: chatId });
    showRes('resultBookmarks', btn, r);
  };
  global.apiGetBookmarks = async function(btn) {
    setLoading(btn, true);
    var r = await apiCall('GET', '/api/bookmarks');
    showRes('resultBookmarks', btn, r);
  };
  global.apiRemoveBookmark = async function(btn) {
    setLoading(btn, true);
    var msgId = idVal('bmDelMsgId');
    if (!msgId) { toast('Enter message ID', 'error'); setLoading(btn, false); return; }
    var r = await apiCall('DELETE', '/api/bookmarks/' + encodeURIComponent(msgId));
    showRes('resultBookmarks', btn, r);
  };

  // Reports
  global.apiCreateReport = async function(btn) {
    setLoading(btn, true);
    var msgId = idVal('repMsgId');
    var reason = idVal('repReason');
    var desc = idVal('repDesc');
    if (!msgId || !reason) { toast('Fill message ID and reason', 'error'); setLoading(btn, false); return; }
    var r = await apiCall('POST', '/api/reports', { message_id: msgId, reason: reason, description: desc });
    showRes('resultReports', btn, r);
  };
  global.apiListReports = async function(btn) {
    setLoading(btn, true);
    var status = idVal('repStatusFilter');
    var r = await apiCall('GET', '/api/reports?status=' + encodeURIComponent(status));
    showRes('resultReports', btn, r);
  };
  global.apiResolveReport = async function(btn) {
    setLoading(btn, true);
    var id = idVal('repId');
    var status = idVal('repResolveStatus');
    if (!id || !status) { toast('Fill report ID and status', 'error'); setLoading(btn, false); return; }
    var r = await apiCall('POST', '/api/reports/' + encodeURIComponent(id) + '/resolve', { status: status });
    showRes('resultReports', btn, r);
  };

  // Self-Destruct
  global.apiSetSelfDestruct = async function(btn) {
    setLoading(btn, true);
    var msgId = idVal('sdMsgId');
    var seconds = idVal('sdSeconds');
    if (!msgId || !seconds) { toast('Fill message ID and seconds', 'error'); setLoading(btn, false); return; }
    var r = await apiCall('POST', '/api/messages/self-destruct', { message_id: msgId, delete_after: parseInt(seconds) });
    showRes('resultMessages', btn, r);
  };

  // Edit History
  global.apiGetEditHistory = async function(btn) {
    setLoading(btn, true);
    var msgId = idVal('editHistMsgId');
    if (!msgId) { toast('Enter message ID', 'error'); setLoading(btn, false); return; }
    var r = await apiCall('GET', '/api/messages/' + encodeURIComponent(msgId) + '/history');
    showRes('resultMessages', btn, r);
  };

  // Link Preview
  global.apiGetLinkPreview = async function(btn) {
    setLoading(btn, true);
    var url = idVal('previewUrl');
    if (!url) { toast('Enter URL', 'error'); setLoading(btn, false); return; }
    var r = await apiCall('GET', '/api/preview?url=' + encodeURIComponent(url));
    showRes('resultPreview', btn, r);
  };

  // Admin
  global.apiAdminDashboard = async function(btn) {
    setLoading(btn, true);
    var r = await apiCall('GET', '/api/admin/dashboard');
    showRes('resultAdmin', btn, r);
  };
  global.apiAdminUsers = async function(btn) {
    setLoading(btn, true);
    var r = await apiCall('GET', '/api/admin/users');
    showRes('resultAdmin', btn, r);
  };
  global.apiAdminMessages = async function(btn) {
    setLoading(btn, true);
    var r = await apiCall('GET', '/api/admin/messages');
    showRes('resultAdmin', btn, r);
  };
  global.apiAdminReadMessage = async function(btn) {
    setLoading(btn, true);
    var id = idVal('adminMsgId');
    if (!id) { toast('Enter message ID', 'error'); setLoading(btn, false); return; }
    var r = await apiCall('GET', '/api/admin/messages/' + encodeURIComponent(id));
    showRes('resultAdmin', btn, r);
  };
  global.apiAdminBanUser = async function(btn) {
    setLoading(btn, true);
    var uid = idVal('adminBanUserId');
    var reason = idVal('adminBanReason');
    if (!uid || !reason) { toast('Fill user ID and reason', 'error'); setLoading(btn, false); return; }
    var r = await apiCall('POST', '/api/admin/users/ban', { user_id: uid, reason: reason });
    showRes('resultAdmin', btn, r);
  };
  global.apiAdminUnbanUser = async function(btn) {
    setLoading(btn, true);
    var uid = idVal('adminUnbanUserId');
    if (!uid) { toast('Enter user ID', 'error'); setLoading(btn, false); return; }
    var r = await apiCall('POST', '/api/admin/users/unban/' + encodeURIComponent(uid));
    showRes('resultAdmin', btn, r);
  };
  global.apiAdminSettings = async function(btn) {
    setLoading(btn, true);
    var r = await apiCall('GET', '/api/admin/settings');
    showRes('resultAdmin', btn, r);
  };
  global.apiAdminUpdateSetting = async function(btn) {
    setLoading(btn, true);
    var key = idVal('adminSetKey');
    var val = idVal('adminSetVal');
    if (!key || !val) { toast('Fill key and value', 'error'); setLoading(btn, false); return; }
    var r = await apiCall('PUT', '/api/admin/settings', { key: key, value: val });
    showRes('resultAdmin', btn, r);
  };
  global.apiAdminLogs = async function(btn) {
    setLoading(btn, true);
    var r = await apiCall('GET', '/api/admin/logs');
    showRes('resultAdmin', btn, r);
  };
  global.apiAdminIPBlocks = async function(btn) {
    setLoading(btn, true);
    var r = await apiCall('GET', '/api/admin/ip-blocks');
    showRes('resultAdmin', btn, r);
  };
  global.apiAdminUnblockIP = async function(btn) {
    setLoading(btn, true);
    var ip = idVal('adminUnblockIP');
    if (!ip) { toast('Enter IP address', 'error'); setLoading(btn, false); return; }
    var r = await apiCall('POST', '/api/admin/ip-blocks/' + encodeURIComponent(ip) + '/unblock');
    showRes('resultAdmin', btn, r);
  };

  // Expose internals for app.js
  global._api = { getToken: function() { return authToken; }, setToken: function(t) { authToken = t; localStorage.setItem('messenger_token', t); }, clearToken: function() { authToken = ''; localStorage.removeItem('messenger_token'); } };

  // Helpers
  function idVal(id) { return (document.getElementById(id) || {}).value || ''; }
  function idChk(id) { return (document.getElementById(id) || {}).checked || false; }
  function idFile(id) { var el = document.getElementById(id); return (el && el.files && el.files[0]) || null; }
})(window);
