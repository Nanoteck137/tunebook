<script lang="ts">
  import { goto } from "$app/navigation";
  import { Artist } from "$lib/api/types";
  import Image from "$lib/components/Image.svelte";
  import Pagination from "$lib/components/Pagination.svelte";
  import Spacer from "$lib/components/Spacer.svelte";
  import { cn } from "$lib/utils";
  import {
    buttonVariants,
    DropdownMenu,
    Separator,
  } from "@nanoteck137/nano-ui";
  import { EllipsisVertical, Users } from "lucide-svelte";

  let { data } = $props();
</script>

<!-- <form method="GET">
  <div class="flex flex-col gap-2">
    <Input
      type="text"
      name="filter"
      placeholder="Filter"
      value={data.filter ?? ""}
    />

    <Input
      type="text"
      name="sort"
      placeholder="Sort"
      value={data.sort ?? ""}
    />
  </div>

  {#if data.filterError}
    <p class="text-red-400">{data.filterError}</p>
  {/if}
  {#if data.sortError}
    <p class="text-red-400">{data.sortError}</p>
  {/if}
  <div class="h-2"></div>
  <Button type="submit">
    <Filter />
    Filter Tracks
  </Button>
</form>

<div class="h-2"></div> -->

{#snippet artistItem(artist: Artist)}
  <div class="py-2">
    <div class="relative flex items-center gap-2 rounded pr-2">
      <a href={`/artists/${artist.id}`}>
        <Image class="w-14 min-w-14" src={artist.coverArt.small} alt="cover" />
      </a>
      <div class="flex flex-grow flex-col">
        <div class="flex items-center gap-1">
          <a
            class="line-clamp-1 w-fit text-sm font-medium"
            href={`/artists/${artist.id}`}
            title={artist.name}
          >
            {artist.name}
          </a>
        </div>

        <p class="line-clamp-1 text-xs text-muted-foreground">
          {#if artist.tags.length > 0}
            {artist.tags.join(", ")}
          {:else}
            No Tags
          {/if}
        </p>
      </div>
      <div class="flex items-center">
        <DropdownMenu.Root>
          <DropdownMenu.Trigger
            class={cn(
              buttonVariants({ variant: "ghost", size: "icon-lg" }),
              "rounded-full",
            )}
          >
            <EllipsisVertical />
          </DropdownMenu.Trigger>
          <DropdownMenu.Content align="end">
            <DropdownMenu.Group>
              <DropdownMenu.Item
                onSelect={() => {
                  goto(`/artists/${artist.id}`);
                }}
              >
                <Users />
                Go to Artist
              </DropdownMenu.Item>
            </DropdownMenu.Group>
          </DropdownMenu.Content>
        </DropdownMenu.Root>
      </div>
    </div>
  </div>
{/snippet}

<div class="flex items-center justify-between">
  <p class="text-bold text-xl">Artists</p>
  <p class="text-sm">{data.page.totalItems} artist(s)</p>
</div>

<div class="flex flex-col">
  {#each data.artists as artist}
    {@render artistItem(artist)}
    <Separator />
  {/each}
</div>

<Spacer size="lg" />

<Pagination page={data.page} />
