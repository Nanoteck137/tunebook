<script lang="ts">
	import { goto } from "$app/navigation";
	import { page } from "$app/state";
	import {
		Button,
		Separator,
		Input,
		DropdownMenu,
		buttonVariants,
	} from "$lib/components/ui";
	import {
		Play,
		Shuffle,
		X,
		ListSortAscendingIcon,
		CheckIcon,
		ListFilter,
		Heart,
	} from "@lucide/svelte";
	import { getMusicManager } from "$lib/music-manager.svelte";
	import { getApiClient, handleApiError } from "$lib";
	import type { Track } from "$lib/api/types";
	import InfiniteScroll from "$lib/components/InfiniteScroll.svelte";
	import { InfiniteScrollController } from "$lib/infinite-scroll.svelte";
	import TrackList from "$lib/components/track-list/TrackList.svelte";
	import Spacer from "$lib/components/Spacer.svelte";
	import SectionHeader from "$lib/components/SectionHeader.svelte";
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

	let searchQuery = $state(page.url.searchParams.get("query") ?? "");
	function updateSearch() {
		const query = page.url.searchParams;
		query.delete("query");

		if (searchQuery) {
			query.set("query", searchQuery);
		}

		goto("?" + query.toString(), { invalidateAll: true });
	}

	function clearSearch() {
		searchQuery = "";
		updateSearch();
	}

	let sort = $state(
		(page.url.searchParams.get("sort") as SortType) ?? defaultSort,
	);
	function updateSort(value: string) {
		sort = value as SortType;

		const query = page.url.searchParams;
		query.delete("sort");
		query.set("sort", sort);

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
	<section>
		<SectionHeader count={data.page?.totalItems ?? 0}>
			<Heart />
			Favorites

			{#snippet actions()}
				<div class="flex items-center gap-2">
					<Button size="sm" onclick={() => playAll()}>
						<Play />
						Play All
					</Button>
					<Button
						size="sm"
						variant="outline"
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
						Shuffle
					</Button>
				</div>
			{/snippet}
		</SectionHeader>
	</section>

	<!-- Toolbar -->
	<div class="flex flex-wrap items-center justify-between gap-2 px-2">
		<div class="relative flex-1 md:max-w-64">
			<Input
				class="pr-8"
				placeholder="Search favorites..."
				bind:value={searchQuery}
				onkeydown={(e) => {
					if (e.key === "Enter") {
						updateSearch();
					}
				}}
			/>
			{#if searchQuery}
				<button
					class="absolute top-1/2 right-1.5 -translate-y-1/2 rounded-full p-0.5 text-muted-foreground hover:text-foreground"
					onclick={clearSearch}
					aria-label="Clear search"
				>
					<X size={14} />
				</button>
			{/if}
		</div>

		<div class="flex items-center gap-1">
			<DropdownMenu.Root>
				<DropdownMenu.Trigger
					class={buttonVariants({ variant: "ghost", size: "icon" })}
					title="Sort"
					aria-label="Sort"
				>
					<ListSortAscendingIcon />
				</DropdownMenu.Trigger>
				<DropdownMenu.Content align="end">
					<DropdownMenu.Group>
						{#each sortTypes as ty (ty.value)}
							{@const selected = sort === ty.value}
							<DropdownMenu.Item
								onSelect={() => updateSort(ty.value)}
								class={selected ? "bg-accent text-foreground" : ""}
							>
								{#if selected}
									<CheckIcon />
								{/if}
								{ty.label}
							</DropdownMenu.Item>
						{/each}
					</DropdownMenu.Group>
				</DropdownMenu.Content>
			</DropdownMenu.Root>
		</div>
	</div>

	<Separator />
</div>

<div
	class="flex flex-wrap items-center justify-between gap-2 rounded-lg border bg-muted/40 px-3 py-2"
>
	<div class="flex flex-wrap items-center gap-1.5">
		<span
			class="mr-1.5 flex items-center gap-1.5 text-xs font-medium text-muted-foreground"
		>
			<ListFilter size={12} />
			Saved Filters
		</span>

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
			class={buttonVariants({ variant: "ghost", size: "sm" })}
		>
			<ListFilter size={14} />
			Manage Filters
		</a>

		{#if filterId}
			<Button variant="ghost" size="sm" onclick={clearFilter}>
				<X size={14} />
				Clear
			</Button>
		{/if}
	</div>
</div>

<Spacer size="lg" />

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
