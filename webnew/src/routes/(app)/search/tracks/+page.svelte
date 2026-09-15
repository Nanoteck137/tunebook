<script lang="ts">
  import { goto } from "$app/navigation";
  import { Separator } from "$lib/components/ui";
  import Pagination from "$lib/components/Pagination.svelte";
  import Spacer from "$lib/components/Spacer.svelte";
  import TrackList from "$lib/components/track-list/TrackList.svelte";
  import { onMount } from "svelte";
  import SearchBarHeader from "../SearchBarHeader.svelte";

  let { data } = $props();

  async function doSearch(query: string) {
    await goto(`/search/tracks?query=${query}`, {
      invalidateAll: true,
      keepFocus: true,
      replaceState: true,
    });
  }

  function clearSearch() {
    value = "";
    doSearch("");
  }

  let value = $state("");

  onMount(() => {
    value = data.query;
  });
</script>

<svelte:head>
  <title>Search Tracks - Tunebook</title>
</svelte:head>

<div class="flex flex-col gap-6">
  <SearchBarHeader
    searchBarPlaceholder="Search tracks..."
    {value}
    setValue={(v) => {
      value = v;
    }}
    search={doSearch}
    searchWithValue={() => {
      doSearch(value);
    }}
    {clearSearch}
  />

  {#if data.query && data.tracks.length === 0}
    <p class="py-12 text-center text-sm text-muted-foreground">
      No tracks found for "{data.query}".
    </p>
  {/if}

  {#if data.tracks.length > 0}
    {#if data.page}
      <div class="flex items-baseline gap-2">
        <span class="text-sm text-muted-foreground">
          {data.page.totalItems} track(s)
        </span>
      </div>
    {/if}

    <TrackList
      totalTracks={data.tracks.length}
      tracks={data.tracks}
      onPlay={() => {}}
    />

    {#if data.page}
      <Spacer size="lg" />
      <Separator />
      <Spacer size="lg" />

      <Pagination page={data.page} />
    {/if}
  {/if}
</div>
