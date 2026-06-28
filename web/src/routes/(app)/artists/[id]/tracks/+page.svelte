<script lang="ts">
  import { goto } from "$app/navigation";
  import { page } from "$app/state";
  import { Breadcrumb, Button, Select, Separator } from "@nanoteck137/nano-ui";
  import { Play, Shuffle } from "lucide-svelte";
  import TrackList from "$lib/components/track-list/TrackList.svelte";
  import { getMusicManager } from "$lib/music-manager.svelte";
  import Pagination from "$lib/components/Pagination.svelte";
  import Spacer from "$lib/components/Spacer.svelte";
  import { defineEnumTypes } from "$lib/utils";

  let { data } = $props();
  const musicManager = getMusicManager();

  const { sortTypes, defaultSort } = defineEnumTypes(
    [
      { label: "Name (A-Z)", value: "name-a-z" },
      { label: "Name (Z-A)", value: "name-z-a" },
      { label: "Artist", value: "artist" },
      { label: "Album", value: "album" },
      { label: "Duration", value: "duration" },
      { label: "Year", value: "year" },
      { label: "Added (New–Old)", value: "created-new" },
      { label: "Added (Old-New)", value: "created-old" },
    ] as const,
    "name-a-z",
  );

  type SortType = (typeof sortTypes)[number]["value"];

  let sort = $state(
    (page.url.searchParams.get("sort") as SortType) ?? defaultSort,
  );

  function updateSort(value: string) {
    sort = value as SortType;

    const query = page.url.searchParams;
    query.delete("sort");

    if (sort !== defaultSort) {
      query.set("sort", sort);
    }

    goto("?" + query.toString(), { invalidateAll: true });
  }
</script>

<div class="flex flex-col gap-4">
  <Breadcrumb.Root>
    <Breadcrumb.List>
      <Breadcrumb.Item>
        <Breadcrumb.Link href="/artists">Artists</Breadcrumb.Link>
      </Breadcrumb.Item>
      <Breadcrumb.Separator />
      <Breadcrumb.Item>
        <Breadcrumb.Link href="/artists/{data.artist.id}">
          {data.artist.name}
        </Breadcrumb.Link>
      </Breadcrumb.Item>
      <Breadcrumb.Separator />
      <Breadcrumb.Item>
        <Breadcrumb.Page>Tracks</Breadcrumb.Page>
      </Breadcrumb.Item>
    </Breadcrumb.List>
  </Breadcrumb.Root>

  <div
    class="flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between"
  >
    <div class="flex items-baseline gap-2">
      <h1 class="text-xl font-bold">Tracks</h1>
      {#if data.page}
        <span class="text-sm text-muted-foreground"
          >{data.page.totalItems}</span
        >
      {/if}
    </div>

    <div class="flex items-center gap-2">
      <Button
        variant="outline"
        size="sm"
        onclick={async () => {
          await musicManager.queueRequest(
            { type: "addArtist", artistId: data.artist.id },
            { shuffle: true },
          );
        }}
      >
        <Shuffle size={14} />
        Shuffle
      </Button>
      <Button
        size="sm"
        onclick={async () => {
          await musicManager.queueRequest(
            { type: "addArtist", artistId: data.artist.id },
            {},
          );
        }}
      >
        <Play size={14} />
        Play All
      </Button>
    </div>

    <Select.Root
      type="single"
      allowDeselect={false}
      value={sort}
      onValueChange={updateSort}
    >
      <Select.Trigger class="h-9 w-full sm:w-40">
        {sortTypes.find((i) => i.value === sort)?.label ?? "Sort"}
      </Select.Trigger>
      <Select.Content>
        {#each sortTypes as ty (ty.value)}
          <Select.Item value={ty.value} label={ty.label} />
        {/each}
      </Select.Content>
    </Select.Root>
  </div>

  <TrackList
    totalTracks={data.tracks.length}
    tracks={data.tracks}
    onPlay={async (trackId) => {
      // TODO: queue artist tracks and start from trackId
    }}
  />

  <Spacer size="lg" />
  <Separator />
  <Spacer size="lg" />

  <Pagination page={data.page} />
</div>
