<script lang="ts">
	import type { Track } from "$lib/api/types";
	import { getMusicManager } from "$lib/music-manager.svelte";
	import FavoriteButton from "$lib/components/FavoriteButton.svelte";
	import QuickAddButton from "$lib/components/QuickAddButton.svelte";
	import Image from "$lib/components/Image.svelte";
	import ArtistList from "$lib/components/ArtistList.svelte";
	import { getFavorites } from "$lib/favorites.svelte";
	import { getQuickPlaylist } from "$lib/quick-playlist.svelte";
	import { Heart, Play, Star, Music2 } from "@lucide/svelte";

	type Props = {
		tracks: Track[];
		limit?: number;
	};

	let { tracks, limit = 8 }: Props = $props();
	const musicManager = getMusicManager();
	const favorites = getFavorites();
	const quickPlaylist = getQuickPlaylist();

	const items = $derived(tracks.slice(0, limit));

	async function handlePlay(trackId: string) {
		await musicManager.addTracks({
			trackIds: tracks.map((t) => t.id),
			trackId,
		});
	}

	function formatDuration(seconds: number) {
		const m = Math.floor(seconds / 60);
		const s = seconds % 60;
		return `${m}:${s.toString().padStart(2, "0")}`;
	}

	function tagChips(track: Track) {
		return track.tags.slice(0, 3);
	}
</script>

<svelte:head>
	<title>Track variants</title>
</svelte:head>

<div class="flex flex-col gap-8">
	<!-- ============================================================ -->
	<!-- Variant A: Minimal flat rows                                 -->
	<!-- ============================================================ -->
	<section class="flex flex-col gap-2">
		<h2 class="text-lg font-bold">Variant A — Minimal rows</h2>
		<p class="text-xs text-muted-foreground">
			Flat rows, hover highlight, duration column.
		</p>
		<div class="divide-y">
			{#each items as track (track.id)}
				<div
					class="group flex items-center gap-3 rounded-lg px-2 py-2 transition-colors hover:bg-accent hover:text-accent-foreground"
				>
					<button
						class="shrink-0 overflow-hidden rounded-md"
						aria-label="Play {track.name}"
						onclick={() => handlePlay(track.id)}
					>
						<div class="relative h-12 w-12">
							<Image class="h-12 w-12" src={track.coverArt.small} alt="" />
							<div
								class="absolute inset-0 flex items-center justify-center bg-black/55 opacity-0 transition-opacity group-hover:opacity-100"
							>
								<Play size={18} class="text-white" />
							</div>
						</div>
					</button>

					<div class="min-w-0 flex-1">
						<p class="truncate text-sm font-medium">{track.name}</p>
						<ArtistList
							class="text-xs text-muted-foreground"
							artists={track.artists}
						/>
					</div>

					<span
						class="hidden shrink-0 text-xs tabular-nums text-muted-foreground sm:block"
					>
						{formatDuration(track.duration)}
					</span>

					<div class="hidden shrink-0 items-center gap-0.5 sm:flex">
						<FavoriteButton show trackId={track.id} />
						<QuickAddButton trackId={track.id} />
					</div>
				</div>
			{/each}
		</div>
	</section>

	<!-- ============================================================ -->
	<!-- Variant B: Cover-forward tile grid (like Albums page)         -->
	<!-- ============================================================ -->
	<section class="flex flex-col gap-2">
		<h2 class="text-lg font-bold">Variant B — Tile grid</h2>
		<p class="text-xs text-muted-foreground">
			Album-style cards, hover play overlay.
		</p>
		<div
			class="grid grid-cols-2 gap-3 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-5 xl:grid-cols-6"
		>
			{#each items as track (track.id)}
				<div class="group flex flex-col">
					<div class="relative">
						<button
							class="block w-full overflow-hidden rounded-lg"
							aria-label="Play {track.name}"
							onclick={() => handlePlay(track.id)}
						>
							<img
								src={track.coverArt.medium}
								alt=""
								class="aspect-square w-full object-cover transition-transform duration-300 group-hover:scale-105"
							/>
						</button>

						<button
							class="absolute right-2 bottom-2 hidden h-10 w-10 items-center justify-center rounded-full bg-primary text-primary-foreground opacity-0 shadow-lg transition-all group-hover:scale-105 group-hover:opacity-100 hover:scale-110 sm:flex"
							aria-label="Play {track.name}"
							onclick={() => handlePlay(track.id)}
						>
							<Play size={18} />
						</button>
					</div>

					<div class="flex flex-col gap-0.5 pt-2">
						<p class="truncate text-sm font-medium">{track.name}</p>
						<ArtistList
							class="text-xs text-muted-foreground"
							artists={track.artists}
						/>
						{#if tagChips(track).length > 0}
							<div class="flex gap-1.5 pt-0.5">
								{#each tagChips(track) as tag (tag)}
									<span
										class="rounded-full bg-secondary/50 px-2 py-px text-[10px] text-muted-foreground"
										>{tag}</span
									>
								{/each}
							</div>
						{/if}
					</div>
				</div>
			{/each}
		</div>
	</section>

	<!-- ============================================================ -->
	<!-- Variant C: Rich card rows                                     -->
	<!-- ============================================================ -->
	<section class="flex flex-col gap-2">
		<h2 class="text-lg font-bold">Variant C — Rich card rows</h2>
		<p class="text-xs text-muted-foreground">
			Bordered cards per row, larger cover.
		</p>
		<div class="flex flex-col gap-2">
			{#each items as track, i (track.id)}
				<div
					class="group flex items-center gap-3 rounded-lg border bg-card p-3 transition-colors hover:bg-accent hover:text-accent-foreground"
				>
					<span
						class="hidden w-5 shrink-0 text-center text-sm font-medium tabular-nums text-muted-foreground sm:block"
					>
						{i + 1}
					</span>

					<button
						class="shrink-0 overflow-hidden rounded-md"
						aria-label="Play {track.name}"
						onclick={() => handlePlay(track.id)}
					>
						<div class="relative h-14 w-14">
							<Image class="h-14 w-14" src={track.coverArt.small} alt="" />
							<div
								class="absolute inset-0 flex items-center justify-center bg-black/55 opacity-0 transition-opacity group-hover:opacity-100"
							>
								<Play size={20} class="text-white" />
							</div>
						</div>
					</button>

					<div class="min-w-0 flex-1">
						<div class="flex items-center gap-1">
							{#if favorites.hasTrack(track.id)}
								<Heart
									size={12}
									class="shrink-0 fill-primary text-primary"
								/>
							{/if}
							{#if quickPlaylist.hasTrack(track.id)}
								<Star size={12} class="shrink-0 fill-primary text-primary" />
							{/if}
							<p class="truncate text-sm font-medium">{track.name}</p>
						</div>
						<ArtistList
							class="text-xs text-muted-foreground"
							artists={track.artists}
						/>
						{#if track.tags.length > 0}
							<div class="mt-0.5 flex flex-wrap gap-1.5">
								{#each track.tags as tag (tag)}
									<span
										class="rounded-full bg-secondary/50 px-2 py-px text-[10px] text-muted-foreground"
										>{tag}</span
									>
								{/each}
							</div>
						{/if}
					</div>

					<span
						class="hidden shrink-0 text-xs tabular-nums text-muted-foreground sm:block"
					>
						{formatDuration(track.duration)}
					</span>

					<div class="flex shrink-0 items-center gap-0.5 sm:gap-1">
						<FavoriteButton show trackId={track.id} />
						<QuickAddButton trackId={track.id} />
					</div>
				</div>
			{/each}
		</div>
	</section>

	<!-- ============================================================ -->
	<!-- Variant D: Tiny compact rows                                  -->
	<!-- ============================================================ -->
	<section class="flex flex-col gap-2">
		<h2 class="text-lg font-bold">Variant D — Tiny rows</h2>
		<p class="text-xs text-muted-foreground">
			One-line name + artist, always-visible play.
		</p>
		<div class="divide-y">
			{#each items as track (track.id)}
				<div
					class="group flex items-center gap-3 rounded-lg px-2 py-1.5 transition-colors hover:bg-accent hover:text-accent-foreground"
				>
					<button
						class="relative h-10 w-10 shrink-0"
						aria-label="Play {track.name}"
						onclick={() => handlePlay(track.id)}
					>
						<Image class="h-10 w-10" src={track.coverArt.small} alt="" />
						<div
							class="absolute inset-0 flex items-center justify-center bg-black/60 opacity-0 transition-opacity group-hover:opacity-100"
						>
							<Play size={14} class="text-white" />
						</div>
					</button>

					<div class="min-w-0 flex-1 truncate text-sm">
						<span class="font-medium">{track.name}</span>
						<span class="text-muted-foreground">
							{" "}·{" "}{track.artists.map((a) => a.name).join(", ")}
						</span>
					</div>

					<span
						class="hidden shrink-0 text-xs tabular-nums text-muted-foreground sm:block"
					>
						{formatDuration(track.duration)}
					</span>

					<Music2 class="hidden h-3.5 w-3.5 shrink-0 text-muted-foreground sm:block" />
				</div>
			{/each}
		</div>
	</section>
</div>