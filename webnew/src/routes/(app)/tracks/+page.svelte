<script lang="ts">
	import { goto } from "$app/navigation";
	import { page } from "$app/state";
	import { onMount } from "svelte";
	import { getApiClient, handleApiError } from "$lib";
	import type { Track } from "$lib/api/types";
	import { Button, Checkbox } from "$lib/components/ui";
	import {
		SortableHeader,
		SortToggleDropdown,
		type Column,
	} from "$lib/components/sort";
	import { Music, Play, Shuffle } from "@lucide/svelte";
	import HeroIcon from "$lib/components/HeroIcon.svelte";
	import HeroCard from "$lib/components/HeroCard.svelte";
	import TrackList from "$lib/components/track-list/TrackList.svelte";
	import { getMusicManager } from "$lib/music-manager.svelte";
	import InfiniteScroll from "$lib/components/InfiniteScroll.svelte";
	import { InfiniteScrollController } from "$lib/infinite-scroll.svelte";
	import DebouncedSearchInput from "$lib/components/DebouncedSearchInput.svelte";
	import SavedFilterButton from "$lib/components/SavedFilterButton.svelte";
	import SavedFilterCard from "$lib/components/SavedFilterCard.svelte";
	import {
		sortTypes,
		defaultSort,
		constructFilterSort,
		type SortType,
	} from "./types";
	import Spacer from "$lib/components/Spacer.svelte";

	let { data } = $props();
	const musicManager = getMusicManager();
	const apiClient = getApiClient();

	let filterId = $derived(page.url.searchParams.get("filterId"));

	let filterOpen = $state(false);

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

			const res = await apiClient.getTracks({ query });
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

	let selectedTracks = $state<string[]>([]);

	function toggleSelectAll() {
		if (selectedTracks.length > 0) {
			selectedTracks = [];
			return;
		}

		selectedTracks = scroll.items.map((track) => track.id);
	}

	async function playTracks(options: { shuffle?: boolean } = {}) {
		if (filterId) {
			await musicManager.queueRequest(
				{ type: "addFilter", filterId },
				options,
			);
		} else {
			await musicManager.addTracks({
				trackIds: scroll.items.map((t) => t.id),
				clear: true,
			});
		}
	}

	const columns: Column[] = [
		{
			label: "Title",
			asc: "name-a-z",
			desc: "name-z-a",
			className:
				"-ml-1 flex min-w-0 flex-1 items-center gap-1 rounded-sm px-1 py-0.5 transition-colors hover:text-foreground",
		},
		{
			label: "Album",
			asc: "album",
			desc: "album-desc",
			className:
				"-ml-1 hidden shrink-0 items-center gap-1 rounded-sm px-1 py-0.5 transition-colors hover:text-foreground md:flex",
		},
		{
			label: "Duration",
			asc: "duration",
			desc: "duration-desc",
			className:
				"hidden shrink-0 items-center gap-1 rounded-sm px-1 py-0.5 transition-colors hover:text-foreground md:flex",
		},
		{
			label: "Added",
			asc: "created-new",
			desc: "created-old",
			className:
				"-ml-1 hidden shrink-0 items-center gap-1 rounded-sm px-1 py-0.5 transition-colors hover:text-foreground lg:flex",
		},
	];
</script>

<div class="flex flex-col gap-4">
	<HeroCard
		class="section-tracks"
		innerClass="sm:flex-row sm:items-center sm:justify-between"
	>
		<div class="flex items-center gap-4">
			<HeroIcon>
				<Music />
			</HeroIcon>
			<div class="flex min-w-0 flex-col">
				<h1 class="text-2xl font-bold">Tracks</h1>
				<p class="text-sm text-muted-foreground">
					Every track in your library
					{#if data.page}
						&middot; {data.page.totalItems}
					{/if}
				</p>
			</div>
		</div>

		<div class="flex gap-2">
			<Button size="sm" onclick={() => playTracks()}>
				<Play />
				Play
			</Button>
			<Button
				size="icon-sm"
				variant="ghost"
				onclick={() => playTracks({ shuffle: true })}
			>
				<Shuffle />
			</Button>
		</div>
	</HeroCard>

	<div class="flex items-center justify-between gap-2 sm:hidden">
		<DebouncedSearchInput
			class="flex-1"
			placeholder="Search tracks..."
			{value}
			setValue={(v) => (value = v)}
			{search}
		/>

		<div class="flex items-center gap-2 pr-2">
			<SavedFilterButton
				bind:filterOpen
				hasFilters={data.filters && data.filters.length > 0}
			/>

			<SortToggleDropdown
				types={sortTypes}
				{sort}
				{defaultSort}
				onSortChange={updateSort}
			/>

			<Checkbox
				title="Select all"
				aria-label="Select all"
				checked={selectedTracks.length > 0}
				onCheckedChange={() => toggleSelectAll()}
			></Checkbox>
		</div>
	</div>

	<SortableHeader {sort} onSortChange={updateSort} {columns}>
		<div
			class="flex shrink-0 items-center gap-2 border-l border-border/40 pl-3"
		>
			<SavedFilterButton
				bind:filterOpen
				hasFilters={data.filters && data.filters.length > 0}
			/>

			<DebouncedSearchInput
				inputClass="h-7 w-44 pr-6"
				iconSize={12}
				placeholder="Search tracks..."
				{value}
				setValue={(v) => (value = v)}
				{search}
			/>

			<Checkbox
				title="Select all"
				aria-label="Select all"
				checked={selectedTracks.length > 0}
				onCheckedChange={() => toggleSelectAll()}
			/>
		</div>
	</SortableHeader>

	<SavedFilterCard {filterOpen} filters={data.filters} />
</div>

<Spacer size="md" />

<InfiniteScroll controller={scroll}>
	<TrackList
		tracks={scroll.items}
		bind:selectedTracks
		onPlay={async (trackId) => {
			if (filterId) {
				await musicManager.queueRequest(
					{ type: "addFilter", filterId },
					{ queueIndexToTrackId: trackId },
				);
			} else {
				await musicManager.addTracks({
					trackIds: scroll.items.map((t) => t.id),
					trackId,
					clear: true,
				});
			}
		}}
	/>
</InfiniteScroll>
