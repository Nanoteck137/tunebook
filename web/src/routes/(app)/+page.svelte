<script lang="ts">
  import Spinner from "$lib/components/Spinner.svelte";
  import { Button } from "@nanoteck137/nano-ui";

  let { data } = $props();
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
    <!-- <GradientButton href="/login">Login</GradientButton> -->
    <Button size="lg" href="/login">Login</Button>
  </div>
{:else}
  <div class="flex flex-col gap-8 p-4">
    <div class="flex items-center gap-4">
      <a
        class="bg-gradient-to-tr from-logo-1 via-logo-2 to-logo-3 bg-clip-text text-2xl font-medium text-transparent"
        href="/"
      >
        Tunebook
      </a>
      <span class="text-xl font-bold">|</span>
      <h1 class="text-2xl font-bold">
        Welcome{`, ${data.user.displayName}`}!
      </h1>
    </div>

    <section>
      <a
        class="mb-4 text-xl font-semibold hover:cursor-pointer hover:underline"
        href="/playlists"
      >
        Your Playlists
      </a>

      <div class="flex gap-2 overflow-x-auto pb-2">
        {#await data.playlists}
          <Spinner />
        {:then playlists}
          {#each playlists as playlist}
            <div class="flex shrink-0 flex-col items-center">
              <a
                href="/playlists/{playlist.id}"
                class="group w-40 cursor-pointer"
              >
                <img
                  class="aspect-square w-40 rounded-lg object-cover"
                  src={playlist.coverArt.medium}
                  alt=""
                />
                <div
                  class="mt-2 w-40 truncate text-sm font-medium group-hover:underline"
                >
                  {playlist.name}
                </div>
                <div class="mt-1 text-xs text-muted-foreground">
                  {playlist.trackCount} tracks
                </div>
              </a>
            </div>
          {:else}
            <p>No Playlists</p>
          {/each}
        {/await}
      </div>
    </section>

    <section>
      <a
        class="mb-4 text-xl font-semibold hover:cursor-pointer hover:underline"
        href="/users/{data.user.id}/favorites"
      >
        Favorites
      </a>
      <div class="flex gap-2 overflow-x-auto pb-2">
        {#await data.favorites}
          <Spinner />
        {:then favorites}
          {#each favorites as track}
            <div class="group flex shrink-0 flex-col items-center">
              <a href="/tracks" class="w-40 cursor-pointer">
                <img
                  class="aspect-square w-40 rounded-lg object-cover"
                  src={track.coverArt.medium}
                  alt=""
                />
                <div
                  class="mt-2 w-40 truncate text-sm font-medium group-hover:underline"
                >
                  {track.name}
                </div>
              </a>

              <a
                href="/artists/{track.artists[0].id}"
                class="mt-1 w-40 truncate text-xs text-muted-foreground hover:underline"
              >
                {track.artists[0].name}
              </a>
            </div>
          {/each}
        {/await}
      </div>
    </section>

    <section>
      <h2 class="mb-4 text-xl font-semibold">Recently Added Albums</h2>
      <div class="flex gap-2 overflow-x-auto pb-2">
        {#await data.recentAlbums}
          <Spinner />
        {:then albums}
          {#each albums as album}
            <div class="group flex shrink-0 flex-col items-center">
              <a href="/albums/{album.id}" class="w-40 cursor-pointer">
                <img
                  class="aspect-square w-40 rounded-lg object-cover"
                  src={album.coverArt.medium}
                  alt=""
                />
                <div
                  class="mt-2 w-40 truncate text-sm font-medium group-hover:underline"
                >
                  {album.name}
                </div>
              </a>
              <a
                href="/artists/{album.artists[0].id}"
                class="mt-1 w-40 truncate text-xs text-muted-foreground hover:underline"
              >
                {album.artists[0].name}
              </a>
            </div>
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
