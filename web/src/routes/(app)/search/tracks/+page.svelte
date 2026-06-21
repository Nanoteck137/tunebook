<script lang="ts">
  import { goto } from "$app/navigation";
  import { page } from "$app/state";
  import { Button, Input } from "@nanoteck137/nano-ui";
  import { Search, X } from "lucide-svelte";
  import Pagination from "$lib/components/Pagination.svelte";
  import TrackList from "$lib/components/track-list/TrackList.svelte";
  import { getMusicManager } from "$lib/music-manager.svelte.js";
  import { cn } from "$lib/utils";
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

  function clearSearch() {
    value = "";
    doSearch("");
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

  const tabs = [
    { label: "All", href: "/search" },
    { label: "Tracks", href: "/search/tracks" },
    { label: "Artists", href: "/search/artists" },
    { label: "Albums", href: "/search/albums" },
    { label: "Playlists", href: "/search/playlists" },
    { label: "Users", href: "/search/users" },
  ];
</script>

<svelte:head>
  <title>Search Tracks - Tunebook</title>
</svelte:head>

<div class="flex flex-col gap-6">
  <div>
    <h1 class="mb-4 text-xl font-bold">Search Tracks</h1>

    <form
      action=""
      method="get"
      onsubmit={(e) => {
        e.preventDefault();
        clearTimeout(timer);
        doSearch(value);
      }}
    >
      <div class="rounded-lg border bg-card p-3">
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
          {#if data.query}
            <Button variant="ghost" size="icon" onclick={clearSearch}>
              <X size={16} />
            </Button>
          {/if}
        </div>
      </div>
    </form>
  </div>

  <nav class="flex flex-wrap gap-1">
    {#each tabs as { label, href }}
      <a
        href="{href}?query={data.query}"
        class={cn(
          "rounded-md px-3 py-1.5 text-sm font-medium transition-colors",
          page.url.pathname === href
            ? "bg-primary text-primary-foreground"
            : "text-muted-foreground hover:bg-accent hover:text-accent-foreground",
        )}
      >
        {label}
      </a>
    {/each}
  </nav>

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
