<script lang="ts">
	import type { Snippet } from "svelte";
	import { cn } from "$lib/utils";

	let {
		// "gradient": 1px gradient hairline around the card.
		// "outline": standard border on the card instead.
		variant = "gradient",

		// Extra classes on the outer wrapper (root element). Use this to scope
		// the section colors, e.g. class="section-tracks".
		class: wrapperClass = "",

		// Extra classes on the inner card (layout, padding, responsive
		// variants such as "md:flex-row md:items-end md:gap-8").
		innerClass = "",

		// Inner card reference (scroll detection, ...).
		ref = $bindable(),

		children,
	}: {
		variant?: "gradient" | "outline";
		class?: string;
		innerClass?: string;
		ref?: HTMLElement | null;
		children: Snippet;
	} = $props();
</script>

<div
	bind:this={ref}
	class={cn(
		"rounded-lg",
		variant === "gradient" &&
			"bg-linear-to-b from-section-hero-to to-section-hero-from p-px",
		wrapperClass,
	)}
>
	<section
		class={cn(
			"flex h-full w-full flex-col gap-4 rounded-lg bg-linear-to-b from-section-hero-from to-section-hero-to p-4 shadow-sm sm:p-6",
			variant === "outline" && "border",
			innerClass,
		)}
	>
		{@render children()}
	</section>
</div>
