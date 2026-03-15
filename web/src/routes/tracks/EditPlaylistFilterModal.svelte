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
  import type { TrackFilter } from "$lib/api/types";

  const Schema = z.object({
    name: z.string().min(1),
    filter: z.string().min(1),
  });

  export type Props = {
    open: boolean;
    filter: TrackFilter;
  };

  let { open = $bindable(), filter }: Props = $props();
  const apiClient = getApiClient();

  $effect(() => {
    if (open) {
      reset({
        data: {
          name: filter.name,
          filter: filter.filter,
        },
      });
    }
  });

  // TODO(patrik): Move to utils
  const toFieldError = (val: unknown): [string] | undefined =>
    typeof val === "string" && val ? [val] : undefined;

  const { form, errors, enhance, reset, submitting } = superForm(
    defaults(zod(Schema)),
    {
      id: "edit-playlist-filter-modal",
      SPA: true,
      validators: zod(Schema),
      dataType: "json",
      async onUpdate({ form, cancel }) {
        if (form.valid) {
          const formData = form.data;

          const res = await apiClient.editTrackFilter(filter.filterId, {
            name: formData.name,
            filter: formData.filter,
          });
          if (!res.success) {
            cancel();

            if (res.error.type === "VALIDATION_ERROR" && res.error.extra) {
              // TODO(patrik): Figure out a better way to handle errors
              // for the api client
              const rec = res.error.extra as Record<string, unknown>;
              $errors.name = toFieldError(rec["name"]);
              $errors.filter = toFieldError(rec["filter"]);

              return;
            }

            return handleApiError(res.error);
          }

          open = false;

          toast.success("Successfully updated playlist filter");
          invalidateAll();
        }
      },
    },
  );
</script>

<Dialog.Root bind:open>
  <Dialog.Content class="max-h-[420px] overflow-y-scroll">
    <Dialog.Header>
      <Dialog.Title>Edit playlist filter</Dialog.Title>
    </Dialog.Header>

    <form class="flex flex-col gap-4 px-[1px]" use:enhance>
      <FormItem>
        <Label for="name">Name</Label>
        <Input id="name" name="name" type="text" bind:value={$form.name} />
        <Errors errors={$errors.name} />
      </FormItem>

      <FormItem>
        <Label for="filter">Filter</Label>
        <Input
          id="filter"
          name="filter"
          type="text"
          bind:value={$form.filter}
        />
        <Errors errors={$errors.filter} />
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
