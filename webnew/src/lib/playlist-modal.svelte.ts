import { handleApiError } from "$lib";
import type { ApiClient } from "$lib/api/client";
import type { Playlist } from "$lib/api/types";

export type ShowPlaylistModalOptions = {
	selectedId?: string;
	title?: string;
	description?: string;
};

function sanitizeFilterValue(value: string): string {
	return value.replace(/'/g, "\\'");
}

class PlaylistModalManager {
	apiClient: ApiClient;

	open = $state(false);
	playlists = $state<Playlist[]>([]);
	selectedId: string | null = $state(null);
	searchQuery = $state("");
	title = $state("Select Playlist");
	description = $state("Choose a playlist to add tracks to");
	// eslint-disable-next-line no-unused-vars
	resolve: ((playlist: Playlist | null) => void) | null = null;
	promise: Promise<Playlist | null> | null = null;
	private searchTimeout: ReturnType<typeof setTimeout> | null = null;

	constructor(apiClient: ApiClient) {
		this.apiClient = apiClient;
	}

	async show(options?: ShowPlaylistModalOptions): Promise<Playlist | null> {
		if (this.promise) {
			return this.promise;
		}

		this.promise = this._show(options);
		const result = await this.promise;
		this.promise = null;
		return result;
	}

	private async _show(
		options?: ShowPlaylistModalOptions,
	): Promise<Playlist | null> {
		const res = await this.apiClient.getPlaylists();
		if (!res.success) {
			handleApiError(res.error);
			return null;
		}

		this.playlists = res.data.playlists;
		this.selectedId = options?.selectedId ?? null;
		this.searchQuery = "";
		this.title = options?.title ?? "Select Playlist";
		this.description =
			options?.description ?? "Choose a playlist to add tracks to";
		this.open = true;

		return new Promise((resolve) => {
			this.resolve = resolve;
		});
	}

	async search(query: string) {
		this.searchQuery = query;

		if (this.searchTimeout) {
			clearTimeout(this.searchTimeout);
		}

		this.searchTimeout = setTimeout(async () => {
			const options: Parameters<typeof this.apiClient.getPlaylists>[0] = {};
			if (query.trim()) {
				options.query = {
					filter: `name contains "%${sanitizeFilterValue(query.trim())}%"`,
				};
			}

			const res = await this.apiClient.getPlaylists(options);
			if (res.success) {
				this.playlists = res.data.playlists;
			}
		}, 200);
	}

	async enter() {
		if (this.playlists.length > 0) {
			const playlist = this.playlists[0]
			this.select(playlist)
		}
	}

	select(playlist: Playlist) {
		const resolve = this.resolve;
		this.resolve = null;
		this.open = false;
		resolve?.(playlist);
	}
}

let manager: PlaylistModalManager | null = null;

export function initPlaylistModalManager(apiClient: ApiClient) {
	manager = new PlaylistModalManager(apiClient);
}

export function getManager(): PlaylistModalManager {
	if (!manager) {
		throw new Error("PlaylistModalManager not initialized");
	}
	return manager;
}

export async function showPlaylistModal(
	options?: ShowPlaylistModalOptions,
): Promise<Playlist | null> {
	return getManager().show(options);
}
