<script lang="ts">
  import { getManager } from "$lib/playlist-modal.svelte";
  import { Button, Dialog, Input, ScrollArea } from "@nanoteck137/nano-ui";
  import { Check, ListMusic, Search } from "lucide-svelte";

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
  <Dialog.Content class="overflow-hidden sm:max-w-md">
    <div class="relative">
      <div
        class="absolute -right-16 -top-16 h-40 w-40 rounded-full bg-gradient-to-tr from-logo-1/10 via-logo-2/10 to-logo-3/10 blur-xl"
      ></div>

      <Dialog.Header class="relative text-left">
        <div class="flex items-center gap-3">
          <div
            class="flex h-10 w-10 shrink-0 items-center justify-center rounded-full bg-gradient-to-tr from-logo-1 via-logo-2 to-logo-3"
          >
            <ListMusic size={18} class="text-white" />
          </div>
          <div>
            <Dialog.Title class="text-xl sm:text-2xl">
              <span
                class="bg-gradient-to-tr from-logo-1 via-logo-2 to-logo-3 bg-clip-text text-transparent"
              >
                {manager.title}
              </span>
            </Dialog.Title>
            <Dialog.Description>
              {manager.description}
            </Dialog.Description>
          </div>
        </div>
      </Dialog.Header>
    </div>

    <div class="relative mb-3">
      <Search
        size={16}
        class="pointer-events-none absolute left-3 top-1/2 -translate-y-1/2 text-muted-foreground"
      />
      <Input
        class="h-9 w-full pl-9"
        placeholder="Search playlists..."
        bind:value={manager.searchQuery}
        oninput={() => manager.search(manager.searchQuery)}
        tabindex={-1}
      />
    </div>

    <ScrollArea class="max-h-[320px]">
      <div class="flex flex-col gap-1">
        {#each manager.playlists as playlist (playlist.id)}
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
        {:else}
          <p class="py-8 text-center text-sm text-muted-foreground">
            No playlists found
          </p>
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
