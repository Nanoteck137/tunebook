<script lang="ts">
  import AlbumSkeletonTile from "$lib/components/tiles/AlbumSkeletonTile.svelte";
  import AlbumTile from "$lib/components/tiles/AlbumTile.svelte";
  import PlaylistSkeletonTile from "$lib/components/tiles/PlaylistSkeletonTile.svelte";
  import PlaylistTile from "$lib/components/tiles/PlaylistTile.svelte";
  import TrackSkeletonTile from "$lib/components/tiles/TrackSkeletonTile.svelte";
  import TrackTile from "$lib/components/tiles/TrackTile.svelte";
  import { Button } from "@nanoteck137/nano-ui";
  import { ChevronRight } from "lucide-svelte";
  import { getMusicManager } from "$lib/music-manager.svelte";
  let { data } = $props();

  const musicManager = getMusicManager();
</script>

{#if !data.user}
  <div class="mt-16 flex flex-col items-center justify-center gap-8 p-4">
    <div
      class="flex h-32 w-32 items-center justify-center rounded-full bg-gradient-to-tr from-logo-1 via-logo-2 to-logo-3"
    >
      <svg
        xmlns="http://www.w3.org/2000/svg"
        class="h-16 w-16 text-black"
        viewBox="0 0 24 24"
        fill="none"
        stroke="currentColor"
        stroke-width="2"
        stroke-linecap="round"
        stroke-linejoin="round"
      >
        <circle cx="12" cy="12" r="10" /><circle cx="12" cy="12" r="3" />
      </svg>
    </div>
    <h1
      class="bg-gradient-to-tr from-logo-1 via-logo-2 to-logo-3 bg-clip-text text-5xl font-bold text-transparent"
    >
      Tunebook
    </h1>
    <p class="text-muted-foreground">Your personal music streaming server</p>
    <Button size="lg" href="/login">Login</Button>
  </div>
{:else}
  <div class="flex flex-col gap-8">
    <div class="flex flex-col items-center gap-4 border-b pb-4">
      <a
        class="bg-gradient-to-tr from-logo-1 via-logo-2 to-logo-3 bg-clip-text px-8 text-2xl font-medium text-transparent"
        href="/"
      >
        Tunebook
      </a>

      <h1 class="text-center text-2xl font-bold">
        Welcome, {data.user.displayName}!
      </h1>
    </div>

    <section>
      <a
        class="flex items-center gap-1 text-xl font-semibold hover:cursor-pointer hover:underline"
        href="/playlists"
      >
        Your Playlists
        <ChevronRight />
      </a>

      <div class="h-4"></div>

      <div class="flex gap-2 overflow-x-auto pb-4">
        {#await data.playlists}
          {#each Array(5) as _i}
            <PlaylistSkeletonTile />
          {/each}
        {:then playlists}
          {#each playlists as playlist (playlist.id)}
            <PlaylistTile
              id={playlist.id}
              cover={playlist.coverArt.medium}
              name={playlist.name}
              trackCount={playlist.trackCount}
            />
          {:else}
            <p>No Playlists</p>
          {/each}
        {/await}
      </div>
    </section>

    <section>
      <a
        class="flex items-center gap-1 text-xl font-semibold hover:cursor-pointer hover:underline"
        href="/users/{data.user.id}/favorites"
      >
        Favorites
        <ChevronRight />
      </a>

      <div class="h-4"></div>

      <div class="flex gap-2 overflow-x-auto pb-4">
        {#await data.favorites}
          {#each Array(5) as _i}
            <TrackSkeletonTile />
          {/each}
        {:then favorites}
          {#each favorites as track (track.id)}
            <TrackTile
              id={track.id}
              cover={track.coverArt.medium}
              name={track.name}
              artists={track.artists}
            />
          {/each}
        {/await}
      </div>
    </section>

    <section>
      <a
        class="flex items-center gap-1 text-xl font-semibold hover:cursor-pointer hover:underline"
        href="/albums?sort=created-new"
      >
        Recently Added Albums
        <ChevronRight />
      </a>

      <div class="h-4"></div>

      <div class="flex gap-2 overflow-x-auto pb-4">
        {#await data.recentAlbums}
          {#each Array(5) as _i}
            <AlbumSkeletonTile />
          {/each}
        {:then albums}
          {#each albums as album (album.id)}
            <AlbumTile
              id={album.id}
              cover={album.coverArt.medium}
              name={album.name}
              artists={album.artists}
            />
          {/each}
        {/await}
      </div>
    </section>

    <footer class="mt-8 border-t pt-4">
      <div
        class="flex flex-col items-center justify-center gap-2 text-sm text-muted-foreground"
      >
        <a
          class="bg-gradient-to-tr from-logo-1 via-logo-2 to-logo-3 bg-clip-text text-lg font-medium text-transparent"
          href="/"
        >
          Tunebook
        </a>
        <p>Your personal music streaming server</p>
      </div>
    </footer>
  </div>
{/if}
