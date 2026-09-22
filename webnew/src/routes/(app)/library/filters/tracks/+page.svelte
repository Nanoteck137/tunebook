<script lang="ts">
	import { goto, invalidateAll } from "$app/navigation";
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
		Edit,
		ExternalLink,
		FileMusic,
		Heart,
		ListFilter,
		Plus,
		Trash,
		X,
	} from "@lucide/svelte";
	import { toast } from "svelte-sonner";
	import NewFilterModal from "../../../tracks/NewFilterModal.svelte";
	import EditFilterModal from "../../../tracks/EditFilterModal.svelte";

	let { data } = $props();
	const apiClient = getApiClient();

	let openNewFilterModal = $state(false);

	let editingFilter = $state<TrackFilter | null>(null);
	let openEditFilterModal = $state(false);

	let deletingFilter = $state<TrackFilter | null>(null);
	let openConfirmDelete = $state(false);

	let selectedFilters = $state<string[]>([]);

	function startEdit(filter: TrackFilter) {
		editingFilter = filter;
		openEditFilterModal = true;
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
		invalidateAll();
	}
</script>

<div class="flex flex-col gap-4">
	<div class="flex items-baseline justify-between gap-2">
		<div class="flex items-baseline gap-2">
			<h1 class="text-xl font-bold">Track Filters</h1>
			<span class="text-sm text-muted-foreground">{data.filters.length}</span>
		</div>

		<Button variant="ghost" onclick={() => (openNewFilterModal = true)}>
			<Plus />
			New Filter
		</Button>
	</div>

	{#if selectedFilters.length > 0}
		<div class="flex items-center gap-2">
			<Button
				class="rounded-full"
				variant="ghost"
				size="icon-lg"
				onclick={() => {
					selectedFilters = [];
				}}
			>
				<X />
			</Button>

			<Button
				class="rounded-full"
				variant="default"
				onclick={() => handleReorder(null)}
			>
				<ChevronDown />
				Insert after
			</Button>
		</div>
	{/if}

	{#if data.filters.length > 0}
		<div class="flex flex-col gap-2">
			{#each data.filters as filter (filter.filterId)}
				<div
					class="flex items-center justify-between gap-2 rounded-lg border bg-card p-3"
				>
					<div class="flex min-w-0 items-center gap-2">
						{#if selectedFilters.length > 0}
							<Checkbox
								checked={selectedFilters.includes(filter.filterId)}
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
						{#if selectedFilters.length > 0}
							<Button
								class="rounded-full"
								variant="ghost"
								size="icon"
								onclick={() => handleReorder(filter.filterId)}
								title="Move selected after this filter"
								aria-label={`Move selected after ${filter.name}`}
							>
								<ChevronDown />
							</Button>
						{:else}
							<DropdownMenu.Root>
								<DropdownMenu.Trigger
									class={buttonVariants({ variant: "ghost", size: "icon" })}
									title="Filter options"
									aria-label="Filter options"
								>
									<ExternalLink />
								</DropdownMenu.Trigger>
								<DropdownMenu.Content align="end">
									<DropdownMenu.Group>
										<DropdownMenu.Item
											onSelect={() => {
												selectedFilters = [filter.filterId];
											}}
										>
											Select filter
										</DropdownMenu.Item>
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
								</DropdownMenu.Content>
							</DropdownMenu.Root>
							<Button
								variant="ghost"
								size="icon"
								title="Edit"
								aria-label="Edit"
								onclick={() => startEdit(filter)}
							>
								<Edit />
							</Button>
							<Button
								variant="ghost"
								size="icon"
								class="text-destructive hover:text-destructive"
								title="Delete"
								aria-label="Delete"
								onclick={() => startDelete(filter)}
							>
								<Trash />
							</Button>
						{/if}
					</div>
				</div>
			{/each}
		</div>
	{:else}
		<p class="text-sm text-muted-foreground">No filters yet</p>
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
