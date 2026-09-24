<script lang="ts">
	import { goto } from "$app/navigation";
	import { page } from "$app/state";
	import { Breadcrumb, Button } from "$lib/components/ui";
	import { Play, Shuffle } from "@lucide/svelte";
	import { SortToggleDropdown } from "$lib/components/sort";
	import TrackList from "$lib/components/track-list/TrackList.svelte";
	import { getMusicManager } from "$lib/music-manager.svelte";
	import { getApiClient, handleApiError } from "$lib";
	import type { Track } from "$lib/api/types";
	import InfiniteScroll from "$lib/components/InfiniteScroll.svelte";
	import { InfiniteScrollController } from "$lib/infinite-scroll.svelte";
	import { sortTypes, defaultSort, applySort, type SortType } from "./types";

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
				filter: `artistId = "${data.artist.id}" or featuringArtists has "${data.artist.id}"`,
			};

			applySort(data.filter, query);

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
</script>

<div class="flex flex-col gap-4">
	<Breadcrumb.Root>
		<Breadcrumb.List>
			<Breadcrumb.Item>
				<Breadcrumb.Link href="/artists">Artists</Breadcrumb.Link>
			</Breadcrumb.Item>
			<Breadcrumb.Separator />
			<Breadcrumb.Item>
				<Breadcrumb.Link href="/artists/{data.artist.id}">
					{data.artist.name}
				</Breadcrumb.Link>
			</Breadcrumb.Item>
			<Breadcrumb.Separator />
			<Breadcrumb.Item>
				<Breadcrumb.Page>Tracks</Breadcrumb.Page>
			</Breadcrumb.Item>
		</Breadcrumb.List>
	</Breadcrumb.Root>

	<div
		class="flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between"
	>
		<div class="flex items-baseline gap-2">
			<h1 class="text-xl font-bold">Tracks</h1>
			{#if data.page}
				<span class="text-sm text-muted-foreground"
					>{data.page.totalItems}</span
				>
			{/if}
		</div>

		<div class="flex items-center gap-2">
			<Button
				variant="outline"
				size="sm"
				onclick={async () => {
					await musicManager.queueRequest(
						{ type: "addArtist", artistId: data.artist.id },
						{ shuffle: true },
					);
				}}
			>
				<Shuffle size={14} />
				Shuffle
			</Button>
			<Button
				size="sm"
				onclick={async () => {
					await musicManager.queueRequest(
						{ type: "addArtist", artistId: data.artist.id },
						{},
					);
				}}
			>
				<Play size={14} />
				Play All
			</Button>
		</div>

		<SortToggleDropdown types={sortTypes} {sort} onSortChange={updateSort} />
	</div>

	<InfiniteScroll controller={scroll}>
		<TrackList
			totalTracks={scroll.items.length}
			tracks={scroll.items}
			onPlay={async (trackId) => {
				await musicManager.queueRequest(
					{ type: "addArtist", artistId: data.artist.id },
					{ queueIndexToTrackId: trackId },
				);
			}}
		/>
	</InfiniteScroll>
</div>
