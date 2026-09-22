<script lang="ts">
	import { page } from "$app/state";
	import { afterNavigate } from "$app/navigation";
	import {
		Breadcrumb,
		Button,
		buttonVariants,
		DropdownMenu,
		Input,
	} from "$lib/components/ui";
	import {
		EllipsisVertical,
		ListPlus,
		ListSortAscendingIcon,
		CheckIcon,
		Play,
		Shuffle,
	} from "@lucide/svelte";
	import TrackList from "$lib/components/track-list/TrackList.svelte";
	import { getMusicManager } from "$lib/music-manager.svelte.js";
	import Image from "$lib/components/Image.svelte";
	import ArtistList from "$lib/components/ArtistList.svelte";

	let { data } = $props();
	const musicManager = getMusicManager();

	let highlightTrackId = $derived(page.url.searchParams.get("track"));

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

	let durationLabel = $derived.by(() => {
		const totalMs = data.tracks.reduce((sum, t) => sum + t.duration, 0);
		const s = Math.floor(totalMs / 1000);
		const h = Math.floor(s / 3600);
		const m = Math.floor((s % 3600) / 60);
		return h > 0 ? `${h}h ${m}m` : `${m}m`;
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
	<Image
		class="w-40 min-w-40 self-center rounded-xl shadow-2xl ring-1 ring-black/15 transition-transform duration-300 hover:scale-[1.02] md:w-52 md:min-w-52 dark:ring-white/10"
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
			&middot; {durationLabel}
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

<div class="h-4"></div>

<div class="flex flex-wrap items-center justify-between gap-2">
	<div class="relative w-full md:max-w-56">
		<Input class="pr-8" placeholder="Search tracks..." />
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
					<DropdownMenu.Item>
						<CheckIcon />
						Name (A-Z)
					</DropdownMenu.Item>
					<DropdownMenu.Item>Name (Z-A)</DropdownMenu.Item>
					<DropdownMenu.Item>Disc / Track number</DropdownMenu.Item>
				</DropdownMenu.Group>
			</DropdownMenu.Content>
		</DropdownMenu.Root>
	</div>
</div>

<div class="h-4"></div>

<TrackList
	isAlbumShowcase={true}
	totalTracks={data.tracks.length}
	tracks={data.tracks}
	highlightId={highlightTrackId}
	onPlay={async (trackId) => {
		await musicManager.addAlbumTracks({
			albumId: data.album.id,
			trackId,
		});
	}}
/>
