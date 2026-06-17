<script lang="ts">
  import {
    Button,
    buttonVariants,
    DropdownMenu,
  } from "@nanoteck137/nano-ui";
  import { Check, EllipsisVertical, FileHeart, Plus } from "lucide-svelte";
  import { cn } from "$lib/utils";
  import { invalidateAll } from "$app/navigation";
  import { getApiClient, handleApiError } from "$lib";
  import NewPlaylistModal from "./NewPlaylistModal.svelte";
  import PlaylistTile from "$lib/components/tiles/PlaylistTile.svelte";
  import Spacer from "$lib/components/Spacer.svelte";
  import Pagination from "$lib/components/Pagination.svelte";
  import { Separator } from "@nanoteck137/nano-ui";

  let { data } = $props();
  const apiClient = getApiClient();

  let openNewPlaylistModal = $state(false);
</script>

<div class="flex flex-col gap-4">
  <div class="flex items-center justify-between">
    <div class="flex items-center gap-2">
      <p class="text-bold text-xl">Playlists</p>
      {#if data.page}
        <span class="text-sm text-muted-foreground">
          ({data.page.totalItems} playlists)
        </span>
      {/if}
    </div>
    <Button
      variant="ghost"
      onclick={() => {
        openNewPlaylistModal = true;
      }}
    >
      <Plus />
      New Playlist
    </Button>
  </div>
</div>

<Spacer size="lg" />

<div class="flex flex-shrink flex-wrap justify-center gap-4">
  {#each data.playlists as playlist}
    <div class="relative flex shrink-0 flex-col items-center">
      {#if data.user?.quickPlaylist === playlist.id}
        <div
          class="absolute left-2 top-2 z-10 flex items-center gap-1 rounded-full bg-primary px-2 py-0.5 text-xs text-primary-foreground"
        >
          <Check size={12} />
          Quick
        </div>
      {/if}

      <PlaylistTile
        id={playlist.id}
        cover={playlist.coverArt.medium}
        name={playlist.name}
        trackCount={playlist.trackCount}
      />

      <div class="absolute right-2 top-2">
        <DropdownMenu.Root>
          <DropdownMenu.Trigger
            class={cn(
              buttonVariants({ variant: "secondary", size: "icon" }),
              "rounded-full opacity-70 hover:opacity-100",
            )}
          >
            <EllipsisVertical size={14} />
          </DropdownMenu.Trigger>
          <DropdownMenu.Content align="end">
            <DropdownMenu.Group>
              {#if data.user?.quickPlaylist !== playlist.id}
                <DropdownMenu.Item
                  onSelect={async () => {
                    const res = await apiClient.setQuickPlaylist({
                      playlistId: playlist.id,
                    });
                    if (!res.success) {
                      handleApiError(res.error);
                      return;
                    }

                    await invalidateAll();
                  }}
                >
                  <FileHeart />
                  Set as Quick Playlist
                </DropdownMenu.Item>
              {/if}
            </DropdownMenu.Group>
          </DropdownMenu.Content>
        </DropdownMenu.Root>
      </div>
    </div>
  {/each}
</div>

<Spacer size="lg" />

<Separator />

<Spacer size="lg" />

<Pagination page={data.page} />

<NewPlaylistModal bind:open={openNewPlaylistModal} />
