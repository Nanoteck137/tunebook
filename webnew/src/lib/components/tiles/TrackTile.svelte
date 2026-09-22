<script lang="ts">
	import { Play } from "@lucide/svelte";
	import ArtistList, { type Artist } from "$lib/components/ArtistList.svelte";
	import { getMusicManager } from "$lib/music-manager.svelte";

	type Props = {
		id: string;
		albumId?: string;
		cover: string;
		name: string;
		artists: Artist[];
	};

	const { id, albumId, cover, name, artists }: Props = $props();
	const musicManager = getMusicManager();

	let trackHref = $derived(albumId ? `/albums/${albumId}?track=${id}` : "#");

	async function play() {
		await musicManager.addTracks({
			trackIds: [id],
			trackId: id,
			clear: true,
		});
	}
</script>

<div class="group relative flex w-40 shrink-0 flex-col">
	<div class="relative overflow-hidden rounded-lg">
		<!-- svelte-ignore a11y_invalid_attribute -->
		<a href={trackHref} class="block">
			<img
				class="aspect-square w-40 object-cover transition-transform duration-300 group-hover:scale-105"
				src={cover}
				alt=""
				title={name}
			/>
		</a>

		<button
			class="absolute right-2 bottom-2 hidden h-9 w-9 translate-y-2 cursor-pointer items-center justify-center rounded-full bg-primary text-primary-foreground opacity-0 shadow-lg transition-all duration-300 group-hover:translate-y-0 group-hover:opacity-100 hover:scale-110 sm:flex"
			title="Play track"
			aria-label={`Play ${name}`}
			onclick={play}
		>
			<Play size={16} />
		</button>
	</div>

	<!-- svelte-ignore a11y_invalid_attribute -->
	<a
		href={trackHref}
		class="mt-2 w-40 truncate text-sm font-medium hover:underline"
		title={name}
	>
		{name}
	</a>

	<ArtistList class="w-40 justify-start text-muted-foreground" {artists} />
</div>
