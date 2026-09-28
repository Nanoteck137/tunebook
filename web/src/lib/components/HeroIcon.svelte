<script lang="ts" module>
	import { type VariantProps, tv } from "tailwind-variants";

	export const heroIconVariants = tv({
		base: "flex shrink-0 items-center justify-center rounded-xl bg-linear-to-tr from-section-gradiant-1 via-section-gradiant-2 to-section-gradiant-3 text-white [&_svg:not([class*='size-'])]:size-6",
		variants: {
			size: {
				default: "h-12 w-12",
				sm: "h-9 w-9",
				lg: "h-14 w-14",
			},
		},
		defaultVariants: {
			size: "default",
		},
	});

	export type HeroIconSize = VariantProps<typeof heroIconVariants>["size"];
</script>

<script lang="ts">
	import type { Snippet } from "svelte";
	import { cn } from "$lib/utils";

	let {
		size = "default",
		class: className = "",
		children,
	}: {
		size?: HeroIconSize;
		class?: string;
		children: Snippet;
	} = $props();
</script>

<div class={cn(heroIconVariants({ size }), className)}>
	{@render children()}
</div>
