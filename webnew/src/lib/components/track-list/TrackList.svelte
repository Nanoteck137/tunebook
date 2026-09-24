<script lang="ts">
	import TrackListItem from "./TrackListItem.svelte";
	import {
		ChevronDown,
		DiscAlbum,
		EllipsisVertical,
		Heart,
		Info,
		ListChecks,
		ListPlus,
		Play,
		Shuffle,
		Star,
		User,
		X,
	} from "@lucide/svelte";
	import { getApiClient, handleApiError } from "$lib";
	import { showPlaylistModal } from "$lib/playlist-modal.svelte";
	import { cn } from "$lib/utils";
	import type { Track } from "$lib/api/types";
	import { goto, invalidateAll } from "$app/navigation";
	import { getFavorites } from "$lib/favorites.svelte";
	import { getQuickPlaylist } from "$lib/quick-playlist.svelte";
	import FavoriteButton from "$lib/components/FavoriteButton.svelte";
	import QuickAddButton from "$lib/components/QuickAddButton.svelte";
	import { toast } from "svelte-sonner";
	import {
		Dialog,
		DropdownMenu,
		Button,
		Checkbox,
		buttonVariants,
		Separator,
	} from "$lib/components/ui";

	type Props = {
		totalTracks: number;
		tracks: Track[];

		isAlbumShowcase?: boolean;
		displayOrder?: boolean;
		highlightId?: string | null;

		selectedTracks?: string[];

		// eslint-disable-next-line no-unused-vars
		onPlay: (trackId: string, shuffle: boolean) => void;
		// eslint-disable-next-line no-unused-vars
		onReorder?: (items: string[], anchor: string | null) => void;
	};

	let {
		isAlbumShowcase,
		tracks,
		displayOrder,
		highlightId,
		selectedTracks = $bindable([]),
		onPlay,
		onReorder,
	}: Props = $props();
	const apiClient = getApiClient();
	const favoritesManager = getFavorites();
	const quickPlaylistManager = getQuickPlaylist();

	let infoTrackId = $state<string | null>(null);
	let infoOpen = $state(false);

	let infoTrack = $derived(
		infoTrackId ? (tracks.find((t) => t.id === infoTrackId) ?? null) : null,
	);

	function formatDuration(seconds: number) {
		const m = Math.floor(seconds / 60);
		const s = seconds % 60;
		return `${m}:${s.toString().padStart(2, "0")}`;
	}

	function formatDate(iso: string) {
		return new Date(iso).toLocaleDateString(undefined, {
			year: "numeric",
			month: "short",
			day: "numeric",
		});
	}

	let allSelectedFavorited = $derived(
		selectedTracks.length > 0 &&
			selectedTracks.every((id) => favoritesManager.hasTrack(id)),
	);

	function toggleTrackSelection(id: string) {
		if (selectedTracks.includes(id)) {
			selectedTracks = selectedTracks.filter((v) => v !== id);
		} else {
			selectedTracks = [...selectedTracks, id];
		}
	}

	async function favoriteSelectedTracks() {
		const ids = selectedTracks;
		selectedTracks = [];

		if (allSelectedFavorited) {
			await favoritesManager.unfavoriteTracks(ids);
		} else {
			await favoritesManager.favoriteTracks(ids);
		}
	}

	async function addSelectedToPlaylist() {
		const ids = selectedTracks;
		const playlist = await showPlaylistModal();
		if (!playlist) return;

		const t = toast.success(`Adding (0 / ${ids.length}) to ${playlist.name}`);

		let num = 0;

		for (const id of ids) {
			const res = await apiClient.addItemToPlaylist(playlist.id, {
				trackId: id,
			});
			if (!res.success) {
				if (res.error.type !== "PLAYLIST_ALREADY_HAS_TRACK") {
					handleApiError(res.error);
					return;
				}
			}

			num++;

			toast.info(`Adding (${num} / ${ids.length}) to ${playlist.name}`, {
				id: t,
			});
		}

		toast.success(
			`Added ${ids.length} track${ids.length === 1 ? "" : "s"} to ${playlist.name}`,
			{
				id: t,
			},
		);
		selectedTracks = [];
	}
</script>

<div class="flex flex-col">
	{#if selectedTracks.length > 0}
		<div
			class="flex h-12 shrink-0 items-center gap-1 rounded-lg border border-dashed border-border/60 px-1 text-xs font-medium text-muted-foreground"
		>
			<button
				class="rounded-full p-1.5 transition-colors hover:bg-muted hover:text-foreground"
				onclick={() => {
					selectedTracks = [];
				}}
				aria-label="Clear selection"
				title="Clear selection"
			>
				<X size={16} />
			</button>

			{#if onReorder}
				<button
					class="group flex h-full min-w-0 flex-1 items-center justify-center gap-2 rounded-md border-2 border-dotted border-transparent transition-colors hover:border-primary/60 hover:bg-accent/50 hover:text-foreground"
					onclick={() => {
						onReorder(selectedTracks, null);
						selectedTracks = [];
					}}
				>
					<ChevronDown
						class="text-muted-foreground/60 transition-colors group-hover:text-foreground"
					/>
					Place {selectedTracks.length} here
				</button>
			{:else}
				<span class="flex-1 text-center">
					{selectedTracks.length}
					{selectedTracks.length === 1 ? "track" : "tracks"} selected
				</span>
			{/if}

			<DropdownMenu.Root>
				<DropdownMenu.Trigger
					class={buttonVariants({
						variant: "ghost",
						size: "icon-lg",
						class: "rounded-full",
					})}
					title="Actions"
					aria-label="Actions"
				>
					<EllipsisVertical />
				</DropdownMenu.Trigger>
				<DropdownMenu.Content align="end">
					<DropdownMenu.Group>
						<DropdownMenu.Item onSelect={favoriteSelectedTracks}>
							<Heart />
							{allSelectedFavorited
								? "Remove from favorites"
								: "Add to favorites"}
						</DropdownMenu.Item>
						<DropdownMenu.Item onSelect={addSelectedToPlaylist}>
							<ListPlus />
							Save to Playlist
						</DropdownMenu.Item>
					</DropdownMenu.Group>
				</DropdownMenu.Content>
			</DropdownMenu.Root>
		</div>
	{/if}

	<div class="flex flex-col">
		{#each tracks as track (track.id)}
			<div class="group" id="track-{track.id}">
				<TrackListItem
					class={highlightId === track.id
						? "row-flash row-highlight"
						: undefined}
					showNumber={isAlbumShowcase}
					{displayOrder}
					{track}
					selectionMode={selectedTracks.length > 0}
					selected={selectedTracks.includes(track.id)}
					onSelect={() => toggleTrackSelection(track.id)}
					onPlayClicked={() => {
						onPlay(track.id, false);
					}}
				>
					{#if selectedTracks.length > 0}
						<div class="flex h-11 w-11 items-center justify-center">
							<Checkbox
								checked={selectedTracks.includes(track.id)}
								onCheckedChange={(checked) => {
									if (checked) {
										selectedTracks = [...selectedTracks, track.id];
									} else {
										selectedTracks = selectedTracks.filter(
											(id) => track.id !== id,
										);
									}
								}}
							/>
						</div>

						<Button
							class="rounded-full"
							variant="ghost"
							size="icon-lg"
							onclick={() => {
								onReorder?.(selectedTracks, track.id);
								selectedTracks = [];
							}}
						>
							<ChevronDown />
						</Button>
					{/if}

					{#if selectedTracks.length <= 0}
						<div class="hidden sm:flex sm:items-center sm:gap-0.5">
							<FavoriteButton show trackId={track.id} />
							<QuickAddButton trackId={track.id} />
						</div>

						<DropdownMenu.Root>
							<DropdownMenu.Trigger
								class={cn(
									buttonVariants({ variant: "ghost", size: "icon-lg" }),
									"rounded-full",
								)}
							>
								<EllipsisVertical />
							</DropdownMenu.Trigger>
							<DropdownMenu.Content align="end">
								<DropdownMenu.Group>
									<DropdownMenu.Item
										onSelect={() => {
											selectedTracks = [...selectedTracks, track.id];
										}}
									>
										<ListChecks />
										Select track
									</DropdownMenu.Item>
								</DropdownMenu.Group>

								<DropdownMenu.Separator />

								<DropdownMenu.Group>
									<DropdownMenu.Item
										onSelect={() => {
											onPlay(track.id, false);
										}}
									>
										<Play />
										Play
									</DropdownMenu.Item>

									<DropdownMenu.Item
										onSelect={() => {
											onPlay(track.id, true);
										}}
									>
										<Shuffle />
										Shuffle play
									</DropdownMenu.Item>
								</DropdownMenu.Group>

								<DropdownMenu.Separator />

								{#if !isAlbumShowcase}
									<DropdownMenu.Item
										onSelect={() => {
											goto(`/albums/${track.albumId}?track=${track.id}`);
										}}
									>
										<DiscAlbum />
										Go to Album
									</DropdownMenu.Item>
								{/if}

								<DropdownMenu.Sub>
									<DropdownMenu.SubTrigger>
										<User />
										Go to artist
									</DropdownMenu.SubTrigger>
									<DropdownMenu.SubContent>
										{#each track.artists as artist (artist.id)}
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

								<DropdownMenu.Group>
									<DropdownMenu.Item
										onSelect={async () => {
											const wasFav = favoritesManager.hasTrack(track.id);
											await favoritesManager.toggleTrack(track.id);
											toast.success(
												wasFav
													? "Removed from favorites"
													: "Added to favorites",
											);
										}}
									>
										{#if favoritesManager.hasTrack(track.id)}
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
												const wasIn = quickPlaylistManager.hasTrack(track.id);
												await quickPlaylistManager.toggleTrack(track.id);
												toast.success(
													wasIn
														? "Removed from quick playlist"
														: "Added to quick playlist",
												);
											}}
										>
											{#if quickPlaylistManager.hasTrack(track.id)}
												<Star class="fill-primary stroke-primary" />
												Remove from Quick
											{:else}
												<Star />
												Quick Add
											{/if}
										</DropdownMenu.Item>
									{/if}

									<DropdownMenu.Item
										onSelect={async () => {
											const playlist = await showPlaylistModal();
											if (!playlist) return;

											const res = await apiClient.addItemToPlaylist(playlist.id, {
												trackId: track.id,
											});
											if (!res.success) {
												handleApiError(res.error);
												return;
											}

											await invalidateAll();
										}}
									>
										<ListPlus />
										Save to Playlist
									</DropdownMenu.Item>
								</DropdownMenu.Group>

								<DropdownMenu.Separator />

								<DropdownMenu.Group>
									<DropdownMenu.Item
										onSelect={() => {
											infoTrackId = track.id;
											infoOpen = true;
										}}
									>
										<Info />
										Show more info
									</DropdownMenu.Item>
								</DropdownMenu.Group>
							</DropdownMenu.Content>
						</DropdownMenu.Root>
					{/if}
				</TrackListItem>

				<Separator />
			</div>
		{/each}
	</div>
</div>

<Dialog.Root open={infoOpen} onOpenChange={(v) => (infoOpen = v)}>
	<Dialog.Content class="sm:max-w-lg">
		<Dialog.Header>
			<Dialog.Title>Track Info</Dialog.Title>
			<Dialog.Description>
				Detailed information about the track
			</Dialog.Description>
		</Dialog.Header>

		{#if infoTrack}
			<div class="flex flex-col gap-4 sm:flex-row">
				<div class="flex shrink-0 justify-center sm:block">
					<img
						src={infoTrack.coverArt.large}
						alt={infoTrack.name}
						class="h-48 w-48 rounded-lg border object-cover sm:h-44 sm:w-44"
					/>
				</div>

				<div class="flex min-w-0 flex-1 flex-col gap-2">
					<div>
						<p class="text-lg leading-tight font-semibold">{infoTrack.name}</p>
						<p class="text-sm text-muted-foreground">
							{#each infoTrack.artists as artist, i (artist.id)}
								{#if i > 0}{", "}{/if}
								<a
									href="/artists/{artist.id}"
									class="hover:underline"
									title={artist.name}
								>
									{artist.name}
								</a>
							{/each}
						</p>
					</div>

					<div class="grid grid-cols-[auto_1fr] gap-x-3 gap-y-1 text-sm">
						<span class="text-muted-foreground">Album</span>
						<a href="/albums/{infoTrack.albumId}" class="hover:underline">
							{infoTrack.albumName}
						</a>

						{#if infoTrack.number}
							<span class="text-muted-foreground">Track</span>
							<span>#{infoTrack.number}</span>
						{/if}

						<span class="text-muted-foreground">Duration</span>
						<span>{formatDuration(infoTrack.duration)}</span>

						{#if infoTrack.year}
							<span class="text-muted-foreground">Year</span>
							<span>{infoTrack.year}</span>
						{/if}

						{#if infoTrack.tags.length > 0}
							<span class="text-muted-foreground">Tags</span>
							<div class="flex flex-wrap gap-1">
								{#each infoTrack.tags as tag (tag)}
									<span class="rounded-md bg-secondary px-1.5 py-0.5 text-xs"
										>{tag}</span
									>
								{/each}
							</div>
						{/if}

						<span class="text-muted-foreground">Added</span>
						<span>{formatDate(infoTrack.created)}</span>

						<span class="text-muted-foreground">Updated</span>
						<span>{formatDate(infoTrack.updated)}</span>
					</div>
				</div>
			</div>
		{:else}
			<p class="text-sm text-muted-foreground">Track not found.</p>
		{/if}

		<Dialog.Footer>
			<Button variant="outline" onclick={() => (infoOpen = false)}>
				Close
			</Button>
		</Dialog.Footer>
	</Dialog.Content>
</Dialog.Root>
