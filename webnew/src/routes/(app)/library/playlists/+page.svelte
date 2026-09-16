<script lang="ts">
	import { goto, invalidateAll } from "$app/navigation";
	import { page } from "$app/state";
	import {
		Button,
		buttonVariants,
		Checkbox,
		DropdownMenu,
		Input,
		Separator,
	} from "$lib/components/ui";
	import {
		CheckIcon,
		ChevronDown,
		EllipsisVertical,
		FileHeart,
		ListSortAscendingIcon,
		Play,
		Plus,
		Shuffle,
		Star,
		X,
	} from "@lucide/svelte";
	import { cn, formatPlayTime } from "$lib/utils";
	import { getApiClient, handleApiError } from "$lib";
	import type { Playlist } from "$lib/api/types";
	import InfiniteScroll from "$lib/components/InfiniteScroll.svelte";
	import { InfiniteScrollController } from "$lib/infinite-scroll.svelte";
	import { getMusicManager } from "$lib/music-manager.svelte";
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
	const musicManager = getMusicManager();

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
		query.set("sort", sort);

		goto("?" + query.toString(), { invalidateAll: true });
	}

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

	async function playPlaylist(playlistId: string, shuffle = false) {
		await musicManager.queueRequest(
			{ type: "addPlaylist", playlistId },
			{ shuffle },
		);
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

		if (sort !== "custom") {
			updateSort("custom");
		} else {
			invalidateAll();
		}
	}
</script>

<div class="flex flex-col gap-4">
	<div class="flex items-baseline justify-between gap-2 px-2">
		<div class="flex items-baseline gap-2">
			<h1 class="text-xl font-bold">Playlists</h1>
			{#if data.page}
				<span class="text-sm text-muted-foreground">
					{data.page.totalItems}
				</span>
			{/if}
		</div>

		<Button size="sm" onclick={() => (openNewPlaylistModal = true)}>
			<Plus size={14} />
			New Playlist
		</Button>
	</div>

	<!-- Toolbar -->
	<div class="flex flex-wrap items-center justify-between gap-2 px-2">
		{#if selectedPlaylists.length > 0}
			<div class="flex items-center gap-2">
				<Button
					class="rounded-full"
					variant="ghost"
					size="icon-lg"
					onclick={() => {
						selectedPlaylists = [];
					}}
				>
					<X />
				</Button>

				<Button
					class="rounded-full"
					variant="default"
					onclick={() => handleReorder(null)}
				>
					<ChevronDown />
					Insert after
				</Button>
			</div>
		{:else}
			<div class="relative flex-1 md:max-w-64">
				<Input
					class="pr-8"
					placeholder="Search playlists..."
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
		{/if}

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

	<InfiniteScroll controller={scroll}>
		<div
			class="grid grid-cols-2 gap-3 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-5 xl:grid-cols-6 2xl:grid-cols-7"
		>
			{#each scroll.items as playlist (playlist.id)}
				{@const isQuick = data.user?.quickPlaylist === playlist.id}
				<div class="group relative flex flex-col">
					<div class="relative">
						<a
							href="/playlists/{playlist.id}"
							class="block overflow-hidden rounded-lg"
						>
							<img
								src={playlist.coverArt.medium}
								alt={playlist.name}
								class="aspect-square w-full object-cover transition-transform duration-300 group-hover:scale-105"
							/>
						</a>

						{#if selectedPlaylists.length > 0}
							<div class="absolute top-1.5 left-1.5 z-10">
								<Checkbox
									checked={selectedPlaylists.includes(playlist.id)}
									onCheckedChange={(checked) => {
										if (checked) {
											selectedPlaylists = [...selectedPlaylists, playlist.id];
										} else {
											selectedPlaylists = selectedPlaylists.filter(
												(id) => playlist.id !== id,
											);
										}
									}}
								/>
							</div>

							<Button
								class="absolute right-2 bottom-2 z-10 hidden h-10 w-10 items-center justify-center rounded-full opacity-0 shadow-lg transition-all group-hover:scale-105 group-hover:opacity-100 hover:scale-110 sm:flex"
								variant="default"
								onclick={() => handleReorder(playlist.id)}
								title="Move selected after this playlist"
								aria-label={`Move selected after ${playlist.name}`}
							>
								<ChevronDown size={18} />
							</Button>
						{:else}
							<button
								class="absolute top-1.5 right-1.5 z-10 flex h-7 w-7 items-center justify-center rounded-full transition-all {isQuick
									? 'bg-primary text-primary-foreground shadow-md'
									: 'border bg-background/70 text-muted-foreground backdrop-blur-sm hover:scale-105 hover:text-foreground'}"
								title={isQuick
									? "Unset as quick playlist"
									: "Set as quick playlist"}
								aria-label={isQuick
									? `Unset ${playlist.name} as quick playlist`
									: `Set ${playlist.name} as quick playlist`}
								onclick={() => toggleQuick(playlist.id)}
							>
								<Star size={14} class={isQuick ? "fill-current" : ""} />
							</button>

							<button
								class="absolute right-2 bottom-2 hidden h-10 w-10 items-center justify-center rounded-full bg-primary text-primary-foreground opacity-0 shadow-lg transition-all group-hover:scale-105 group-hover:opacity-100 hover:scale-110 sm:flex"
								title="Play playlist"
								aria-label={`Play ${playlist.name}`}
								onclick={() => playPlaylist(playlist.id)}
							>
								<Play size={18} />
							</button>
						{/if}
					</div>

					<div class="flex flex-col gap-0.5 pt-2">
						<div class="flex items-center gap-1">
							<a
								href="/playlists/{playlist.id}"
								class="min-w-0 flex-1 truncate text-sm font-medium hover:underline"
								title={playlist.name}
							>
								{playlist.name}
							</a>

							{#if selectedPlaylists.length <= 0}
								<DropdownMenu.Root>
									<DropdownMenu.Trigger
										class={cn(
											buttonVariants({ variant: "ghost", size: "icon-sm" }),
											"-mr-1 shrink-0 rounded-full text-muted-foreground",
										)}
										aria-label={`More options for ${playlist.name}`}
									>
										<EllipsisVertical size={14} />
									</DropdownMenu.Trigger>
									<DropdownMenu.Content align="end">
										<DropdownMenu.Group>
											<DropdownMenu.Item
												onSelect={() => playPlaylist(playlist.id)}
											>
												<Play size={14} />
												Play
											</DropdownMenu.Item>
											<DropdownMenu.Item
												onSelect={() => playPlaylist(playlist.id, true)}
											>
												<Shuffle size={14} />
												Shuffle play
											</DropdownMenu.Item>
										</DropdownMenu.Group>

										<DropdownMenu.Separator />

										<DropdownMenu.Group>
											<DropdownMenu.Item
												onSelect={() => {
													selectedPlaylists = [playlist.id];
												}}
											>
												Select playlist
											</DropdownMenu.Item>

											{#if data.user?.quickPlaylist !== playlist.id}
												<DropdownMenu.Item
													onSelect={async () => {
														const res = await apiClient.setQuickPlaylist({
															playlistId: playlist.id,
														});
														if (!res.success) {
															handleApiError(res.error);
															return;
														}

														await invalidateAll();
													}}
												>
													<FileHeart />
													Set as Quick Playlist
												</DropdownMenu.Item>
											{:else}
												<DropdownMenu.Item
													onSelect={async () => {
														const res = await apiClient.setQuickPlaylist({
															playlistId: "",
														});
														if (!res.success) {
															handleApiError(res.error);
															return;
														}

														await invalidateAll();
													}}
												>
													<FileHeart />
													Remove Quick Playlist
												</DropdownMenu.Item>
											{/if}
										</DropdownMenu.Group>
									</DropdownMenu.Content>
								</DropdownMenu.Root>
							{/if}
						</div>

<p class="truncate text-xs text-muted-foreground">
						{playlist.trackCount} track{playlist.trackCount !== 1 ? "s" : ""}
						{#if playlist.playTime > 0}
							&middot; {formatPlayTime(playlist.playTime)}
						{/if}
					</p>
					</div>
				</div>
			{/each}
		</div>
	</InfiniteScroll>
</div>

<NewPlaylistModal bind:open={openNewPlaylistModal} />
