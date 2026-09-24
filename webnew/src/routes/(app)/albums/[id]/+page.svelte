<script lang="ts">
	import { page } from "$app/state";
	import { afterNavigate, goto } from "$app/navigation";
	import {
		Breadcrumb,
		Button,
		buttonVariants,
		Checkbox,
		DropdownMenu,
		Input,
	} from "$lib/components/ui";
	import {
		EllipsisVertical,
		ListPlus,
		Play,
		Shuffle,
		X,
	} from "@lucide/svelte";
	import {
		SortDropdown,
		SortableHeader,
		defaultSort,
		trackColumns,
		trackSortTypes,
		type SortType,
	} from "$lib/components/sort";
	import TrackList from "$lib/components/track-list/TrackList.svelte";
	import { getMusicManager } from "$lib/music-manager.svelte.js";
	import Image from "$lib/components/Image.svelte";
	import ArtistList from "$lib/components/ArtistList.svelte";
	import { formatPlayTime } from "$lib/utils";
	import SectionImage from "$lib/components/SectionImage.svelte";
	import Spacer from "$lib/components/Spacer.svelte";

	let { data } = $props();
	const musicManager = getMusicManager();

	let highlightTrackId = $derived(page.url.searchParams.get("track"));

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

	let searchQuery = $state(page.url.searchParams.get("query") ?? "");

	let searchTimer: ReturnType<typeof setTimeout>;

	function onSearchInput(e: Event) {
		const target = e.target as HTMLInputElement;
		const current = target.value;
		searchQuery = current;

		clearTimeout(searchTimer);
		searchTimer = setTimeout(() => {
			updateSearch();
		}, 500);
	}

	function updateSearch() {
		clearTimeout(searchTimer);

		const query = page.url.searchParams;
		query.delete("query");

		if (searchQuery) {
			query.set("query", searchQuery);
		}

		goto("?" + query.toString(), {
			invalidateAll: true,
			keepFocus: true,
			replaceState: true,
		});
	}

	function clearSearch() {
		searchQuery = "";
		updateSearch();
	}

	let selectedTracks = $state<string[]>([]);

	function toggleSelectAll() {
		if (selectedTracks.length > 0) {
			selectedTracks = [];
			return;
		}

		selectedTracks = data.tracks.map((track) => track.id);
	}

	let didHighlight = false;

	afterNavigate(() => {
		if (didHighlight || !highlightTrackId) return;
		didHighlight = true;

		requestAnimationFrame(() => {
			requestAnimationFrame(() => {
				document
					.getElementById(`track-${highlightTrackId}`)
					?.scrollIntoView({ behavior: "smooth", block: "center" });
			});
		});
	});

	let playTimeLabel = $derived.by(() => {
		const s = data.album.playTime;
		return s > 0 ? formatPlayTime(s) : "";
	});
</script>

<div class="py-2">
	<Breadcrumb.Root>
		<Breadcrumb.List>
			<Breadcrumb.Item>
				<Breadcrumb.Link href="/albums">Albums</Breadcrumb.Link>
			</Breadcrumb.Item>
			<Breadcrumb.Separator />
			<Breadcrumb.Item>
				<Breadcrumb.Page>{data.album.name}</Breadcrumb.Page>
			</Breadcrumb.Item>
		</Breadcrumb.List>
	</Breadcrumb.Root>
</div>

<div
	class="section-albums flex flex-col gap-6 rounded-lg border bg-linear-to-b from-section-hero-from to-section-hero-to p-4 shadow-sm sm:p-6 md:flex-row md:items-end md:gap-8"
>
	<SectionImage
		class="w-40 min-w-40 self-center rounded-xl shadow-2xl md:w-52 md:min-w-52"
		src={data.album.coverArt.large}
		alt={data.album.name}
	/>

	<div class="flex min-w-0 flex-col gap-2">
		<p
			class="text-xs font-semibold tracking-wider text-muted-foreground uppercase"
		>
			{data.album.albumType}
			{#if data.album.year}
				&middot; {data.album.year}
			{/if}
		</p>

		<h1 class="line-clamp-2 text-2xl font-bold md:text-4xl">
			{data.album.name}
		</h1>

		<div
			class="flex flex-wrap items-center gap-x-1 text-sm text-muted-foreground"
		>
			<ArtistList artists={data.album.artists} class="text-sm" />
		</div>

		<p class="text-sm text-muted-foreground">
			{data.tracks.length}
			{data.tracks.length === 1 ? "song" : "songs"}
			{#if playTimeLabel}&middot; {playTimeLabel}{/if}
		</p>

		{#if data.album.tags.length > 0}
			<div class="flex flex-wrap gap-1">
				{#each data.album.tags as tag (tag)}
					<span
						class="rounded-full bg-secondary px-2.5 py-0.5 text-xs text-secondary-foreground"
						>{tag}</span
					>
				{/each}
			</div>
		{/if}

		<div class="flex gap-2 pt-2">
			<Button
				onclick={async () => {
					await musicManager.addAlbumTracks({
						albumId: data.album.id,
						clear: true,
					});
				}}
			>
				<Play />
				Play
			</Button>

			<Button
				variant="ghost"
				size="icon"
				onclick={async () => {
					await musicManager.addAlbumTracks({
						albumId: data.album.id,
						clear: true,
						shuffle: true,
					});
				}}
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
								// await musicManager.queueRequest(
								//   { type: "addAlbum", albumId: data.album.id },
								//   { append: "back" },
								// );
							}}
						>
							<ListPlus />
							Append to Queue
						</DropdownMenu.Item>
					</DropdownMenu.Group>
				</DropdownMenu.Content>
			</DropdownMenu.Root>
		</div>
	</div>
</div>

<Spacer />

<div class="flex items-center justify-between gap-2 sm:hidden">
	<div class="relative flex-1">
		<Input
			class="pr-8"
			placeholder="Search tracks..."
			value={searchQuery}
			oninput={onSearchInput}
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

	<div class="flex items-center gap-2 pr-2">
		<SortDropdown
			types={trackSortTypes.album}
			{sort}
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

<SortableHeader
	{sort}
	onSortChange={updateSort}
	columns={trackColumns("album")}
>
	<div class="flex shrink-0 items-center gap-2 border-l border-border/40 pl-3">
		<div class="relative">
			<Input
				class="h-7 w-44 pr-6"
				placeholder="Search tracks..."
				value={searchQuery}
				oninput={onSearchInput}
				onkeydown={(e) => {
					if (e.key === "Enter") {
						updateSearch();
					}
				}}
			/>
			{#if searchQuery}
				<button
					class="absolute top-1/2 right-1 -translate-y-1/2 rounded-full p-0.5 text-muted-foreground hover:text-foreground"
					onclick={clearSearch}
					aria-label="Clear search"
				>
					<X size={12} />
				</button>
			{/if}
		</div>
		<Checkbox
			title="Select all"
			aria-label="Select all"
			checked={selectedTracks.length > 0}
			onCheckedChange={() => toggleSelectAll()}
		/>
	</div>
</SortableHeader>

<Spacer />

<TrackList
	isAlbumShowcase={true}
	totalTracks={data.tracks.length}
	tracks={data.tracks}
	highlightId={highlightTrackId}
	bind:selectedTracks
	onPlay={async (trackId) => {
		await musicManager.addAlbumTracks({
			albumId: data.album.id,
			trackId,
		});
	}}
/>
