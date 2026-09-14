<script lang="ts">
	import { goto, invalidateAll } from "$app/navigation";
	import { getApiClient, handleApiError } from "$lib";
	import ConfirmModal from "$lib/components/new-modals/ConfirmModal.svelte";
	import Image from "$lib/components/Image.svelte";
	import Spinner from "$lib/components/Spinner.svelte";
	import TrackList from "$lib/components/track-list/TrackList.svelte";
	import { getMusicManager } from "$lib/music-manager.svelte.js";
	import type { Track } from "$lib/api/types";
	import {
		Breadcrumb,
		Button,
		buttonVariants,
		Checkbox,
		DropdownMenu,
		Input,
		Select,
		Separator,
	} from "$lib/components/ui";
	import {
		CheckIcon,
		EllipsisVertical,
		ListPlus,
		ListSortAscendingIcon,
		Pencil,
		Play,
		Shuffle,
		Trash,
		Upload,
		Wand2,
	} from "@lucide/svelte";
	import EditPlaylistModal from "./EditPlaylistModal.svelte";
	import UploadPlaylistCoverModal from "./UploadPlaylistCoverModal.svelte";
	import { toast } from "svelte-sonner";

	const { data } = $props();
	const musicManager = getMusicManager();
	const apiClient = getApiClient();

	let openConfirmDelete = $state(false);
	let openEditPlaylistModal = $state(false);
	let openUploadCoverModal = $state(false);

	const sortOptions = [
		{ label: "Playlist Order", value: "playlist-order" },
		{ label: "Reverse Playlist Order", value: "playlist-order-reverse" },
		{ label: "Name (A-Z)", value: "name-a-z" },
		{ label: "Name (Z-A)", value: "name-z-a" },
		{ label: "Recently Added (New–Old)", value: "added-new" },
		{ label: "Recently Added (Old-New)", value: "added-old" },
	] as const;
	type SortOption = (typeof sortOptions)[number]["value"];

	let selectedSort = $state<SortOption>("playlist-order");

	const isOwner = $derived(data.user?.id === data.playlist.ownerId);

	let heroRef = $state<HTMLDivElement | null>(null);
	let showCompactHeader = $state(false);

	$effect(() => {
		const el = heroRef;
		if (!el) return;

		const observer = new IntersectionObserver(
			([entry]) => {
				showCompactHeader = !entry.isIntersecting;
			},
			{ rootMargin: "-56px 0px 0px 0px", threshold: 0 },
		);

		observer.observe(el);
		return () => observer.disconnect();
	});

	$effect(() => {
		if (showCompactHeader) {
			document.body.setAttribute("data-hero", "true");
		} else {
			document.body.setAttribute("data-hero", "false");
		}

		return () => {
			document.body.removeAttribute("data-hero");
		};
	});

	async function playPlaylist(opts: { shuffle?: boolean } = {}) {
		await musicManager.queueRequest(
			{ type: "addPlaylist", playlistId: data.playlist.id },
			{ shuffle: opts.shuffle },
		);
	}

	let tracks = $state<Track[]>([]);
	let totalItems = $state(0);
	let currentPage = $state(0);
	let trackIds = $state<Set<string>>(new Set());
	let loadingMore = $state(false);
	let loadMoreError = $state(false);
	let sentinel = $state<HTMLDivElement | null>(null);

	$effect(() => {
		tracks = [...data.items];
		totalItems = data.page.totalItems;
		currentPage = data.page.page;
		trackIds = new Set(data.items.map((t) => t.id));
	});

	const hasMore = $derived(tracks.length < totalItems);

	async function loadMore() {
		if (loadingMore || !hasMore) return;

		loadingMore = true;
		loadMoreError = false;

		const res = await apiClient.getPlaylistItems(data.playlist.id, {
			query: {
				page: String(currentPage + 1),
				perPage: String(data.page.perPage),
			},
		});

		if (res.success) {
			const newItems = res.data.items.filter((t) => !trackIds.has(t.id));
			for (const track of newItems) {
				trackIds.add(track.id);
			}
			tracks = [...tracks, ...newItems];
			currentPage = res.data.page.page;
			totalItems = res.data.page.totalItems;
		} else {
			handleApiError(res.error);
			loadMoreError = true;
		}

		loadingMore = false;
	}

	$effect(() => {
		const el = sentinel;
		if (!el) return;

		const observer = new IntersectionObserver(
			([entry]) => {
				if (entry.isIntersecting) {
					loadMore();
				}
			},
			{ rootMargin: "200px 0px 0px 0px" },
		);

		observer.observe(el);
		return () => observer.disconnect();
	});
</script>

<div class="py-2">
	<Breadcrumb.Root>
		<Breadcrumb.List>
			<Breadcrumb.Item>
				<Breadcrumb.Link href="/library/playlists">Playlists</Breadcrumb.Link>
			</Breadcrumb.Item>
			<Breadcrumb.Separator />
			<Breadcrumb.Item>
				<Breadcrumb.Page>{data.playlist.name}</Breadcrumb.Page>
			</Breadcrumb.Item>
		</Breadcrumb.List>
	</Breadcrumb.Root>
</div>

<div
	bind:this={heroRef}
	class="flex flex-col gap-6 rounded-lg border bg-linear-to-b from-accent to-background p-4 sm:p-6 md:flex-row md:items-end md:gap-8 dark:from-zinc-900 dark:to-background"
>
	<Image
		class="w-40 min-w-40 self-center shadow-lg md:w-52 md:min-w-52"
		src={data.playlist.coverArt.large}
		alt={data.playlist.name}
	/>

	<div class="flex min-w-0 flex-col gap-2">
		<p
			class="text-xs font-semibold tracking-wider text-muted-foreground uppercase"
		>
			Playlist
		</p>

		<h1 class="line-clamp-2 text-2xl font-bold md:text-4xl">
			{data.playlist.name}
		</h1>

		<div
			class="flex flex-wrap items-center gap-x-1 text-sm text-muted-foreground"
		>
			<span class="font-medium text-foreground"
				>{data.playlist.ownerDisplayName}</span
			>
			<span>&middot; {data.playlist.trackCount} tracks</span>
		</div>

		<div class="flex gap-2 pt-2">
			<Button onclick={() => playPlaylist()}>
				<Play size={14} />
				Play
			</Button>

			<Button variant="ghost" onclick={() => playPlaylist({ shuffle: true })}>
				<Shuffle size={14} />
			</Button>

			<DropdownMenu.Root>
				<DropdownMenu.Trigger
					class={buttonVariants({ variant: "ghost", size: "icon" })}
				>
					<EllipsisVertical />
				</DropdownMenu.Trigger>
				<DropdownMenu.Content align="start">
					<DropdownMenu.Group>
						<DropdownMenu.Item
							onSelect={async () => {
								await musicManager.queueRequest(
									{ type: "addPlaylist", playlistId: data.playlist.id },
									{ append: "back" },
								);
							}}
						>
							<ListPlus />
							Append to Queue
						</DropdownMenu.Item>

						{#if isOwner}
							<DropdownMenu.Separator />

							<DropdownMenu.Item
								onSelect={() => {
									openEditPlaylistModal = true;
								}}
							>
								<Pencil />
								Edit Playlist
							</DropdownMenu.Item>

							<DropdownMenu.Item
								onSelect={() => {
									openUploadCoverModal = true;
								}}
							>
								<Upload />
								Upload Cover
							</DropdownMenu.Item>

							<DropdownMenu.Item
								onSelect={async () => {
									const res = await apiClient.generatePlaylistImage(
										data.playlist.id,
									);
									if (!res.success) {
										handleApiError(res.error);
									} else {
										toast.success("Generating cover");
									}
								}}
							>
								<Wand2 />
								Generate Cover
							</DropdownMenu.Item>

							<DropdownMenu.Separator />

							<DropdownMenu.Item
								onSelect={() => {
									openConfirmDelete = true;
								}}
							>
								<Trash />
								Delete Playlist
							</DropdownMenu.Item>
						{/if}
					</DropdownMenu.Group>
				</DropdownMenu.Content>
			</DropdownMenu.Root>
		</div>
	</div>
</div>

<div class="h-4"></div>

<div
	role="button"
	tabindex={showCompactHeader ? 0 : -1}
	onclick={() => window.scrollTo({ top: 0, behavior: "smooth" })}
	onkeydown={(e) => {
		if (e.key === "Enter" || e.key === " ") {
			e.preventDefault();
			window.scrollTo({ top: 0, behavior: "smooth" });
		}
	}}
	class="fixed inset-x-0 top-14 z-40 cursor-pointer border-b bg-background/95 backdrop-blur transition-opacity duration-300 supports-backdrop-filter:bg-background/60 {showCompactHeader
		? 'opacity-100'
		: 'pointer-events-none opacity-0'}"
>
	<div class="old-container mx-auto flex h-14 items-center gap-3 px-6 sm:px-8">
		<Image
			class="h-10 w-10 shrink-0 shadow-sm"
			src={data.playlist.coverArt.small}
			alt={data.playlist.name}
		/>

		<p class="truncate text-sm font-semibold">{data.playlist.name}</p>

		<div class="grow"></div>

		<Button
			size="icon-sm"
			title="Play"
			aria-label="Play"
			onclick={(e) => {
				e.stopPropagation();
				playPlaylist();
			}}
		>
			<Play />
		</Button>

		<Button
			size="icon-sm"
			variant="ghost"
			title="Shuffle"
			aria-label="Shuffle"
			onclick={(e) => {
				e.stopPropagation();
				playPlaylist({ shuffle: true });
			}}
		>
			<Shuffle />
		</Button>
	</div>
</div>

<div class="h-4"></div>

<div class="flex items-center justify-between gap-2 px-2">
	<Input class="md:max-w-64" placeholder="Search tracks..." />

	<div class="flex items-center gap-2 pr-2">
		<!-- <Select.Root -->
		<!-- 	type="single" -->
		<!-- 	allowDeselect={false} -->
		<!-- 	value={selectedSort} -->
		<!-- 	onValueChange={(v) => (selectedSort = v as SortOption)} -->
		<!-- > -->
		<!-- 	<Select.Trigger hideIcon> -->
		<!-- 		<ListSortAscendingIcon /> -->
		<!-- 	</Select.Trigger> -->
		<!-- 	<Select.Trigger size="sm" class="w-44"> -->
		<!-- 	  {sortOptions.find((i) => i.value === selectedSort)?.label ?? "Sort"} -->
		<!-- 	</Select.Trigger> -->
		<!-- 	<Select.Content align="end"> -->
		<!-- 		{#each sortOptions as opt (opt.value)} -->
		<!-- 			<Select.Item value={opt.value} label={opt.label} /> -->
		<!-- 		{/each} -->
		<!-- 	</Select.Content> -->
		<!-- </Select.Root> -->

		<DropdownMenu.Root>
			<DropdownMenu.Trigger
				class={buttonVariants({ variant: "ghost", size: "icon" })}
			>
				<ListSortAscendingIcon />
			</DropdownMenu.Trigger>
			<DropdownMenu.Content align="end">
				<DropdownMenu.Group>
					{#each sortOptions as opt, i (opt.value)}
						{@const selected = i === 1}
						<DropdownMenu.Item
							onSelect={async () => {}}
							class={selected ? "bg-accent text-foreground" : ""}
						>
							{#if selected}
								<CheckIcon />
							{/if}
							{opt.label}
						</DropdownMenu.Item>
					{/each}
				</DropdownMenu.Group>
			</DropdownMenu.Content>
		</DropdownMenu.Root>

		<Checkbox></Checkbox>
	</div>
</div>

<div class="h-2"></div>
<Separator />

<TrackList
	displayOrder
	totalTracks={totalItems}
	{tracks}
	onPlay={async (trackId, shuffle) => {
		await musicManager.queueRequest(
			{ type: "addPlaylist", playlistId: data.playlist.id },
			{ queueIndexToTrackId: trackId, shuffle },
		);
	}}
	onReorder={async (items, anchor) => {
		const res = await apiClient.reorderPlaylistItems(data.playlist.id, {
			before: false,
			anchorTrackId: anchor ?? "",
			trackIds: items,
		});
		if (!res.success) {
			return handleApiError(res.error);
		}

		toast.success("Updated playlist");
		invalidateAll();
	}}
/>

<div class="h-4"></div>

<div bind:this={sentinel}></div>

{#if loadingMore}
	<div class="flex justify-center py-6">
		<Spinner />
	</div>
{/if}

{#if loadMoreError}
	<div class="flex flex-col items-center gap-2 py-6">
		<p class="text-sm text-muted-foreground">Failed to load more tracks</p>
		<Button size="sm" variant="outline" onclick={loadMore}>Retry</Button>
	</div>
{/if}

<ConfirmModal
	bind:open={openConfirmDelete}
	removeTrigger
	confirmDelete
	onResult={async () => {
		const res = await apiClient.deletePlaylist(data.playlist.id);
		if (!res.success) {
			handleApiError(res.error);
			invalidateAll();
			return;
		}

		toast.success("Deleted playlist");
		goto("/library/playlists", { invalidateAll: true });
	}}
/>

<EditPlaylistModal
	bind:open={openEditPlaylistModal}
	playlist={data.playlist}
/>

<UploadPlaylistCoverModal
	bind:open={openUploadCoverModal}
	playlistId={data.playlist.id}
/>
