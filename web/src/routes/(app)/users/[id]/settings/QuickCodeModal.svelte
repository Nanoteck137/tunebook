<script lang="ts">
  import { getApiClient, handleApiError } from "$lib";
  import Errors from "$lib/components/Errors.svelte";
  import FormItem from "$lib/components/FormItem.svelte";
  import { Button, Dialog, Input, Label } from "@nanoteck137/nano-ui";
  import { QrCode } from "lucide-svelte";
  import { zod } from "sveltekit-superforms/adapters";
  import { defaults, superForm } from "sveltekit-superforms/client";
  import { z } from "zod";
  import Spinner from "$lib/components/Spinner.svelte";
  import toast from "svelte-5-french-toast";

  const Schema = z.object({
    code: z.string().min(1),
  });

  export type Props = {
    open: boolean;
  };

  let { open = $bindable() }: Props = $props();
  const apiClient = getApiClient();

  let codeInput: HTMLInputElement | undefined = $state();

  $effect(() => {
    if (open) {
      reset({});
      codeInput?.focus();
    }
  });

  const { form, errors, enhance, reset, submitting } = superForm(
    defaults(zod(Schema)),
    {
      SPA: true,
      validators: zod(Schema),
      dataType: "json",
      resetForm: true,
      async onUpdate({ form, cancel }) {
        if (form.valid) {
          const formData = form.data;

          const res = await apiClient.authClaimQuickConnectCode({
            code: formData.code,
          });
          if (!res.success) {
            cancel();
            return handleApiError(res.error);
          }

          open = false;

          toast.success("Successfully logged in");
          reset({ data: {} });
        }
      },
    },
  );
</script>

<Dialog.Root bind:open>
  <Dialog.Content class="overflow-hidden sm:max-w-md">
    <div class="relative">
      <div
        class="absolute -right-16 -top-16 h-40 w-40 rounded-full bg-gradient-to-tr from-logo-1/10 via-logo-2/10 to-logo-3/10 blur-xl"
      ></div>

      <Dialog.Header class="relative text-left">
        <div class="flex items-center gap-3">
          <div class="flex h-10 w-10 shrink-0 items-center justify-center rounded-full bg-gradient-to-tr from-logo-1 via-logo-2 to-logo-3">
            <QrCode size={18} class="text-white" />
          </div>
          <div>
            <Dialog.Title class="text-xl sm:text-2xl">
              <span
                class="bg-gradient-to-tr from-logo-1 via-logo-2 to-logo-3 bg-clip-text text-transparent"
              >
                Enter QuickCode
              </span>
            </Dialog.Title>
            <Dialog.Description>
              Enter the QuickCode from the login screen
            </Dialog.Description>
          </div>
        </div>
      </Dialog.Header>
    </div>

    <form class="flex flex-col gap-4" use:enhance>
      <FormItem>
        <Label for="code">Code</Label>
        <Input
          id="code"
          name="code"
          type="text"
          bind:value={$form.code}
          autocomplete="off"
          placeholder="Enter code..."
          ref={codeInput}
        />
        <Errors errors={$errors.code} />
      </FormItem>

      <Dialog.Footer class="gap-2 sm:gap-0">
        <Button
          variant="outline"
          onclick={() => {
            open = false;
          }}
        >
          Close
        </Button>

        <Button type="submit" disabled={$submitting}>
          Sign In
          {#if $submitting}
            <Spinner />
          {/if}
        </Button>
      </Dialog.Footer>
    </form>
  </Dialog.Content>
</Dialog.Root>
