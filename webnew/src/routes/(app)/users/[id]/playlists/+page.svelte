<script lang="ts">
	import { goto, invalidateAll } from "$app/navigation";
	import { page } from "$app/state";
	import { onMount } from "svelte";
	import { ChevronDown, ListMusic, Plus, X } from "@lucide/svelte";
	import { Button } from "$lib/components/ui";
	import { getApiClient, handleApiError } from "$lib";
	import type { Playlist } from "$lib/api/types";
	import SectionHeader from "$lib/components/SectionHeader.svelte";
	import DebouncedSearchInput from "$lib/components/DebouncedSearchInput.svelte";
	import { SortToggleDropdown } from "$lib/components/sort";
	import PlaylistTile from "$lib/components/tiles/PlaylistTile.svelte";
	import TileGrid from "$lib/components/tiles/TileGrid.svelte";
	import InfiniteScroll from "$lib/components/InfiniteScroll.svelte";
	import { InfiniteScrollController } from "$lib/infinite-scroll.svelte";
	import NewPlaylistModal from "../../../playlists/NewPlaylistModal.svelte";
	import {
		sortTypes,
		defaultSort,
		type SortType,
		constructFilterSort,
	} from "../../../playlists/types";
	import { toast } from "svelte-sonner";
	import Spacer from "$lib/components/Spacer.svelte";

	let { data } = $props();
	const apiClient = getApiClient();

	let isOwner = $derived(data.userData.id === data.user?.id);

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

			constructFilterSort(data.filter, query, data.userData.id);

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

	function toggleSelect(playlistId: string, checked: boolean) {
		if (checked) {
			selectedPlaylists = [...selectedPlaylists, playlistId];
		} else {
			selectedPlaylists = selectedPlaylists.filter((id) => id !== playlistId);
		}
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
	<SectionHeader count={data.page.totalItems}>
		<ListMusic />
		Playlists

		{#snippet actions()}
			{#if isOwner}
				<Button size="sm" onclick={() => (openNewPlaylistModal = true)}>
					<Plus size={14} />
					New Playlist
				</Button>
			{/if}
		{/snippet}
	</SectionHeader>

	<div class="flex flex-wrap items-center justify-between gap-2">
		<DebouncedSearchInput
			class="flex-1 md:max-w-64"
			placeholder="Search playlists..."
			{value}
			setValue={(v) => (value = v)}
			{search}
		/>

		<SortToggleDropdown
			types={sortTypes}
			{sort}
			{defaultSort}
			onSortChange={(value) => updateSort(value)}
		/>
	</div>
</div>

<Spacer size="md" />

<InfiniteScroll controller={scroll}>
	{#if scroll.items.length === 0}
		<div class="flex flex-col items-center gap-2 rounded-lg border py-16">
			<ListMusic size={32} class="text-muted-foreground/40" />
			<p class="text-sm text-muted-foreground">
				{data.filter.query
					? "No playlists match your search"
					: isOwner
						? "You haven't created any playlists yet"
						: "No playlists yet"}
			</p>
		</div>
	{:else}
		<TileGrid class="xl:grid-cols-5 2xl:grid-cols-5">
			{#if isOwner && selectedPlaylists.length > 0}
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
					selectionMode={isOwner && selectedPlaylists.length > 0}
					selected={selectedPlaylists.includes(playlist.id)}
					onSelectChange={isOwner
						? (checked) => toggleSelect(playlist.id, checked)
						: undefined}
					onMoveAfter={isOwner ? () => handleReorder(playlist.id) : undefined}
					isQuick={isOwner && data.user?.quickPlaylist === playlist.id}
					onToggleQuick={isOwner ? () => toggleQuick(playlist.id) : undefined}
				/>
			{/each}
		</TileGrid>
	{/if}
</InfiniteScroll>

{#if isOwner}
	<NewPlaylistModal bind:open={openNewPlaylistModal} />
{/if}
