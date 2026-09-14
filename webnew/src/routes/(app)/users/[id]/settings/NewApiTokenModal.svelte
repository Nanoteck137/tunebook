<script lang="ts">
	import { invalidateAll } from "$app/navigation";
	import { getApiClient, handleApiError } from "$lib";
	import Errors from "$lib/components/Errors.svelte";
	import FormItem from "$lib/components/FormItem.svelte";
	import { Button, Dialog, Input, Label } from "$lib/components/ui";
	import { Key } from "@lucide/svelte";
	import { zod4 } from "sveltekit-superforms/adapters";
	import { defaults, superForm } from "sveltekit-superforms/client";
	import { z } from "zod";
	import Spinner from "$lib/components/Spinner.svelte";
	import { toast } from "svelte-sonner";

	const Schema = z.object({
		name: z.string().min(1),
	});

	export type Props = {
		open: boolean;
	};

	let { open = $bindable() }: Props = $props();
	const apiClient = getApiClient();

	let nameInput: HTMLInputElement | undefined = $state();

	$effect(() => {
		if (open) {
			reset({});
			nameInput?.focus();
		}
	});

	const { form, errors, enhance, reset, submitting } = superForm(
		defaults(zod4(Schema)),
		{
			SPA: true,
			validators: zod4(Schema),
			dataType: "json",
			resetForm: true,
			async onUpdate({ form, cancel }) {
				if (form.valid) {
					const formData = form.data;
					const res = await apiClient.createApiToken({
						name: formData.name,
					});
					if (!res.success) {
						cancel();
						return handleApiError(res.error);
					}

					open = false;

					toast.success("Successfully created new API token");
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
				class="absolute -top-16 -right-16 h-40 w-40 rounded-full bg-linear-to-tr from-logo-1/10 via-logo-2/10 to-logo-3/10 blur-xl"
			></div>

			<Dialog.Header class="relative text-left">
				<div class="flex items-center gap-3">
					<div
						class="flex h-10 w-10 shrink-0 items-center justify-center rounded-full bg-linear-to-tr from-logo-1 via-logo-2 to-logo-3"
					>
						<Key size={18} class="text-white" />
					</div>
					<div>
						<Dialog.Title class="text-xl sm:text-2xl">
							<span
								class="bg-linear-to-tr from-logo-1 via-logo-2 to-logo-3 bg-clip-text text-transparent"
							>
								New API Token
							</span>
						</Dialog.Title>
						<Dialog.Description>
							Create a new API token for external access
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
					placeholder="e.g. My App, Script..."
					ref={nameInput}
				/>
				<Errors errors={$errors.name} />
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
					Create Token
					{#if $submitting}
						<Spinner />
					{/if}
				</Button>
			</Dialog.Footer>
		</form>
	</Dialog.Content>
</Dialog.Root>
