<script lang="ts">
	import { getApiClient, handleApiError } from "$lib";
	import type { UserData } from "$lib/api/types";
	import InfiniteScroll from "$lib/components/InfiniteScroll.svelte";
	import { InfiniteScrollController } from "$lib/infinite-scroll.svelte";
	import Spacer from "$lib/components/Spacer.svelte";

	let { data } = $props();

	const apiClient = getApiClient();

	const scroll = new InfiniteScrollController<UserData>({
		initialLoad: () => ({
			items: data.users,
			hasMore: data.page.page + 1 < data.page.totalPages,
			page: data.page.page,
		}),
		load: async (page) => {
			const res = await apiClient.searchUsers({
				query: {
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

<div class="flex flex-col gap-4">
	<div class="flex items-center justify-between">
		<div class="flex items-center gap-2">
			<p class="text-bold text-xl">Users</p>
			{#if data.page}
				<span class="text-sm text-muted-foreground">
					({data.page.totalItems} users)
				</span>
			{/if}
		</div>
	</div>
</div>

<Spacer size="lg" />

<InfiniteScroll controller={scroll}>
	<div class="flex shrink flex-wrap justify-center gap-4">
		{#each scroll.items as user (user.id)}
			<div class="flex shrink-0 flex-col items-center">
				<a href="/users/{user.id}" class="group w-40 cursor-pointer">
					<img
						class="aspect-square w-40 rounded-full object-cover"
						src={user.picture.medium}
						alt=""
						title={user.displayName}
					/>

					<div
						class="mt-2 w-40 truncate text-center text-sm font-medium group-hover:underline"
						title={user.displayName}
					>
						{user.displayName}
					</div>

					<div class="mt-1 text-center text-xs text-muted-foreground">
						{user.role}
					</div>
				</a>
			</div>
		{/each}
	</div>
</InfiniteScroll>
