<script lang="ts">
	import { EllipsisVertical, Play, Shuffle } from "@lucide/svelte";
	import { buttonVariants, DropdownMenu } from "$lib/components/ui";
	import { getApiClient, handleApiError } from "$lib";
	import type { Playlist } from "$lib/api/types";
	import { cn, formatPlayTime } from "$lib/utils";
	import { getMusicManager } from "$lib/music-manager.svelte";
	import InfiniteScroll from "$lib/components/InfiniteScroll.svelte";
	import { InfiniteScrollController } from "$lib/infinite-scroll.svelte";
	import Spacer from "$lib/components/Spacer.svelte";

	let { data } = $props();

	const musicManager = getMusicManager();
	const apiClient = getApiClient();

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
				filter: `ownerId = "${data.userData.id}"`,
				sort: "position",
			};

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

	function playPlaylist(playlistId: string, shuffle = false) {
		return musicManager.queueRequest(
			{ type: "addPlaylist", playlistId },
			{ shuffle },
		);
	}
</script>

<div class="flex flex-col gap-6">
	<div
		class="flex flex-col gap-6 rounded-lg border bg-linear-to-b from-[oklch(0.93_0.045_75)] to-background p-4 shadow-sm sm:p-6 md:flex-row md:items-end md:gap-8 dark:from-[oklch(0.24_0.03_80)] dark:to-background"
	>
		<div class="flex min-w-0 flex-col gap-2">
			<p
				class="text-xs font-semibold tracking-wider text-muted-foreground uppercase"
			>
				Playlists
			</p>

			<h1 class="line-clamp-2 text-2xl font-bold md:text-4xl">
				{data.userData.displayName}
			</h1>

			<p class="text-sm text-muted-foreground">
				{data.page.totalItems}
				{data.page.totalItems === 1 ? "playlist" : "playlists"}
			</p>
		</div>
	</div>
</div>

<Spacer size="lg" />

<InfiniteScroll controller={scroll}>
	<div
		class="grid grid-cols-2 gap-3 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-5 xl:grid-cols-6 2xl:grid-cols-7"
	>
		{#each scroll.items as playlist (playlist.id)}
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

					<button
						class="absolute right-2 bottom-2 hidden h-10 w-10 translate-y-2 items-center justify-center rounded-full bg-primary text-primary-foreground opacity-0 shadow-lg transition-all duration-300 group-hover:translate-y-0 group-hover:scale-105 group-hover:opacity-100 hover:scale-110 sm:flex"
						title="Play playlist"
						aria-label={`Play ${playlist.name}`}
						onclick={() => playPlaylist(playlist.id)}
					>
						<Play size={18} />
					</button>
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
							</DropdownMenu.Content>
						</DropdownMenu.Root>
					</div>

					<p class="truncate text-xs text-muted-foreground">
						{playlist.trackCount}
						{playlist.trackCount !== 1 ? "tracks" : "track"}
						{#if playlist.playTime > 0}
							&middot; {formatPlayTime(playlist.playTime)}
						{/if}
					</p>
				</div>
			</div>
		{/each}
	</div>

	{#if scroll.items.length === 0}
		<p class="px-2 text-sm text-muted-foreground">No playlists yet.</p>
	{/if}
</InfiniteScroll>
