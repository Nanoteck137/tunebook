<script lang="ts">
	import { goto } from "$app/navigation";
	import { page } from "$app/state";
	import {
		DiscAlbum,
		EllipsisVertical,
		Heart,
		History,
		Play,
		Shuffle,
		Star,
		X,
		Users,
	} from "@lucide/svelte";
	import Image from "$lib/components/Image.svelte";
	import Pagination from "$lib/components/Pagination.svelte";
	import { Button, DropdownMenu, Separator, buttonVariants } from "$lib/components/ui";
	import { getFavorites } from "$lib/favorites.svelte";
	import { getMusicManager } from "$lib/music-manager.svelte";
	import { getQuickPlaylist } from "$lib/quick-playlist.svelte";
	import { toast } from "svelte-sonner";

	let { data } = $props();

	const musicManager = getMusicManager();
	const favoritesManager = getFavorites();
	const quickPlaylistManager = getQuickPlaylist();

	let yearParam = $derived(page.url.searchParams.get("year"));

	function formatRelativeTime(millis: number): string {
		const unixSeconds = Math.floor(millis / 1000);
		const now = Math.floor(Date.now() / 1000);
		const diff = now - unixSeconds;

		if (diff < 60) return "just now";
		if (diff < 3600) {
			const m = Math.floor(diff / 60);
			return `${m}m ago`;
		}
		if (diff < 86400) {
			const h = Math.floor(diff / 3600);
			return `${h}h ago`;
		}
		if (diff < 604800) {
			const d = Math.floor(diff / 86400);
			return `${d}d ago`;
		}
		return new Date(unixSeconds * 1000).toLocaleDateString(undefined, {
			month: "short",
			day: "numeric",
		});
	}

	function statusLabel(status: string) {
		if (status === "completed") return "Completed";
		if (status === "skipped") return "Skipped";
		return "In Progress";
	}

	function statusClass(status: string) {
		if (status === "completed") return "bg-green-500/10 text-green-500 ring-green-500/25";
		if (status === "skipped") return "bg-muted text-muted-foreground ring-foreground/10";
		return "bg-yellow-500/10 text-yellow-500 ring-yellow-500/25";
	}

	function percentColor(pct: number): string {
		if (pct >= 80) return "bg-green-500";
		if (pct >= 40) return "bg-yellow-500";
		return "bg-muted-foreground/40";
	}

	function clearYearFilter() {
		const query = page.url.searchParams;
		query.delete("year");
		goto(`?${query.toString()}`, { invalidateAll: true });
	}

	async function playAll() {
		await musicManager.addTracks({
			trackIds: data.history.map((e) => e.track.id),
		});
	}

	async function shufflePlay() {
		const ids = data.history.map((e) => e.track.id);
		for (let i = ids.length - 1; i > 0; i--) {
			const j = Math.floor(Math.random() * (i + 1));
			[ids[i], ids[j]] = [ids[j], ids[i]];
		}
		await musicManager.addTracks({ trackIds: ids });
	}

	function playTrack(trackId: string) {
		const trackIds = data.history.map((e) => e.track.id);
		musicManager.addTracks({ trackIds, trackId });
	}
</script>

<div class="flex flex-col gap-6">
	<div
		class="flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between"
	>
		<div class="flex items-baseline gap-2">
			<h1 class="text-xl font-bold">
				{#if yearParam}
					History for {yearParam}
				{:else}
					Listening History
				{/if}
			</h1>
			<span class="text-sm text-muted-foreground">{data.page.totalItems}</span>
		</div>

		<div class="flex items-center gap-2">
			{#if yearParam}
				<Button variant="outline" size="sm" onclick={clearYearFilter}>
					<X size={14} />
					Clear year
				</Button>
			{/if}
			<Button size="sm" onclick={() => playAll()}>
				<Play size={14} />
				Play All
			</Button>
			<Button
				variant="outline"
				size="sm"
				onclick={() => shufflePlay()}
			>
				<Shuffle size={14} />
				Shuffle
			</Button>
		</div>
	</div>

	{#if data.history.length === 0}
		<div class="flex flex-col items-center gap-2 rounded-lg border py-16">
			<History size={32} class="text-muted-foreground/40" />
			<p class="text-sm text-muted-foreground">No listening history yet</p>
		</div>
	{:else}
		<div class="flex flex-col gap-1.5">
			{#each data.history as entry (entry.id)}
				<div
					class="group flex items-center gap-3 rounded-lg border bg-card p-2.5 transition-colors hover:bg-accent hover:text-accent-foreground has-data-[state='open']:bg-accent has-data-[state='open']:text-accent-foreground"
				>
					<button
						class="shrink-0 overflow-hidden rounded-md"
						onclick={() => playTrack(entry.track.id)}
						aria-label="Play {entry.track.name}"
					>
						<div class="relative h-12 w-12">
							<Image class="h-12 w-12" src={entry.track.coverArt.small} alt="" />
							<div
								class="absolute inset-0 flex items-center justify-center bg-black/55 opacity-0 transition-opacity group-hover:opacity-100"
							>
								<Play size={18} class="text-white" />
							</div>
						</div>
					</button>

					<div class="flex min-w-0 flex-1 flex-col gap-0.5">
						<div class="flex items-center gap-2">
							<span class="truncate text-sm font-medium" title={entry.track.name}>
								{entry.track.name}
							</span>
							{#if favoritesManager.hasTrack(entry.track.id)}
								<Heart
									size={12}
									class="shrink-0 fill-primary text-primary"
								/>
							{/if}
							{#if quickPlaylistManager.hasTrack(entry.track.id)}
								<Star size={12} class="shrink-0 fill-primary text-primary" />
							{/if}
							<span
								class="shrink-0 rounded-full px-2 py-px text-[10px] font-medium ring-1 {statusClass(
									entry.status,
								)}"
							>
								{statusLabel(entry.status)}
							</span>
						</div>

						<div class="flex min-w-0 items-center gap-1.5 text-xs text-muted-foreground">
							<span class="truncate" title={entry.track.artists.map((a) => a.name).join(", ")}>
								{#each entry.track.artists as artist, i (artist.id)}
									{#if i > 0}{", "}{/if}
									<a class="hover:underline" href="/artists/{artist.id}">
										{artist.name}
									</a>
								{/each}
							</span>
							{#if entry.track.albumName}
								<span class="shrink-0">&middot;</span>
								<span class="truncate">{entry.track.albumName}</span>
							{/if}
						</div>
					</div>

					<div class="hidden shrink-0 flex-col items-end gap-1 sm:flex">
						<span class="text-xs tabular-nums text-muted-foreground">
							{formatRelativeTime(entry.listenedAt)}
						</span>
						<div class="flex items-center gap-1.5">
								<div class="h-1.5 w-16 overflow-hidden rounded-full bg-muted">
									<div
										class="h-full rounded-full transition-all {percentColor(
											entry.percentPlayed,
										)}"
										style="width: {entry.percentPlayed}%"
									></div>
								</div>
								<span class="w-8 text-right text-[11px] tabular-nums text-muted-foreground"
									>{Math.round(entry.percentPlayed)}%</span
								>
							</div>
					</div>

					<span
						class="shrink-0 text-[11px] text-muted-foreground sm:hidden"
					>
						{formatRelativeTime(entry.listenedAt)}
					</span>

					<DropdownMenu.Root>
						<DropdownMenu.Trigger
							class={buttonVariants({ variant: "ghost", size: "icon" })}
							title="More options"
							aria-label="More options"
						>
							<EllipsisVertical />
						</DropdownMenu.Trigger>
						<DropdownMenu.Content align="end">
							<DropdownMenu.Group>
								<DropdownMenu.Item
									onSelect={() => playTrack(entry.track.id)}
								>
									<Play />
									Play
								</DropdownMenu.Item>
								<DropdownMenu.Item
									onSelect={() => {
										const trackIds = data.history.map((e) => e.track.id);
										for (let i = trackIds.length - 1; i > 0; i--) {
											const j = Math.floor(Math.random() * (i + 1));
											[trackIds[i], trackIds[j]] = [
												trackIds[j],
												trackIds[i],
											];
										}
										musicManager.addTracks({ trackIds });
									}}
								>
									<Shuffle />
									Shuffle play
								</DropdownMenu.Item>
							</DropdownMenu.Group>

							<DropdownMenu.Separator />

							<DropdownMenu.Item
								onSelect={() => goto(`/albums/${entry.track.albumId}`)}
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
									{#each entry.track.artists as artist (artist.id)}
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

							<DropdownMenu.Item
								onSelect={async () => {
									const wasFav = favoritesManager.hasTrack(
										entry.track.id,
									);
									await favoritesManager.toggleTrack(entry.track.id);
									toast.success(
										wasFav
											? "Removed from favorites"
											: "Added to favorites",
									);
								}}
							>
								{#if favoritesManager.hasTrack(entry.track.id)}
									<Heart class="fill-primary stroke-primary" />
									Unfavorite
								{:else}
									<Heart />
									Favorite
								{/if}
							</DropdownMenu.Item>
							{#if quickPlaylistManager.playlist !== null}
								<DropdownMenu.Item
									onSelect={async () => {
										const wasIn = quickPlaylistManager.hasTrack(
											entry.track.id,
										);
										await quickPlaylistManager.toggleTrack(
											entry.track.id,
										);
										toast.success(
											wasIn
												? "Removed from quick playlist"
												: "Added to quick playlist",
										);
									}}
								>
									{#if quickPlaylistManager.hasTrack(entry.track.id)}
										<Star class="fill-primary stroke-primary" />
										Remove from Quick
									{:else}
										<Star />
										Quick Add
									{/if}
								</DropdownMenu.Item>
							{/if}
						</DropdownMenu.Content>
					</DropdownMenu.Root>
				</div>
			{/each}
		</div>
	{/if}

	<Separator />

	<Pagination page={data.page} />
</div>