<script lang="ts">
  import { goto } from "$app/navigation";
  import { page } from "$app/state";
  import { Button, Input, Select, Separator } from "@nanoteck137/nano-ui";
  import { Play, Shuffle, Plus, X } from "lucide-svelte";
  import TrackList from "$lib/components/track-list/TrackList.svelte";
  import { getMusicManager } from "$lib/music-manager.svelte";
  import Spacer from "$lib/components/Spacer.svelte";
  import NewFilterModal from "./NewFilterModal.svelte";
  import FilterButton from "./FilterButton.svelte";
  import Pagination from "$lib/components/Pagination.svelte";
  import { sortTypes, defaultSort, type SortType } from "./types";

  let { data } = $props();
  const musicManager = getMusicManager();

  let openNewFilterModal = $state(false);
  let selectedSort = $state<SortType>(defaultSort);

  let tagInput = $state("");
  let tagMode = $state<"include" | "exclude">("include");
  let tags = $state<{ value: string; mode: "include" | "exclude" }[]>([]);

  function addTag() {
    const t = tagInput.trim();
    if (!t) return;

    if (!tags.some((x) => x.value === t && x.mode === tagMode)) {
      tags = [...tags, { value: t, mode: tagMode }];
    }

    tagInput = "";
  }

  function removeTag(value: string, mode: "include" | "exclude") {
    tags = tags.filter((t) => !(t.value === value && t.mode === mode));
  }

  let filterId = $derived(page.url.searchParams.get("filterId"));

  function clearFilter() {
    const query = page.url.searchParams;
    query.delete("filterId");
    goto("?" + query.toString(), {
      invalidateAll: true,
      replaceState: true,
    });
  }

  async function playTracks(options: { shuffle?: boolean } = {}) {
    if (filterId) {
      await musicManager.queueRequest(
        { type: "addFilter", filterId },
        options,
      );
    } else {
      await musicManager.addTracks({
        trackIds: data.tracks.map((t) => t.id),
        clear: true,
      });
    }
  }
</script>

<div class="flex flex-col gap-4">
  <div
    class="flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between"
  >
    <div class="flex items-baseline gap-2">
      <h1 class="text-xl font-bold">Tracks</h1>
      {#if data.page}
        <span class="text-sm text-muted-foreground"
          >{data.page.totalItems}</span
        >
      {/if}
    </div>

    <div class="flex items-center gap-2">
      <Button variant="outline" size="sm" onclick={() => playTracks({ shuffle: true })}>
        <Shuffle size={14} />
        Shuffle
      </Button>
      <Button size="sm" onclick={() => playTracks()}>
        <Play size={14} />
        Play All
      </Button>
    </div>
  </div>

  <div class="rounded-lg border bg-card p-3">
    <div
      class="flex flex-col gap-3 sm:flex-row sm:items-end sm:justify-between"
    >
      <div class="flex flex-1 flex-col gap-2 sm:flex-row sm:items-center">
        <Input class="h-9 sm:w-56" placeholder="Search tracks..." disabled />
        <Select.Root type="single" allowDeselect={false} onValueChange={(v) => selectedSort = v as SortType}>
          <Select.Trigger class="h-9 w-full sm:w-40">
            {sortTypes.find((i) => i.value === selectedSort)?.label ?? "Sort"}
          </Select.Trigger>
          <Select.Content>
            {#each sortTypes as ty (ty.value)}
              <Select.Item value={ty.value} label={ty.label} />
            {/each}
          </Select.Content>
        </Select.Root>
      </div>

    </div>

    <div class="mt-3 flex flex-wrap items-center gap-1.5">
      <span class="text-xs font-medium text-muted-foreground">Tags</span>

      <div class="flex items-center gap-1">
        <button
          class="rounded-l-md border px-1.5 py-1 text-xs font-medium transition-colors {tagMode ===
          'include'
            ? 'border-primary bg-primary text-primary-foreground'
            : 'bg-transparent text-muted-foreground hover:text-foreground'}"
          onclick={() => (tagMode = "include")}
        >
          + Inc
        </button>
        <button
          class="-ml-px rounded-r-md border px-1.5 py-1 text-xs font-medium transition-colors {tagMode ===
          'exclude'
            ? 'border-destructive bg-destructive text-destructive-foreground'
            : 'bg-transparent text-muted-foreground hover:text-foreground'}"
          onclick={() => (tagMode = "exclude")}
        >
          - Exc
        </button>
      </div>

      <Input
        class="h-7 w-28 text-xs"
        placeholder="Tag name..."
        bind:value={tagInput}
        onkeydown={(e) => {
          if (e.key === "Enter") {
            addTag();
          }
        }}
      />

      <button
        class="flex h-7 w-7 items-center justify-center rounded-md text-muted-foreground hover:text-foreground"
        onclick={addTag}
      >
        <Plus size={14} />
      </button>

      {#each tags as t (t.value + t.mode)}
        <span
          class="flex items-center gap-0.5 rounded-full px-2 py-0.5 text-xs {t.mode ===
          'include'
            ? 'bg-primary/10 text-primary'
            : 'bg-destructive/10 text-destructive'}"
        >
          {t.mode === "include" ? "+" : "-"}{t.value}
          <button
            class="hover:text-inherit/80"
            onclick={() => removeTag(t.value, t.mode)}
          >
            <X size={11} />
          </button>
        </span>
      {/each}

      {#if tags.length > 0}
        <button
          class="text-xs text-muted-foreground hover:text-foreground"
          onclick={() => (tags = [])}
        >
          Clear
        </button>
      {/if}
    </div>

    <div class="mt-3 flex flex-wrap items-center gap-2">
      {#if data.filters && data.filters.length > 0}
        <span class="text-xs font-medium text-muted-foreground"
          >Saved Filters</span
        >
        {#each data.filters as filter (filter.filterId)}
          <FilterButton {filter} />
        {/each}
      {/if}

      <Button
        variant="ghost"
        size="sm"
        onclick={() => (openNewFilterModal = true)}
      >
        <Plus size={14} />
        New Filter
      </Button>

      {#if page.url.searchParams.has("filterId")}
        <Button
          variant="ghost"
          size="sm"
          onclick={clearFilter}
        >
          <X size={14} />
          Clear
        </Button>
      {/if}
    </div>
  </div>
</div>

<Spacer size="md" />

<TrackList
  totalTracks={data.page.totalItems}
  tracks={data.tracks}
  onPlay={async (trackId) => {
    if (filterId) {
      await musicManager.queueRequest(
        { type: "addFilter", filterId },
        { queueIndexToTrackId: trackId },
      );
    } else {
      await musicManager.addTracks({
        trackIds: data.tracks.map((t) => t.id),
        trackId,
        clear: true,
      });
    }
  }}
/>

<Spacer size="lg" />
<Separator />
<Spacer size="lg" />

<Pagination page={data.page} />

<NewFilterModal bind:open={openNewFilterModal} />
