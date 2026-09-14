<script lang="ts">
  import { page } from "$app/state";
  import { Heart, ListFilter, ListMusic } from "@lucide/svelte";
  import { Button } from "$lib/components/ui";
  import { cn } from "$lib/utils";

  let { children } = $props();

  const tabs = [
    {
      label: "Playlists",
      href: "/library/playlists",
      icon: ListMusic,
      match: "/library/playlists",
    },
    {
      label: "Favorites",
      href: "/library/favorites",
      icon: Heart,
      match: "/library/favorites",
    },
    {
      label: "Filters",
      href: "/library/filters/tracks",
      icon: ListFilter,
      match: "/library/filters",
    },
  ];

  let pathname = $derived(page.url.pathname);

  function isActive(tab: (typeof tabs)[number]) {
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