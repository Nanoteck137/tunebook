<script lang="ts">
  import { invalidateAll } from "$app/navigation";
  import { getApiClient, handleApiError } from "$lib";
  import Errors from "$lib/components/Errors.svelte";
  import FormItem from "$lib/components/FormItem.svelte";
  import { Button, Dialog, Label } from "@nanoteck137/nano-ui";
  import { ImageUp } from "lucide-svelte";
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
      async onUpdate({ form, cancel }) {
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
            cancel();
            return handleApiError(res.error);
          }

          open = false;

          toast.success("Successfully uploaded cover");
          invalidateAll();
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
            <ImageUp size={18} class="text-white" />
          </div>
          <div>
            <Dialog.Title class="text-xl sm:text-2xl">
              <span
                class="bg-gradient-to-tr from-logo-1 via-logo-2 to-logo-3 bg-clip-text text-transparent"
              >
                Upload Cover
              </span>
            </Dialog.Title>
            <Dialog.Description>
              Upload a new cover image for your playlist
            </Dialog.Description>
          </div>
        </div>
      </Dialog.Header>
    </div>

    <form
      class="flex flex-col gap-4"
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
          }}
        >
          Close
        </Button>

        <Button type="submit" disabled={$submitting}>
          Upload
          {#if $submitting}
            <Spinner />
          {/if}
        </Button>
      </Dialog.Footer>
    </form>
  </Dialog.Content>
</Dialog.Root>
