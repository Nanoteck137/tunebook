<script lang="ts">
	import { goto } from "$app/navigation";
	import { page } from "$app/state";
	import { onMount } from "svelte";
	import { fly } from "svelte/transition";
	import { cn } from "$lib/utils";
	import { Button, Card, Separator, buttonVariants } from "$lib/components/ui";
	import {
		Play,
		Shuffle,
		X,
		ListFilter,
		Heart,
		ExternalLink,
	} from "@lucide/svelte";
	import { getMusicManager } from "$lib/music-manager.svelte";
	import { getApiClient, handleApiError } from "$lib";
	import type { Track } from "$lib/api/types";
	import InfiniteScroll from "$lib/components/InfiniteScroll.svelte";
	import { InfiniteScrollController } from "$lib/infinite-scroll.svelte";
	import TrackList from "$lib/components/track-list/TrackList.svelte";
	import Spacer from "$lib/components/Spacer.svelte";
	import SectionHeader from "$lib/components/SectionHeader.svelte";
	import DebouncedSearchInput from "$lib/components/DebouncedSearchInput.svelte";
	import { SortToggleDropdown } from "$lib/components/sort";
	import FilterButton from "../../tracks/FilterButton.svelte";
	import {
		sortTypes,
		defaultSort,
		type SortType,
		constructFilterSort,
	} from "./types";
	import SavedFilterButton from "$lib/components/SavedFilterButton.svelte";
	import SavedFilterCard from "$lib/components/SavedFilterCard.svelte";

	let { data } = $props();

	const musicManager = getMusicManager();
	const apiClient = getApiClient();
	const user = $derived(data.user!);

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

	function clearFilter() {
		const query = page.url.searchParams;
		query.delete("filterId");
		goto("?" + query.toString(), {
			invalidateAll: true,
			replaceState: true,
		});
	}

	async function playAll() {
		await musicManager.queueRequest({
			type: "addFavorites",
			userId: user.id,
			filterId: filterId ?? undefined,
		});
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

			const res = await apiClient.getUserTrackFavoritesById(user.id, {
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
</script>

<div class="flex flex-col gap-4">
	<SectionHeader count={data.page?.totalItems ?? 0}>
		<Heart />
		Favorites

		{#snippet actions()}
			<div class="flex items-center gap-2">
				<Button size="sm" onclick={() => playAll()}>
					<Play />
					Play
				</Button>
				<Button
					size="icon-sm"
					variant="ghost"
					onclick={async () => {
						await musicManager.queueRequest(
							{
								type: "addFavorites",
								userId: user.id,
								filterId: filterId ?? undefined,
							},
							{ shuffle: true },
						);
					}}
				>
					<Shuffle />
				</Button>
			</div>
		{/snippet}
	</SectionHeader>

	<!-- Toolbar -->
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

<InfiniteScroll controller={scroll}>
	<TrackList
		tracks={scroll.items}
		onPlay={async (trackId) => {
			await musicManager.queueRequest(
				{
					type: "addFavorites",
					userId: user.id,
					filterId: filterId ?? undefined,
				},
				{ queueIndexToTrackId: trackId },
			);
		}}
	/>
</InfiniteScroll>
