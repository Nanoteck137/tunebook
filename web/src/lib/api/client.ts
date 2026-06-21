import { z } from "zod";
import * as api from "./types";
import { BaseApiClient, createUrl, type ExtraOptions } from "./base-client";


export class ApiClient extends BaseApiClient {
  url: ClientUrls;

  constructor(baseUrl: string) {
    super(baseUrl);
    this.url = new ClientUrls(baseUrl);
  }
  
  addItemToPlaylist(playlistId: string, body: api.AddItemToPlaylistBody, options?: ExtraOptions) {
    return this.request(`/api/v1/playlists/${playlistId}/items`, "POST", z.undefined(), z.any(), body, options)
  }
  
  addQueueItems(body: api.AddQueueItemsBody, options?: ExtraOptions) {
    return this.request("/api/v1/queue/items", "POST", z.undefined(), z.any(), body, options)
  }
  
  
  authClaimQuickConnectCode(body: api.AuthClaimQuickConnectCodeBody, options?: ExtraOptions) {
    return this.request("/api/v1/auth/quick-connect/claim", "POST", z.undefined(), z.any(), body, options)
  }
  
  authFinishProvider(body: api.AuthFinishProviderBody, options?: ExtraOptions) {
    return this.request("/api/v1/auth/providers/finish", "POST", api.AuthFinishProvider, z.any(), body, options)
  }
  
  authFinishQuickConnect(body: api.AuthFinishQuickConnectBody, options?: ExtraOptions) {
    return this.request("/api/v1/auth/quick-connect/finish", "POST", api.AuthFinishQuickConnect, z.any(), body, options)
  }
  
  authGetProviderStatus(body: api.AuthGetProviderStatusBody, options?: ExtraOptions) {
    return this.request("/api/v1/auth/provider/status", "POST", api.AuthGetProviderStatus, z.any(), body, options)
  }
  
  authGetProviders(options?: ExtraOptions) {
    return this.request("/api/v1/auth/providers", "GET", api.GetAuthProviders, z.any(), undefined, options)
  }
  
  authGetQuickConnectStatus(body: api.AuthGetQuickConnectStatusBody, options?: ExtraOptions) {
    return this.request("/api/v1/auth/quick-connect/status", "POST", api.AuthGetQuickConnectStatus, z.any(), body, options)
  }
  
  authProviderInitiate(body: api.AuthInitiateBody, options?: ExtraOptions) {
    return this.request("/api/v1/auth/providers/initiate", "POST", api.AuthInitiate, z.any(), body, options)
  }
  
  authQuickConnectInitiate(options?: ExtraOptions) {
    return this.request("/api/v1/auth/quick-connect/initiate", "POST", api.AuthQuickConnectInitiate, z.any(), undefined, options)
  }
  
  clearQueue(options?: ExtraOptions) {
    return this.request("/api/v1/queue", "DELETE", z.undefined(), z.any(), undefined, options)
  }
  
  createApiToken(body: api.CreateApiTokenBody, options?: ExtraOptions) {
    return this.request("/api/v1/me/apitokens", "POST", api.CreateApiToken, z.any(), body, options)
  }
  
  createPlaylist(body: api.CreatePlaylistBody, options?: ExtraOptions) {
    return this.request("/api/v1/playlists", "POST", api.CreatePlaylist, z.any(), body, options)
  }
  
  createTrackFilter(body: api.CreateTrackFilterBody, options?: ExtraOptions) {
    return this.request("/api/v1/me/filters/tracks", "POST", z.undefined(), z.any(), body, options)
  }
  
  deleteApiToken(tokenId: string, options?: ExtraOptions) {
    return this.request(`/api/v1/me/apitokens/${tokenId}`, "DELETE", z.undefined(), z.any(), undefined, options)
  }
  
  deletePlaylist(playlistId: string, options?: ExtraOptions) {
    return this.request(`/api/v1/playlists/${playlistId}`, "DELETE", z.undefined(), z.any(), undefined, options)
  }
  
  deleteTrackFilter(filterId: string, options?: ExtraOptions) {
    return this.request(`/api/v1/me/filters/tracks/${filterId}`, "DELETE", z.undefined(), z.any(), undefined, options)
  }
  
  editPlaylist(playlistId: string, body: api.EditPlaylistBody, options?: ExtraOptions) {
    return this.request(`/api/v1/playlists/${playlistId}`, "PATCH", z.undefined(), z.any(), body, options)
  }
  
  favoriteTrack(trackId: string, options?: ExtraOptions) {
    return this.request(`/api/v1/me/favorites/tracks/${trackId}`, "POST", z.undefined(), z.any(), undefined, options)
  }
  
  generatePlaylistImage(playlistId: string, options?: ExtraOptions) {
    return this.request(`/api/v1/playlists/${playlistId}/images/generate`, "POST", z.undefined(), z.any(), undefined, options)
  }
  
  getAlbumById(id: string, options?: ExtraOptions) {
    return this.request(`/api/v1/albums/${id}`, "GET", api.GetAlbumById, z.any(), undefined, options)
  }
  
  
  getAlbumTracks(id: string, options?: ExtraOptions) {
    return this.request(`/api/v1/albums/${id}/tracks`, "GET", api.GetAlbumTracks, z.any(), undefined, options)
  }
  
  getAlbums(options?: ExtraOptions) {
    return this.request("/api/v1/albums", "GET", api.GetAlbums, z.any(), undefined, options)
  }
  
  getApiTokens(options?: ExtraOptions) {
    return this.request("/api/v1/me/apitokens", "GET", api.GetApiTokens, z.any(), undefined, options)
  }
  
  getArtistAlbums(id: string, options?: ExtraOptions) {
    return this.request(`/api/v1/artists/${id}/albums`, "GET", api.GetArtistAlbumsById, z.any(), undefined, options)
  }
  
  getArtistById(id: string, options?: ExtraOptions) {
    return this.request(`/api/v1/artists/${id}`, "GET", api.GetArtistById, z.any(), undefined, options)
  }
  
  
  getArtists(options?: ExtraOptions) {
    return this.request("/api/v1/artists", "GET", api.GetArtists, z.any(), undefined, options)
  }
  
  getFavoriteTrackIds(options?: ExtraOptions) {
    return this.request("/api/v1/me/favorites/tracks/ids", "GET", api.GetFavoriteTrackIds, z.any(), undefined, options)
  }
  
  getMe(options?: ExtraOptions) {
    return this.request("/api/v1/auth/me", "GET", api.GetMe, z.any(), undefined, options)
  }
  
  getMediaSettings(options?: ExtraOptions) {
    return this.request("/api/v1/media/settings", "GET", api.GetMediaSettings, z.any(), undefined, options)
  }
  
  getPlaylistById(playlistId: string, options?: ExtraOptions) {
    return this.request(`/api/v1/playlists/${playlistId}`, "GET", api.GetPlaylistById, z.any(), undefined, options)
  }
  
  
  getPlaylistItemIds(playlistId: string, options?: ExtraOptions) {
    return this.request(`/api/v1/playlists/${playlistId}/ids`, "GET", api.GetPlaylistItemIds, z.any(), undefined, options)
  }
  
  getPlaylistItems(playlistId: string, options?: ExtraOptions) {
    return this.request(`/api/v1/playlists/${playlistId}/items`, "GET", api.GetPlaylistItems, z.any(), undefined, options)
  }
  
  getPlaylists(options?: ExtraOptions) {
    return this.request("/api/v1/playlists", "GET", api.GetPlaylists, z.any(), undefined, options)
  }
  
  getQueue(options?: ExtraOptions) {
    return this.request("/api/v1/queue", "GET", api.GetQueue, z.any(), undefined, options)
  }
  
  getQueueIds(options?: ExtraOptions) {
    return this.request("/api/v1/queue/ids", "GET", api.GetQueueIds, z.any(), undefined, options)
  }
  
  getQueueItemAtIndex(position: string, options?: ExtraOptions) {
    return this.request(`/api/v1/queue/items/${position}`, "GET", api.GetQueueItem, z.any(), undefined, options)
  }
  
  getQuickPlaylistIds(options?: ExtraOptions) {
    return this.request("/api/v1/me/quickplaylist", "GET", api.GetQuickPlaylistIds, z.any(), undefined, options)
  }
  
  getSystemInfo(options?: ExtraOptions) {
    return this.request("/api/v1/system/info", "GET", api.GetSystemInfo, z.any(), undefined, options)
  }
  
  getTrackById(id: string, options?: ExtraOptions) {
    return this.request(`/api/v1/tracks/${id}`, "GET", api.GetTrackById, z.any(), undefined, options)
  }
  
  getTrackHistory(options?: ExtraOptions) {
    return this.request("/api/v1/history/tracks", "GET", api.GetTrackHistory, z.any(), undefined, options)
  }
  
  getTrackHistoryById(id: string, options?: ExtraOptions) {
    return this.request(`/api/v1/history/tracks/${id}`, "GET", api.GetHistoryById, z.any(), undefined, options)
  }
  
  getTracks(options?: ExtraOptions) {
    return this.request("/api/v1/tracks", "GET", api.GetTracks, z.any(), undefined, options)
  }
  
  getUser(userId: string, options?: ExtraOptions) {
    return this.request(`/api/v1/users/${userId}`, "GET", api.GetUser, z.any(), undefined, options)
  }
  
  
  getUserStats(userId: string, options?: ExtraOptions) {
    return this.request(`/api/v1/users/${userId}/stats`, "GET", api.GetUserStats, z.any(), undefined, options)
  }
  
  getUserTrackFavorites(userId: string, options?: ExtraOptions) {
    return this.request(`/api/v1/users/${userId}/favorites/tracks`, "GET", api.GetUserFavorites, z.any(), undefined, options)
  }
  
  getUserTrackFilters(userId: string, options?: ExtraOptions) {
    return this.request(`/api/v1/users/${userId}/filters/tracks`, "GET", api.GetTrackFilters, z.any(), undefined, options)
  }
  
  pushTrackHistory(body: api.PushTrackHistoryBody, options?: ExtraOptions) {
    return this.request("/api/v1/history/tracks", "POST", api.PushTrackHistory, z.any(), body, options)
  }
  
  removePlaylistItem(playlistId: string, body: api.RemovePlaylistItemBody, options?: ExtraOptions) {
    return this.request(`/api/v1/playlists/${playlistId}/items`, "DELETE", z.undefined(), z.any(), body, options)
  }
  
  removeQueueItem(itemId: string, options?: ExtraOptions) {
    return this.request(`/api/v1/queue/items/${itemId}`, "DELETE", z.undefined(), z.any(), undefined, options)
  }
  
  reorderPlaylistItems(playlistId: string, body: api.ReorderPlaylistItemsBody, options?: ExtraOptions) {
    return this.request(`/api/v1/playlists/${playlistId}/items/reorder`, "POST", z.undefined(), z.any(), body, options)
  }
  
  replaceQueue(body: api.ReplaceQueueBody, options?: ExtraOptions) {
    return this.request("/api/v1/queue", "PUT", z.undefined(), z.any(), body, options)
  }
  
  runTask(taskName: string, options?: ExtraOptions) {
    return this.request(`/api/v1/system/task/${taskName}`, "POST", z.undefined(), z.any(), undefined, options)
  }
  
  searchAlbums(options?: ExtraOptions) {
    return this.request("/api/v1/search/albums", "GET", api.SearchAlbums, z.any(), undefined, options)
  }
  
  searchArtists(options?: ExtraOptions) {
    return this.request("/api/v1/search/artists", "GET", api.SearchArtists, z.any(), undefined, options)
  }
  
  searchPlaylists(options?: ExtraOptions) {
    return this.request("/api/v1/search/playlists", "GET", api.SearchPlaylists, z.any(), undefined, options)
  }
  
  searchTracks(options?: ExtraOptions) {
    return this.request("/api/v1/search/tracks", "GET", api.SearchTracks, z.any(), undefined, options)
  }
  
  searchUsers(options?: ExtraOptions) {
    return this.request("/api/v1/search/users", "GET", api.SearchUsers, z.any(), undefined, options)
  }
  
  setQueuePosition(body: api.SetQueuePositionBody, options?: ExtraOptions) {
    return this.request("/api/v1/queue/position", "PATCH", z.undefined(), z.any(), body, options)
  }
  
  setQuickPlaylist(body: api.SetQuickPlaylistBody, options?: ExtraOptions) {
    return this.request("/api/v1/me/quickplaylist", "POST", z.undefined(), z.any(), body, options)
  }
  
  
  
  unfavoriteTrack(trackId: string, options?: ExtraOptions) {
    return this.request(`/api/v1/me/favorites/tracks/${trackId}`, "DELETE", z.undefined(), z.any(), undefined, options)
  }
  
  updateMe(body: api.UpdateMeBody, options?: ExtraOptions) {
    return this.request("/api/v1/me", "PATCH", z.undefined(), z.any(), body, options)
  }
  
  updateTrackFilter(filterId: string, body: api.UpdateTrackFilterBody, options?: ExtraOptions) {
    return this.request(`/api/v1/me/filters/tracks/${filterId}`, "PATCH", z.undefined(), z.any(), body, options)
  }
  
  uploadPlaylistImage(playlistId: string, body: FormData, options?: ExtraOptions) {
    return this.requestForm(`/api/v1/playlists/${playlistId}/image/upload`, "POST", z.undefined(), z.any(), body, options)
  }
}

export class ClientUrls {
  baseUrl: string;

  constructor(baseUrl: string) {
    this.baseUrl = baseUrl;
  }
  
  addItemToPlaylist(playlistId: string) {
    return createUrl(this.baseUrl, `/api/v1/playlists/${playlistId}/items`)
  }
  
  addQueueItems() {
    return createUrl(this.baseUrl, "/api/v1/queue/items")
  }
  
  authCallback() {
    return createUrl(this.baseUrl, "/api/v1/auth/providers/callback")
  }
  
  authClaimQuickConnectCode() {
    return createUrl(this.baseUrl, "/api/v1/auth/quick-connect/claim")
  }
  
  authFinishProvider() {
    return createUrl(this.baseUrl, "/api/v1/auth/providers/finish")
  }
  
  authFinishQuickConnect() {
    return createUrl(this.baseUrl, "/api/v1/auth/quick-connect/finish")
  }
  
  authGetProviderStatus() {
    return createUrl(this.baseUrl, "/api/v1/auth/provider/status")
  }
  
  authGetProviders() {
    return createUrl(this.baseUrl, "/api/v1/auth/providers")
  }
  
  authGetQuickConnectStatus() {
    return createUrl(this.baseUrl, "/api/v1/auth/quick-connect/status")
  }
  
  authProviderInitiate() {
    return createUrl(this.baseUrl, "/api/v1/auth/providers/initiate")
  }
  
  authQuickConnectInitiate() {
    return createUrl(this.baseUrl, "/api/v1/auth/quick-connect/initiate")
  }
  
  clearQueue() {
    return createUrl(this.baseUrl, "/api/v1/queue")
  }
  
  createApiToken() {
    return createUrl(this.baseUrl, "/api/v1/me/apitokens")
  }
  
  createPlaylist() {
    return createUrl(this.baseUrl, "/api/v1/playlists")
  }
  
  createTrackFilter() {
    return createUrl(this.baseUrl, "/api/v1/me/filters/tracks")
  }
  
  deleteApiToken(tokenId: string) {
    return createUrl(this.baseUrl, `/api/v1/me/apitokens/${tokenId}`)
  }
  
  deletePlaylist(playlistId: string) {
    return createUrl(this.baseUrl, `/api/v1/playlists/${playlistId}`)
  }
  
  deleteTrackFilter(filterId: string) {
    return createUrl(this.baseUrl, `/api/v1/me/filters/tracks/${filterId}`)
  }
  
  editPlaylist(playlistId: string) {
    return createUrl(this.baseUrl, `/api/v1/playlists/${playlistId}`)
  }
  
  favoriteTrack(trackId: string) {
    return createUrl(this.baseUrl, `/api/v1/me/favorites/tracks/${trackId}`)
  }
  
  generatePlaylistImage(playlistId: string) {
    return createUrl(this.baseUrl, `/api/v1/playlists/${playlistId}/images/generate`)
  }
  
  getAlbumById(id: string) {
    return createUrl(this.baseUrl, `/api/v1/albums/${id}`)
  }
  
  getAlbumImage(albumId: string, image: string) {
    return createUrl(this.baseUrl, `/files/albums/images/${albumId}/${image}`)
  }
  
  getAlbumTracks(id: string) {
    return createUrl(this.baseUrl, `/api/v1/albums/${id}/tracks`)
  }
  
  getAlbums() {
    return createUrl(this.baseUrl, "/api/v1/albums")
  }
  
  getApiTokens() {
    return createUrl(this.baseUrl, "/api/v1/me/apitokens")
  }
  
  getArtistAlbums(id: string) {
    return createUrl(this.baseUrl, `/api/v1/artists/${id}/albums`)
  }
  
  getArtistById(id: string) {
    return createUrl(this.baseUrl, `/api/v1/artists/${id}`)
  }
  
  getArtistImage(artistId: string, image: string) {
    return createUrl(this.baseUrl, `/files/artists/images/${artistId}/${image}`)
  }
  
  getArtists() {
    return createUrl(this.baseUrl, "/api/v1/artists")
  }
  
  getFavoriteTrackIds() {
    return createUrl(this.baseUrl, "/api/v1/me/favorites/tracks/ids")
  }
  
  getMe() {
    return createUrl(this.baseUrl, "/api/v1/auth/me")
  }
  
  getMediaSettings() {
    return createUrl(this.baseUrl, "/api/v1/media/settings")
  }
  
  getPlaylistById(playlistId: string) {
    return createUrl(this.baseUrl, `/api/v1/playlists/${playlistId}`)
  }
  
  getPlaylistImage(playlistId: string, image: string) {
    return createUrl(this.baseUrl, `/files/playlists/images/${playlistId}/${image}`)
  }
  
  getPlaylistItemIds(playlistId: string) {
    return createUrl(this.baseUrl, `/api/v1/playlists/${playlistId}/ids`)
  }
  
  getPlaylistItems(playlistId: string) {
    return createUrl(this.baseUrl, `/api/v1/playlists/${playlistId}/items`)
  }
  
  getPlaylists() {
    return createUrl(this.baseUrl, "/api/v1/playlists")
  }
  
  getQueue() {
    return createUrl(this.baseUrl, "/api/v1/queue")
  }
  
  getQueueIds() {
    return createUrl(this.baseUrl, "/api/v1/queue/ids")
  }
  
  getQueueItemAtIndex(position: string) {
    return createUrl(this.baseUrl, `/api/v1/queue/items/${position}`)
  }
  
  getQuickPlaylistIds() {
    return createUrl(this.baseUrl, "/api/v1/me/quickplaylist")
  }
  
  getSystemInfo() {
    return createUrl(this.baseUrl, "/api/v1/system/info")
  }
  
  getTrackById(id: string) {
    return createUrl(this.baseUrl, `/api/v1/tracks/${id}`)
  }
  
  getTrackHistory() {
    return createUrl(this.baseUrl, "/api/v1/history/tracks")
  }
  
  getTrackHistoryById(id: string) {
    return createUrl(this.baseUrl, `/api/v1/history/tracks/${id}`)
  }
  
  getTracks() {
    return createUrl(this.baseUrl, "/api/v1/tracks")
  }
  
  getUser(userId: string) {
    return createUrl(this.baseUrl, `/api/v1/users/${userId}`)
  }
  
  getUserImage(userId: string, image: string) {
    return createUrl(this.baseUrl, `/files/users/images/${userId}/${image}`)
  }
  
  getUserStats(userId: string) {
    return createUrl(this.baseUrl, `/api/v1/users/${userId}/stats`)
  }
  
  getUserTrackFavorites(userId: string) {
    return createUrl(this.baseUrl, `/api/v1/users/${userId}/favorites/tracks`)
  }
  
  getUserTrackFilters(userId: string) {
    return createUrl(this.baseUrl, `/api/v1/users/${userId}/filters/tracks`)
  }
  
  pushTrackHistory() {
    return createUrl(this.baseUrl, "/api/v1/history/tracks")
  }
  
  removePlaylistItem(playlistId: string) {
    return createUrl(this.baseUrl, `/api/v1/playlists/${playlistId}/items`)
  }
  
  removeQueueItem(itemId: string) {
    return createUrl(this.baseUrl, `/api/v1/queue/items/${itemId}`)
  }
  
  reorderPlaylistItems(playlistId: string) {
    return createUrl(this.baseUrl, `/api/v1/playlists/${playlistId}/items/reorder`)
  }
  
  replaceQueue() {
    return createUrl(this.baseUrl, "/api/v1/queue")
  }
  
  runTask(taskName: string) {
    return createUrl(this.baseUrl, `/api/v1/system/task/${taskName}`)
  }
  
  searchAlbums() {
    return createUrl(this.baseUrl, "/api/v1/search/albums")
  }
  
  searchArtists() {
    return createUrl(this.baseUrl, "/api/v1/search/artists")
  }
  
  searchPlaylists() {
    return createUrl(this.baseUrl, "/api/v1/search/playlists")
  }
  
  searchTracks() {
    return createUrl(this.baseUrl, "/api/v1/search/tracks")
  }
  
  searchUsers() {
    return createUrl(this.baseUrl, "/api/v1/search/users")
  }
  
  setQueuePosition() {
    return createUrl(this.baseUrl, "/api/v1/queue/position")
  }
  
  setQuickPlaylist() {
    return createUrl(this.baseUrl, "/api/v1/me/quickplaylist")
  }
  
  sseHandler() {
    return createUrl(this.baseUrl, "/api/v1/system/sse")
  }
  
  streamTrack(trackId: string) {
    return createUrl(this.baseUrl, `/api/v1/media/stream/tracks/${trackId}`)
  }
  
  unfavoriteTrack(trackId: string) {
    return createUrl(this.baseUrl, `/api/v1/me/favorites/tracks/${trackId}`)
  }
  
  updateMe() {
    return createUrl(this.baseUrl, "/api/v1/me")
  }
  
  updateTrackFilter(filterId: string) {
    return createUrl(this.baseUrl, `/api/v1/me/filters/tracks/${filterId}`)
  }
  
  uploadPlaylistImage(playlistId: string) {
    return createUrl(this.baseUrl, `/api/v1/playlists/${playlistId}/image/upload`)
  }
}
