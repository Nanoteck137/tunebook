<script lang="ts">
  import { PUBLIC_COMMIT, PUBLIC_VERSION } from "$env/static/public";
  import { getApiClient, handleApiError } from "$lib";
  import { formatDuration } from "$lib/utils.js";
  import { Button } from "@nanoteck137/nano-ui";
  import { Play } from "lucide-svelte";
  import { onMount } from "svelte";
  import toast from "svelte-5-french-toast";
  import { z } from "zod";

  const { data } = $props();
  const apiClient = getApiClient();

  let errors = $state<string[]>([]);
  let numArtists = $state(0);
  let numAlbums = $state(0);
  let numTracks = $state(0);

  let missingArtists = $state<MissingItemTy[]>([]);
  let missingAlbums = $state<MissingItemTy[]>([]);
  let missingTracks = $state<MissingItemTy[]>([]);

  let artistSyncTime = $state(0);
  let albumSyncTime = $state(0);
  let trackSyncTime = $state(0);
  let totalSyncTime = $state(0);

  const MissingItem = z.object({
    id: z.string(),
    name: z.string(),
  });
  type MissingItemTy = z.infer<typeof MissingItem>;

  const LibrarySyncStateEvent = z.object({
    errors: z.array(z.string()),

    numArtists: z.number(),
    numAlbums: z.number(),
    numTracks: z.number(),

    missingArtists: z.array(MissingItem),
    missingAlbums: z.array(MissingItem),
    missingTracks: z.array(MissingItem),

    artistsSyncDurationMs: z.number(),
    albumsSyncDurationMs: z.number(),
    tracksSyncDurationMs: z.number(),
    totalSyncDurationMs: z.number(),
  });

  const TaskSyncStateEventTask = z.object({
    name: z.string(),
    isRunning: z.boolean(),
  });
  type TaskSyncStateEventTaskTy = z.infer<typeof TaskSyncStateEventTask>;

  const TaskSyncStateEvent = z.object({
    tasks: z.array(TaskSyncStateEventTask),
  });

  let tasks = $state<TaskSyncStateEventTaskTy[]>([]);

  onMount(() => {
    const eventSource = new EventSource(apiClient.url.sseHandler());

    eventSource.addEventListener("connected", () => {
      console.log("Connected to SSE handler");
    });

    eventSource.addEventListener("library-sync-state", (e) => {
      console.log(JSON.parse(e.data));
      const data = LibrarySyncStateEvent.parse(JSON.parse(e.data));

      console.log("library state", data);

      errors = data.errors;
      numArtists = data.numArtists;
      numAlbums = data.numAlbums;
      numTracks = data.numTracks;

      missingArtists = data.missingArtists;
      missingAlbums = data.missingAlbums;
      missingTracks = data.missingTracks;

      artistSyncTime = data.artistsSyncDurationMs;
      albumSyncTime = data.albumsSyncDurationMs;
      trackSyncTime = data.tracksSyncDurationMs;
      totalSyncTime = data.totalSyncDurationMs;
    });

    eventSource.addEventListener("task-sync-state", (e) => {
      const data = TaskSyncStateEvent.parse(JSON.parse(e.data));

      console.log("tasks state", data);
      tasks = data.tasks;

      // isSyncing = data.isRunning;
      // errors = data.errors;
      // numArtists = data.numArtists;
      // numAlbums = data.numAlbums;
      // numTracks = data.numTracks;

      // artistSyncTime = data.artistsSyncDurationMs;
      // albumSyncTime = data.albumsSyncDurationMs;
      // trackSyncTime = data.tracksSyncDurationMs;
      // totalSyncTime = data.totalSyncDurationMs;
    });

    return () => {
      eventSource.close();
    };
  });
</script>

<p>Server Page (W.I.P)</p>

<p>Version: {PUBLIC_VERSION}</p>
<p>Commit: {PUBLIC_COMMIT}</p>

{#each tasks as task}
  <div class="flex items-center gap-2">
    <p>{task.name} - Running: {task.isRunning}</p>
    {#if !task.isRunning}
      <Button
        variant="ghost"
        size="icon"
        onclick={async () => {
          const res = await apiClient.runTask(task.name);
          if (!res.success) {
            return handleApiError(res.error);
          }

          toast.success("Dispatched task");
        }}
      >
        <Play />
      </Button>
    {/if}
  </div>
{/each}

<p>Num Artists: {numArtists}</p>
<p>Num Albums: {numAlbums}</p>
<p>Num Tracks: {numTracks}</p>

<br />

<p>Artist Sync Time: {formatDuration(artistSyncTime)}</p>
<p>Album Sync Time: {formatDuration(albumSyncTime)}</p>
<p>Track Sync Time: {formatDuration(trackSyncTime)}</p>
<p>Total Sync Time: {formatDuration(totalSyncTime)}</p>

<p>Missing Artists:</p>
<div class="flex flex-col">
  {#each missingArtists as artist}
    <a href="/artists/{artist.id}">{artist.name}</a>
  {/each}
</div>

<p>Missing Albums:</p>
<div class="flex flex-col">
  {#each missingAlbums as album}
    <a href="/albums/{album.id}">{album.name}</a>
  {/each}
</div>

<p>Missing Tracks:</p>
<div class="flex flex-col">
  {#each missingTracks as track}
    <p>{track.name}</p>
  {/each}
</div>

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
