<script lang="ts">
  import { goto } from "$app/navigation";
  import { Button, Input } from "@nanoteck137/nano-ui";
  import { Search } from "lucide-svelte";
  import Pagination from "$lib/components/Pagination.svelte";
  import Image from "$lib/components/Image.svelte";
  import { onMount } from "svelte";

  let { data } = $props();

  async function doSearch(query: string) {
    await goto(`/search/users?query=${query}`, {
      invalidateAll: true,
      keepFocus: true,
      replaceState: true,
    });
  }

  let initialValue = $state("");
  let value = "";

  onMount(() => {
    initialValue = data.query;
    value = data.query;
  });

  let timer: ReturnType<typeof setTimeout>;
  function onInput(e: Event) {
    const target = e.target as HTMLInputElement;
    const current = target.value;
    value = current;

    clearTimeout(timer);
    timer = setTimeout(() => {
      doSearch(current);
    }, 500);
  }
</script>

<svelte:head>
  <title>Search Users - Tunebook</title>
</svelte:head>

<div class="flex flex-col gap-6">
  <div class="flex flex-col gap-2">
    <h1 class="text-xl font-bold">Search Users</h1>

    <form
      action=""
      method="get"
      onsubmit={(e) => {
        e.preventDefault();
        clearTimeout(timer);
        doSearch(value);
      }}
    >
      <div class="flex items-center gap-2">
        <div class="relative flex-1">
          <Search
            size={16}
            class="pointer-events-none absolute left-3 top-1/2 -translate-y-1/2 text-muted-foreground"
          />
          <Input
            id="query"
            name="query"
            placeholder="Search users..."
            autocomplete="off"
            value={initialValue}
            oninput={onInput}
            class="pl-9"
          />
        </div>
        <Button type="submit">Search</Button>
      </div>
    </form>
  </div>

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
      {#each data.users as user}
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
      <Pagination page={data.page} />
    {/if}
  {/if}
</div>
