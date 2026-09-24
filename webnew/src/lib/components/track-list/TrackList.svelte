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
	import SectionImage from "$lib/components/SectionImage.svelte";
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

											const res = await apiClient.addItemToPlaylist(
												playlist.id,
												{
													trackId: track.id,
												},
											);
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
	<Dialog.Content class="max-w-md gap-0 overflow-hidden p-0">
		{#if infoTrack}
			<div class="relative overflow-hidden">
				<img
					src={infoTrack.coverArt.original}
					alt=""
					aria-hidden="true"
					class="h-44 w-full scale-105 object-cover blur-sm"
				/>
				<div
					class="absolute inset-0 bg-linear-to-t from-black/90 via-black/40 to-black/30"
				></div>
				<div class="absolute inset-x-0 bottom-0 flex items-end gap-4 p-4">
					<SectionImage
						variant="full"
						src={infoTrack.coverArt.large}
						alt={infoTrack.name}
						class="section-tracks aspect-square w-20 shrink-0 rounded-md p-0.5 shadow-lg"
					/>
					<div class="min-w-0 flex-1 pb-0.5">
						<p
							class="text-lg leading-tight font-bold text-ellipsis text-white"
						>
							{infoTrack.name}
						</p>
						<p class="text-sm text-ellipsis text-white/80">
							{#each infoTrack.artists as artist, i (artist.id)}
								{#if i > 0}
									{", "}
								{/if}
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
				</div>
			</div>

			<div class="flex flex-col gap-4 p-4 pt-3">
				<div class="flex items-center gap-2">
					<Button
						size="sm"
						class="flex-1"
						onclick={() => onPlay(infoTrack.id, false)}
					>
						<Play />
						Play
					</Button>

					<Button
						size="sm"
						variant="outline"
						class="flex-1"
						onclick={() => onPlay(infoTrack.id, true)}
					>
						<Shuffle />
						Shuffle play
					</Button>
				</div>

				<div class="rounded-lg border bg-card p-3 text-sm">
					<div class="grid grid-cols-[auto_1fr] gap-x-4 gap-y-1.5">
						<span class="text-muted-foreground">Album</span>
						<a
							href="/albums/{infoTrack.albumId}"
							class="font-medium hover:underline"
						>
							{infoTrack.albumName}
						</a>

						{#if infoTrack.number}
							<span class="text-muted-foreground">Track</span>
							<span class="font-medium">#{infoTrack.number}</span>
						{/if}

						<span class="text-muted-foreground">Duration</span>
						<span class="font-medium"
							>{formatDuration(infoTrack.duration)}</span
						>

						{#if infoTrack.year}
							<span class="text-muted-foreground">Year</span>
							<span class="font-medium">{infoTrack.year}</span>
						{/if}

						<span class="text-muted-foreground">Added</span>
						<span class="font-medium">{formatDate(infoTrack.created)}</span>
						<span class="text-muted-foreground">Updated</span>
						<span class="font-medium">{formatDate(infoTrack.updated)}</span>
					</div>
				</div>

				{#if infoTrack.tags.length > 0}
					<div class="flex flex-wrap gap-1.5">
						{#each infoTrack.tags as tag (tag)}
							<span class="rounded-full bg-secondary px-2.5 py-0.5 text-xs">
								{tag}
							</span>
						{/each}
					</div>
				{/if}
			</div>
		{:else}
			<p class="p-4 text-sm text-muted-foreground">Track not found.</p>
		{/if}
	</Dialog.Content>
</Dialog.Root>
