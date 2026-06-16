<script lang="ts">
  import { getQuickPlaylist } from "$lib/quick-playlist.svelte";
  import { Button } from "@nanoteck137/nano-ui";
  import { Star } from "lucide-svelte";

  type Props = {
    trackId: string;
  };

  const { trackId }: Props = $props();
  const quickPlaylist = getQuickPlaylist();
</script>

{#if quickPlaylist.playlist !== null}
  <Button
    type="submit"
    class="rounded-full"
    variant="ghost"
    size="icon-lg"
    onclick={() => {
      quickPlaylist.toggleTrack(trackId);
    }}
    title={"Quick Add: " + quickPlaylist.playlist.name}
  >
    {#if quickPlaylist.hasTrack(trackId)}
      <Star class="fill-primary" />
    {:else}
      <Star />
    {/if}
  </Button>
{/if}
