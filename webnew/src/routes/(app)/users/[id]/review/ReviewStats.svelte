<script lang="ts">
	import { formatPlayTime } from "$lib/utils";

	let {
		trackCount,
		uniqueTracks,
		listeningTime,
		avgCompletion,
		skipCount,
		favoritePlays,
	}: {
		trackCount: number;
		uniqueTracks: number;
		listeningTime: number;
		avgCompletion: number;
		skipCount: number;
		favoritePlays: number;
	} = $props();

	function formatPercent(rate: number): string {
		return `${rate.toLocaleString("en-US", { maximumFractionDigits: 1 })}%`;
	}

	let skipRate = $derived(trackCount > 0 ? (skipCount / trackCount) * 100 : 0);

	let favoritesOverlap = $derived(
		trackCount > 0 ? Math.round((favoritePlays / trackCount) * 100) : 0,
	);
</script>

<dl class="grid grid-cols-2 gap-x-4 gap-y-2.5 text-xs sm:grid-cols-3">
	<div class="flex flex-col">
		<dt class="text-muted-foreground">Plays</dt>
		<dd class="font-semibold tabular-nums">{trackCount.toLocaleString()}</dd>
	</div>

	<div class="flex flex-col">
		<dt class="text-muted-foreground">Unique</dt>
		<dd class="font-semibold tabular-nums">{uniqueTracks.toLocaleString()}</dd>
	</div>

	<div class="flex flex-col">
		<dt class="text-muted-foreground">Play time</dt>
		<dd class="font-semibold tabular-nums">{formatPlayTime(listeningTime)}</dd>
	</div>

	<div class="flex flex-col">
		<dt class="text-muted-foreground">Avg completion</dt>
		<dd class="font-semibold tabular-nums">{formatPercent(avgCompletion)}</dd>
	</div>

	<div class="flex flex-col">
		<dt class="text-muted-foreground">Skip rate</dt>
		<dd class="font-semibold tabular-nums">{formatPercent(skipRate)}</dd>
	</div>

	<div class="flex flex-col">
		<dt class="text-muted-foreground">Favorites</dt>
		<dd class="font-semibold tabular-nums">{favoritesOverlap}%</dd>
	</div>
</dl>
