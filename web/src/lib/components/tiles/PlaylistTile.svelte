<script lang="ts">
  type Props = {
    id: string;
    cover: string;
    name: string;
    trackCount: number;
  };

  const { id, cover, name, trackCount }: Props = $props();
</script>

<div class="flex shrink-0 flex-col items-center">
  <a href="/playlists/{id}" class="group w-40 cursor-pointer">
    <img
      class="aspect-square w-40 rounded-lg object-cover"
      src={cover}
      alt=""
      title={name}
    />

    <div class="mt-2 w-40 truncate text-sm font-medium group-hover:underline">
      {name}
    </div>

    <div class="mt-1 text-xs text-muted-foreground">
      {trackCount} track(s)
    </div>
  </a>
</div>
