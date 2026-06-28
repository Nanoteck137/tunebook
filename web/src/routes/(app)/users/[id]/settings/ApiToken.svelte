<script lang="ts">
  import { invalidateAll } from "$app/navigation";
  import { getApiClient, handleApiError } from "$lib";
  import type { ApiToken as ApiTokenType } from "$lib/api/types";
  import ConfirmModal from "$lib/components/new-modals/ConfirmModal.svelte";
  import { Button, Dialog, Input, Label, Separator } from "@nanoteck137/nano-ui";
  import { Copy, Eye, Trash } from "lucide-svelte";
  import toast from "svelte-5-french-toast";

  type Props = {
    token: ApiTokenType;
  };

  const { token }: Props = $props();
  const apiClient = getApiClient();

  let openTokenShow = $state(false);
  let openDeleteModal = $state(false);

  let createdString = $derived(
    new Date(token.created).toLocaleDateString(),
  );
</script>

<div class="flex items-center justify-between py-3">
  <div class="flex flex-col">
    <span class="text-sm font-medium">{token.name}</span>
    <span class="text-xs text-muted-foreground">Created {createdString}</span>
  </div>
  <div class="flex gap-2">
    <Button
      variant="ghost"
      size="icon"
      onclick={() => {
        openTokenShow = true;
      }}
    >
      <Eye size={16} />
    </Button>

    <Button
      size="icon"
      variant="ghost"
      class="text-destructive hover:text-destructive"
      onclick={() => {
        openDeleteModal = true;
      }}
    >
      <Trash size={16} />
    </Button>
  </div>
</div>

<Separator />

<Dialog.Root bind:open={openTokenShow}>
  <Dialog.Content>
    <Dialog.Header>
      <Dialog.Title>{token.name}</Dialog.Title>
      <Dialog.Description>
        Created {createdString}
      </Dialog.Description>
    </Dialog.Header>

    <div class="flex flex-col gap-2">
      <Label for="token-value">Token</Label>
      <div class="flex gap-2">
        <Input
          id="token-value"
          class="font-mono"
          value={token.id}
          readonly
          onclick={(e) => {
            (e.target as HTMLInputElement).select();
          }}
        />
        <Button
          variant="outline"
          size="icon"
          onclick={() => {
            navigator.clipboard.writeText(token.id);
            toast.success("Copied to clipboard");
          }}
        >
          <Copy size={16} />
        </Button>
      </div>
    </div>

    <Dialog.Footer>
      <Button
        variant="destructive"
        onclick={() => {
          openTokenShow = false;
          openDeleteModal = true;
        }}
      >
        <Trash size={16} />
        Delete Token
      </Button>
    </Dialog.Footer>
  </Dialog.Content>
</Dialog.Root>

<ConfirmModal
  bind:open={openDeleteModal}
  removeTrigger
  title="Delete Token?"
  description="Are you sure you want to delete this token? This action cannot be undone."
  confirmDelete
  onResult={async () => {
    const res = await apiClient.deleteApiToken(token.id);
    if (!res.success) {
      return handleApiError(res.error);
    }

    toast.success("Successfully deleted api token");
    invalidateAll();
  }}
/>
