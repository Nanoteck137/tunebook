<script lang="ts">
	import { Search, X, ListFilter, Mic } from "@lucide/svelte";
	import HeroCard from "$lib/components/HeroCard.svelte";
	import HeroIcon from "$lib/components/HeroIcon.svelte";
	import { getApiClient, handleApiError } from "$lib";
	import type { Artist } from "$lib/api/types";
	import InfiniteScroll from "$lib/components/InfiniteScroll.svelte";
	import { InfiniteScrollController } from "$lib/infinite-scroll.svelte";
	import { Button } from "$lib/components/ui";
	import { SortToggleDropdown } from "$lib/components/sort";
	import DebouncedSearchInput from "$lib/components/DebouncedSearchInput.svelte";
	import { cn } from "$lib/utils";
	import { goto } from "$app/navigation";
	import { page } from "$app/state";
	import { onMount } from "svelte";
	import {
		sortTypes,
		defaultSort,
		type SortType,
		constructFilterSort,
	} from "./types";
	import ArtistTile from "$lib/components/tiles/ArtistTile.svelte";
	import TileGrid from "$lib/components/tiles/TileGrid.svelte";
	import Spacer from "$lib/components/Spacer.svelte";

	let { data } = $props();

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

	let value = $state("");

	onMount(() => {
		value = page.url.searchParams.get("query") ?? "";
	});

	async function search(query: string) {
		const params = page.url.searchParams;
		params.delete("query");

		if (query) {
			params.set("query", query);
		}

		await goto("?" + params.toString(), {
			invalidateAll: true,
			keepFocus: true,
			replaceState: true,
		});
	}

	function clearFilters() {
		sort = defaultSort;

		const query = page.url.searchParams;
		query.delete("sort");

		goto("?" + query.toString(), { invalidateAll: true });
	}

	const apiClient = getApiClient();

	const scroll = new InfiniteScrollController<Artist>({
		initialLoad: () => ({
			items: data.artists,
			hasMore: data.page.page + 1 < data.page.totalPages,
			page: data.page.page,
		}),
		load: async (nextPage) => {
			const query: Record<string, string> = {
				page: String(nextPage),
				perPage: String(data.page.perPage),
			};

			constructFilterSort(data.filter, query);

			const res = await apiClient.getArtists({ query });
			if (!res.success) {
				handleApiError(res.error);
				return null;
			}

			return {
				items: res.data.artists,
				hasMore: res.data.page.page + 1 < res.data.page.totalPages,
			};
		},
		itemKey: (artist) => artist.id,
	});
</script>

<HeroCard
	class="section-artists"
	innerClass="sm:flex-row sm:items-center sm:justify-between"
>
	<div class="flex items-center gap-4">
		<HeroIcon>
			<Mic />
		</HeroIcon>

		<div class="flex min-w-0 flex-col">
			<h1 class="text-2xl font-bold">Artists</h1>
			<p class="text-sm text-muted-foreground">
				Every artist in your library
				{#if data.page}
					&middot; {data.page.totalItems}
				{/if}
			</p>
		</div>
	</div>
</HeroCard>

<Spacer size="md" />

<!-- Toolbar -->
<div class="flex flex-wrap items-center justify-between gap-2">
	<DebouncedSearchInput
		class="flex-1 md:max-w-64"
		placeholder="Search artists..."
		{value}
		setValue={(v) => (value = v)}
		{search}
	/>

	<div class="flex items-center gap-1">
		<Button
			variant="ghost"
			size="icon"
			href="/search/artists"
			title="Advanced search"
		>
			<Search />
		</Button>

		<SortToggleDropdown
			types={sortTypes}
			{sort}
			{defaultSort}
			onSortChange={(value) => updateSort(value)}
		/>
	</div>
</div>

<Spacer size="md" />

<InfiniteScroll controller={scroll}>
	<TileGrid>
		{#each scroll.items as artist (artist.id)}
			<ArtistTile {artist} />
		{/each}
	</TileGrid>
</InfiniteScroll>
