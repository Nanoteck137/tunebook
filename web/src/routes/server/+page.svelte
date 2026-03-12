<script lang="ts">
  import { PUBLIC_COMMIT, PUBLIC_VERSION } from "$env/static/public";
  import { getApiClient, handleApiError } from "$lib";
  import { formatDuration } from "$lib/utils.js";
  import { Button } from "@nanoteck137/nano-ui";
  import { onDestroy, onMount } from "svelte";
  import { z } from "zod";

  const { data } = $props();
  const apiClient = getApiClient();

  let isSyncing = $state(false);
  let errors = $state<string[]>([]);
  let numArtists = $state(0);
  let numAlbums = $state(0);
  let numTracks = $state(0);

  let artistSyncTime = $state(0);
  let albumSyncTime = $state(0);
  let trackSyncTime = $state(0);
  let totalSyncTime = $state(0);

  const SyncError = z.object({
    type: z.string(),
    message: z.string(),
    fullMessage: z.string().optional(),
  });

  const MissingAlbum = z.object({
    id: z.string(),
    name: z.string(),
    artistName: z.string(),
  });

  const MissingTrack = z.object({
    id: z.string(),
    name: z.string(),
    albumName: z.string(),
    artistName: z.string(),
  });

  const SyncState = z.object({
    isSyncing: z.boolean(),
    isRetrivingPaths: z.boolean(),

    paths: z.array(
      z.object({
        name: z.string(),
        path: z.string(),
        isDir: z.boolean(),
        depth: z.number(),
      }),
    ),

    report: z.object({
      syncErrors: z.array(SyncError).nullable(),
      missingAlbums: z.array(MissingAlbum).nullable(),
      missingTracks: z.array(MissingTrack).nullable(),
    }),
  });
  type SyncStateTy = z.infer<typeof SyncState>;

  const Event = z.discriminatedUnion("type", [
    z.object({
      type: z.literal("sync-state"),
      data: SyncState,
    }),
    z.object({
      type: z.literal("syncing"),
      data: z.object({
        syncing: z.boolean(),
      }),
    }),
    z.object({
      type: z.literal("report"),
      data: z.object({
        syncErrors: z.array(SyncError).nullable(),
        missingAlbums: z.array(MissingAlbum).nullable(),
        missingTracks: z.array(MissingTrack).nullable(),
      }),
    }),
  ]);

  const LibrarySyncStateEvent = z.object({
    isRunning: z.boolean(),
    errors: z.array(z.string()),

    numArtists: z.number(),
    numAlbums: z.number(),
    numTracks: z.number(),

    artistsSyncDurationMs: z.number(),
    albumsSyncDurationMs: z.number(),
    tracksSyncDurationMs: z.number(),
    totalSyncDurationMs: z.number(),
  });

  onMount(() => {
    const eventSource = new EventSource(apiClient.url.sseHandler());

    eventSource.addEventListener("connected", () => {
      console.log("Connected to SSE handler");
    });

    eventSource.addEventListener("library-sync-state", (e) => {
      const data = LibrarySyncStateEvent.parse(JSON.parse(e.data));

      console.log("library state", data);

      isSyncing = data.isRunning;
      errors = data.errors;
      numArtists = data.numArtists;
      numAlbums = data.numAlbums;
      numTracks = data.numTracks;

      artistSyncTime = data.artistsSyncDurationMs;
      albumSyncTime = data.albumsSyncDurationMs;
      trackSyncTime = data.tracksSyncDurationMs;
      totalSyncTime = data.totalSyncDurationMs;
    });

    return () => {
      eventSource.close();
    };
  });
</script>

<p>Server Page (W.I.P)</p>

<p>Version: {PUBLIC_VERSION}</p>
<p>Commit: {PUBLIC_COMMIT}</p>

<p>Library Syncing: {isSyncing}</p>

<Button
  onclick={async () => {
    const res = await apiClient.syncLibrary({});
    if (!res.success) {
      handleApiError(res.error);
      return;
    }
  }}
>
  Sync Library
</Button>

<p>Num Artists: {numArtists}</p>
<p>Num Albums: {numAlbums}</p>
<p>Num Tracks: {numTracks}</p>

<br />

<p>Artist Sync Time: {formatDuration(artistSyncTime)}</p>
<p>Album Sync Time: {formatDuration(albumSyncTime)}</p>
<p>Track Sync Time: {formatDuration(trackSyncTime)}</p>
<p>Total Sync Time: {formatDuration(totalSyncTime)}</p>

<p>Errors:</p>
{#each errors as err}
  <p>Error: {err}</p>
{/each}

<!-- {#if syncState?.report.syncErrors}
  {#each syncState?.report.syncErrors as err}
    <p class="whitespace-pre font-mono">{err.fullMessage}</p>
    <br />
  {/each}
{/if} -->

<!-- {#if syncState?.paths}
  <div class="flex flex-col items-start">
    {#each syncState?.paths as path}
      <button
        class="hover:underline"
        style={`padding-left: ${path.depth * 20}px`}
        onclick={async () => {
          const res = await apiClient.syncLibrary({ path: path.path });
          if (!res.success) {
            handleApiError(res.error);
            return;
          }
        }}
      >
        {path.path}
      </button>
    {/each}
  </div>
{/if} -->

<div class="h-10"></div>

<p>Media Formats:</p>
<div class="flex flex-col gap-2">
  {#each data.mediaSettings.formats as format}
    <div>
      <p>Name: {format.name}</p>
      <p>Format: {format.format}</p>
      <p>Ext: {format.ext}</p>
      <p>High: {format.qualityHighBitrate}</p>
      <p>Medium: {format.qualityMediumBitrate}</p>
      <p>Low: {format.qualityLowBitrate}</p>
    </div>
  {/each}
</div>

<div class="h-10"></div>

<p>Media Device Specs:</p>
<div class="flex flex-col gap-2">
  {#each data.mediaSettings.deviceSpecs as spec}
    <div>
      <p>Name: {spec.name}</p>
      <p>Prefered Format: {spec.preferedFormat}</p>
      <p>Allowed Formats: {spec.allowedFormats.join(",")}</p>
    </div>
  {/each}
</div>
