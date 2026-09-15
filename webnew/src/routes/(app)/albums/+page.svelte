<script lang="ts">
	import {
		Search,
		X,
		Plus,
		EllipsisVertical,
		Info,
		Play,
		ListFilter,
		ListSortAscendingIcon,
		CheckIcon,
		Shuffle,
		User,
	} from "@lucide/svelte";
	import Spacer from "$lib/components/Spacer.svelte";
	import Pagination from "$lib/components/Pagination.svelte";
	import {
		Separator,
		Button,
		Input,
		Dialog,
		DropdownMenu,
		buttonVariants,
	} from "$lib/components/ui";
	import { cn } from "$lib/utils";
	import { goto } from "$app/navigation";
	import { getMusicManager } from "$lib/music-manager.svelte";
	import {
		sortTypes,
		decadeTypes,
		defaultSort,
		defaultDecade,
		type SortType,
		type DecadeType,
	} from "./types";
	import { page } from "$app/state";
	import { fly } from "svelte/transition";

	let { data } = $props();

	const musicManager = getMusicManager();

	let sort = $state(
		(page.url.searchParams.get("sort") as SortType) ?? defaultSort,
	);
	function updateSort(value: string) {
		sort = value as SortType;

		const query = page.url.searchParams;
		query.delete("sort");
		query.set("sort", sort);

		goto("?" + query.toString(), { invalidateAll: true });
	}

	let searchQuery = $state(page.url.searchParams.get("query") ?? "");
	function updateSearch() {
		const query = page.url.searchParams;
		query.delete("query");

		if (searchQuery) {
			query.set("query", searchQuery);
		}

		goto("?" + query.toString(), { invalidateAll: true });
	}

	function clearSearch() {
		searchQuery = "";
		updateSearch();
	}

	let decade = $state(
		(page.url.searchParams.get("decade") as DecadeType) ?? defaultDecade,
	);
	function updateDecade(value: string) {
		decade = value as DecadeType;

		const query = page.url.searchParams;
		query.delete("decade");

		if (decade !== "none") {
			query.set("decade", decade);
		}

		goto("?" + query.toString(), { invalidateAll: true });
	}

	let tagInput = $state("");
	let includeTags = $state(
		page.url.searchParams.get("tags")?.split(",").filter(Boolean) ?? [],
	);
	let excludeTags = $state(
		page.url.searchParams.get("excludeTags")?.split(",").filter(Boolean) ?? [],
	);
	let tagMode = $state<"include" | "exclude">("include");

	function addTag() {
		const tag = tagInput.trim();
		if (!tag) return;

		if (tagMode === "include") {
			if (includeTags.includes(tag)) return;
			includeTags = [...includeTags, tag];
		} else {
			if (excludeTags.includes(tag)) return;
			excludeTags = [...excludeTags, tag];
		}

		tagInput = "";
		applyTagFilters();
	}

	function removeIncludeTag(tag: string) {
		includeTags = includeTags.filter((t) => t !== tag);
		applyTagFilters();
	}

	function removeExcludeTag(tag: string) {
		excludeTags = excludeTags.filter((t) => t !== tag);
		applyTagFilters();
	}

	function applyTagFilters() {
		const query = page.url.searchParams;
		query.delete("tags");
		query.delete("excludeTags");

		if (includeTags.length > 0) {
			query.set("tags", includeTags.join(","));
		}

		if (excludeTags.length > 0) {
			query.set("excludeTags", excludeTags.join(","));
		}

		goto("?" + query.toString(), { invalidateAll: true });
	}

	function clearFilters() {
		sort = "name-a-z";
		decade = "none";
		includeTags = [];
		excludeTags = [];
		tagInput = "";

		goto("/albums", { invalidateAll: true });
	}

	let hasActiveFilters = $derived(
		decade !== "none" || includeTags.length > 0 || excludeTags.length > 0,
	);

	let filterOpen = $state(false);

	let infoAlbumId = $state<string | null>(null);
	let infoOpen = $state(false);

	let infoAlbum = $derived(
		infoAlbumId
			? (data.albums.find((a) => a.id === infoAlbumId) ?? null)
			: null,
	);

	function showInfo(id: string) {
		infoAlbumId = id;
		infoOpen = true;
	}

	function formatDate(iso: string) {
		return new Date(iso).toLocaleDateString(undefined, {
			year: "numeric",
			month: "short",
			day: "numeric",
		});
	}

	async function playAlbum(albumId: string, shuffle = false) {
		await musicManager.queueRequest(
			{ type: "addAlbum", albumId },
			{ shuffle },
		);
	}
</script>

<div class="flex flex-col gap-4">
	<div class="flex items-baseline gap-2 px-2">
		<h1 class="text-xl font-bold">Albums</h1>
		{#if data.page}
			<span class="text-sm text-muted-foreground">{data.page.totalItems}</span>
		{/if}
	</div>

	<!-- Toolbar -->
	<div class="flex flex-wrap items-center justify-between gap-2 px-2">
		<div class="relative flex-1 md:max-w-64">
			<Input
				class="pr-8"
				placeholder="Search albums..."
				bind:value={searchQuery}
				onkeydown={(e) => {
					if (e.key === "Enter") {
						updateSearch();
					}
				}}
			/>
			{#if searchQuery}
				<button
					class="absolute top-1/2 right-1.5 -translate-y-1/2 rounded-full p-0.5 text-muted-foreground hover:text-foreground"
					onclick={clearSearch}
					aria-label="Clear search"
				>
					<X size={14} />
				</button>
			{/if}
		</div>

		<div class="flex items-center gap-1">
			<Button variant="ghost" size="icon" href="/search/albums" title="Advanced search">
				<Search />
			</Button>

			<Button
				variant="ghost"
				size="icon"
				title="Filters"
				aria-label="Filters"
					class={cn(
						"relative transition-opacity",
						filterOpen && "bg-accent text-accent-foreground",
					)}
					onclick={() => (filterOpen = !filterOpen)}
				>
					<ListFilter />
					{#if hasActiveFilters}
						<span
							class="absolute top-1.5 right-1.5 h-2 w-2 rounded-full bg-primary"
						></span>
					{/if}
				</Button>

				<DropdownMenu.Root>
					<DropdownMenu.Trigger
						class={buttonVariants({ variant: "ghost", size: "icon" })}
						title="Sort"
						aria-label="Sort"
					>
						<ListSortAscendingIcon />
					</DropdownMenu.Trigger>
					<DropdownMenu.Content align="end">
						<DropdownMenu.Group>
							{#each sortTypes as ty (ty.value)}
								{@const selected = sort === ty.value}
								<DropdownMenu.Item
									onSelect={() => updateSort(ty.value)}
									class={selected ? "bg-accent text-foreground" : ""}
								>
									{#if selected}
										<CheckIcon />
									{/if}
									{ty.label}
								</DropdownMenu.Item>
							{/each}
						</DropdownMenu.Group>
					</DropdownMenu.Content>
				</DropdownMenu.Root>

				{#if hasActiveFilters}
					<Button variant="ghost" size="sm" onclick={clearFilters} class="pr-1">
						<X size={14} />
						Clear
					</Button>
				{/if}
		</div>
	</div>

	<Separator />

	{#if filterOpen}
		<div transition:fly={{ y: -6, duration: 150 }} class="flex flex-col gap-3 px-2">
			<div class="flex flex-wrap items-center gap-1.5">
				<span class="text-xs font-medium text-muted-foreground">Decade</span>
				{#each decadeTypes as d (d.value)}
					<button
						class="rounded-md border px-2 py-1 text-xs transition-colors {decade ===
						d.value
							? 'border-primary bg-primary text-primary-foreground'
							: 'bg-transparent text-muted-foreground hover:text-foreground'}"
						onclick={() => updateDecade(d.value)}
					>
						{d.label}
					</button>
				{/each}
			</div>

			<div class="flex flex-wrap items-center gap-1.5">
				<span class="text-xs font-medium text-muted-foreground">Tags</span>

				<div class="flex items-center gap-1">
					<button
						class="rounded-l-md border px-1.5 py-1 text-xs font-medium transition-colors {tagMode ===
						'include'
							? 'border-primary bg-primary text-primary-foreground'
							: 'bg-transparent text-muted-foreground hover:text-foreground'}"
						onclick={() => (tagMode = "include")}
					>
						+ Inc
					</button>
					<button
						class="-ml-px rounded-r-md border px-1.5 py-1 text-xs font-medium transition-colors {tagMode ===
						'exclude'
							? 'text-destructive-foreground border-destructive bg-destructive'
							: 'bg-transparent text-muted-foreground hover:text-foreground'}"
						onclick={() => (tagMode = "exclude")}
					>
						- Exc
					</button>
				</div>

				<Input
					class="h-7 w-28 text-xs"
					placeholder="Tag name..."
					bind:value={tagInput}
					onkeydown={(e) => {
						if (e.key === "Enter") {
							addTag();
						}
					}}
				/>

				<Button variant="ghost" size="icon" class="h-7 w-7" onclick={addTag}>
					<Plus size={14} />
				</Button>

				{#each includeTags as tag (tag)}
					<span
						class="flex items-center gap-0.5 rounded-full bg-primary/10 px-2 py-0.5 text-xs text-primary"
					>
						+{tag}
						<button
							class="hover:text-primary/80"
							onclick={() => removeIncludeTag(tag)}
						>
							<X size={11} />
						</button>
					</span>
				{/each}
				{#each excludeTags as tag (tag)}
					<span
						class="flex items-center gap-0.5 rounded-full bg-destructive/10 px-2 py-0.5 text-xs text-destructive"
					>
						-{tag}
						<button
							class="hover:text-destructive/80"
							onclick={() => removeExcludeTag(tag)}
						>
							<X size={11} />
						</button>
					</span>
				{/each}
			</div>
		</div>

		<Separator />
	{/if}
</div>

<Spacer size="lg" />

<div
	class="grid grid-cols-2 gap-3 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-5 xl:grid-cols-6 2xl:grid-cols-7"
>
	{#each data.albums as album (album.id)}
		<div class="group relative flex flex-col">
			<div class="relative">
				<a href="/albums/{album.id}" class="block overflow-hidden rounded-lg">
					<img
						src={album.coverArt.medium}
						alt={album.name}
						class="aspect-square w-full object-cover transition-transform duration-300 group-hover:scale-105"
					/>
				</a>

				<button
					class="absolute right-2 bottom-2 hidden h-10 w-10 translate-y-2 items-center justify-center rounded-full bg-primary text-primary-foreground opacity-0 shadow-lg transition-all duration-300 group-hover:translate-y-0 group-hover:scale-105 group-hover:opacity-100 hover:scale-110 sm:flex"
					title="Play album"
					aria-label={`Play ${album.name}`}
					onclick={() => playAlbum(album.id)}
				>
					<Play size={18} />
				</button>
			</div>

			<div class="flex flex-col gap-0.5 pt-2">
				<div class="flex items-center gap-1">
					<a
						href="/albums/{album.id}"
						class="min-w-0 flex-1 truncate text-sm font-medium hover:underline"
						title={album.name}
					>
						{album.name}
					</a>

					<DropdownMenu.Root>
						<DropdownMenu.Trigger
							class={cn(
								buttonVariants({ variant: "ghost", size: "icon-sm" }),
								"-mr-1 shrink-0 rounded-full text-muted-foreground",
							)}
							aria-label={`More options for ${album.name}`}
						>
							<EllipsisVertical size={14} />
						</DropdownMenu.Trigger>
						<DropdownMenu.Content align="center">
							<DropdownMenu.Group>
								<DropdownMenu.Item
									onclick={() => playAlbum(album.id)}
								>
									<Play size={14} />
									Play
								</DropdownMenu.Item>
								<DropdownMenu.Item
									onclick={() => playAlbum(album.id, true)}
								>
									<Shuffle size={14} />
									Shuffle play
								</DropdownMenu.Item>
							</DropdownMenu.Group>

							<DropdownMenu.Separator />

							<DropdownMenu.Group>
								<DropdownMenu.Sub>
									<DropdownMenu.SubTrigger>
										<User size={14} />
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
								<DropdownMenu.Item onclick={() => showInfo(album.id)}>
									<Info size={14} />
									Show more info
								</DropdownMenu.Item>
							</DropdownMenu.Group>
						</DropdownMenu.Content>
					</DropdownMenu.Root>
				</div>

				<p
					class="truncate text-xs text-muted-foreground"
					title={album.artists.map((a) => a.name).join(", ")}
				>
					{#if album.year}
						{album.year} &middot;
					{/if}
					{album.artists.map((a) => a.name).join(", ")}
				</p>
			</div>
		</div>
	{/each}
</div>

<Spacer size="lg" />

<Separator />

<Spacer size="lg" />

<Pagination page={data.page} />

<Dialog.Root open={infoOpen} onOpenChange={(v) => (infoOpen = v)}>
	<Dialog.Content class="max-w-md gap-0 overflow-hidden p-0">
		{#if infoAlbum}
			<div class="relative">
				<img
					src={infoAlbum.coverArt.original}
					alt=""
					aria-hidden="true"
					class="h-44 w-full object-cover"
				/>
				<div
					class="absolute inset-0 bg-gradient-to-t from-black/90 via-black/40 to-black/30"
				></div>
				<div class="absolute inset-x-0 bottom-0 flex items-end gap-4 p-4">
					<img
						src={infoAlbum.coverArt.large}
						alt={infoAlbum.name}
						class="h-20 w-20 shrink-0 rounded-md object-cover shadow-lg"
					/>
					<div class="min-w-0 flex-1 pb-0.5">
						<p class="truncate text-lg leading-tight font-bold text-white">
							{infoAlbum.name}
						</p>
						<p class="truncate text-sm text-white/80">
							{#each infoAlbum.artists as artist, i (artist.id)}
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
					<Button
						size="sm"
						class="flex-1"
						onclick={() => playAlbum(infoAlbum.id)}
					>
						<Play size={14} />
						Play
					</Button>
					<Button
						size="sm"
						variant="outline"
						class="flex-1"
						onclick={() => playAlbum(infoAlbum.id, true)}
					>
						<Shuffle size={14} />
						Shuffle
					</Button>
				</div>

				<div class="rounded-lg border bg-card p-3 text-sm">
					<div class="grid grid-cols-[auto_1fr] gap-x-4 gap-y-1.5">
						{#if infoAlbum.year}
							<span class="text-muted-foreground">Year</span>
							<span class="font-medium">{infoAlbum.year}</span>
						{/if}
						<span class="text-muted-foreground">Added</span>
						<span class="font-medium">{formatDate(infoAlbum.created)}</span>
						<span class="text-muted-foreground">Updated</span>
						<span class="font-medium">{formatDate(infoAlbum.updated)}</span>
					</div>
				</div>

				{#if infoAlbum.tags.length > 0}
					<div class="flex flex-wrap gap-1.5">
						{#each infoAlbum.tags as tag (tag)}
							<span class="rounded-full bg-secondary px-2.5 py-0.5 text-xs"
								>{tag}</span
							>
						{/each}
					</div>
				{/if}
			</div>
		{:else}
			<p class="p-4 text-sm text-muted-foreground">Album not found.</p>
		{/if}
	</Dialog.Content>
</Dialog.Root>
