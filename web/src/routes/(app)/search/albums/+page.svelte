<script lang="ts">
  import { goto } from "$app/navigation";
  import { page } from "$app/state";
  import {
    Button,
    buttonVariants,
    Dialog,
    DropdownMenu,
    Input,
    Separator,
  } from "@nanoteck137/nano-ui";
  import { EllipsisVertical, Info, Search, X } from "lucide-svelte";
  import Pagination from "$lib/components/Pagination.svelte";
  import Spacer from "$lib/components/Spacer.svelte";
  import { cn } from "$lib/utils";
  import { onMount } from "svelte";

  let { data } = $props();

  async function doSearch(query: string) {
    await goto(`/search/albums?query=${query}`, {
      invalidateAll: true,
      keepFocus: true,
      replaceState: true,
    });
  }

  function clearSearch() {
    value = "";
    doSearch("");
  }

  let initialValue = $state("");
  let value = "";

  onMount(() => {
    initialValue = data.query;
    value = data.query;
  });

  let timer: ReturnType<typeof setTimeout>;
  function onInput(e: Event) {
    const target = e.target as HTMLInputElement;
    const current = target.value;
    value = current;

    clearTimeout(timer);
    timer = setTimeout(() => {
      doSearch(current);
    }, 500);
  }

  let infoAlbumId = $state<string | null>(null);
  let infoOpen = $state(false);

  let infoAlbum = $derived(
    infoAlbumId
      ? (data.albums.find((a) => a.id === infoAlbumId) ?? null)
      : null,
  );

  function showInfo(id: string) {
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

  const tabs = [
    { label: "All", href: "/search" },
    { label: "Tracks", href: "/search/tracks" },
    { label: "Artists", href: "/search/artists" },
    { label: "Albums", href: "/search/albums" },
    { label: "Playlists", href: "/search/playlists" },
    { label: "Users", href: "/search/users" },
  ];
</script>

<svelte:head>
  <title>Search Albums - Tunebook</title>
</svelte:head>

<div class="flex flex-col gap-6">
  <div>
    <h1 class="mb-4 text-xl font-bold">Search Albums</h1>

    <form
      action=""
      method="get"
      onsubmit={(e) => {
        e.preventDefault();
        clearTimeout(timer);
        doSearch(value);
      }}
    >
      <div class="rounded-lg border bg-card p-3">
        <div class="flex items-center gap-2">
          <div class="relative flex-1">
            <Search
              size={16}
              class="pointer-events-none absolute left-3 top-1/2 -translate-y-1/2 text-muted-foreground"
            />
            <Input
              id="query"
              name="query"
              placeholder="Search albums..."
              autocomplete="off"
              value={initialValue}
              oninput={onInput}
              class="pl-9"
            />
          </div>
          <Button type="submit">Search</Button>
          {#if data.query}
            <Button variant="ghost" size="icon" onclick={clearSearch}>
              <X size={16} />
            </Button>
          {/if}
        </div>
      </div>
    </form>
  </div>

  <nav class="flex flex-wrap gap-1">
    {#each tabs as { label, href } ({ label })}
      <a
        href="{href}?query={data.query}"
        class={cn(
          "rounded-md px-3 py-1.5 text-sm font-medium transition-colors",
          page.url.pathname === href
            ? "bg-primary text-primary-foreground"
            : "text-muted-foreground hover:bg-accent hover:text-accent-foreground",
        )}
      >
        {label}
      </a>
    {/each}
  </nav>

  {#if data.query && data.albums.length === 0}
    <p class="py-12 text-center text-sm text-muted-foreground">
      No albums found for "{data.query}".
    </p>
  {/if}

  {#if data.albums.length > 0}
    {#if data.page}
      <div class="flex items-baseline gap-2">
        <span class="text-sm text-muted-foreground">
          {data.page.totalItems} album(s)
        </span>
      </div>
    {/if}

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

    {#if data.page}
      <Spacer size="lg" />
      <Separator />
      <Spacer size="lg" />

      <Pagination page={data.page} />
    {/if}
  {/if}
</div>

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
