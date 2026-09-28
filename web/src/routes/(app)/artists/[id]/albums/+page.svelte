<script lang="ts">
	import { goto } from "$app/navigation";
	import { page } from "$app/state";
	import { Disc } from "@lucide/svelte";
	import { SortToggleDropdown } from "$lib/components/sort";
	import Image from "$lib/components/Image.svelte";
	import SectionHeader from "$lib/components/SectionHeader.svelte";
	import { getApiClient, handleApiError } from "$lib";
	import type { Album } from "$lib/api/types";
	import InfiniteScroll from "$lib/components/InfiniteScroll.svelte";
	import { InfiniteScrollController } from "$lib/infinite-scroll.svelte";
	import { sortTypes, defaultSort, applySort, type SortType } from "./types";
	import TileGrid from "$lib/components/tiles/TileGrid.svelte";
	import Spacer from "$lib/components/Spacer.svelte";
	import AlbumTile from "$lib/components/tiles/AlbumTile.svelte";

	let { data } = $props();
	const apiClient = getApiClient();

	let sort = $state(
		(page.url.searchParams.get("sort") as SortType) ?? defaultSort,
	);

	function updateSort(value: string) {
		sort = value as SortType;

		const query = page.url.searchParams;
		query.delete("sort");

		if (sort !== defaultSort) {
			query.set("sort", sort);
		}

		goto("?" + query.toString(), { invalidateAll: true });
	}

	const scroll = new InfiniteScrollController<Album>({
		initialLoad: () => ({
			items: data.albums,
			hasMore: data.page.page + 1 < data.page.totalPages,
			page: data.page.page,
		}),
		load: async (nextPage) => {
			const query: Record<string, string> = {
				page: String(nextPage),
				perPage: String(data.page.perPage),
				filter: `artistId = "${data.artist.id}" or featuringArtists has "${data.artist.id}"`,
			};

			applySort(data.filter, query);

			const res = await apiClient.getAlbums({ query });
			if (!res.success) {
				handleApiError(res.error);
				return null;
			}

			return {
				items: res.data.albums,
				hasMore: res.data.page.page + 1 < res.data.page.totalPages,
			};
		},
		itemKey: (album) => album.id,
	});
</script>

<SectionHeader count={data.page.totalItems}>
	<Disc />
	Albums

	{#snippet actions()}
		<SortToggleDropdown
			types={sortTypes}
			{sort}
			{defaultSort}
			onSortChange={updateSort}
		/>
	{/snippet}
</SectionHeader>

<Spacer size="md" />

<InfiniteScroll controller={scroll}>
	<TileGrid>
		{#each scroll.items as album (album.id)}
			<AlbumTile {album} />
		{/each}
	</TileGrid>
</InfiniteScroll>
