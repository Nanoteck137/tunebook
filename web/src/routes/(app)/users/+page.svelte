<script lang="ts">
  import Spacer from "$lib/components/Spacer.svelte";
  import Pagination from "$lib/components/Pagination.svelte";
  import { Separator } from "@nanoteck137/nano-ui";

  let { data } = $props();
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

<div class="flex flex-shrink flex-wrap justify-center gap-4">
  {#each data.users as user}
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

<Spacer size="lg" />

<Separator />

<Spacer size="lg" />

<Pagination page={data.page} />
