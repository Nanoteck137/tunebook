<script lang="ts">
  import {
    Button,
    buttonVariants,
    Checkbox,
    DropdownMenu,
    Separator,
  } from "@nanoteck137/nano-ui";
  import TrackListItem from "./TrackListItem.svelte";
  import {
    ChevronDown,
    DiscAlbum,
    EllipsisVertical,
    ListPlus,
    Users,
    X,
  } from "lucide-svelte";
  import { cn } from "$lib/utils";
  import type { Playlist, Track } from "$lib/api/types";
  import QuickAddButton from "$lib/components/QuickAddButton.svelte";
  import { goto, invalidateAll } from "$app/navigation";

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
    totalTracks,
    tracks,
    displayOrder,
    userPlaylists,
    quickPlaylist,
    onPlay,
    onReorder,
  }: Props = $props();

  let selectedTracks = $state<string[]>([]);
</script>

<div class="flex flex-col">
  <div class="flex items-center justify-between">
    <p class="text-bold text-xl">Tracks</p>
    <p class="text-sm">{totalTracks} track(s)</p>
  </div>

  {#if selectedTracks.length > 0}
    <div class="flex h-16 items-center justify-end bg-red-200 px-2">
      <Button
        class="rounded-full"
        variant="ghost"
        size="icon-lg"
        onclick={() => {
          onReorder?.(selectedTracks, null);
          selectedTracks = [];
        }}
      >
        <ChevronDown />
      </Button>

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
    </div>
  {/if}

  <div class="flex flex-col">
    {#each tracks as track}
      <div class="group">
        <TrackListItem
          class="group-even:bg-off-background2 group-even:hover:bg-off-background1"
          showNumber={isAlbumShowcase}
          {displayOrder}
          {track}
          onPlayClicked={() => {
            onPlay(track.id);
          }}
        >
          {#if selectedTracks.length > 0}
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

            <div class="flex h-11 w-11 items-center justify-center">
              <Checkbox
                class=""
                checked={selectedTracks.includes(track.id)}
                controlledChecked={true}
                onCheckedChange={(checked) => {
                  if (checked) {
                    selectedTracks.push(track.id);
                  } else {
                    selectedTracks = selectedTracks.filter(
                      (id) => track.id !== id,
                    );
                  }
                }}
              />
            </div>
          {/if}

          {#if selectedTracks.length <= 0}
            <QuickAddButton show={!!quickPlaylist} trackId={track.id} />

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
                      selectedTracks.push(track.id);
                    }}
                  >
                    Select track
                  </DropdownMenu.Item>

                  <DropdownMenu.Item
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
                      if (!userPlaylists) return;

                      // await openAddToPlaylist({
                      //   playlists: userPlaylists,
                      //   track,
                      // });

                      await invalidateAll();
                    }}
                  >
                    <ListPlus />
                    Save to Playlist
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
