<script lang="ts">
  import { goto } from "$app/navigation";
  import { page } from "$app/stores";
  import type { TrackFilter } from "$lib/api/types";
  import { Button } from "@nanoteck137/nano-ui";
  import { Edit } from "lucide-svelte";
  import EditPlaylistFilterModal from "./EditPlaylistFilterModal.svelte";

  type Props = {
    filter: TrackFilter;
  };

  const { filter }: Props = $props();

  let openFilterEditModal = $state(false);
</script>

<div class="flex gap-2">
  <Button
    onclick={() => {
      const query = $page.url.searchParams;
      query.set("filterId", filter.filterId);
      goto("?" + query.toString(), {
        invalidateAll: true,
        replaceState: true,
      });
    }}
  >
    Filter: {filter.filterId} - {filter.name}
  </Button>

  <Button
    variant="outline"
    size="icon"
    onclick={() => {
      openFilterEditModal = true;
    }}
  >
    <Edit />
  </Button>
</div>

<EditPlaylistFilterModal bind:open={openFilterEditModal} {filter} />
