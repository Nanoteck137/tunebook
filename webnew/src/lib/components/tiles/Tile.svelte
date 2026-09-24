<script lang="ts">
	import { EllipsisVertical, Play } from "@lucide/svelte";
	import type { Snippet } from "svelte";
	import SectionImage, {
		type SectionImageHover,
	} from "../SectionImage.svelte";
	import { Button, DropdownMenu, buttonVariants } from "$lib/components/ui";
	import { cn } from "$lib/utils";

	type Props = {
		href: string;
		src: string;
		alt?: string;
		name: string;

		// Scopes the section colors, e.g. "section-albums".
		sectionClass?: string;

		hover?: SectionImageHover;

		// Shows a play overlay button in the corner of the image.
		onPlay?: () => void;

		// Items for the ellipsis dropdown. Omit it and the button is hidden.
		menu?: Snippet;

		// Content rendered below the name.
		subtitle?: Snippet;

		// Content rendered over the image (checkboxes, quick-playlist toggle,
		// ...). Position it with absolute classes.
		overlay?: Snippet;

		// Extra classes on the name link (line-clamp, ...).
		nameClass?: string;

		// Blocks navigation on the image/name links (selection mode).
		disabled?: boolean;

		// Called instead of navigating when `disabled` (e.g. toggle selection).
		onDisabledClick?: () => void;

		class?: string;
	};

	let {
		href,
		src,
		alt = "",
		name,
		sectionClass,
		hover = "normal",
		onPlay,
		menu,
		subtitle,
		overlay,
		nameClass,
		disabled = false,
		onDisabledClick,
		class: className = "",
	}: Props = $props();
</script>

<div class={cn("group relative flex shrink-0 flex-col", className)}>
	<div class="relative">
		<a
			{href}
			class="block overflow-hidden rounded-lg"
			onclick={(e) => {
				if (disabled) {
					e.preventDefault();
					onDisabledClick?.();
				}
			}}
			aria-disabled={disabled || undefined}
		>
			<SectionImage
				variant="tile"
				hover={disabled ? "off" : hover}
				class={cn("aspect-square w-full", sectionClass)}
				{src}
				{alt}
				title={alt}
			/>
		</a>

		{#if onPlay}
			<Button
				class="absolute right-2 bottom-2 hidden translate-y-2 rounded-full opacity-0 shadow-lg transition-all duration-300 group-hover:translate-y-0 group-hover:scale-105 group-hover:opacity-100 hover:scale-110 sm:inline-flex"
				size="icon"
				onclick={onPlay}
			>
				<Play />
			</Button>
		{/if}

		{#if overlay}
			{@render overlay()}
		{/if}
	</div>

	<div class="flex flex-col gap-0.5 pt-2">
		<div class="flex items-center gap-1">
			<a
				{href}
				class="min-w-0 flex-1 truncate text-sm font-medium group-hover:underline {nameClass}"
				title={name}
				onclick={(e) => {
					if (disabled) {
						e.preventDefault();
						onDisabledClick?.();
					}
				}}
				aria-disabled={disabled || undefined}
			>
				{name}
			</a>

			{#if menu}
				<DropdownMenu.Root>
					<DropdownMenu.Trigger
						class={cn(
							buttonVariants({ variant: "ghost", size: "icon-sm" }),
							"-mr-1 shrink-0 rounded-full text-muted-foreground",
						)}
						aria-label={`More options for ${name}`}
					>
						<EllipsisVertical size={14} />
					</DropdownMenu.Trigger>
					<DropdownMenu.Content align="center">
						{@render menu()}
					</DropdownMenu.Content>
				</DropdownMenu.Root>
			{/if}
		</div>

		{#if subtitle}
			{@render subtitle()}
		{/if}
	</div>
</div>
