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
		Mic,
	} from "@lucide/svelte";
	import SectionImage from "$lib/components/SectionImage.svelte";
	import HeroCard from "$lib/components/HeroCard.svelte";
	import HeroIcon from "$lib/components/HeroIcon.svelte";
	import { getApiClient, handleApiError } from "$lib";
	import type { Artist } from "$lib/api/types";
	import InfiniteScroll from "$lib/components/InfiniteScroll.svelte";
	import { InfiniteScrollController } from "$lib/infinite-scroll.svelte";
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
	import { page } from "$app/state";
	import { getMusicManager } from "$lib/music-manager.svelte";
	import {
		sortTypes,
		defaultSort,
		type SortType,
		constructFilterSort,
	} from "./types";
	import { fly } from "svelte/transition";
	import ArtistTile from "$lib/components/tiles/ArtistTile.svelte";
	import Spacer from "$lib/components/Spacer.svelte";

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
		sort = defaultSort;
		includeTags = [];
		excludeTags = [];
		tagInput = "";

		const query = page.url.searchParams;
		query.delete("sort");
		query.delete("tags");
		query.delete("excludeTags");

		goto("?" + query.toString(), { invalidateAll: true });
	}

	let hasActiveFilters = $derived(
		includeTags.length > 0 || excludeTags.length > 0,
	);

	let filterOpen = $state(false);

	const apiClient = getApiClient();

	const scroll = new InfiniteScrollController<Artist>({
		initialLoad: () => ({
			items: data.artists,
			hasMore: data.page.page + 1 < data.page.totalPages,
			page: data.page.page,
		}),
		load: async (nextPage) => {
			const query: Record<string, string> = {
				page: String(nextPage),
				perPage: String(data.page.perPage),
			};

			constructFilterSort(data.filter, query);

			const res = await apiClient.getArtists({ query });
			if (!res.success) {
				handleApiError(res.error);
				return null;
			}

			return {
				items: res.data.artists,
				hasMore: res.data.page.page + 1 < res.data.page.totalPages,
			};
		},
		itemKey: (artist) => artist.id,
	});
</script>

<div class="flex flex-col gap-4">
	<HeroCard
		class="section-artists"
		innerClass="sm:flex-row sm:items-center sm:justify-between"
	>
		<div class="flex items-center gap-4">
			<HeroIcon>
				<Mic />
			</HeroIcon>

			<div class="flex min-w-0 flex-col">
				<h1 class="text-2xl font-bold">Artists</h1>
				<p class="text-sm text-muted-foreground">
					Every artist in your library
					{#if data.page}
						&middot; {data.page.totalItems}
					{/if}
				</p>
			</div>
		</div>
	</HeroCard>

	<!-- Toolbar -->
	<div class="flex flex-wrap items-center justify-between gap-2">
		<div class="relative flex-1 md:max-w-64">
			<Input
				class="pr-8"
				placeholder="Search artists..."
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
			<Button
				variant="ghost"
				size="icon"
				href="/search/artists"
				title="Advanced search"
			>
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

	{#if filterOpen}
		<Separator />

		<div
			transition:fly={{ y: -6, duration: 150 }}
			class="flex flex-col gap-3 px-2"
		>
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
	{/if}
</div>

<Spacer size="lg" />

<InfiniteScroll controller={scroll}>
	<div
		class="grid grid-cols-2 gap-3 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-5 xl:grid-cols-6 2xl:grid-cols-7"
	>
		{#each scroll.items as artist (artist.id)}
			<ArtistTile {artist} />
		{/each}
	</div>
</InfiniteScroll>
