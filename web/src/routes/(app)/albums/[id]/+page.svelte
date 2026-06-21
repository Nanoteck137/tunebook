<script lang="ts">
  import {
    Breadcrumb,
    Button,
    buttonVariants,
    DropdownMenu,
  } from "@nanoteck137/nano-ui";
  import { EllipsisVertical, ListPlus, Play, Shuffle } from "lucide-svelte";
  import TrackList from "$lib/components/track-list/TrackList.svelte";
  import { getMusicManager } from "$lib/music-manager.svelte.js";
  import Image from "$lib/components/Image.svelte";
  import ArtistList from "$lib/components/ArtistList.svelte";

  let { data } = $props();
  const musicManager = getMusicManager();
</script>

<div class="py-2">
  <Breadcrumb.Root>
    <Breadcrumb.List>
      <Breadcrumb.Item>
        <Breadcrumb.Link href="/albums">Albums</Breadcrumb.Link>
      </Breadcrumb.Item>
      <Breadcrumb.Separator />
      <Breadcrumb.Item>
        <Breadcrumb.Page>{data.album.name}</Breadcrumb.Page>
      </Breadcrumb.Item>
    </Breadcrumb.List>
  </Breadcrumb.Root>
</div>

<div
  class="flex flex-col gap-6 rounded-lg border bg-gradient-to-b from-zinc-900 to-background p-4 sm:p-6 md:flex-row md:items-end md:gap-8"
>
  <Image
    class="w-40 min-w-40 self-center shadow-lg md:w-52 md:min-w-52"
    src={data.album.coverArt.large}
    alt={data.album.name}
  />

  <div class="flex min-w-0 flex-col gap-2">
    <p
      class="text-xs font-semibold uppercase tracking-wider text-muted-foreground"
    >
      {data.album.albumType}
    </p>

    <h1 class="line-clamp-2 text-2xl font-bold md:text-4xl">
      {data.album.name}
    </h1>

    <div
      class="flex flex-wrap items-center gap-x-1 text-sm text-muted-foreground"
    >
      <ArtistList artists={data.album.artists} class="text-sm" />

      {#if data.album.year}
        <span>&middot; {data.album.year}</span>
      {/if}
    </div>

    {#if data.album.tags.length > 0}
      <div class="flex flex-wrap gap-1">
        {#each data.album.tags as tag}
          <span
            class="rounded-full bg-secondary px-2.5 py-0.5 text-xs text-secondary-foreground"
            >{tag}</span
          >
        {/each}
      </div>
    {/if}

    <div class="flex gap-2 pt-2">
      <Button
        onclick={async () => {
          await musicManager.addAlbumTracks({
            albumId: data.album.id,
            clear: true,
          });
        }}
      >
        <Play />
        Play
      </Button>

      <Button
        variant="outline"
        onclick={async () => {
          await musicManager.addAlbumTracks({
            albumId: data.album.id,
            clear: true,
          });
          // TODO: shuffle after adding
        }}
      >
        <Shuffle />
        Shuffle
      </Button>

      <DropdownMenu.Root>
        <DropdownMenu.Trigger
          class={buttonVariants({ variant: "outline", size: "icon" })}
        >
          <EllipsisVertical />
        </DropdownMenu.Trigger>
        <DropdownMenu.Content align="start">
          <DropdownMenu.Group>
            <DropdownMenu.Item
              onSelect={async () => {
                // await musicManager.queueRequest(
                //   { type: "addAlbum", albumId: data.album.id },
                //   { append: "back" },
                // );
              }}
            >
              <ListPlus />
              Append to Queue
            </DropdownMenu.Item>
          </DropdownMenu.Group>
        </DropdownMenu.Content>
      </DropdownMenu.Root>
    </div>
  </div>
</div>

<div class="h-4"></div>

<TrackList
  isAlbumShowcase={true}
  totalTracks={data.tracks.length}
  tracks={data.tracks}
  onPlay={async (trackId) => {
    await musicManager.addAlbumTracks({
      albumId: data.album.id,
      trackId,
    });
  }}
/>
