<script lang="ts">
  import {
    ChevronDown,
    Pause,
    Play,
    SkipBack,
    SkipForward,
    Volume2,
    VolumeX,
    X,
  } from "lucide-svelte";
  import { crossfade, fade } from "svelte/transition";
  import { formatTime } from "$lib/utils";
  import SeekSlider from "$lib/components/SeekSlider.svelte";
  import { getMusicManager, type MediaItem } from "$lib/music-manager.svelte";
  import Image from "$lib/components/Image.svelte";

  const musicManager = getMusicManager();

  const [send, receive] = crossfade({ duration: 250 });

  let currentMediaItem = $state<MediaItem | null>(null);
  let nextItems = $state<MediaItem[]>([]);

  $effect(() => {
    currentMediaItem = musicManager.currentItem;

    musicManager.queue.getNextItems(0, 50).then((items) => {
      nextItems = items;
    });
  });

  let open = $state(false);
  let closing = $state(false);
  let activeTab = $state<string | null>(null);

  let seekValue = $derived(
    Number.isNaN(musicManager.duration)
      ? 0
      : musicManager.currentTime / musicManager.duration,
  );

  function openPlayer() {
    open = true;
    closing = false;
    activeTab = null;
  }

  function closePlayer() {
    closing = true;
    setTimeout(() => {
      open = false;
      closing = false;
    }, 200);
  }

  function switchTab(tab: string) {
    activeTab = activeTab === tab ? null : tab;
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
</script>

<!-- Bottom mini bar -->
{#if !open}
  <div class="fixed bottom-14 z-[60] w-full md:hidden">
    <button
      class="flex w-full items-center gap-3 border-t bg-background px-3 py-2 text-foreground"
      onclick={openPlayer}
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
    </button>
  </div>
{/if}

<!-- Full-screen player -->
{#if open}
  <div
    class="fixed inset-0 z-50 flex flex-col bg-gradient-to-b from-zinc-950 to-background transition-opacity duration-200"
    class:opacity-100={!closing}
    class:opacity-0={closing}
    role="dialog"
    aria-label="Now Playing"
  >
    <!-- Header -->
    <div class="flex items-center gap-4 px-4 pb-2 pt-2">
      <button
        class="flex items-center gap-1 text-muted-foreground transition-colors hover:text-foreground"
        onclick={closePlayer}
      >
        <ChevronDown size="24" />
      </button>
      <p
        class="truncate text-center text-sm font-medium text-muted-foreground"
      >
        Now Playing
      </p>
    </div>

    <!-- Content area -->
    <div class="flex-1 overflow-y-auto overscroll-contain px-4 pb-2">
      {#if !activeTab}
        <div in:fade={{ duration: 150 }} out:fade={{ duration: 150 }}>
          <div class="flex flex-col items-center gap-3 pb-4">
            <!-- Cover art -->
            <div class="flex w-full justify-center pt-2">
              <div in:receive={{ key: "cover" }} out:send={{ key: "cover" }}>
                <Image
                  class="w-72 max-w-full shadow-xl"
                  src={currentMediaItem?.coverArt}
                  alt="Track Cover Art"
                  loading="eager"
                />
              </div>
            </div>

            <!-- Track info -->
            <div class="flex w-full flex-col items-center text-center">
              <p class="text-xl font-medium">{currentMediaItem?.name}</p>
              <p class="text-sm text-muted-foreground">
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

            <!-- Seek bar -->
            <div class="flex w-full flex-col gap-1">
              <SeekSlider
                value={seekValue}
                onValue={(p) => {
                  musicManager.setPosition(p * musicManager.duration);
                }}
                buffered={musicManager.buffered}
              />
              <div class="flex justify-between text-xs text-muted-foreground">
                <span class="tabular-nums"
                  >{formatTime(musicManager.currentTime)}</span
                >
                <span class="tabular-nums"
                  >{formatTime(
                    Number.isNaN(musicManager.duration)
                      ? 0
                      : musicManager.duration,
                  )}</span
                >
              </div>
            </div>

            <!-- Controls -->
            <div class="flex w-full items-center justify-center gap-8 py-1">
              <button
                class="text-muted-foreground transition-colors hover:text-foreground"
                onclick={() => musicManager.previousTrack()}
              >
                <SkipBack size="26" />
              </button>

              {#if musicManager.loading}
                <div class="h-16 w-16 animate-pulse rounded-full bg-muted" />
              {:else if musicManager.playing}
                <button
                  class="flex h-16 w-16 items-center justify-center rounded-full bg-foreground text-background shadow-lg transition-colors hover:scale-105"
                  onclick={() => musicManager.pause()}
                >
                  <Pause size="30" />
                </button>
              {:else}
                <button
                  class="flex h-16 w-16 items-center justify-center rounded-full bg-foreground text-background shadow-lg transition-colors hover:scale-105"
                  onclick={() => musicManager.play()}
                >
                  <Play size="30" />
                </button>
              {/if}

              <button
                class="text-muted-foreground transition-colors hover:text-foreground"
                onclick={() => musicManager.nextTrack()}
              >
                <SkipForward size="26" />
              </button>
            </div>

            <!-- Volume -->
            <div class="flex w-full max-w-56 items-center gap-2">
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
        </div>
      {:else}
        <div in:fade={{ duration: 150 }} out:fade={{ duration: 150 }}>
          <!-- Mini player at top -->
          <div class="flex items-center gap-3 py-2">
            <div>
              <Image
                class="w-12 min-w-12 shrink-0"
                src={currentMediaItem?.coverArt}
                alt="cover"
                loading="eager"
              />
            </div>
            <div class="flex min-w-0 flex-1 flex-col">
              <p class="truncate text-sm font-medium">
                {currentMediaItem?.name ?? "No track playing"}
              </p>
              <p class="truncate text-xs text-muted-foreground">
                {currentMediaItem?.artists[0]?.name ?? ""}
              </p>
            </div>
            {#if musicManager.loading}
              <div class="h-10 w-10 animate-pulse rounded-full bg-muted" />
            {:else if musicManager.playing}
              <button
                class="flex h-10 w-10 items-center justify-center rounded-full bg-foreground text-background"
                onclick={() => musicManager.pause()}
              >
                <Pause size="20" />
              </button>
            {:else}
              <button
                class="flex h-10 w-10 items-center justify-center rounded-full bg-foreground text-background"
                onclick={() => musicManager.play()}
              >
                <Play size="20" />
              </button>
            {/if}
          </div>

          {#if activeTab === "queue"}
            <div class="flex items-center justify-between px-2 py-2">
              <p class="text-sm font-medium text-muted-foreground">Up next</p>
              <a
                href="/queue"
                class="text-sm font-medium text-primary transition-colors hover:underline"
                onclick={() => closePlayer()}
              >
                View full queue
              </a>
            </div>

            {#if nextItems.length > 0}
              <div class="flex w-full flex-col gap-1 pb-4">
                {#each nextItems as mediaItem, i (mediaItem.trackId)}
                  {@const queueIndex = musicManager.queue.index + 1 + i}
                  <div
                    class="group flex items-center gap-3 rounded-md p-2 transition-colors hover:bg-accent/50"
                  >
                    <span
                      class="w-6 text-right text-xs tabular-nums text-muted-foreground"
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
                      <p class="truncate text-sm font-medium">
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
            {:else}
              <p class="py-8 text-center text-sm text-muted-foreground">
                Queue is empty
              </p>
            {/if}
          {:else if activeTab === "lyrics"}
            <div class="flex items-center justify-center py-16">
              <p class="text-sm text-muted-foreground">Lyrics not available</p>
            </div>
          {:else if activeTab === "related"}
            <div class="flex items-center justify-center py-16">
              <p class="text-sm text-muted-foreground">
                Related content not available
              </p>
            </div>
          {/if}
        </div>
      {/if}
    </div>

    <!-- Tab bar -->
    <div class="flex border-t border-border px-4 pt-2">
      <button
        class="flex-1 pb-2 text-sm font-medium transition-colors {activeTab ===
        'queue'
          ? 'border-b-2 border-foreground text-foreground'
          : 'text-muted-foreground'}"
        onclick={() => switchTab("queue")}
      >
        Queue
      </button>
      <button
        class="flex-1 pb-2 text-sm font-medium transition-colors {activeTab ===
        'lyrics'
          ? 'border-b-2 border-foreground text-foreground'
          : 'text-muted-foreground'}"
        onclick={() => switchTab("lyrics")}
      >
        Lyrics
      </button>
      <button
        class="flex-1 pb-2 text-sm font-medium transition-colors {activeTab ===
        'related'
          ? 'border-b-2 border-foreground text-foreground'
          : 'text-muted-foreground'}"
        onclick={() => switchTab("related")}
      >
        Related
      </button>
    </div>
  </div>
{/if}
