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
    cover: z.instanceof(FileList),
  });

  export type Props = {
    open: boolean;
    playlistId: string;
  };

  let { open = $bindable(), playlistId }: Props = $props();
  const apiClient = getApiClient();

  $effect(() => {
    if (open) {
      reset({});
    }
  });

  const { form, errors, enhance, reset, submitting } = superForm(
    defaults(zod(Schema)),
    {
      id: "upload-playlist-cover-modal",
      SPA: true,
      validators: zod(Schema),
      dataType: "json",
      resetForm: true,
      async onUpdate({ form }) {
        if (form.valid) {
          const formData = form.data;
          console.log("UPLOAD COVER", formData);

          const file = formData.cover.item(0);
          if (!file) {
            return;
          }

          const body = new FormData();
          body.append("image", file, file.name);

          const res = await apiClient.uploadPlaylistImage(playlistId, body);
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

    <form
      class="flex flex-col gap-4 px-[1px]"
      enctype="multipart/form-data"
      use:enhance
    >
      <FormItem>
        <Label for="cover">Cover File</Label>
        <input id="cover" name="cover" type="file" bind:files={$form.cover} />
        <Errors errors={$errors.cover} />
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
