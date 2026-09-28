<script lang="ts">
	import { Music } from "@lucide/svelte";
	import type { RankedTrack } from "$lib/api/types";
	import RankedItem from "$lib/components/RankedItem.svelte";
	import SectionHeader from "$lib/components/SectionHeader.svelte";
	import { Breadcrumb } from "$lib/components/ui";
	import InfiniteScroll from "$lib/components/InfiniteScroll.svelte";
	import { InfiniteScrollController } from "$lib/infinite-scroll.svelte";
	import { getApiClient, handleApiError } from "$lib";

	let { data } = $props();
	const apiClient = getApiClient();

	const monthNames = [
		"January",
		"February",
		"March",
		"April",
		"May",
		"June",
		"July",
		"August",
		"September",
		"October",
		"November",
		"December",
	];

	let monthName = $derived(monthNames[data.month - 1]);

	const scroll = new InfiniteScrollController<RankedTrack>({
		initialLoad: () => ({
			items: data.tracks,
			hasMore: data.page.page + 1 < data.page.totalPages,
			page: data.page.page,
		}),
		load: async (nextPage) => {
			const res = await apiClient.getUserYearReviewMonthTracks(
				data.userData.id,
				String(data.year),
				String(data.month),
				{
					query: {
						page: String(nextPage),
						perPage: String(data.page.perPage),
					},
				},
			);
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
				<Breadcrumb.Link href="/users/{data.userData.id}/review">
					Year in Review
				</Breadcrumb.Link>
			</Breadcrumb.Item>
			<Breadcrumb.Separator />
			<Breadcrumb.Item>
				<Breadcrumb.Link href="/users/{data.userData.id}/review/{data.year}">
					{data.year}
				</Breadcrumb.Link>
			</Breadcrumb.Item>
			<Breadcrumb.Separator />
			<Breadcrumb.Item>
				<Breadcrumb.Link
					href="/users/{data.userData
						.id}/review/{data.year}/months/{data.month}"
				>
					{monthName}
				</Breadcrumb.Link>
			</Breadcrumb.Item>
			<Breadcrumb.Separator />
			<Breadcrumb.Item>
				<Breadcrumb.Page>Top Tracks</Breadcrumb.Page>
			</Breadcrumb.Item>
		</Breadcrumb.List>
	</Breadcrumb.Root>

	<SectionHeader count={data.page.totalItems}>
		<Music />
		Top Tracks
	</SectionHeader>

	<p class="text-sm text-muted-foreground">
		{data.page.totalItems.toLocaleString()}
		{data.tracks.length === 1 ? "track" : "tracks"} ranked by plays
	</p>

	{#if data.page.totalItems > 0}
		<InfiniteScroll controller={scroll} className="gap-2">
			{#each scroll.items as item (item.id)}
				<RankedItem
					rank={item.rank}
					name={item.name}
					playCount={item.playCount}
					href="/albums/{item.albumId}"
					subtitle={item.artists.map((a) => a.name).join(", ")}
					art={item.coverArt.small}
				/>
			{/each}
		</InfiniteScroll>
	{:else}
		<div
			class="flex flex-col items-center gap-2 rounded-lg border py-16 text-center"
		>
			<Music size={32} class="text-muted-foreground/40" />
			<p class="text-sm font-medium">No top tracks yet</p>
			<p class="max-w-sm text-sm text-muted-foreground">
				tracks ranked by play count will appear here once your monthly review
				is generated.
			</p>
		</div>
	{/if}
</div>
