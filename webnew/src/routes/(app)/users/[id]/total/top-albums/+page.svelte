<script lang="ts">
	import { DiscAlbum } from "@lucide/svelte";
	import type { RankedAlbum } from "$lib/api/types";
	import RankedItem from "../RankedItem.svelte";
	import SectionHeader from "$lib/components/SectionHeader.svelte";
	import { Breadcrumb } from "$lib/components/ui";
	import InfiniteScroll from "$lib/components/InfiniteScroll.svelte";
	import { InfiniteScrollController } from "$lib/infinite-scroll.svelte";
	import { getApiClient, handleApiError } from "$lib";

	let { data } = $props();
	const apiClient = getApiClient();

	function albumArtistNames(album: RankedAlbum): string {
		return album.artists.map((a) => a.name).join(", ");
	}

	const scroll = new InfiniteScrollController<RankedAlbum>({
		initialLoad: () => ({
			items: data.albums,
			hasMore: data.page.page + 1 < data.page.totalPages,
			page: data.page.page,
		}),
		load: async (nextPage) => {
			const res = await apiClient.getUserTotalReviewAlbums(data.userData.id, {
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
				items: res.data.albums,
				hasMore: res.data.page.page + 1 < res.data.page.totalPages,
			};
		},
		itemKey: (album) => album.id,
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
				<Breadcrumb.Page>Top Albums</Breadcrumb.Page>
			</Breadcrumb.Item>
		</Breadcrumb.List>
	</Breadcrumb.Root>

	<SectionHeader count={data.page.totalItems}>
		<DiscAlbum />
		Top Albums
	</SectionHeader>

	<p class="text-sm text-muted-foreground">
		{data.page.totalItems.toLocaleString()}
		{data.albums.length === 1 ? "album" : "albums"} ranked by plays
	</p>

	{#if data.page.totalItems > 0}
		<InfiniteScroll controller={scroll} className="gap-2">
			{#each scroll.items as item (item.id)}
				<RankedItem
					rank={item.rank}
					name={item.name}
					playCount={item.playCount}
					href="/albums/{item.id}"
					subtitle={albumArtistNames(item)}
					art={item.coverArt.small}
				/>
			{/each}
		</InfiniteScroll>
	{:else}
		<div
			class="flex flex-col items-center gap-2 rounded-lg border py-16 text-center"
		>
			<DiscAlbum size={32} class="text-muted-foreground/40" />
			<p class="text-sm font-medium">No top albums yet</p>
			<p class="max-w-sm text-sm text-muted-foreground">
				Albums ranked by play count will appear here once your all-time review is
				generated.
			</p>
		</div>
	{/if}
</div>
