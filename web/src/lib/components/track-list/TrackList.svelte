<script lang="ts">
  import {
    Button,
    buttonVariants,
    Checkbox,
    Dialog,
    DropdownMenu,
    Separator,
  } from "@nanoteck137/nano-ui";
  import TrackListItem from "./TrackListItem.svelte";
  import {
    ChevronDown,
    DiscAlbum,
    EllipsisVertical,
    Heart,
    Info,
    ListPlus,
    Star,
    Users,
    X,
  } from "lucide-svelte";
  import { getApiClient, handleApiError } from "$lib";
  import { showPlaylistModal } from "$lib/playlist-modal.svelte";
  import { cn } from "$lib/utils";
  import type { Playlist, Track } from "$lib/api/types";
  import { goto, invalidateAll } from "$app/navigation";
  import { getFavorites } from "$lib/favorites.svelte";
  import { getQuickPlaylist } from "$lib/quick-playlist.svelte";
  import FavoriteButton from "$lib/components/FavoriteButton.svelte";
  import QuickAddButton from "$lib/components/QuickAddButton.svelte";
  import toast from "svelte-5-french-toast";

  type Props = {
    totalTracks: number;
    tracks: Track[];

    isAlbumShowcase?: boolean;
    displayOrder?: boolean;

    userPlaylists?: Playlist[] | null;
    quickPlaylist?: string | null;

    onPlay: (trackId: string) => void;
    onReorder?: (items: string[], anchor: string | null) => void;
  };

  const {
    isAlbumShowcase,
    tracks,
    displayOrder,
    userPlaylists,
    quickPlaylist,
    onPlay,
    onReorder,
  }: Props = $props();
  const apiClient = getApiClient();
  const favoritesManager = getFavorites();
  const quickPlaylistManager = getQuickPlaylist();

  let selectedTracks = $state<string[]>([]);

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
</script>

<div class="flex flex-col">
  {#if selectedTracks.length > 0}
    <div class="flex h-14 items-center justify-end gap-1 px-2">
      <Button
        class="rounded-full"
        variant="ghost"
        size="icon-lg"
        onclick={() => {
          selectedTracks = [];
        }}
      >
        <X />
      </Button>

      <Button
        class="rounded-full"
        variant="default"
        onclick={() => {
          onReorder?.(selectedTracks, null);
          selectedTracks = [];
        }}
      >
        <ChevronDown />
        Insert after
      </Button>
    </div>
  {/if}

  <div class="flex flex-col">
    {#each tracks as track}
      <div class="group">
        <TrackListItem
          showNumber={isAlbumShowcase}
          {displayOrder}
          {track}
          onPlayClicked={() => {
            onPlay(track.id);
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
                    Select track
                  </DropdownMenu.Item>

                  <DropdownMenu.Item
                    class="sm:hidden"
                    onSelect={async () => {
                      const wasFav = favoritesManager.hasTrack(track.id);
                      await favoritesManager.toggleTrack(track.id);
                      toast.success(wasFav ? "Removed from favorites" : "Added to favorites");
                    }}
                  >
                    {#if favoritesManager.hasTrack(track.id)}
                      <Heart class="fill-primary" />
                      Unfavorite
                    {:else}
                      <Heart />
                      Favorite
                    {/if}
                  </DropdownMenu.Item>

                  {#if quickPlaylistManager.playlist !== null}
                    <DropdownMenu.Item
                      class="sm:hidden"
                      onSelect={async () => {
                        const wasIn = quickPlaylistManager.hasTrack(track.id);
                        await quickPlaylistManager.toggleTrack(track.id);
                        toast.success(wasIn ? "Removed from quick playlist" : "Added to quick playlist");
                      }}
                    >
                      {#if quickPlaylistManager.hasTrack(track.id)}
                        <Star class="fill-primary" />
                        Remove from Quick
                      {:else}
                        <Star />
                        Quick Add
                      {/if}
                    </DropdownMenu.Item>
                  {/if}

                  <DropdownMenu.Item
                    class="sm:hidden"
                    onSelect={() => {
                      goto(`/artists/${track.artists[0].id}`);
                    }}
                  >
                    <Users />
                    Go to Artist
                  </DropdownMenu.Item>
                  {#if !isAlbumShowcase}
                    <DropdownMenu.Item
                      onSelect={() => {
                        goto(`/albums/${track.albumId}`);
                      }}
                    >
                      <DiscAlbum />
                      Go to Album
                    </DropdownMenu.Item>
                  {/if}
                  <DropdownMenu.Item
                    onSelect={async () => {
                      const id = await showPlaylistModal();
                      if (!id) return;

                      const res = await apiClient.addItemToPlaylist(id, {
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

                  <DropdownMenu.Item
                    onSelect={() => {
                      infoTrackId = track.id;
                      infoOpen = true;
                    }}
                  >
                    <Info size={14} />
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
            <p class="text-lg font-semibold leading-tight">{infoTrack.name}</p>
            <p class="text-sm text-muted-foreground">
              {#each infoTrack.artists as artist, i}
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
                {#each infoTrack.tags as tag}
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
