<script lang="ts">
  import Image from "$lib/components/Image.svelte";
  import TrackListItem from "$lib/components/track-list/TrackListItem.svelte";
  import { getMusicManager } from "$lib/music-manager.svelte";
  import { isRoleAdmin } from "$lib/utils.js";
  import {
    Breadcrumb,
    Button,
    buttonVariants,
    DropdownMenu,
    Separator,
  } from "@nanoteck137/nano-ui";
  import {
    ChevronRight,
    EllipsisVertical,
    ListPlus,
    Pencil,
    Play,
    Shuffle,
  } from "lucide-svelte";

  const { data } = $props();
  const musicManager = getMusicManager();

  // TODO: Replace with real data from API once backend supports it
  let showFullBio = $state(false);
  let mockBio = $state(
    "No biography available. This section will display artist information when the backend provides it.",
  );
  let mockRelatedArtists = $state<{ id: string; name: string; coverArt: string }[]>([]);
</script>

<div class="py-2">
  <Breadcrumb.Root>
    <Breadcrumb.List>
      <Breadcrumb.Item>
        <Breadcrumb.Link href="/artists">Artists</Breadcrumb.Link>
      </Breadcrumb.Item>
      <Breadcrumb.Separator />
      <Breadcrumb.Item>
        <Breadcrumb.Page>{data.artist.name}</Breadcrumb.Page>
      </Breadcrumb.Item>
    </Breadcrumb.List>
  </Breadcrumb.Root>
</div>

<div class="flex flex-col gap-6 rounded-lg border bg-gradient-to-b from-zinc-900 to-background p-4 sm:p-6 md:flex-row md:items-end md:gap-8">
  <Image
    class="w-40 min-w-40 self-center shadow-lg md:w-52 md:min-w-52"
    src={data.artist.coverArt.large}
    alt={data.artist.name}
  />

  <div class="flex min-w-0 flex-col gap-2">
    <p class="text-xs font-semibold uppercase tracking-wider text-muted-foreground">
      Artist
    </p>

    <h1 class="line-clamp-2 text-2xl font-bold md:text-4xl">
      {data.artist.name}
    </h1>

    {#if data.artist.tags.length > 0}
      <div class="flex flex-wrap gap-1">
        {#each data.artist.tags as tag}
          <span class="rounded-full bg-secondary px-2.5 py-0.5 text-xs text-secondary-foreground">{tag}</span>
        {/each}
      </div>
    {/if}

    <div class="flex gap-2 pt-2">
      <Button
        onclick={async () => {
          await musicManager.queueRequest(
            { type: "addArtist", artistId: data.artist.id },
            {},
          );
        }}
      >
        <Play />
        Play
      </Button>

      <Button
        variant="outline"
        onclick={async () => {
          await musicManager.queueRequest(
            { type: "addArtist", artistId: data.artist.id },
            { shuffle: true },
          );
        }}
      >
        <Shuffle />
        Shuffle
      </Button>

      <DropdownMenu.Root>
        <DropdownMenu.Trigger class={buttonVariants({ variant: "outline", size: "icon" })}>
          <EllipsisVertical />
        </DropdownMenu.Trigger>
        <DropdownMenu.Content align="start">
          <DropdownMenu.Group>
            <DropdownMenu.Item
              onSelect={async () => {
                await musicManager.queueRequest(
                  { type: "addArtist", artistId: data.artist.id },
                  { append: "back" },
                );
              }}
            >
              <ListPlus />
              Append to Queue
            </DropdownMenu.Item>
            {#if isRoleAdmin(data.user?.role || "")}
              <DropdownMenu.Link href="/artists/{data.artist.id}/edit">
                <Pencil />
                Edit Artist
              </DropdownMenu.Link>
            {/if}
          </DropdownMenu.Group>
        </DropdownMenu.Content>
      </DropdownMenu.Root>
    </div>
  </div>
</div>

<div class="flex flex-col gap-6">
  <div class="rounded-lg border bg-card">
    <div class="flex items-center justify-between border-b px-4 py-3">
      <p class="font-semibold">Tracks</p>
      <Button href="/artists/{data.artist.id}/tracks" variant="outline" size="sm">
        Show All
        <ChevronRight class="h-4 w-4" />
      </Button>
    </div>
    {#each data.tracks as track}
      <TrackListItem {track} />
      {#if track !== data.tracks[data.tracks.length - 1]}
        <Separator />
      {/if}
    {/each}
  </div>

  <div class="rounded-lg border bg-card">
    <div class="flex items-center justify-between border-b px-4 py-3">
      <p class="font-semibold">Albums</p>
      <Button href="/artists/{data.artist.id}/albums" variant="outline" size="sm">
        Show All
        <ChevronRight class="h-4 w-4" />
      </Button>
    </div>
    <div class="grid grid-cols-2 gap-3 p-4 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-5">
      {#each data.albums as album}
        <a
          href="/albums/{album.id}"
          class="group flex flex-col overflow-hidden rounded-lg border bg-card transition-shadow hover:shadow-md"
        >
          <Image
            class="aspect-square w-full rounded-none border-0"
            src={album.coverArt.medium}
            alt={album.name}
          />
          <div class="flex flex-col gap-0.5 p-2">
            <p class="truncate text-sm font-medium group-hover:underline" title={album.name}>
              {album.name}
            </p>
            <p class="truncate text-xs text-muted-foreground" title={album.artists.map((a) => a.name).join(", ")}>
              {album.artists.map((a) => a.name).join(", ")}
            </p>
          </div>
        </a>
      {/each}
    </div>
  </div>

  <div class="rounded-lg border bg-card">
    <div class="border-b px-4 py-3">
      <p class="font-semibold">About</p>
    </div>
    <div class="p-4">
      <p class="text-sm text-muted-foreground">
        {#if mockBio.length > 200 && !showFullBio}
          {mockBio.slice(0, 200)}&hellip;
        {:else}
          {mockBio}
        {/if}
      </p>
      {#if mockBio.length > 200}
        <Button variant="link" class="h-auto p-0 text-xs" onclick={() => (showFullBio = !showFullBio)}>
          {showFullBio ? "Show less" : "Read more"}
        </Button>
      {/if}
    </div>
  </div>

  {#if mockRelatedArtists.length > 0}
    <div class="rounded-lg border bg-card">
      <div class="border-b px-4 py-3">
        <p class="font-semibold">Related Artists</p>
      </div>
      <div class="flex gap-3 overflow-x-auto p-4 pb-2">
        {#each mockRelatedArtists as related}
          <a href="/artists/{related.id}" class="flex w-32 shrink-0 flex-col items-center gap-1">
            <Image class="w-32" src={related.coverArt} alt={related.name} />
            <p class="line-clamp-2 text-center text-xs font-medium">{related.name}</p>
          </a>
        {/each}
      </div>
    </div>
  {/if}
</div>
