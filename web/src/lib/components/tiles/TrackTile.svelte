<script lang="ts">
  import ArtistList, { type Artist } from "$lib/components/ArtistList.svelte";

  type Props = {
    id: string;
    cover: string;
    name: string;
    artists: Artist[];
  };

  const { cover, name, artists }: Props = $props();
</script>

<div class="flex shrink-0 flex-col items-center">
  <!-- svelte-ignore a11y_invalid_attribute -->
  <a href="#" class="group w-40 cursor-pointer">
    <img
      class="aspect-square w-40 rounded-lg object-cover"
      src={cover}
      alt=""
      title={name}
    />

    <div
      class="mt-2 w-40 truncate text-sm font-medium group-hover:underline"
      title={name}
    >
      {name}
    </div>
  </a>

  <ArtistList class="justify-start text-muted-foreground" {artists} />
</div>
