<script lang="ts">
  import { invalidateAll } from "$app/navigation";
  import { getApiClient, handleApiError } from "$lib";
  import Errors from "$lib/components/Errors.svelte";
  import FormItem from "$lib/components/FormItem.svelte";
  import { Button, Dialog, Input, Label } from "@nanoteck137/nano-ui";
  import toast from "svelte-5-french-toast";
  import { zod } from "sveltekit-superforms/adapters";
  import { defaults, superForm } from "sveltekit-superforms/client";
  import { z } from "zod";
  import Spinner from "$lib/components/Spinner.svelte";
  import type { Playlist } from "$lib/api/types";

  const Schema = z.object({
    name: z.string().min(1),
  });

  export type Props = {
    open: boolean;
    playlist: Playlist;
  };

  let { open = $bindable(), playlist }: Props = $props();
  const apiClient = getApiClient();

  $effect(() => {
    if (open) {
      reset({
        data: {
          name: playlist.name,
        },
      });
    }
  });

  const { form, errors, enhance, reset, submitting } = superForm(
    defaults(zod(Schema)),
    {
      id: "edit-playlist-modal",
      SPA: true,
      validators: zod(Schema),
      dataType: "json",
      resetForm: true,
      async onUpdate({ form }) {
        if (form.valid) {
          const formData = form.data;

          const res = await apiClient.editPlaylist(playlist.id, {
            name: formData.name,
          });
          if (!res.success) {
            return handleApiError(res.error);
          }

          open = false;

          toast.success("Successfully updated playlist");
          invalidateAll();
        }
      },
    },
  );
</script>

<Dialog.Root bind:open>
  <Dialog.Content class="max-h-[420px] overflow-y-scroll">
    <Dialog.Header>
      <Dialog.Title>Edit playlist</Dialog.Title>
    </Dialog.Header>

    <form class="flex flex-col gap-4 px-[1px]" use:enhance>
      <FormItem>
        <Label for="name">Name</Label>
        <Input id="name" name="name" type="text" bind:value={$form.name} />
        <Errors errors={$errors.name} />
      </FormItem>

      <Dialog.Footer class="gap-2 sm:gap-0">
        <Button
          variant="outline"
          onclick={() => {
            open = false;
            reset();
          }}
        >
          Close
        </Button>

        <Button type="submit" disabled={$submitting}>
          Save
          {#if $submitting}
            <Spinner />
          {/if}
        </Button>
      </Dialog.Footer>
    </form>
  </Dialog.Content>
</Dialog.Root>
