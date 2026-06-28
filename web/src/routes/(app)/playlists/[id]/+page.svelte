<script lang="ts">
  import { goto, invalidateAll } from "$app/navigation";
  import { page } from "$app/state";
  import { getApiClient, handleApiError } from "$lib";
  import ConfirmModal from "$lib/components/new-modals/ConfirmModal.svelte";
  import Image from "$lib/components/Image.svelte";
  import TrackList from "$lib/components/track-list/TrackList.svelte";
  import { getMusicManager } from "$lib/music-manager.svelte.js";
  import {
    Breadcrumb,
    Button,
    buttonVariants,
    DropdownMenu,
    Pagination,
  } from "@nanoteck137/nano-ui";
  import {
    EllipsisVertical,
    ListPlus,
    Pencil,
    Play,
    Shuffle,
    Trash,
    Upload,
    Wand2,
  } from "lucide-svelte";
  import toast from "svelte-5-french-toast";
  import EditPlaylistModal from "./EditPlaylistModal.svelte";
  import UploadPlaylistCoverModal from "./UploadPlaylistCoverModal.svelte";

  const { data } = $props();
  const musicManager = getMusicManager();
  const apiClient = getApiClient();

  let openConfirmDelete = $state(false);
  let openEditPlaylistModal = $state(false);
  let openUploadCoverModal = $state(false);

  const isOwner = $derived(data.user?.id === data.playlist.ownerId);
</script>

<div class="py-2">
  <Breadcrumb.Root>
    <Breadcrumb.List>
      <Breadcrumb.Item>
        <Breadcrumb.Link href="/playlists">Playlists</Breadcrumb.Link>
      </Breadcrumb.Item>
      <Breadcrumb.Separator />
      <Breadcrumb.Item>
        <Breadcrumb.Page>{data.playlist.name}</Breadcrumb.Page>
      </Breadcrumb.Item>
    </Breadcrumb.List>
  </Breadcrumb.Root>
</div>

<div
  class="flex flex-col gap-6 rounded-lg border bg-gradient-to-b from-zinc-900 to-background p-4 sm:p-6 md:flex-row md:items-end md:gap-8"
>
  <Image
    class="w-40 min-w-40 self-center shadow-lg md:w-52 md:min-w-52"
    src={data.playlist.coverArt.large}
    alt={data.playlist.name}
  />

  <div class="flex min-w-0 flex-col gap-2">
    <p
      class="text-xs font-semibold uppercase tracking-wider text-muted-foreground"
    >
      Playlist
    </p>

    <h1 class="line-clamp-2 text-2xl font-bold md:text-4xl">
      {data.playlist.name}
    </h1>

    <div
      class="flex flex-wrap items-center gap-x-1 text-sm text-muted-foreground"
    >
      <span class="font-medium text-foreground"
        >{data.playlist.ownerDisplayName}</span
      >
      <span>&middot; {data.playlist.trackCount} tracks</span>
    </div>

    <div class="flex gap-2 pt-2">
      <Button
        size="sm"
        onclick={async () => {
          await musicManager.queueRequest(
            { type: "addPlaylist", playlistId: data.playlist.id },
            {},
          );
        }}
      >
        <Play size={14} />
        Play
      </Button>

      <Button
        variant="outline"
        size="sm"
        onclick={async () => {
          await musicManager.queueRequest(
            { type: "addPlaylist", playlistId: data.playlist.id },
            { shuffle: true },
          );
        }}
      >
        <Shuffle size={14} />
        Shuffle
      </Button>

      <DropdownMenu.Root>
        <DropdownMenu.Trigger
          class={buttonVariants({ variant: "outline", size: "icon" })}
        >
          <EllipsisVertical />
        </DropdownMenu.Trigger>
        <DropdownMenu.Content align="start">
          <DropdownMenu.Group>
            <DropdownMenu.Item
              onSelect={async () => {
                await musicManager.queueRequest(
                  { type: "addPlaylist", playlistId: data.playlist.id },
                  { append: "back" },
                );
              }}
            >
              <ListPlus />
              Append to Queue
            </DropdownMenu.Item>

            {#if isOwner}
              <DropdownMenu.Separator />

              <DropdownMenu.Item
                onSelect={() => {
                  openEditPlaylistModal = true;
                }}
              >
                <Pencil />
                Edit Playlist
              </DropdownMenu.Item>

              <DropdownMenu.Item
                onSelect={() => {
                  openUploadCoverModal = true;
                }}
              >
                <Upload />
                Upload Cover
              </DropdownMenu.Item>

              <DropdownMenu.Item
                onSelect={async () => {
                  const res = await apiClient.generatePlaylistImage(
                    data.playlist.id,
                  );
                  if (!res.success) {
                    handleApiError(res.error);
                  } else {
                    toast.success("Generating cover");
                  }
                }}
              >
                <Wand2 />
                Generate Cover
              </DropdownMenu.Item>

              <DropdownMenu.Separator />

              <DropdownMenu.Item
                onSelect={() => {
                  openConfirmDelete = true;
                }}
              >
                <Trash />
                Delete Playlist
              </DropdownMenu.Item>
            {/if}
          </DropdownMenu.Group>
        </DropdownMenu.Content>
      </DropdownMenu.Root>
    </div>
  </div>
</div>

<div class="h-4"></div>

<TrackList
  displayOrder
  totalTracks={data.page.totalItems}
  tracks={data.items}
  onPlay={async (trackId) => {
    await musicManager.queueRequest(
      { type: "addPlaylist", playlistId: data.playlist.id },
      { queueIndexToTrackId: trackId },
    );
  }}
  onReorder={async (items, anchor) => {
    const res = await apiClient.reorderPlaylistItems(data.playlist.id, {
      before: false,
      anchorTrackId: anchor ?? "",
      trackIds: items,
    });
    if (!res.success) {
      return handleApiError(res.error);
    }

    toast.success("Updated playlist");
    invalidateAll();
  }}
/>

<div class="h-4"></div>

<Pagination.Root
  page={data.page.page + 1}
  count={data.page.totalItems}
  perPage={data.page.perPage}
  siblingCount={1}
  onPageChange={(p) => {
    const query = page.url.searchParams;
    query.set("page", (p - 1).toString());

    goto(`?${query.toString()}`, { invalidateAll: true, keepFocus: true });
  }}
>
  {#snippet children({ pages, currentPage })}
    <Pagination.Content>
      <Pagination.Item>
        <Pagination.PrevButton />
      </Pagination.Item>
      {#each pages as page (page.key)}
        {#if page.type === "ellipsis"}
          <Pagination.Item>
            <Pagination.Ellipsis />
          </Pagination.Item>
        {:else}
          <Pagination.Item>
            <Pagination.Link
              href="?page={page.value}"
              {page}
              isActive={currentPage === page.value}
            >
              {page.value}
            </Pagination.Link>
          </Pagination.Item>
        {/if}
      {/each}
      <Pagination.Item>
        <Pagination.NextButton />
      </Pagination.Item>
    </Pagination.Content>
  {/snippet}
</Pagination.Root>

<ConfirmModal
  bind:open={openConfirmDelete}
  removeTrigger
  confirmDelete
  onResult={async () => {
    const res = await apiClient.deletePlaylist(data.playlist.id);
    if (!res.success) {
      handleApiError(res.error);
      invalidateAll();
      return;
    }

    toast.success("Deleted playlist");
    goto("/playlists", { invalidateAll: true });
  }}
/>

<EditPlaylistModal
  bind:open={openEditPlaylistModal}
  playlist={data.playlist}
/>

<UploadPlaylistCoverModal
  bind:open={openUploadCoverModal}
  playlistId={data.playlist.id}
/>
