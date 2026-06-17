<script lang="ts">
  import {
    ChevronDown,
    Pause,
    Play,
    SkipBack,
    SkipForward,
    Volume2,
    VolumeX,
    ListMusic,
    X,
    ListX,
    Music2,
    ArrowLeft,
  } from "lucide-svelte";
  import { formatTime } from "$lib/utils";
  import { Button, Sheet } from "@nanoteck137/nano-ui";
  import SeekSlider from "$lib/components/SeekSlider.svelte";
  import { fly } from "svelte/transition";
  import { getMusicManager, type MediaItem } from "$lib/music-manager.svelte";
  import Image from "$lib/components/Image.svelte";

  const musicManager = getMusicManager();

  let currentMediaItem = $state<MediaItem | null>(null);

  $effect(() => {
    currentMediaItem = musicManager.queue.getCurrentMediaItem();
  });

  let previousItems = $derived(musicManager.queue.getPreviousItems());
  let currentQueueItem = $derived(musicManager.queue.getCurrentItem());
  let nextItems = $derived(musicManager.queue.getNextItems());

  let view = $state<"player" | "queue">("player");
</script>

<div
  class="z-30 border-t bg-background text-foreground md:hidden"
  transition:fly={{ y: 200 }}
>
  <Sheet.Root>
    <Sheet.Trigger class="flex w-full items-center gap-3 px-3 py-2">
      <Image
        class="w-10 min-w-10 shrink-0"
        src={currentMediaItem?.coverArt}
        alt="cover"
        loading="eager"
      />

      <div class="flex min-w-0 flex-col items-start text-left">
        <p class="w-full truncate text-sm font-medium">{currentMediaItem?.name ?? "No track playing"}</p>
        <p class="w-full truncate text-xs text-muted-foreground">{currentMediaItem?.artists[0]?.name ?? ""}</p>
      </div>

      <div class="flex-grow"></div>

      {#if musicManager.loading}
        <div class="h-8 w-8 animate-pulse rounded-full bg-muted" />
      {:else if musicManager.playing}
        <button
          class="flex h-8 w-8 items-center justify-center"
          onclick={(e) => {
            e.stopPropagation();
            musicManager.pause();
          }}
        >
          <Pause size="22" />
        </button>
      {:else}
        <button
          class="flex h-8 w-8 items-center justify-center"
          onclick={(e) => {
            e.stopPropagation();
            musicManager.play();
          }}
        >
          <Play size="22" />
        </button>
      {/if}

      <ChevronDown size="20" class="shrink-0 text-muted-foreground" />
    </Sheet.Trigger>

    <Sheet.Content
      side="bottom"
      class="flex h-full max-h-full flex-col bg-gradient-to-b from-zinc-950 to-background"
    >
      {#if view === "player"}
        <div class="flex flex-1 flex-col px-4 pb-8 pt-2">
          <!-- Header -->
          <div class="flex items-center justify-between pb-4">
            <Sheet.Close class="flex items-center gap-1 text-sm text-muted-foreground transition-colors hover:text-foreground">
              <ChevronDown size="20" />
            </Sheet.Close>
            <p class="truncate text-center text-sm font-medium">
              {currentMediaItem?.name ?? ""}
            </p>
            <button
              class="text-muted-foreground transition-colors hover:text-foreground"
              onclick={() => (view = "queue")}
            >
              <ListMusic size="20" />
            </button>
          </div>

          <!-- Cover art + info -->
          <div class="flex flex-1 flex-col items-center justify-center gap-4">
            <Image
              class="w-72 max-w-full shadow-xl"
              src={currentMediaItem?.coverArt}
              alt="Track Cover Art"
              loading="eager"
            />

            <div class="flex w-full flex-col items-center text-center">
              <p class="text-lg font-medium">{currentMediaItem?.name}</p>
              <p class="text-sm text-muted-foreground">
                {#if currentMediaItem}
                  {#each currentMediaItem.artists as artist, i}
                    {#if i > 0}, {/if}
                    <a href="/artists/{artist.id}" class="hover:underline">{artist.name}</a>
                  {/each}
                {/if}
              </p>
            </div>
          </div>

          <!-- Seek bar -->
          <div class="flex w-full flex-col gap-1 pb-4">
            <SeekSlider
              value={Number.isNaN(musicManager.duration) ? 0 : musicManager.currentTime / musicManager.duration}
              onValue={(p) => {
                musicManager.setPosition(p * musicManager.duration);
              }}
              buffered={musicManager.buffered}
            />
            <div class="flex justify-between text-xs text-muted-foreground">
              <span class="tabular-nums">{formatTime(musicManager.currentTime)}</span>
              <span class="tabular-nums">{formatTime(Number.isNaN(musicManager.duration) ? 0 : musicManager.duration)}</span>
            </div>
          </div>

          <!-- Controls -->
          <div class="flex items-center justify-center gap-6 pb-4">
            <button
              class="text-muted-foreground transition-colors hover:text-foreground"
              onclick={() => musicManager.previousTrack()}
            >
              <SkipBack size="28" />
            </button>

            {#if musicManager.loading}
              <div class="h-14 w-14 animate-pulse rounded-full bg-muted" />
            {:else if musicManager.playing}
              <button
                class="flex h-14 w-14 items-center justify-center rounded-full bg-foreground text-background transition-colors hover:scale-105"
                onclick={() => musicManager.pause()}
              >
                <Pause size="28" />
              </button>
            {:else}
              <button
                class="flex h-14 w-14 items-center justify-center rounded-full bg-foreground text-background transition-colors hover:scale-105"
                onclick={() => musicManager.play()}
              >
                <Play size="28" />
              </button>
            {/if}

            <button
              class="text-muted-foreground transition-colors hover:text-foreground"
              onclick={() => musicManager.nextTrack()}
            >
              <SkipForward size="28" />
            </button>
          </div>

          <!-- Volume -->
          <div class="flex w-full max-w-48 items-center gap-2 self-center">
            <button
              class="shrink-0 text-muted-foreground transition-colors hover:text-foreground"
              onclick={() => {
                musicManager.muted = !musicManager.muted;
              }}
            >
              {#if musicManager.muted || musicManager.volume === 0}
                <VolumeX size="18" />
              {:else}
                <Volume2 size="18" />
              {/if}
            </button>

            <div class="flex-1">
              <SeekSlider
                growOnHover={false}
                value={musicManager.muted ? 0 : musicManager.volume}
                onValue={(p) => {
                  musicManager.volume = p;
                  musicManager.muted = p === 0;
                }}
              />
            </div>
          </div>
        </div>
      {:else}
        <!-- Queue view -->
        <div class="flex flex-1 flex-col px-4 pb-8 pt-2">
          <div class="flex items-center gap-3 pb-4 pt-1">
            <button
              class="text-muted-foreground transition-colors hover:text-foreground"
              onclick={() => (view = "player")}
            >
              <ArrowLeft size="22" />
            </button>
            <p class="text-base font-semibold">Queue</p>
            <div class="flex-1" />
            <Button
              class="rounded-full"
              variant="ghost"
              size="icon"
              onclick={() => musicManager.clearQueue()}
            >
              <ListX size="16" />
            </Button>
          </div>

          <div class="flex-1 overflow-y-auto overscroll-contain">
            <div class="flex flex-col gap-1 pr-2">
              <!-- Played -->
              {#if previousItems.length > 0}
                <p class="py-1 text-xs font-semibold uppercase tracking-wider text-muted-foreground">
                  Played
                </p>
                {#each previousItems as mediaItem, i (mediaItem.trackId)}
                  {@const queueIndex = i}
                  <div class="group flex items-center gap-3 rounded-md p-2 transition-colors hover:bg-accent/50">
                    <button
                      class="shrink-0"
                      onclick={async () => {
                        await musicManager.setQueueIndex(queueIndex);
                        musicManager.play();
                      }}
                    >
                      <Image class="w-10 min-w-10 rounded" src={mediaItem.coverArt} alt="cover" />
                    </button>
                    <div class="flex min-w-0 flex-1 flex-col">
                      <p class="truncate text-sm font-medium">{mediaItem.name}</p>
                      <p class="truncate text-xs text-muted-foreground">
                        {#each mediaItem.artists as artist, j}
                          {#if j > 0}, {/if}{artist.name}
                        {/each}
                      </p>
                    </div>
                  </div>
                {/each}
              {/if}

              <!-- Now Playing -->
              {#if currentQueueItem}
                <p class="pt-3 pb-1 text-xs font-semibold uppercase tracking-wider text-muted-foreground">
                  Now Playing
                </p>
                <div class="flex items-center gap-3 rounded-md bg-accent/50 p-2">
                  <div class="relative shrink-0">
                    <Image class="w-10 min-w-10 rounded" src={currentQueueItem.coverArt} alt="cover" />
                    <div class="absolute inset-0 flex items-center justify-center">
                      <Music2 class="text-primary" size="14" />
                    </div>
                  </div>
                  <div class="flex min-w-0 flex-col">
                    <p class="truncate text-sm font-medium">{currentQueueItem.name}</p>
                    <p class="truncate text-xs text-muted-foreground">
                      {#each currentQueueItem.artists as artist, j}
                        {#if j > 0}, {/if}{artist.name}
                      {/each}
                    </p>
                  </div>
                </div>
              {/if}

              <!-- Next in queue -->
              {#if nextItems.length > 0}
                <p class="pt-3 pb-1 text-xs font-semibold uppercase tracking-wider text-muted-foreground">
                  Next in queue
                </p>
                {#each nextItems as mediaItem, i (mediaItem.trackId)}
                  {@const queueIndex = musicManager.queue.index + 1 + i}
                  <div class="group flex items-center gap-3 rounded-md p-2 transition-colors hover:bg-accent/50">
                    <span class="w-5 text-right text-xs tabular-nums text-muted-foreground">{i + 1}</span>
                    <button
                      class="shrink-0"
                      onclick={async () => {
                        await musicManager.setQueueIndex(queueIndex);
                        musicManager.play();
                      }}
                    >
                      <Image class="w-10 min-w-10 rounded" src={mediaItem.coverArt} alt="cover" />
                    </button>
                    <div class="flex min-w-0 flex-1 flex-col">
                      <p class="truncate text-sm font-medium">{mediaItem.name}</p>
                      <p class="truncate text-xs text-muted-foreground">
                        {#each mediaItem.artists as artist, j}
                          {#if j > 0}, {/if}{artist.name}
                        {/each}
                      </p>
                    </div>
                    <button
                      class="shrink-0 rounded-full p-1 text-muted-foreground opacity-0 transition-opacity hover:text-foreground group-hover:opacity-100"
                      onclick={() => musicManager.removeQueueItem(queueIndex)}
                    >
                      <X size="14" />
                    </button>
                  </div>
                {/each}
              {/if}

              {#if !currentQueueItem}
                <p class="py-4 text-center text-sm text-muted-foreground">Queue is empty</p>
              {/if}
            </div>
          </div>
        </div>
      {/if}
    </Sheet.Content>
  </Sheet.Root>
</div>
