<script lang="ts">
	import {
		ChevronDown,
		ChevronLeft,
		ChevronRight,
		FileHeart,
		ListChecks,
		Play,
		Shuffle,
		Star,
	} from "@lucide/svelte";
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
		<div class="absolute top-2 right-2 z-10">
			<Checkbox
				checked={selected}
				onCheckedChange={(checked) => onSelectChange?.(checked)}
			/>
		</div>

		{#if onMoveAfter}
			<Button
				class="absolute right-2 bottom-2 hidden translate-y-2 rounded-full opacity-0 shadow-lg transition-all duration-300 group-hover:translate-y-0 group-hover:scale-105 group-hover:opacity-100 hover:scale-110 sm:inline-flex"
				size="icon"
				onclick={onMoveAfter}
				title="Move selected after this playlist"
				aria-label={`Move selected after ${playlist.name}`}
			>
				<ChevronRight />
			</Button>
		{/if}
	{:else if onToggleQuick}
		<Button
			class="absolute top-2 right-2 rounded-full opacity-80 shadow-lg transition-all duration-300 group-hover:scale-105 group-hover:opacity-100 hover:scale-110"
			size="icon"
			variant="secondary"
			title={isQuick ? "Unset as quick playlist" : "Set as quick playlist"}
			aria-label={isQuick
				? `Unset ${playlist.name} as quick playlist`
				: `Set ${playlist.name} as quick playlist`}
			onclick={onToggleQuick}
		>
			<Star class={isQuick ? "fill-primary stroke-primary" : ""} />
		</Button>
	{/if}
{/snippet}

{#snippet menu()}
	{#if onSelectChange}
		<DropdownMenu.Group>
			<DropdownMenu.Item onSelect={() => onSelectChange(true)}>
				<ListChecks />
				Select playlist
			</DropdownMenu.Item>
		</DropdownMenu.Group>

		<DropdownMenu.Separator />
	{/if}

	<DropdownMenu.Group>
		<DropdownMenu.Item onSelect={() => play()}>
			<Play />
			Play
		</DropdownMenu.Item>
		<DropdownMenu.Item onSelect={() => play(true)}>
			<Shuffle />
			Shuffle play
		</DropdownMenu.Item>
	</DropdownMenu.Group>

	{#if onToggleQuick}
		<DropdownMenu.Separator />

		<DropdownMenu.Group>
			<DropdownMenu.Item onSelect={onToggleQuick}>
				<FileHeart />
				{isQuick ? "Remove Quick Playlist" : "Set as Quick Playlist"}
			</DropdownMenu.Item>
		</DropdownMenu.Group>
	{/if}
{/snippet}

<Tile
	href="/playlists/{playlist.id}"
	src={playlist.coverArt.medium}
	alt={playlist.name}
	name={playlist.name}
	hover="group"
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
