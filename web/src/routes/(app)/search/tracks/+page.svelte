<script lang="ts">
  import { goto } from "$app/navigation";
  import { Button, Input } from "@nanoteck137/nano-ui";
  import { Search } from "lucide-svelte";
  import Pagination from "$lib/components/Pagination.svelte";
  import TrackList from "$lib/components/track-list/TrackList.svelte";
  import { getMusicManager } from "$lib/music-manager.svelte.js";
  import { onMount } from "svelte";

  let { data } = $props();

  const musicManager = getMusicManager();

  async function doSearch(query: string) {
    await goto(`/search/tracks?query=${query}`, {
      invalidateAll: true,
      keepFocus: true,
      replaceState: true,
    });
  }

  let initialValue = $state("");
  let value = "";

  onMount(() => {
    initialValue = data.query;
    value = data.query;
  });

  let timer: ReturnType<typeof setTimeout>;
  function onInput(e: Event) {
    const target = e.target as HTMLInputElement;
    const current = target.value;
    value = current;

    clearTimeout(timer);
    timer = setTimeout(() => {
      doSearch(current);
    }, 500);
  }
</script>

<svelte:head>
  <title>Search Tracks - Tunebook</title>
</svelte:head>

<div class="flex flex-col gap-6">
  <div class="flex flex-col gap-2">
    <h1 class="text-xl font-bold">Search Tracks</h1>

    <form
      action=""
      method="get"
      onsubmit={(e) => {
        e.preventDefault();
        clearTimeout(timer);
        doSearch(value);
      }}
    >
      <div class="flex items-center gap-2">
        <div class="relative flex-1">
          <Search
            size={16}
            class="pointer-events-none absolute left-3 top-1/2 -translate-y-1/2 text-muted-foreground"
          />
          <Input
            id="query"
            name="query"
            placeholder="Search tracks..."
            autocomplete="off"
            value={initialValue}
            oninput={onInput}
            class="pl-9"
          />
        </div>
        <Button type="submit">Search</Button>
      </div>
    </form>
  </div>

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
      <Pagination page={data.page} />
    {/if}
  {/if}
</div>
