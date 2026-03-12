import { z } from "zod";
import * as api from "./types";
import { BaseApiClient, createUrl, type ExtraOptions } from "./base-client";


export class ApiClient extends BaseApiClient {
  url: ClientUrls;

  constructor(baseUrl: string) {
    super(baseUrl);
    this.url = new ClientUrls(baseUrl);
  }
  
  addItemToPlaylist(id: string, body: api.AddItemToPlaylistBody, options?: ExtraOptions) {
    return this.request(`/api/v1/playlists/${id}/items`, "POST", z.undefined(), z.any(), body, options)
  }
  
  addPlaylistFilter(playlistId: string, body: api.AddPlaylistFilterBody, options?: ExtraOptions) {
    return this.request(`/api/v1/playlists/${playlistId}/filters`, "POST", api.AddPlaylistFilter, z.any(), body, options)
  }
  
  addToUserQuickPlaylist(body: api.TrackId, options?: ExtraOptions) {
    return this.request("/api/v1/user/quickplaylist", "POST", z.undefined(), z.any(), body, options)
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
  
  clearPlaylist(id: string, options?: ExtraOptions) {
    return this.request(`/api/v1/playlists/${id}/items/all`, "DELETE", z.undefined(), z.any(), undefined, options)
  }
  
  createApiToken(body: api.CreateApiTokenBody, options?: ExtraOptions) {
    return this.request("/api/v1/user/apitoken", "POST", api.CreateApiToken, z.any(), body, options)
  }
  
  createPlaylist(body: api.CreatePlaylistBody, options?: ExtraOptions) {
    return this.request("/api/v1/playlists", "POST", api.CreatePlaylist, z.any(), body, options)
  }
  
  createPlaylistFromFilter(body: api.PostPlaylistFilterBody, options?: ExtraOptions) {
    return this.request("/api/v1/playlists/filter", "POST", api.CreatePlaylist, z.any(), body, options)
  }
  
  createVirtualPlaylist(body: api.CreateVirtualPlaylistBody, options?: ExtraOptions) {
    return this.request("/api/v1/virtual-playlists", "POST", api.CreateVirtualPlaylist, z.any(), body, options)
  }
  
  deleteApiToken(id: string, options?: ExtraOptions) {
    return this.request(`/api/v1/user/apitoken/${id}`, "DELETE", z.undefined(), z.any(), undefined, options)
  }
  
  deletePlaylist(id: string, options?: ExtraOptions) {
    return this.request(`/api/v1/playlists/${id}`, "DELETE", z.undefined(), z.any(), undefined, options)
  }
  
  deleteVirtualPlaylist(id: string, options?: ExtraOptions) {
    return this.request(`/api/v1/virtual-playlists/${id}`, "DELETE", z.undefined(), z.any(), undefined, options)
  }
  
  editPlaylist(id: string, body: api.EditPlaylistBody, options?: ExtraOptions) {
    return this.request(`/api/v1/playlists/${id}`, "PATCH", z.undefined(), z.any(), body, options)
  }
  
  generatePlaylistImage(id: string, options?: ExtraOptions) {
    return this.request(`/api/v1/playlists/${id}/images/generate`, "POST", z.undefined(), z.any(), undefined, options)
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
  
  getMediaFromPlaylist(playlistId: string, body: api.GetMediaFromPlaylistBody, options?: ExtraOptions) {
    return this.request(`/api/v1/media/playlist/${playlistId}`, "POST", api.GetMedia, z.any(), body, options)
  }
  
  getPlaylistById(id: string, options?: ExtraOptions) {
    return this.request(`/api/v1/playlists/${id}`, "GET", api.GetPlaylistById, z.any(), undefined, options)
  }
  
  getPlaylistFilters(playlistId: string, options?: ExtraOptions) {
    return this.request(`/api/v1/playlists/${playlistId}/filters`, "GET", api.GetPlaylistFilters, z.any(), undefined, options)
  }
  
  
  getPlaylistItems(id: string, options?: ExtraOptions) {
    return this.request(`/api/v1/playlists/${id}/items`, "GET", api.GetPlaylistItems, z.any(), undefined, options)
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
  
  getTracks(options?: ExtraOptions) {
    return this.request("/api/v1/tracks", "GET", api.GetTracks, z.any(), undefined, options)
  }
  
  getUser(id: string, options?: ExtraOptions) {
    return this.request(`/api/v1/users/${id}`, "GET", api.GetUser, z.any(), undefined, options)
  }
  
  getUserQuickPlaylistItemIds(options?: ExtraOptions) {
    return this.request("/api/v1/user/quickplaylist", "GET", api.GetUserQuickPlaylistItemIds, z.any(), undefined, options)
  }
  
  getVirtualPlaylistById(id: string, options?: ExtraOptions) {
    return this.request(`/api/v1/virtual-playlists/${id}`, "GET", api.GetVirtualPlaylistById, z.any(), undefined, options)
  }
  
  getVirtualPlaylistTracks(id: string, options?: ExtraOptions) {
    return this.request(`/api/v1/virtual-playlists/${id}/tracks`, "GET", api.GetVirtualPlaylistTracks, z.any(), undefined, options)
  }
  
  getVirtualPlaylists(options?: ExtraOptions) {
    return this.request("/api/v1/virtual-playlists", "GET", api.GetVirtualPlaylists, z.any(), undefined, options)
  }
  
  getVirtualPlaylistsForPlaylist(playlistId: string, options?: ExtraOptions) {
    return this.request(`/api/v1/virtual-playlists/playlists/${playlistId}`, "GET", api.GetVirtualPlaylists, z.any(), undefined, options)
  }
  
  removeItemFromUserQuickPlaylist(body: api.TrackId, options?: ExtraOptions) {
    return this.request("/api/v1/user/quickplaylist", "DELETE", z.undefined(), z.any(), body, options)
  }
  
  removePlaylistItem(id: string, body: api.RemovePlaylistItemBody, options?: ExtraOptions) {
    return this.request(`/api/v1/playlists/${id}/items`, "DELETE", z.undefined(), z.any(), body, options)
  }
  
  searchAlbums(options?: ExtraOptions) {
    return this.request("/api/v1/albums/search", "GET", api.SearchAlbums, z.any(), undefined, options)
  }
  
  searchArtists(options?: ExtraOptions) {
    return this.request("/api/v1/artists/search", "GET", api.SearchArtists, z.any(), undefined, options)
  }
  
  searchTracks(options?: ExtraOptions) {
    return this.request("/api/v1/tracks/search", "GET", api.SearchTracks, z.any(), undefined, options)
  }
  
  
  
  syncLibrary(options?: ExtraOptions) {
    return this.request("/api/v1/system/library", "POST", z.undefined(), z.any(), undefined, options)
  }
  
  updateUserSettings(body: api.UpdateUserSettingsBody, options?: ExtraOptions) {
    return this.request("/api/v1/user/settings", "PATCH", z.undefined(), z.any(), body, options)
  }
  
  updateVirtualPlaylist(id: string, body: api.UpdateVirtualPlaylistBody, options?: ExtraOptions) {
    return this.request(`/api/v1/virtual-playlists/${id}`, "PATCH", z.undefined(), z.any(), body, options)
  }
}

export class ClientUrls {
  baseUrl: string;

  constructor(baseUrl: string) {
    this.baseUrl = baseUrl;
  }
  
  addItemToPlaylist(id: string) {
    return createUrl(this.baseUrl, `/api/v1/playlists/${id}/items`)
  }
  
  addPlaylistFilter(playlistId: string) {
    return createUrl(this.baseUrl, `/api/v1/playlists/${playlistId}/filters`)
  }
  
  addToUserQuickPlaylist() {
    return createUrl(this.baseUrl, "/api/v1/user/quickplaylist")
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
  
  clearPlaylist(id: string) {
    return createUrl(this.baseUrl, `/api/v1/playlists/${id}/items/all`)
  }
  
  createApiToken() {
    return createUrl(this.baseUrl, "/api/v1/user/apitoken")
  }
  
  createPlaylist() {
    return createUrl(this.baseUrl, "/api/v1/playlists")
  }
  
  createPlaylistFromFilter() {
    return createUrl(this.baseUrl, "/api/v1/playlists/filter")
  }
  
  createVirtualPlaylist() {
    return createUrl(this.baseUrl, "/api/v1/virtual-playlists")
  }
  
  deleteApiToken(id: string) {
    return createUrl(this.baseUrl, `/api/v1/user/apitoken/${id}`)
  }
  
  deletePlaylist(id: string) {
    return createUrl(this.baseUrl, `/api/v1/playlists/${id}`)
  }
  
  deleteVirtualPlaylist(id: string) {
    return createUrl(this.baseUrl, `/api/v1/virtual-playlists/${id}`)
  }
  
  editPlaylist(id: string) {
    return createUrl(this.baseUrl, `/api/v1/playlists/${id}`)
  }
  
  generatePlaylistImage(id: string) {
    return createUrl(this.baseUrl, `/api/v1/playlists/${id}/images/generate`)
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
  
  getMediaFromPlaylist(playlistId: string) {
    return createUrl(this.baseUrl, `/api/v1/media/playlist/${playlistId}`)
  }
  
  getPlaylistById(id: string) {
    return createUrl(this.baseUrl, `/api/v1/playlists/${id}`)
  }
  
  getPlaylistFilters(playlistId: string) {
    return createUrl(this.baseUrl, `/api/v1/playlists/${playlistId}/filters`)
  }
  
  getPlaylistImage(playlistId: string, image: string) {
    return createUrl(this.baseUrl, `/files/playlists/images/${playlistId}/${image}`)
  }
  
  getPlaylistItems(id: string) {
    return createUrl(this.baseUrl, `/api/v1/playlists/${id}/items`)
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
  
  getTracks() {
    return createUrl(this.baseUrl, "/api/v1/tracks")
  }
  
  getUser(id: string) {
    return createUrl(this.baseUrl, `/api/v1/users/${id}`)
  }
  
  getUserQuickPlaylistItemIds() {
    return createUrl(this.baseUrl, "/api/v1/user/quickplaylist")
  }
  
  getVirtualPlaylistById(id: string) {
    return createUrl(this.baseUrl, `/api/v1/virtual-playlists/${id}`)
  }
  
  getVirtualPlaylistTracks(id: string) {
    return createUrl(this.baseUrl, `/api/v1/virtual-playlists/${id}/tracks`)
  }
  
  getVirtualPlaylists() {
    return createUrl(this.baseUrl, "/api/v1/virtual-playlists")
  }
  
  getVirtualPlaylistsForPlaylist(playlistId: string) {
    return createUrl(this.baseUrl, `/api/v1/virtual-playlists/playlists/${playlistId}`)
  }
  
  removeItemFromUserQuickPlaylist() {
    return createUrl(this.baseUrl, "/api/v1/user/quickplaylist")
  }
  
  removePlaylistItem(id: string) {
    return createUrl(this.baseUrl, `/api/v1/playlists/${id}/items`)
  }
  
  searchAlbums() {
    return createUrl(this.baseUrl, "/api/v1/albums/search")
  }
  
  searchArtists() {
    return createUrl(this.baseUrl, "/api/v1/artists/search")
  }
  
  searchTracks() {
    return createUrl(this.baseUrl, "/api/v1/tracks/search")
  }
  
  sseHandler() {
    return createUrl(this.baseUrl, "/api/v1/system/sse")
  }
  
  streamTrack(trackId: string) {
    return createUrl(this.baseUrl, `/media/tracks/${trackId}/stream`)
  }
  
  syncLibrary() {
    return createUrl(this.baseUrl, "/api/v1/system/library")
  }
  
  updateUserSettings() {
    return createUrl(this.baseUrl, "/api/v1/user/settings")
  }
  
  updateVirtualPlaylist(id: string) {
    return createUrl(this.baseUrl, `/api/v1/virtual-playlists/${id}`)
  }
}
