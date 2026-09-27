<script lang="ts">
	import { Users } from "@lucide/svelte";
	import type { RankedArtist } from "$lib/api/types";
	import RankedItem from "$lib/components/RankedItem.svelte";
	import SectionHeader from "$lib/components/SectionHeader.svelte";
	import { Breadcrumb } from "$lib/components/ui";
	import InfiniteScroll from "$lib/components/InfiniteScroll.svelte";
	import { InfiniteScrollController } from "$lib/infinite-scroll.svelte";
	import { getApiClient, handleApiError } from "$lib";

	let { data } = $props();
	const apiClient = getApiClient();

	const scroll = new InfiniteScrollController<RankedArtist>({
		initialLoad: () => ({
			items: data.artists,
			hasMore: data.page.page + 1 < data.page.totalPages,
			page: data.page.page,
		}),
		load: async (nextPage) => {
			const res = await apiClient.getUserTotalReviewArtists(data.userData.id, {
				query: {
					page: String(nextPage),
					perPage: String(data.page.perPage),
				},
			});
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

<div class="flex flex-col gap-4">
	<Breadcrumb.Root>
		<Breadcrumb.List>
			<Breadcrumb.Item>
				<Breadcrumb.Link href="/users/{data.userData.id}/total">
					All time
				</Breadcrumb.Link>
			</Breadcrumb.Item>
			<Breadcrumb.Separator />
			<Breadcrumb.Item>
				<Breadcrumb.Page>Top Artists</Breadcrumb.Page>
			</Breadcrumb.Item>
		</Breadcrumb.List>
	</Breadcrumb.Root>

	<SectionHeader count={data.page.totalItems}>
		<Users />
		Top Artists
	</SectionHeader>

	<p class="text-sm text-muted-foreground">
		{data.page.totalItems.toLocaleString()}
		{data.artists.length === 1 ? "artist" : "artists"} ranked by plays
	</p>

	{#if data.page.totalItems > 0}
		<InfiniteScroll controller={scroll} className="gap-2">
			{#each scroll.items as item (item.id)}
				<RankedItem
					rank={item.rank}
					name={item.name}
					playCount={item.playCount}
					href="/artists/{item.id}"
					art={item.coverArt.small}
				/>
			{/each}
		</InfiniteScroll>
	{:else}
		<div
			class="flex flex-col items-center gap-2 rounded-lg border py-16 text-center"
		>
			<Users size={32} class="text-muted-foreground/40" />
			<p class="text-sm font-medium">No top artists yet</p>
			<p class="max-w-sm text-sm text-muted-foreground">
				Artists ranked by play count will appear here once your all-time review
				is generated.
			</p>
		</div>
	{/if}
</div>
