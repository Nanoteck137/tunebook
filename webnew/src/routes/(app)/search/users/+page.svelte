<script lang="ts">
	import { goto } from "$app/navigation";
	import { getApiClient, handleApiError } from "$lib";
	import type { UserData } from "$lib/api/types";
	import Image from "$lib/components/Image.svelte";
	import InfiniteScroll from "$lib/components/InfiniteScroll.svelte";
	import { InfiniteScrollController } from "$lib/infinite-scroll.svelte";
	import { onMount } from "svelte";
	import SearchBarHeader from "../SearchBarHeader.svelte";

	let { data } = $props();

	const apiClient = getApiClient();

	async function doSearch(query: string) {
		await goto(`/search/users?query=${query}`, {
			invalidateAll: true,
			keepFocus: true,
			replaceState: true,
		});
	}

	function clearSearch() {
		value = "";
		doSearch("");
	}

	let value = $state("");

	onMount(() => {
		value = data.query;
	});

	const scroll = new InfiniteScrollController<UserData>({
		initialLoad: () => {
			if (!data.page) return { items: [], hasMore: false, page: 0 };

			return {
				items: data.users,
				hasMore: data.page.page + 1 < data.page.totalPages,
				page: data.page.page,
			};
		},
		load: async (page) => {
			const res = await apiClient.searchUsers({
				query: {
					query: data.query,
					page: String(page),
					perPage: String(data.page?.perPage ?? 30),
				},
			});

			if (!res.success) {
				handleApiError(res.error);
				return null;
			}

			return {
				items: res.data.users,
				hasMore: res.data.page.page + 1 < res.data.page.totalPages,
			};
		},
		itemKey: (user) => user.id,
	});
</script>

<svelte:head>
	<title>Search Users - Tunebook</title>
</svelte:head>

<div class="flex flex-col gap-6">
	<SearchBarHeader
		searchBarPlaceholder="Search users..."
		{value}
		setValue={(v) => {
			value = v;
		}}
		search={doSearch}
		searchWithValue={() => {
			doSearch(value);
		}}
		{clearSearch}
	/>

	{#if data.query && scroll.items.length === 0}
		<p class="py-12 text-center text-sm text-muted-foreground">
			No users found for "{data.query}".
		</p>
	{/if}

	{#if scroll.items.length > 0}
		{#if data.page}
			<div class="flex items-baseline gap-2">
				<span class="text-sm text-muted-foreground">
					{data.page.totalItems} user(s)
				</span>
			</div>
		{/if}

		<InfiniteScroll controller={scroll}>
			<div
				class="grid grid-cols-2 gap-3 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-5 xl:grid-cols-6 2xl:grid-cols-7"
			>
				{#each scroll.items as user (user.id)}
					<div
						class="group relative flex flex-col overflow-hidden rounded-lg border bg-card transition-shadow hover:shadow-md"
					>
						<a href="/users/{user.id}">
							<Image
								class="aspect-square w-full rounded-none border-0"
								src={user.picture.medium}
								alt={user.displayName}
							/>
						</a>

						<div class="flex flex-col gap-0.5 p-2">
							<a
								href="/users/{user.id}"
								class="truncate text-sm font-medium hover:underline"
								title={user.displayName}
							>
								{user.displayName}
							</a>
							<p class="truncate text-xs text-muted-foreground">
								{user.role}
							</p>
						</div>
					</div>
				{/each}
			</div>
		</InfiniteScroll>
	{/if}
</div>
