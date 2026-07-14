<script lang="ts">
  import SeekSlider from "$lib/components/SeekSlider.svelte";
  import { Button, ScrollArea, Sheet } from "@nanoteck137/nano-ui";
  import { formatTime } from "$lib/utils";
  import {
    ListX,
    ListMusic,
    Pause,
    Play,
    SkipBack,
    SkipForward,
    Volume2,
    VolumeX,
    X,
    Music2,
  } from "lucide-svelte";
  import { getMusicManager, type MediaItem } from "$lib/music-manager.svelte";
  import Spinner from "$lib/components/Spinner.svelte";
  import Image from "$lib/components/Image.svelte";

  const musicManager = getMusicManager();

  let currentMediaItem = $state<MediaItem | null>(null);
  let previousItems = $state<MediaItem[]>([]);
  let previousStartIndex = $state(0);
  let currentQueueItem = $state<MediaItem | null>(null);
  let nextItems = $state<MediaItem[]>([]);

  $effect(() => {
    currentMediaItem = musicManager.currentItem;
    currentQueueItem = musicManager.currentItem;

    musicManager.queue.getPreviousItems(0, 10).then((result) => {
      previousItems = result.items;
      previousStartIndex = result.startIndex;
    });
    musicManager.queue.getNextItems(0, 50).then((items) => {
      nextItems = items;
    });
  });
</script>

{#snippet queueSheet()}
  <Sheet.Root>
    <Sheet.Trigger>
      <ListMusic size="20" />
    </Sheet.Trigger>
    <Sheet.Content side="right">
      <div class="flex items-center justify-between pb-4">
        <p class="text-base font-semibold">Queue</p>
        <div class="flex items-center">
          <Button variant="ghost" size="sm" href="/queue">View all</Button>
          <Button
            class="rounded-full"
            variant="ghost"
            size="icon"
            onclick={async () => {
              await musicManager.clearQueue();
            }}
          >
            <ListX />
          </Button>
        </div>
      </div>

      <ScrollArea class="h-full pb-6">
        <div class="mr-3 flex flex-col gap-3">
          <!-- Played -->
          {#if previousItems.length > 0}
            <div>
              <p
                class="mb-2 text-xs font-semibold uppercase tracking-wider text-muted-foreground"
              >
                Played
              </p>
              <div class="flex flex-col gap-1">
                {#each previousItems as mediaItem, i (mediaItem.trackId)}
                  {@const queueIndex = previousStartIndex + i}
                  <div
                    class="group flex items-center gap-3 rounded-md p-2 transition-colors hover:bg-accent/50"
                  >
                    <button
                      class="shrink-0"
                      onclick={async () => {
                        await musicManager.setQueueIndex(queueIndex);
                        musicManager.play();
                      }}
                    >
                      <Image
                        class="w-10 min-w-10 rounded"
                        src={mediaItem.coverArt}
                        alt="cover"
                      />
                    </button>
                    <div class="flex min-w-0 flex-1 flex-col">
                      <p
                        class="truncate text-sm font-medium"
                        title={mediaItem.name}
                      >
                        {mediaItem.name}
                      </p>
                      <p class="truncate text-xs text-muted-foreground">
                        {#each mediaItem.artists as artist, j (artist.id)}
                          {#if j > 0},
                          {/if}{artist.name}
                        {/each}
                      </p>
                    </div>
                  </div>
                {/each}
              </div>
            </div>
          {/if}

          <!-- Now Playing -->
          {#if currentQueueItem}
            <div>
              <p
                class="mb-2 text-xs font-semibold uppercase tracking-wider text-muted-foreground"
              >
                Now Playing
              </p>
              <div class="flex items-center gap-3 rounded-md bg-accent/50 p-3">
                <div class="relative shrink-0">
                  <Image
                    class="w-12 min-w-12 rounded"
                    src={currentQueueItem.coverArt}
                    alt="cover"
                  />
                  <div
                    class="absolute inset-0 flex items-center justify-center"
                  >
                    <Music2 class="text-primary" size="16" />
                  </div>
                </div>
                <div class="flex min-w-0 flex-col">
                  <p
                    class="truncate text-sm font-medium"
                    title={currentQueueItem.name}
                  >
                    {currentQueueItem.name}
                  </p>
                  <p class="truncate text-xs text-muted-foreground">
                    {#each currentQueueItem.artists as artist, j (artist.id)}
                      {#if j > 0},
                      {/if}{artist.name}
                    {/each}
                  </p>
                </div>
              </div>
            </div>
          {/if}

          <!-- Next in queue -->
          {#if nextItems.length > 0}
            <div>
              <p
                class="mb-2 text-xs font-semibold uppercase tracking-wider text-muted-foreground"
              >
                Next in queue
              </p>
              <div class="flex flex-col gap-1">
                {#each nextItems as mediaItem, i (mediaItem.trackId)}
                  {@const queueIndex = musicManager.queue.index + 1 + i}
                  <div
                    class="group flex items-center gap-3 rounded-md p-2 transition-colors hover:bg-accent/50"
                  >
                    <span
                      class="w-5 text-right text-xs tabular-nums text-muted-foreground"
                      >{i + 1}</span
                    >
                    <button
                      class="shrink-0"
                      onclick={async () => {
                        await musicManager.setQueueIndex(queueIndex);
                        musicManager.play();
                      }}
                    >
                      <Image
                        class="w-10 min-w-10 rounded"
                        src={mediaItem.coverArt}
                        alt="cover"
                      />
                    </button>
                    <div class="flex min-w-0 flex-1 flex-col">
                      <p
                        class="truncate text-sm font-medium"
                        title={mediaItem.name}
                      >
                        {mediaItem.name}
                      </p>
                      <p class="truncate text-xs text-muted-foreground">
                        {#each mediaItem.artists as artist, j (artist.id)}
                          {#if j > 0},
                          {/if}{artist.name}
                        {/each}
                      </p>
                    </div>
                    <button
                      class="shrink-0 rounded-full p-1 text-muted-foreground opacity-0 transition-opacity hover:text-foreground group-hover:opacity-100"
                      onclick={async () => {
                        await musicManager.removeQueueItem(queueIndex);
                      }}
                    >
                      <X size="14" />
                    </button>
                  </div>
                {/each}
              </div>
            </div>
          {/if}

          {#if !currentQueueItem}
            <p class="py-8 text-center text-sm text-muted-foreground">
              Queue is empty
            </p>
          {/if}
        </div>
      </ScrollArea>
    </Sheet.Content>
  </Sheet.Root>
{/snippet}

<div
  class="container relative z-30 hidden h-[72px] border-t bg-background md:block"
>
  <div class="absolute -top-1.5 left-0 right-0 px-0">
    <SeekSlider
      value={Number.isNaN(musicManager.duration)
        ? 0
        : musicManager.currentTime / musicManager.duration}
      onValue={(p) => {
        musicManager.setPosition(p * musicManager.duration);
      }}
      buffered={musicManager.buffered}
    />
  </div>

  <div class="flex h-full items-center justify-between gap-4">
    <!-- Left: Track info -->
    <div class="flex min-w-0 basis-1/4 items-center gap-3">
      <a
        href={currentMediaItem ? `/albums/${currentMediaItem.album.id}` : "#"}
        class="shrink-0"
      >
        <Image
          class="w-12 min-w-12"
          src={currentMediaItem?.coverArt}
          alt="cover"
          loading="eager"
        />
      </a>
      <div class="flex min-w-0 flex-col">
        <a
          href={currentMediaItem
            ? `/albums/${currentMediaItem.album.id}`
            : "#"}
          class="truncate text-sm font-medium hover:underline"
          title={currentMediaItem?.name}
        >
          {currentMediaItem?.name ?? "No track playing"}
        </a>
        <p class="truncate text-xs text-muted-foreground">
          {#if currentMediaItem}
            {#each currentMediaItem.artists as artist, i (artist.id)}
              {#if i > 0},
              {/if}
              <a href="/artists/{artist.id}" class="hover:underline"
                >{artist.name}</a
              >
            {/each}
          {/if}
        </p>
      </div>
    </div>

    <!-- Center: Controls + time -->
    <div class="flex flex-col items-center gap-0.5">
      <div class="flex items-center gap-3">
        <button
          class="text-muted-foreground transition-colors hover:text-foreground"
          onclick={() => musicManager.previousTrack()}
        >
          <SkipBack size="20" />
        </button>

        {#if musicManager.loading}
          <Spinner class="h-8 w-8" />
        {:else if musicManager.playing}
          <button
            class="flex h-8 w-8 items-center justify-center rounded-full bg-foreground text-background transition-colors hover:scale-105"
            onclick={() => musicManager.pause()}
          >
            <Pause size="18" />
          </button>
        {:else}
          <button
            class="flex h-8 w-8 items-center justify-center rounded-full bg-foreground text-background transition-colors hover:scale-105"
            onclick={() => musicManager.play()}
          >
            <Play size="18" />
          </button>
        {/if}

        <button
          class="text-muted-foreground transition-colors hover:text-foreground"
          onclick={() => musicManager.nextTrack()}
        >
          <SkipForward size="20" />
        </button>
      </div>

      <div class="flex items-center gap-1 text-xs text-muted-foreground">
        <span class="min-w-[32px] text-right tabular-nums"
          >{formatTime(musicManager.currentTime)}</span
        >
        <span class="text-[10px]">/</span>
        <span class="min-w-[32px] text-left tabular-nums"
          >{formatTime(
            Number.isNaN(musicManager.duration) ? 0 : musicManager.duration,
          )}</span
        >
      </div>
    </div>

    <!-- Right: Volume + Queue -->
    <div class="flex basis-1/4 items-center justify-end gap-2">
      <button
        class="text-muted-foreground transition-colors hover:text-foreground"
        onclick={() => {
          musicManager.muted = !musicManager.muted;
        }}
      >
        {#if musicManager.muted || musicManager.volume === 0}
          <VolumeX size="20" />
        {:else}
          <Volume2 size="20" />
        {/if}
      </button>

      <SeekSlider
        class="w-24"
        growOnHover={false}
        value={musicManager.muted ? 0 : musicManager.volume}
        onValue={(p) => {
          musicManager.volume = p;
          musicManager.muted = p === 0;
        }}
      />

      {@render queueSheet()}
    </div>
  </div>
</div>
