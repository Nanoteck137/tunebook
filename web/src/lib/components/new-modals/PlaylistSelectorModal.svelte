<script lang="ts">
  import { getManager } from "$lib/playlist-modal.svelte";
  import { Button, Dialog, ScrollArea } from "@nanoteck137/nano-ui";
  import { Check } from "lucide-svelte";

  const manager = getManager();

  $effect(() => {
    if (!manager.open && manager.resolve) {
      const resolve = manager.resolve;
      manager.resolve = null;
      resolve(null);
    }
  });
</script>

<Dialog.Root bind:open={manager.open}>
  <Dialog.Content>
    <Dialog.Header>
      <Dialog.Title>Select Playlist</Dialog.Title>
    </Dialog.Header>

    <ScrollArea class="max-h-[320px]">
      <div class="flex flex-col gap-1">
        {#each manager.playlists as playlist}
          <button
            class="flex w-full items-center gap-3 rounded-lg px-2 py-2 text-left hover:bg-accent disabled:opacity-50"
            disabled={manager.selectedId === playlist.id}
            onclick={() => manager.select(playlist.id)}
          >
            <img
              class="h-10 w-10 shrink-0 rounded object-cover"
              src={playlist.coverArt.small}
              alt=""
            />
            <div class="flex min-w-0 flex-1 flex-col">
              <span class="truncate text-sm font-medium">{playlist.name}</span>
              <span class="text-xs text-muted-foreground"
                >{playlist.trackCount} track(s)</span
              >
            </div>
            {#if manager.selectedId === playlist.id}
              <Check class="shrink-0 text-primary" size={18} />
            {/if}
          </button>
        {/each}
      </div>
    </ScrollArea>

    <Dialog.Footer>
      <Button
        variant="outline"
        onclick={() => {
          manager.open = false;
        }}
      >
        Close
      </Button>
    </Dialog.Footer>
  </Dialog.Content>
</Dialog.Root>
