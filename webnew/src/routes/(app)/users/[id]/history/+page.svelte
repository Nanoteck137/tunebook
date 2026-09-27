<script lang="ts">
	import { goto, invalidateAll } from "$app/navigation";
	import { page } from "$app/state";
	import { onMount } from "svelte";
	import {
		Check,
		DiscAlbum,
		EllipsisVertical,
		Heart,
		History,
		ListFilter,
		ListPlus,
		Play,
		Star,
		User,
	} from "@lucide/svelte";
	import Image from "$lib/components/Image.svelte";
	import { DropdownMenu, buttonVariants } from "$lib/components/ui";
	import { getFavorites } from "$lib/favorites.svelte";
	import { getMusicManager } from "$lib/music-manager.svelte";
	import { getQuickPlaylist } from "$lib/quick-playlist.svelte";
	import { showPlaylistModal } from "$lib/playlist-modal.svelte";
	import InfiniteScroll from "$lib/components/InfiniteScroll.svelte";
	import { InfiniteScrollController } from "$lib/infinite-scroll.svelte";
	import { getApiClient, handleApiError } from "$lib";
	import type { TrackHistory } from "$lib/api/types";
	import SectionHeader from "$lib/components/SectionHeader.svelte";
	import Spacer from "$lib/components/Spacer.svelte";
	import DebouncedSearchInput from "$lib/components/DebouncedSearchInput.svelte";
	import { SortToggleDropdown } from "$lib/components/sort";
	import { toast } from "svelte-sonner";
	import {
		sortTypes,
		statusTypes,
		defaultSort,
		defaultStatus,
		type SortType,
		type StatusType,
		constructFilterSort,
	} from "./types";

	let { data } = $props();

	const musicManager = getMusicManager();
	const favoritesManager = getFavorites();
	const quickPlaylistManager = getQuickPlaylist();
	const apiClient = getApiClient();

	const scroll = new InfiniteScrollController<TrackHistory>({
		initialLoad: () => ({
			items: data.history,
			hasMore: data.page.page + 1 < data.page.totalPages,
			page: data.page.page,
		}),
		load: async (nextPage) => {
			const query: Record<string, string> = {
				page: String(nextPage),
				perPage: String(data.page.perPage),
			};

			constructFilterSort(data.filter, query);

			const res = await apiClient.getTrackHistory({ query });
			if (!res.success) {
				handleApiError(res.error);
				return null;
			}

			return {
				items: res.data.history,
				hasMore: res.data.page.page + 1 < res.data.page.totalPages,
			};
		},
		itemKey: (entry) => entry.id,
	});

	let value = $state("");

	onMount(() => {
		value = page.url.searchParams.get("query") ?? "";
	});

	async function search(query: string) {
		const params = page.url.searchParams;
		params.delete("query");

		if (query) {
			params.set("query", query);
		}

		await goto("?" + params.toString(), {
			invalidateAll: true,
			keepFocus: true,
			replaceState: true,
		});
	}

	let sort = $state(
		(page.url.searchParams.get("sort") as SortType) ?? defaultSort,
	);

	function updateSort(value: string) {
		sort = value as SortType;

		const query = page.url.searchParams;
		query.delete("sort");

		if (sort !== defaultSort) {
			query.set("sort", sort);
		}

		goto("?" + query.toString(), { invalidateAll: true });
	}

	let status = $state(
		(page.url.searchParams.get("status") as StatusType) ?? defaultStatus,
	);

	function updateStatus(value: StatusType) {
		status = value;

		const query = page.url.searchParams;
		query.delete("status");

		if (status !== defaultStatus) {
			query.set("status", status);
		}

		goto("?" + query.toString(), { invalidateAll: true });
	}

	let hasActiveFilter = $derived(
		data.filter.query !== "" || data.filter.status !== defaultStatus,
	);

	function formatRelativeTime(millis: number): string {
		const unixSeconds = Math.floor(millis / 1000);
		const now = Math.floor(Date.now() / 1000);
		const diff = now - unixSeconds;

		if (diff < 60) return "just now";
		if (diff < 3600) {
			const m = Math.floor(diff / 60);
			return `${m}m ago`;
		}
		if (diff < 86400) {
			const h = Math.floor(diff / 3600);
			return `${h}h ago`;
		}
		if (diff < 604800) {
			const d = Math.floor(diff / 86400);
			return `${d}d ago`;
		}
		return new Date(unixSeconds * 1000).toLocaleDateString(undefined, {
			month: "short",
			day: "numeric",
		});
	}

	function statusLabel(status: string) {
		return status === "completed" ? "Completed" : "Skipped";
	}

	function statusClass(status: string) {
		return status === "completed"
			? "bg-green-500/10 text-green-500 ring-green-500/25"
			: "bg-muted text-muted-foreground ring-border";
	}

	function progressBg(pct: number): string {
		if (pct >= 80) return "bg-green-500/20";
		if (pct >= 40) return "bg-yellow-500/20";
		return "bg-primary/15";
	}

	function playTrack(trackId: string) {
		const trackIds = scroll.items.map((e) => e.track.id);
		musicManager.addTracks({ trackIds, trackId });
	}

	async function saveToPlaylist(trackId: string) {
		const playlist = await showPlaylistModal();
		if (!playlist) return;

		const res = await apiClient.addItemToPlaylist(playlist.id, { trackId });
		if (!res.success) {
			if (res.error.type !== "PLAYLIST_ALREADY_HAS_TRACK") {
				handleApiError(res.error);
			}
			return;
		}

		await invalidateAll();
	}
</script>

{#snippet statusMenu()}
	<DropdownMenu.Content align="end">
		<DropdownMenu.Group>
			{#each statusTypes as type (type.value)}
				<DropdownMenu.Item
					onSelect={() => updateStatus(type.value)}
					class={status === type.value ? "bg-accent text-foreground" : ""}
				>
					<Check
						class={status === type.value ? "text-primary" : "opacity-0"}
					/>
					{type.label}
				</DropdownMenu.Item>
			{/each}
		</DropdownMenu.Group>
	</DropdownMenu.Content>
{/snippet}

<div class="flex flex-col gap-4">
	<SectionHeader count={data.page.totalItems}>
		<History />
		Listening History
	</SectionHeader>

	<div class="flex flex-wrap items-center justify-between gap-2">
		<DebouncedSearchInput
			class="flex-1 md:max-w-64"
			placeholder="Search history..."
			{value}
			setValue={(v) => (value = v)}
			{search}
		/>

		<div class="flex items-center gap-1">
			<DropdownMenu.Root>
				<DropdownMenu.Trigger
					class={buttonVariants({
						variant: "ghost",
						size: "icon",
						class: "relative",
					})}
					title="Status"
					aria-label="Filter by status"
				>
					<ListFilter />
					{#if status !== defaultStatus}
						<span
							class="absolute top-1.5 right-1.5 h-2 w-2 rounded-full bg-primary"
						></span>
					{/if}
				</DropdownMenu.Trigger>
				{@render statusMenu()}
			</DropdownMenu.Root>

			<SortToggleDropdown
				types={sortTypes}
				{sort}
				{defaultSort}
				onSortChange={(value) => updateSort(value)}
			/>
		</div>
	</div>
</div>

<Spacer size="md" />

<InfiniteScroll controller={scroll}>
	{#if scroll.items.length === 0}
		<div class="flex flex-col items-center gap-2 rounded-lg border py-16">
			<History size={32} class="text-muted-foreground/40" />
			<p class="text-sm text-muted-foreground">
				{hasActiveFilter
					? "No plays match the current search or filter"
					: "No listening history yet"}
			</p>
		</div>
	{:else}
		<div class="flex flex-col gap-1.5">
			{#each scroll.items as entry (entry.id)}
				<div
					class="group relative overflow-hidden rounded-lg border bg-card transition-colors hover:bg-accent hover:text-accent-foreground has-data-[state='open']:bg-accent has-data-[state='open']:text-accent-foreground"
				>
					<div
						class="absolute inset-y-0 left-0 transition-[width] duration-300 {progressBg(
							entry.percentPlayed,
						)}"
						style="width: {entry.percentPlayed}%"
					></div>

					<div class="relative z-10 flex items-center gap-3 p-2.5">
						<button
							class="shrink-0 overflow-hidden rounded-md"
							onclick={() => playTrack(entry.track.id)}
							aria-label="Play {entry.track.name}"
						>
							<div class="relative h-12 w-12">
								<Image
									class="h-12 w-12"
									src={entry.track.coverArt.small}
									alt=""
								/>
								<div
									class="absolute inset-0 flex items-center justify-center bg-black/55 opacity-0 transition-opacity group-hover:opacity-100"
								>
									<Play size={18} class="text-white" />
								</div>
							</div>
						</button>

						<div class="flex min-w-0 flex-1 flex-col gap-0.5">
							<div class="flex items-center gap-2">
								<span
									class="truncate text-sm font-medium"
									title={entry.track.name}
								>
									{entry.track.name}
								</span>
								{#if favoritesManager.hasTrack(entry.track.id)}
									<Heart
										size={12}
										class="shrink-0 fill-primary text-primary"
									/>
								{/if}
								{#if quickPlaylistManager.hasTrack(entry.track.id)}
									<Star size={12} class="shrink-0 fill-primary text-primary" />
								{/if}
								<span
									class="shrink-0 rounded-full px-2 py-px text-[10px] font-medium ring-1 {statusClass(
										entry.status,
									)}"
								>
									{statusLabel(entry.status)}
								</span>
							</div>

							<div
								class="flex min-w-0 items-center gap-1.5 text-xs text-muted-foreground"
							>
								<span
									class="truncate"
									title={entry.track.artists.map((a) => a.name).join(", ")}
								>
									{#each entry.track.artists as artist, i (artist.id)}
										{#if i > 0},
										{/if}
										<a class="hover:underline" href="/artists/{artist.id}">
											{artist.name}
										</a>
									{/each}
								</span>
								{#if entry.track.albumName}
									<span class="shrink-0">&middot;</span>
									<span class="truncate">{entry.track.albumName}</span>
								{/if}
							</div>
						</div>

						<div class="hidden shrink-0 flex-col items-end gap-1 sm:flex">
							<span class="text-xs text-muted-foreground tabular-nums">
								{formatRelativeTime(entry.listenedAt)}
							</span>
							<span
								class="w-8 text-right text-[11px] text-muted-foreground tabular-nums"
							>
								{Math.round(entry.percentPlayed)}%
							</span>
						</div>

						<span class="shrink-0 text-[11px] text-muted-foreground sm:hidden">
							{formatRelativeTime(entry.listenedAt)}
						</span>

						<DropdownMenu.Root>
							<DropdownMenu.Trigger
								class={buttonVariants({ variant: "ghost", size: "icon" })}
								title="More options"
								aria-label="More options"
							>
								<EllipsisVertical />
							</DropdownMenu.Trigger>
							<DropdownMenu.Content align="end">
								<DropdownMenu.Group>
									<DropdownMenu.Item
										onSelect={() => playTrack(entry.track.id)}
									>
										<Play />
										Play
									</DropdownMenu.Item>
								</DropdownMenu.Group>

								<DropdownMenu.Separator />

								<DropdownMenu.Item
									onSelect={() => goto(`/albums/${entry.track.albumId}`)}
								>
									<DiscAlbum />
									Go to Album
								</DropdownMenu.Item>
								<DropdownMenu.Sub>
									<DropdownMenu.SubTrigger>
										<User />
										Go to artist
									</DropdownMenu.SubTrigger>
									<DropdownMenu.SubContent>
										{#each entry.track.artists as artist (artist.id)}
											<a
												href="/artists/{artist.id}"
												class="flex items-center gap-2 rounded-sm px-3 py-1.5 text-sm text-popover-foreground hover:bg-accent hover:text-accent-foreground"
											>
												{artist.name}
											</a>
										{/each}
									</DropdownMenu.SubContent>
								</DropdownMenu.Sub>

								<DropdownMenu.Separator />

								<DropdownMenu.Item
									onSelect={async () => {
										const wasFav = favoritesManager.hasTrack(entry.track.id);
										await favoritesManager.toggleTrack(entry.track.id);
										toast.success(
											wasFav ? "Removed from favorites" : "Added to favorites",
										);
									}}
								>
									{#if favoritesManager.hasTrack(entry.track.id)}
										<Heart class="fill-primary stroke-primary" />
										Unfavorite
									{:else}
										<Heart />
										Favorite
									{/if}
								</DropdownMenu.Item>
								{#if quickPlaylistManager.playlist !== null}
									<DropdownMenu.Item
										onSelect={async () => {
											const wasIn = quickPlaylistManager.hasTrack(
												entry.track.id,
											);
											await quickPlaylistManager.toggleTrack(entry.track.id);
											toast.success(
												wasIn
													? "Removed from quick playlist"
													: "Added to quick playlist",
											);
										}}
									>
										{#if quickPlaylistManager.hasTrack(entry.track.id)}
											<Star class="fill-primary stroke-primary" />
											Remove from Quick
										{:else}
											<Star />
											Quick Add
										{/if}
									</DropdownMenu.Item>
								{/if}
								<DropdownMenu.Item
									onSelect={() => saveToPlaylist(entry.track.id)}
								>
									<ListPlus />
									Save to Playlist
								</DropdownMenu.Item>
							</DropdownMenu.Content>
						</DropdownMenu.Root>
					</div>
				</div>
			{/each}
		</div>
	{/if}
</InfiniteScroll>
