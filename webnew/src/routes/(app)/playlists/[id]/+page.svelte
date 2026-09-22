<script lang="ts">
	import { goto, invalidateAll } from "$app/navigation";
	import { getApiClient, handleApiError } from "$lib";
	import { formatPlayTime } from "$lib/utils";
	import ConfirmModal from "$lib/components/new-modals/ConfirmModal.svelte";
	import Image from "$lib/components/Image.svelte";
	import TrackList from "$lib/components/track-list/TrackList.svelte";
	import { getMusicManager } from "$lib/music-manager.svelte";
	import { InfiniteScrollController } from "$lib/infinite-scroll.svelte";
	import InfiniteScroll from "$lib/components/InfiniteScroll.svelte";
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

	const totalTracks = $derived(data.page.totalItems);

	const scroll = new InfiniteScrollController<Track>({
		initialLoad: () => ({
			items: data.items,
			hasMore: data.page.page + 1 < data.page.totalPages,
			page: data.page.page,
		}),
		load: async (page) => {
			const res = await apiClient.getPlaylistItems(data.playlist.id, {
				query: { page: String(page), perPage: String(data.page.perPage) },
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
	class="section-playlists flex flex-col gap-6 rounded-lg border bg-linear-to-b from-section-hero-from to-section-hero-to p-4 shadow-sm sm:p-6 md:flex-row md:items-end md:gap-8"
>
	<Image
		class="w-40 min-w-40 self-center rounded-xl shadow-2xl ring-1 ring-black/15 transition-transform duration-300 hover:scale-[1.02] md:w-52 md:min-w-52 dark:ring-white/10"
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
			<a
				href="/users/{data.playlist.ownerId}"
				class="font-medium text-foreground hover:underline"
				title={data.playlist.ownerDisplayName}
			>
				{data.playlist.ownerDisplayName}
			</a>
			<span>&middot; {data.playlist.trackCount} tracks</span>
			{#if data.playlist.playTime > 0}
				<span>&middot; {formatPlayTime(data.playlist.playTime)}</span>
			{/if}
		</div>

		<div class="flex gap-2 pt-2">
			<Button onclick={() => playPlaylist()}>
				<Play />
				Play
			</Button>

			<Button
				variant="ghost"
				size="icon"
				onclick={() => playPlaylist({ shuffle: true })}
				title="Shuffle"
				aria-label="Shuffle"
			>
				<Shuffle />
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

<InfiniteScroll controller={scroll} errorMessage="Failed to load more tracks">
	<TrackList
		displayOrder
		{totalTracks}
		tracks={scroll.items}
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
</InfiniteScroll>

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
