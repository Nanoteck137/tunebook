<script lang="ts">
	import { page } from "$app/state";
	import { LayoutDashboard, Library, ListChecks } from "@lucide/svelte";
	import { Button } from "$lib/components/ui";
	import { cn } from "$lib/utils";

	let { children } = $props();

	const tabs = [
		{
			label: "Main",
			href: "/server",
			icon: LayoutDashboard,
			match: "/server",
			exact: true,
		},
		{
			label: "Library",
			href: "/server/library",
			icon: Library,
			match: "/server/library",
		},
		{
			label: "Jobs",
			href: "/server/jobs",
			icon: ListChecks,
			match: "/server/jobs",
		},
	];

	let pathname = $derived(page.url.pathname);

	function isActive(tab: (typeof tabs)[number]) {
		if (tab.exact) {
			return pathname === tab.match;
		}

		return pathname.startsWith(tab.match);
	}
</script>

<div class="flex flex-col gap-4">
	<div class="flex items-center gap-2">
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

	{@render children()}
</div>
