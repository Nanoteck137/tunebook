<script lang="ts">
  import { goto } from "$app/navigation";
  import { page } from "$app/state";
  import { Button, Input, Separator } from "@nanoteck137/nano-ui";
  import { Search, X } from "lucide-svelte";
  import Pagination from "$lib/components/Pagination.svelte";
  import Spacer from "$lib/components/Spacer.svelte";
  import Image from "$lib/components/Image.svelte";
  import { cn } from "$lib/utils";
  import { onMount } from "svelte";

  let { data } = $props();

  async function doSearch(query: string) {
    await goto(`/search/playlists?query=${query}`, {
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
  <title>Search Playlists - Tunebook</title>
</svelte:head>

<div class="flex flex-col gap-6">
  <div>
    <h1 class="mb-4 text-xl font-bold">Search Playlists</h1>

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
              placeholder="Search playlists..."
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
            <p
              class="flex items-center gap-1 truncate text-xs text-muted-foreground"
            >
              {#if playlist.ownerPicture}
                <img
                  src={playlist.ownerPicture.small}
                  alt=""
                  class="h-4 w-4 rounded-full object-cover"
                />
              {/if}
              {playlist.ownerDisplayName}
            </p>
          </div>
        </div>
      {/each}
    </div>

    {#if data.page}
      <Spacer size="lg" />
      <Separator />
      <Spacer size="lg" />

      <Pagination page={data.page} />
    {/if}
  {/if}
</div>
