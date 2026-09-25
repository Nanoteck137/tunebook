<script lang="ts">
	import { goto } from "$app/navigation";
	import { page } from "$app/state";
	import { EllipsisVertical, Play, Shuffle, X } from "@lucide/svelte";
	import { Button, buttonVariants, DropdownMenu } from "$lib/components/ui";
	import { getMusicManager } from "$lib/music-manager.svelte";
	import { getApiClient, handleApiError } from "$lib";
	import type { Track } from "$lib/api/types";
	import InfiniteScroll from "$lib/components/InfiniteScroll.svelte";
	import { InfiniteScrollController } from "$lib/infinite-scroll.svelte";
	import TrackList from "$lib/components/track-list/TrackList.svelte";
	import Spacer from "$lib/components/Spacer.svelte";
	// import FilterButton from "../../../tracks/FilterButton.svelte";

	let { data } = $props();

	const musicManager = getMusicManager();
	const apiClient = getApiClient();

	let filterId = $derived(page.url.searchParams.get("filterId"));

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
			userId: data.userData.id,
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

			if (filterId) {
				query["filterId"] = filterId;
			}

			const res = await apiClient.getUserTrackFavoritesById(data.userData.id, {
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

<div class="flex flex-col gap-6">
	<div
		class="flex flex-col gap-6 rounded-lg border bg-linear-to-b from-[oklch(0.93_0.045_75)] to-background p-4 shadow-sm sm:p-6 md:flex-row md:items-end md:gap-8 dark:from-[oklch(0.24_0.03_80)] dark:to-background"
	>
		<div class="flex min-w-0 flex-col gap-2">
			<p
				class="text-xs font-semibold tracking-wider text-muted-foreground uppercase"
			>
				Favorites
			</p>

			<h1 class="line-clamp-2 text-2xl font-bold md:text-4xl">
				{data.userData.displayName}
			</h1>

			<p class="text-sm text-muted-foreground">
				{data.page.totalItems}
				{data.page.totalItems === 1 ? "track" : "tracks"}
			</p>

			<div class="flex gap-2 pt-2">
				<Button onclick={() => playAll()}>
					<Play />
					Play
				</Button>
				<Button
					variant="ghost"
					size="icon"
					onclick={async () => {
						await musicManager.queueRequest(
							{
								type: "addFavorites",
								userId: data.userData.id,
								filterId: filterId ?? undefined,
							},
							{ shuffle: true },
						);
					}}
				>
					<Shuffle />
				</Button>
				{#if filterId}
					<DropdownMenu.Root>
						<DropdownMenu.Trigger
							class={buttonVariants({ variant: "ghost", size: "icon" })}
						>
							<EllipsisVertical />
						</DropdownMenu.Trigger>
						<DropdownMenu.Content align="start">
							<DropdownMenu.Group>
								<DropdownMenu.Item onSelect={clearFilter}>
									<X />
									Clear filter
								</DropdownMenu.Item>
							</DropdownMenu.Group>
						</DropdownMenu.Content>
					</DropdownMenu.Root>
				{/if}
			</div>
		</div>
	</div>

	{#if data.filters && data.filters.length > 0}
		<div class="flex flex-wrap items-center gap-2">
			<!-- {#each data.filters as filter (filter.filterId)} -->
			<!--   <FilterButton {filter} /> -->
			<!-- {/each} -->
		</div>
	{/if}
</div>

<Spacer size="lg" />

<InfiniteScroll controller={scroll}>
	<TrackList
		tracks={scroll.items}
		onPlay={async (trackId) => {
			await musicManager.queueRequest(
				{
					type: "addFavorites",
					userId: data.userData.id,
					filterId: filterId ?? undefined,
				},
				{ queueIndexToTrackId: trackId },
			);
		}}
	/>
</InfiniteScroll>
