<script lang="ts">
  import { invalidateAll } from "$app/navigation";
  import { getApiClient, handleApiError } from "$lib";
  import Errors from "$lib/components/Errors.svelte";
  import FormItem from "$lib/components/FormItem.svelte";
  import { Button, Dialog, Input, Label } from "@nanoteck137/nano-ui";
  import { ListMusic } from "lucide-svelte";
  import toast from "svelte-5-french-toast";
  import { zod } from "sveltekit-superforms/adapters";
  import { defaults, superForm } from "sveltekit-superforms/client";
  import { z } from "zod";
  import Spinner from "$lib/components/Spinner.svelte";
  import type { Playlist } from "$lib/api/types";

  const Schema = z.object({
    name: z.string().min(1),
    coverUrl: z.string(),
  });

  export type Props = {
    open: boolean;
    playlist: Playlist;
  };

  let { open = $bindable(), playlist }: Props = $props();
  const apiClient = getApiClient();

  let nameInput: HTMLInputElement | undefined = $state();

  $effect(() => {
    if (open) {
      reset({
        data: {
          name: playlist.name,
        },
      });
      nameInput?.focus();
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
      async onUpdate({ form, cancel }) {
        if (form.valid) {
          const formData = form.data;

          const res = await apiClient.editPlaylist(playlist.id, {
            name: formData.name,
            coverUrl: formData.coverUrl !== "" ? formData.coverUrl : null,
          });
          if (!res.success) {
            cancel();
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
  <Dialog.Content class="overflow-hidden sm:max-w-md">
    <div class="relative">
      <div
        class="absolute -right-16 -top-16 h-40 w-40 rounded-full bg-gradient-to-tr from-logo-1/10 via-logo-2/10 to-logo-3/10 blur-xl"
      ></div>

      <Dialog.Header class="relative text-left">
        <div class="flex items-center gap-3">
          <div class="flex h-10 w-10 shrink-0 items-center justify-center rounded-full bg-gradient-to-tr from-logo-1 via-logo-2 to-logo-3">
            <ListMusic size={18} class="text-white" />
          </div>
          <div>
            <Dialog.Title class="text-xl sm:text-2xl">
              <span
                class="bg-gradient-to-tr from-logo-1 via-logo-2 to-logo-3 bg-clip-text text-transparent"
              >
                Edit Playlist
              </span>
            </Dialog.Title>
            <Dialog.Description>
              Update your playlist details
            </Dialog.Description>
          </div>
        </div>
      </Dialog.Header>
    </div>

    <form class="flex flex-col gap-4" use:enhance>
      <FormItem>
        <Label for="name">Name</Label>
        <Input
          id="name"
          name="name"
          type="text"
          bind:value={$form.name}
          autocomplete="off"
          placeholder="e.g. Chill Vibes, Road Trip..."
          ref={nameInput}
        />
        <Errors errors={$errors.name} />
      </FormItem>

      <FormItem>
        <Label for="coverUrl">Cover URL</Label>
        <Input
          id="coverUrl"
          name="coverUrl"
          type="text"
          bind:value={$form.coverUrl}
          autocomplete="off"
          placeholder="https://example.com/cover.jpg"
        />
        <Errors errors={$errors.coverUrl} />
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
          Save Changes
          {#if $submitting}
            <Spinner />
          {/if}
        </Button>
      </Dialog.Footer>
    </form>
  </Dialog.Content>
</Dialog.Root>
