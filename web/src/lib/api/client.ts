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
  
  createApiToken(body: api.CreateApiTokenBody, options?: ExtraOptions) {
    return this.request("/api/v1/user/apitoken", "POST", api.CreateApiToken, z.any(), body, options)
  }
  
  createPlaylist(body: api.CreatePlaylistBody, options?: ExtraOptions) {
    return this.request("/api/v1/playlists", "POST", api.CreatePlaylist, z.any(), body, options)
  }
  
  createPlaylistFilter(playlistId: string, body: api.AddPlaylistFilterBody, options?: ExtraOptions) {
    return this.request(`/api/v1/playlists/${playlistId}/filters`, "POST", api.AddPlaylistFilter, z.any(), body, options)
  }
  
  createTrackFilter(body: api.CreateTrackFilterBody, options?: ExtraOptions) {
    return this.request("/api/v1/user/tracks/filter", "POST", api.CreateTrackFilter, z.any(), body, options)
  }
  
  deleteApiToken(id: string, options?: ExtraOptions) {
    return this.request(`/api/v1/user/apitoken/${id}`, "DELETE", z.undefined(), z.any(), undefined, options)
  }
  
  deletePlaylist(playlistId: string, options?: ExtraOptions) {
    return this.request(`/api/v1/playlists/${playlistId}`, "DELETE", z.undefined(), z.any(), undefined, options)
  }
  
  deleteTrackFilter(filterId: string, options?: ExtraOptions) {
    return this.request(`/api/v1/user/tracks/filter/${filterId}`, "DELETE", z.undefined(), z.any(), undefined, options)
  }
  
  editPlaylist(playlistId: string, body: api.EditPlaylistBody, options?: ExtraOptions) {
    return this.request(`/api/v1/playlists/${playlistId}`, "PATCH", z.undefined(), z.any(), body, options)
  }
  
  editPlaylistFilter(playlistId: string, filterId: string, body: api.EditPlaylistFilterBody, options?: ExtraOptions) {
    return this.request(`/api/v1/playlists/${playlistId}/filters/${filterId}`, "PATCH", z.undefined(), z.any(), body, options)
  }
  
  editTrackFilter(filterId: string, body: api.EditTrackFilterBody, options?: ExtraOptions) {
    return this.request(`/api/v1/user/tracks/filter/${filterId}`, "PATCH", z.undefined(), z.any(), body, options)
  }
  
  editUser(body: api.EditUserBody, options?: ExtraOptions) {
    return this.request("/api/v1/user", "PATCH", z.undefined(), z.any(), body, options)
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
  
  getAllApiTokens(options?: ExtraOptions) {
    return this.request("/api/v1/user/apitoken", "GET", api.GetAllApiTokens, z.any(), undefined, options)
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
  
  getMe(options?: ExtraOptions) {
    return this.request("/api/v1/auth/me", "GET", api.GetMe, z.any(), undefined, options)
  }
  
  getMediaFromAlbum(albumId: string, body: api.GetMediaFromAlbumBody, options?: ExtraOptions) {
    return this.request(`/api/v1/media/album/${albumId}`, "POST", api.GetMedia, z.any(), body, options)
  }
  
  getMediaFromArtist(artistId: string, body: api.GetMediaFromArtistBody, options?: ExtraOptions) {
    return this.request(`/api/v1/media/artist/${artistId}`, "POST", api.GetMedia, z.any(), body, options)
  }
  
  getMediaFromFilter(body: api.GetMediaFromFilterBody, options?: ExtraOptions) {
    return this.request("/api/v1/media/filter", "POST", api.GetMedia, z.any(), body, options)
  }
  
  getMediaFromIds(body: api.GetMediaFromIdsBody, options?: ExtraOptions) {
    return this.request("/api/v1/media/ids", "POST", api.GetMedia, z.any(), body, options)
  }
  
  getMediaSettings(options?: ExtraOptions) {
    return this.request("/api/v1/media/settings", "GET", api.GetMediaSettings, z.any(), undefined, options)
  }
  
  getPlaylistById(playlistId: string, options?: ExtraOptions) {
    return this.request(`/api/v1/playlists/${playlistId}`, "GET", api.GetPlaylistById, z.any(), undefined, options)
  }
  
  getPlaylistFilters(playlistId: string, options?: ExtraOptions) {
    return this.request(`/api/v1/playlists/${playlistId}/filters`, "GET", api.GetPlaylistFilters, z.any(), undefined, options)
  }
  
  
  getPlaylistItems(playlistId: string, options?: ExtraOptions) {
    return this.request(`/api/v1/playlists/${playlistId}/items`, "GET", api.GetPlaylistItems, z.any(), undefined, options)
  }
  
  getPlaylists(options?: ExtraOptions) {
    return this.request("/api/v1/playlists", "GET", api.GetPlaylists, z.any(), undefined, options)
  }
  
  getSystemInfo(options?: ExtraOptions) {
    return this.request("/api/v1/system/info", "GET", api.GetSystemInfo, z.any(), undefined, options)
  }
  
  getTrackById(id: string, options?: ExtraOptions) {
    return this.request(`/api/v1/tracks/${id}`, "GET", api.GetTrackById, z.any(), undefined, options)
  }
  
  getTrackFilters(userId: string, options?: ExtraOptions) {
    return this.request(`/api/v1/user/${userId}/tracks/filter`, "GET", api.GetTrackFilters, z.any(), undefined, options)
  }
  
  getTracks(options?: ExtraOptions) {
    return this.request("/api/v1/tracks", "GET", api.GetTracks, z.any(), undefined, options)
  }
  
  getUser(id: string, options?: ExtraOptions) {
    return this.request(`/api/v1/users/${id}`, "GET", api.GetUser, z.any(), undefined, options)
  }
  
  
  getUserQuickPlaylistItemIds(options?: ExtraOptions) {
    return this.request("/api/v1/user/quickplaylist", "GET", api.GetUserQuickPlaylistItemIds, z.any(), undefined, options)
  }
  
  recordTrack(trackId: string, body: api.RecordTrackBody, options?: ExtraOptions) {
    return this.request(`/api/v1/media/record/track/${trackId}`, "POST", z.undefined(), z.any(), body, options)
  }
  
  removePlaylistItem(playlistId: string, body: api.RemovePlaylistItemBody, options?: ExtraOptions) {
    return this.request(`/api/v1/playlists/${playlistId}/items`, "DELETE", z.undefined(), z.any(), body, options)
  }
  
  reorderPlaylistItems(playlistId: string, body: api.ReorderPlaylistItemsBody, options?: ExtraOptions) {
    return this.request(`/api/v1/playlists/${playlistId}/items/reorder`, "POST", z.undefined(), z.any(), body, options)
  }
  
  runJob(jobName: string, options?: ExtraOptions) {
    return this.request(`/api/v1/system/job/${jobName}`, "POST", z.undefined(), z.any(), undefined, options)
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
  
  
  
  updateUserSettings(body: api.UpdateUserSettingsBody, options?: ExtraOptions) {
    return this.request("/api/v1/user/settings", "PATCH", z.undefined(), z.any(), body, options)
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
  
  createApiToken() {
    return createUrl(this.baseUrl, "/api/v1/user/apitoken")
  }
  
  createPlaylist() {
    return createUrl(this.baseUrl, "/api/v1/playlists")
  }
  
  createPlaylistFilter(playlistId: string) {
    return createUrl(this.baseUrl, `/api/v1/playlists/${playlistId}/filters`)
  }
  
  createTrackFilter() {
    return createUrl(this.baseUrl, "/api/v1/user/tracks/filter")
  }
  
  deleteApiToken(id: string) {
    return createUrl(this.baseUrl, `/api/v1/user/apitoken/${id}`)
  }
  
  deletePlaylist(playlistId: string) {
    return createUrl(this.baseUrl, `/api/v1/playlists/${playlistId}`)
  }
  
  deleteTrackFilter(filterId: string) {
    return createUrl(this.baseUrl, `/api/v1/user/tracks/filter/${filterId}`)
  }
  
  editPlaylist(playlistId: string) {
    return createUrl(this.baseUrl, `/api/v1/playlists/${playlistId}`)
  }
  
  editPlaylistFilter(playlistId: string, filterId: string) {
    return createUrl(this.baseUrl, `/api/v1/playlists/${playlistId}/filters/${filterId}`)
  }
  
  editTrackFilter(filterId: string) {
    return createUrl(this.baseUrl, `/api/v1/user/tracks/filter/${filterId}`)
  }
  
  editUser() {
    return createUrl(this.baseUrl, "/api/v1/user")
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
  
  getAllApiTokens() {
    return createUrl(this.baseUrl, "/api/v1/user/apitoken")
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
  
  getMe() {
    return createUrl(this.baseUrl, "/api/v1/auth/me")
  }
  
  getMediaFromAlbum(albumId: string) {
    return createUrl(this.baseUrl, `/api/v1/media/album/${albumId}`)
  }
  
  getMediaFromArtist(artistId: string) {
    return createUrl(this.baseUrl, `/api/v1/media/artist/${artistId}`)
  }
  
  getMediaFromFilter() {
    return createUrl(this.baseUrl, "/api/v1/media/filter")
  }
  
  getMediaFromIds() {
    return createUrl(this.baseUrl, "/api/v1/media/ids")
  }
  
  getMediaSettings() {
    return createUrl(this.baseUrl, "/api/v1/media/settings")
  }
  
  getPlaylistById(playlistId: string) {
    return createUrl(this.baseUrl, `/api/v1/playlists/${playlistId}`)
  }
  
  getPlaylistFilters(playlistId: string) {
    return createUrl(this.baseUrl, `/api/v1/playlists/${playlistId}/filters`)
  }
  
  getPlaylistImage(playlistId: string, image: string) {
    return createUrl(this.baseUrl, `/files/playlists/images/${playlistId}/${image}`)
  }
  
  getPlaylistItems(playlistId: string) {
    return createUrl(this.baseUrl, `/api/v1/playlists/${playlistId}/items`)
  }
  
  getPlaylists() {
    return createUrl(this.baseUrl, "/api/v1/playlists")
  }
  
  getSystemInfo() {
    return createUrl(this.baseUrl, "/api/v1/system/info")
  }
  
  getTrackById(id: string) {
    return createUrl(this.baseUrl, `/api/v1/tracks/${id}`)
  }
  
  getTrackFilters(userId: string) {
    return createUrl(this.baseUrl, `/api/v1/user/${userId}/tracks/filter`)
  }
  
  getTracks() {
    return createUrl(this.baseUrl, "/api/v1/tracks")
  }
  
  getUser(id: string) {
    return createUrl(this.baseUrl, `/api/v1/users/${id}`)
  }
  
  getUserImage(userId: string, image: string) {
    return createUrl(this.baseUrl, `/files/users/images/${userId}/${image}`)
  }
  
  getUserQuickPlaylistItemIds() {
    return createUrl(this.baseUrl, "/api/v1/user/quickplaylist")
  }
  
  recordTrack(trackId: string) {
    return createUrl(this.baseUrl, `/api/v1/media/record/track/${trackId}`)
  }
  
  removePlaylistItem(playlistId: string) {
    return createUrl(this.baseUrl, `/api/v1/playlists/${playlistId}/items`)
  }
  
  reorderPlaylistItems(playlistId: string) {
    return createUrl(this.baseUrl, `/api/v1/playlists/${playlistId}/items/reorder`)
  }
  
  runJob(jobName: string) {
    return createUrl(this.baseUrl, `/api/v1/system/job/${jobName}`)
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
  
  sseHandler() {
    return createUrl(this.baseUrl, "/api/v1/system/sse")
  }
  
  streamTrack(trackId: string) {
    return createUrl(this.baseUrl, `/api/v1/media/stream/tracks/${trackId}`)
  }
  
  updateUserSettings() {
    return createUrl(this.baseUrl, "/api/v1/user/settings")
  }
  
  uploadPlaylistImage(playlistId: string) {
    return createUrl(this.baseUrl, `/api/v1/playlists/${playlistId}/image/upload`)
  }
}
