<script lang="ts">
	import { goto } from "$app/navigation";
	import { page } from "$app/state";
	import { onMount } from "svelte";
	import { Button, Checkbox } from "$lib/components/ui";
	import { Music, Play, Shuffle } from "@lucide/svelte";
	import { SortToggleDropdown } from "$lib/components/sort";
	import TrackList from "$lib/components/track-list/TrackList.svelte";
	import { getMusicManager } from "$lib/music-manager.svelte";
	import { getApiClient, handleApiError } from "$lib";
	import type { Track } from "$lib/api/types";
	import InfiniteScroll from "$lib/components/InfiniteScroll.svelte";
	import { InfiniteScrollController } from "$lib/infinite-scroll.svelte";
	import SectionHeader from "$lib/components/SectionHeader.svelte";
	import DebouncedSearchInput from "$lib/components/DebouncedSearchInput.svelte";
	import {
		sortTypes,
		defaultSort,
		buildArtistTracksQuery,
		type SortType,
	} from "./types";
	import Spacer from "$lib/components/Spacer.svelte";

	let { data } = $props();
	const musicManager = getMusicManager();
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

			buildArtistTracksQuery(data.filter, data.artist.id, query);

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
</script>

<div class="flex flex-col gap-4">
	<SectionHeader count={data.page.totalItems}>
		<Music />
		Tracks

		{#snippet actions()}
			<div class="flex items-center gap-2">
				<Button
					size="sm"
					onclick={async () => {
						await musicManager.queueRequest(
							{ type: "addArtist", artistId: data.artist.id },
							{},
						);
					}}
				>
					<Play />
					Play
				</Button>
				<Button
					size="icon-sm"
					variant="ghost"
					onclick={async () => {
						await musicManager.queueRequest(
							{ type: "addArtist", artistId: data.artist.id },
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
			placeholder="Search tracks..."
			{value}
			setValue={(v) => (value = v)}
			{search}
		/>

		<div class="flex items-center gap-1 pr-2">
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
			/>
		</div>
	</div>
</div>

<Spacer size="md" />

<InfiniteScroll controller={scroll}>
	<TrackList
		totalTracks={data.page.totalItems}
		tracks={scroll.items}
		bind:selectedTracks
		onPlay={async (trackId) => {
			await musicManager.queueRequest(
				{ type: "addArtist", artistId: data.artist.id },
				{ queueIndexToTrackId: trackId },
			);
		}}
	/>
</InfiniteScroll>
