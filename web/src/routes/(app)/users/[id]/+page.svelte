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
  import type { Track } from "$lib/api/types";
  import TrackListItem from "$lib/components/track-list/TrackListItem.svelte";

  type YearStats = {
    year: number;
    trackCount: number;
    listeningTime: number;
  };

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
  let topTracks: Track[] = [
    {
      id: "mock-1",
      name: "Bohemian Rhapsody",
      order: 1,
      duration: 354,
      number: 1,
      year: 1975,
      coverArt: { original: "", small: "", medium: "", large: "" },
      albumId: "mock-album-1",
      albumName: "A Night at the Opera",
      artists: [{ id: "mock-artist-1", name: "Queen" }],
      tags: ["rock", "classic"],
      created: "2024-01-01T00:00:00Z",
      updated: "2024-01-01T00:00:00Z",
    },
    {
      id: "mock-2",
      name: "Stairway to Heaven",
      order: 2,
      duration: 482,
      number: 1,
      year: 1971,
      coverArt: { original: "", small: "", medium: "", large: "" },
      albumId: "mock-album-2",
      albumName: "Led Zeppelin IV",
      artists: [{ id: "mock-artist-2", name: "Led Zeppelin" }],
      tags: ["rock", "classic"],
      created: "2024-01-01T00:00:00Z",
      updated: "2024-01-01T00:00:00Z",
    },
    {
      id: "mock-3",
      name: "Hotel California",
      order: 3,
      duration: 391,
      number: 1,
      year: 1976,
      coverArt: { original: "", small: "", medium: "", large: "" },
      albumId: "mock-album-3",
      albumName: "Hotel California",
      artists: [{ id: "mock-artist-3", name: "Eagles" }],
      tags: ["rock", "classic"],
      created: "2024-01-01T00:00:00Z",
      updated: "2024-01-01T00:00:00Z",
    },
    {
      id: "mock-4",
      name: "Smells Like Teen Spirit",
      order: 4,
      duration: 301,
      number: 1,
      year: 1991,
      coverArt: { original: "", small: "", medium: "", large: "" },
      albumId: "mock-album-4",
      albumName: "Nevermind",
      artists: [{ id: "mock-artist-4", name: "Nirvana" }],
      tags: ["grunge", "rock"],
      created: "2024-01-01T00:00:00Z",
      updated: "2024-01-01T00:00:00Z",
    },
    {
      id: "mock-5",
      name: "Imagine",
      order: 5,
      duration: 187,
      number: 1,
      year: 1971,
      coverArt: { original: "", small: "", medium: "", large: "" },
      albumId: "mock-album-5",
      albumName: "Imagine",
      artists: [{ id: "mock-artist-5", name: "John Lennon" }],
      tags: ["rock", "classic"],
      created: "2024-01-01T00:00:00Z",
      updated: "2024-01-01T00:00:00Z",
    },
  ];

  // TODO: Replace with API data when available
  let yearStats: YearStats[] = [
    { year: 2026, trackCount: 2847, listeningTime: 184320 },
    { year: 2025, trackCount: 12563, listeningTime: 783600 },
    { year: 2024, trackCount: 8931, listeningTime: 542700 },
    { year: 2023, trackCount: 4215, listeningTime: 261900 },
  ];

  let maxTrackCount = $derived(
    Math.max(...yearStats.map((y) => y.trackCount)),
  );

  function formatListeningTime(seconds: number): string {
    const hours = Math.floor(seconds / 3600);
    const minutes = Math.floor((seconds % 3600) / 60);
    return `${hours}h ${minutes}m`;
  }
</script>

<div class="flex flex-col gap-6">
  <div class="grid grid-cols-2 gap-4 lg:grid-cols-4">
    {#each stats as stat}
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

  <!-- TODO: Replace with API data when available -->
  <!-- <Card.Root>
    <div class="p-6">
      <div class="flex items-center gap-2">
        <Clock size={18} />
        <h2 class="text-lg font-semibold">Recently Played</h2>
      </div>
      <p class="mt-1 text-sm text-muted-foreground">
        Recently played tracks.
      </p>

      <Separator class="my-4" />

      {#if recentTracks.length === 0}
        <div class="flex flex-col items-center gap-2 py-8 text-sm text-muted-foreground">
          <Play size={24} class="text-muted-foreground/50" />
          <p>No listening history yet</p>
          <p class="text-xs">Start playing music to see your history here.</p>
        </div>
      {/if}
    </div>
  </Card.Root> -->

  <Card.Root>
    <div class="p-6">
      <div class="flex items-center gap-2">
        <TrendingUp size={18} />
        <h2 class="text-lg font-semibold">Top Tracks</h2>
      </div>
      <p class="mt-1 text-sm text-muted-foreground">
        This section will show the tracks {data.userData.displayName} listens to
        the most.
      </p>

      <Separator class="my-4" />

      <div class="flex flex-col">
        {#each topTracks as track}
          <TrackListItem {track} onPlayClicked={() => {}} />
        {/each}
      </div>
    </div>
  </Card.Root>

  <!-- TODO: Replace with API data when available -->
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
        {#each yearStats as stat}
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
                  style="width: {(stat.trackCount / maxTrackCount) * 100}%"
                ></div>
              </div>
            </div>
          </a>
        {/each}
      </div>
    </div>
  </Card.Root>
</div>
