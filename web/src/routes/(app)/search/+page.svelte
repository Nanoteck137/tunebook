<script lang="ts">
  import { goto, invalidateAll, onNavigate } from "$app/navigation";
  import AlbumListItem from "$lib/components/AlbumListItem.svelte";
  import ArtistListItem from "$lib/components/ArtistListItem.svelte";
  import Image from "$lib/components/Image.svelte";
  import TrackList from "$lib/components/track-list/TrackList.svelte";
  import { Button, Input, Label } from "@nanoteck137/nano-ui";
  import { onMount } from "svelte";

  const { data } = $props();

  async function search(query: string) {
    await goto(`?query=${query}`, {
      invalidateAll: true,
      keepFocus: true,
      replaceState: true,
    });
  }

  let initialValue = $state();
  let value = "";

  onMount(() => {
    initialValue = data.query;
  });

  // eslint-disable-next-line no-undef
  let timer: NodeJS.Timeout;
  function onInput(e: Event) {
    const target = e.target as HTMLInputElement;
    const current = target.value;
    value = current;

    clearTimeout(timer);
    timer = setTimeout(async () => {
      search(current);
    }, 500);
  }

  // NOTE(patrik): Fix for clicking the search button
  onNavigate((e) => {
    if (e.type === "link" && e.from?.url.pathname === "/search") {
      invalidateAll();
    }
  });

  function formatError(err: { type: string; code: number; message: string }) {
    // TODO(patrik): Better error
    return err.message;
  }
</script>

<form
  action=""
  method="get"
  onsubmit={(e) => {
    e.preventDefault();
    clearTimeout(timer);
    search(value);
  }}
>
  <div class="flex flex-col gap-4">
    <div class="flex flex-col gap-2">
      <Label for="query">Search Query</Label>
      <Input
        id="query"
        name="query"
        autocomplete="off"
        value={initialValue}
        oninput={onInput}
      />
    </div>
    <Button type="submit">Search</Button>
  </div>
</form>

<div class="h-4"></div>

{#if data.artistError}
  <p class="text-red-400">
    Artist Query Error: {formatError(data.artistError)}
  </p>
{/if}

{#if data.albumError}
  <p class="text-red-400">
    Album Query Error: {formatError(data.albumError)}
  </p>
{/if}

{#if data.trackError}
  <p class="text-red-400">
    Track Query Error: {formatError(data.trackError)}
  </p>
{/if}

{#if data.artists.length > 0}
  <div class="flex items-center justify-between">
    <p class="text-bold">Artists</p>
    <p class="text-xs">{data.artists.length} artist(s)</p>
  </div>

  {#each data.artists as artist}
    <ArtistListItem {artist} link />
  {/each}
{/if}

{#if data.albums.length > 0}
  <div class="flex items-center justify-between">
    <p class="text-bold">Albums</p>
    <p class="text-xs">{data.albums.length} album(s)</p>
  </div>

  {#each data.albums as album}
    <AlbumListItem {album} link />
  {/each}
{/if}

{#if data.tracks.length > 0}
  <TrackList
    totalTracks={data.tracks.length}
    tracks={data.tracks}
    userPlaylists={data.userPlaylists}
    quickPlaylist={data.user?.quickPlaylist}
    onPlay={() => {}}
  />
  <!-- onTrackPlay={() => {}} -->
{/if}

{#if data.playlists.length > 0}
  <div class="flex items-center justify-between">
    <p class="text-bold">Playlists</p>
    <p class="text-xs">{data.playlists.length} playlist(s)</p>
  </div>

  {#each data.playlists as playlist}
    <div class="flex items-center gap-2 border-b py-2 pr-2">
      <div class="group relative">
        <Image
          class="w-14 min-w-14"
          src={playlist.coverArt.small}
          alt="cover"
        />
      </div>
      <div class="flex flex-grow flex-col">
        <div class="flex items-center gap-1">
          <a
            class="line-clamp-1 w-fit text-sm font-medium hover:underline"
            title={playlist.name}
            href="/playlists/{playlist.id}"
          >
            {playlist.name}
          </a>

          <!-- <p>•</p> -->

          <!-- <ArtistList artists={album.artists} /> -->
          <!-- <a
            class="line-clamp-1 text-xs font-light hover:underline"
            title={playlist.ownerDisplayName}
            href={`/artists/${playlist.ownerId}`}
          >
            {playlist.ownerDisplayName}
          </a> -->
        </div>

        <div class="flex items-center gap-2">
          <Image
            class="w-6 min-w-6"
            src={playlist.ownerPicture.small}
            alt="cover"
          />

          <a
            class="line-clamp-1 text-xs hover:underline"
            title={playlist.ownerDisplayName}
            href={`/artists/${playlist.ownerId}`}
          >
            {playlist.ownerDisplayName}
          </a>
        </div>
      </div>
      <div class="flex items-center">
        <!-- {@render children?.()} -->
      </div>
    </div>
  {/each}
{/if}

{#if data.users.length > 0}
  <div class="flex items-center justify-between">
    <p class="text-bold">Users</p>
    <p class="text-xs">{data.users.length} user(s)</p>
  </div>

  {#each data.users as user}
    <div class="flex items-center gap-2 border-b py-2 pr-2">
      <div class="group relative">
        <Image class="w-14 min-w-14" src={user.picture.small} alt="cover" />
      </div>
      <div class="flex flex-grow flex-col">
        <div class="flex items-center gap-1">
          <a
            class="line-clamp-1 w-fit text-sm font-medium hover:underline"
            title={user.displayName}
            href="/users/{user.id}"
          >
            {user.displayName}
          </a>

          <!-- <p>•</p> -->

          <!-- <ArtistList artists={album.artists} /> -->
          <!-- <a
            class="line-clamp-1 text-xs font-light hover:underline"
            title={playlist.ownerDisplayName}
            href={`/artists/${playlist.ownerId}`}
          >
            {playlist.ownerDisplayName}
          </a> -->
        </div>

        <!-- <div class="flex items-center gap-2">
          <Image
            class="w-6 min-w-6"
            src={playlist.ownerPicture.small}
            alt="cover"
          />

          <a
            class="line-clamp-1 text-xs hover:underline"
            title={playlist.ownerDisplayName}
            href={`/artists/${playlist.ownerId}`}
          >
            {playlist.ownerDisplayName}
          </a>
        </div> -->
      </div>
      <div class="flex items-center">
        <!-- {@render children?.()} -->
      </div>
    </div>
  {/each}
{/if}
