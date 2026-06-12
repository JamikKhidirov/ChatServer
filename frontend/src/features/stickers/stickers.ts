import { apiCall } from '../../api/client'

export async function createStickerPack(data: { name: string }) {
  return apiCall('POST', '/api/stickers/packs', data)
}

export async function listStickerPacks() {
  return apiCall('GET', '/api/stickers/packs')
}

export async function getMyStickerPacks() {
  return apiCall('GET', '/api/stickers/packs/my')
}

export async function getStickerPack(packId: string) {
  return apiCall('GET', `/api/stickers/packs/${packId}`)
}

export async function addSticker(packId: string, data: { emoji: string; imageUrl?: string }) {
  return apiCall('POST', `/api/stickers/packs/${packId}/stickers`, data)
}

export async function deleteStickerPack(packId: string) {
  return apiCall('DELETE', `/api/stickers/packs/${packId}`)
}

export async function getStickerLibrary() {
  return apiCall('GET', '/api/stickers/library')
}

export async function addStickerToLibrary(stickerId: string) {
  return apiCall('POST', '/api/stickers/library', { stickerId })
}
