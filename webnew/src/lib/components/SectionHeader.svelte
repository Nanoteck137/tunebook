<script lang="ts">
	import type { Snippet } from "svelte";
	import { ChevronRight } from "@lucide/svelte";

	let {
		children,
		count,
		viewAllHref,
	}: {
		children: Snippet;
		count?: number;
		viewAllHref?: string;
	} = $props();
</script>

<div class="mb-3 flex items-baseline justify-between px-2">
	<div class="flex items-baseline gap-2">
		<h2
			class="text-lg font-bold [&_svg:not([class*='size-'])]:size-4 [&_svg]:mr-1.5 [&_svg]:inline [&_svg]:shrink-0 [&_svg]:text-muted-foreground [&_svg]:align-[-0.125em]"
		>
			{@render children()}
		</h2>
		{#if count !== undefined}
			<span class="text-sm text-muted-foreground">{count}</span>
		{/if}
	</div>
	{#if viewAllHref}
		<a
			href={viewAllHref}
			class="flex items-center gap-0.5 text-sm font-medium text-primary hover:underline"
		>
			View all
			<ChevronRight class="h-4 w-4" />
		</a>
	{/if}
</div>
