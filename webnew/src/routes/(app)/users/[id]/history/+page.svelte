<script lang="ts">
	import { goto } from "$app/navigation";
	import { page } from "$app/state";
	import {
		DiscAlbum,
		EllipsisVertical,
		Heart,
		History,
		Play,
		Shuffle,
		Star,
		X,
		Users,
	} from "@lucide/svelte";
	import Image from "$lib/components/Image.svelte";
	import { Button, buttonVariants, DropdownMenu } from "$lib/components/ui";
	import { getFavorites } from "$lib/favorites.svelte";
	import { getMusicManager } from "$lib/music-manager.svelte";
	import { getQuickPlaylist } from "$lib/quick-playlist.svelte";
	import InfiniteScroll from "$lib/components/InfiniteScroll.svelte";
	import { InfiniteScrollController } from "$lib/infinite-scroll.svelte";
	import { getApiClient, handleApiError } from "$lib";
	import type { TrackHistory } from "$lib/api/types";
	import { toast } from "svelte-sonner";

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
			const res = await apiClient.getTrackHistory({
				query: {
					page: String(nextPage),
					perPage: String(data.page.perPage),
				},
			});
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

	let yearParam = $derived(page.url.searchParams.get("year"));

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
		if (status === "completed") return "Completed";
		if (status === "skipped") return "Skipped";
		return "In Progress";
	}

	function statusClass(status: string) {
		if (status === "completed")
			return "bg-green-500/10 text-green-500 ring-green-500/25";
		if (status === "skipped")
			return "bg-muted text-muted-foreground ring-foreground/10";
		return "bg-yellow-500/10 text-yellow-500 ring-yellow-500/25";
	}

	function progressBg(pct: number): string {
		if (pct >= 80) return "bg-green-500/20";
		if (pct >= 40) return "bg-yellow-500/20";
		return "bg-primary/15";
	}

	let totalListeningTime = $derived(
		scroll.items.reduce(
			(sum, e) => sum + e.track.duration * (e.percentPlayed / 100),
			0,
		),
	);

	function formatDuration(seconds: number): string {
		const h = Math.floor(seconds / 3600);
		const m = Math.floor((seconds % 3600) / 60);
		if (h > 0) return `${h}h ${m}m`;
		return `${m}m`;
	}

	function clearYearFilter() {
		const query = page.url.searchParams;
		query.delete("year");
		goto(`?${query.toString()}`, { invalidateAll: true });
	}

	async function playAll() {
		await musicManager.addTracks({
			trackIds: scroll.items.map((e) => e.track.id),
		});
	}

	async function shufflePlay() {
		const ids = scroll.items.map((e) => e.track.id);
		for (let i = ids.length - 1; i > 0; i--) {
			const j = Math.floor(Math.random() * (i + 1));
			[ids[i], ids[j]] = [ids[j], ids[i]];
		}
		await musicManager.addTracks({ trackIds: ids });
	}

	function playTrack(trackId: string) {
		const trackIds = scroll.items.map((e) => e.track.id);
		musicManager.addTracks({ trackIds, trackId });
	}
</script>

<div class="flex flex-col gap-6">
	<div
		class="flex flex-col gap-6 rounded-lg border bg-linear-to-b from-[oklch(0.93_0.045_75)] to-background p-4 shadow-sm sm:p-6 md:flex-row md:items-end md:gap-8 dark:from-[oklch(0.24_0.03_80)] dark:to-background"
	>
		<div class="flex min-w-0 flex-col gap-2">
			<p
				class="text-xs font-semibold tracking-wider text-muted-foreground uppercase"
			>
				Listening History
				{#if yearParam}
					&middot; {yearParam}
				{/if}
			</p>

			<h1 class="line-clamp-2 text-2xl font-bold md:text-4xl">
				{#if yearParam}
					History for {yearParam}
				{:else}
					All Tracks
				{/if}
			</h1>

			<p class="text-sm text-muted-foreground">
				{data.page.totalItems}
				{data.page.totalItems === 1 ? "play" : "plays"}
				&middot; {formatDuration(totalListeningTime)}
			</p>

			<div class="flex gap-2 pt-2">
				<Button onclick={() => playAll()}>
					<Play />
					Play
				</Button>
				<Button variant="ghost" size="icon" onclick={() => shufflePlay()}>
					<Shuffle />
				</Button>
				<DropdownMenu.Root>
					<DropdownMenu.Trigger
						class={buttonVariants({ variant: "ghost", size: "icon" })}
					>
						<EllipsisVertical />
					</DropdownMenu.Trigger>
					<DropdownMenu.Content align="start">
						<DropdownMenu.Group>
							{#if yearParam}
								<DropdownMenu.Item onSelect={clearYearFilter}>
									<X />
									Clear year filter
								</DropdownMenu.Item>
							{/if}
						</DropdownMenu.Group>
					</DropdownMenu.Content>
				</DropdownMenu.Root>
			</div>
		</div>
	</div>

	{#if scroll.items.length === 0}
		<div class="flex flex-col items-center gap-2 rounded-lg border py-16">
			<History size={32} class="text-muted-foreground/40" />
			<p class="text-sm text-muted-foreground">No listening history yet</p>
		</div>
	{:else}
		<InfiniteScroll controller={scroll}>
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
										<Star
											size={12}
											class="shrink-0 fill-primary text-primary"
										/>
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
											{#if i > 0}{", "}{/if}
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

							<span
								class="shrink-0 text-[11px] text-muted-foreground sm:hidden"
							>
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
										<DropdownMenu.Item
											onSelect={() => {
												const trackIds = scroll.items.map((e) => e.track.id);
												for (let i = trackIds.length - 1; i > 0; i--) {
													const j = Math.floor(Math.random() * (i + 1));
													[trackIds[i], trackIds[j]] = [
														trackIds[j],
														trackIds[i],
													];
												}
												musicManager.addTracks({ trackIds });
											}}
										>
											<Shuffle />
											Shuffle play
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
											<Users />
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
												wasFav
													? "Removed from favorites"
													: "Added to favorites",
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
								</DropdownMenu.Content>
							</DropdownMenu.Root>
						</div>
					</div>
				{/each}
			</div>
		</InfiniteScroll>
	{/if}
</div>
