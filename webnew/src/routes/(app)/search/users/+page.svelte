<script lang="ts">
  import { goto } from "$app/navigation";
  import { Separator } from "$lib/components/ui";
  import Pagination from "$lib/components/Pagination.svelte";
  import Spacer from "$lib/components/Spacer.svelte";
  import Image from "$lib/components/Image.svelte";
  import { onMount } from "svelte";
  import SearchBarHeader from "../SearchBarHeader.svelte";

  let { data } = $props();

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

  {#if data.query && data.users.length === 0}
    <p class="py-12 text-center text-sm text-muted-foreground">
      No users found for "{data.query}".
    </p>
  {/if}

  {#if data.users.length > 0}
    {#if data.page}
      <div class="flex items-baseline gap-2">
        <span class="text-sm text-muted-foreground">
          {data.page.totalItems} user(s)
        </span>
      </div>
    {/if}

    <div
      class="grid grid-cols-2 gap-3 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-5 xl:grid-cols-6 2xl:grid-cols-7"
    >
      {#each data.users as user (user.id)}
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

    {#if data.page}
      <Spacer size="lg" />
      <Separator />
      <Spacer size="lg" />

      <Pagination page={data.page} />
    {/if}
  {/if}
</div>
