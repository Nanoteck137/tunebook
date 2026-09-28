<script lang="ts">
	import { Button } from "$lib/components/ui";
	import { showQuickCodeModal } from "$lib/quick-code-modal.svelte";
	import { QrCode } from "@lucide/svelte";
</script>

<div class="flex flex-wrap items-center justify-between gap-3">
	<p class="max-w-md text-sm text-muted-foreground">
		Sign in on another device without a password. Start from the login screen
		there, then enter the code it generates.
	</p>

	<Button variant="outline" size="sm" onclick={showQuickCodeModal}>
		<QrCode size={14} />
		Enter Code
	</Button>
</div>
