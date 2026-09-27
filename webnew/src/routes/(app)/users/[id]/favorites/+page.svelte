<script lang="ts">
	import { goto } from "$app/navigation";
	import { page } from "$app/state";
	import { onMount } from "svelte";
	import { Heart, Play, Shuffle } from "@lucide/svelte";
	import { Button } from "$lib/components/ui";
	import { getMusicManager } from "$lib/music-manager.svelte";
	import { getApiClient, handleApiError } from "$lib";
	import type { Track } from "$lib/api/types";
	import InfiniteScroll from "$lib/components/InfiniteScroll.svelte";
	import { InfiniteScrollController } from "$lib/infinite-scroll.svelte";
	import TrackList from "$lib/components/track-list/TrackList.svelte";
	import SectionHeader from "$lib/components/SectionHeader.svelte";
	import DebouncedSearchInput from "$lib/components/DebouncedSearchInput.svelte";
	import { SortToggleDropdown } from "$lib/components/sort";
	import SavedFilterButton from "$lib/components/SavedFilterButton.svelte";
	import SavedFilterCard from "$lib/components/SavedFilterCard.svelte";
	import {
		sortTypes,
		defaultSort,
		type SortType,
		constructFilterSort,
	} from "../../../library/favorites/types";
	import Spacer from "$lib/components/Spacer.svelte";

	let { data } = $props();

	const musicManager = getMusicManager();
	const apiClient = getApiClient();

	let filterId = $derived(page.url.searchParams.get("filterId"));

	let filterOpen = $state(false);

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

	function play(shuffle = false) {
		return musicManager.queueRequest(
			{
				type: "addFavorites",
				userId: data.userData.id,
				filterId: filterId ?? undefined,
			},
			{ shuffle },
		);
	}

	const scroll = new InfiniteScrollController<Track>({
		initialLoad: () => ({
			items: data.tracks,
			hasMore: data.page.page + 1 < data.page.totalPages,
			page: data.page.page,
		}),
		load: async (nextPage) => {
			const query: Record<string, string> = {
				page: String(nextPage),
				perPage: String(data.page.perPage),
			};

			constructFilterSort(data.filter, query);

			if (filterId) {
				query["filterId"] = filterId;
			}

			const res = await apiClient.getUserTrackFavoritesById(data.userData.id, {
				query,
			});
			if (!res.success) {
				handleApiError(res.error);
				return null;
			}

			return {
				items: res.data.items,
				hasMore: res.data.page.page + 1 < res.data.page.totalPages,
			};
		},
		itemKey: (track) => track.id,
	});

	let hasActiveFilter = $derived(
		data.filter.query !== "" || filterId !== null,
	);
</script>

<div class="flex flex-col gap-4">
	<SectionHeader count={data.page.totalItems}>
		<Heart />
		Favorites

		{#snippet actions()}
			<div class="flex items-center gap-2">
				<Button size="sm" onclick={() => play()}>
					<Play />
					Play
				</Button>
				<Button
					size="icon-sm"
					variant="ghost"
					onclick={() => play(true)}
					title="Shuffle play"
					aria-label="Shuffle play"
				>
					<Shuffle />
				</Button>
			</div>
		{/snippet}
	</SectionHeader>

	<div class="flex flex-wrap items-center justify-between gap-2">
		<DebouncedSearchInput
			class="flex-1 md:max-w-64"
			placeholder="Search favorites..."
			{value}
			setValue={(v) => (value = v)}
			{search}
		/>

		<div class="flex items-center gap-1">
			<SavedFilterButton
				bind:filterOpen
				hasFilters={data.filters && data.filters.length > 0}
			/>

			<SortToggleDropdown
				types={sortTypes}
				{sort}
				{defaultSort}
				onSortChange={(value) => updateSort(value)}
			/>
		</div>
	</div>

	<SavedFilterCard {filterOpen} filters={data.filters} />
</div>

<Spacer size="md" />

<InfiniteScroll controller={scroll}>
	{#if scroll.items.length === 0}
		<div class="flex flex-col items-center gap-2 rounded-lg border py-16">
			<Heart size={32} class="text-muted-foreground/40" />
			<p class="text-sm text-muted-foreground">
				{hasActiveFilter
					? "No favorites match the current search or filter"
					: "No favorites yet"}
			</p>
		</div>
	{:else}
		<TrackList
			tracks={scroll.items}
			onPlay={async (trackId) => {
				await musicManager.queueRequest(
					{
						type: "addFavorites",
						userId: data.userData.id,
						filterId: filterId ?? undefined,
					},
					{ queueIndexToTrackId: trackId },
				);
			}}
		/>
	{/if}
</InfiniteScroll>
