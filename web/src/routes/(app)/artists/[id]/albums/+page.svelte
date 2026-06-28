<script lang="ts">
  import { goto } from "$app/navigation";
  import { page } from "$app/state";
  import { Breadcrumb, Select, Separator } from "@nanoteck137/nano-ui";
  import Image from "$lib/components/Image.svelte";
  import Pagination from "$lib/components/Pagination.svelte";
  import Spacer from "$lib/components/Spacer.svelte";
  import { defineEnumTypes } from "$lib/utils";

  let { data } = $props();

  const { sortTypes, defaultSort } = defineEnumTypes(
    [
      { label: "Name (A-Z)", value: "name-a-z" },
      { label: "Name (Z-A)", value: "name-z-a" },
      { label: "Year (New–Old)", value: "year-new" },
      { label: "Year (Old-New)", value: "year-old" },
      { label: "Created (New–Old)", value: "created-new" },
      { label: "Created (Old-New)", value: "created-old" },
      { label: "Updated (New–Old)", value: "updated-new" },
      { label: "Updated (Old-New)", value: "updated-old" },
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
        <Breadcrumb.Page>Albums</Breadcrumb.Page>
      </Breadcrumb.Item>
    </Breadcrumb.List>
  </Breadcrumb.Root>

  <div
    class="flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between"
  >
    <div class="flex items-baseline gap-2">
      <h1 class="text-xl font-bold">Albums</h1>
      {#if data.page}
        <span class="text-sm text-muted-foreground"
          >{data.page.totalItems}</span
        >
      {/if}
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

  <div
    class="grid grid-cols-2 gap-3 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-5 xl:grid-cols-6 2xl:grid-cols-7"
  >
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
          <p
            class="truncate text-sm font-medium group-hover:underline"
            title={album.name}
          >
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

  <Spacer size="lg" />
  <Separator />
  <Spacer size="lg" />

  <Pagination page={data.page} />
</div>
