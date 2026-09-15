<script lang="ts">
	import { goto } from "$app/navigation";
	import { onMount } from "svelte";
	import { getApiClient, handleApiError } from "$lib";
	import type { Track } from "$lib/api/types";
	import TrackList from "$lib/components/track-list/TrackList.svelte";
	import InfiniteScroll from "$lib/components/InfiniteScroll.svelte";
	import { InfiniteScrollController } from "$lib/infinite-scroll.svelte";
	import SearchBarHeader from "../SearchBarHeader.svelte";

	let { data } = $props();

	const apiClient = getApiClient();

	async function doSearch(query: string) {
		await goto(`/search/tracks?query=${query}`, {
			invalidateAll: true,
			keepFocus: true,
			replaceState: true,
		});
	}

	function clearSearch() {
		value = "";
		doSearch("");
	}

	let value = $state("");

	onMount(() => {
		value = data.query;
	});

	const scroll = new InfiniteScrollController<Track>({
		initialLoad: () => {
			if (!data.page) return { items: [], hasMore: false, page: 0 };

			return {
				items: data.tracks,
				hasMore: data.page.page + 1 < data.page.totalPages,
				page: data.page.page,
			};
		},
		load: async (page) => {
			const res = await apiClient.searchTracks({
				query: {
					query: data.query,
					page: String(page),
					perPage: String(data.page?.perPage ?? 50),
				},
			});

			if (!res.success) {
				handleApiError(res.error);
				return null;
			}

			return {
				items: res.data.tracks,
				hasMore: res.data.page.page + 1 < res.data.page.totalPages,
			};
		},
		itemKey: (track) => track.id,
	});
</script>

<svelte:head>
	<title>Search Tracks - Tunebook</title>
</svelte:head>

<div class="flex flex-col gap-6">
	<SearchBarHeader
		searchBarPlaceholder="Search tracks..."
		{value}
		setValue={(v) => {
			value = v;
		}}
		search={doSearch}
		searchWithValue={() => {
			doSearch(value);
		}}
		{clearSearch}
	/>

	{#if data.query && scroll.items.length === 0}
		<p class="py-12 text-center text-sm text-muted-foreground">
			No tracks found for "{data.query}".
		</p>
	{/if}

	{#if scroll.items.length > 0}
		{#if data.page}
			<div class="flex items-baseline gap-2">
				<span class="text-sm text-muted-foreground">
					{data.page.totalItems} track(s)
				</span>
			</div>
		{/if}

		<InfiniteScroll controller={scroll}>
			<TrackList
				totalTracks={scroll.items.length}
				tracks={scroll.items}
				onPlay={() => {}}
			/>
		</InfiniteScroll>
	{/if}
</div>
