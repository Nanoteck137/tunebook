<script lang="ts">
	import { goto } from "$app/navigation";
	import { page } from "$app/state";
	import { getApiClient, handleApiError } from "$lib";
	import type { Track } from "$lib/api/types";
	import {
		Button,
		Input,
		Select,
		Separator,
		buttonVariants,
	} from "$lib/components/ui";
	import { Music, Play, Shuffle, Plus, X, ListFilter } from "@lucide/svelte";
	import HeroIcon from "$lib/components/HeroIcon.svelte";
	import TrackList from "$lib/components/track-list/TrackList.svelte";
	import TrackVariants from "$lib/components/track-list/TrackVariants.svelte";
	import { getMusicManager } from "$lib/music-manager.svelte";
	import InfiniteScroll from "$lib/components/InfiniteScroll.svelte";
	import { InfiniteScrollController } from "$lib/infinite-scroll.svelte";
	import Spacer from "$lib/components/Spacer.svelte";
	import FilterButton from "./FilterButton.svelte";
	import { sortTypes, defaultSort, type SortType } from "./types";

	let { data } = $props();
	const musicManager = getMusicManager();
	const apiClient = getApiClient();

	let selectedSort = $state<SortType>(defaultSort);

	let tagInput = $state("");
	let tagMode = $state<"include" | "exclude">("include");
	let tags = $state<{ value: string; mode: "include" | "exclude" }[]>([]);

	function addTag() {
		const t = tagInput.trim();
		if (!t) return;

		if (!tags.some((x) => x.value === t && x.mode === tagMode)) {
			tags = [...tags, { value: t, mode: tagMode }];
		}

		tagInput = "";
	}

	function removeTag(value: string, mode: "include" | "exclude") {
		tags = tags.filter((t) => !(t.value === value && t.mode === mode));
	}

	let filterId = $derived(page.url.searchParams.get("filterId"));

	function clearFilter() {
		const query = page.url.searchParams;
		query.delete("filterId");
		goto("?" + query.toString(), {
			invalidateAll: true,
			replaceState: true,
		});
	}

	const scroll = new InfiniteScrollController<Track>({
		initialLoad: () => ({
			items: data.tracks,
			hasMore: data.page.page + 1 < data.page.totalPages,
			page: data.page.page,
		}),
		load: async (nextPage) => {
			const query: Record<string, string> = {
				page: String(nextPage),
				perPage: String(data.page.perPage),
			};

			if (filterId) {
				query["filterId"] = filterId;
			}

			const res = await apiClient.getTracks({ query });
			if (!res.success) {
				handleApiError(res.error);
				return null;
			}

			return {
				items: res.data.tracks,
				hasMore: res.data.page.page + 1 < res.data.page.totalPages,
			};
		},
		itemKey: (track) => track.id,
	});

	async function playTracks(options: { shuffle?: boolean } = {}) {
		if (filterId) {
			await musicManager.queueRequest(
				{ type: "addFilter", filterId },
				options,
			);
		} else {
			await musicManager.addTracks({
				trackIds: scroll.items.map((t) => t.id),
				clear: true,
			});
		}
	}
</script>

<div class="section-tracks flex flex-col gap-6">
	<section
		class="rounded-lg border bg-linear-to-b from-section-hero-from to-section-hero-to p-4 shadow-sm sm:p-6"
	>
		<div class="flex flex-col gap-4">
			<div class="flex items-center gap-4">
				<HeroIcon><Music /></HeroIcon>
				<div class="flex min-w-0 flex-col">
					<h1 class="text-2xl font-bold">Tracks</h1>
					<p class="text-sm text-muted-foreground">
						Every track in your library
						{#if data.page}
							&middot; {data.page.totalItems}
						{/if}
					</p>
				</div>
			</div>

			<div class="flex gap-2 pt-1">
				<Button size="sm" onclick={() => playTracks({ shuffle: true })}>
					<Shuffle size={14} />
					Shuffle
				</Button>
				<Button size="sm" onclick={() => playTracks()}>
					<Play size={14} />
					Play All
				</Button>
			</div>
		</div>
	</section>

	<div class="rounded-lg border bg-card p-3">
		<div
			class="flex flex-col gap-3 sm:flex-row sm:items-end sm:justify-between"
		>
			<div class="flex flex-1 flex-col gap-2 sm:flex-row sm:items-center">
				<Input class="sm:w-56" placeholder="Search tracks..." disabled />
				<Select.Root
					type="single"
					allowDeselect={false}
					onValueChange={(v) => (selectedSort = v as SortType)}
				>
					<Select.Trigger class="h-9 w-full sm:w-40">
						{sortTypes.find((i) => i.value === selectedSort)?.label ?? "Sort"}
					</Select.Trigger>
					<Select.Content>
						{#each sortTypes as ty (ty.value)}
							<Select.Item value={ty.value} label={ty.label} />
						{/each}
					</Select.Content>
				</Select.Root>
			</div>
		</div>

		<div class="mt-3 flex flex-wrap items-center gap-1.5">
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

			<button
				class="flex h-7 w-7 items-center justify-center rounded-md text-muted-foreground hover:text-foreground"
				onclick={addTag}
			>
				<Plus size={14} />
			</button>

			{#each tags as t (t.value + t.mode)}
				<span
					class="flex items-center gap-0.5 rounded-full px-2 py-0.5 text-xs {t.mode ===
					'include'
						? 'bg-primary/10 text-primary'
						: 'bg-destructive/10 text-destructive'}"
				>
					{t.mode === "include" ? "+" : "-"}{t.value}
					<button
						class="hover:text-inherit/80"
						onclick={() => removeTag(t.value, t.mode)}
					>
						<X size={11} />
					</button>
				</span>
			{/each}

			{#if tags.length > 0}
				<button
					class="text-xs text-muted-foreground hover:text-foreground"
					onclick={() => (tags = [])}
				>
					Clear
				</button>
			{/if}
		</div>
	</div>

	<div
		class="flex flex-wrap items-center justify-between gap-2 rounded-lg border bg-muted/40 px-3 py-2"
	>
		<div class="flex flex-wrap items-center gap-1.5">
			<span
				class="mr-1.5 flex items-center gap-1.5 text-xs font-medium text-muted-foreground"
			>
				<ListFilter size={12} />
				Saved Filters
			</span>

			{#if data.filters && data.filters.length > 0}
				{#each data.filters as filter (filter.filterId)}
					<FilterButton {filter} />
				{/each}
			{:else}
				<span class="text-sm text-muted-foreground">None saved yet</span>
			{/if}
		</div>

		<div class="flex items-center gap-1">
			<a
				href="/library/filters/tracks"
				class={buttonVariants({ variant: "ghost", size: "sm" })}
			>
				<ListFilter size={14} />
				Manage Filters
			</a>

			{#if page.url.searchParams.has("filterId")}
				<Button variant="ghost" size="sm" onclick={clearFilter}>
					<X size={14} />
					Clear
				</Button>
			{/if}
		</div>
	</div>
</div>

<Spacer size="md" />

<section class="rounded-lg border bg-card p-4">
	<h2 class="mb-2 text-lg font-bold">Design Variants</h2>
	<TrackVariants tracks={data.tracks} />
</section>

<Spacer size="lg" />

<Separator />

<Spacer size="md" />

<div class="flex items-baseline gap-2 px-2">
	<h2 class="text-lg font-bold">Current list</h2>
</div>

<Spacer size="md" />

<InfiniteScroll controller={scroll}>
	<TrackList
		totalTracks={data.page.totalItems}
		tracks={scroll.items}
		onPlay={async (trackId) => {
			if (filterId) {
				await musicManager.queueRequest(
					{ type: "addFilter", filterId },
					{ queueIndexToTrackId: trackId },
				);
			} else {
				await musicManager.addTracks({
					trackIds: scroll.items.map((t) => t.id),
					trackId,
					clear: true,
				});
			}
		}}
	/>
</InfiniteScroll>
