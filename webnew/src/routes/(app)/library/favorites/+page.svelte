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
			<Button
				variant="ghost"
				size="icon"
				title="Saved Filters"
				aria-label="Saved Filters"
				class={cn(
					"relative transition-opacity",
					filterOpen && "bg-accent text-accent-foreground",
				)}
				onclick={() => (filterOpen = !filterOpen)}
			>
				<ListFilter />
				{#if data.filters && data.filters.length > 0}
					<span
						class="absolute top-1.5 right-1.5 h-2 w-2 rounded-full bg-primary"
					></span>
				{/if}
			</Button>

			<SortToggleDropdown
				types={sortTypes}
				{sort}
				{defaultSort}
				onSortChange={(value) => updateSort(value)}
			/>
		</div>
	</div>

	{#if filterOpen}
		<div transition:fly={{ y: -6, duration: 150 }}>
			<Card.Root class="py-2">
				<Card.Content class="flex-warp flex items-center justify-between px-2">
					<div class="flex flex-wrap items-center gap-1.5">
						{#if data.filters && data.filters.length > 0}
							{#each data.filters as filter (filter.filterId)}
								<FilterButton {filter} />
							{/each}
						{:else}
							<span class="text-sm text-muted-foreground">None saved yet</span>
						{/if}
					</div>

					<div class="flex items-center gap-1">
						<a
							href="/library/filters/tracks"
							class={buttonVariants({ variant: "ghost", size: "icon" })}
							title="Manage Filters"
						>
							<ExternalLink />
						</a>

						{#if filterId}
							<Button variant="ghost" size="icon" onclick={clearFilter}>
								<X />
							</Button>
						{/if}
					</div>
				</Card.Content>
			</Card.Root>
		</div>
	{/if}
</div>

<InfiniteScroll controller={scroll}>
	<TrackList
		totalTracks={data.page.totalItems}
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
