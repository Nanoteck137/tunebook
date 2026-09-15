<script lang="ts">
	import { page } from "$app/state";
	import { Calendar } from "@lucide/svelte";
	import { Breadcrumb, Button } from "$lib/components/ui";

	const { data, children } = $props();

	let createdString = $derived(
		new Date(data.userData.created).toLocaleDateString(undefined, {
			year: "numeric",
			month: "long",
			day: "numeric",
		}),
	);

	let roleLabel = $derived(
		(
			{
				super_user: "Super User",
				admin: "Admin",
				user: "User",
			} as Record<string, string>
		)[data.userData.role] ?? data.userData.role,
	);

	const tabs = $derived([
		{ label: "Overview", href: `/users/${data.userData.id}` },
		{ label: "Top", href: `/users/${data.userData.id}/top` },
		{ label: "Favorites", href: `/users/${data.userData.id}/favorites` },
		{ label: "History", href: `/users/${data.userData.id}/history` },
		{ label: "Review", href: `/users/${data.userData.id}/review` },
		{ label: "Settings", href: `/users/${data.userData.id}/settings` },
	]);
</script>

<div class="flex flex-col gap-6">
	<div class="py-2">
		<Breadcrumb.Root>
			<Breadcrumb.List>
				<Breadcrumb.Item>
					<Breadcrumb.Link href="/users">Users</Breadcrumb.Link>
				</Breadcrumb.Item>
				<Breadcrumb.Separator />
				<Breadcrumb.Item>
					<Breadcrumb.Page>{data.userData.displayName}</Breadcrumb.Page>
				</Breadcrumb.Item>
			</Breadcrumb.List>
		</Breadcrumb.Root>
	</div>

	<div
		class="flex flex-col gap-6 rounded-lg border bg-linear-to-b from-[oklch(0.93_0.045_75)] to-background p-4 shadow-sm sm:p-6 md:flex-row md:items-end md:gap-8 dark:from-[oklch(0.24_0.03_80)] dark:to-background"
	>
		<img
			class="h-40 w-40 self-center rounded-full shadow-2xl ring-1 ring-black/15 transition-transform duration-300 hover:scale-[1.02] md:h-52 md:w-52 dark:ring-white/10"
			src={data.userData.picture.large}
			alt=""
		/>

		<div class="flex min-w-0 flex-col gap-2">
			<p
				class="text-xs font-semibold uppercase tracking-wider text-muted-foreground"
			>
				Profile
			</p>

			<h1 class="line-clamp-2 text-2xl font-bold md:text-4xl">
				{data.userData.displayName}
			</h1>

			<p class="font-medium text-foreground">{roleLabel}</p>

			<p class="flex items-center gap-1.5 text-sm text-muted-foreground">
				<Calendar size={14} />
				<span>Member since {createdString}</span>
			</p>
		</div>
	</div>

	<nav class="flex flex-wrap gap-1">
		{#each tabs as tab (tab.href)}
			{#if tab.href === `/users/${data.userData.id}` || data.userData.id === data.user?.id}
				<Button
					class="transition-colors hover:bg-accent hover:text-accent-foreground dark:hover:bg-accent {page
						.url.pathname === tab.href
						? 'bg-accent text-accent-foreground'
						: 'text-muted-foreground'}"
					variant="ghost"
					href={tab.href}
				>
					{tab.label}
				</Button>
			{/if}
		{/each}
	</nav>

	<div class="min-w-0 flex-1">
		{@render children()}
	</div>
</div>