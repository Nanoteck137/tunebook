<script lang="ts">
	import type { Track } from "$lib/api/types";
	import ArtistList from "$lib/components/ArtistList.svelte";
	import Image from "$lib/components/Image.svelte";
	import { cn } from "$lib/utils";
	import { Heart, Play, Star } from "@lucide/svelte";
	import type { Snippet } from "svelte";
	import { getFavorites } from "$lib/favorites.svelte";
	import { getQuickPlaylist } from "$lib/quick-playlist.svelte";

	type Props = {
		class?: string;
		showNumber?: boolean;
		displayOrder?: boolean;
		track: Track;
		children?: Snippet;

		onPlayClicked?: () => void;
	};

	const {
		class: className,
		showNumber,
		displayOrder,
		track,
		children,
		onPlayClicked,
	}: Props = $props();
	const favoriteManager = getFavorites();
	const quickPlaylistManager = getQuickPlaylist();

	let isFav = $derived(favoriteManager.hasTrack(track.id));
	let isQuick = $derived(
		quickPlaylistManager.playlist !== null &&
			quickPlaylistManager.hasTrack(track.id),
	);

	function formatDuration(seconds: number) {
		const m = Math.floor(seconds / 60);
		const s = seconds % 60;
		return `${m}:${s.toString().padStart(2, "0")}`;
	}
</script>

<div
	class={cn(
		"group hover:bg-accent hover:text-accent-foreground has-data-[state='open']:bg-accent has-data-[state='open']:text-accent-foreground flex items-center gap-3 rounded-lg p-2 transition-colors",
		className,
	)}
>
	<button
		class="shrink-0"
		onclick={() => onPlayClicked?.()}
		aria-label="Play {track.name}"
	>
		{#if showNumber}
			<div class="flex h-12 w-12 items-center justify-center overflow-hidden rounded-md">
				<span
					class="text-sm font-medium tabular-nums group-hover:hidden"
				>
					{track.number}.
				</span>
				<div class="hidden items-center justify-center group-hover:flex">
					<Play size={20} />
				</div>
			</div>
		{:else}
			<div class="relative h-12 w-12">
				<Image class="h-12 w-12" src={track.coverArt.small} alt="" />
				<div
					class="absolute inset-0 flex items-center justify-center rounded-md bg-black/55 opacity-0 transition-opacity group-hover:opacity-100"
				>
					<Play size={18} class="text-white" />
				</div>
			</div>
		{/if}
	</button>

	<div class="flex min-w-0 flex-1 flex-col gap-0.5">
		<p class="truncate text-sm font-medium" title={track.name}>
			{#if isFav}
				<Heart
					size={12}
					class="mr-0.5 inline fill-primary text-primary sm:hidden"
				/>
			{/if}
			{#if isQuick}
				<Star
					size={12}
					class="mr-0.5 inline fill-primary text-primary sm:hidden"
				/>
			{/if}
			{#if displayOrder}
				{track.order}.
			{/if}
			{track.name}
		</p>

		<ArtistList class="text-muted-foreground" artists={track.artists} />

		{#if track.tags.length > 0}
			<div class="hidden min-w-0 gap-1.5 sm:flex">
				{#each track.tags as tag (tag)}
					<span
						class="shrink-0 rounded-full bg-secondary/50 px-2 py-px text-[10px] text-muted-foreground"
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
		{@render children?.()}
	</div>
</div>