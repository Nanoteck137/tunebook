<script lang="ts">
	import {
		Activity,
		CalendarDays,
		CalendarRange,
		ChevronLeft,
		ChevronRight,
		DiscAlbum,
		Heart,
		Music,
		Sparkles,
		Tags,
		Users,
	} from "@lucide/svelte";
	import AlbumTile from "$lib/components/tiles/AlbumTile.svelte";
	import ArtistTile from "$lib/components/tiles/ArtistTile.svelte";
	import SectionHeader from "$lib/components/SectionHeader.svelte";
	import Spacer from "$lib/components/Spacer.svelte";
	import { Breadcrumb, Button } from "$lib/components/ui";
	import { formatPlayTime } from "$lib/utils";

	let { data } = $props();

	const monthNames = [
		"January",
		"February",
		"March",
		"April",
		"May",
		"June",
		"July",
		"August",
		"September",
		"October",
		"November",
		"December",
	];

	let currentYear = $derived(new Date().getFullYear());
	let monthName = $derived(monthNames[data.monthNum - 1]);
	let isCurrentMonth = $derived(
		data.year === currentYear && data.monthNum === new Date().getMonth() + 1,
	);

	let skipRate = $derived(
		data.month.playCount > 0
			? (data.month.skipCount / data.month.playCount) * 100
			: 0,
	);
	let favoritesOverlap = $derived(
		data.month.playCount > 0
			? (data.month.favoritePlays / data.month.playCount) * 100
			: 0,
	);
	let repeatRate = $derived(
		data.month.uniqueTracks > 0
			? data.month.playCount / data.month.uniqueTracks
			: 0,
	);
	let annualShare = $derived(
		data.review.review.trackCount > 0
			? (data.month.playCount / data.review.review.trackCount) * 100
			: 0,
	);

	function formatPercent(value: number): string {
		return `${value.toLocaleString("en-US", { maximumFractionDigits: 1 })}%`;
	}

	function formatTagSlug(slug: string): string {
		const words = slug.split("-");
		const sentence = words.join(" ");
		return sentence.charAt(0).toUpperCase() + sentence.slice(1);
	}
</script>

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
				<Breadcrumb.Link href="/users/{data.userData.id}/review/{data.year}">
					{data.year}
				</Breadcrumb.Link>
			</Breadcrumb.Item>
			<Breadcrumb.Separator />
			<Breadcrumb.Item>
				<Breadcrumb.Page>{monthName}</Breadcrumb.Page>
			</Breadcrumb.Item>
		</Breadcrumb.List>
	</Breadcrumb.Root>

	<SectionHeader count={data.month.playCount}>
		<CalendarDays />
		{monthName}

		{#if isCurrentMonth}
			<span
				class="rounded-full bg-primary/10 px-2.5 py-0.5 text-xs font-medium text-primary ring-1 ring-primary/25"
			>
				In progress
			</span>
		{/if}

		{#snippet actions()}
			<div class="flex items-center gap-1">
				<Button
					variant="ghost"
					size="icon-sm"
					disabled={data.monthNum <= 1}
					href={data.monthNum > 1
						? `/users/${data.userData.id}/review/${data.year}/months/${data.monthNum - 1}`
						: undefined}
					aria-label="Previous month"
				>
					<ChevronLeft />
				</Button>
				<Button
					variant="ghost"
					size="icon-sm"
					disabled={data.monthNum >= 12}
					href={data.monthNum < 12
						? `/users/${data.userData.id}/review/${data.year}/months/${data.monthNum + 1}`
						: undefined}
					aria-label="Next month"
				>
					<ChevronRight />
				</Button>
			</div>
		{/snippet}
	</SectionHeader>

	{#if data.month.playCount > 0}
		<p class="text-sm text-muted-foreground">
			{formatPlayTime(data.month.playTime)} listened across
			{data.month.uniqueTracks.toLocaleString()} unique tracks
		</p>

		<section>
			<SectionHeader>
				<Activity />
				Listening Stats
			</SectionHeader>
			<Spacer />

			<div class="grid grid-cols-2 gap-3 md:grid-cols-4">
				{@render statCard("Plays", data.month.playCount.toLocaleString())}
				{@render statCard(
					"Listening Time",
					formatPlayTime(data.month.playTime),
				)}
				{@render statCard(
					"Unique Tracks",
					data.month.uniqueTracks.toLocaleString(),
				)}
				{@render statCard("Skipped", data.month.skipCount.toLocaleString())}
				{@render statCard(
					"Avg Completion",
					formatPercent(data.month.avgCompletion),
				)}
				{@render statCard("Skip Rate", formatPercent(skipRate))}
				{@render statCard("Repeat Rate", `${repeatRate.toFixed(1)}x`)}
				{@render statCard("Annual Share", formatPercent(annualShare))}
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

		{#if data.topTracks.length > 0}
			<section>
				<SectionHeader
					count={data.topTrackCount}
					viewAllHref="/users/{data.userData
						.id}/review/{data.year}/months/{data.monthNum}/top-tracks"
				>
					<Music />
					Top Tracks
				</SectionHeader>
				<Spacer />

				<div class="flex flex-col">
					{#each data.topTracks as item (item.id)}
						<div class="flex items-center gap-3 rounded-lg border bg-card p-2">
							<span
								class="w-7 shrink-0 text-center text-sm font-bold text-muted-foreground"
							>
								{item.rank}
							</span>
							<img
								src={item.coverArt.small}
								alt=""
								class="size-10 shrink-0 rounded object-cover"
							/>
							<span class="min-w-0 flex-1">
								<span class="block truncate text-sm font-medium">
									{item.name}
								</span>
								<span class="block truncate text-xs text-muted-foreground">
									{item.artists.map((a) => a.name).join(", ")}
								</span>
							</span>
							<span class="shrink-0 text-xs text-muted-foreground">
								{item.playCount.toLocaleString()} plays
							</span>
						</div>
					{/each}
				</div>
			</section>
		{/if}

		{#if data.topAlbums.length > 0}
			<section>
				<SectionHeader
					count={data.topAlbumCount}
					viewAllHref="/users/{data.userData
						.id}/review/{data.year}/months/{data.monthNum}/top-albums"
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
		{/if}

		{#if data.topArtists.length > 0}
			<section>
				<SectionHeader
					count={data.topArtistCount}
					viewAllHref="/users/{data.userData
						.id}/review/{data.year}/months/{data.monthNum}/top-artists"
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
		{/if}

		{#if data.topTags.length > 0}
			<section>
				<SectionHeader
					count={data.topTagCount}
					viewAllHref="/users/{data.userData
						.id}/review/{data.year}/months/{data.monthNum}/top-tags"
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
		{/if}

		{#if data.topDecades.length > 0}
			<section>
				<SectionHeader
					count={data.topDecadeCount}
					viewAllHref="/users/{data.userData
						.id}/review/{data.year}/months/{data.monthNum}/top-decades"
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
	{:else}
		<div
			class="flex flex-col items-center gap-2 rounded-lg border py-16 text-center"
		>
			<Sparkles size={32} class="text-muted-foreground/40" />
			<p class="text-sm font-medium">No listening data for {monthName}</p>
			<p class="max-w-sm text-sm text-muted-foreground">
				This month has no plays in the {data.year} review. Use the arrows to browse
				the other months.
			</p>
		</div>
	{/if}
</div>
