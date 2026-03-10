<script lang="ts">
  import { Button } from "@nanoteck137/nano-ui";
  import { Plus } from "lucide-svelte";
  import NewVirtualPlaylistModal from "./NewVirtualPlaylistModal.svelte";

  const { data } = $props();

  let openNewVirtualPlaylistModal = $state(false);
</script>

<Button
  variant="ghost"
  onclick={() => {
    openNewVirtualPlaylistModal = true;
  }}
>
  <Plus />
  New Taglist
</Button>

<div class="flex flex-col gap-2">
  {#each data.virtualPlaylists as virtualPlaylist}
    <a href="/virtual-playlists/{virtualPlaylist.id}">{virtualPlaylist.name}</a
    >
  {/each}
</div>

<NewVirtualPlaylistModal bind:open={openNewVirtualPlaylistModal} />
