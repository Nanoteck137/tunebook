<script lang="ts">
  import Slider from "$lib/components/SeekSlider.svelte";
  import { Button, ScrollArea, Separator, Sheet } from "@nanoteck137/nano-ui";
  import { formatTime } from "$lib/utils";
  import {
    ListX,
    Logs,
    Pause,
    Play,
    SkipBack,
    SkipForward,
    TestTube,
    Volume2,
    VolumeX,
  } from "lucide-svelte";
  import { getMusicManager } from "$lib/music-manager.svelte";
  import Spinner from "$lib/components/Spinner.svelte";
  import type { MediaItem } from "$lib/api/types";
  import Image from "$lib/components/Image.svelte";
  import { getApiClient, handleApiError } from "$lib";
  import toast from "svelte-5-french-toast";

  const musicManager = getMusicManager();
  const apiClient = getApiClient();

  let currentMediaItem = $state<MediaItem | null>(null);

  $effect(() => {
    currentMediaItem = musicManager.queue.getCurrentMediaItem();
  });
</script>

{#snippet queueSheet()}
  <Sheet.Root>
    <Sheet.Trigger>
      <Logs size="24" />
    </Sheet.Trigger>
    <Sheet.Content side="right">
      <div class="flex items-center gap-2 pb-2">
        <p>Queue</p>
        <Button
          class="rounded-full"
          variant="ghost"
          size="icon"
          onclick={() => {
            musicManager.clearQueue();
          }}
        >
          <ListX />
        </Button>
      </div>

      <ScrollArea class="h-full pb-6">
        <div class="mr-3 flex flex-col gap-2">
          {#each musicManager.queue.items as mediaItem, i}
            <div
              class={`flex items-center gap-2 rounded p-1 ${musicManager.queue.index === i ? "bg-accent text-accent-foreground" : ""}`}
            >
              <div class="group relative">
                <Image
                  class="w-12 min-w-12"
                  src={mediaItem?.coverArt.small}
                  alt="cover"
                />
                {#if i == musicManager.queue.index}
                  <div
                    class="absolute bottom-0 left-0 right-0 top-0 flex items-center justify-center rounded border bg-muted/80"
                  >
                    <Play class="text-muted-foreground" size="25" />
                  </div>
                {:else}
                  <button
                    class={`absolute bottom-0 left-0 right-0 top-0 hidden items-center justify-center rounded border bg-muted/80 group-hover:flex`}
                    onclick={async () => {
                      await musicManager.setQueueIndex(i);
                      musicManager.play();
                    }}
                  >
                    <Play class="text-muted-foreground" size="25" />
                  </button>
                {/if}
              </div>
              <div class="flex flex-col">
                <p class="line-clamp-1 text-sm" title={mediaItem?.name}>
                  {mediaItem?.name}
                </p>
                <p
                  class="line-clamp-1 text-xs font-light"
                  title={mediaItem?.artists[0].name}
                >
                  {mediaItem?.artists[0].name}
                </p>
              </div>
            </div>
            <Separator />
          {/each}
        </div>
      </ScrollArea>
    </Sheet.Content>
  </Sheet.Root>
{/snippet}

<div
  class="container z-30 hidden h-16 bg-background text-foreground transition-transform duration-500 md:block"
>
  <div class="absolute -top-1.5 left-0 right-0">
    <Slider
      value={musicManager.currentTime / musicManager.duration}
      onValue={(p) => {
        musicManager.setPosition(p * musicManager.duration);
        // onSeek(p * musicManager.duration);
      }}
    />
  </div>

  <div class="grid h-full grid-cols-footer">
    <div class="flex items-center gap-2">
      <div class="flex items-center gap-2">
        <!-- <button
          onclick={async () => {
            if (!mediaItem) return;

            const res = await apiClient.recordTrack(mediaItem.track.id, {
              source: "unknown",
              duration: currentTime,
            });
            if (!res.success) {
              return handleApiError(res.error);
            }

            toast.success("Recorded track");
          }}
        >
          <TestTube size={32} />
        </button> -->

        <button
          onclick={() => {
            musicManager.previousTrack();
          }}
        >
          <SkipBack size={32} />
        </button>

        {#if musicManager.loading}
          <Spinner class="h-8 w-8" />
        {:else if musicManager.playing}
          <button
            onclick={() => {
              musicManager.pause();
            }}
          >
            <Pause size={32} />
          </button>
        {:else}
          <button
            onclick={() => {
              musicManager.play();
            }}
          >
            <Play size={32} />
          </button>
        {/if}

        <button
          onclick={() => {
            musicManager.nextTrack();
          }}
        >
          <SkipForward size={32} />
        </button>
      </div>

      <p class="hidden min-w-20 text-xs font-medium lg:block">
        {formatTime(musicManager.currentTime)} /{" "}
        {formatTime(
          Number.isNaN(musicManager.duration) ? 0 : musicManager.duration,
        )}
      </p>

      <div class="flex items-center justify-center gap-2 align-middle">
        <Image
          class="w-12 min-w-12"
          src={currentMediaItem?.coverArt.small}
          alt="cover"
          loading="eager"
        />
        <div class="flex flex-col">
          <p
            class="line-clamp-1 text-ellipsis text-sm"
            title={currentMediaItem?.name}
          >
            {currentMediaItem?.name}
          </p>

          <p class="line-clamp-1 min-w-80 text-ellipsis text-xs">
            {currentMediaItem?.artists[0].name}
          </p>
        </div>
      </div>
    </div>

    <div class="flex items-center justify-evenly">
      <div class="flex w-full items-center gap-4 p-4">
        <Slider
          value={musicManager.volume}
          onValue={(p) => {
            // musicManager.set
            // onVolumeChanged(p);
          }}
        />
        <button
          onclick={() => {
            // onToggleMuted();
          }}
        >
          {#if musicManager.muted}
            <VolumeX size="25" />
          {:else}
            <Volume2 size="25" />
          {/if}
        </button>

        {@render queueSheet()}
      </div>
    </div>
  </div>
</div>
