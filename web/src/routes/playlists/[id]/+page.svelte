<script lang="ts">
  import { goto, invalidateAll } from "$app/navigation";
  import { page } from "$app/stores";
  import { getApiClient, handleApiError } from "$lib";
  import ConfirmModal from "$lib/components/new-modals/ConfirmModal.svelte";
  import Spacer from "$lib/components/Spacer.svelte";
  import TrackList from "$lib/components/track-list/TrackList.svelte";
  import TrackListHeader from "$lib/components/track-list/TrackListHeader.svelte";
  import { getMusicManager } from "$lib/music-manager.svelte.js";
  import {
    Breadcrumb,
    Button,
    DropdownMenu,
    Pagination,
  } from "@nanoteck137/nano-ui";
  import { Pencil, Trash } from "lucide-svelte";
  import toast from "svelte-5-french-toast";
  import NewFilterModal from "./NewFilterModal.svelte";

  const { data } = $props();
  const musicManager = getMusicManager();
  const apiClient = getApiClient();

  let openFilterModal = $state(false);
  let openConfirmDeleteAlbum = $state(false);
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

<Button
  onclick={async () => {
    const res = await apiClient.generatePlaylistImage(data.playlist.id);
    if (!res.success) {
      return handleApiError(res.error);
    }
  }}
>
  Test Generate
</Button>

<Button
  onclick={async () => {
    openFilterModal = true;
  }}
>
  New Virtual Playlist
</Button>

<p>Track Count: {data.playlist.trackCount}</p>

<TrackListHeader
  name={data.playlist.name}
  image={data.playlist.coverArt.medium}
  onPlay={async (shuffle) => {
    await musicManager.queueRequest(
      { type: "addPlaylist", playlistId: data.playlist.id },
      { shuffle },
    );
  }}
>
  {#snippet more()}
    <DropdownMenu.Group>
      <DropdownMenu.Item onSelect={() => {}}>
        <Pencil />
        Edit Playlist
      </DropdownMenu.Item>

      <DropdownMenu.Item
        onSelect={() => {
          openConfirmDeleteAlbum = true;
        }}
      >
        <Trash />
        Delete Playlist
      </DropdownMenu.Item>
    </DropdownMenu.Group>
  {/snippet}
</TrackListHeader>

<Spacer size="md" />

{#each data.filters as filter}
  <button
    onclick={() => {
      const query = $page.url.searchParams;
      query.set("filterId", filter.filterId);
      goto("?" + query.toString(), {
        invalidateAll: true,
        replaceState: true,
      });
    }}
  >
    Filter: {filter.filterId} - {filter.name}
  </button>
{/each}

<Spacer size="md" />

<TrackList
  displayOrder
  totalTracks={data.page.totalItems}
  tracks={data.items}
  userPlaylists={data.userPlaylists}
  quickPlaylist={data.user?.quickPlaylist}
  onPlay={async (trackId) => {
    await musicManager.queueRequest(
      { type: "addPlaylist", playlistId: data.playlist.id },
      { queueIndexToTrackId: trackId },
    );
  }}
/>

<Spacer size="md" />

<Pagination.Root
  page={data.page.page + 1}
  count={data.page.totalItems}
  perPage={data.page.perPage}
  siblingCount={1}
  onPageChange={(p) => {
    const query = $page.url.searchParams;
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

<NewFilterModal bind:open={openFilterModal} playlistId={data.playlist.id} />

<ConfirmModal
  bind:open={openConfirmDeleteAlbum}
  removeTrigger
  confirmDelete
  onResult={async () => {
    const res = await apiClient.deletePlaylist(data.playlist.id);
    if (!res.success) {
      handleApiError(res.error);
      invalidateAll();
      return;
    }

    toast.success("Successfully deleted album");
    goto("/playlists", { invalidateAll: true });
  }}
/>
