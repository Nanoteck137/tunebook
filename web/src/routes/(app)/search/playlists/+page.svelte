<script lang="ts">
  import { goto } from "$app/navigation";
  import { Button, Input } from "@nanoteck137/nano-ui";
  import { Search } from "lucide-svelte";
  import Pagination from "$lib/components/Pagination.svelte";
  import Image from "$lib/components/Image.svelte";
  import { onMount } from "svelte";

  let { data } = $props();

  async function doSearch(query: string) {
    await goto(`/search/playlists?query=${query}`, {
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
  <title>Search Playlists - Tunebook</title>
</svelte:head>

<div class="flex flex-col gap-6">
  <div class="flex flex-col gap-2">
    <h1 class="text-xl font-bold">Search Playlists</h1>

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
            placeholder="Search playlists..."
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

  {#if data.query && data.playlists.length === 0}
    <p class="py-12 text-center text-sm text-muted-foreground">
      No playlists found for "{data.query}".
    </p>
  {/if}

  {#if data.playlists.length > 0}
    {#if data.page}
      <div class="flex items-baseline gap-2">
        <span class="text-sm text-muted-foreground">
          {data.page.totalItems} playlist(s)
        </span>
      </div>
    {/if}

    <div
      class="grid grid-cols-2 gap-3 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-5 xl:grid-cols-6 2xl:grid-cols-7"
    >
      {#each data.playlists as playlist}
        <div
          class="group relative flex flex-col overflow-hidden rounded-lg border bg-card transition-shadow hover:shadow-md"
        >
          <a href="/playlists/{playlist.id}">
            <Image
              class="aspect-square w-full rounded-none border-0"
              src={playlist.coverArt.medium}
              alt={playlist.name}
            />
          </a>

          <div class="flex flex-col gap-0.5 p-2">
            <a
              href="/playlists/{playlist.id}"
              class="truncate text-sm font-medium hover:underline"
              title={playlist.name}
            >
              {playlist.name}
            </a>
            <p class="truncate text-xs text-muted-foreground">
              {playlist.trackCount} track{playlist.trackCount !== 1 ? "s" : ""}
            </p>
          </div>
        </div>
      {/each}
    </div>

    {#if data.page}
      <Pagination page={data.page} />
    {/if}
  {/if}
</div>
