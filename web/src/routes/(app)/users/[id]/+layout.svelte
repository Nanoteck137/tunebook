<script lang="ts">
  import { Calendar } from "lucide-svelte";

  const { data, children } = $props();

  let createdString = $derived(
    new Date(data.userData.created).toLocaleDateString(undefined, {
      year: "numeric",
      month: "long",
      day: "numeric",
    }),
  );

  let roleLabel = $derived(
    ({
      super_user: "Super User",
      admin: "Admin",
      user: "User",
    } as Record<string, string>)[data.userData.role] ?? data.userData.role,
  );
</script>

<div class="flex flex-col gap-8 md:flex-row">
  <div
    class="flex shrink-0 flex-col items-center gap-4 md:sticky md:top-24 md:w-48 md:items-start md:self-start"
  >
    <img
      class="h-24 w-24 rounded-full object-cover"
      src={data.userData.picture.large}
      alt=""
    />

    <div class="text-center md:text-left">
      <h1 class="text-2xl font-bold">{data.userData.displayName}</h1>
      <p class="text-sm text-muted-foreground">{roleLabel}</p>
      <div class="flex items-center gap-1.5 text-xs text-muted-foreground">
        <Calendar size={12} />
        <span>Member since {createdString}</span>
      </div>
    </div>

    <div class="flex gap-4 md:flex-col">
      <a
        class="text-sm font-medium hover:underline"
        href="/users/{data.userData.id}"
      >
        Overview
      </a>

      {#if data.userData.id === data.user?.id}
        <a
          class="text-sm font-medium hover:underline"
          href="/users/{data.userData.id}/favorites"
        >
          Favorites
        </a>
      {/if}

      {#if data.userData.id === data.user?.id}
        <a
          class="text-sm font-medium hover:underline"
          href="/users/{data.userData.id}/history"
        >
          History
        </a>
      {/if}

      {#if data.userData.id === data.user?.id}
        <a
          class="text-sm font-medium hover:underline"
          href="/users/{data.userData.id}/settings"
        >
          Settings
        </a>
      {/if}
    </div>
  </div>

  <div class="min-w-0 flex-1">
    {@render children()}
  </div>
</div>
