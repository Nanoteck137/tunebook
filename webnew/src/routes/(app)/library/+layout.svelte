<script lang="ts">
	import { page } from "$app/state";
	import { Heart, Library, ListFilter, ListMusic } from "@lucide/svelte";
	import { Button } from "$lib/components/ui";
	import { cn } from "$lib/utils";

	let { children } = $props();

	const tabs = [
		{
			label: "Playlists",
			href: "/library/playlists",
			icon: ListMusic,
			match: "/library/playlists",
		},
		{
			label: "Favorites",
			href: "/library/favorites",
			icon: Heart,
			match: "/library/favorites",
		},
		{
			label: "Filters",
			href: "/library/filters/tracks",
			icon: ListFilter,
			match: "/library/filters",
		},
	];

	let pathname = $derived(page.url.pathname);

	function isActive(tab: (typeof tabs)[number]) {
		return pathname.startsWith(tab.match);
	}
</script>

<div class="section-library flex flex-col gap-6">
	<div
		class="flex flex-col gap-4 rounded-lg border bg-linear-to-b from-section-hero-from to-section-hero-to p-4 shadow-sm sm:p-6"
	>
		<div class="flex items-center gap-4">
			<div
				class="flex h-12 w-12 shrink-0 items-center justify-center rounded-xl bg-linear-to-tr from-section-gradiant-1 via-section-gradiant-2 to-section-gradiant-3 text-white"
			>
				<Library size={24} />
			</div>
			<div class="flex min-w-0 flex-col">
				<h1 class="text-2xl font-bold">Library</h1>
				<p class="text-sm text-muted-foreground">
					Your playlists, favorites and saved filters
				</p>
			</div>
		</div>

		<div
			class="flex flex-wrap items-center gap-1 border-t border-border/40 pt-3"
		>
			{#each tabs as tab (tab.href)}
				<Button
					variant="ghost"
					size="sm"
					href={tab.href}
					class={cn(
						"text-foreground hover:bg-accent hover:text-accent-foreground",
						isActive(tab)
							? "bg-accent text-accent-foreground"
							: "text-muted-foreground",
					)}
				>
					<tab.icon />
					{tab.label}
				</Button>
			{/each}
		</div>
	</div>

	{@render children()}
</div>
