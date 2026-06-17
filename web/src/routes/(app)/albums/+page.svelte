<script lang="ts">
  import { Search } from "lucide-svelte";
  import AlbumTile from "$lib/components/tiles/AlbumTile.svelte";
  import Spacer from "$lib/components/Spacer.svelte";
  import Pagination from "$lib/components/Pagination.svelte";
  import { Separator, Button, Select } from "@nanoteck137/nano-ui";
  import { goto } from "$app/navigation";
  import { sortTypes, type SortType } from "./types";
  import { page } from "$app/state";

  let { data } = $props();

  let sort = $state(data.filter.sort);
  function updateSort(value: string) {
    sort = value as SortType;

    const query = page.url.searchParams;
    query.delete("sort");

    query.set("sort", sort);

    goto("?" + query.toString(), { invalidateAll: true });
  }
</script>

<div class="flex flex-col gap-4">
  <div class="flex items-center justify-between">
    <div class="flex items-center gap-2">
      <p class="text-bold text-xl">Albums</p>
      {#if data.page}
        <span class="text-sm text-muted-foreground">
          ({data.page.totalItems} albums)
        </span>
      {/if}
    </div>
    <div class="flex items-center gap-2">
      <Select.Root
        type="single"
        allowDeselect={false}
        value={sort}
        onValueChange={updateSort}
      >
        <Select.Trigger class="w-48">
          {sortTypes.find((i) => i.value === sort)?.label ?? "Sort"}
        </Select.Trigger>
        <Select.Content>
          {#each sortTypes as ty (ty.value)}
            <Select.Item value={ty.value} label={ty.label} />
          {/each}
        </Select.Content>
      </Select.Root>

      <Button variant="outline" href="/search/albums">
        <Search />
        Search
      </Button>
    </div>
  </div>
</div>

<Spacer size="lg" />

<!-- <div
    class="grid grid-cols-2 gap-2 md:grid-cols-4 lg:grid-cols-5 xl:grid-cols-6 2xl:grid-cols-7"
  > -->
<div class="flex flex-shrink flex-wrap justify-center gap-4">
  {#each data.albums as album}
    <AlbumTile
      id={album.id}
      cover={album.coverArt.medium}
      name={album.name}
      artists={album.artists}
    />
  {/each}
</div>

<Spacer size="lg" />

<Separator />

<Spacer size="lg" />

<Pagination page={data.page} />
