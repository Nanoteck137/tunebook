<script lang="ts">
  import { goto, invalidateAll } from "$app/navigation";
  import { page } from "$app/state";
  import type { TrackFilter } from "$lib/api/types";
  import { ListFilter } from "@lucide/svelte";
  import { cn } from "$lib/utils";

  type Props = {
    filter: TrackFilter;
  };

  const { filter }: Props = $props();

  let active = $derived(
    page.url.searchParams.get("filterId") === filter.filterId,
  );

  function selectFilter() {
    const query = page.url.searchParams;
    query.set("filterId", filter.filterId);
    goto("?" + query.toString(), {
      invalidateAll: true,
      replaceState: true,
    });
  }
</script>

<button
  class={cn(
    "flex items-center gap-1.5 rounded-full border py-1 pr-2.5 pl-2.5 text-sm font-medium transition-colors",
    active
      ? "border-primary bg-primary text-primary-foreground hover:bg-primary/90"
      : "border-border bg-card text-foreground hover:bg-accent",
  )}
  onclick={selectFilter}
  title={active ? "Active filter" : "Use this filter"}
>
  <ListFilter size={12} class={active ? "shrink-0" : "shrink-0 text-muted-foreground"} />
  <span class="max-w-32 truncate">{filter.name}</span>
</button>