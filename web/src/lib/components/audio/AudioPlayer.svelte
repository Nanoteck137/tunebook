<script lang="ts">
  import LargePlayer from "$lib/components/audio/LargePlayer.svelte";
  import { getMusicManager } from "$lib/music-manager.svelte";

  const musicManager = getMusicManager();

  $effect(() => {
    if (musicManager.showPlayer) {
      document.body.setAttribute("data-player", "true");
    } else {
      document.body.setAttribute("data-player", "false");
    }
  });
</script>

{#if musicManager.showPlayer}
  <LargePlayer />
{/if}
