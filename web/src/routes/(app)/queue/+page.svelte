<script lang="ts">
  import { getMusicManager } from "$lib/music-manager.svelte";
  import Image from "$lib/components/Image.svelte";
  import { handleApiError } from "$lib";
  import { Play } from "lucide-svelte";
  import { onMount } from "svelte";

  let { data } = $props();

  const musicManager = getMusicManager();
  const deviceId = localStorage.getItem("device-id") ?? "";

  type QueueListItem = {
    queueItemId: string;
    trackId: string;
    name: string;
    albumName: string;
    artistNames: string;
    coverArt: string;
    duration: number;
  };

  let items = $state<QueueListItem[]>([]);
  let currentIndex = $state(0);
  let totalItems = $state(0);
  let loading = $state(false);
  let hasMore = $state(true);
  let page = $state(0);
  let sentinel: HTMLDivElement | null = $state(null);

  $effect(() => {
    currentIndex = musicManager.queue.index;
  });

  async function loadMore() {
    if (loading || !hasMore) return;
    loading = true;

    const res = await data.apiClient.getQueue(deviceId, {
      query: { page: page.toString() },
    });

    if (!res.success) {
      handleApiError(res.error);
      loading = false;
      return;
    }

    const newItems: QueueListItem[] = res.data.items.map((item) => ({
      queueItemId: item.queueItemId,
      trackId: item.track.id,
      name: item.track.name,
      albumName: item.track.albumName,
      artistNames: item.track.artists.map((a) => a.name).join(", "),
      coverArt: item.track.coverArt.small,
      duration: item.track.duration,
    }));

    items = [...items, ...newItems];
    currentIndex = res.data.currentIndex;
    totalItems = res.data.page.totalItems;
    page++;
    hasMore = items.length < totalItems;
    loading = false;
  }

  function formatDuration(seconds: number): string {
    const m = Math.floor(seconds / 60);
    const s = Math.floor(seconds % 60);
    return `${m}:${s.toString().padStart(2, "0")}`;
  }

  onMount(() => {
    loadMore();

    if (!sentinel) return;

    const observer = new IntersectionObserver(
      (entries) => {
        if (entries[0].isIntersecting) {
          loadMore();
        }
      },
      { rootMargin: "200px" },
    );

    observer.observe(sentinel);

    return () => observer.disconnect();
  });
</script>

<svelte:head>
  <title>Queue - Tunebook</title>
</svelte:head>

<div class="container mx-auto px-4 py-6">
  <div class="mb-6">
    <h1 class="text-3xl font-bold">Queue</h1>
    <p class="text-sm text-muted-foreground">
      {#if totalItems > 0}
        Track {currentIndex + 1} of {totalItems}
      {:else}
        No tracks in queue
      {/if}
    </p>
  </div>

  {#if items.length === 0 && !loading}
    <p class="py-12 text-center text-muted-foreground">The queue is empty.</p>
  {:else}
    <div class="flex flex-col gap-1">
      {#each items as item, index (item.queueItemId)}
        {@const isCurrent = index === currentIndex}
        <button
          class="group flex w-full items-center gap-3 rounded-md p-2 text-left transition-colors hover:bg-accent {isCurrent
            ? 'bg-accent/80'
            : ''}"
          onclick={async () => {
            await musicManager.setQueueIndex(index);
            musicManager.play();
          }}
        >
          <span
            class="w-8 flex-shrink-0 text-right text-sm tabular-nums text-muted-foreground"
          >
            {index + 1}
          </span>

          <Image
            class="h-12 w-12 flex-shrink-0 rounded object-cover"
            src={item.coverArt}
            alt={item.name}
          />

          <div class="flex min-w-0 flex-1 flex-col">
            <p class="truncate font-medium">{item.name}</p>
            <p class="truncate text-sm text-muted-foreground">
              {item.artistNames}
            </p>
          </div>

          <span class="flex-shrink-0 text-sm text-muted-foreground">
            {formatDuration(item.duration)}
          </span>

          {#if isCurrent}
            <Play class="ml-2 h-5 w-5 flex-shrink-0 text-primary" />
          {:else}
            <Play
              class="ml-2 h-5 w-5 flex-shrink-0 text-muted-foreground opacity-0 transition-opacity group-hover:opacity-100"
            />
          {/if}
        </button>
      {/each}
    </div>
  {/if}

  <div bind:this={sentinel} class="h-4"></div>

  {#if loading}
    <p class="py-6 text-center text-sm text-muted-foreground">Loading...</p>
  {/if}
</div>
