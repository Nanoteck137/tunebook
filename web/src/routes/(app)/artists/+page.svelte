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
  import { page } from "$app/state";
  import {
    sortTypes,
    defaultSort,
    type SortType,
  } from "./types";

  let { data } = $props();

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

  let searchQuery = $state(page.url.searchParams.get("query") ?? "");
  function updateSearch() {
    const query = page.url.searchParams;
    query.delete("query");

    if (searchQuery) {
      query.set("query", searchQuery);
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
    searchQuery = "";
    sort = defaultSort;
    includeTags = [];
    excludeTags = [];
    const query = page.url.searchParams;
    query.delete("query");
    query.delete("sort");
    query.delete("tags");
    query.delete("excludeTags");
    goto("?" + query.toString(), { invalidateAll: true });
  }

  let hasActiveFilters = $derived(
    searchQuery !== "" ||
      sort !== defaultSort ||
      includeTags.length > 0 ||
      excludeTags.length > 0,
  );

  let infoArtistId = $state<string | null>(null);
  let infoOpen = $state(false);

  let infoArtist = $derived(
    infoArtistId
      ? (data.artists.find((a) => a.id === infoArtistId) ?? null)
      : null,
  );

  function showInfo(id: string) {
    infoArtistId = id;
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
  <div class="flex items-baseline justify-between gap-2">
    <div class="flex items-baseline gap-2">
      <h1 class="text-xl font-bold">Artists</h1>
      {#if data.page}
        <span class="text-sm text-muted-foreground">{data.page.totalItems}</span>
      {/if}
    </div>

    <Button variant="outline" size="icon" href="/search/artists">
      <Search size={16} />
    </Button>
  </div>

  <div class="rounded-lg border bg-card p-3">
    <div
      class="flex flex-col gap-3 sm:flex-row sm:items-end sm:justify-between"
    >
      <div class="flex flex-1 flex-col gap-2 sm:flex-row sm:items-center">
        <Input
          class="h-9 sm:w-56"
          placeholder="Search artists..."
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
        <Button variant="outline" size="icon" onclick={updateSearch}>
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
            ? 'border-primary bg-primary text-primary-foreground'
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
          if (e.key === "Enter") addTag();
        }}
      />
      <Button variant="ghost" size="icon" class="h-7 w-7" onclick={addTag}>
        <Plus size={14} />
      </Button>
    </div>

    {#if includeTags.length > 0 || excludeTags.length > 0}
      <div class="mt-2 flex flex-wrap items-center gap-1.5">
        {#each includeTags as tag (tag)}
          <span
            class="flex items-center gap-1 rounded-md bg-primary/10 px-2 py-0.5 text-xs text-primary"
          >
            +{tag}
            <button onclick={() => removeIncludeTag(tag)} class="hover:text-destructive">
              <X size={12} />
            </button>
          </span>
        {/each}
        {#each excludeTags as tag (tag)}
          <span
            class="flex items-center gap-1 rounded-md bg-destructive/10 px-2 py-0.5 text-xs text-destructive"
          >
            -{tag}
            <button onclick={() => removeExcludeTag(tag)} class="hover:text-destructive">
              <X size={12} />
            </button>
          </span>
        {/each}
      </div>
    {/if}
  </div>
</div>

<Spacer size="lg" />

<div
  class="grid grid-cols-2 gap-3 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-5 xl:grid-cols-6 2xl:grid-cols-7"
>
  {#each data.artists as artist (artist.id)}
    <div
      class="group relative flex flex-col overflow-hidden rounded-lg border bg-card transition-shadow hover:shadow-md"
    >
      <a href="/artists/{artist.id}">
        <img
          src={artist.coverArt.medium}
          alt={artist.name}
          class="aspect-square w-full object-cover"
        />
      </a>

      <div class="flex flex-col gap-0.5 p-2">
        <a
          href="/artists/{artist.id}"
          class="truncate text-sm font-medium hover:underline"
          title={artist.name}
        >
          {artist.name}
        </a>
      </div>

      <div class="absolute right-1.5 top-1.5">
        <DropdownMenu.Root>
          <DropdownMenu.Trigger
            class={cn(
              buttonVariants({ variant: "secondary", size: "icon" }),
              "rounded-full opacity-70 hover:opacity-100",
            )}
          >
            <EllipsisVertical size={14} />
          </DropdownMenu.Trigger>
          <DropdownMenu.Content align="end">
            <DropdownMenu.Group>
              <DropdownMenu.Item onclick={() => showInfo(artist.id)}>
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
  <Dialog.Content class="sm:max-w-lg">
    <Dialog.Header>
      <Dialog.Title>Artist Info</Dialog.Title>
      <Dialog.Description>
        Detailed information about the artist
      </Dialog.Description>
    </Dialog.Header>

    {#if infoArtist}
      <div class="flex flex-col gap-4 sm:flex-row">
        <div class="flex shrink-0 justify-center sm:block">
          <img
            src={infoArtist.coverArt.large}
            alt={infoArtist.name}
            class="h-48 w-48 rounded-lg border object-cover sm:h-44 sm:w-44"
          />
        </div>

        <div class="flex min-w-0 flex-1 flex-col gap-2">
          <div>
            <p class="text-lg font-semibold leading-tight">
              {infoArtist.name}
            </p>
          </div>

          <div class="grid grid-cols-[auto_1fr] gap-x-3 gap-y-1 text-sm">
            {#if infoArtist.tags.length > 0}
              <span class="text-muted-foreground">Tags</span>
              <div class="flex flex-wrap gap-1">
                {#each infoArtist.tags as tag (tag)}
                  <span class="rounded-md bg-secondary px-1.5 py-0.5 text-xs"
                    >{tag}</span
                  >
                {/each}
              </div>
            {/if}

            <span class="text-muted-foreground">Added</span>
            <span>{formatDate(infoArtist.created)}</span>

            <span class="text-muted-foreground">Updated</span>
            <span>{formatDate(infoArtist.updated)}</span>
          </div>
        </div>
      </div>
    {:else}
      <p class="text-sm text-muted-foreground">Artist not found.</p>
    {/if}

    <Dialog.Footer>
      <Button variant="outline" onclick={() => (infoOpen = false)}>
        Close
      </Button>
    </Dialog.Footer>
  </Dialog.Content>
</Dialog.Root>
