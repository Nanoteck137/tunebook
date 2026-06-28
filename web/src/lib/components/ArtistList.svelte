<script lang="ts">
  import { cn } from "$lib/utils";

  export type Artist = {
    id: string;
    name: string;
  };

  type Props = {
    class?: string;
    artists: Artist[];
  };

  const { class: className, artists }: Props = $props();
</script>

<p
  class={cn("line-clamp-1 w-full text-xs", className)}
  title={artists.map((a) => a.name).join(", ")}
>
  {#each artists as artist, i (artist.id)}
    {#if i > 0}
      {", "}
    {/if}
    <a class="hover:underline" href="/artists/{artist.id}" title={artist.name}>
      {artist.name}
    </a>
  {/each}
</p>
