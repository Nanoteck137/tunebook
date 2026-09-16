<script lang="ts">
	import SeekSlider from "$lib/components/SeekSlider.svelte";
	import { Button, ScrollArea, Sheet } from "$lib/components/ui";
	import { formatTime } from "$lib/utils";
	import {
		ListX,
		ListMusic,
		Pause,
		Play,
		SkipBack,
		SkipForward,
		Volume2,
		VolumeX,
		X,
		Music2,
	} from "@lucide/svelte";
	import { getMusicManager, type MediaItem } from "$lib/music-manager.svelte";
	import Spinner from "$lib/components/Spinner.svelte";
	import Image from "$lib/components/Image.svelte";
	import FavoriteButton from "$lib/components/FavoriteButton.svelte";
	import QuickAddButton from "$lib/components/QuickAddButton.svelte";

	const musicManager = getMusicManager();

	const iconButton =
		"flex size-8 items-center justify-center rounded-full text-muted-foreground transition-colors hover:bg-accent hover:text-foreground";

	let currentMediaItem = $state<MediaItem | null>(null);
	let previousItems = $state<MediaItem[]>([]);
	let previousStartIndex = $state(0);
	let currentQueueItem = $state<MediaItem | null>(null);
	let nextItems = $state<MediaItem[]>([]);

	$effect(() => {
		currentMediaItem = musicManager.currentItem;
		currentQueueItem = musicManager.currentItem;

		musicManager.queue.getPreviousItems(0, 10).then((result) => {
			previousItems = result.items;
			previousStartIndex = result.startIndex;
		});
		musicManager.queue.getNextItems(0, 50).then((items) => {
			nextItems = items;
		});
	});
</script>

{#snippet queueSheet()}
	<Sheet.Root>
		<Sheet.Trigger class={iconButton} title="Queue" aria-label="Queue">
			<ListMusic size="20" />
		</Sheet.Trigger>
		<Sheet.Content side="right" class="gap-0" showCloseButton={false}>
			<Sheet.Title class="sr-only">Queue</Sheet.Title>
			<div
				class="flex shrink-0 items-center justify-between gap-2 border-b px-4 py-3"
			>
				<p class="text-sm font-semibold">Queue</p>
				<div class="flex items-center gap-1">
					<Button variant="ghost" size="sm" href="/queue">View all</Button>
					<Button
						class="rounded-full"
						variant="ghost"
						size="icon"
						title="Clear queue"
						aria-label="Clear queue"
						onclick={async () => {
							await musicManager.clearQueue();
						}}
					>
						<ListX />
					</Button>
					<Sheet.Close>
						{#snippet child({ props })}
							<Button
								class="rounded-full"
								variant="ghost"
								size="icon"
								title="Close queue"
								aria-label="Close queue"
								{...props}
							>
								<X size={18} />
							</Button>
						{/snippet}
					</Sheet.Close>
				</div>
			</div>

			<ScrollArea class="min-h-0 flex-1">
				<div class="flex flex-col gap-4 px-4 py-4 pr-3">
					<!-- Played -->
					{#if previousItems.length > 0}
						<div>
							<p
								class="mb-2 text-xs font-semibold tracking-wider text-muted-foreground uppercase"
							>
								Played
							</p>
							<div class="flex flex-col gap-1">
								{#each previousItems as mediaItem, i (mediaItem.trackId)}
									{@const queueIndex = previousStartIndex + i}
									<div
										class="group flex items-center gap-3 rounded-md p-2 transition-colors hover:bg-accent/50"
									>
										<button
											class="shrink-0"
											onclick={async () => {
												await musicManager.setQueueIndex(queueIndex);
												musicManager.play();
											}}
										>
											<Image
												class="w-10 min-w-10 rounded"
												src={mediaItem.coverArt}
												alt="cover"
											/>
										</button>
										<div class="flex min-w-0 flex-1 flex-col">
											<p
												class="truncate text-sm font-medium"
												title={mediaItem.name}
											>
												{mediaItem.name}
											</p>
											<p class="truncate text-xs text-muted-foreground">
												{#each mediaItem.artists as artist, j (artist.id)}
													{#if j > 0},
													{/if}{artist.name}
												{/each}
											</p>
										</div>
									</div>
								{/each}
							</div>
						</div>
					{/if}

					<!-- Now Playing -->
					{#if currentQueueItem}
						<div>
							<p
								class="mb-2 text-xs font-semibold tracking-wider text-muted-foreground uppercase"
							>
								Now Playing
							</p>
							<div class="flex items-center gap-3 rounded-md bg-accent/50 p-3">
								<div class="relative shrink-0">
									<Image
										class="w-12 min-w-12 rounded"
										src={currentQueueItem.coverArt}
										alt="cover"
									/>
									<div
										class="absolute inset-0 flex items-center justify-center"
									>
										<Music2 class="text-primary" size="16" />
									</div>
								</div>
								<div class="flex min-w-0 flex-col">
									<p
										class="truncate text-sm font-medium"
										title={currentQueueItem.name}
									>
										{currentQueueItem.name}
									</p>
									<p class="truncate text-xs text-muted-foreground">
										{#each currentQueueItem.artists as artist, j (artist.id)}
											{#if j > 0},
											{/if}{artist.name}
										{/each}
									</p>
								</div>
							</div>
						</div>
					{/if}

					<!-- Next in queue -->
					{#if nextItems.length > 0}
						<div>
							<p
								class="mb-2 text-xs font-semibold tracking-wider text-muted-foreground uppercase"
							>
								Next in queue
							</p>
							<div class="flex flex-col gap-1">
								{#each nextItems as mediaItem, i (mediaItem.trackId)}
									{@const queueIndex = musicManager.queue.index + 1 + i}
									<div
										class="group flex items-center gap-3 rounded-md p-2 transition-colors hover:bg-accent/50"
									>
										<span
											class="w-5 text-right text-xs text-muted-foreground tabular-nums"
											>{i + 1}</span
										>
										<button
											class="shrink-0"
											onclick={async () => {
												await musicManager.setQueueIndex(queueIndex);
												musicManager.play();
											}}
										>
											<Image
												class="w-10 min-w-10 rounded"
												src={mediaItem.coverArt}
												alt="cover"
											/>
										</button>
										<div class="flex min-w-0 flex-1 flex-col">
											<p
												class="truncate text-sm font-medium"
												title={mediaItem.name}
											>
												{mediaItem.name}
											</p>
											<p class="truncate text-xs text-muted-foreground">
												{#each mediaItem.artists as artist, j (artist.id)}
													{#if j > 0},
													{/if}{artist.name}
												{/each}
											</p>
										</div>
										<button
											class="shrink-0 rounded-full p-1 text-muted-foreground opacity-0 transition-opacity group-hover:opacity-100 hover:text-foreground"
											onclick={async () => {
												await musicManager.removeQueueItem(queueIndex);
											}}
										>
											<X size="14" />
										</button>
									</div>
								{/each}
							</div>
						</div>
					{/if}

					{#if !currentQueueItem}
						<p class="py-8 text-center text-sm text-muted-foreground">
							Queue is empty
						</p>
					{/if}
				</div>
			</ScrollArea>
		</Sheet.Content>
	</Sheet.Root>
{/snippet}

<div
	class="relative hidden h-20 border-t border-border/60 bg-background/95 backdrop-blur-md supports-backdrop-filter:bg-background/80 md:block"
>
	<div
		class="old-container absolute inset-x-0 -top-2 z-10 max-w-screen-2xl px-4 sm:px-8"
	>
		<SeekSlider
			value={Number.isNaN(musicManager.duration)
				? 0
				: musicManager.currentTime / musicManager.duration}
			onValue={(p) => {
				musicManager.setPosition(p * musicManager.duration);
			}}
			buffered={musicManager.buffered}
			tooltip
			formatTooltip={(p) => formatTime(p * musicManager.duration)}
		/>
	</div>

	<div
		class="old-container flex h-full max-w-screen-2xl items-center gap-6 px-4 sm:px-8"
	>
		<!-- Left: Now playing -->
		<div class="flex min-w-0 flex-1 basis-1/3 items-center gap-3">
			<a
				href={currentMediaItem ? `/albums/${currentMediaItem.album.id}` : "#"}
				class="shrink-0"
				title={currentMediaItem?.album.name ?? ""}
			>
				<Image
					class="size-14 min-w-14 rounded-lg border border-border/60 shadow-sm transition-transform duration-200 hover:scale-105"
					src={currentMediaItem?.coverArt}
					alt="cover"
					loading="eager"
				/>
			</a>

			<div class="flex min-w-0 flex-1 flex-col">
				<a
					href={currentMediaItem
						? `/albums/${currentMediaItem.album.id}`
						: "#"}
					class="truncate text-sm font-semibold hover:underline"
					title={currentMediaItem?.name ?? ""}
				>
					{currentMediaItem?.name ?? "No track playing"}
				</a>
				<p class="truncate text-xs text-muted-foreground">
					{#if currentMediaItem}
						{#each currentMediaItem.artists as artist, i (artist.id)}
							{#if i > 0},
							{/if}
							<a href="/artists/{artist.id}" class="hover:underline"
								>{artist.name}</a
							>
						{/each}
					{/if}
				</p>
			</div>

			{#if currentMediaItem}
				<div class="flex shrink-0 items-center gap-0.5">
					<FavoriteButton show trackId={currentMediaItem.trackId} />
					<QuickAddButton trackId={currentMediaItem.trackId} />
				</div>
			{/if}
		</div>

		<!-- Center: Transport -->
		<div class="flex shrink-0 flex-col items-center gap-1">
			<div class="flex items-center gap-3">
				<button
					class={iconButton}
					title="Previous track"
					aria-label="Previous track"
					onclick={() => musicManager.previousTrack()}
				>
					<SkipBack size="18" />
				</button>

				{#if musicManager.loading}
					<Spinner class="h-8 w-8" />
				{:else if musicManager.playing}
					<button
						class="flex h-9 w-9 items-center justify-center rounded-full bg-primary text-primary-foreground shadow-sm transition hover:scale-105 hover:bg-primary/90 active:scale-95"
						title="Pause"
						aria-label="Pause"
						onclick={() => musicManager.pause()}
					>
						<Pause size="18" />
					</button>
				{:else}
					<button
						class="flex h-9 w-9 items-center justify-center rounded-full bg-primary text-primary-foreground shadow-sm transition hover:scale-105 hover:bg-primary/90 active:scale-95"
						title="Play"
						aria-label="Play"
						onclick={() => musicManager.play()}
					>
						<Play size="18" />
					</button>
				{/if}

				<button
					class={iconButton}
					title="Next track"
					aria-label="Next track"
					onclick={() => musicManager.nextTrack()}
				>
					<SkipForward size="18" />
				</button>
			</div>

			<div
				class="text-[11px] leading-none tracking-tight text-muted-foreground tabular-nums"
			>
				<span class="inline-block min-w-8 text-right"
					>{formatTime(musicManager.currentTime)}</span
				>
				<span class="px-0.5 text-muted-foreground/60">/</span>
				<span class="inline-block min-w-8 text-left"
					>{formatTime(
						Number.isNaN(musicManager.duration) ? 0 : musicManager.duration,
					)}</span
				>
			</div>
		</div>

		<!-- Right: Volume + Queue -->
		<div class="flex flex-1 basis-1/3 items-center justify-end gap-2">
			<button
				class={iconButton}
				title={musicManager.muted || musicManager.volume === 0
					? "Unmute"
					: "Mute"}
				aria-label="Toggle mute"
				onclick={() => {
					musicManager.muted = !musicManager.muted;
				}}
			>
				{#if musicManager.muted || musicManager.volume === 0}
					<VolumeX size="18" />
				{:else}
					<Volume2 size="18" />
				{/if}
			</button>

			<SeekSlider
				class="w-28"
				growOnHover={false}
				tooltip
				value={musicManager.muted ? 0 : musicManager.volume}
				onValue={(p) => {
					musicManager.volume = p;
					musicManager.muted = p === 0;
				}}
			/>

			{@render queueSheet()}
		</div>
	</div>
</div>
