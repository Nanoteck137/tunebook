<script lang="ts">
	import { invalidateAll } from "$app/navigation";
	import { getApiClient, handleApiError } from "$lib";
	import Errors from "$lib/components/Errors.svelte";
	import FormItem from "$lib/components/FormItem.svelte";
	import { Button, Input, Label } from "$lib/components/ui";
	import { zod4 } from "sveltekit-superforms/adapters";
	import { defaults, superForm } from "sveltekit-superforms/client";
	import { z } from "zod";
	import { untrack } from "svelte";
	import Spinner from "$lib/components/Spinner.svelte";
	import { toast } from "svelte-sonner";

	let { currentName }: { currentName: string } = $props();

	const Schema = z.object({
		displayName: z.string().min(1),
	});

	const apiClient = getApiClient();

	// Seeded once on mount. Reseeding on every revalidation would discard
	// whatever the user has typed but not yet submitted.
	const initialName = untrack(() => currentName);

	const f = superForm(defaults({ displayName: initialName }, zod4(Schema)), {
		id: "change-display-name",
		SPA: true,
		validators: zod4(Schema),
		dataType: "json",
		async onUpdate({ form }) {
			if (form.valid) {
				const formData = form.data;
				const res = await apiClient.updateMe({
					displayName: formData.displayName,
				});
				if (!res.success) {
					return handleApiError(res.error);
				}

				toast.success("Successfully changed display name");
				invalidateAll();
			}
		},
	});
	const { form, errors, enhance, submitting } = f;
</script>

<form class="flex flex-col gap-4" use:enhance>
	<FormItem>
		<Label for="displayName">Display Name</Label>
		<Input
			id="displayName"
			name="displayName"
			type="text"
			class="max-w-sm"
			bind:value={$form.displayName}
		/>
		<Errors errors={$errors.displayName} />
	</FormItem>

	<div>
		<Button
			type="submit"
			disabled={$submitting || $form.displayName.trim() === currentName}
		>
			Update
			{#if $submitting}
				<Spinner />
			{/if}
		</Button>
	</div>
</form>
