<script lang="ts">
	import { getApiClient, handleApiError } from "$lib";
	import Errors from "$lib/components/Errors.svelte";
	import FormItem from "$lib/components/FormItem.svelte";
	import { Button, Dialog, Input, Label } from "$lib/components/ui";
	import { getQuickCodeModalManager } from "$lib/quick-code-modal.svelte";
	import { QrCode } from "@lucide/svelte";
	import { zod4 } from "sveltekit-superforms/adapters";
	import { defaults, superForm } from "sveltekit-superforms/client";
	import { z } from "zod";
	import Spinner from "$lib/components/Spinner.svelte";
	import { toast } from "svelte-sonner";

	const Schema = z.object({
		code: z.string().min(1),
	});

	const manager = getQuickCodeModalManager();
	const apiClient = getApiClient();

	let codeInput: HTMLInputElement | undefined = $state();

	$effect(() => {
		if (manager.open) {
			reset({});
			codeInput?.focus();
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

					const res = await apiClient.authClaimQuickConnectCode({
						code: formData.code,
					});
					if (!res.success) {
						cancel();
						return handleApiError(res.error);
					}

					manager.open = false;

					toast.success("Successfully logged in");
					reset({ data: {} });
				}
			},
		},
	);
</script>

<Dialog.Root bind:open={manager.open}>
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
						<QrCode size={18} class="text-white" />
					</div>
					<div>
						<Dialog.Title class="text-xl sm:text-2xl">
							<span
								class="bg-linear-to-tr from-logo-1 via-logo-2 to-logo-3 bg-clip-text text-transparent"
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

			<Dialog.Footer>
				<Button
					variant="outline"
					onclick={() => {
						manager.open = false;
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
