<script lang="ts">
  import type { Track } from "$lib/api/types";
  import ArtistList from "$lib/components/ArtistList.svelte";
  import Image from "$lib/components/Image.svelte";
  import { cn } from "$lib/utils";
  import { Heart, Play, Star } from "lucide-svelte";
  import type { Snippet } from "svelte";
  import { getFavorites } from "$lib/favorites.svelte";
  import { getQuickPlaylist } from "$lib/quick-playlist.svelte";

  type Props = {
    class?: string;
    showNumber?: boolean;
    displayOrder?: boolean;
    track: Track;
    children?: Snippet;

    onPlayClicked?: () => void;
  };

  const {
    class: className,
    showNumber,
    displayOrder,
    track,
    children,
    onPlayClicked,
  }: Props = $props();
  const favoriteManager = getFavorites();
  const quickPlaylistManager = getQuickPlaylist();

  let isFav = $derived(favoriteManager.hasTrack(track.id));
  let isQuick = $derived(
    quickPlaylistManager.playlist !== null &&
      quickPlaylistManager.hasTrack(track.id),
  );
</script>

<div
  class={cn(
    "group flex items-center gap-2 p-2 hover:bg-off-background1 group-even:bg-off-background2 group-even:hover:bg-off-background1 has-[[data-state='open']]:bg-off-background1 group-even:has-[[data-state='open']]:bg-off-background1",
    className,
  )}
>
  <button
    class="shrink-0"
    onclick={() => onPlayClicked?.()}
    aria-label="Play {track.name}"
  >
    {#if showNumber}
      <div class="flex min-h-10 min-w-10 flex-col items-end justify-center">
        <p class="text-right text-sm font-medium group-hover:hidden">
          {track.number}.
        </p>
        <div class="hidden items-center justify-center group-hover:flex">
          <Play size={20} />
        </div>
      </div>
    {:else}
      <div class="relative">
        <Image class="w-14 min-w-14" src={track.coverArt.small} alt="cover" />
        <div
          class="absolute inset-0 flex items-center justify-center rounded border bg-black/60 opacity-0 transition-opacity group-hover:opacity-100"
        >
          <Play size={20} class="text-white" />
        </div>
      </div>
    {/if}
  </button>

  <div class="flex min-w-0 flex-1 flex-col gap-0.5">
    <p class="truncate text-sm font-medium" title={track.name}>
      {#if isFav}
        <Heart size={12} class="mr-0.5 inline fill-primary text-primary sm:hidden" />
      {/if}
      {#if isQuick}
        <Star size={12} class="mr-0.5 inline fill-primary text-primary sm:hidden" />
      {/if}
      {#if displayOrder}
        {track.order}.
      {/if}
      {track.name}
    </p>

    <ArtistList class="text-muted-foreground" artists={track.artists} />

    {#if track.tags.length > 0}
      <div
        class="flex min-w-0 gap-1 overflow-hidden text-ellipsis whitespace-nowrap"
      >
        {#each track.tags as tag (tag)}
          <span
            class="shrink-0 rounded bg-secondary/50 px-1 py-0.5 text-[10px] text-muted-foreground"
            >{tag}</span
          >
        {/each}
      </div>
    {/if}
  </div>

  <div class="flex shrink-0 items-center gap-0.5 sm:gap-1">
    {@render children?.()}
  </div>
</div>
