<script lang="ts" generics="T">
	import type { Snippet } from "svelte";
	import type { InfiniteScrollController } from "$lib/infinite-scroll.svelte";
	import { cn } from "$lib/utils";
	import Spinner from "./Spinner.svelte";
	import { Button } from "./ui";

	let {
		controller,
		rootMargin = "200px 0px 0px 0px",
		errorMessage = "Failed to load more items",
		className,
		children,
		loadingSnippet,
		errorSnippet,
	}: {
		controller: InfiniteScrollController<T>;
		rootMargin?: string;
		errorMessage?: string;
		className?: string;
		children?: Snippet;
		loadingSnippet?: Snippet;
		errorSnippet?: Snippet;
	} = $props();

	let sentinel = $state<HTMLElement | null>(null);

	$effect(() => {
		const el = sentinel;
		if (!el) return;

		const observer = new IntersectionObserver(
			([entry]) => {
				if (entry.isIntersecting) {
					controller.loadMore();
				}
			},
			{ rootMargin },
		);

		observer.observe(el);
		return () => observer.disconnect();
	});
</script>

<div class={cn("flex flex-col", className)}>
	{@render children?.()}

	<div bind:this={sentinel} class="w-full"></div>

	{#if controller.loading}
		{#if loadingSnippet}
			{@render loadingSnippet()}
		{:else}
			<div class="flex justify-center py-6">
				<Spinner />
			</div>
		{/if}
	{/if}

	{#if controller.error}
		{#if errorSnippet}
			{@render errorSnippet()}
		{:else}
			<div class="flex flex-col items-center gap-2 py-6">
				<p class="text-sm text-muted-foreground">{errorMessage}</p>
				<Button
					size="sm"
					variant="outline"
					onclick={() => controller.loadMore()}
				>
					Retry
				</Button>
			</div>
		{/if}
	{/if}
</div>
