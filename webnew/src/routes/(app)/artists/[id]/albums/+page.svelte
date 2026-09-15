<script lang="ts">
	import { goto } from "$app/navigation";
	import { page } from "$app/state";
	import { Breadcrumb, Select } from "$lib/components/ui";
	import Image from "$lib/components/Image.svelte";
	import { getApiClient, handleApiError } from "$lib";
	import type { Album } from "$lib/api/types";
	import InfiniteScroll from "$lib/components/InfiniteScroll.svelte";
	import { InfiniteScrollController } from "$lib/infinite-scroll.svelte";
	import { sortTypes, defaultSort, applySort, type SortType } from "./types";

	let { data } = $props();
	const apiClient = getApiClient();

	let sort = $state(
		(page.url.searchParams.get("sort") as SortType) ?? defaultSort,
	);

	function updateSort(value: string) {
		sort = value as SortType;

		const query = page.url.searchParams;
		query.delete("sort");

		if (sort !== defaultSort) {
			query.set("sort", sort);
		}

		goto("?" + query.toString(), { invalidateAll: true });
	}

	const scroll = new InfiniteScrollController<Album>({
		initialLoad: () => ({
			items: data.albums,
			hasMore: data.page.page + 1 < data.page.totalPages,
			page: data.page.page,
		}),
		load: async (nextPage) => {
			const query: Record<string, string> = {
				page: String(nextPage),
				perPage: String(data.page.perPage),
				filter: `artistId = "${data.artist.id}" or featuringArtists has "${data.artist.id}"`,
			};

			applySort(data.filter, query);

			const res = await apiClient.getAlbums({ query });
			if (!res.success) {
				handleApiError(res.error);
				return null;
			}

			return {
				items: res.data.albums,
				hasMore: res.data.page.page + 1 < res.data.page.totalPages,
			};
		},
		itemKey: (album) => album.id,
	});
</script>

<div class="flex flex-col gap-4">
	<Breadcrumb.Root>
		<Breadcrumb.List>
			<Breadcrumb.Item>
				<Breadcrumb.Link href="/artists">Artists</Breadcrumb.Link>
			</Breadcrumb.Item>
			<Breadcrumb.Separator />
			<Breadcrumb.Item>
				<Breadcrumb.Link href="/artists/{data.artist.id}">
					{data.artist.name}
				</Breadcrumb.Link>
			</Breadcrumb.Item>
			<Breadcrumb.Separator />
			<Breadcrumb.Item>
				<Breadcrumb.Page>Albums</Breadcrumb.Page>
			</Breadcrumb.Item>
		</Breadcrumb.List>
	</Breadcrumb.Root>

	<div
		class="flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between"
	>
		<div class="flex items-baseline gap-2">
			<h1 class="text-xl font-bold">Albums</h1>
			{#if data.page}
				<span class="text-sm text-muted-foreground"
					>{data.page.totalItems}</span
				>
			{/if}
		</div>

		<Select.Root
			type="single"
			allowDeselect={false}
			value={sort}
			onValueChange={updateSort}
		>
			<Select.Trigger class="h-9 w-full sm:w-40">
				{sortTypes.find((i) => i.value === sort)?.label ?? "Sort"}
			</Select.Trigger>
			<Select.Content>
				{#each sortTypes as ty (ty.value)}
					<Select.Item value={ty.value} label={ty.label} />
				{/each}
			</Select.Content>
		</Select.Root>
	</div>

	<InfiniteScroll controller={scroll}>
		<div
			class="grid grid-cols-2 gap-3 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-5 xl:grid-cols-6 2xl:grid-cols-7"
		>
			{#each scroll.items as album (album.id)}
				<a
					href="/albums/{album.id}"
					class="group flex flex-col overflow-hidden rounded-lg border bg-card transition-shadow hover:shadow-md"
				>
					<Image
						class="aspect-square w-full rounded-none border-0"
						src={album.coverArt.medium}
						alt={album.name}
					/>
					<div class="flex flex-col gap-0.5 p-2">
						<p
							class="truncate text-sm font-medium group-hover:underline"
							title={album.name}
						>
							{album.name}
						</p>
						<p
							class="truncate text-xs text-muted-foreground"
							title={album.artists.map((a) => a.name).join(", ")}
						>
							{album.artists.map((a) => a.name).join(", ")}
						</p>
					</div>
				</a>
			{/each}
		</div>
	</InfiniteScroll>
</div>
