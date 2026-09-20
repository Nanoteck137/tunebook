<script lang="ts">
	import Image from "$lib/components/Image.svelte";
	import TrackListItem from "$lib/components/track-list/TrackListItem.svelte";
	import { getMusicManager } from "$lib/music-manager.svelte";
	import {
		Breadcrumb,
		Button,
		buttonVariants,
		DropdownMenu,
		Separator,
	} from "$lib/components/ui";
	import {
		Disc,
		EllipsisVertical,
		ListPlus,
		Music,
		Play,
		Shuffle,
		TrendingUp,
		Users,
	} from "@lucide/svelte";
	import SectionHeader from "$lib/components/SectionHeader.svelte";
	import TopTrackItem from "$lib/components/track-list/TopTrackItem.svelte";

	const { data } = $props();
	const musicManager = getMusicManager();

	let madeStats = $derived(
		`${data.albumPage.totalItems} ${data.albumPage.totalItems === 1 ? "album" : "albums"}` +
			` · ${data.trackPage.totalItems} ${data.trackPage.totalItems === 1 ? "track" : "tracks"}`,
	);

	let featuredCount = $derived(
		data.featuredAlbumPage.totalItems + data.featuredTrackPage.totalItems,
	);

	function otherArtists(artists: { id: string; name: string }[]) {
		return artists.filter((a) => a.id !== data.artist.id).map((a) => a.name);
	}
</script>

<div class="py-2">
	<Breadcrumb.Root>
		<Breadcrumb.List>
			<Breadcrumb.Item>
				<Breadcrumb.Link href="/artists">Artists</Breadcrumb.Link>
			</Breadcrumb.Item>
			<Breadcrumb.Separator />
			<Breadcrumb.Item>
				<Breadcrumb.Page>{data.artist.name}</Breadcrumb.Page>
			</Breadcrumb.Item>
		</Breadcrumb.List>
	</Breadcrumb.Root>
</div>

<div
	class="flex flex-col gap-6 rounded-lg border bg-linear-to-b from-[oklch(0.93_0.045_75)] to-background p-4 shadow-sm sm:p-6 md:flex-row md:items-end md:gap-8 dark:from-[oklch(0.24_0.03_80)] dark:to-background"
>
	<Image
		class="w-40 min-w-40 self-center rounded-xl shadow-2xl ring-1 ring-black/15 transition-transform duration-300 hover:scale-[1.02] md:w-52 md:min-w-52 dark:ring-white/10"
		src={data.artist.coverArt.large}
		alt={data.artist.name}
	/>

	<div class="flex min-w-0 flex-col gap-2">
		<p
			class="text-xs font-semibold tracking-wider text-muted-foreground uppercase"
		>
			Artist
		</p>

		<h1 class="line-clamp-2 text-2xl font-bold md:text-4xl">
			{data.artist.name}
		</h1>

		<p class="text-sm text-muted-foreground">
			{madeStats}
			{#if featuredCount > 0}
				&middot; Appears on {featuredCount}
			{/if}
		</p>

		{#if data.artist.tags.length > 0}
			<div class="flex flex-wrap gap-1">
				{#each data.artist.tags as tag (tag)}
					<span
						class="rounded-full bg-secondary px-2.5 py-0.5 text-xs text-secondary-foreground"
					>
						{tag}
					</span>
				{/each}
			</div>
		{/if}

		<div class="flex gap-2 pt-2">
			<Button
				onclick={async () => {
					await musicManager.queueRequest(
						{ type: "addArtist", artistId: data.artist.id },
						{},
					);
				}}
			>
				<Play />
				Play
			</Button>

			<Button
				variant="ghost"
				size="icon"
				onclick={async () => {
					await musicManager.queueRequest(
						{ type: "addArtist", artistId: data.artist.id },
						{ shuffle: true },
					);
				}}
				title="Shuffle"
				aria-label="Shuffle"
			>
				<Shuffle />
			</Button>

			<DropdownMenu.Root>
				<DropdownMenu.Trigger
					class={buttonVariants({ variant: "ghost", size: "icon" })}
					title="More options"
					aria-label="More options"
				>
					<EllipsisVertical />
				</DropdownMenu.Trigger>
				<DropdownMenu.Content align="start">
					<DropdownMenu.Group>
						<DropdownMenu.Item
							onSelect={async () => {
								await musicManager.queueRequest(
									{ type: "addArtist", artistId: data.artist.id },
									{ append: "back" },
								);
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

<div class="flex flex-col gap-10">
	{#if data.userTopTracks.length > 0}
		<section>
			<SectionHeader
				count={data.userTopTrackPage?.totalItems ?? data.userTopTracks.length}
				viewAllHref={data.user
					? `/users/${data.user.id}/review/total/top-tracks`
					: undefined}
			>
				<TrendingUp />
				Your Top Tracks
			</SectionHeader>

			<div class="flex flex-col">
				{#each data.userTopTracks as track (track.id)}
					<TopTrackItem
						rank={track.rank}
						{track}
						playCount={track.playCount}
					/>
				{/each}
			</div>
		</section>
	{/if}

	{#if data.tracks.length > 0}
		<section>
			<SectionHeader
				count={data.trackPage.totalItems}
				viewAllHref={data.trackPage.totalItems > data.tracks.length
					? `/artists/${data.artist.id}/tracks`
					: undefined}
			>
				<Music />
				Songs
			</SectionHeader>

			<div class="flex flex-col">
				{#each data.tracks as track, i (track.id)}
					<TrackListItem {track} />
					{#if i < data.tracks.length - 1}
						<Separator />
					{/if}
				{/each}
			</div>
		</section>
	{/if}

	{#if data.albums.length > 0}
		<section>
			<SectionHeader
				count={data.albumPage.totalItems}
				viewAllHref={data.albumPage.totalItems > data.albums.length
					? `/artists/${data.artist.id}/albums`
					: undefined}
			>
				<Disc />
				Albums
			</SectionHeader>

			<div
				class="grid grid-cols-2 gap-3 px-2 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-5"
			>
				{#each data.albums as album (album.id)}
					<a
						href="/albums/{album.id}"
						class="group flex flex-col"
						title={album.name}
					>
						<div class="relative overflow-hidden rounded-lg">
							<Image
								class="aspect-square w-full rounded-none transition-transform duration-300 group-hover:scale-105"
								src={album.coverArt.medium}
								alt={album.name}
							/>
						</div>
						<div class="flex flex-col gap-0.5 pt-2">
							<p
								class="truncate text-sm font-medium group-hover:underline"
								title={album.name}
							>
								{album.name}
							</p>
							<p
								class="truncate text-xs text-muted-foreground"
								title={album.artists.map((a) => a.name).join(", ")}
							>
								{#if album.year}
									{album.year}
								{:else}
									{album.albumType}
								{/if}
							</p>
						</div>
					</a>
				{/each}
			</div>
		</section>
	{/if}

	{#if featuredCount > 0}
		<Separator />

		<section>
			<SectionHeader count={featuredCount}>
				<Users />
				Appears on
			</SectionHeader>

			{#if data.featuredAlbums.length > 0}
				<div
					class="grid grid-cols-2 gap-3 px-2 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-5"
				>
					{#each data.featuredAlbums as album (album.id)}
						<a
							href="/albums/{album.id}"
							class="group flex flex-col"
							title={album.name}
						>
							<div class="relative overflow-hidden rounded-lg">
								<Image
									class="aspect-square w-full rounded-none transition-transform duration-300 group-hover:scale-105"
									src={album.coverArt.medium}
									alt={album.name}
								/>
							</div>
							<div class="flex flex-col gap-0.5 pt-2">
								<p
									class="truncate text-sm font-medium group-hover:underline"
									title={album.name}
								>
									{album.name}
								</p>
								<p
									class="truncate text-xs text-muted-foreground"
									title={otherArtists(album.artists).join(", ")}
								>
									{otherArtists(album.artists).join(", ")}
								</p>
							</div>
						</a>
					{/each}
				</div>
			{/if}

			{#if data.featuredTracks.length > 0}
				<div class="mt-4 flex flex-col">
					{#each data.featuredTracks as track, i (track.id)}
						<TrackListItem {track} />
						{#if i < data.featuredTracks.length - 1}
							<Separator />
						{/if}
					{/each}
				</div>
			{/if}
		</section>
	{/if}
</div>
