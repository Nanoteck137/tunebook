<script lang="ts">
	import { formatDuration } from "$lib/utils.js";
	import { getApiClient, handleApiError } from "$lib";
	import {
		CircleAlert,
		Disc3,
		DiscAlbum,
		FileMusic,
		LoaderCircle,
		RefreshCw,
		Trash,
		Users,
	} from "@lucide/svelte";
	import { onMount } from "svelte";
	import { toast } from "svelte-sonner";
	import { AlertDialog, Button } from "$lib/components/ui";
	import SectionHeader from "$lib/components/SectionHeader.svelte";
	import type { LibrarySyncStateEventTy } from "../events";
	import {
		connectServerSse,
		JobSyncStateEvent,
		LibrarySyncStateEvent,
	} from "../events";
	import Spacer from "$lib/components/Spacer.svelte";

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

	let cleanupOpen = $state(false);

	const missingCount = $derived(
		syncState.missingArtists.length +
			syncState.missingAlbums.length +
			syncState.missingTracks.length,
	);

	async function runSync() {
		const res = await apiClient.runTask("library-sync");
		if (!res.success) {
			return handleApiError(res.error);
		}

		toast.success("Dispatched library sync");
	}

	async function runCleanup() {
		const res = await apiClient.runTask("library-cleanup");
		if (!res.success) {
			return handleApiError(res.error);
		}

		toast.success("Dispatched library cleanup");
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

<div class="flex flex-col gap-4">
	<section>
		<SectionHeader>
			<RefreshCw />
			Library sync

			{#snippet actions()}
				<span class="text-xs text-muted-foreground">
					Last sync: {formatDuration(syncState.totalSyncDurationMs)}
				</span>

				<Button onclick={runSync} disabled={syncing}>
					{#if syncing}
						<LoaderCircle size={14} class="animate-spin" />
						Syncing...
					{:else}
						Sync now
					{/if}
				</Button>
			{/snippet}
		</SectionHeader>

		<Spacer />

		<div class="flex flex-col gap-4 rounded-lg border bg-card p-4">
			{#if syncing && syncState.currentItem}
				<div class="flex items-center gap-3 rounded-lg border p-3">
					<LoaderCircle
						size={16}
						class="shrink-0 animate-spin text-muted-foreground"
					/>
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

			<div class="grid grid-cols-3 gap-3">
				<div class="flex flex-col gap-2 rounded-lg border bg-card p-4">
					<div class="flex items-center gap-2 text-sm text-muted-foreground">
						<Users size={16} />
						<span>Artists</span>
					</div>
					<span class="text-2xl font-bold">{syncState.numArtists}</span>
					<span class="text-xs text-muted-foreground">
						{formatDuration(syncState.artistsSyncDurationMs)}
					</span>
				</div>

				<div class="flex flex-col gap-2 rounded-lg border bg-card p-4">
					<div class="flex items-center gap-2 text-sm text-muted-foreground">
						<DiscAlbum size={16} />
						<span>Albums</span>
					</div>
					<span class="text-2xl font-bold">{syncState.numAlbums}</span>
					<span class="text-xs text-muted-foreground">
						{formatDuration(syncState.albumsSyncDurationMs)}
					</span>
				</div>

				<div class="flex flex-col gap-2 rounded-lg border bg-card p-4">
					<div class="flex items-center gap-2 text-sm text-muted-foreground">
						<FileMusic size={16} />
						<span>Tracks</span>
					</div>
					<span class="text-2xl font-bold">{syncState.numTracks}</span>
					<span class="text-xs text-muted-foreground">
						{formatDuration(syncState.tracksSyncDurationMs)}
					</span>
				</div>
			</div>
		</div>
	</section>

	<section>
		<SectionHeader count={missingCount}>
			<Disc3 />
			Missing Items

			{#snippet actions()}
				<Button
					variant="outline"
					size="sm"
					onclick={() => (cleanupOpen = true)}
					disabled={missingCount === 0}
				>
					<Trash />
					Cleanup
				</Button>
			{/snippet}
		</SectionHeader>

		<Spacer />

		<div class="flex flex-col gap-4 rounded-lg border bg-card p-4">
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
	</section>

	{#if syncState.errors.length > 0}
		<section>
			<div class="mb-3 flex items-center px-2">
				<CircleAlert class="text-destructive" size={16} />

				<h2 class="ml-1.5 text-lg font-bold text-destructive">Errors</h2>
			</div>

			<div
				class="flex flex-col gap-1 rounded-lg border border-destructive/30 bg-card p-4"
			>
				{#each syncState.errors as err (err)}
					<p class="font-mono text-sm text-destructive">{err}</p>
				{/each}
			</div>
		</section>
	{/if}
</div>

<AlertDialog.Root bind:open={cleanupOpen}>
	<AlertDialog.Content>
		<AlertDialog.Header>
			<AlertDialog.Title>Clean up missing items?</AlertDialog.Title>
			<AlertDialog.Description>
				This permanently deletes the {missingCount} missing artists, albums and tracks
				from the database. This cannot be undone.
			</AlertDialog.Description>
		</AlertDialog.Header>
		<AlertDialog.Footer>
			<AlertDialog.Cancel>Cancel</AlertDialog.Cancel>
			<AlertDialog.Action onclick={runCleanup}>Clean up</AlertDialog.Action>
		</AlertDialog.Footer>
	</AlertDialog.Content>
</AlertDialog.Root>
