<script lang="ts">
  import { Search, X, Plus, EllipsisVertical, Info } from "lucide-svelte";
  import Spacer from "$lib/components/Spacer.svelte";
  import Pagination from "$lib/components/Pagination.svelte";
  import {
    Separator,
    Button,
    Select,
    Input,
    Dialog,
    DropdownMenu,
    buttonVariants,
  } from "@nanoteck137/nano-ui";
  import { cn } from "$lib/utils";
  import { goto } from "$app/navigation";
  import {
    sortTypes,
    decadeTypes,
    defaultSort,
    defaultDecade,
    type SortType,
    type DecadeType,
  } from "./types";
  import { page } from "$app/state";

  let { data } = $props();

  let sort = $state(
    (page.url.searchParams.get("sort") as SortType) ?? defaultSort,
  );
  function updateSort(value: string) {
    sort = value as SortType;

    const query = page.url.searchParams;
    query.delete("sort");
    query.set("sort", sort);

    goto("?" + query.toString(), { invalidateAll: true });
  }

  let searchQuery = $state(page.url.searchParams.get("query") ?? "");
  function updateSearch() {
    const query = page.url.searchParams;
    query.delete("query");

    if (searchQuery) {
      query.set("query", searchQuery);
    }

    goto("?" + query.toString(), { invalidateAll: true });
  }

  let decade = $state(
    (page.url.searchParams.get("decade") as DecadeType) ?? defaultDecade,
  );
  function updateDecade(value: string) {
    decade = value as DecadeType;

    const query = page.url.searchParams;
    query.delete("decade");

    if (decade !== "none") {
      query.set("decade", decade);
    }

    goto("?" + query.toString(), { invalidateAll: true });
  }

  let tagInput = $state("");
  let includeTags = $state(
    page.url.searchParams.get("tags")?.split(",").filter(Boolean) ?? [],
  );
  let excludeTags = $state(
    page.url.searchParams.get("excludeTags")?.split(",").filter(Boolean) ?? [],
  );
  let tagMode = $state<"include" | "exclude">("include");

  function addTag() {
    const tag = tagInput.trim();
    if (!tag) return;

    if (tagMode === "include") {
      if (includeTags.includes(tag)) return;
      includeTags = [...includeTags, tag];
    } else {
      if (excludeTags.includes(tag)) return;
      excludeTags = [...excludeTags, tag];
    }

    tagInput = "";
    applyTagFilters();
  }

  function removeIncludeTag(tag: string) {
    includeTags = includeTags.filter((t) => t !== tag);
    applyTagFilters();
  }

  function removeExcludeTag(tag: string) {
    excludeTags = excludeTags.filter((t) => t !== tag);
    applyTagFilters();
  }

  function applyTagFilters() {
    const query = page.url.searchParams;
    query.delete("tags");
    query.delete("excludeTags");

    if (includeTags.length > 0) {
      query.set("tags", includeTags.join(","));
    }

    if (excludeTags.length > 0) {
      query.set("excludeTags", excludeTags.join(","));
    }

    goto("?" + query.toString(), { invalidateAll: true });
  }

  function clearFilters() {
    sort = "name-a-z";
    decade = "none";
    includeTags = [];
    excludeTags = [];
    tagInput = "";

    goto("/albums", { invalidateAll: true });
  }

  let hasActiveFilters = $derived(
    decade !== "none" || includeTags.length > 0 || excludeTags.length > 0,
  );

  let infoAlbumId = $state<string | null>(null);
  let infoOpen = $state(false);

  let infoAlbum = $derived(
    infoAlbumId
      ? (data.albums.find((a) => a.id === infoAlbumId) ?? null)
      : null,
  );

  function showInfo(id: string) {
    console.log("show", id);
    infoAlbumId = id;
    infoOpen = true;
  }

  function formatDate(iso: string) {
    return new Date(iso).toLocaleDateString(undefined, {
      year: "numeric",
      month: "short",
      day: "numeric",
    });
  }
</script>

<div class="flex flex-col gap-4">
  <div class="flex items-baseline gap-2">
    <h1 class="text-xl font-bold">Albums</h1>
    {#if data.page}
      <span class="text-sm text-muted-foreground">{data.page.totalItems}</span>
    {/if}
  </div>

  <div class="rounded-lg border bg-card p-3">
    <div
      class="flex flex-col gap-3 sm:flex-row sm:items-end sm:justify-between"
    >
      <div class="flex flex-1 flex-col gap-2 sm:flex-row sm:items-center">
        <Input
          class="h-9 sm:w-56"
          placeholder="Search albums..."
          bind:value={searchQuery}
          onkeydown={(e) => {
            if (e.key === "Enter") {
              updateSearch();
            }
          }}
        />
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

      <div class="flex items-center gap-1.5">
        <Button variant="outline" size="icon" href="/search/albums">
          <Search size={16} />
        </Button>
        {#if hasActiveFilters}
          <Button variant="ghost" size="sm" onclick={clearFilters}>
            <X size={14} />
            Clear
          </Button>
        {/if}
      </div>
    </div>

    <div class="mt-3 flex flex-wrap items-center gap-1.5">
      <span class="text-xs font-medium text-muted-foreground">Decade</span>
      {#each decadeTypes as d (d.value)}
        <button
          class="rounded-md border px-2 py-1 text-xs transition-colors {decade ===
          d.value
            ? 'border-primary bg-primary text-primary-foreground'
            : 'bg-transparent text-muted-foreground hover:text-foreground'}"
          onclick={() => updateDecade(d.value)}
        >
          {d.label}
        </button>
      {/each}
    </div>

    <div class="mt-3 flex flex-wrap items-center gap-1.5">
      <span class="text-xs font-medium text-muted-foreground">Tags</span>

      <div class="flex items-center gap-1">
        <button
          class="rounded-l-md border px-1.5 py-1 text-xs font-medium transition-colors {tagMode ===
          'include'
            ? 'border-primary bg-primary text-primary-foreground'
            : 'bg-transparent text-muted-foreground hover:text-foreground'}"
          onclick={() => (tagMode = "include")}
        >
          + Inc
        </button>
        <button
          class="-ml-px rounded-r-md border px-1.5 py-1 text-xs font-medium transition-colors {tagMode ===
          'exclude'
            ? 'border-destructive bg-destructive text-destructive-foreground'
            : 'bg-transparent text-muted-foreground hover:text-foreground'}"
          onclick={() => (tagMode = "exclude")}
        >
          - Exc
        </button>
      </div>

      <Input
        class="h-7 w-28 text-xs"
        placeholder="Tag name..."
        bind:value={tagInput}
        onkeydown={(e) => {
          if (e.key === "Enter") {
            addTag();
          }
        }}
      />

      <Button variant="ghost" size="icon" class="h-7 w-7" onclick={addTag}>
        <Plus size={14} />
      </Button>

      {#each includeTags as tag (tag)}
        <span
          class="flex items-center gap-0.5 rounded-full bg-primary/10 px-2 py-0.5 text-xs text-primary"
        >
          +{tag}
          <button
            class="hover:text-primary/80"
            onclick={() => removeIncludeTag(tag)}
          >
            <X size={11} />
          </button>
        </span>
      {/each}
      {#each excludeTags as tag (tag)}
        <span
          class="flex items-center gap-0.5 rounded-full bg-destructive/10 px-2 py-0.5 text-xs text-destructive"
        >
          -{tag}
          <button
            class="hover:text-destructive/80"
            onclick={() => removeExcludeTag(tag)}
          >
            <X size={11} />
          </button>
        </span>
      {/each}
    </div>
  </div>
</div>

<Spacer size="lg" />

<div
  class="grid grid-cols-2 gap-3 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-5 xl:grid-cols-6 2xl:grid-cols-7"
>
  {#each data.albums as album (album.id)}
    <div
      class="group relative flex flex-col overflow-hidden rounded-lg border bg-card transition-shadow hover:shadow-md"
    >
      <a href="/albums/{album.id}">
        <img
          src={album.coverArt.medium}
          alt={album.name}
          class="aspect-square w-full object-cover"
        />
      </a>

      <div class="flex flex-col gap-0.5 p-2">
        <a
          href="/albums/{album.id}"
          class="truncate text-sm font-medium hover:underline"
          title={album.name}
        >
          {album.name}
        </a>
        <p
          class="truncate text-xs text-muted-foreground"
          title={album.artists.map((a) => a.name).join(", ")}
        >
          {album.artists.map((a) => a.name).join(", ")}
        </p>
      </div>

      <div class="absolute right-1.5 top-1.5">
        <DropdownMenu.Root>
          <DropdownMenu.Trigger
            class={cn(
              buttonVariants({ variant: "secondary", size: "icon" }),
              "rounded-full",
            )}
          >
            <EllipsisVertical size={14} />
          </DropdownMenu.Trigger>
          <DropdownMenu.Content align="end">
            <DropdownMenu.Group>
              <DropdownMenu.Item onclick={() => showInfo(album.id)}>
                <Info size={14} />
                Show more info
              </DropdownMenu.Item>
            </DropdownMenu.Group>
          </DropdownMenu.Content>
        </DropdownMenu.Root>
      </div>
    </div>
  {/each}
</div>

<Spacer size="lg" />

<Separator />

<Spacer size="lg" />

<Pagination page={data.page} />

<Dialog.Root open={infoOpen} onOpenChange={(v) => (infoOpen = v)}>
  <Dialog.Content class="">
    <Dialog.Header>
      <Dialog.Title>Album Info</Dialog.Title>
      <Dialog.Description>
        Detailed information about the album
      </Dialog.Description>
    </Dialog.Header>

    {#if infoAlbum}
      <div class="flex flex-col gap-4 sm:flex-row">
        <div class="flex shrink-0 justify-center sm:block">
          <img
            src={infoAlbum.coverArt.large}
            alt={infoAlbum.name}
            class="h-48 w-48 rounded-lg border object-cover sm:h-44 sm:w-44"
          />
        </div>

        <div class="flex min-w-0 flex-1 flex-col gap-2">
          <div>
            <p class="text-lg font-semibold leading-tight">{infoAlbum.name}</p>
            <p class="text-sm text-muted-foreground">
              {#each infoAlbum.artists as artist, i (artist.id)}
                {#if i > 0}
                  {", "}
                {/if}
                <a
                  href="/artists/{artist.id}"
                  class="hover:underline"
                  title={artist.name}
                >
                  {artist.name}
                </a>
              {/each}
            </p>
          </div>

          <div class="grid grid-cols-[auto_1fr] gap-x-3 gap-y-1 text-sm">
            {#if infoAlbum.year}
              <span class="text-muted-foreground">Year</span>
              <span>{infoAlbum.year}</span>
            {/if}

            {#if infoAlbum.tags.length > 0}
              <span class="text-muted-foreground">Tags</span>
              <div class="flex flex-wrap gap-1">
                {#each infoAlbum.tags as tag (tag)}
                  <span class="rounded-md bg-secondary px-1.5 py-0.5 text-xs"
                    >{tag}</span
                  >
                {/each}
              </div>
            {/if}

            <span class="text-muted-foreground">Added</span>
            <span>{formatDate(infoAlbum.created)}</span>

            <span class="text-muted-foreground">Updated</span>
            <span>{formatDate(infoAlbum.updated)}</span>
          </div>
        </div>
      </div>
    {:else}
      <p class="text-sm text-muted-foreground">Album not found.</p>
    {/if}

    <Dialog.Footer>
      <Button variant="outline" onclick={() => (infoOpen = false)}>
        Close
      </Button>
    </Dialog.Footer>
  </Dialog.Content>
</Dialog.Root>
