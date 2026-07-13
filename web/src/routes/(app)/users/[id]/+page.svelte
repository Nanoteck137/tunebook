<script lang="ts">
  import { Card, Separator } from "@nanoteck137/nano-ui";
  import {
    BarChart3,
    Clock,
    Heart,
    ListMusic,
    Play,
    TrendingUp,
  } from "lucide-svelte";
  import TrackListItem from "$lib/components/track-list/TrackListItem.svelte";

  let { data } = $props();

  let stats = $derived([
    {
      icon: Play,
      label: "Tracks Played",
      value: data.stats.numTracksPlayed.toLocaleString(),
      note: "",
    },
    {
      icon: Clock,
      label: "Listening Time",
      value: formatListeningTime(data.stats.listeningTime),
      note: "",
    },
    {
      icon: ListMusic,
      label: "Playlists",
      value: data.stats.numPlaylistsCreated.toLocaleString(),
      note: "",
    },
    {
      icon: Heart,
      label: "Favorites",
      value: data.stats.numFavoriteTracks.toLocaleString(),
      note: "",
    },
  ]);

  // TODO: Replace with API data when available
  // let recentTracks: Track[] = [];

  let maxTrackCount = $derived(
    Math.max(...data.yearStats.map((y) => y.trackCount), 0),
  );

  function formatListeningTime(seconds: number): string {
    const hours = Math.floor(seconds / 3600);
    const minutes = Math.floor((seconds % 3600) / 60);
    return `${hours}h ${minutes}m`;
  }
</script>

<div class="flex flex-col gap-6">
  <div class="grid grid-cols-2 gap-4 lg:grid-cols-4">
    {#each stats as stat (stat.label)}
      <Card.Root>
        <div class="flex flex-col gap-2 p-4">
          <div class="flex items-center gap-2">
            <stat.icon size={18} class="text-muted-foreground" />
            <span class="text-sm text-muted-foreground">{stat.label}</span>
          </div>
          <span class="text-2xl font-bold">{stat.value}</span>
          {#if stat.note}
            <span class="text-xs text-muted-foreground">{stat.note}</span>
          {/if}
        </div>
      </Card.Root>
    {/each}
  </div>

  <!-- Top Tracks -->
  <Card.Root>
    <div class="p-6">
      <div class="flex items-center gap-2">
        <TrendingUp size={18} />
        <h2 class="text-lg font-semibold">Top Tracks</h2>
      </div>
      <p class="mt-1 text-sm text-muted-foreground">
        The tracks {data.userData.displayName} listens to the most.
      </p>

      <Separator class="my-4" />

      <div class="flex flex-col">
        {#if data.topTracks.length === 0}
          <p class="text-sm text-muted-foreground">No tracks played yet.</p>
        {:else}
          {#each data.topTracks as track (track.id)}
            <TrackListItem {track} onPlayClicked={() => {}} />
          {/each}
        {/if}
      </div>
    </div>
  </Card.Root>

  <!-- Year in Review -->
  <Card.Root>
    <div class="p-6">
      <div class="flex items-center gap-2">
        <BarChart3 size={18} />
        <h2 class="text-lg font-semibold">Year in Review</h2>
      </div>
      <p class="mt-1 text-sm text-muted-foreground">
        Year-over-year listening statistics.
      </p>

      <Separator class="my-4" />

      <div class="flex flex-col gap-4">
        {#if data.yearStats.length === 0}
          <p class="text-sm text-muted-foreground">No listening data yet.</p>
        {:else}
          {#each data.yearStats as stat (stat.year)}
            <a
              href="/users/{data.userData.id}/history?year={stat.year}"
              class="flex items-center gap-4 rounded-lg px-2 py-1.5 transition-colors hover:bg-muted/50"
            >
              <span class="w-12 text-sm font-medium">{stat.year}</span>
              <div class="flex flex-1 flex-col gap-1">
                <div
                  class="flex items-center justify-between text-xs text-muted-foreground"
                >
                  <span>{stat.trackCount.toLocaleString()} tracks</span>
                  <span>{formatListeningTime(stat.listeningTime)}</span>
                </div>
                <div class="h-2 w-full rounded-full bg-muted">
                  <div
                    class="h-2 rounded-full bg-gradient-to-r from-logo-1 to-logo-3"
                    style="width: {maxTrackCount > 0
                      ? (stat.trackCount / maxTrackCount) * 100
                      : 0}%"
                  ></div>
                </div>
              </div>
            </a>
          {/each}
        {/if}
      </div>
    </div>
  </Card.Root>
</div>
