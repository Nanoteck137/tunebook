<script lang="ts">
	import { Info, Play, Shuffle, User } from "@lucide/svelte";
	import { getMusicManager } from "$lib/music-manager.svelte";
	import SectionImage from "$lib/components/SectionImage.svelte";
	import type { Album } from "$lib/api/types";
	import { Button, Dialog, DropdownMenu } from "$lib/components/ui";
	import { formatDate } from "$lib/utils";
	import Tile from "./Tile.svelte";

	type Props = {
		album: Album;
		size?: "default" | "sm";
		class?: string;
	};

	const { album, size = "default", class: className }: Props = $props();
	const musicManager = getMusicManager();

	let openShowMore = $state(false);

	async function play(shuffle: boolean = false) {
		await musicManager.queueRequest(
			{ type: "addAlbum", albumId: album.id },
			{
				shuffle,
			},
		);
	}
</script>

{#snippet subtitle()}
	<p
		class="line-clamp-1 text-xs text-ellipsis text-muted-foreground"
		title={album.artists.map((a) => a.name).join(", ")}
	>
		{#if album.year && size === "default"}
			{album.year} &middot;
		{/if}
		{album.artists.map((a) => a.name).join(", ")}
	</p>
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
		<DropdownMenu.Sub>
			<DropdownMenu.SubTrigger>
				<User />
				Go to artist
			</DropdownMenu.SubTrigger>
			<DropdownMenu.SubContent>
				{#each album.artists as artist (artist.id)}
					<a
						href="/artists/{artist.id}"
						class="flex items-center gap-2 rounded-sm px-3 py-1.5 text-sm text-popover-foreground hover:bg-accent hover:text-accent-foreground"
					>
						{artist.name}
					</a>
				{/each}
			</DropdownMenu.SubContent>
		</DropdownMenu.Sub>
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
	href="/albums/{album.id}"
	src={size === "sm" ? album.coverArt.small : album.coverArt.medium}
	alt={album.name}
	name={album.name}
	sectionClass="section-albums"
	nameClass={size === "sm" ? "line-clamp-1" : "line-clamp-2"}
	onPlay={() => play()}
	menu={size === "default" ? menu : undefined}
	{subtitle}
	class={className}
/>

{#if size === "default"}
	<Dialog.Root bind:open={openShowMore}>
		<Dialog.Content class="max-w-md gap-0 overflow-hidden p-0">
			<div class="relative overflow-hidden">
				<img
					src={album.coverArt.original}
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
						src={album.coverArt.large}
						alt={album.name}
						class="section-albums aspect-square w-20 shrink-0 rounded-md p-0.5 shadow-lg"
					/>
					<div class="min-w-0 flex-1 pb-0.5">
						<p
							class="text-lg leading-tight font-bold text-ellipsis text-white"
						>
							{album.name}
						</p>
						<p class="text-sm text-ellipsis text-white/80">
							{#each album.artists as artist, i (artist.id)}
								{#if i > 0}
									{", "}
								{/if}
								<a
									href="/artists/{artist.id}"
									class="hover:underline"
									title={artist.name}
								>
									{artist.name}
								</a>
							{/each}
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
						{#if album.year}
							<span class="text-muted-foreground">Year</span>
							<span class="font-medium">{album.year}</span>
						{/if}
						<span class="text-muted-foreground">Added</span>
						<span class="font-medium">{formatDate(album.created)}</span>
						<span class="text-muted-foreground">Updated</span>
						<span class="font-medium">{formatDate(album.updated)}</span>
					</div>
				</div>

				{#if album.tags.length > 0}
					<div class="flex flex-wrap gap-1.5">
						{#each album.tags as tag (tag)}
							<span class="rounded-full bg-secondary px-2.5 py-0.5 text-xs">
								{tag}
							</span>
						{/each}
					</div>
				{/if}
			</div>
		</Dialog.Content>
	</Dialog.Root>
{/if}
