;(function(global) {
  'use strict';

  var toastContainer = document.getElementById('toastContainer');
  var tokenSection = document.getElementById('tokenSection');
  var tokenDisplay = document.getElementById('tokenDisplay');
  var tokenInfo = document.getElementById('tokenInfo');
  var tokenInfoText = document.getElementById('tokenInfoText');
  var connectionStatus = document.getElementById('connectionStatus');
  var statusIndicator = document.getElementById('statusIndicator');

  init();

  function init() {
    updateTokenDisplay();
    updateConnectionStatus();
    bindTabSwitching();
  }

  function bindTabSwitching() {
    document.querySelectorAll('.sidebar .tab').forEach(function(tab) {
      tab.addEventListener('click', function() {
        document.querySelectorAll('.sidebar .tab').forEach(function(t) { t.classList.remove('active'); });
        document.querySelectorAll('.tab-content').forEach(function(c) { c.classList.remove('active'); });
        tab.classList.add('active');
        var target = document.getElementById('tab-' + tab.dataset.tab);
        if (target) target.classList.add('active');
      });
    });
  }

  function updateTokenDisplay() {
    var token = global._api.getToken();
    if (token) {
      tokenSection.classList.add('visible');
      tokenDisplay.textContent = token;
      try {
        var parts = token.split('.');
        if (parts.length === 3) {
          var payload = JSON.parse(atob(parts[1]));
          tokenInfo.style.display = 'flex';
          var exp = payload.exp ? new Date(payload.exp * 1000).toLocaleString() : 'N/A';
          tokenInfoText.textContent = 'Subject: ' + (payload.sub || payload.id || 'N/A') + ' | Expires: ' + exp;
        }
      } catch(e) {}
    } else {
      tokenSection.classList.remove('visible');
      tokenDisplay.textContent = '';
      tokenInfo.style.display = 'none';
      tokenInfoText.textContent = '';
      localStorage.removeItem('messenger_token');
    }
  }

  function clearToken() {
    global._api.clearToken();
    updateTokenDisplay();
    updateConnectionStatus();
    toast('Token cleared', 'info');
  }

  function copyToken() {
    var token = global._api.getToken();
    if (!token) { toast('No token to copy', 'error'); return; }
    navigator.clipboard.writeText(token).then(function() {
      toast('Token copied to clipboard', 'success');
    }).catch(function() {
      toast('Failed to copy token', 'error');
    });
  }

  function updateConnectionStatus() {
    if (global._api.getToken()) {
      statusIndicator.className = 'status-dot online';
      connectionStatus.textContent = 'Authenticated';
    } else {
      statusIndicator.className = 'status-dot offline';
      connectionStatus.textContent = 'Not connected';
    }
  }

  function toast(msg, type) {
    type = type || 'info';
    var icons = { error: '\u26A0\uFE0F', success: '\u2705', info: '\uD83D\uDCAC' };
    var iconHtml = icons[type] || icons.info;
    var t = document.createElement('div');
    t.className = 'toast ' + type;
    t.innerHTML = '<span class="toast-icon">' + iconHtml + '</span><span class="toast-msg">' + escapeHtml(msg) + '</span>';
    toastContainer.appendChild(t);
    setTimeout(function() {
      t.style.opacity = '0';
      t.style.transform = 'translateX(100%)';
      t.style.transition = 'all .3s ease';
      setTimeout(function() { t.remove(); }, 300);
    }, 3500);
  }

  function setLoading(btn, loading) {
    if (!btn) return;
    if (loading) {
      btn._origContent = btn.innerHTML;
      btn.disabled = true;
      btn.innerHTML = '<span class="loading-spinner"></span> Loading...';
    } else {
      btn.disabled = false;
      btn.innerHTML = btn._origContent || btn.textContent;
    }
  }

  function escapeHtml(s) {
    if (typeof s !== 'string') s = String(s);
    return s.replace(/&/g,'&amp;').replace(/</g,'&lt;').replace(/>/g,'&gt;').replace(/"/g,'&quot;').replace(/'/g,'&#39;');
  }

  function syntaxHighlight(obj, indent) {
    indent = indent || 0;
    var pad = '  '.repeat(indent);
    if (obj === null || obj === undefined) return '<span class="hl-null">null</span>';
    if (typeof obj === 'string') return '<span class="hl-str">"' + escapeHtml(obj) + '"</span>';
    if (typeof obj === 'number') return '<span class="hl-num">' + obj + '</span>';
    if (typeof obj === 'boolean') return '<span class="hl-bool">' + obj + '</span>';
    if (Array.isArray(obj)) {
      if (obj.length === 0) return '<span class="hl-bracket">[</span><span class="hl-bracket">]</span>';
      var items = obj.map(function(item) {
        return pad + '  ' + syntaxHighlight(item, indent + 1);
      });
      return '<span class="hl-bracket">[</span>\n' + items.join(',\n') + '\n' + pad + '<span class="hl-bracket">]</span>';
    }
    if (typeof obj === 'object') {
      var keys = Object.keys(obj);
      if (keys.length === 0) return '<span class="hl-bracket">{</span><span class="hl-bracket">}</span>';
      var items = keys.map(function(k) {
        return pad + '  <span class="hl-key">"' + escapeHtml(k) + '"</span><span class="hl-punct">: </span>' + syntaxHighlight(obj[k], indent + 1);
      });
      return '<span class="hl-bracket">{</span>\n' + items.join(',\n') + '\n' + pad + '<span class="hl-bracket">}</span>';
    }
    return escapeHtml(String(obj));
  }

  function showResult(id, data) {
    var el = document.getElementById(id);
    if (!el) return;
    if (data === undefined || data === null) { el.innerHTML = '<span class="hl-null">No response</span>'; return; }
    if (typeof data === 'string') { try { data = JSON.parse(data); } catch(e) {} }
    if (data && typeof data === 'object' && data.error) {
      el.innerHTML = '<span class="hl-bracket">{</span>\n  <span class="hl-key">"error"</span>: <span class="hl-str">"' + escapeHtml(data.data ? (typeof data.data === 'object' ? JSON.stringify(data.data) : '' + data.data) : 'Request failed') + '"</span>\n<span class="hl-bracket">}</span>';
      return;
    }
    el.innerHTML = syntaxHighlight(data);
  }

  function setEmoji(btn, emoji) {
    document.getElementById('reactEmoji').value = emoji;
    document.querySelectorAll('.emoji-grid button').forEach(function(b) { b.classList.remove('selected'); });
    if (btn) btn.classList.add('selected');
  }

  function toggleEditMsg() {
    var row = document.getElementById('editMsgRow');
    row.style.display = row.style.display === 'none' ? 'flex' : 'none';
  }

  // Export to window
  global.updateTokenDisplay = updateTokenDisplay;
  global.updateConnectionStatus = updateConnectionStatus;
  global.clearToken = clearToken;
  global.copyToken = copyToken;
  global.toast = toast;
  global.setLoading = setLoading;
  global.showResult = showResult;
  global.setEmoji = setEmoji;
  global.toggleEditMsg = toggleEditMsg;
})(window);
