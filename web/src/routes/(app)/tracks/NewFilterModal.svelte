<script lang="ts">
  import { invalidateAll } from "$app/navigation";
  import { getApiClient, handleApiError } from "$lib";
  import Errors from "$lib/components/Errors.svelte";
  import FormItem from "$lib/components/FormItem.svelte";
  import { Button, Dialog, Input, Label } from "@nanoteck137/nano-ui";
  import { ListFilter } from "lucide-svelte";
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

  let nameInput: HTMLInputElement | undefined = $state();

  const toFieldError = (val: unknown): [string] | undefined =>
    typeof val === "string" && val ? [val] : undefined;

  $effect(() => {
    if (open) {
      reset({});
      nameInput?.focus();
    }
  });

  const { form, errors, enhance, reset, submitting } = superForm(
    defaults(zod(Schema)),
    {
      id: "new-filter-modal",
      SPA: true,
      validators: zod(Schema),
      dataType: "json",
      resetForm: true,
      async onUpdate({ form, cancel }) {
        if (form.valid) {
          const formData = form.data;
          const res = await apiClient.createTrackFilter({
            name: formData.name,
            filter: formData.filter,
          });
          if (!res.success) {
            cancel();

            if (res.error.type === "VALIDATION_ERROR" && res.error.extra) {
              const rec = res.error.extra as Record<string, unknown>;
              if (rec["filter"]) {
                $errors.filter = toFieldError(rec["filter"]);
              }
              return;
            }

            return handleApiError(res.error);
          }

          open = false;

          toast.success("Successfully created new filter");
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
          <div
            class="flex h-10 w-10 shrink-0 items-center justify-center rounded-full bg-gradient-to-tr from-logo-1 via-logo-2 to-logo-3"
          >
            <ListFilter size={18} class="text-white" />
          </div>
          <div>
            <Dialog.Title class="text-xl sm:text-2xl">
              <span
                class="bg-gradient-to-tr from-logo-1 via-logo-2 to-logo-3 bg-clip-text text-transparent"
              >
                New Filter
              </span>
            </Dialog.Title>
            <Dialog.Description>
              Create a new filter to organize your tracks
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
          placeholder="e.g. Jazz Favorites, 80s Rock..."
          ref={nameInput}
        />
        <Errors errors={$errors.name} />
      </FormItem>

      <FormItem>
        <Label for="filter">Filter</Label>
        <Input
          id="filter"
          name="filter"
          type="text"
          bind:value={$form.filter}
          autocomplete="off"
          placeholder="e.g. genre == 'Jazz'"
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
          Create Filter
          {#if $submitting}
            <Spinner />
          {/if}
        </Button>
      </Dialog.Footer>
    </form>
  </Dialog.Content>
</Dialog.Root>
