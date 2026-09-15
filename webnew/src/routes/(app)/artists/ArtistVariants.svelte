<script lang="ts">
  import { ArrowUpRight, ChevronRight } from "@lucide/svelte";
  import type { Artist } from "$lib/api/types";

  type Props = {
    artists: Artist[];
  };

  const { artists }: Props = $props();

  let items = $derived(artists.slice(0, 8));

  const gridTight =
    "grid grid-cols-2 gap-3 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-5 xl:grid-cols-6 2xl:grid-cols-7";
</script>

<div class="flex flex-col gap-8">
  <div>
    <p class="mb-2 text-sm font-semibold text-muted-foreground">A · Clean tiles</p>
    <div class={gridTight}>
      {#each items as artist (artist.id)}
        <a
          href="/artists/{artist.id}"
          class="group flex flex-col gap-1.5"
          title={artist.name}
        >
          <div
            class="relative overflow-hidden rounded-xl shadow-sm ring-1 ring-black/10 transition-shadow group-hover:shadow-lg dark:ring-white/10"
          >
            <img
              class="aspect-square w-full object-cover transition-transform duration-300 group-hover:scale-105"
              src={artist.coverArt.medium}
              alt={artist.name}
            />
          </div>
          <p class="truncate px-0.5 text-sm font-medium group-hover:underline">
            {artist.name}
          </p>
          {#if artist.tags.length > 0}
            <p class="truncate px-0.5 text-xs text-muted-foreground">
              {artist.tags.slice(0, 2).join(" · ")}
            </p>
          {/if}
        </a>
      {/each}
    </div>
  </div>

  <div>
    <p class="mb-2 text-sm font-semibold text-muted-foreground">
      B · Overlay name on image
    </p>
    <div class={gridTight}>
      {#each items as artist (artist.id)}
        <a
          href="/artists/{artist.id}"
          class="group relative overflow-hidden rounded-xl transition-shadow hover:shadow-lg"
          title={artist.name}
        >
          <img
            class="aspect-square w-full object-cover transition-transform duration-300 group-hover:scale-105"
            src={artist.coverArt.medium}
            alt={artist.name}
          />
          <div
            class="absolute inset-0 bg-linear-to-t from-black/75 via-black/10 to-transparent"
          ></div>
          <div class="absolute inset-x-0 bottom-0 flex flex-col gap-0.5 p-3">
            <p class="truncate text-sm font-semibold text-white">
              {artist.name}
            </p>
            {#if artist.tags.length > 0}
              <p class="truncate text-xs text-white/70">
                {artist.tags.slice(0, 2).join(" · ")}
              </p>
            {/if}
          </div>
        </a>
      {/each}
    </div>
  </div>

  <div>
    <p class="mb-2 text-sm font-semibold text-muted-foreground">
      C · Row list
    </p>
    <div class="flex flex-col divide-y overflow-hidden rounded-lg border bg-card">
      {#each items as artist (artist.id)}
        <a
          href="/artists/{artist.id}"
          class="group flex items-center gap-3 p-2 transition-colors hover:bg-muted/50"
          title={artist.name}
        >
          <img
            class="h-12 w-12 shrink-0 rounded-md object-cover ring-1 ring-black/10 dark:ring-white/10"
            src={artist.coverArt.medium}
            alt={artist.name}
          />
          <div class="min-w-0 flex-1">
            <p class="truncate text-sm font-medium group-hover:underline">
              {artist.name}
            </p>
            {#if artist.tags.length > 0}
              <p class="truncate text-xs text-muted-foreground">
                {artist.tags.join(", ")}
              </p>
            {/if}
          </div>
          <ChevronRight
            size={16}
            class="shrink-0 text-muted-foreground transition-transform group-hover:translate-x-0.5"
          />
        </a>
      {/each}
    </div>
  </div>

  <div>
    <p class="mb-2 text-sm font-semibold text-muted-foreground">
      D · Featured tiles
    </p>
    <div class="grid grid-cols-2 gap-6 sm:grid-cols-3 md:grid-cols-4">
      {#each items as artist (artist.id)}
        <a
          href="/artists/{artist.id}"
          class="group relative flex flex-col gap-2"
          title={artist.name}
        >
          <div
            class="relative overflow-hidden rounded-2xl shadow-lg ring-1 ring-black/10 dark:ring-white/10"
          >
            <img
              class="aspect-square w-full object-cover transition-transform duration-300 group-hover:scale-105"
              src={artist.coverArt.medium}
              alt={artist.name}
            />
            <div
              class="absolute inset-0 flex items-center justify-center bg-black/40 opacity-0 transition-opacity group-hover:opacity-100"
            >
              <span
                class="flex h-10 w-10 items-center justify-center rounded-full bg-white/20 backdrop-blur-sm"
              >
                <ArrowUpRight size={18} class="text-white" />
              </span>
            </div>
          </div>
          <p class="truncate px-0.5 text-base font-semibold group-hover:underline">
            {artist.name}
          </p>
          {#if artist.tags.length > 0}
            <p class="truncate px-0.5 text-xs text-muted-foreground">
              {artist.tags.slice(0, 2).join(" · ")}
            </p>
          {/if}
        </a>
      {/each}
    </div>
  </div>
</div>