<script lang="ts">
  import { PUBLIC_COMMIT, PUBLIC_VERSION } from "$env/static/public";
  import { getApiClient, handleApiError } from "$lib";
  import { formatDuration } from "$lib/utils.js";
  import { Button, Card, Separator } from "@nanoteck137/nano-ui";
  import {
    AlertCircle,
    DiscAlbum,
    FileMusic,
    Play,
    RefreshCw,
    Server,
    Users,
  } from "lucide-svelte";
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
    displayName: z.string(),
    isRunning: z.boolean(),
  });
  type TaskSyncStateEventTaskTy = z.infer<typeof TaskSyncStateEventTask>;

  const TaskSyncStateEvent = z.object({
    tasks: z.array(TaskSyncStateEventTask),
  });

  let tasks = $state<TaskSyncStateEventTaskTy[]>([]);

  async function setupEventSource(): Promise<EventSource | null> {
    const res = await apiClient.createSseToken();
    if (!res.success) {
      handleApiError(res.error);
      return null;
    }

    const url = apiClient.url.sseHandler();
    url.searchParams.set("token", res.data.token);

    const eventSource = new EventSource(url);

    eventSource.addEventListener("connected", () => {
      console.log("Connected to SSE handler");
    });

    eventSource.addEventListener("library-sync-state", (e) => {
      const data = LibrarySyncStateEvent.parse(JSON.parse(e.data));

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

      tasks = data.tasks;
    });

    return eventSource;
  }

  let eventSource = $state<EventSource | null>(null);

  onMount(() => {
    setupEventSource().then((e) => {
      eventSource = e;
    });

    return () => {
      eventSource?.close();
    };
  });
</script>

<div class="flex flex-col gap-6">
  <div class="flex items-center gap-3">
    <Server size={24} />
    <h1 class="text-xl font-bold">Server</h1>
    <span class="text-xs text-muted-foreground">
      Server: {data.systemInfo.version} | Client: v{PUBLIC_VERSION} ({PUBLIC_COMMIT})
    </span>
  </div>

  <div class="grid grid-cols-3 gap-4">
    <Card.Root>
      <div class="flex flex-col gap-2 p-4">
        <div class="flex items-center gap-2 text-sm text-muted-foreground">
          <Users size={16} />
          <span>Artists</span>
        </div>
        <span class="text-2xl font-bold">{numArtists}</span>
        <span class="text-xs text-muted-foreground">
          Last sync: {formatDuration(artistSyncTime)}
        </span>
      </div>
    </Card.Root>

    <Card.Root>
      <div class="flex flex-col gap-2 p-4">
        <div class="flex items-center gap-2 text-sm text-muted-foreground">
          <DiscAlbum size={16} />
          <span>Albums</span>
        </div>
        <span class="text-2xl font-bold">{numAlbums}</span>
        <span class="text-xs text-muted-foreground">
          Last sync: {formatDuration(albumSyncTime)}
        </span>
      </div>
    </Card.Root>

    <Card.Root>
      <div class="flex flex-col gap-2 p-4">
        <div class="flex items-center gap-2 text-sm text-muted-foreground">
          <FileMusic size={16} />
          <span>Tracks</span>
        </div>
        <span class="text-2xl font-bold">{numTracks}</span>
        <span class="text-xs text-muted-foreground">
          Last sync: {formatDuration(trackSyncTime)}
        </span>
      </div>
    </Card.Root>
  </div>

  <div class="grid grid-cols-1 gap-6 lg:grid-cols-2">
    <Card.Root>
      <div class="p-6">
        <div class="flex items-center gap-2">
          <RefreshCw size={18} />
          <h2 class="text-lg font-semibold">Tasks</h2>
        </div>

        <Separator class="my-4" />

        <div class="flex flex-col gap-2">
          {#each tasks as task (task.name)}
            <div
              class="flex items-center justify-between rounded-lg border p-3"
            >
              <div class="flex flex-col">
                <span class="text-sm font-medium">{task.displayName}</span>
                <span class="text-xs text-muted-foreground">
                  {task.isRunning ? "Running..." : "Idle"}
                </span>
              </div>
              {#if !task.isRunning}
                <Button
                  variant="outline"
                  size="sm"
                  onclick={async () => {
                    const res = await apiClient.runTask(task.name);
                    if (!res.success) {
                      return handleApiError(res.error);
                    }

                    toast.success("Dispatched task");
                  }}
                >
                  <Play size={14} />
                  Run
                </Button>
              {/if}
            </div>
          {/each}
        </div>
      </div>
    </Card.Root>

    <Card.Root>
      <div class="p-6">
        <div class="flex items-center gap-2">
          <AlertCircle size={18} />
          <h2 class="text-lg font-semibold">Missing Items</h2>
        </div>

        <Separator class="my-4" />

        <div class="flex flex-col gap-4">
          {#if missingArtists.length > 0}
            <div>
              <span class="text-xs font-medium text-muted-foreground"
                >Artists ({missingArtists.length})</span
              >
              <div class="mt-1 flex flex-col">
                {#each missingArtists as artist (artist.id)}
                  <a
                    href="/artists/{artist.id}"
                    class="text-sm hover:underline">{artist.name}</a
                  >
                {/each}
              </div>
            </div>
          {/if}

          {#if missingAlbums.length > 0}
            <div>
              <span class="text-xs font-medium text-muted-foreground"
                >Albums ({missingAlbums.length})</span
              >
              <div class="mt-1 flex flex-col">
                {#each missingAlbums as album (album.id)}
                  <a href="/albums/{album.id}" class="text-sm hover:underline"
                    >{album.name}</a
                  >
                {/each}
              </div>
            </div>
          {/if}

          {#if missingTracks.length > 0}
            <div>
              <span class="text-xs font-medium text-muted-foreground"
                >Tracks ({missingTracks.length})</span
              >
              <div class="mt-1 flex flex-col">
                {#each missingTracks as track (track.name)}
                  <span class="text-sm">{track.name}</span>
                {/each}
              </div>
            </div>
          {/if}

          {#if missingArtists.length === 0 && missingAlbums.length === 0 && missingTracks.length === 0}
            <p class="text-sm text-muted-foreground">No missing items.</p>
          {/if}
        </div>
      </div>
    </Card.Root>
  </div>

  {#if errors.length > 0}
    <Card.Root>
      <div class="p-6">
        <div class="flex items-center gap-2 text-destructive">
          <AlertCircle size={18} />
          <h2 class="text-lg font-semibold">Errors</h2>
        </div>

        <Separator class="my-4" />

        <div class="flex flex-col gap-1">
          {#each errors as err (err)}
            <p class="font-mono text-sm text-destructive">{err}</p>
          {/each}
        </div>
      </div>
    </Card.Root>
  {/if}

  <Card.Root>
    <div class="p-6">
      <div class="flex items-center gap-2">
        <Server size={18} />
        <h2 class="text-lg font-semibold">Media Configuration</h2>
      </div>

      <Separator class="my-4" />

      <div class="flex flex-col gap-6">
        <div>
          <h3 class="mb-2 text-sm font-medium text-muted-foreground">
            Formats
          </h3>
          <div class="overflow-x-auto">
            <table class="w-full text-left text-sm">
              <thead>
                <tr class="text-muted-foreground">
                  <th class="px-3 py-2 font-medium">Name</th>
                  <th class="px-3 py-2 font-medium">Format</th>
                  <th class="px-3 py-2 font-medium">Ext</th>
                  <th class="px-3 py-2 font-medium">High</th>
                  <th class="px-3 py-2 font-medium">Medium</th>
                  <th class="px-3 py-2 font-medium">Low</th>
                </tr>
              </thead>
              <tbody>
                {#each data.mediaSettings.formats as format (format.name)}
                  <tr class="border-t">
                    <td class="px-3 py-2">{format.name}</td>
                    <td class="px-3 py-2 font-mono text-xs">{format.format}</td
                    >
                    <td class="px-3 py-2 font-mono text-xs">{format.ext}</td>
                    <td class="px-3 py-2">{format.qualityHighBitrate}kbps</td>
                    <td class="px-3 py-2">{format.qualityMediumBitrate}kbps</td
                    >
                    <td class="px-3 py-2">{format.qualityLowBitrate}kbps</td>
                  </tr>
                {/each}
              </tbody>
            </table>
          </div>
        </div>

        <Separator />

        <div>
          <h3 class="mb-2 text-sm font-medium text-muted-foreground">
            Device Specs
          </h3>
          <div class="overflow-x-auto">
            <table class="w-full text-left text-sm">
              <thead>
                <tr class="text-muted-foreground">
                  <th class="px-3 py-2 font-medium">Name</th>
                  <th class="px-3 py-2 font-medium">Preferred Format</th>
                  <th class="px-3 py-2 font-medium">Allowed Formats</th>
                </tr>
              </thead>
              <tbody>
                {#each data.mediaSettings.deviceSpecs as spec (spec.name)}
                  <tr class="border-t">
                    <td class="px-3 py-2">{spec.name}</td>
                    <td class="px-3 py-2 font-mono text-xs"
                      >{spec.preferedFormat}</td
                    >
                    <td class="px-3 py-2 font-mono text-xs"
                      >{spec.allowedFormats.join(", ")}</td
                    >
                  </tr>
                {/each}
              </tbody>
            </table>
          </div>
        </div>
      </div>
    </div>
  </Card.Root>
</div>
