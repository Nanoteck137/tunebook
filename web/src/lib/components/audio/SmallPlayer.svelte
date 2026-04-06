<script lang="ts">
  import {
    ChevronUp,
    Pause,
    Play,
    SkipBack,
    SkipForward,
  } from "lucide-svelte";
  import { formatTime } from "$lib/utils";
  import { ScrollArea, Sheet, buttonVariants } from "@nanoteck137/nano-ui";
  import SeekSlider from "$lib/components/SeekSlider.svelte";
  import { fly } from "svelte/transition";
  import { getMusicManager, type MediaItem } from "$lib/music-manager.svelte";
  import Image from "$lib/components/Image.svelte";

  const musicManager = getMusicManager();

  let currentMediaItem = $state<MediaItem | null>(null);

  $effect(() => {
    currentMediaItem = musicManager.queue.getCurrentMediaItem();
  });
</script>

{#snippet queueSheet()}
  <Sheet.Root>
    <Sheet.Trigger class={buttonVariants({ variant: "outline" })}>
      Queue
    </Sheet.Trigger>
    <Sheet.Content side="bottom">
      <p class="pb-2">Queue</p>
      <ScrollArea class="h-96">
        <div class="flex flex-col gap-2">
          {#each musicManager.queue.items as mediaItem, i}
            <div class="flex items-center gap-2">
              <div class="group relative">
                <Image
                  class="w-12 min-w-12"
                  src={mediaItem.coverArt}
                  alt="cover"
                />
                {#if i == musicManager.queue.index}
                  <div
                    class="absolute bottom-0 left-0 right-0 top-0 flex items-center justify-center border bg-black/80"
                  >
                    <Play size="25" />
                  </div>
                {:else}
                  <button
                    class={`absolute bottom-0 left-0 right-0 top-0 hidden items-center justify-center border bg-black/80 group-hover:flex`}
                    onclick={() => {
                      musicManager.setQueueIndex(i);
                      musicManager.play();
                    }}
                  >
                    <Play size="25" />
                  </button>
                {/if}
              </div>
              <div class="flex flex-col">
                <p class="line-clamp-1 text-sm" title={mediaItem.name}>
                  {mediaItem.name}
                </p>
                <p
                  class="line-clamp-1 text-xs"
                  title={mediaItem.artists[0].name}
                >
                  {mediaItem.artists[0].name}
                </p>
              </div>
            </div>
          {/each}
        </div>
      </ScrollArea>
    </Sheet.Content>
  </Sheet.Root>
{/snippet}

<div
  class="z-30 h-16 border-t bg-background text-foreground md:hidden"
  transition:fly={{ y: 200 }}
>
  <div class="flex items-center">
    {#if musicManager.playing}
      <button class="p-4" onclick={() => musicManager.pause()}>
        <Pause size="24" />
      </button>
    {:else}
      <button class="p-4" onclick={() => musicManager.play()}>
        <Play size="24" />
      </button>
    {/if}

    <Sheet.Root>
      <Sheet.Trigger class="flex grow items-center">
        <Image
          class="w-12 min-w-12"
          src={currentMediaItem?.coverArt}
          alt="cover"
          loading="eager"
        />

        <div class="flex flex-col items-start justify-center px-2">
          <p class="line-clamp-1 text-sm">{currentMediaItem?.name}</p>
          <p class="line-clamp-1 text-xs">
            {currentMediaItem?.artists[0].name}
          </p>
        </div>

        <div class="flex-grow"></div>
        <div class="flex h-16 min-w-16 items-center justify-center">
          <ChevronUp size="30" />
        </div>
      </Sheet.Trigger>
      <Sheet.Content side="bottom">
        <div class="relative flex flex-col items-center justify-center gap-2">
          {@render queueSheet()}

          <Image
            class="w-64"
            src={currentMediaItem?.coverArt}
            alt="Track Cover Art"
            loading="eager"
          />

          <div class="flex flex-col items-center">
            <p class="font-medium">{currentMediaItem?.name}</p>
            <p class="text-xs">{currentMediaItem?.artists[0].name}</p>
          </div>

          <div class="flex w-full flex-col gap-1 px-4 py-2">
            <SeekSlider
              value={musicManager.currentTime / musicManager.duration}
              onValue={(p) => {
                musicManager.setPosition(p * musicManager.duration);
              }}
            />

            <div class="flex justify-between">
              <p class="text-sm">
                {formatTime(musicManager.currentTime)}
              </p>

              <p class="text-sm">
                {formatTime(
                  Number.isNaN(musicManager.duration)
                    ? 0
                    : musicManager.duration,
                )}
              </p>
            </div>
          </div>

          <div class="flex w-full items-center gap-4 px-4">
            <div class="flex gap-4">
              <button
                onclick={() => {
                  musicManager.previousTrack();
                }}
              >
                <SkipBack size="38" />
              </button>

              {#if musicManager.loading}
                <p>Loading...</p>
              {:else if musicManager.playing}
                <button
                  onclick={() => {
                    musicManager.pause();
                  }}
                >
                  <Pause size={46} />
                </button>
              {:else}
                <button
                  onclick={() => {
                    musicManager.play();
                  }}
                >
                  <Play size={46} />
                </button>
              {/if}

              <button
                onclick={() => {
                  musicManager.nextTrack();
                }}
              >
                <SkipForward size="38" />
              </button>
            </div>

            <div class="flex-grow"></div>

            <div class="flex w-full max-w-56 items-center gap-4">
              <!-- <Slider
                bind:value={vol}
                onValueChange={(e) => onVolumeChanged(e[0] / 100)}
              />
              <button
                onclick={() => {
                  onToggleMuted();
                }}
              >
                {#if audioMuted}
                  <VolumeX size="30" />
                {:else}
                  <Volume2 size="30" />
                {/if}
              </button> -->
            </div>
          </div>
        </div>
      </Sheet.Content>
    </Sheet.Root>
  </div>
</div>
