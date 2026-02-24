<script lang="ts">
  import { Button } from "@nanoteck137/nano-ui";
  import { Plus } from "lucide-svelte";
  import NewTaglistModal from "./NewTaglistModal.svelte";

  const { data } = $props();

  let openNewTaglistModal = $state(false);
</script>

<Button
  variant="ghost"
  onclick={() => {
    openNewTaglistModal = true;
  }}
>
  <Plus />
  New Taglist
</Button>

<div class="flex flex-col gap-2">
  {#each data.taglists as taglist}
    <a href="/taglists/{taglist.id}">{taglist.name}</a>
  {/each}
</div>

<NewTaglistModal bind:open={openNewTaglistModal} />
