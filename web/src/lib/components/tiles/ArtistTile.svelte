<script lang="ts">
	import { Info, Play, Shuffle } from "@lucide/svelte";
	import { getMusicManager } from "$lib/music-manager.svelte";
	import SectionImage from "$lib/components/SectionImage.svelte";
	import type { Artist } from "$lib/api/types";
	import { Button, Dialog, DropdownMenu } from "$lib/components/ui";
	import { formatDate } from "$lib/utils";
	import Tile from "./Tile.svelte";

	type Props = {
		artist: Artist;
	};

	const { artist }: Props = $props();
	const musicManager = getMusicManager();

	let openShowMore = $state(false);

	async function play(shuffle: boolean = false) {
		await musicManager.queueRequest(
			{ type: "addArtist", artistId: artist.id },
			{ shuffle },
		);
	}
</script>

{#snippet subtitle()}
	{#if artist.tags.length > 0}
		<p
			class="line-clamp-1 text-xs text-ellipsis text-muted-foreground"
			title={artist.tags.join(", ")}
		>
			{artist.tags.join(", ")}
		</p>
	{/if}
{/snippet}

{#snippet menu()}
	<DropdownMenu.Group>
		<DropdownMenu.Item onclick={() => play()}>
			<Play />
			Play
		</DropdownMenu.Item>
		<DropdownMenu.Item onclick={() => play(true)}>
			<Shuffle />
			Shuffle play
		</DropdownMenu.Item>
	</DropdownMenu.Group>

	<DropdownMenu.Separator />

	<DropdownMenu.Group>
		<DropdownMenu.Item onclick={() => (openShowMore = true)}>
			<Info />
			Show more info
		</DropdownMenu.Item>
	</DropdownMenu.Group>
{/snippet}

<Tile
	href="/artists/{artist.id}"
	src={artist.coverArt.medium}
	alt={artist.name}
	name={artist.name}
	sectionClass="section-artists"
	hover="group"
	onPlay={() => play()}
	{menu}
	{subtitle}
/>

<Dialog.Root bind:open={openShowMore}>
	<Dialog.Content class="max-w-md gap-0 overflow-hidden p-0">
		<div class="relative overflow-hidden">
			<img
				src={artist.coverArt.original}
				alt=""
				aria-hidden="true"
				class="h-44 w-full scale-105 object-cover blur-sm"
			/>
			<div
				class="absolute inset-0 bg-linear-to-t from-black/90 via-black/40 to-black/30"
			></div>
			<div class="absolute inset-x-0 bottom-0 flex items-end gap-4 p-4">
				<SectionImage
					variant="full"
					src={artist.coverArt.large}
					alt={artist.name}
					class="section-artists aspect-square w-20 shrink-0 rounded-md p-0.5 shadow-lg"
				/>
				<div class="min-w-0 flex-1 pb-0.5">
					<p class="text-lg leading-tight font-bold text-ellipsis text-white">
						{artist.name}
					</p>
				</div>
			</div>
		</div>

		<div class="flex flex-col gap-4 p-4 pt-3">
			<div class="flex items-center gap-2">
				<Button size="sm" class="flex-1" onclick={() => play()}>
					<Play />
					Play
				</Button>

				<Button
					size="sm"
					variant="outline"
					class="flex-1"
					onclick={() => play(true)}
				>
					<Shuffle />
					Shuffle
				</Button>
			</div>

			<div class="rounded-lg border bg-card p-3 text-sm">
				<div class="grid grid-cols-[auto_1fr] gap-x-4 gap-y-1.5">
					<span class="text-muted-foreground">Added</span>
					<span class="font-medium">{formatDate(artist.created)}</span>
					<span class="text-muted-foreground">Updated</span>
					<span class="font-medium">{formatDate(artist.updated)}</span>
				</div>
			</div>

			{#if artist.tags.length > 0}
				<div class="flex flex-wrap gap-1.5">
					{#each artist.tags as tag (tag)}
						<span class="rounded-full bg-secondary px-2.5 py-0.5 text-xs">
							{tag}
						</span>
					{/each}
				</div>
			{/if}
		</div>
	</Dialog.Content>
</Dialog.Root>
