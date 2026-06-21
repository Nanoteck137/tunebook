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

  let { data } = $props();

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
  <div class="flex items-baseline gap-2">
    <h1 class="text-xl font-bold">Artists</h1>
    {#if data.page}
      <span class="text-sm text-muted-foreground">{data.page.totalItems}</span>
    {/if}
  </div>

  <div class="rounded-lg border bg-card p-3">
    <div
      class="flex flex-col gap-3 sm:flex-row sm:items-end sm:justify-between"
    >
      <div class="flex flex-1 flex-col gap-2 sm:flex-row sm:items-center">
        <Input class="h-9 sm:w-56" placeholder="Search artists..." disabled />
        <Select.Root type="single" allowDeselect={false}>
          <Select.Trigger class="h-9 w-full sm:w-40">
            Name (A-Z)
          </Select.Trigger>
          <Select.Content>
            <Select.Item value="name-a-z" label="Name (A-Z)" />
            <Select.Item value="name-z-a" label="Name (Z-A)" />
            <Select.Item value="created-new" label="Created (New–Old)" />
            <Select.Item value="created-old" label="Created (Old-New)" />
          </Select.Content>
        </Select.Root>
      </div>

      <Button variant="outline" size="icon" href="/search/artists">
        <Search size={16} />
      </Button>
    </div>

    <div class="mt-3 flex flex-wrap items-center gap-1.5">
      <span class="text-xs font-medium text-muted-foreground">Tags</span>
      <div class="flex items-center gap-1">
        <button
          class="rounded-l-md border border-primary bg-primary px-1.5 py-1 text-xs font-medium text-primary-foreground"
        >
          + Inc
        </button>
        <button
          class="-ml-px rounded-r-md border bg-transparent px-1.5 py-1 text-xs font-medium text-muted-foreground"
        >
          - Exc
        </button>
      </div>
      <Input class="h-7 w-28 text-xs" placeholder="Tag name..." disabled />
      <Button variant="ghost" size="icon" class="h-7 w-7" disabled>
        <Plus size={14} />
      </Button>
    </div>
  </div>
</div>

<Spacer size="lg" />

<div
  class="grid grid-cols-2 gap-3 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-5 xl:grid-cols-6 2xl:grid-cols-7"
>
  {#each data.artists as artist}
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
                {#each infoArtist.tags as tag}
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
