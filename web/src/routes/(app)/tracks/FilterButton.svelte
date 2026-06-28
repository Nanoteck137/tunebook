<script lang="ts">
  import { goto, invalidateAll } from "$app/navigation";
  import { page } from "$app/state";
  import { getApiClient, handleApiError } from "$lib";
  import type { TrackFilter } from "$lib/api/types";
  import { Edit, Trash } from "lucide-svelte";
  import EditFilterModal from "./EditFilterModal.svelte";
  import toast from "svelte-5-french-toast";

  type Props = {
    filter: TrackFilter;
  };

  const { filter }: Props = $props();
  const apiClient = getApiClient();

  let openFilterEditModal = $state(false);
  let openDropdown = $state(false);

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

  async function deleteFilter() {
    const res = await apiClient.deleteTrackFilter(filter.filterId);
    if (!res.success) {
      return handleApiError(res.error);
    }

    const query = page.url.searchParams;
    if (query.get("filterId") === filter.filterId) {
      query.delete("filterId");
    }

    toast.success("Filter deleted");
    invalidateAll();
  }

  function handleClickOutside(e: MouseEvent) {
    const target = e.target as HTMLElement;
    if (!target.closest(".filter-dropdown")) {
      openDropdown = false;
    }
  }

  $effect(() => {
    if (openDropdown) {
      document.addEventListener("click", handleClickOutside);
      return () => document.removeEventListener("click", handleClickOutside);
    }
  });
</script>

<div class="filter-dropdown flex items-center gap-1">
  <button
    class="rounded-full border px-3 py-1 text-sm font-medium transition-colors {active
      ? 'border-primary bg-primary text-primary-foreground'
      : 'border-border bg-card text-foreground hover:bg-accent'}"
    onclick={selectFilter}
  >
    {filter.name}
  </button>

  <div class="relative">
    <button
      class="flex h-6 w-6 items-center justify-center rounded-full text-muted-foreground hover:bg-accent hover:text-foreground"
      aria-label="Filter options"
      onclick={() => (openDropdown = !openDropdown)}
    >
      <svg
        width="12"
        height="12"
        viewBox="0 0 24 24"
        fill="none"
        stroke="currentColor"
        stroke-width="2"
        stroke-linecap="round"
        stroke-linejoin="round"
      >
        <polyline points="6 9 12 15 18 9" />
      </svg>
    </button>

    {#if openDropdown}
      <div
        class="absolute right-0 top-full z-50 mt-1 w-36 overflow-hidden rounded-md border bg-popover shadow-md"
      >
        <button
          class="flex w-full items-center gap-2 px-3 py-2 text-left text-sm hover:bg-accent"
          onclick={() => {
            openDropdown = false;
            openFilterEditModal = true;
          }}
        >
          <Edit size={14} />
          Edit
        </button>
        <button
          class="flex w-full items-center gap-2 px-3 py-2 text-left text-sm text-destructive hover:bg-accent"
          onclick={async () => {
            openDropdown = false;
            await deleteFilter();
          }}
        >
          <Trash size={14} />
          Delete
        </button>
      </div>
    {/if}
  </div>
</div>

<EditFilterModal bind:open={openFilterEditModal} {filter} />
