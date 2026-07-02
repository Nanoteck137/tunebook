<script lang="ts">
  import { Card, Separator, Button } from "@nanoteck137/nano-ui";
  import { History, X } from "lucide-svelte";
  import { goto } from "$app/navigation";
  import { page } from "$app/state";
  import { getMusicManager } from "$lib/music-manager.svelte";
  import Pagination from "$lib/components/Pagination.svelte";

  let { data } = $props();

  const musicManager = getMusicManager();

  let yearParam = $derived(page.url.searchParams.get("year"));

  function formatRelativeTime(millis: number): string {
    const unixSeconds = Math.floor(millis / 1000);
    const now = Math.floor(Date.now() / 1000);
    const diff = now - unixSeconds;

    if (diff < 0) return "just now";
    if (diff < 60) return "just now";
    if (diff < 3600) {
      const m = Math.floor(diff / 60);
      return `${m}m ago`;
    }
    if (diff < 86400) {
      const h = Math.floor(diff / 3600);
      return `${h}h ago`;
    }
    if (diff < 604800) {
      const d = Math.floor(diff / 86400);
      return `${d}d ago`;
    }
    return new Date(unixSeconds * 1000).toLocaleDateString(undefined, {
      month: "short",
      day: "numeric",
    });
  }

  function statusLabel(status: string) {
    if (status === "completed") return "Completed";
    if (status === "skipped") return "Skipped";
    return "In Progress";
  }

  function percentColor(pct: number): string {
    if (pct >= 80) return "bg-green-500";
    if (pct >= 40) return "bg-yellow-500";
    return "bg-muted-foreground/40";
  }

  function clearYearFilter() {
    const query = page.url.searchParams;
    query.delete("year");
    goto(`?${query.toString()}`, { invalidateAll: true });
  }

  function playTrack(trackId: string) {
    const trackIds = data.history.map((e) => e.track.id);
    musicManager.addTracks({ trackIds, trackId });
  }
</script>

<div class="flex flex-col gap-6">
  <div class="flex items-center justify-between">
    <div class="flex items-baseline gap-2">
      <h1 class="text-xl font-bold">
        {#if yearParam}
          History for {yearParam}
        {:else}
          Listening History
        {/if}
      </h1>
      <span class="text-sm text-muted-foreground">{data.page.totalItems}</span>
    </div>

    {#if yearParam}
      <Button variant="outline" size="sm" onclick={clearYearFilter}>
        <X size={14} />
        Clear year
      </Button>
    {/if}
  </div>

  <Card.Root>
    {#if data.history.length === 0}
      <div class="flex flex-col items-center gap-2 py-16">
        <History size={32} class="text-muted-foreground/40" />
        <p class="text-sm text-muted-foreground">No listening history yet</p>
      </div>
    {:else}
      {#each data.history as entry (entry.id)}
        <div
          class="group flex items-center gap-3 px-4 py-2 transition-colors hover:bg-muted/30"
        >
          <button
            class="shrink-0"
            onclick={() => playTrack(entry.track.id)}
            aria-label="Play {entry.track.name}"
          >
            <img
              src={entry.track.coverArt.small}
              alt={entry.track.name}
              class="h-10 w-10 rounded object-cover"
            />
          </button>

          <div class="flex min-w-0 flex-1 flex-col">
            <div class="flex items-center gap-2">
              <span class="truncate text-sm font-medium" title={entry.track.name}              >
                {entry.track.name}
              </span>
              <span
                class="shrink-0 rounded bg-muted px-1.5 py-0.5 text-[10px] font-medium {entry.status ===
                'completed'
                  ? 'text-green-500'
                  : entry.status === 'skipped'
                    ? 'text-muted-foreground'
                    : 'text-yellow-500'}"
              >
                {statusLabel(entry.status)}
              </span>
            </div>
            <div
              class="flex items-center gap-1.5 text-xs text-muted-foreground"
            >
              <span class="truncate">
                {entry.track.artists.map((a) => a.name).join(", ")}
              </span>
              {#if entry.track.albumName}
                <span class="shrink-0">&middot;</span>
                <span class="truncate">{entry.track.albumName}</span>
              {/if}
            </div>
          </div>

          <div class="hidden items-center gap-3 sm:flex">
            {#if entry.status !== "skipped"}
              <div class="flex items-center gap-1.5">
                <div class="h-1.5 w-16 overflow-hidden rounded-full bg-muted">
                  <div
                    class="h-full rounded-full transition-all {percentColor(
                      entry.percentPlayed,
                    )}"
                    style="width: {entry.percentPlayed}%"
                  ></div>
                </div>
                <span class="text-[11px] text-muted-foreground"
                  >{Math.round(entry.percentPlayed)}%</span
                >
              </div>
            {/if}
          </div>

          <span class="shrink-0 text-[11px] text-muted-foreground">
            {formatRelativeTime(entry.listenedAt)}
          </span>
        </div>
      {/each}
    {/if}
  </Card.Root>

  <Separator />

  <Pagination page={data.page} />
</div>
