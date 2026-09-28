<script lang="ts">
	import { goto, invalidateAll } from "$app/navigation";
	import { page } from "$app/state";
	import { onMount } from "svelte";
	import { Button, Separator } from "$lib/components/ui";
	import { ChevronDown, ListMusic, Plus, X } from "@lucide/svelte";
	import { cn } from "$lib/utils";
	import SectionHeader from "$lib/components/SectionHeader.svelte";
	import DebouncedSearchInput from "$lib/components/DebouncedSearchInput.svelte";
	import { SortToggleDropdown } from "$lib/components/sort";
	import { getApiClient, handleApiError } from "$lib";
	import type { Playlist } from "$lib/api/types";
	import PlaylistTile from "$lib/components/tiles/PlaylistTile.svelte";
	import TileGrid from "$lib/components/tiles/TileGrid.svelte";
	import InfiniteScroll from "$lib/components/InfiniteScroll.svelte";
	import { InfiniteScrollController } from "$lib/infinite-scroll.svelte";
	import NewPlaylistModal from "../../playlists/NewPlaylistModal.svelte";
	import {
		sortTypes,
		defaultSort,
		type SortType,
		constructFilterSort,
	} from "../../playlists/types";
	import { toast } from "svelte-sonner";

	let { data } = $props();
	const apiClient = getApiClient();

	let openNewPlaylistModal = $state(false);
	let selectedPlaylists = $state<string[]>([]);

	const scroll = new InfiniteScrollController<Playlist>({
		initialLoad: () => ({
			items: data.playlists,
			hasMore: data.page.page + 1 < data.page.totalPages,
			page: data.page.page,
		}),
		load: async (nextPage) => {
			const query: Record<string, string> = {
				page: String(nextPage),
				perPage: String(data.page.perPage),
			};

			constructFilterSort(data.filter, query, data.user?.id);

			const res = await apiClient.getPlaylists({ query });
			if (!res.success) {
				handleApiError(res.error);
				return null;
			}

			return {
				items: res.data.playlists,
				hasMore: res.data.page.page + 1 < res.data.page.totalPages,
			};
		},
		itemKey: (playlist) => playlist.id,
	});

	let sort = $state(
		(page.url.searchParams.get("sort") as SortType) ?? defaultSort,
	);

	function updateSort(value: string) {
		sort = value as SortType;

		selectedPlaylists = [];

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

	async function toggleQuick(playlistId: string) {
		const isQuick = data.user?.quickPlaylist === playlistId;

		const res = await apiClient.setQuickPlaylist({
			playlistId: isQuick ? "" : playlistId,
		});
		if (!res.success) {
			handleApiError(res.error);
			return;
		}

		await invalidateAll();
	}

	async function handleReorder(anchorPlaylistId: string | null) {
		const res = await apiClient.reorderPlaylists({
			before: false,
			anchorPlaylistId: anchorPlaylistId ?? "",
			playlistIds: selectedPlaylists,
		});
		if (!res.success) {
			return handleApiError(res.error);
		}

		selectedPlaylists = [];
		toast.success("Updated playlists");

		if (sort !== "position") {
			updateSort("position");
		} else {
			invalidateAll();
		}
	}
</script>

<div class="flex flex-col gap-4">
	<SectionHeader count={data.page?.totalItems}>
		<ListMusic />
		Playlists

		{#snippet actions()}
			<Button size="sm" onclick={() => (openNewPlaylistModal = true)}>
				<Plus size={14} />
				New Playlist
			</Button>
		{/snippet}
	</SectionHeader>

	<!-- Toolbar -->
	<div class="flex flex-wrap items-center justify-between gap-2">
		<DebouncedSearchInput
			class="flex-1 md:max-w-64"
			placeholder="Search playlists..."
			{value}
			setValue={(v) => (value = v)}
			{search}
		/>

		<div class="flex items-center gap-1">
			<SortToggleDropdown
				types={sortTypes}
				{sort}
				{defaultSort}
				onSortChange={(value) => updateSort(value)}
			/>
		</div>
	</div>

	<InfiniteScroll controller={scroll}>
		<TileGrid>
			{#if selectedPlaylists.length > 0}
				<div class="group relative flex shrink-0 flex-col">
					<button
						class="flex aspect-square w-full items-center justify-center overflow-hidden rounded-lg border-2 border-dashed border-border/60 bg-muted/40 text-muted-foreground transition-colors group-hover:border-primary/60 group-hover:bg-accent/50 group-hover:text-foreground"
						onclick={() => handleReorder(null)}
						aria-label={`Move ${selectedPlaylists.length} selected playlists to the beginning`}
					>
						<div class="flex flex-col items-center gap-1.5">
							<ChevronDown
								class="h-5 w-5 text-muted-foreground/60 transition-colors group-hover:text-foreground"
							/>
							<span class="px-2 text-xs font-medium">
								Place {selectedPlaylists.length} here
							</span>
						</div>
					</button>

					<button
						class="absolute top-1.5 left-1.5 flex h-7 w-7 items-center justify-center rounded-full border bg-background/70 text-muted-foreground backdrop-blur-sm transition-colors hover:text-foreground"
						onclick={() => (selectedPlaylists = [])}
						aria-label="Clear selection"
						title="Clear selection"
					>
						<X size={14} />
					</button>

					<div class="flex flex-col gap-0.5 pt-2">
						<span
							class="truncate text-sm font-medium text-muted-foreground transition-colors group-hover:text-foreground"
						>
							Move to beginning
						</span>
					</div>
				</div>
			{/if}

			{#each scroll.items as playlist (playlist.id)}
				<PlaylistTile
					{playlist}
					selectionMode={selectedPlaylists.length > 0}
					selected={selectedPlaylists.includes(playlist.id)}
					onSelectChange={(checked) => {
						if (checked) {
							selectedPlaylists = [...selectedPlaylists, playlist.id];
						} else {
							selectedPlaylists = selectedPlaylists.filter(
								(id) => playlist.id !== id,
							);
						}
					}}
					onMoveAfter={() => handleReorder(playlist.id)}
					isQuick={data.user?.quickPlaylist === playlist.id}
					onToggleQuick={() => toggleQuick(playlist.id)}
				/>
			{/each}
		</TileGrid>
	</InfiniteScroll>
</div>

<NewPlaylistModal bind:open={openNewPlaylistModal} />
