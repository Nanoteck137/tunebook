<script lang="ts">
	import { goto } from "$app/navigation";
	import {
		Activity,
		BarChart3,
		CalendarDays,
		CalendarRange,
		ChartColumn,
		DiscAlbum,
		Heart,
		ListPlus,
		Music,
		Play,
		Shuffle,
		Sparkles,
		Tags,
		Users,
	} from "@lucide/svelte";
	import AlbumTile from "$lib/components/tiles/AlbumTile.svelte";
	import ArtistTile from "$lib/components/tiles/ArtistTile.svelte";
	import SectionHeader from "$lib/components/SectionHeader.svelte";
	import Spacer from "$lib/components/Spacer.svelte";
	import TopTrackItem from "$lib/components/track-list/TopTrackItem.svelte";
	import { Breadcrumb, DropdownMenu } from "$lib/components/ui";
	import { getFavorites } from "$lib/favorites.svelte";
	import { formatPlayTime } from "$lib/utils";
	import { toast } from "svelte-sonner";

	let { data } = $props();
	const favoritesManager = getFavorites();

	let currentYear = $derived(new Date().getFullYear());
	let review = $derived(data.review);
	let topTracks = $derived(data.topTracks);

	let topArtist = $derived(data.topArtists[0] ?? null);
	let topAlbum = $derived(data.topAlbums[0] ?? null);
	let topTrack = $derived(topTracks[0] ?? null);

	const MONTH_LABELS = [
		"Jan",
		"Feb",
		"Mar",
		"Apr",
		"May",
		"Jun",
		"Jul",
		"Aug",
		"Sep",
		"Oct",
		"Nov",
		"Dec",
	];

	let months = $derived(
		MONTH_LABELS.map((label, i) => {
			const row = data.months.find((m) => m.month === i + 1);

			return {
				label,
				month: i + 1,
				hours: (row?.playTime ?? 0) / 3600,
				playCount: row?.playCount ?? 0,
			};
		}),
	);

	let maxMonthlyHours = $derived(
		Math.max(...months.map((m) => m.hours), 1 / 60),
	);

	function formatHours(hours: number): string {
		if (hours >= 1) {
			return `${Math.round(hours)}h`;
		}

		return `${Math.round(hours * 60)}m`;
	}

	function formatPercent(value: number): string {
		return `${value.toLocaleString("en-US", { maximumFractionDigits: 1 })}%`;
	}

	let repeatRate = $derived(
		review.uniqueTracks > 0 ? review.trackCount / review.uniqueTracks : 0,
	);
	let skipRate = $derived(
		review.trackCount > 0 ? (review.skipCount / review.trackCount) * 100 : 0,
	);
	let favoritesOverlap = $derived(
		review.trackCount > 0
			? (review.favoritePlays / review.trackCount) * 100
			: 0,
	);

	function formatTagSlug(slug: string): string {
		const words = slug.split("-");
		const sentence = words.join(" ");
		return sentence.charAt(0).toUpperCase() + sentence.slice(1);
	}
</script>

{#snippet highlightCard(
	label: string,
	href: string | null,
	name: string,
	art: string | undefined,
)}
	<div class="flex flex-col gap-1.5 rounded-lg border bg-card p-4">
		<span class="text-xs text-muted-foreground">{label}</span>
		{#if name && href}
			<a {href} class="flex items-center gap-2">
				{#if art}
					<img
						src={art}
						alt=""
						class="size-10 shrink-0 rounded object-cover"
					/>
				{/if}
				<span class="truncate text-sm font-medium">{name}</span>
			</a>
		{:else}
			<span class="truncate text-sm font-medium text-muted-foreground">
				—
			</span>
		{/if}
	</div>
{/snippet}

{#snippet statCard(label: string, value: string)}
	<div class="flex flex-col gap-1.5 rounded-lg border bg-card p-4">
		<span class="text-xs text-muted-foreground">{label}</span>
		<span class="text-2xl font-bold">{value}</span>
	</div>
{/snippet}

<div class="flex flex-col gap-6">
	<Breadcrumb.Root>
		<Breadcrumb.List>
			<Breadcrumb.Item>
				<Breadcrumb.Link href="/users/{data.userData.id}/review">
					Year in Review
				</Breadcrumb.Link>
			</Breadcrumb.Item>
			<Breadcrumb.Separator />
			<Breadcrumb.Item>
				<Breadcrumb.Page>{data.year}</Breadcrumb.Page>
			</Breadcrumb.Item>
		</Breadcrumb.List>
	</Breadcrumb.Root>

	<SectionHeader count={review.trackCount}>
		<Sparkles />
		{data.year}

		{#if data.year === currentYear}
			<span
				class="rounded-full bg-primary/10 px-2.5 py-0.5 text-xs font-medium text-primary ring-1 ring-primary/25"
			>
				In progress
			</span>
		{/if}
	</SectionHeader>

	{#if review.trackCount === 0}
		<div
			class="flex flex-col items-center gap-2 rounded-lg border py-16 text-center"
		>
			<BarChart3 size={32} class="text-muted-foreground/40" />
			<p class="text-sm font-medium">No review for {data.year}</p>
			<p class="max-w-sm text-sm text-muted-foreground">
				A year in review is generated once there is listening data for the
				year.
				{#if data.year === currentYear}
					Listening so far this year will show up here as it is played.
				{/if}
			</p>
		</div>
	{:else}
		<p class="text-sm text-muted-foreground">
			{formatPlayTime(review.listeningTime)} listened across
			{review.uniqueTracks.toLocaleString()} unique tracks
		</p>

		<section>
			<SectionHeader>
				<Sparkles />
				Highlights
			</SectionHeader>
			<Spacer />

			<div class="grid gap-3 sm:grid-cols-2 lg:grid-cols-3">
				{@render highlightCard(
					"Top Artist",
					topArtist ? `/artists/${topArtist.id}` : null,
					topArtist?.name ?? "",
					topArtist?.coverArt.small,
				)}
				{@render highlightCard(
					"Top Album",
					topAlbum ? `/albums/${topAlbum.id}` : null,
					topAlbum?.name ?? "",
					topAlbum?.coverArt.small,
				)}
				{@render highlightCard(
					"Top Song",
					topTrack ? `/albums/${topTrack.albumId}` : null,
					topTrack?.name ?? "",
					topTrack?.coverArt.small,
				)}
			</div>
		</section>

		<section>
			<SectionHeader>
				<Activity />
				Listening Stats
			</SectionHeader>
			<Spacer />

			<div class="grid grid-cols-2 gap-3 md:grid-cols-4">
				{@render statCard("Plays", review.trackCount.toLocaleString())}
				{@render statCard(
					"Listening Time",
					formatPlayTime(review.listeningTime),
				)}
				{@render statCard(
					"Unique Tracks",
					review.uniqueTracks.toLocaleString(),
				)}
				{@render statCard("Skipped", review.skipCount.toLocaleString())}
				{@render statCard(
					"Avg Completion",
					formatPercent(review.avgCompletion),
				)}
				{@render statCard("Skip Rate", formatPercent(skipRate))}
				{@render statCard("Repeat Rate", `${repeatRate.toFixed(1)}x`)}
				<div class="flex flex-col gap-1.5 rounded-lg border bg-card p-4">
					<span class="text-xs text-muted-foreground">Favorites Overlap</span>
					<div class="flex items-baseline gap-1.5">
						<span class="text-2xl font-bold">
							{formatPercent(favoritesOverlap)}
						</span>
						<Heart size={16} class="shrink-0 fill-primary stroke-primary" />
					</div>
				</div>
			</div>
		</section>

		<section>
			<SectionHeader
				count={data.topTrackCount}
				viewAllHref="/users/{data.userData.id}/review/{data.year}/top-tracks"
			>
				<Music />
				Top Tracks
			</SectionHeader>
			<Spacer />

			<div class="flex flex-col">
				{#each topTracks as item (item.id)}
					<TopTrackItem
						rank={item.rank}
						track={item}
						playCount={item.playCount}
					>
						{#snippet menuItems()}
							<DropdownMenu.Group>
								<DropdownMenu.Item>
									<Play />
									Play
								</DropdownMenu.Item>
								<DropdownMenu.Item>
									<Shuffle />
									Shuffle play
								</DropdownMenu.Item>
							</DropdownMenu.Group>

							<DropdownMenu.Separator />

							<DropdownMenu.Item
								onSelect={() => {
									goto(`/albums/${item.albumId}`);
								}}
							>
								<DiscAlbum />
								Go to Album
							</DropdownMenu.Item>
							<DropdownMenu.Sub>
								<DropdownMenu.SubTrigger>
									<Users />
									Go to artist
								</DropdownMenu.SubTrigger>
								<DropdownMenu.SubContent>
									{#each item.artists as artist (artist.id)}
										<a
											href="/artists/{artist.id}"
											class="flex items-center gap-2 rounded-sm px-3 py-1.5 text-sm text-popover-foreground hover:bg-accent hover:text-accent-foreground"
										>
											{artist.name}
										</a>
									{/each}
								</DropdownMenu.SubContent>
							</DropdownMenu.Sub>
							<DropdownMenu.Separator />
							<DropdownMenu.Item>
								<ListPlus />
								Add to Playlist
							</DropdownMenu.Item>
							<DropdownMenu.Item
								onSelect={async () => {
									const wasFav = favoritesManager.hasTrack(item.id);
									await favoritesManager.toggleTrack(item.id);
									toast.success(
										wasFav ? "Removed from favorites" : "Added to favorites",
									);
								}}
							>
								{#if favoritesManager.hasTrack(item.id)}
									<Heart class="fill-primary stroke-primary" />
									Unfavorite
								{:else}
									<Heart />
									Favorite
								{/if}
							</DropdownMenu.Item>
						{/snippet}
					</TopTrackItem>
				{/each}
			</div>
		</section>

		<section>
			<SectionHeader
				count={data.topAlbumCount}
				viewAllHref="/users/{data.userData.id}/review/{data.year}/top-albums"
			>
				<DiscAlbum />
				Top Albums
			</SectionHeader>
			<Spacer />

			<div class="flex items-start gap-4 overflow-x-auto pb-2">
				{#each data.topAlbums as album (album.id)}
					<div class="flex w-40 shrink-0 flex-col">
						<AlbumTile size="sm" {album} />
					</div>
				{/each}
			</div>
		</section>

		<section>
			<SectionHeader
				count={data.topArtistCount}
				viewAllHref="/users/{data.userData.id}/review/{data.year}/top-artists"
			>
				<Users />
				Top Artists
			</SectionHeader>
			<Spacer />

			<div class="flex items-start gap-4 overflow-x-auto pb-2">
				{#each data.topArtists as artist (artist.id)}
					<div class="flex w-40 shrink-0 flex-col">
						<ArtistTile {artist} />
					</div>
				{/each}
			</div>
		</section>

		<section>
			<SectionHeader>
				<ChartColumn />
				Monthly Listening
			</SectionHeader>
			<Spacer />

			<div class="flex h-56 items-end gap-2 sm:h-48">
				{#each months as m (m.month)}
					<div class="group flex h-full flex-1 flex-col items-center gap-1.5">
						<span class="text-[10px] text-muted-foreground">
							{formatHours(m.hours)}
						</span>
						<div class="flex w-full flex-1 items-end">
							<div
								class="relative w-full"
								style="height: {m.hours > 0
									? Math.max((m.hours / maxMonthlyHours) * 100, 4)
									: 0}%"
							>
								{#if m.playCount > 0}
									<div
										class="pointer-events-none absolute bottom-full left-1/2 z-10 mb-1 -translate-x-1/2 rounded-md bg-foreground px-2 py-1 text-xs font-medium whitespace-nowrap text-background opacity-0 shadow-md transition-opacity group-hover:opacity-100"
									>
										{m.label}: {formatHours(m.hours)}
										&middot;
										{m.playCount.toLocaleString()} plays
									</div>
								{/if}
								<div
									class="h-full w-full rounded-t bg-linear-to-t from-logo-3 to-logo-1 transition-all group-hover:opacity-80"
								></div>
							</div>
						</div>
						<span class="text-[10px] text-muted-foreground">
							{m.label}
						</span>
					</div>
				{/each}
			</div>
		</section>

		<section>
			<SectionHeader>
				<CalendarDays />
				Monthly Breakdown
			</SectionHeader>
			<Spacer />

			<div class="grid grid-cols-4 gap-2 sm:flex sm:flex-wrap">
				{#each months as m (m.month)}
					<a
						href="/users/{data.userData
							.id}/review/{data.year}/months/{m.month}"
						class="flex min-w-0 flex-col items-center gap-0.5 rounded-lg border bg-card px-2 py-2 transition-colors hover:bg-accent sm:min-w-[88px] sm:flex-none"
						class:pointer-events-none={m.playCount === 0}
						class:opacity-40={m.playCount === 0}
						aria-disabled={m.playCount === 0}
					>
						<span class="text-sm font-medium">{m.label}</span>
						<span class="text-xs text-muted-foreground">
							{m.playCount.toLocaleString()}
						</span>
					</a>
				{/each}
			</div>
		</section>

		<section>
			<SectionHeader
				count={data.topTagCount}
				viewAllHref="/users/{data.userData.id}/review/{data.year}/top-tags"
			>
				<Tags />
				Top Tags
			</SectionHeader>
			<Spacer />

			<div class="flex flex-wrap gap-2">
				{#each data.topTags as tag (tag.tagSlug)}
					<div
						class="flex items-center gap-2 rounded-full border bg-card px-3 py-1.5"
					>
						<span class="text-sm">
							{tag.rank}. {formatTagSlug(tag.tagSlug)}
						</span>
						<span class="text-xs text-muted-foreground">
							{tag.playCount.toLocaleString()} plays
						</span>
					</div>
				{/each}
			</div>
		</section>

		<section>
			<SectionHeader
				count={data.topDecadeCount}
				viewAllHref="/users/{data.userData.id}/review/{data.year}/top-decades"
			>
				<CalendarRange />
				Decades
			</SectionHeader>
			<Spacer />

			<div class="flex flex-col gap-2">
				{#each data.topDecades as decade (decade.decade)}
					{@const count = decade.playCount}
					{@const maxDecadeCount = data.topDecades[0]?.playCount ?? 1}
					<div class="flex items-center gap-3">
						<span class="w-14 shrink-0 text-sm font-medium">
							{decade.decade}s
						</span>
						<div class="h-6 flex-1 overflow-hidden rounded bg-muted">
							<div
								class="h-full bg-linear-to-r from-logo-3 to-logo-1"
								style="width: {Math.max((count / maxDecadeCount) * 100, 4)}%"
							></div>
						</div>
						<span
							class="w-20 shrink-0 text-right text-xs text-muted-foreground"
						>
							{count.toLocaleString()} plays
						</span>
					</div>
				{/each}
			</div>
		</section>
	{/if}
</div>
