<script lang="ts">
	import type { Snippet } from "svelte";
	import { cn } from "$lib/utils";

	let {
		class: className = "",
		children,
	}: {
		class?: string;
		children: Snippet;
	} = $props();
</script>

<div
	class={cn(
		"grid grid-cols-2 gap-3 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-5 xl:grid-cols-6 2xl:grid-cols-7",
		className,
	)}
>
	{@render children()}
</div>
