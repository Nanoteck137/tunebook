<script lang="ts">
  import { goto } from "$app/navigation";
  import { page } from "$app/state";
  import {
    Button,
    buttonVariants,
    DropdownMenu,
    Input,
    Select,
    Separator,
  } from "@nanoteck137/nano-ui";
  import {
    Check,
    EllipsisVertical,
    FileHeart,
    Plus,
    Search,
    X,
  } from "lucide-svelte";
  import { cn } from "$lib/utils";
  import Spacer from "$lib/components/Spacer.svelte";
  import { invalidateAll } from "$app/navigation";
  import { getApiClient, handleApiError } from "$lib";
  import NewPlaylistModal from "./NewPlaylistModal.svelte";
  import Image from "$lib/components/Image.svelte";
  import Pagination from "$lib/components/Pagination.svelte";
  import { sortTypes, defaultSort, type SortType } from "./types";

  let { data } = $props();
  const apiClient = getApiClient();

  let openNewPlaylistModal = $state(false);

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

  let showAll = $state(page.url.searchParams.get("all") === "true");

  function clearFilters() {
    searchQuery = "";
    showAll = false;
    sort = defaultSort;
    const query = page.url.searchParams;
    query.delete("query");
    query.delete("all");
    query.delete("sort");
    goto("?" + query.toString(), { invalidateAll: true });
  }

  let hasActiveFilters = $derived(searchQuery !== "" || showAll);

  function toggleAll() {
    showAll = !showAll;

    const query = page.url.searchParams;
    query.delete("all");

    if (showAll) {
      query.set("all", "true");
    }

    goto("?" + query.toString(), { invalidateAll: true });
  }
</script>

<div class="flex flex-col gap-4">
  <div class="flex items-baseline justify-between gap-2">
    <div class="flex items-baseline gap-2">
      <h1 class="text-xl font-bold">Playlists</h1>
      {#if data.page}
        <span class="text-sm text-muted-foreground">
          {data.page.totalItems}
        </span>
      {/if}
    </div>

    <Button variant="ghost" onclick={() => (openNewPlaylistModal = true)}>
      <Plus />
      New Playlist
    </Button>
  </div>

  <div class="rounded-lg border bg-card p-3">
    <div
      class="flex flex-col gap-3 sm:flex-row sm:items-end sm:justify-between"
    >
      <div class="flex flex-1 flex-col gap-2 sm:flex-row sm:items-center">
        <Input
          class="h-9 sm:w-56"
          placeholder="Search playlists..."
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
      <span class="text-xs font-medium text-muted-foreground">Owner</span>
      <button
        class="rounded-md border px-2 py-1 text-xs transition-colors {!showAll
          ? 'border-primary bg-primary text-primary-foreground'
          : 'bg-transparent text-muted-foreground hover:text-foreground'}"
        onclick={() => {
          if (showAll) toggleAll();
        }}
      >
        My Playlists
      </button>
      <button
        class="rounded-md border px-2 py-1 text-xs transition-colors {showAll
          ? 'border-primary bg-primary text-primary-foreground'
          : 'bg-transparent text-muted-foreground hover:text-foreground'}"
        onclick={() => {
          if (!showAll) toggleAll();
        }}
      >
        All Playlists
      </button>
    </div>
  </div>

  <div
    class="grid grid-cols-2 gap-3 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-5 xl:grid-cols-6 2xl:grid-cols-7"
  >
    {#each data.playlists as playlist (playlist.id)}
      <div
        class="group relative flex flex-col overflow-hidden rounded-lg border bg-card transition-shadow hover:shadow-md"
      >
        <a href="/playlists/{playlist.id}">
          <Image
            class="aspect-square w-full rounded-none border-0"
            src={playlist.coverArt.medium}
            alt={playlist.name}
          />
        </a>

        <div class="flex flex-col gap-0.5 p-2">
          <a
            href="/playlists/{playlist.id}"
            class="truncate text-sm font-medium hover:underline"
            title={playlist.name}
          >
            {playlist.name}
          </a>
          <p class="truncate text-xs text-muted-foreground">
            {playlist.trackCount} track{playlist.trackCount !== 1 ? "s" : ""}
          </p>
          <p
            class="flex items-center gap-1 truncate text-xs text-muted-foreground"
          >
            {#if playlist.ownerPicture}
              <img
                src={playlist.ownerPicture.small}
                alt=""
                class="h-4 w-4 rounded-full object-cover"
              />
            {/if}
            {playlist.ownerDisplayName}
          </p>
        </div>

        {#if data.user?.quickPlaylist === playlist.id}
          <div
            class="absolute left-1.5 top-1.5 flex items-center gap-1 rounded-full bg-primary px-2 py-0.5 text-xs text-primary-foreground"
          >
            <Check size={12} />
            Quick
          </div>
        {/if}

        <div class="absolute right-1.5 top-1.5">
          <DropdownMenu.Root>
            <DropdownMenu.Trigger
              class={cn(
                buttonVariants({ variant: "secondary", size: "icon" }),
                "h-7 w-7 rounded-full opacity-0 transition-opacity group-hover:opacity-100",
              )}
            >
              <EllipsisVertical size={14} />
            </DropdownMenu.Trigger>
            <DropdownMenu.Content align="end">
              <DropdownMenu.Group>
                {#if data.user?.quickPlaylist !== playlist.id}
                  <DropdownMenu.Item
                    onSelect={async () => {
                      const res = await apiClient.setQuickPlaylist({
                        playlistId: playlist.id,
                      });
                      if (!res.success) {
                        handleApiError(res.error);
                        return;
                      }

                      await invalidateAll();
                    }}
                  >
                    <FileHeart />
                    Set as Quick Playlist
                  </DropdownMenu.Item>
                {/if}
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
</div>

<NewPlaylistModal bind:open={openNewPlaylistModal} />
