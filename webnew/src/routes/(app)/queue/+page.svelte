<script lang="ts">
	import { getMusicManager, type MediaItem } from "$lib/music-manager.svelte";
	import Image from "$lib/components/Image.svelte";
	import { getApiClient, handleApiError } from "$lib";
	import { InfiniteScrollController } from "$lib/infinite-scroll.svelte";
	import InfiniteScroll from "$lib/components/InfiniteScroll.svelte";
	import {
		DiscAlbum,
		EllipsisVertical,
		Heart,
		ListPlus,
		ListX,
		Play,
		Star,
		User,
	} from "@lucide/svelte";
	import { fly } from "svelte/transition";
	import { Button, buttonVariants, DropdownMenu } from "$lib/components/ui";
	import { cn } from "$lib/utils";
	import HeroCard from "$lib/components/HeroCard.svelte";
	import { goto } from "$app/navigation";
	import { getFavorites } from "$lib/favorites.svelte";
	import { getQuickPlaylist } from "$lib/quick-playlist.svelte";
	import { showPlaylistModal } from "$lib/playlist-modal.svelte";
	import { toast } from "svelte-sonner";
	import FavoriteButton from "$lib/components/FavoriteButton.svelte";
	import QuickAddButton from "$lib/components/QuickAddButton.svelte";
	import type { ArtistInfo } from "$lib/api/types";

	let { data } = $props();

	const musicManager = getMusicManager();
	const apiClient = getApiClient();
	const favoritesManager = getFavorites();
	const quickPlaylistManager = getQuickPlaylist();
	const deviceId = localStorage.getItem("device-id") ?? "";

	type QueueListItem = {
		queueItemId: string;
		trackId: string;
		name: string;
		albumId: string;
		albumName: string;
		artists: ArtistInfo[];
		artistNames: string;
		coverArt: string;
		duration: number;
	};

	let currentIndex = $state(0);
	let totalItems = $state(0);

	let currentMediaItem = $state<MediaItem | null>(null);

	$effect(() => {
		currentIndex = musicManager.queue.index;
		currentMediaItem = musicManager.currentItem;
	});

	const scroll = new InfiniteScrollController<QueueListItem>({
		initialLoad: () => ({ items: [], hasMore: true, page: -1 }),
		load: async (page) => {
			const res = await data.apiClient.getQueue(deviceId, {
				query: { page: page.toString() },
			});

			if (!res.success) {
				handleApiError(res.error);
				return null;
			}

			currentIndex = res.data.currentIndex;
			totalItems = res.data.page.totalItems;

			const next = res.data.items.map((item) => ({
				queueItemId: item.queueItemId,
				trackId: item.track.id,
				name: item.track.name,
				albumId: item.track.albumId,
				albumName: item.track.albumName,
				artists: item.track.artists,
				artistNames: item.track.artists.map((a) => a.name).join(", "),
				coverArt: item.track.coverArt.small,
				duration: item.track.duration,
			}));

			return { items: next, hasMore: page < res.data.page.totalPages - 1 };
		},
		itemKey: (item) => item.queueItemId,
	});

	function formatDuration(seconds: number): string {
		const m = Math.floor(seconds / 60);
		const s = Math.floor(seconds % 60);
		return `${m}:${s.toString().padStart(2, "0")}`;
	}

	async function playQueueItem(index: number) {
		await musicManager.setQueueIndex(index);
		musicManager.play();
	}

	async function clearQueue() {
		await musicManager.clearQueue();
		totalItems = 0;
		currentIndex = 0;
		scroll.setInitial({ items: [], hasMore: false, page: 0 });
	}

	async function toggleFavorite(trackId: string) {
		const wasFav = favoritesManager.hasTrack(trackId);
		await favoritesManager.toggleTrack(trackId);
		toast.success(wasFav ? "Removed from favorites" : "Added to favorites");
	}

	async function toggleQuickPlaylist(trackId: string) {
		const wasIn = quickPlaylistManager.hasTrack(trackId);
		await quickPlaylistManager.toggleTrack(trackId);
		toast.success(
			wasIn ? "Removed from quick playlist" : "Added to quick playlist",
		);
	}

	async function saveToPlaylist(trackId: string) {
		const playlist = await showPlaylistModal();
		if (!playlist) return;

		const res = await apiClient.addItemToPlaylist(playlist.id, { trackId });
		if (!res.success) {
			handleApiError(res.error);
		}
	}
</script>

<svelte:head>
	<title>Queue - Tunebook</title>
</svelte:head>

<div class="section-queue flex flex-col gap-4">
	<HeroCard>
		<div class="flex min-w-0 flex-col gap-2">
			<p
				class="text-xs font-semibold tracking-wider text-muted-foreground uppercase"
			>
				Up Next
			</p>

			<h1 class="text-2xl font-bold sm:text-4xl">Queue</h1>

			<p class="text-sm text-muted-foreground">
				{#if totalItems > 0}
					Track {currentIndex + 1} of {totalItems}
				{:else}
					No tracks in queue
				{/if}
			</p>

			{#if totalItems > 0}
				<div class="flex items-center gap-2 pt-2">
					<Button
						variant="outline"
						size="sm"
						onclick={clearQueue}
						class="shrink-0"
					>
						<ListX />
						<span class="hidden sm:inline">Clear queue</span>
					</Button>
				</div>
			{/if}
		</div>
	</HeroCard>

	{#if currentMediaItem}
		<div
			class="relative overflow-hidden rounded-xl border bg-accent/30"
			in:fly={{ y: 12, duration: 250 }}
		>
			<Image
				class="pointer-events-none absolute inset-0 h-full w-full scale-110 object-cover opacity-25 blur-2xl"
				src={currentMediaItem.coverArt}
				alt=""
			/>
			<div
				class="pointer-events-none absolute inset-0 bg-linear-to-br from-background/80 via-background/40 to-background/90"
			></div>

			<div class="relative flex items-center gap-4 p-4 sm:p-5">
				<button
					onclick={async () => {
						await musicManager.setQueueIndex(currentIndex);
						musicManager.play();
					}}
					class="shrink-0 cursor-pointer"
					aria-label="Go to current track"
				>
					<Image
						class="h-20 w-20 rounded-lg shadow-lg sm:h-28 sm:w-28"
						src={currentMediaItem.coverArt}
						alt={currentMediaItem.name}
					/>
				</button>

				<div class="flex min-w-0 flex-1 flex-col">
					<p
						class="text-xs font-semibold tracking-widest text-primary uppercase"
					>
						Now Playing
					</p>
					<p class="mt-0.5 truncate text-lg font-bold sm:text-2xl">
						{currentMediaItem.name}
					</p>
					<p class="truncate text-sm text-muted-foreground">
						{#each currentMediaItem.artists as artist, i (artist.id)}
							{#if i > 0},
							{/if}{artist.name}
						{/each}
						{#if currentMediaItem.album.name}
							<span class="text-muted-foreground/70">
								· {currentMediaItem.album.name}
							</span>
						{/if}
					</p>
					<div class="mt-2 flex items-center gap-1">
						<FavoriteButton show trackId={currentMediaItem.trackId} />
						<QuickAddButton trackId={currentMediaItem.trackId} />
					</div>
				</div>
			</div>
		</div>
	{/if}

	<InfiniteScroll controller={scroll} rootMargin="200px">
		{#snippet loadingSnippet()}
			<p class="py-6 text-center text-sm text-muted-foreground">Loading...</p>
		{/snippet}

		{#if scroll.items.length === 0 && !scroll.loading}
			<p class="py-12 text-center text-muted-foreground">
				The queue is empty.
			</p>
		{:else}
			<div class="flex flex-col gap-1">
				{#each scroll.items as item, index (item.queueItemId)}
					{@const isCurrent = index === currentIndex}
					{#if index === 0 && currentIndex > 0}
						<p
							class="px-2 pt-4 pb-1 text-xs font-semibold tracking-wider text-muted-foreground uppercase"
						>
							Played
						</p>
					{/if}
					{#if index === currentIndex + 1}
						<p
							class="px-2 pt-4 pb-1 text-xs font-semibold tracking-wider text-muted-foreground uppercase"
						>
							Up Next
						</p>
					{/if}

					<div
						class="group flex w-full items-center gap-2 rounded-md px-2 py-1 transition-colors hover:bg-accent has-data-[state='open']:bg-accent {isCurrent
							? 'bg-accent/80'
							: ''} sm:gap-3 sm:p-2"
						in:fly={{ y: 8, duration: 200 }}
					>
						<button
							class="flex min-w-0 flex-1 cursor-pointer items-center gap-2 text-left sm:gap-3"
							onclick={() => playQueueItem(index)}
						>
							<span
								class="hidden w-6 flex-shrink-0 text-right text-xs text-muted-foreground tabular-nums sm:block sm:w-8 sm:text-sm"
							>
								{index + 1}
							</span>

							<Image
								class="h-9 w-9 flex-shrink-0 rounded object-cover sm:h-12 sm:w-12"
								src={item.coverArt}
								alt={item.name}
							/>

							<div class="flex min-w-0 flex-1 flex-col">
								<p class="truncate text-sm font-medium sm:text-base">
									{item.name}
								</p>
								<p class="truncate text-xs text-muted-foreground sm:text-sm">
									{item.artistNames}
									{#if item.albumName}
										<span class="text-muted-foreground/70">
											· {item.albumName}
										</span>
									{/if}
								</p>
							</div>

							<span
								class="hidden flex-shrink-0 text-xs text-muted-foreground tabular-nums sm:block sm:text-sm"
							>
								{formatDuration(item.duration)}
							</span>

							{#if isCurrent && musicManager.playing}
								<span
									class="ml-1 flex h-4 flex-shrink-0 items-end gap-[3px] text-primary sm:ml-2"
									aria-label="Playing"
								>
									<span class="eq-bar"></span>
									<span class="eq-bar" style="animation-delay: 150ms"></span>
									<span class="eq-bar" style="animation-delay: 300ms"></span>
								</span>
							{:else if isCurrent}
								<Play
									class="ml-1 h-4 w-4 flex-shrink-0 text-primary sm:ml-2 sm:h-5 sm:w-5"
								/>
							{:else}
								<Play
									class="ml-1 h-4 w-4 flex-shrink-0 text-muted-foreground opacity-0 transition-opacity group-hover:opacity-100 sm:ml-2 sm:h-5 sm:w-5"
								/>
							{/if}
						</button>

						<DropdownMenu.Root>
							<DropdownMenu.Trigger
								class={cn(
									buttonVariants({ variant: "ghost", size: "icon-sm" }),
									"-mr-1 shrink-0 rounded-full text-muted-foreground",
								)}
								title="More options"
								aria-label="More options"
							>
								<EllipsisVertical />
							</DropdownMenu.Trigger>
							<DropdownMenu.Content align="end">
								<DropdownMenu.Group>
									<DropdownMenu.Item onSelect={() => playQueueItem(index)}>
										<Play />
										Play
									</DropdownMenu.Item>
								</DropdownMenu.Group>
								<DropdownMenu.Separator />

								<DropdownMenu.Group>
									<DropdownMenu.Item
										onSelect={() => {
											goto(`/albums/${item.albumId}`);
										}}
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
											{#each item.artists as artist (artist.id)}
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
										onSelect={() => saveToPlaylist(item.trackId)}
									>
										<ListPlus />
										Add to Playlist
									</DropdownMenu.Item>
									<DropdownMenu.Item
										onSelect={() => toggleFavorite(item.trackId)}
									>
										{#if favoritesManager.hasTrack(item.trackId)}
											<Heart class="fill-primary stroke-primary" />
											Unfavorite
										{:else}
											<Heart />
											Favorite
										{/if}
									</DropdownMenu.Item>
									{#if quickPlaylistManager.playlist !== null}
										<DropdownMenu.Item
											onSelect={() => toggleQuickPlaylist(item.trackId)}
										>
											{#if quickPlaylistManager.hasTrack(item.trackId)}
												<Star class="fill-primary stroke-primary" />
												Remove from Quick
											{:else}
												<Star />
												Quick Add
											{/if}
										</DropdownMenu.Item>
									{/if}
								</DropdownMenu.Group>
							</DropdownMenu.Content>
						</DropdownMenu.Root>
					</div>
				{/each}
			</div>
		{/if}
	</InfiniteScroll>
</div>

<style>
	.eq-bar {
		width: 3px;
		height: 14px;
		border-radius: 2px;
		background: currentColor;
		transform-origin: bottom;
		animation: eq-bounce 0.9s ease-in-out infinite;
	}

	@keyframes eq-bounce {
		0%,
		100% {
			transform: scaleY(0.25);
		}
		50% {
			transform: scaleY(1);
		}
	}
</style>
