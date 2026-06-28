<script lang="ts">
  import type { Modal } from "$lib/components/new-modals";
  import { Button, Dialog } from "@nanoteck137/nano-ui";
  import { TriangleAlert } from "lucide-svelte";

  export type Props = {
    open?: boolean;

    removeTrigger?: boolean;

    title?: string;
    description?: string;
    confirmDelete?: boolean;
  };

  let {
    open = $bindable(false),

    removeTrigger,

    title,
    description,
    confirmDelete,

    class: className,
    children,
    onResult,
  }: Props & Modal<unknown> = $props();
</script>

<Dialog.Root bind:open>
  {#if !removeTrigger}
    <Dialog.Trigger class={className}>
      {@render children?.()}
    </Dialog.Trigger>
  {/if}

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
            <TriangleAlert size={18} class="text-white" />
          </div>
          <div>
            <Dialog.Title class="text-xl sm:text-2xl">
              <span
                class="bg-gradient-to-tr from-logo-1 via-logo-2 to-logo-3 bg-clip-text text-transparent"
              >
                {title ?? "Are you sure?"}
              </span>
            </Dialog.Title>
            {#if description}
              <Dialog.Description>
                {description}
              </Dialog.Description>
            {/if}
          </div>
        </div>
      </Dialog.Header>
    </div>

    <form
      class="flex flex-col gap-4"
      onsubmit={(e) => {
        e.preventDefault();
        onResult(true);
        open = false;
      }}
    >
      <Dialog.Footer>
        <Button
          variant="outline"
          onclick={() => {
            open = false;
          }}
        >
          Close
        </Button>
        <Button
          variant={confirmDelete ? "destructive" : "default"}
          type="submit"
        >
          {confirmDelete ? "Delete" : "Confirm"}
        </Button>
      </Dialog.Footer>
    </form>
  </Dialog.Content>
</Dialog.Root>
