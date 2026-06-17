<script lang="ts">
  import { Button } from "@nanoteck137/nano-ui";
  import TrackList from "$lib/components/track-list/TrackList.svelte";
  import { getMusicManager } from "$lib/music-manager.svelte";
  import TrackListHeader from "$lib/components/track-list/TrackListHeader.svelte";
  import Spacer from "$lib/components/Spacer.svelte";
  import NewFilterModal from "./NewFilterModal.svelte";
  import FilterButton from "./FilterButton.svelte";
  import Pagination from "$lib/components/Pagination.svelte";

  let { data } = $props();
  const musicManager = getMusicManager();

  let openNewFilterModal = $state(false);
</script>

<Button
  onclick={() => {
    openNewFilterModal = true;
  }}
>
  New Filter
</Button>

<!-- <form method="GET">
  <div class="flex flex-col gap-2">
    <Input
      type="text"
      name="filter"
      placeholder="Filter"
      value={data.filter ?? ""}
    />

    <Input
      type="text"
      name="sort"
      placeholder="Sort"
      value={data.sort ?? ""}
    />
  </div>

  {#if data.filterError}
    <p class="text-red-400">{data.filterError}</p>
  {/if}
  {#if data.sortError}
    <p class="text-red-400">{data.sortError}</p>
  {/if}
  <div class="h-2"></div>
  <Button type="submit">
    <Filter />
    Filter Tracks
  </Button>
</form>

<div class="h-2"></div> -->

<TrackListHeader
  name="Tracks"
  onPlay={async (shuffle) => {
    // await musicManager.queueRequest(
    //   {
    //     type: "addFilter",
    //     // filter: data.filter ?? "",
    //     filter: "",
    //   },
    //   { shuffle },
    // );
  }}
/>

<Spacer size="md" />

{#each data.filters as filter}
  <FilterButton {filter} />
{/each}

<Spacer size="md" />

<TrackList
  totalTracks={data.page.totalItems}
  tracks={data.tracks}
  userPlaylists={data.userPlaylists}
  quickPlaylist={data.user?.quickPlaylist}
  onPlay={async (trackId) => {
    await musicManager.addTracks({ trackId, clear: true });

    // await musicManager.queueRequest(
    //   {
    //     type: "addFilter",
    //     // filter: data.filter ?? "",
    //     filter: "",
    //   },
    //   { queueIndexToTrackId: trackId },
    // );
  }}
/>

<Spacer size="lg" />

<Pagination page={data.page} />

<NewFilterModal bind:open={openNewFilterModal} />
