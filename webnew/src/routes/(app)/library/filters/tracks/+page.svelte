<script lang="ts">
	import { goto, invalidateAll } from "$app/navigation";
	import { page } from "$app/state";
	import { onMount } from "svelte";
	import { getApiClient, handleApiError } from "$lib";
	import type { TrackFilter } from "$lib/api/types";
	import {
		Button,
		buttonVariants,
		Checkbox,
		DropdownMenu,
	} from "$lib/components/ui";
	import ConfirmModal from "$lib/components/new-modals/ConfirmModal.svelte";
	import {
		ChevronDown,
		EllipsisVertical,
		FileMusic,
		Heart,
		ListChecks,
		ListFilter,
		Plus,
		SquarePen,
		Trash,
		X,
	} from "@lucide/svelte";
	import { toast } from "svelte-sonner";
	import { cn } from "$lib/utils";
	import NewFilterModal from "../../../tracks/NewFilterModal.svelte";
	import EditFilterModal from "../../../tracks/EditFilterModal.svelte";
	import SectionHeader from "$lib/components/SectionHeader.svelte";
	import DebouncedSearchInput from "$lib/components/DebouncedSearchInput.svelte";
	import { SortToggleDropdown } from "$lib/components/sort";
	import {
		sortTypes,
		defaultSort,
		type SortType,
		searchFilters,
		sortFilters,
	} from "./types";

	let { data } = $props();
	const apiClient = getApiClient();

	let openNewFilterModal = $state(false);

	let editingFilter = $state<TrackFilter | null>(null);
	let openEditFilterModal = $state(false);

	let deletingFilter = $state<TrackFilter | null>(null);
	let openConfirmDelete = $state(false);

	let selectedFilters = $state<string[]>([]);

	let value = $state("");

	onMount(() => {
		value = page.url.searchParams.get("query") ?? "";
	});

	let sort = $state(
		(page.url.searchParams.get("sort") as SortType) ?? defaultSort,
	);

	let filters = $derived(
		sortFilters(searchFilters(data.filters, value), sort),
	);

	function startEdit(filter: TrackFilter) {
		editingFilter = filter;
		openEditFilterModal = true;
	}

	function toggleFilterSelection(filterId: string) {
		if (selectedFilters.includes(filterId)) {
			selectedFilters = selectedFilters.filter((id) => id !== filterId);
		} else {
			selectedFilters = [...selectedFilters, filterId];
		}
	}

	function startDelete(filter: TrackFilter) {
		deletingFilter = filter;
		openConfirmDelete = true;
	}

	async function deleteFilter() {
		if (!deletingFilter) return;

		const res = await apiClient.deleteTrackFilter(deletingFilter.filterId);
		if (!res.success) {
			handleApiError(res.error);
			invalidateAll();
			return;
		}

		toast.success("Filter deleted");
		invalidateAll();
	}

	function updateSort(sortValue: string) {
		sort = sortValue as SortType;

		selectedFilters = [];

		const query = page.url.searchParams;
		query.delete("sort");

		if (sort !== defaultSort) {
			query.set("sort", sort);
		}

		goto("?" + query.toString(), { invalidateAll: true });
	}

	async function search(query: string) {
		const params = page.url.searchParams;
		params.delete("query");

		if (query) {
			params.set("query", query);
		}

		await goto("?" + params.toString(), {
			invalidateAll: true,
			keepFocus: true,
			replaceState: true,
		});
	}

	async function handleReorder(anchorFilterId: string | null) {
		const res = await apiClient.reorderTrackFilters({
			before: false,
			anchorFilterId: anchorFilterId ?? "",
			filterIds: selectedFilters,
		});
		if (!res.success) {
			return handleApiError(res.error);
		}

		selectedFilters = [];
		toast.success("Updated filters");

		if (sort !== "position") {
			updateSort("position");
		} else {
			invalidateAll();
		}
	}
</script>

<div class="flex flex-col gap-4">
	<SectionHeader count={data.filters.length}>
		<ListFilter />
		Track Filters

		{#snippet actions()}
			<Button size="sm" onclick={() => (openNewFilterModal = true)}>
				<Plus size={14} />
				New Filter
			</Button>
		{/snippet}
	</SectionHeader>

	<!-- Toolbar -->
	<div class="flex flex-wrap items-center justify-between gap-2">
		<DebouncedSearchInput
			class="flex-1 md:max-w-64"
			placeholder="Search filters..."
			{value}
			setValue={(v) => (value = v)}
			{search}
		/>

		<div class="flex items-center gap-1">
			<SortToggleDropdown
				types={sortTypes}
				{sort}
				{defaultSort}
				onSortChange={(value) => updateSort(value)}
			/>
		</div>
	</div>

	{#if selectedFilters.length > 0}
		<div
			class="flex h-12 shrink-0 items-center gap-1 rounded-lg border border-dashed border-border/60 px-1 text-xs font-medium text-muted-foreground"
		>
			<button
				class="rounded-full p-1.5 transition-colors hover:bg-muted hover:text-foreground"
				onclick={() => {
					selectedFilters = [];
				}}
				aria-label="Clear selection"
				title="Clear selection"
			>
				<X size={16} />
			</button>

			<button
				class="group flex h-full min-w-0 flex-1 items-center justify-center gap-2 rounded-md border-2 border-dotted border-transparent transition-colors hover:border-primary/60 hover:bg-accent/50 hover:text-foreground"
				onclick={() => handleReorder(null)}
			>
				<ChevronDown
					class="text-muted-foreground/60 transition-colors group-hover:text-foreground"
				/>
				Place {selectedFilters.length} here
			</button>
		</div>
	{/if}

	{#if filters.length > 0}
		<div class="flex flex-col gap-2">
			{#each filters as filter (filter.filterId)}
				{@const selected = selectedFilters.includes(filter.filterId)}
				{@const selectionMode = selectedFilters.length > 0}
				<div
					class={cn(
						"flex items-center justify-between gap-2 rounded-lg border bg-card p-3 transition-colors",
						selectionMode && "cursor-pointer hover:bg-accent/50",
						selected && "border-primary/60 bg-accent",
					)}
					role="button"
					tabindex={selectionMode ? 0 : -1}
					aria-pressed={selectionMode ? selected : undefined}
					onclick={() => {
						if (selectionMode) {
							toggleFilterSelection(filter.filterId);
						}
					}}
					onkeydown={(e) => {
						if (selectionMode && (e.key === "Enter" || e.key === " ")) {
							e.preventDefault();
							toggleFilterSelection(filter.filterId);
						}
					}}
				>
					<div class="flex min-w-0 items-center gap-2">
						{#if selectionMode}
							<Checkbox
								checked={selected}
								onclick={(e) => e.stopPropagation()}
								onCheckedChange={(checked) => {
									if (checked) {
										selectedFilters = [...selectedFilters, filter.filterId];
									} else {
										selectedFilters = selectedFilters.filter(
											(id) => filter.filterId !== id,
										);
									}
								}}
							/>
						{/if}

						<div class="min-w-0">
							<p
								class="flex items-center gap-2 truncate text-sm font-medium"
								title={filter.name}
							>
								<ListFilter class="h-4 w-4 shrink-0 text-muted-foreground" />
								{filter.name}
							</p>
							<p
								class="truncate font-mono text-xs text-muted-foreground"
								title={filter.filter}
							>
								{filter.filter}
							</p>
						</div>
					</div>

					<div class="flex shrink-0 items-center gap-1">
						{#if selectionMode}
							<Button
								class="rounded-full"
								variant="ghost"
								size="icon"
								onclick={(e) => {
									e.stopPropagation();
									handleReorder(filter.filterId);
								}}
								title="Move selected after this filter"
								aria-label={`Move selected after ${filter.name}`}
							>
								<ChevronDown />
							</Button>
						{:else}
							<DropdownMenu.Root>
								<DropdownMenu.Trigger
									class={cn(
										buttonVariants({ variant: "ghost", size: "icon-sm" }),
										"-mr-1 shrink-0 rounded-full text-muted-foreground",
									)}
									aria-label={`More options for ${filter.name}`}
								>
									<EllipsisVertical size={14} />
								</DropdownMenu.Trigger>
								<DropdownMenu.Content align="end">
									<DropdownMenu.Group>
										<DropdownMenu.Item
											onSelect={() => {
												selectedFilters = [filter.filterId];
											}}
										>
											<ListChecks />
											Select filter
										</DropdownMenu.Item>
									</DropdownMenu.Group>

									<DropdownMenu.Separator />

									<DropdownMenu.Group>
										<DropdownMenu.Item
											onSelect={() =>
												goto(`/tracks?filterId=${filter.filterId}`)}
										>
											<FileMusic />
											View in Tracks
										</DropdownMenu.Item>
										<DropdownMenu.Item
											onSelect={() =>
												goto(`/library/favorites?filterId=${filter.filterId}`)}
										>
											<Heart />
											View in Favorites
										</DropdownMenu.Item>
									</DropdownMenu.Group>

									<DropdownMenu.Separator />

									<DropdownMenu.Group>
										<DropdownMenu.Item onSelect={() => startEdit(filter)}>
											<SquarePen />
											Edit Filter
										</DropdownMenu.Item>
										<DropdownMenu.Item onSelect={() => startDelete(filter)}>
											<Trash />
											Delete Filter
										</DropdownMenu.Item>
									</DropdownMenu.Group>
								</DropdownMenu.Content>
							</DropdownMenu.Root>
						{/if}
					</div>
				</div>
			{/each}
		</div>
	{:else}
		<p class="text-sm text-muted-foreground">
			{#if value}
				No filters match your search
			{:else}
				No filters yet
			{/if}
		</p>
	{/if}
</div>

<NewFilterModal bind:open={openNewFilterModal} />

{#if editingFilter}
	<EditFilterModal bind:open={openEditFilterModal} filter={editingFilter} />
{/if}

<ConfirmModal
	bind:open={openConfirmDelete}
	removeTrigger
	confirmDelete
	title="Delete Filter"
	description={deletingFilter
		? `Are you sure you want to delete "${deletingFilter.name}"?`
		: undefined}
	onResult={async () => {
		await deleteFilter();
	}}
/>
