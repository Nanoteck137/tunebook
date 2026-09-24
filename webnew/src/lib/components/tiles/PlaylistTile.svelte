<script lang="ts">
	import { ChevronDown, FileHeart, Play, Shuffle, Star } from "@lucide/svelte";
	import { getMusicManager } from "$lib/music-manager.svelte";
	import type { Playlist } from "$lib/api/types";
	import { Button, Checkbox, DropdownMenu } from "$lib/components/ui";
	import { formatPlayTime } from "$lib/utils";
	import Tile from "./Tile.svelte";

	type Props = {
		playlist: Playlist;
		size?: "default" | "sm";
		class?: string;

		// Reorder selection mode.
		selectionMode?: boolean;
		selected?: boolean;
		onSelectChange?: (checked: boolean) => void;
		onMoveAfter?: () => void;

		// Quick playlist (starred) tile.
		isQuick?: boolean;
		onToggleQuick?: () => void;
	};

	const {
		playlist,
		size = "default",
		class: className,
		selectionMode = false,
		selected = false,
		onSelectChange,
		onMoveAfter,
		isQuick = false,
		onToggleQuick,
	}: Props = $props();
	const musicManager = getMusicManager();

	async function play(shuffle: boolean = false) {
		await musicManager.queueRequest(
			{ type: "addPlaylist", playlistId: playlist.id },
			{ shuffle },
		);
	}
</script>

{#snippet subtitle()}
	<p class="truncate text-xs text-muted-foreground">
		{playlist.trackCount} track{playlist.trackCount !== 1 ? "s" : ""}
		{#if playlist.playTime > 0}
			&middot; {formatPlayTime(playlist.playTime)}
		{/if}
	</p>
{/snippet}

{#snippet overlay()}
	{#if selectionMode}
		<div class="absolute top-1.5 left-1.5 z-10">
			<Checkbox
				checked={selected}
				onCheckedChange={(checked) => onSelectChange?.(checked)}
			/>
		</div>

		{#if onMoveAfter}
			<Button
				class="absolute right-2 bottom-2 z-10 hidden h-10 w-10 items-center justify-center rounded-full opacity-0 shadow-lg transition-all group-hover:scale-105 group-hover:opacity-100 hover:scale-110 sm:flex"
				variant="default"
				onclick={onMoveAfter}
				title="Move selected after this playlist"
				aria-label={`Move selected after ${playlist.name}`}
			>
				<ChevronDown size={18} />
			</Button>
		{/if}
	{:else if onToggleQuick}
		<button
			class="absolute top-1.5 right-1.5 z-10 flex h-7 w-7 items-center justify-center rounded-full transition-all {isQuick
				? 'bg-primary text-primary-foreground shadow-md'
				: 'border bg-background/70 text-muted-foreground backdrop-blur-sm hover:scale-105 hover:text-foreground'}"
			title={isQuick ? "Unset as quick playlist" : "Set as quick playlist"}
			aria-label={isQuick
				? `Unset ${playlist.name} as quick playlist`
				: `Set ${playlist.name} as quick playlist`}
			onclick={onToggleQuick}
		>
			<Star size={14} class={isQuick ? "fill-current" : ""} />
		</button>
	{/if}
{/snippet}

{#snippet menu()}
	<DropdownMenu.Group>
		<DropdownMenu.Item onclick={() => play()}>
			<Play size={14} />
			Play
		</DropdownMenu.Item>
		<DropdownMenu.Item onclick={() => play(true)}>
			<Shuffle size={14} />
			Shuffle play
		</DropdownMenu.Item>
	</DropdownMenu.Group>

	<DropdownMenu.Separator />

	<DropdownMenu.Group>
		{#if onSelectChange}
			<DropdownMenu.Item onclick={() => onSelectChange(true)}>
				Select playlist
			</DropdownMenu.Item>
		{/if}

		{#if onToggleQuick}
			<DropdownMenu.Item onclick={onToggleQuick}>
				<FileHeart />
				{isQuick ? "Remove Quick Playlist" : "Set as Quick Playlist"}
			</DropdownMenu.Item>
		{/if}
	</DropdownMenu.Group>
{/snippet}

<Tile
	href="/playlists/{playlist.id}"
	src={playlist.coverArt.medium}
	alt={playlist.name}
	name={playlist.name}
	sectionClass="section-playlists"
	nameClass={size === "sm" ? "line-clamp-1" : "line-clamp-2"}
	onPlay={selectionMode ? undefined : () => play()}
	disabled={selectionMode}
	onDisabledClick={() => onSelectChange?.(!selected)}
	menu={selectionMode || (!onSelectChange && !onToggleQuick)
		? undefined
		: menu}
	{overlay}
	{subtitle}
	class={className}
/>
