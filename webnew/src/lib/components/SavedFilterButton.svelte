<script lang="ts">
	import { Button } from "$lib/components/ui";
	import { cn } from "$lib/utils";
	import { ListFilter } from "@lucide/svelte";

	let {
		class: className = "",

		filterOpen = $bindable(),
		hasFilters = false,
	}: {
		class?: string;
		filterOpen: boolean;
		hasFilters?: boolean;
	} = $props();
</script>

<Button
	class={cn(
		"relative transition-opacity",
		filterOpen && "bg-accent text-accent-foreground",
		className,
	)}
	variant="ghost"
	size="icon"
	title="Saved Filters"
	aria-label="Saved Filters"
	onclick={() => (filterOpen = !filterOpen)}
>
	<ListFilter />
	{#if hasFilters}
		<span class="absolute top-1.5 right-1.5 h-2 w-2 rounded-full bg-primary"
		></span>
	{/if}
</Button>
