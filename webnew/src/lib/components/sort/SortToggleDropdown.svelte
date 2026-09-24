<script lang="ts">
	import { buttonVariants, DropdownMenu } from "$lib/components/ui";
	import { ArrowDown, ArrowUp, ListSortAscendingIcon } from "@lucide/svelte";
	import type { SortToggleType } from "./types";

	let {
		types,
		sort,
		onSortChange,
		align,
	}: {
		types: readonly SortToggleType[];
		sort: string;
		onSortChange: (sort: string) => void;
		align?: "start" | "center" | "end";
	} = $props();

	function select(ty: SortToggleType) {
		if (sort === ty.value) {
			onSortChange(ty.reverse);
		} else {
			onSortChange(ty.value);
		}
	}
</script>

<DropdownMenu.Root>
	<DropdownMenu.Trigger
		class={buttonVariants({ variant: "ghost", size: "icon" })}
		title="Sort"
		aria-label="Sort"
	>
		<ListSortAscendingIcon />
	</DropdownMenu.Trigger>

	<DropdownMenu.Content align={align ?? "end"}>
		<DropdownMenu.Group>
			{#each types as ty (ty.value)}
				{@const reversed = sort === ty.reverse}
				{@const active = sort === ty.value || reversed}
				{@const dir = active
					? reversed
						? ty.direction === "asc"
							? "desc"
							: "asc"
						: ty.direction
					: null}
				<DropdownMenu.Item
					onSelect={() => select(ty)}
					class={active ? "bg-accent text-foreground" : ""}
				>
					{#if dir === "asc"}
						<ArrowUp class="h-4 w-4 text-primary" />
					{:else if dir === "desc"}
						<ArrowDown class="h-4 w-4 text-primary" />
					{:else}
						<ArrowUp class="h-4 w-4 opacity-0" aria-hidden="true" />
					{/if}
					{ty.label}
				</DropdownMenu.Item>
			{/each}
		</DropdownMenu.Group>
	</DropdownMenu.Content>
</DropdownMenu.Root>
