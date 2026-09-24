<script lang="ts">
	import { ArrowDown, ArrowUp, ArrowUpDown } from "@lucide/svelte";
	import type { Snippet } from "svelte";
	import type { SortType } from "./types";
	import { cn } from "$lib/utils";

	export type Column = {
		label: string;
		asc: SortType;
		desc: SortType;
		className: string;
		title?: string;
	};

	let {
		class: className = "",

		sort,
		onSortChange,
		columns,
		children,
	}: {
		class?: string;

		sort: SortType;
		onSortChange: (sort: SortType) => void;
		columns: Column[];
		children?: Snippet;
	} = $props();

	function cycle(asc: SortType, desc: SortType) {
		if (sort === asc) {
			onSortChange(desc);
		} else if (sort === desc) {
			onSortChange(asc);
		} else {
			onSortChange(asc);
		}
	}

	function active(col: Column) {
		return sort === col.asc || sort === col.desc;
	}
</script>

<div
	class={cn("group/header hidden items-center gap-3 pl-2 pr-4 pt-2 pb-1.5 text-xs font-medium text-muted-foreground select-none sm:flex sm:gap-4", className)}
	role="row"
>
	{#each columns as col (col.label)}
		{@const isAsc = sort === col.asc}
		{@const isDesc = sort === col.desc}
		<button
			class="{col.className} {active(col) ? 'text-foreground' : ''}"
			title={col.title ?? `Sort by ${col.label}`}
			aria-label={col.title ?? `Sort by ${col.label}`}
			onclick={() => cycle(col.asc, col.desc)}
		>
			{#if isDesc}
				<ArrowDown size={12} class="shrink-0 text-primary" />
			{:else if isAsc}
				<ArrowUp size={12} class="shrink-0 text-primary" />
			{:else}
				<ArrowUpDown
					size={12}
					class="shrink-0 opacity-0 transition-opacity group-hover/header:opacity-50"
				/>
			{/if}
			{col.label}
		</button>
	{/each}

	{#if children}
		{@render children()}
	{/if}
</div>
