<script lang="ts">
  import { goto } from "$app/navigation";
  import { page } from "$app/state";
  import { Button, Input } from "@nanoteck137/nano-ui";
  import { Search, X } from "lucide-svelte";
  import { onMount } from "svelte";
  import { cn } from "$lib/utils";
  import TrackList from "$lib/components/track-list/TrackList.svelte";
  import Image from "$lib/components/Image.svelte";

  const { data } = $props();

  async function search(query: string) {
    await goto(`?query=${query}`, {
      invalidateAll: true,
      keepFocus: true,
      replaceState: true,
    });
  }

  function clearSearch() {
    value = "";
    search("");
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
      search(current);
    }, 500);
  }

  function formatError(err: { type: string; code: number; message: string }) {
    return err.message;
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
  <title>Search - Tunebook</title>
</svelte:head>

<div class="flex flex-col gap-6">
  <div>
    <h1 class="mb-4 text-xl font-bold">Search</h1>

    <form
      action=""
      method="get"
      onsubmit={(e) => {
        e.preventDefault();
        clearTimeout(timer);
        search(value);
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
              placeholder="Search artists, albums, tracks, playlists, users..."
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

  {#if data.artistError || data.albumError || data.trackError}
    <div class="flex flex-col gap-1 text-sm text-red-400">
      {#if data.artistError}
        <p>Artists: {formatError(data.artistError)}</p>
      {/if}
      {#if data.albumError}
        <p>Albums: {formatError(data.albumError)}</p>
      {/if}
      {#if data.trackError}
        <p>Tracks: {formatError(data.trackError)}</p>
      {/if}
    </div>
  {/if}

  {#if data.query && data.artists.length === 0 && data.albums.length === 0 && data.tracks.length === 0 && data.playlists.length === 0 && data.users.length === 0}
    <p class="py-12 text-center text-sm text-muted-foreground">
      No results found for "{data.query}".
    </p>
  {/if}

  {#if data.tracks.length > 0}
    <div>
      <div class="mb-3 flex items-baseline justify-between">
        <h2 class="text-base font-semibold">Tracks</h2>
        <div class="flex items-center gap-2">
          <span class="text-xs text-muted-foreground"
            >{data.tracks.length}</span
          >
          <a
            href="/search/tracks?query={data.query}"
            class="text-xs font-medium text-primary hover:underline"
          >
            View all
          </a>
        </div>
      </div>

      <TrackList
        totalTracks={data.tracks.length}
        tracks={data.tracks}
        onPlay={() => {}}
      />
    </div>
  {/if}

  {#if data.artists.length > 0}
    <div>
      <div class="mb-3 flex items-baseline justify-between">
        <h2 class="text-base font-semibold">Artists</h2>
        <div class="flex items-center gap-2">
          <span class="text-xs text-muted-foreground"
            >{data.artists.length}</span
          >
          <a
            href="/search/artists?query={data.query}"
            class="text-xs font-medium text-primary hover:underline"
          >
            View all
          </a>
        </div>
      </div>

      <div
        class="grid grid-cols-2 gap-3 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-5 xl:grid-cols-6"
      >
        {#each data.artists.slice(0, 6) as artist}
          <a
            href="/artists/{artist.id}"
            class="flex flex-col overflow-hidden rounded-lg border bg-card transition-shadow hover:shadow-md"
          >
            <img
              src={artist.coverArt.medium}
              alt={artist.name}
              class="aspect-square w-full object-cover"
            />
            <div class="p-2">
              <p class="truncate text-sm font-medium" title={artist.name}>
                {artist.name}
              </p>
            </div>
          </a>
        {/each}
      </div>
    </div>
  {/if}

  {#if data.albums.length > 0}
    <div>
      <div class="mb-3 flex items-baseline justify-between">
        <h2 class="text-base font-semibold">Albums</h2>
        <div class="flex items-center gap-2">
          <span class="text-xs text-muted-foreground"
            >{data.albums.length}</span
          >
          <a
            href="/search/albums?query={data.query}"
            class="text-xs font-medium text-primary hover:underline"
          >
            View all
          </a>
        </div>
      </div>

      <div
        class="grid grid-cols-2 gap-3 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-5 xl:grid-cols-6"
      >
        {#each data.albums.slice(0, 6) as album}
          <a
            href="/albums/{album.id}"
            class="flex flex-col overflow-hidden rounded-lg border bg-card transition-shadow hover:shadow-md"
          >
            <img
              src={album.coverArt.medium}
              alt={album.name}
              class="aspect-square w-full object-cover"
            />
            <div class="p-2">
              <p class="truncate text-sm font-medium" title={album.name}>
                {album.name}
              </p>
              <p
                class="truncate text-xs text-muted-foreground"
                title={album.artists.map((a) => a.name).join(", ")}
              >
                {album.artists.map((a) => a.name).join(", ")}
              </p>
            </div>
          </a>
        {/each}
      </div>
    </div>
  {/if}

  {#if data.playlists.length > 0}
    <div>
      <div class="mb-3 flex items-baseline justify-between">
        <h2 class="text-base font-semibold">Playlists</h2>
        <div class="flex items-center gap-2">
          <span class="text-xs text-muted-foreground"
            >{data.playlists.length}</span
          >
          <a
            href="/search/playlists?query={data.query}"
            class="text-xs font-medium text-primary hover:underline"
          >
            View all
          </a>
        </div>
      </div>

      <div
        class="grid grid-cols-2 gap-3 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-5 xl:grid-cols-6"
      >
        {#each data.playlists.slice(0, 6) as playlist}
          <a
            href="/playlists/{playlist.id}"
            class="flex flex-col overflow-hidden rounded-lg border bg-card transition-shadow hover:shadow-md"
          >
            <Image
              class="aspect-square w-full rounded-none border-0"
              src={playlist.coverArt.medium}
              alt={playlist.name}
            />
            <div class="p-2">
              <p class="truncate text-sm font-medium" title={playlist.name}>
                {playlist.name}
              </p>
              <p class="truncate text-xs text-muted-foreground">
                {playlist.ownerDisplayName}
              </p>
            </div>
          </a>
        {/each}
      </div>
    </div>
  {/if}

  {#if data.users.length > 0}
    <div>
      <div class="mb-3 flex items-baseline justify-between">
        <h2 class="text-base font-semibold">Users</h2>
        <div class="flex items-center gap-2">
          <span class="text-xs text-muted-foreground">{data.users.length}</span
          >
          <a
            href="/search/users?query={data.query}"
            class="text-xs font-medium text-primary hover:underline"
          >
            View all
          </a>
        </div>
      </div>

      <div
        class="grid grid-cols-2 gap-3 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-5 xl:grid-cols-6"
      >
        {#each data.users.slice(0, 6) as user}
          <a
            href="/users/{user.id}"
            class="flex flex-col overflow-hidden rounded-lg border bg-card transition-shadow hover:shadow-md"
          >
            <img
              src={user.picture.medium}
              alt={user.displayName}
              class="aspect-square w-full object-cover"
            />
            <div class="p-2">
              <p class="truncate text-sm font-medium" title={user.displayName}>
                {user.displayName}
              </p>
              <p class="truncate text-xs text-muted-foreground">{user.role}</p>
            </div>
          </a>
        {/each}
      </div>
    </div>
  {/if}
</div>
