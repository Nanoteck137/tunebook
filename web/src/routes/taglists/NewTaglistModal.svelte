<script lang="ts">
  import { goto } from "$app/navigation";
  import { getApiClient, handleApiError } from "$lib";
  import Errors from "$lib/components/Errors.svelte";
  import FormItem from "$lib/components/FormItem.svelte";
  import {
    Button,
    Dialog,
    Input,
    Label,
    Select,
    Separator,
  } from "@nanoteck137/nano-ui";
  import toast from "svelte-5-french-toast";
  import { zod } from "sveltekit-superforms/adapters";
  import { defaults, superForm } from "sveltekit-superforms/client";
  import { z } from "zod";
  import Spinner from "$lib/components/Spinner.svelte";

  const Schema = z.object({
    name: z.string().min(1, "Name cannot be empty"),
    filter: z.string().min(1, "Filter cannot be empty"),
  });

  export type Props = {
    open: boolean;
  };

  let { open = $bindable() }: Props = $props();
  const apiClient = getApiClient();

  $effect(() => {
    if (open) {
      reset({});
    }
  });

  const { form, errors, enhance, reset, submitting } = superForm(
    defaults(zod(Schema)),
    {
      SPA: true,
      validators: zod(Schema),
      dataType: "json",
      resetForm: true,
      async onUpdate({ form }) {
        if (form.valid) {
          const formData = form.data;
          const res = await apiClient.createTaglist({
            name: formData.name,
            filter: formData.filter,
          });
          // TODO(patrik): Handle filter errors
          if (!res.success) {
            return handleApiError(res.error);
          }

          open = false;

          toast.success("Successfully created new taglist");
          goto(`/taglists/${res.data.id}`, { invalidateAll: true });
        }
      },
    },
  );
</script>

<Dialog.Root bind:open>
  <Dialog.Content>
    <Dialog.Header>
      <Dialog.Title>Create new playlist</Dialog.Title>
    </Dialog.Header>

    <form class="flex flex-col gap-4" use:enhance>
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
          }}
        >
          Close
        </Button>

        <Button type="submit" disabled={$submitting}>
          Create
          {#if $submitting}
            <Spinner />
          {/if}
        </Button>
      </Dialog.Footer>
    </form>
  </Dialog.Content>
</Dialog.Root>
