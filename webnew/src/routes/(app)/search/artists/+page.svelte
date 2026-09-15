<script lang="ts">
  import { goto } from "$app/navigation";
  import {
    Button,
    buttonVariants,
    Dialog,
    DropdownMenu,
    Separator,
  } from "$lib/components/ui";
  import { EllipsisVertical, Info } from "@lucide/svelte";
  import Pagination from "$lib/components/Pagination.svelte";
  import Spacer from "$lib/components/Spacer.svelte";
  import { cn } from "$lib/utils";
  import { onMount } from "svelte";
  import SearchBarHeader from "../SearchBarHeader.svelte";

  let { data } = $props();

  async function doSearch(query: string) {
    await goto(`/search/artists?query=${query}`, {
      invalidateAll: true,
      keepFocus: true,
      replaceState: true,
    });
  }

  function clearSearch() {
    value = "";
    doSearch("");
  }

  let value = $state("");

  onMount(() => {
    value = data.query;
  });

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

<svelte:head>
  <title>Search Artists - Tunebook</title>
</svelte:head>

<div class="flex flex-col gap-6">
  <SearchBarHeader
    searchBarPlaceholder="Search artists..."
    {value}
    setValue={(v) => {
      value = v;
    }}
    search={doSearch}
    searchWithValue={() => {
      doSearch(value);
    }}
    {clearSearch}
  />

  {#if data.query && data.artists.length === 0}
    <p class="py-12 text-center text-sm text-muted-foreground">
      No artists found for "{data.query}".
    </p>
  {/if}

  {#if data.artists.length > 0}
    {#if data.page}
      <div class="flex items-baseline gap-2">
        <span class="text-sm text-muted-foreground">
          {data.page.totalItems} artist(s)
        </span>
      </div>
    {/if}

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

    {#if data.page}
      <Spacer size="lg" />
      <Separator />
      <Spacer size="lg" />

      <Pagination page={data.page} />
    {/if}
  {/if}
</div>

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
                  <span class="rounded-md bg-secondary px-1.5 py-0.5 text-xs">
                    {tag}
                  </span>
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
