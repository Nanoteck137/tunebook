<script lang="ts">
	import { formatDuration } from "$lib/utils.js";
	import { getApiClient, handleApiError } from "$lib";
	import {
		AlertCircle,
		DiscAlbum,
		FileMusic,
		Library,
		Loader2,
		Users,
	} from "@lucide/svelte";
	import { onMount } from "svelte";
	import { toast } from "svelte-sonner";
	import { Button, Card, Separator } from "$lib/components/ui";
	import type { LibrarySyncStateEventTy } from "../events";
	import {
		connectServerSse,
		JobSyncStateEvent,
		LibrarySyncStateEvent,
	} from "../events";

	const apiClient = getApiClient();

	const initial: LibrarySyncStateEventTy = {
		errors: [],
		currentItem: null,
		numArtists: 0,
		numAlbums: 0,
		numTracks: 0,
		missingArtists: [],
		missingAlbums: [],
		missingTracks: [],
		artistsSyncDurationMs: 0,
		albumsSyncDurationMs: 0,
		tracksSyncDurationMs: 0,
		totalSyncDurationMs: 0,
	};

	let syncState = $state<LibrarySyncStateEventTy>(initial);

	let syncing = $state(false);
	let eventSource = $state<EventSource | null>(null);

	async function runSync() {
		const res = await apiClient.runTask("library-sync");
		if (!res.success) {
			return handleApiError(res.error);
		}

		toast.success("Dispatched library sync");
	}

	onMount(() => {
		connectServerSse({
			"library-sync-state": (data) => {
				syncState = LibrarySyncStateEvent.parse(data);
			},
			"job-sync-state": (data) => {
				const event = JobSyncStateEvent.parse(data);
				syncing =
					event.jobs.find((job) => job.name === "library-sync")?.status ===
					"running";
			},
		}).then((e) => {
			eventSource = e;
		});

		return () => {
			eventSource?.close();
		};
	});
</script>

<div class="flex flex-col gap-6">
	<div class="flex items-center justify-between">
		<div class="flex items-center gap-3">
			<Library size={24} />
			<div class="flex flex-col">
				<h1 class="text-xl font-bold">Library</h1>
				<p class="text-xs text-muted-foreground">
					Last sync: {formatDuration(syncState.totalSyncDurationMs)}
				</p>
			</div>
		</div>

		<Button onclick={runSync} disabled={syncing}>
			{#if syncing}
				<Loader2 size={14} class="animate-spin" />
				Syncing...
			{:else}
				Sync now
			{/if}
		</Button>
	</div>

	{#if syncing && syncState.currentItem}
		<div class="flex items-center gap-3 rounded border p-3">
			<Loader2 size={16} class="shrink-0 animate-spin text-muted-foreground" />
			<span class="text-sm">
				Syncing {syncState.currentItem.kind}:
				<span class="font-medium">{syncState.currentItem.name}</span>
			</span>
			{#if syncState.currentItem.info}
				<span
					class="ml-auto max-w-64 truncate font-mono text-xs text-muted-foreground"
					title={syncState.currentItem.info}
				>
					{syncState.currentItem.info}
				</span>
			{/if}
		</div>
	{/if}

	<div class="grid grid-cols-3 gap-4">
		<Card.Root>
			<Card.Content class="flex flex-col gap-2">
				<div class="flex items-center gap-2 text-sm text-muted-foreground">
					<Users size={16} />
					<span>Artists</span>
				</div>
				<span class="text-2xl font-bold">{syncState.numArtists}</span>
				<span class="text-xs text-muted-foreground">
					Last sync: {formatDuration(syncState.artistsSyncDurationMs)}
				</span>
			</Card.Content>
		</Card.Root>

		<Card.Root>
			<Card.Content class="flex flex-col gap-2">
				<div class="flex items-center gap-2 text-sm text-muted-foreground">
					<DiscAlbum size={16} />
					<span>Albums</span>
				</div>
				<span class="text-2xl font-bold">{syncState.numAlbums}</span>
				<span class="text-xs text-muted-foreground">
					Last sync: {formatDuration(syncState.albumsSyncDurationMs)}
				</span>
			</Card.Content>
		</Card.Root>

		<Card.Root>
			<Card.Content class="flex flex-col gap-2">
				<div class="flex items-center gap-2 text-sm text-muted-foreground">
					<FileMusic size={16} />
					<span>Tracks</span>
				</div>
				<span class="text-2xl font-bold">{syncState.numTracks}</span>
				<span class="text-xs text-muted-foreground">
					Last sync: {formatDuration(syncState.tracksSyncDurationMs)}
				</span>
			</Card.Content>
		</Card.Root>
	</div>

	<Card.Root>
		<Card.Content>
			<div class="flex items-center gap-2">
				<DiscAlbum size={18} />
				<h2 class="text-lg font-semibold">Missing Items</h2>
			</div>

			<Separator class="my-4" />

			<div class="flex flex-col gap-4">
				{#if syncState.missingArtists.length > 0}
					<div>
						<span class="text-xs font-medium text-muted-foreground"
							>Artists ({syncState.missingArtists.length})</span
						>
						<div class="mt-1 flex flex-col">
							{#each syncState.missingArtists as artist (artist.id)}
								<a href="/artists/{artist.id}" class="text-sm hover:underline"
									>{artist.name}</a
								>
							{/each}
						</div>
					</div>
				{/if}

				{#if syncState.missingAlbums.length > 0}
					<div>
						<span class="text-xs font-medium text-muted-foreground"
							>Albums ({syncState.missingAlbums.length})</span
						>
						<div class="mt-1 flex flex-col">
							{#each syncState.missingAlbums as album (album.id)}
								<a href="/albums/{album.id}" class="text-sm hover:underline"
									>{album.name}</a
								>
							{/each}
						</div>
					</div>
				{/if}

				{#if syncState.missingTracks.length > 0}
					<div>
						<span class="text-xs font-medium text-muted-foreground"
							>Tracks ({syncState.missingTracks.length})</span
						>
						<div class="mt-1 flex flex-col">
							{#each syncState.missingTracks as track (track.name)}
								<span class="text-sm">{track.name}</span>
							{/each}
						</div>
					</div>
				{/if}

				{#if syncState.missingArtists.length === 0 && syncState.missingAlbums.length === 0 && syncState.missingTracks.length === 0}
					<p class="text-sm text-muted-foreground">No missing items.</p>
				{/if}
			</div>
		</Card.Content>
	</Card.Root>

	{#if syncState.errors.length > 0}
		<Card.Root>
			<Card.Content>
				<div class="flex items-center gap-2 text-destructive">
					<AlertCircle size={18} />
					<h2 class="text-lg font-semibold">Errors</h2>
				</div>

				<Separator class="my-4" />

				<div class="flex flex-col gap-1">
					{#each syncState.errors as err (err)}
						<p class="font-mono text-sm text-destructive">{err}</p>
					{/each}
				</div>
			</Card.Content>
		</Card.Root>
	{/if}
</div>
