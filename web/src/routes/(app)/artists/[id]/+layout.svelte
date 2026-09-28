<script lang="ts">
	import { page } from "$app/state";
	import {
		Breadcrumb,
		Button,
		buttonVariants,
		DropdownMenu,
	} from "$lib/components/ui";
	import {
		Disc,
		EllipsisVertical,
		LayoutDashboard,
		ListPlus,
		Music,
		Play,
		Shuffle,
	} from "@lucide/svelte";
	import { cn } from "$lib/utils";
	import HeroCard from "$lib/components/HeroCard.svelte";
	import SectionImage from "$lib/components/SectionImage.svelte";
	import { getMusicManager } from "$lib/music-manager.svelte";

	const { data, children } = $props();
	const musicManager = getMusicManager();

	const tabs = $derived([
		{
			label: "Overview",
			href: `/artists/${data.artist.id}`,
			icon: LayoutDashboard,
			match: `/artists/${data.artist.id}`,
			exact: true,
		},
		{
			label: "Songs",
			href: `/artists/${data.artist.id}/tracks`,
			icon: Music,
			match: `/artists/${data.artist.id}/tracks`,
		},
		{
			label: "Albums",
			href: `/artists/${data.artist.id}/albums`,
			icon: Disc,
			match: `/artists/${data.artist.id}/albums`,
		},
	]);

	let pathname = $derived(page.url.pathname);

	function isActive(tab: (typeof tabs)[number]) {
		if (tab.exact) {
			return pathname === tab.match;
		}

		return pathname.startsWith(tab.match);
	}
</script>

{#snippet heroActions()}
	<Button
		onclick={async () => {
			await musicManager.queueRequest(
				{ type: "addArtist", artistId: data.artist.id },
				{},
			);
		}}
	>
		<Play />
		Play
	</Button>

	<Button
		variant="ghost"
		size="icon"
		onclick={async () => {
			await musicManager.queueRequest(
				{ type: "addArtist", artistId: data.artist.id },
				{ shuffle: true },
			);
		}}
		title="Shuffle"
		aria-label="Shuffle"
	>
		<Shuffle />
	</Button>

	<DropdownMenu.Root>
		<DropdownMenu.Trigger
			class={buttonVariants({ variant: "ghost", size: "icon" })}
			title="More options"
			aria-label="More options"
		>
			<EllipsisVertical />
		</DropdownMenu.Trigger>
		<DropdownMenu.Content align="start">
			<DropdownMenu.Group>
				<DropdownMenu.Item
					onSelect={async () => {
						await musicManager.queueRequest(
							{ type: "addArtist", artistId: data.artist.id },
							{ append: "back" },
						);
					}}
				>
					<ListPlus />
					Append to Queue
				</DropdownMenu.Item>
			</DropdownMenu.Group>
		</DropdownMenu.Content>
	</DropdownMenu.Root>
{/snippet}

<div class="section-artists flex flex-col gap-4">
	<div class="py-2">
		<Breadcrumb.Root>
			<Breadcrumb.List>
				<Breadcrumb.Item>
					<Breadcrumb.Link href="/artists">Artists</Breadcrumb.Link>
				</Breadcrumb.Item>
				<Breadcrumb.Separator />
				<Breadcrumb.Item>
					<Breadcrumb.Page>{data.artist.name}</Breadcrumb.Page>
				</Breadcrumb.Item>
			</Breadcrumb.List>
		</Breadcrumb.Root>
	</div>

	<HeroCard>
		<div class="flex min-w-0 flex-col gap-2 md:hidden">
			<SectionImage
				class="aspect-square w-40 min-w-40 self-center rounded-xl shadow-2xl"
				src={data.artist.coverArt.large}
				alt={data.artist.name}
			/>

			<p
				class="text-xs font-semibold tracking-wider text-muted-foreground uppercase"
			>
				Artist
			</p>

			<h1 class="line-clamp-2 text-2xl font-bold">
				{data.artist.name}
			</h1>

			{#if data.artist.tags.length > 0}
				<div class="flex flex-wrap gap-1">
					{#each data.artist.tags as tag (tag)}
						<span
							class="rounded-full bg-secondary px-2.5 py-0.5 text-xs text-secondary-foreground"
						>
							{tag}
						</span>
					{/each}
				</div>
			{/if}

			<div class="flex gap-2 pt-2">
				{@render heroActions()}
			</div>
		</div>

		<div class="hidden gap-6 md:flex md:items-end md:gap-8">
			<SectionImage
				class="aspect-square w-40 min-w-40 self-center rounded-xl shadow-2xl md:w-52 md:min-w-52"
				src={data.artist.coverArt.large}
				alt={data.artist.name}
			/>

			<div class="flex min-w-0 flex-col gap-2">
				<p
					class="text-xs font-semibold tracking-wider text-muted-foreground uppercase"
				>
					Artist
				</p>

				<h1 class="line-clamp-2 text-4xl font-bold">
					{data.artist.name}
				</h1>

				{#if data.artist.tags.length > 0}
					<div class="flex flex-wrap gap-1">
						{#each data.artist.tags as tag (tag)}
							<span
								class="rounded-full bg-secondary px-2.5 py-0.5 text-xs text-secondary-foreground"
							>
								{tag}
							</span>
						{/each}
					</div>
				{/if}

				<div class="flex gap-2 pt-2">
					{@render heroActions()}
				</div>
			</div>
		</div>

		<div
			class="flex flex-wrap items-center gap-1 border-t border-border/40 pt-3"
		>
			{#each tabs as tab (tab.href)}
				<Button
					variant="ghost"
					size="sm"
					href={tab.href}
					class={cn(
						"text-foreground hover:bg-accent hover:text-accent-foreground",
						isActive(tab)
							? "bg-accent text-accent-foreground"
							: "text-muted-foreground",
					)}
				>
					<tab.icon />
					{tab.label}
				</Button>
			{/each}
		</div>
	</HeroCard>

	<div class="min-w-0 flex-1">
		{@render children()}
	</div>
</div>
