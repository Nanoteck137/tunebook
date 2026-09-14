<script lang="ts">
  import {
    ChevronDown,
    ListMusic,
    ListX,
    Pause,
    Play,
    SkipBack,
    SkipForward,
    Volume2,
    VolumeX,
    X,
  } from "@lucide/svelte";
  import { crossfade, fade, fly, slide } from "svelte/transition";
  import { formatTime } from "$lib/utils";
  import SeekSlider from "$lib/components/SeekSlider.svelte";
  import { getMusicManager, type MediaItem } from "$lib/music-manager.svelte";
  import Image from "$lib/components/Image.svelte";
  import FavoriteButton from "$lib/components/FavoriteButton.svelte";
  import QuickAddButton from "$lib/components/QuickAddButton.svelte";
  import { Button } from "$lib/components/ui";

  const musicManager = getMusicManager();

  const [send, receive] = crossfade({ duration: 250 });

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

  let open = $state(false);
  let closing = $state(false);

  let seekValue = $derived(
    Number.isNaN(musicManager.duration)
      ? 0
      : musicManager.currentTime / musicManager.duration,
  );

  function openPlayer() {
    open = true;
    closing = false;
  }

  function closePlayer() {
    closing = true;
    setTimeout(() => {
      open = false;
      closing = false;
    }, 200);
  }

  $effect(() => {
    if (open) {
      document.body.style.overflow = "hidden";
    } else {
      document.body.style.overflow = "";
    }
    return () => {
      document.body.style.overflow = "";
    };
  });

  // Swipe gestures: horizontal on cover switches tracks, vertical closes player
  let coverEl: HTMLDivElement | undefined = $state();
  let dragging = $state(false);
  let dragMode: "h" | "v" | null = $state(null);
  let dragX = $state(0);
  let dragY = $state(0);
  let dragStartX = 0;
  let dragStartY = 0;
  let canSwipeH = false;

  function onGesturePointerDown(e: PointerEvent) {
    if (e.pointerType === "mouse" && e.button !== 0) return;

    const target = e.target as HTMLElement;
    if (target.closest("button, a, [role='slider'], [data-no-drag]")) return;

    dragStartX = e.clientX;
    dragStartY = e.clientY;
    dragX = 0;
    dragY = 0;
    dragMode = null;
    canSwipeH = coverEl?.contains(target) ?? false;
    dragging = true;

    (e.currentTarget as HTMLElement).setPointerCapture(e.pointerId);
  }

  function onGesturePointerMove(e: PointerEvent) {
    if (!dragging) return;

    const dx = e.clientX - dragStartX;
    const dy = e.clientY - dragStartY;

    if (!dragMode) {
      if (Math.abs(dx) > 8 || Math.abs(dy) > 8) {
        dragMode =
          Math.abs(dx) > Math.abs(dy) ? "h" : "v";
      } else {
        return;
      }
    }

    if (dragMode === "h" && canSwipeH) {
      dragX = dx;
    } else if (dragMode === "v") {
      dragY = Math.max(0, dy);
    }
  }

  function onGesturePointerEnd() {
    if (!dragging) return;

    if (dragMode === "h" && canSwipeH) {
      if (dragX < -60) {
        musicManager.nextTrack();
      } else if (dragX > 60) {
        musicManager.previousTrack();
      }
    } else if (dragMode === "v" && dragY > 110) {
      closePlayer();
    }

    dragX = 0;
    dragY = 0;
    dragMode = null;
    dragging = false;
  }

  let overlayOpacity = $derived.by(() => {
    if (dragMode === "v") return `${Math.max(0, 1 - dragY / 240)}`;
    if (closing) return "0";
    return null;
  });

  // Drag-up queue sheet
  let queueOpen = $state(false);
  let queueDragY = $state(0);
  let queueDragging = $state(false);
  let queueStartY = 0;

  function openQueue() {
    queueOpen = true;
  }

  function closeQueue() {
    queueOpen = false;
    queueDragY = 0;
  }

  function onQueueHandleDown(e: PointerEvent) {
    queueStartY = e.clientY;
    queueDragY = 0;
    queueDragging = true;
    (e.currentTarget as HTMLElement).setPointerCapture(e.pointerId);
  }

  function onQueueHandleMove(e: PointerEvent) {
    if (!queueDragging) return;
    queueDragY = Math.max(0, e.clientY - queueStartY);
  }

  function onQueueHandleEnd() {
    if (!queueDragging) return;
    if (queueDragY > 110) closeQueue();
    queueDragY = 0;
    queueDragging = false;
  }
</script>

<!-- Bottom mini bar -->
{#if !open}
  <div
    class="fixed bottom-14 z-[60] w-full md:hidden"
    in:fly={{ y: 40, duration: 220, opacity: 0 }}
    out:fly={{ y: 40, duration: 220, opacity: 0 }}
  >
    <div class="w-full border-t bg-background">
      <SeekSlider
        class="h-3"
        growOnHover={false}
        value={seekValue}
        onValue={(p) => {
          musicManager.setPosition(p * musicManager.duration);
        }}
        buffered={musicManager.buffered}
        ariaLabel="Seek"
      />
      <div
        class="flex w-full items-center gap-3 px-3 pb-2 text-foreground"
        role="button"
        tabindex="0"
        onclick={openPlayer}
        onkeydown={(e) => {
          if (e.key === "Enter" || e.key === " ") {
            e.preventDefault();
            openPlayer();
          }
        }}
      >
        <div in:receive={{ key: "cover" }} out:send={{ key: "cover" }}>
          <Image
            class="w-10 min-w-10 shrink-0"
            src={currentMediaItem?.coverArt}
            alt="cover"
            loading="eager"
          />
        </div>

        <div class="flex min-w-0 flex-col items-start text-left">
          <p class="w-full truncate text-sm font-medium">
            {currentMediaItem?.name ?? "No track playing"}
          </p>
          <p class="w-full truncate text-xs text-muted-foreground">
            {currentMediaItem?.artists[0]?.name ?? ""}
          </p>
        </div>

        <div class="flex-grow"></div>

        {#if musicManager.loading}
          <div class="h-8 w-8 animate-pulse rounded-full bg-muted"></div>
        {:else if musicManager.playing}
          <button
            class="flex h-8 w-8 items-center justify-center"
            onclick={(e) => {
              e.stopPropagation();
              musicManager.pause();
            }}
            aria-label="Pause"
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
            aria-label="Play"
          >
            <Play size="22" />
          </button>
        {/if}
      </div>
    </div>
  </div>
{/if}

<!-- Full-screen player -->
{#if open}
  <div
    class="fixed inset-0 z-50 flex flex-col bg-zinc-950 transition-opacity duration-200"
    class:transition-transform={!dragging}
    style:opacity={overlayOpacity}
    style:transform={dragMode === "v" ? `translateY(${dragY}px)` : null}
    in:fade={{ duration: 220 }}
    role="dialog"
    aria-label="Now Playing"
  >
    <!-- Blurred cover backdrop -->
    <div class="pointer-events-none absolute inset-0 overflow-hidden">
      {#if currentMediaItem}
        <Image
          class="absolute inset-0 h-full w-full scale-125 object-cover blur-2xl"
          src={currentMediaItem.coverArt}
          alt=""
          loading="eager"
        />
        <div
          class="absolute inset-0 bg-linear-to-b from-black/70 via-black/40 to-background"
        ></div>
      {/if}
    </div>

    <div class="relative flex min-h-0 flex-1 flex-col">
      <!-- Header -->
      <div class="flex items-center justify-between px-4 pb-2 pt-3" data-no-drag>
        <button
          class="flex items-center gap-1 text-white/70 transition-colors hover:text-white"
          onclick={closePlayer}
          aria-label="Close"
        >
          <ChevronDown size="24" />
        </button>
        <p class="truncate text-sm font-medium uppercase tracking-widest text-white/70">
          Now Playing
        </p>
        {#if currentMediaItem}
          <div class="flex items-center gap-0.5 text-white">
            <FavoriteButton show trackId={currentMediaItem.trackId} />
            <QuickAddButton trackId={currentMediaItem.trackId} />
          </div>
        {:else}
          <div class="w-6"></div>
        {/if}
      </div>

      <!-- Gesture surface -->
      <div
        class="flex min-h-0 flex-1 flex-col"
        role="presentation"
        onpointerdown={onGesturePointerDown}
        onpointermove={onGesturePointerMove}
        onpointerup={onGesturePointerEnd}
        onpointercancel={onGesturePointerEnd}
      >
        <div
          class="flex-1 overflow-y-auto overscroll-contain px-6"
        >
          <div class="flex min-h-full flex-col justify-center gap-4 py-4">
            <div class="flex w-full justify-center">
              <div
                class:transition-transform={!dragging}
                style:transform={dragMode === "h"
                  ? `translateX(${dragX}px) rotate(${(dragX * 0.04).toFixed(2)}deg)`
                  : null}
              >
                <div in:receive={{ key: "cover" }} out:send={{ key: "cover" }}>
                  <div bind:this={coverEl}>
                    <Image
                      class="w-72 max-w-full shadow-2xl"
                      src={currentMediaItem?.coverArt}
                      alt="Track Cover Art"
                      loading="eager"
                    />
                  </div>
                </div>
              </div>
            </div>

            <!-- Track info -->
            <div class="flex w-full flex-col items-center gap-1 text-center">
              <p class="text-2xl font-semibold text-white">
                {currentMediaItem?.name ?? "Nothing playing"}
              </p>
              <p class="text-sm text-white/60">
                {#if currentMediaItem}
                  {#each currentMediaItem.artists as artist, i (artist.id)}
                    {#if i > 0},
                    {/if}
                    <a
                      href="/artists/{artist.id}"
                      class="font-medium text-white/80 hover:underline"
                      >{artist.name}</a
                    >
                  {/each}
                {/if}
              </p>
            </div>

            <!-- Seek bar -->
            <div
              class="flex w-full flex-col gap-0.5"
              data-no-drag
            >
              <SeekSlider
                value={seekValue}
                onValue={(p) => {
                  musicManager.setPosition(p * musicManager.duration);
                }}
                buffered={musicManager.buffered}
              />
              <div
                class="flex justify-between text-xs tabular-nums text-white/60"
              >
                <span>{formatTime(musicManager.currentTime)}</span>
                <span
                  >{formatTime(
                    Number.isNaN(musicManager.duration)
                      ? 0
                      : musicManager.duration,
                  )}</span
                >
              </div>
            </div>

            <!-- Controls -->
            <div
              class="flex w-full items-center justify-center gap-8 py-1"
              data-no-drag
            >
              <button
                class="text-white/80 transition-colors hover:text-white"
                onclick={() => musicManager.previousTrack()}
                aria-label="Previous"
              >
                <SkipBack size="30" />
              </button>

              {#if musicManager.loading}
                <div class="h-16 w-16 animate-pulse rounded-full bg-white/15"></div>
              {:else if musicManager.playing}
                <button
                  class="flex h-16 w-16 items-center justify-center rounded-full bg-white text-black shadow-lg transition-transform hover:scale-105"
                  onclick={() => musicManager.pause()}
                  aria-label="Pause"
                >
                  <Pause size="30" />
                </button>
              {:else}
                <button
                  class="flex h-16 w-16 items-center justify-center rounded-full bg-white text-black shadow-lg transition-transform hover:scale-105"
                  onclick={() => musicManager.play()}
                  aria-label="Play"
                >
                  <Play size="30" />
                </button>
              {/if}

              <button
                class="text-white/80 transition-colors hover:text-white"
                onclick={() => musicManager.nextTrack()}
                aria-label="Next"
              >
                <SkipForward size="30" />
              </button>
            </div>

            <!-- Volume -->
            <div
              class="mx-auto flex w-full max-w-64 items-center gap-2"
              data-no-drag
            >
              <button
                class="shrink-0 text-white/70 transition-colors hover:text-white"
                onclick={() => {
                  musicManager.muted = !musicManager.muted;
                }}
                aria-label="Mute"
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
                  ariaLabel="Volume"
                />
              </div>
            </div>
          </div>
        </div>

        <!-- Queue launcher -->
        <div class="flex items-center justify-center pb-3" data-no-drag>
          <button
            class="flex items-center gap-1.5 rounded-full bg-white/10 px-4 py-1.5 text-sm font-medium text-white backdrop-blur-sm transition-colors hover:bg-white/20"
            onclick={openQueue}
            aria-label="Open queue"
          >
            <ListMusic size="16" />
            Queue
          </button>
        </div>
      </div>
    </div>

    <!-- Queue sheet -->
    {#if queueOpen}
      <button
        class="absolute inset-0 z-[5] bg-black/50"
        in:fade={{ duration: 150 }}
        out:fade={{ duration: 150 }}
        onclick={closeQueue}
        aria-label="Close queue"
        data-no-drag
      ></button>
      <div
        class="absolute inset-x-0 bottom-0 z-[6] flex max-h-[75vh] flex-col rounded-t-2xl bg-background text-foreground shadow-2xl"
        class:transition-transform={!queueDragging}
        style:transform={queueDragY > 0 ? `translateY(${queueDragY}px)` : null}
        in:slide={{ duration: 220 }}
        out:slide={{ duration: 150 }}
        data-no-drag
      >
        <div
          class="flex cursor-grab touch-none flex-col items-center py-2 active:cursor-grabbing"
          role="presentation"
          onpointerdown={onQueueHandleDown}
          onpointermove={onQueueHandleMove}
          onpointerup={onQueueHandleEnd}
          onpointercancel={onQueueHandleEnd}
        >
          <div class="h-1.5 w-10 rounded-full bg-muted"></div>
        </div>

        <div class="flex items-center justify-between px-4 pb-2">
          <p class="text-base font-semibold">Queue</p>
          <div class="flex items-center gap-1">
            <Button
              variant="ghost"
              size="sm"
              href="/queue"
              onclick={() => closePlayer()}
            >
              View all
            </Button>
            <Button
              class="rounded-full"
              variant="ghost"
              size="icon"
              onclick={async () => {
                await musicManager.clearQueue();
              }}
              aria-label="Clear queue"
            >
              <ListX />
            </Button>
          </div>
        </div>

        <div class="max-h-[45vh] overflow-y-auto overscroll-contain px-4 pb-6">
          <div class="flex flex-col gap-3">
            <!-- Played -->
            {#if previousItems.length > 0}
              <div>
                <p
                  class="mb-1.5 text-xs font-semibold uppercase tracking-wider text-muted-foreground"
                >
                  Played
                </p>
                <div class="flex flex-col gap-0.5">
                  {#each previousItems as mediaItem, i (mediaItem.trackId)}
                    {@const queueIndex = previousStartIndex + i}
                    <div
                      class="group flex items-center gap-3 rounded-md p-2 transition-colors hover:bg-accent/50"
                    >
                      <button
                        class="flex min-w-0 flex-1 items-center gap-3 rounded-md text-left"
                        onclick={async () => {
                          await musicManager.setQueueIndex(queueIndex);
                          musicManager.play();
                        }}
                        aria-label={`Play ${mediaItem.name}`}
                      >
                        <Image
                          class="w-10 min-w-10 shrink-0 rounded"
                          src={mediaItem.coverArt}
                          alt="cover"
                        />
                        <div class="flex min-w-0 flex-1 flex-col">
                          <p class="truncate text-sm font-medium">
                            {mediaItem.name}
                          </p>
                          <p class="truncate text-xs text-muted-foreground">
                            {#each mediaItem.artists as artist, j (
                              artist.id
                            )}
                              {#if j > 0},
                              {/if}{artist.name}
                            {/each}
                          </p>
                        </div>
                      </button>
                    </div>
                  {/each}
                </div>
              </div>
            {/if}

            <!-- Now Playing -->
            {#if currentQueueItem}
              <div>
                <p
                  class="mb-1.5 text-xs font-semibold uppercase tracking-wider text-muted-foreground"
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
                      <ListMusic class="text-primary" size="16" />
                    </div>
                  </div>
                  <div class="flex min-w-0 flex-col">
                    <p class="truncate text-sm font-medium">
                      {currentQueueItem.name}
                    </p>
                    <p class="truncate text-xs text-muted-foreground">
                      {#each currentQueueItem.artists as artist, j (
                        artist.id
                      )}
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
                  class="mb-1.5 text-xs font-semibold uppercase tracking-wider text-muted-foreground"
                >
                  Up Next
                </p>
                <div class="flex flex-col gap-0.5">
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
                        class="flex min-w-0 flex-1 items-center gap-3 rounded-md text-left"
                        onclick={async () => {
                          await musicManager.setQueueIndex(queueIndex);
                          musicManager.play();
                        }}
                        aria-label={`Play ${mediaItem.name}`}
                      >
                        <Image
                          class="w-10 min-w-10 shrink-0 rounded"
                          src={mediaItem.coverArt}
                          alt="cover"
                        />
                        <div class="flex min-w-0 flex-1 flex-col">
                          <p class="truncate text-sm font-medium">
                            {mediaItem.name}
                          </p>
                          <p class="truncate text-xs text-muted-foreground">
                            {#each mediaItem.artists as artist, j (
                              artist.id
                            )}
                              {#if j > 0},
                              {/if}{artist.name}
                            {/each}
                          </p>
                        </div>
                      </button>
                      <button
                        class="shrink-0 rounded-full p-1 text-muted-foreground opacity-0 transition-opacity hover:text-foreground group-hover:opacity-100"
                        onclick={async () => {
                          await musicManager.removeQueueItem(queueIndex);
                        }}
                        aria-label="Remove from queue"
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
        </div>
      </div>
    {/if}
  </div>
{/if}