<script lang="ts">
	import { getMusicManager } from "$lib/music-manager.svelte";
	import type { Track } from "$lib/api/types";
	import Tile from "./Tile.svelte";

	type Props = {
		track: Track;
		size?: "default" | "sm";
		class?: string;
	};

	const { track, size = "default", class: className }: Props = $props();
	const musicManager = getMusicManager();

	async function play() {
		await musicManager.addTracks({
			trackIds: [track.id],
			trackId: track.id,
			clear: true,
		});
	}
</script>

{#snippet subtitle()}
	{#if track.artists.length > 0}
		<p
			class="line-clamp-1 text-xs text-ellipsis text-muted-foreground"
			title={track.artists.map((a) => a.name).join(", ")}
		>
			{track.artists.map((a) => a.name).join(", ")}
		</p>
	{/if}
{/snippet}

<Tile
	href="/albums/{track.albumId}?track={track.id}"
	src={track.coverArt.medium}
	name={track.name}
	sectionClass="section-tracks"
	nameClass={size === "sm" ? "line-clamp-1" : "line-clamp-2"}
	onPlay={play}
	{subtitle}
	class={className}
/>
