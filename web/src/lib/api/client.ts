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
  
  addTrackEvent(trackId: string, body: api.AddTrackEventBody, options?: ExtraOptions) {
    return this.request(`/api/v1/media/event/track/${trackId}`, "POST", z.undefined(), z.any(), body, options)
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
    return this.request("/api/v1/me/apitokens", "POST", api.CreateApiToken, z.any(), body, options)
  }
  
  createPlaylist(body: api.CreatePlaylistBody, options?: ExtraOptions) {
    return this.request("/api/v1/playlists", "POST", api.CreatePlaylist, z.any(), body, options)
  }
  
  createPlaylistFilter(playlistId: string, body: api.AddPlaylistFilterBody, options?: ExtraOptions) {
    return this.request(`/api/v1/playlists/${playlistId}/filters`, "POST", api.AddPlaylistFilter, z.any(), body, options)
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
  
  editPlaylistFilter(playlistId: string, filterId: string, body: api.EditPlaylistFilterBody, options?: ExtraOptions) {
    return this.request(`/api/v1/playlists/${playlistId}/filters/${filterId}`, "PATCH", z.undefined(), z.any(), body, options)
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
  
  getQuickPlaylistIds(options?: ExtraOptions) {
    return this.request("/api/v1/me/quickplaylist", "GET", api.GetQuickPlaylistIds, z.any(), undefined, options)
  }
  
  getSystemInfo(options?: ExtraOptions) {
    return this.request("/api/v1/system/info", "GET", api.GetSystemInfo, z.any(), undefined, options)
  }
  
  getTrackById(id: string, options?: ExtraOptions) {
    return this.request(`/api/v1/tracks/${id}`, "GET", api.GetTrackById, z.any(), undefined, options)
  }
  
  getTracks(options?: ExtraOptions) {
    return this.request("/api/v1/tracks", "GET", api.GetTracks, z.any(), undefined, options)
  }
  
  getUser(userId: string, options?: ExtraOptions) {
    return this.request(`/api/v1/users/${userId}`, "GET", api.GetUser, z.any(), undefined, options)
  }
  
  
  getUserTrackFavorites(userId: string, options?: ExtraOptions) {
    return this.request(`/api/v1/users/${userId}/favorites/tracks`, "GET", api.GetUserFavorites, z.any(), undefined, options)
  }
  
  getUserTrackFilters(userId: string, options?: ExtraOptions) {
    return this.request(`/api/v1/users/${userId}/filters/tracks`, "GET", api.GetTrackFilters, z.any(), undefined, options)
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
  
  addTrackEvent(trackId: string) {
    return createUrl(this.baseUrl, `/api/v1/media/event/track/${trackId}`)
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
    return createUrl(this.baseUrl, "/api/v1/me/apitokens")
  }
  
  createPlaylist() {
    return createUrl(this.baseUrl, "/api/v1/playlists")
  }
  
  createPlaylistFilter(playlistId: string) {
    return createUrl(this.baseUrl, `/api/v1/playlists/${playlistId}/filters`)
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
  
  editPlaylistFilter(playlistId: string, filterId: string) {
    return createUrl(this.baseUrl, `/api/v1/playlists/${playlistId}/filters/${filterId}`)
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
  
  getQuickPlaylistIds() {
    return createUrl(this.baseUrl, "/api/v1/me/quickplaylist")
  }
  
  getSystemInfo() {
    return createUrl(this.baseUrl, "/api/v1/system/info")
  }
  
  getTrackById(id: string) {
    return createUrl(this.baseUrl, `/api/v1/tracks/${id}`)
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
  
  getUserTrackFavorites(userId: string) {
    return createUrl(this.baseUrl, `/api/v1/users/${userId}/favorites/tracks`)
  }
  
  getUserTrackFilters(userId: string) {
    return createUrl(this.baseUrl, `/api/v1/users/${userId}/filters/tracks`)
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
