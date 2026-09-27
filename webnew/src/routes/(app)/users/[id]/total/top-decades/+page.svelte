<script lang="ts">
	import { CalendarRange } from "@lucide/svelte";
	import type { RankedDecade } from "$lib/api/types";
	import RankedItem from "../RankedItem.svelte";
	import SectionHeader from "$lib/components/SectionHeader.svelte";
	import { Breadcrumb } from "$lib/components/ui";
	import InfiniteScroll from "$lib/components/InfiniteScroll.svelte";
	import { InfiniteScrollController } from "$lib/infinite-scroll.svelte";
	import { getApiClient, handleApiError } from "$lib";

	let { data } = $props();
	const apiClient = getApiClient();

	const scroll = new InfiniteScrollController<RankedDecade>({
		initialLoad: () => ({
			items: data.decades,
			hasMore: data.page.page + 1 < data.page.totalPages,
			page: data.page.page,
		}),
		load: async (nextPage) => {
			const res = await apiClient.getUserTotalReviewDecades(data.userData.id, {
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
				items: res.data.decades,
				hasMore: res.data.page.page + 1 < res.data.page.totalPages,
			};
		},
		itemKey: (decade) => String(decade.decade),
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
				<Breadcrumb.Page>Decades</Breadcrumb.Page>
			</Breadcrumb.Item>
		</Breadcrumb.List>
	</Breadcrumb.Root>

	<SectionHeader count={data.page.totalItems}>
		<CalendarRange />
		Decades
	</SectionHeader>

	<p class="text-sm text-muted-foreground">
		{data.page.totalItems.toLocaleString()}
		{data.decades.length === 1 ? "decade" : "decades"} ranked by plays
	</p>

	{#if data.page.totalItems > 0}
		<InfiniteScroll controller={scroll} className="gap-2">
			{#each scroll.items as item (item.decade)}
				<RankedItem
					rank={item.rank}
					name="{item.decade}s"
					playCount={item.playCount}
				/>
			{/each}
		</InfiniteScroll>
	{:else}
		<div
			class="flex flex-col items-center gap-2 rounded-lg border py-16 text-center"
		>
			<CalendarRange size={32} class="text-muted-foreground/40" />
			<p class="text-sm font-medium">No decades yet</p>
			<p class="max-w-sm text-sm text-muted-foreground">
				Decades ranked by play count will appear here once your all-time review is
				generated.
			</p>
		</div>
	{/if}
</div>
