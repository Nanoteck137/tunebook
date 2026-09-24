<script lang="ts">
	import { page } from "$app/state";
	import {
		BarChart3,
		Calendar,
		CalendarRange,
		Heart,
		History,
		LayoutDashboard,
		ListMusic,
		Settings,
	} from "@lucide/svelte";
	import { Breadcrumb, Button } from "$lib/components/ui";
	import { cn } from "$lib/utils";

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
		{
			label: "Overview",
			href: `/users/${data.userData.id}`,
			public: true,
			icon: LayoutDashboard,
			match: `/users/${data.userData.id}`,
			exact: true,
		},
		{
			label: "Top",
			href: `/users/${data.userData.id}/top`,
			public: true,
			icon: BarChart3,
			match: `/users/${data.userData.id}/top`,
		},
		{
			label: "Playlists",
			href: `/users/${data.userData.id}/playlists`,
			public: true,
			icon: ListMusic,
			match: `/users/${data.userData.id}/playlists`,
		},
		{
			label: "Favorites",
			href: `/users/${data.userData.id}/favorites`,
			icon: Heart,
			match: `/users/${data.userData.id}/favorites`,
		},
		{
			label: "History",
			href: `/users/${data.userData.id}/history`,
			icon: History,
			match: `/users/${data.userData.id}/history`,
		},
		{
			label: "Review",
			href: `/users/${data.userData.id}/review`,
			icon: CalendarRange,
			match: `/users/${data.userData.id}/review`,
		},
		{
			label: "Settings",
			href: `/users/${data.userData.id}/settings`,
			icon: Settings,
			match: `/users/${data.userData.id}/settings`,
		},
	]);

	let pathname = $derived(page.url.pathname);

	function isActive(tab: (typeof tabs)[number]) {
		if (tab.exact) {
			return pathname === tab.match;
		}

		return pathname.startsWith(tab.match);
	}
</script>

<div class="section-users flex flex-col gap-6">
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
		class="flex min-w-0 flex-col gap-4 rounded-lg border bg-linear-to-b from-section-hero-from to-section-hero-to p-4 shadow-sm sm:p-6 md:gap-8"
	>
		<div class="flex min-w-0 flex-col gap-4 md:hidden">
			<div class="flex items-center gap-4">
				<img
					class="h-20 min-h-20 w-20 min-w-20 shrink-0 rounded-full shadow-2xl ring-1 ring-black/15 transition-transform duration-300 hover:scale-[1.02] dark:ring-white/10"
					src={data.userData.picture.large}
					alt=""
				/>

				<div class="flex min-w-0 flex-col gap-0.5">
					<p
						class="text-xs font-semibold tracking-wider text-muted-foreground uppercase"
					>
						Profile
					</p>

					<h1 class="line-clamp-2 text-2xl font-bold">
						{data.userData.displayName}
					</h1>
				</div>
			</div>

			<div class="flex flex-col items-start gap-1 text-sm">
				<p class="font-medium text-foreground">{roleLabel}</p>

				<p class="flex items-center gap-1.5 text-muted-foreground">
					<Calendar size={14} />
					<span>Member since {createdString}</span>
				</p>
			</div>
		</div>

		<div class="hidden gap-6 md:flex md:items-end md:gap-8">
			<img
				class="h-52 min-h-52 w-52 min-w-52 shrink-0 rounded-full shadow-2xl ring-1 ring-black/15 transition-transform duration-300 hover:scale-[1.02] dark:ring-white/10"
				src={data.userData.picture.large}
				alt=""
			/>

			<div class="flex min-w-0 flex-col gap-2">
				<p
					class="text-xs font-semibold tracking-wider text-muted-foreground uppercase"
				>
					Profile
				</p>

				<h1 class="line-clamp-2 text-4xl font-bold">
					{data.userData.displayName}
				</h1>

				<p class="font-medium text-foreground">{roleLabel}</p>

				<p class="flex items-center gap-1.5 text-sm text-muted-foreground">
					<Calendar size={14} />
					<span>Member since {createdString}</span>
				</p>
			</div>
		</div>

		<div
			class="flex flex-wrap items-center gap-1 border-t border-border/40 pt-3"
		>
			{#each tabs as tab (tab.href)}
				{#if tab.public || data.userData.id === data.user?.id}
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
				{/if}
			{/each}
		</div>
	</div>

	<div class="min-w-0 flex-1">
		{@render children()}
	</div>
</div>
