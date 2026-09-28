<script lang="ts">
	import { Tags } from "@lucide/svelte";
	import type { RankedTag } from "$lib/api/types";
	import RankedItem from "$lib/components/RankedItem.svelte";
	import SectionHeader from "$lib/components/SectionHeader.svelte";
	import { Breadcrumb } from "$lib/components/ui";
	import InfiniteScroll from "$lib/components/InfiniteScroll.svelte";
	import { InfiniteScrollController } from "$lib/infinite-scroll.svelte";
	import { getApiClient, handleApiError } from "$lib";

	let { data } = $props();
	const apiClient = getApiClient();

	function formatTagSlug(slug: string): string {
		const words = slug.split("-");
		const sentence = words.join(" ");
		return sentence.charAt(0).toUpperCase() + sentence.slice(1);
	}

	const scroll = new InfiniteScrollController<RankedTag>({
		initialLoad: () => ({
			items: data.tags,
			hasMore: data.page.page + 1 < data.page.totalPages,
			page: data.page.page,
		}),
		load: async (nextPage) => {
			const res = await apiClient.getUserTotalReviewTags(data.userData.id, {
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
				items: res.data.tags,
				hasMore: res.data.page.page + 1 < res.data.page.totalPages,
			};
		},
		itemKey: (tag) => tag.tagSlug,
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
				<Breadcrumb.Page>Top Tags</Breadcrumb.Page>
			</Breadcrumb.Item>
		</Breadcrumb.List>
	</Breadcrumb.Root>

	<SectionHeader count={data.page.totalItems}>
		<Tags />
		Top Tags
	</SectionHeader>

	<p class="text-sm text-muted-foreground">
		{data.page.totalItems.toLocaleString()}
		{data.tags.length === 1 ? "tag" : "tags"} ranked by plays
	</p>

	{#if data.page.totalItems > 0}
		<InfiniteScroll controller={scroll} className="gap-2">
			{#each scroll.items as item (item.tagSlug)}
				<RankedItem
					rank={item.rank}
					name={formatTagSlug(item.tagSlug)}
					playCount={item.playCount}
				/>
			{/each}
		</InfiniteScroll>
	{:else}
		<div
			class="flex flex-col items-center gap-2 rounded-lg border py-16 text-center"
		>
			<Tags size={32} class="text-muted-foreground/40" />
			<p class="text-sm font-medium">No top tags yet</p>
			<p class="max-w-sm text-sm text-muted-foreground">
				Tags ranked by play count will appear here once your all-time review is
				generated.
			</p>
		</div>
	{/if}
</div>
