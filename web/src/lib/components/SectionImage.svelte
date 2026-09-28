<script lang="ts" module>
	import { type VariantProps, tv } from "tailwind-variants";

	export const sectionImageVariants = tv({
		slots: {
			wrapper: "",
			img: "w-full h-full object-cover",
		},
		variants: {
			variant: {
				tile: {
					wrapper:
						"rounded-lg bg-linear-to-tl from-section-gradiant-1 to-section-hero-to p-1",
					img: "rounded-lg",
				},
				full: {
					wrapper:
						"rounded-lg bg-linear-to-tl from-section-gradiant-1 via-section-gradiant-2 to-section-gradiant-3 p-1",
					img: "rounded-lg",
				},
				// profile: {
				// 	wrapper:
				// 		"rounded-full bg-linear-to-tr from-section-gradiant-1 via-section-gradiant-2 to-section-gradiant-3 p-0.5",
				// 	img: "rounded-full object-cover",
				// },
			},
			// "off": no hover effect.
			// "normal": the wrapper scales on hover, the image inside scales
			// when its ancestor `group` is hovered.
			// "group": same as "normal" but the wrapper also scales on
			// `group-hover` instead of `hover`.
			hover: {
				off: {
					wrapper: "",
					img: "",
				},
				normal: {
					wrapper: "",
					img: "transition-transform duration-300 hover:scale-105",
				},
				group: {
					wrapper:
						"transition-transform duration-300 group-hover:scale-[1.02]",
					img: "transition-transform duration-300 group-hover:scale-105",
				},
			},
		},
		defaultVariants: {
			variant: "tile",
			hover: "normal",
		},
	});

	export type SectionImageVariant = VariantProps<
		typeof sectionImageVariants
	>["variant"];
	export type SectionImageHover = VariantProps<
		typeof sectionImageVariants
	>["hover"];
</script>

<script lang="ts">
	import type { HTMLImgAttributes } from "svelte/elements";
	import {
		cn,
		type WithElementRef,
		type WithoutChildrenOrChild,
	} from "$lib/utils";

	type Props = WithoutChildrenOrChild<WithElementRef<HTMLImgAttributes>> & {
		variant?: SectionImageVariant;
		hover?: SectionImageHover;

		// Extra classes on the image (size, inner border, ...).
		imgClass?: string;
	};

	let {
		variant = "tile",
		hover = "normal",

		// Extra classes on the ring wrapper (root element). Also the place to
		// scope the section colors, e.g. class="section-albums".
		class: wrapperClass = "",

		imgClass = "",

		ref = $bindable(null),

		...restProps
	}: Props = $props();
</script>

<div
	class={cn(sectionImageVariants({ variant, hover }).wrapper(), wrapperClass)}
>
	<img
		bind:this={ref}
		{...restProps}
		class={cn(sectionImageVariants({ variant, hover }).img(), imgClass)}
	/>
</div>
