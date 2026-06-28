<script lang="ts">
  import { Button, Card, Separator } from "@nanoteck137/nano-ui";
  import ChangeDisplayName from "./ChangeDisplayName.svelte";
  import ApiToken from "./ApiToken.svelte";
  import NewApiTokenModal from "./NewApiTokenModal.svelte";
  import { Plus, QrCode } from "lucide-svelte";
  import QuickCodeModal from "./QuickCodeModal.svelte";

  let { data } = $props();

  let openNewApiTokenModal = $state(false);
  let openQuickCodeModal = $state(false);
</script>

<div class="flex flex-col gap-8">
  <div>
    <h1 class="text-xl font-bold">Settings</h1>
  </div>

  <Card.Root>
    <div class="p-6">
      <h2 class="text-lg font-semibold">Display Name</h2>
      <p class="mt-1 text-sm text-muted-foreground">
        Change how your name appears across Tunebook.
      </p>
      <Separator class="my-4" />
      <ChangeDisplayName />
    </div>
  </Card.Root>

  <Card.Root>
    <div class="p-6">
      <div class="flex items-center justify-between">
        <div>
          <h2 class="text-lg font-semibold">Quick Code</h2>
          <p class="mt-1 text-sm text-muted-foreground">
            Claim a quick connect code to log in on another device.
          </p>
        </div>
        <Button
          variant="outline"
          onclick={() => {
            openQuickCodeModal = true;
          }}
        >
          <QrCode />
          Enter Code
        </Button>
      </div>
    </div>
  </Card.Root>

  <Card.Root>
    <div class="p-6">
      <div class="flex items-center justify-between">
        <div>
          <h2 class="text-lg font-semibold">API Tokens</h2>
          <p class="mt-1 text-sm text-muted-foreground">
            Manage API tokens for programmatic access.
          </p>
        </div>
        <Button
          variant="outline"
          onclick={() => {
            openNewApiTokenModal = true;
          }}
        >
          <Plus />
          New Token
        </Button>
      </div>
      <Separator class="my-4" />
      <div class="flex flex-col">
        {#each data.tokens as token (token.id)}
          <ApiToken {token} />
        {/each}
      </div>
    </div>
  </Card.Root>
</div>

<NewApiTokenModal bind:open={openNewApiTokenModal} />
<QuickCodeModal bind:open={openQuickCodeModal} />
