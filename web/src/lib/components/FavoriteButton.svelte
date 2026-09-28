<script lang="ts">
	import { getFavorites } from "$lib/favorites.svelte";
	import { Button } from "$lib/components/ui";
	import { Heart } from "@lucide/svelte";

	type Props = {
		show: boolean;
		trackId: string;
	};

	const { show, trackId }: Props = $props();
	const favorites = getFavorites();
</script>

{#if show}
	<Button
		type="submit"
		class="rounded-full dark:hover:bg-muted hover:bg-muted"
		variant="ghost"
		size="icon-lg"
		onclick={() => {
			favorites.toggleTrack(trackId);
		}}
		title="Favorite"
		disabled={favorites.loading}
	>
		{#if favorites.hasTrack(trackId)}
			<Heart class="fill-primary stroke-primary"/>
		{:else}
			<Heart />
		{/if}
	</Button>
{/if}
