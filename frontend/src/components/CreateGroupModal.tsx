import { useState, useEffect, useRef } from 'react'
import { api } from '../api/client'

interface Props {
  onClose: () => void
  onCreated: () => void
}

export default function CreateGroupModal({ onClose, onCreated }: Props) {
  const [name, setName] = useState('')
  const [description, setDescription] = useState('')
  const [searchQuery, setSearchQuery] = useState('')
  const [searchResults, setSearchResults] = useState<any[]>([])
  const [selectedUsers, setSelectedUsers] = useState<any[]>([])
  const [creating, setCreating] = useState(false)
  const [error, setError] = useState('')
  const inputRef = useRef<HTMLInputElement>(null)

  useEffect(() => { inputRef.current?.focus() }, [])

  const handleSearch = async () => {
    if (!searchQuery.trim()) return
    const res = await api.call('GET', '/api/users/search?q=' + encodeURIComponent(searchQuery))
    const d = res.data as any
    if (!res.error && Array.isArray(d?.data)) {
      setSearchResults(d.data)
    }
  }

  const toggleUser = (user: any) => {
    setSelectedUsers(prev => {
      const exists = prev.some(u => (u.id || u.ID) === (user.id || user.ID))
      if (exists) return prev.filter(u => (u.id || u.ID) !== (user.id || user.ID))
      return [...prev, user]
    })
  }

  const handleCreate = async () => {
    if (!name.trim()) { setError('Group name is required'); return }
    if (selectedUsers.length === 0) { setError('Select at least one participant'); return }
    setCreating(true)
    setError('')
    const ids = selectedUsers.map(u => u.id || u.ID)
    const res = await api.call('POST', '/api/chats', {
      type: 'group', name: name.trim(), description: description.trim() || undefined, participantIds: ids,
    })
    setCreating(false)
    if (!res.error) {
      onCreated()
      onClose()
    } else {
      setError(typeof res.data === 'string' ? res.data : 'Failed to create group')
    }
  }

  return (
    <div className="modal-overlay" onClick={onClose}>
      <div className="modal-content create-group-modal" onClick={e => e.stopPropagation()}>
        <div className="modal-header">
          <h3>New Group</h3>
          <button className="modal-close-btn" onClick={onClose}>✕</button>
        </div>

        <div className="modal-body">
          {error && <div className="modal-error">{error}</div>}

          <div className="form-field">
            <label>Group Name</label>
            <input
              ref={inputRef}
              type="text" value={name}
              onChange={e => setName(e.target.value)}
              placeholder="Enter group name"
              maxLength={64}
            />
          </div>

          <div className="form-field">
            <label>Description (optional)</label>
            <input
              type="text" value={description}
              onChange={e => setDescription(e.target.value)}
              placeholder="What is this group about?"
              maxLength={512}
            />
          </div>

          <div className="form-field">
            <label>Add Members</label>
            <div className="member-search-row">
              <input
                type="text" value={searchQuery}
                onChange={e => setSearchQuery(e.target.value)}
                onKeyDown={e => e.key === 'Enter' && handleSearch()}
                placeholder="Search users..."
              />
              <button className="search-member-btn" onClick={handleSearch}>Search</button>
            </div>
          </div>

          {selectedUsers.length > 0 && (
            <div className="selected-members">
              {selectedUsers.map(u => (
                <span key={u.id || u.ID} className="member-chip" onClick={() => toggleUser(u)}>
                  {u.displayName || u.username} ✕
                </span>
              ))}
            </div>
          )}

          {searchResults.length > 0 && (
            <div className="member-search-results">
              {searchResults.map((u: any) => {
                const isSelected = selectedUsers.some(s => (s.id || s.ID) === (u.id || u.ID))
                return (
                  <div
                    key={u.id || u.ID}
                    className={`member-search-item ${isSelected ? 'selected' : ''}`}
                    onClick={() => toggleUser(u)}
                  >
                    <div className="member-search-avatar">{u.displayName?.charAt(0) || '?'}</div>
                    <div className="member-search-info">
                      <div className="member-search-name">{u.displayName || u.Username}</div>
                      <div className="member-search-username">@{u.username || u.Username}</div>
                    </div>
                    <div className="member-check">{isSelected ? '✓' : ''}</div>
                  </div>
                )
              })}
            </div>
          )}
        </div>

        <div className="modal-footer">
          <button className="modal-cancel-btn" onClick={onClose}>Cancel</button>
          <button className="modal-submit-btn" onClick={handleCreate} disabled={creating || !name.trim() || selectedUsers.length === 0}>
            {creating ? <span className="spinner-sm" /> : 'Create Group'}
          </button>
        </div>
      </div>
    </div>
  )
}
