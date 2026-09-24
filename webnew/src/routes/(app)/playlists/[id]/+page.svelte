<script lang="ts">
	import { goto, invalidateAll } from "$app/navigation";
	import { page } from "$app/state";
	import { onMount } from "svelte";
	import { getApiClient, handleApiError } from "$lib";
	import { formatPlayTime } from "$lib/utils";
	import ConfirmModal from "$lib/components/new-modals/ConfirmModal.svelte";
	import HeroCard from "$lib/components/HeroCard.svelte";
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
	} from "$lib/components/ui";
	import {
		EllipsisVertical,
		ListPlus,
		Pencil,
		Play,
		Shuffle,
		Trash,
		Upload,
		WandSparkles,
	} from "@lucide/svelte";
	import EditPlaylistModal from "./EditPlaylistModal.svelte";
	import UploadPlaylistCoverModal from "./UploadPlaylistCoverModal.svelte";
	import { toast } from "svelte-sonner";
	import SectionImage from "$lib/components/SectionImage.svelte";
	import Spacer from "$lib/components/Spacer.svelte";
	import DebouncedSearchInput from "$lib/components/DebouncedSearchInput.svelte";
	import {
		SortableHeader,
		SortToggleDropdown,
		defaultSort,
		buildTrackQuery,
		trackColumns,
		trackSortTypes,
		type SortType,
	} from "$lib/components/sort";

	const { data } = $props();
	const musicManager = getMusicManager();
	const apiClient = getApiClient();

	let openConfirmDelete = $state(false);
	let openEditPlaylistModal = $state(false);
	let openUploadCoverModal = $state(false);

	let sort = $state(
		(page.url.searchParams.get("sort") as SortType) ?? defaultSort,
	);

	function updateSort(value: string) {
		sort = value as SortType;

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

	let selectedTracks = $state<string[]>([]);

	async function toggleSelectAll() {
		if (selectedTracks.length > 0) {
			selectedTracks = [];
			return;
		}

		const res = await apiClient.getPlaylistItemIds(data.playlist.id);
		if (!res.success) {
			handleApiError(res.error);
			return;
		}

		selectedTracks = res.data.ids;
	}

	const isOwner = $derived(data.user?.id === data.playlist.ownerId);

	let heroRef = $state<HTMLElement | null>(null);
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
			const query: Record<string, string> = {
				page: String(page),
				perPage: String(data.page.perPage),
			};

			buildTrackQuery(data.filter, "playlist", query);

			const res = await apiClient.getPlaylistItems(data.playlist.id, {
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

<HeroCard
	bind:ref={heroRef}
	class="section-playlists"
	innerClass="md:flex-row md:items-end md:gap-8"
>
	<SectionImage
		class="aspect-square w-40 min-w-40 self-center rounded-xl shadow-2xl md:w-52 md:min-w-52"
		src={data.playlist.coverArt.large}
		alt={data.playlist.name}
	/>

	<!-- <Image -->
	<!-- 	class="w-40 min-w-40 self-center rounded-xl shadow-2xl ring-1 ring-black/15 transition-transform duration-300 hover:scale-[1.02] md:w-52 md:min-w-52 dark:ring-white/10" -->
	<!-- 	src={data.playlist.coverArt.large} -->
	<!-- 	alt={data.playlist.name} -->
	<!-- /> -->

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
								<WandSparkles />
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
</HeroCard>

<Spacer />

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

		<p class="line-clamp-1 text-sm font-semibold text-ellipsis">
			{data.playlist.name}
		</p>

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

<Spacer />

<div class="flex items-center justify-between gap-2 sm:hidden">
	<DebouncedSearchInput
		class="flex-1"
		placeholder="Search tracks..."
		{value}
		setValue={(v) => (value = v)}
		{search}
	/>

	<div class="flex items-center gap-2 pr-2">
		<SortToggleDropdown
			types={trackSortTypes.playlist}
			{sort}
			{defaultSort}
			onSortChange={updateSort}
		/>

		<Checkbox
			title="Select all"
			aria-label="Select all"
			checked={selectedTracks.length > 0}
			onCheckedChange={() => toggleSelectAll()}
		></Checkbox>
	</div>
</div>

<Spacer size="sm" />

<SortableHeader
	{sort}
	onSortChange={updateSort}
	columns={trackColumns("playlist")}
>
	<div class="flex shrink-0 items-center gap-2 border-l border-border/40 pl-3">
		<DebouncedSearchInput
			inputClass="h-7 w-44 pr-6"
			iconSize={12}
			placeholder="Search tracks..."
			{value}
			setValue={(v) => (value = v)}
			{search}
		/>
		<Checkbox
			title="Select all"
			aria-label="Select all"
			checked={selectedTracks.length > 0}
			onCheckedChange={() => toggleSelectAll()}
		/>
	</div>
</SortableHeader>

<Spacer size="sm" />

<InfiniteScroll controller={scroll} errorMessage="Failed to load more tracks">
	<TrackList
		displayOrder
		{totalTracks}
		tracks={scroll.items}
		bind:selectedTracks
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
