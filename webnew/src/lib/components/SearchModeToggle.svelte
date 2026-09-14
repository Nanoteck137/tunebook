<script lang="ts">
	import { Database, Sparkles } from "@lucide/svelte";
	import { ToggleGroup } from "$lib/components/ui";

	let { value = $bindable<string>("query") }: { value?: string } = $props();
</script>

<ToggleGroup.Root
	type="single"
	size="sm"
	variant="outline"
	bind:value
>
	<ToggleGroup.Item value="query" title="Search the database directly (Query API)">
		<Database size={14} />
		Query
	</ToggleGroup.Item>
	<ToggleGroup.Item value="search" title="Search the Meilisearch index">
		<Sparkles size={14} />
		Search
	</ToggleGroup.Item>
</ToggleGroup.Root>