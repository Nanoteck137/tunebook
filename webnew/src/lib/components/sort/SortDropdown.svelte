<script lang="ts">
	import { buttonVariants, DropdownMenu } from "$lib/components/ui";
	import { CheckIcon, ListSortAscendingIcon } from "@lucide/svelte";
	import type { SortType } from "./types";

	let {
		types,
		sort,
		onSortChange,
		align,
	}: {
		types: readonly { label: string; value: string }[];
		sort: SortType;
		onSortChange: (sort: SortType) => void;
		align?: "start" | "center" | "end";
	} = $props();
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
				{@const selected = sort === ty.value}
				<DropdownMenu.Item
					onSelect={() => onSortChange(ty.value as SortType)}
					class={selected ? "bg-accent text-foreground" : ""}
				>
					{#if selected}
						<CheckIcon />
					{/if}
					{ty.label}
				</DropdownMenu.Item>
			{/each}
		</DropdownMenu.Group>
	</DropdownMenu.Content>
</DropdownMenu.Root>
