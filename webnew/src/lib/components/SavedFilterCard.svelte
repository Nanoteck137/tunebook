<script lang="ts">
	import { fly } from "svelte/transition";
	import { Button, buttonVariants, Card } from "$lib/components/ui";
	import type { TrackFilter } from "$lib/api/types";
	import { cn } from "$lib/utils";
	import { page } from "$app/state";
	import { goto } from "$app/navigation";
	import { ExternalLink, ListFilter, X } from "@lucide/svelte";

	const {
		class: className = "",

		filterOpen,
		filters,
	}: {
		class?: string;

		filterOpen: boolean;
		filters?: TrackFilter[];
	} = $props();

	let activeFilterId = $derived(page.url.searchParams.get("filterId"));

	function selectFilter(filterId: string) {
		const query = page.url.searchParams;
		query.set("filterId", filterId);

		goto("?" + query.toString(), {
			invalidateAll: true,
			replaceState: true,
		});
	}

	function clearFilter() {
		const query = page.url.searchParams;
		query.delete("filterId");
		goto("?" + query.toString(), {
			invalidateAll: true,
			replaceState: true,
		});
	}
</script>

{#if filterOpen}
	<div class={cn(className)} transition:fly={{ y: -6, duration: 150 }}>
		<Card.Root class="py-2">
			<Card.Content class="flex-warp flex items-center justify-between px-2">
				<div class="flex flex-wrap items-center gap-1.5">
					{#if filters && filters.length > 0}
						{#each filters as filter (filter.filterId)}
							{@const active = activeFilterId === filter.filterId}
							<button
								class={cn(
									"flex items-center gap-1.5 rounded-full border py-1 pr-2.5 pl-2.5 text-sm font-medium transition-colors",
									active
										? "border-primary bg-primary text-primary-foreground hover:bg-primary/90"
										: "border-border bg-card text-foreground hover:bg-accent",
								)}
								onclick={() => selectFilter(filter.filterId)}
								title={active ? "Active filter" : "Use this filter"}
							>
								<ListFilter
									size={12}
									class={active
										? "shrink-0"
										: "shrink-0 text-muted-foreground"}
								/>
								<span class="max-w-32 truncate">{filter.name}</span>
							</button>
						{/each}
					{:else}
						<span class="text-sm text-muted-foreground">None saved yet</span>
					{/if}
				</div>

				<div class="flex items-center gap-1">
					<a
						href="/library/filters/tracks"
						class={buttonVariants({ variant: "ghost", size: "icon" })}
						title="Manage Filters"
					>
						<ExternalLink />
					</a>

					{#if !!activeFilterId}
						<Button variant="ghost" size="icon" onclick={clearFilter}>
							<X />
						</Button>
					{/if}
				</div>
			</Card.Content>
		</Card.Root>
	</div>
{/if}
