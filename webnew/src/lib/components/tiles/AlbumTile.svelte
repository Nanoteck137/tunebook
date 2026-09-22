<script lang="ts">
	import { Play } from "@lucide/svelte";
	import ArtistList, { type Artist } from "$lib/components/ArtistList.svelte";
	import { getMusicManager } from "$lib/music-manager.svelte";

	type Props = {
		id: string;
		cover: string;
		name: string;
		artists: Artist[];
	};

	const { id, cover, name, artists }: Props = $props();
	const musicManager = getMusicManager();

	async function play() {
		await musicManager.queueRequest({ type: "addAlbum", albumId: id }, {});
	}
</script>

<div class="group relative flex w-40 shrink-0 flex-col">
	<div class="relative overflow-hidden rounded-lg">
		<a href="/albums/{id}" class="block">
			<img
				class="aspect-square w-40 object-cover transition-transform duration-300 group-hover:scale-105"
				src={cover}
				alt=""
				title={name}
			/>
		</a>

		<button
			class="absolute right-2 bottom-2 hidden h-9 w-9 translate-y-2 cursor-pointer items-center justify-center rounded-full bg-primary text-primary-foreground opacity-0 shadow-lg transition-all duration-300 group-hover:translate-y-0 group-hover:opacity-100 hover:scale-110 sm:flex"
			title="Play album"
			aria-label={`Play ${name}`}
			onclick={play}
		>
			<Play size={16} />
		</button>
	</div>

	<a
		href="/albums/{id}"
		class="mt-2 truncate text-sm font-medium hover:underline"
		title={name}
	>
		{name}
	</a>

	<ArtistList class="justify-start text-muted-foreground" {artists} />
</div>

