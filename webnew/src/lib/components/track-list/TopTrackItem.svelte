<script lang="ts">
	import type { Track } from "$lib/api/types";
	import ArtistList from "$lib/components/ArtistList.svelte";
	import Image from "$lib/components/Image.svelte";
	import { buttonVariants, DropdownMenu } from "$lib/components/ui";
	import { cn } from "$lib/utils";
	import { EllipsisVertical, Play } from "@lucide/svelte";
	import type { Snippet } from "svelte";

	type Props = {
		rank?: number;
		track: Track;
		playCount?: number;
		onPlayClicked?: () => void;
		menuItems?: Snippet;
	};

	const { rank, track, playCount, onPlayClicked, menuItems }: Props = $props();
</script>

<div
	class="group flex items-center gap-3 rounded-lg p-2 transition-colors hover:bg-accent hover:text-accent-foreground has-data-[state='open']:bg-accent has-data-[state='open']:text-accent-foreground"
>
	<button
		class="shrink-0"
		onclick={() => onPlayClicked?.()}
		aria-label="Play {track.name}"
	>
		{#if rank !== undefined}
			<div
				class="flex h-12 w-12 items-center justify-center overflow-hidden rounded-md"
			>
				<span class="text-sm font-medium tabular-nums group-hover:hidden">
					{rank}
				</span>
				<div class="hidden items-center justify-center group-hover:flex">
					<Play size={20} />
				</div>
			</div>
		{/if}
	</button>

	<Image class="h-12 w-12 rounded" src={track.coverArt.small} alt="" />

	<div class="flex min-w-0 flex-1 flex-col gap-0.5">
		<p class="truncate text-sm font-medium" title={track.name}>
			{track.name}
		</p>

		<ArtistList class="text-muted-foreground" artists={track.artists} />
	</div>

	{#if playCount !== undefined}
		<span
			class="hidden shrink-0 text-xs text-muted-foreground tabular-nums sm:block"
		>
			{playCount.toLocaleString()}
		</span>
	{/if}

	{#if menuItems}
		<DropdownMenu.Root>
			<DropdownMenu.Trigger
				class={cn(
					buttonVariants({ variant: "ghost", size: "icon-sm" }),
					"-mr-1 shrink-0 rounded-full text-muted-foreground",
				)}
				title="More options"
				aria-label="More options"
			>
				<EllipsisVertical />
			</DropdownMenu.Trigger>
			<DropdownMenu.Content align="end">
				<DropdownMenu.Group>
					{@render menuItems()}
				</DropdownMenu.Group>
			</DropdownMenu.Content>
		</DropdownMenu.Root>
	{/if}
</div>

