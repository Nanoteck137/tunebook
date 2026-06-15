<script lang="ts">
  import TrackList from "$lib/components/track-list/TrackList.svelte";

  const { data } = $props();
</script>

{#each data.history as history}
  <p>
    {history.id} - {history.track.id} - {history.status} - {history.percentPlayed}
  </p>
{/each}

<TrackList
  totalTracks={data.page.totalItems}
  tracks={data.history.map((e) => e.track)}
  onPlay={() => {}}
/>
