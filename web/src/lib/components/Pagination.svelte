<script lang="ts">
  import { goto } from "$app/navigation";
  import { page } from "$app/stores";
  import type { Page } from "$lib/api/types";
  import { Pagination } from "@nanoteck137/nano-ui";
  import { ChevronLeft, ChevronRight } from "lucide-svelte";

  type Props = {
    page: Page;
  };

  const { page: pageInfo }: Props = $props();
</script>

<Pagination.Root
  page={pageInfo.page + 1}
  count={pageInfo.totalItems}
  perPage={pageInfo.perPage}
  siblingCount={0}
  onPageChange={(p) => {
    const query = $page.url.searchParams;
    query.set("page", (p - 1).toString());

    goto(`?${query.toString()}`, { invalidateAll: true, keepFocus: true });
  }}
>
  {#snippet children({ pages, currentPage })}
    <Pagination.Content class="overflow-x-auto">
      <Pagination.Item>
        <Pagination.PrevButton class="w-10 sm:w-20">
          {#snippet children()}
            <ChevronLeft size={16} />
            <span class="hidden sm:inline">Previous</span>
          {/snippet}
        </Pagination.PrevButton>
      </Pagination.Item>
      {#each pages as page (page.key)}
        {#if page.type === "ellipsis"}
          <Pagination.Item>
            <Pagination.Ellipsis />
          </Pagination.Item>
        {:else}
          <Pagination.Item>
            <Pagination.Link
              href="?page={page.value}"
              {page}
              isActive={currentPage === page.value}
            >
              {page.value}
            </Pagination.Link>
          </Pagination.Item>
        {/if}
      {/each}
      <Pagination.Item>
        <Pagination.NextButton class="w-10 sm:w-20">
          {#snippet children()}
            <span class="hidden sm:inline">Next</span>
            <ChevronRight size={16} />
          {/snippet}
        </Pagination.NextButton>
      </Pagination.Item>
    </Pagination.Content>
  {/snippet}
</Pagination.Root>
