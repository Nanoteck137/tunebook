<script lang="ts">
	let {
		rank,
		name,
		playCount,
		href,
		subtitle,
		art,
	}: {
		rank: number;
		name: string;
		playCount: number;
		href?: string;
		subtitle?: string;
		art?: string;
	} = $props();
</script>

{#snippet row()}
	<span
		class="w-7 shrink-0 text-center text-sm font-bold text-muted-foreground"
	>
		{rank}
	</span>
	{#if art}
		<img src={art} alt="" class="size-10 shrink-0 rounded object-cover" />
	{/if}
	<span class="min-w-0 flex-1">
		<span class="block truncate text-sm font-medium">{name}</span>
		{#if subtitle}
			<span class="block truncate text-xs text-muted-foreground">
				{subtitle}
			</span>
		{/if}
	</span>
	<span class="shrink-0 text-xs text-muted-foreground">
		{playCount.toLocaleString()} plays
	</span>
{/snippet}

{#if href}
	<a
		{href}
		class="flex items-center gap-3 rounded-lg border bg-card p-2 transition-colors hover:bg-accent hover:text-accent-foreground"
	>
		{@render row()}
	</a>
{:else}
	<div class="flex items-center gap-3 rounded-lg border bg-card p-2">
		{@render row()}
	</div>
{/if}
