<script lang="ts">
	import { invalidateAll } from "$app/navigation";
	import { getApiClient, handleApiError } from "$lib";
	import type { ApiToken as ApiTokenType } from "$lib/api/types";
	import ConfirmModal from "$lib/components/new-modals/ConfirmModal.svelte";
	import { Button } from "$lib/components/ui";
	import { Copy, Eye, EyeOff, Trash } from "@lucide/svelte";
	import { toast } from "svelte-sonner";

	type Props = {
		token: ApiTokenType;
	};

	const { token }: Props = $props();
	const apiClient = getApiClient();

	let revealed = $state(false);
	let openDeleteModal = $state(false);

	let createdString = $derived(new Date(token.created).toLocaleDateString());

	let masked = $derived(
		token.id.length <= 12
			? "•".repeat(token.id.length)
			: `${token.id.slice(0, 6)}${"•".repeat(12)}${token.id.slice(-4)}`,
	);

	async function copyToken() {
		try {
			await navigator.clipboard.writeText(token.id);
			toast.success("Copied to clipboard");
		} catch {
			toast.error("Failed to copy to clipboard");
		}
	}
</script>

<div class="flex flex-wrap items-center justify-between gap-3 py-3">
	<div class="flex min-w-0 flex-col">
		<span class="truncate text-sm font-medium">{token.name}</span>
		<span class="text-xs text-muted-foreground">Created {createdString}</span>
	</div>

	<div class="flex shrink-0 items-center gap-1">
		<button
			type="button"
			class="flex h-8 max-w-56 items-center gap-1.5 rounded-md px-2 font-mono text-xs text-muted-foreground transition-colors hover:bg-accent hover:text-foreground focus-visible:ring-2 focus-visible:ring-ring focus-visible:outline-none"
			onclick={() => (revealed = !revealed)}
			title={revealed ? "Hide token" : "Reveal token"}
		>
			{#if revealed}
				<EyeOff size={14} class="shrink-0" />
			{:else}
				<Eye size={14} class="shrink-0" />
			{/if}
			<span class="truncate">{revealed ? token.id : masked}</span>
		</button>

		<Button
			variant="ghost"
			size="icon"
			title="Copy token"
			aria-label="Copy token"
			onclick={copyToken}
		>
			<Copy size={16} />
		</Button>

		<Button
			size="icon"
			variant="ghost"
			class="text-destructive hover:text-destructive"
			title="Delete token"
			aria-label="Delete token"
			onclick={() => {
				openDeleteModal = true;
			}}
		>
			<Trash size={16} />
		</Button>
	</div>
</div>

<ConfirmModal
	bind:open={openDeleteModal}
	removeTrigger
	title="Delete Token?"
	description={`"${token.name}" will stop working immediately. This action cannot be undone.`}
	confirmDelete
	onResult={async () => {
		const res = await apiClient.deleteApiToken(token.id);
		if (!res.success) {
			return handleApiError(res.error);
		}

		toast.success("Successfully deleted api token");
		invalidateAll();
	}}
/>
