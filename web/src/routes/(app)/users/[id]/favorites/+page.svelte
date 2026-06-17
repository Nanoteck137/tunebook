<script lang="ts">
  import TrackList from "$lib/components/track-list/TrackList.svelte";
  import Spacer from "$lib/components/Spacer.svelte";
  import Pagination from "$lib/components/Pagination.svelte";
  import { Separator } from "@nanoteck137/nano-ui";

  let { data } = $props();
</script>

<TrackList
  totalTracks={data.page.totalItems}
  tracks={data.tracks}
  onPlay={() => {}}
/>

<Spacer size="lg" />

<Separator />

<Spacer size="lg" />

<Pagination page={data.page} />
