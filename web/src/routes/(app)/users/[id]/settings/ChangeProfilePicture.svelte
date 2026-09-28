<script lang="ts">
	import { invalidateAll } from "$app/navigation";
	import { getApiClient, handleApiError } from "$lib";
	import { Button, Input, Label, Separator } from "$lib/components/ui";
	import Spinner from "$lib/components/Spinner.svelte";
	import { toast } from "svelte-sonner";

	let { currentPicture }: { currentPicture: string } = $props();

	const apiClient = getApiClient();

	let selectedFile: File | undefined = $state();
	let uploading = $state(false);

	let pictureUrl = $state("");
	let settingUrl = $state(false);

	// Created in an effect rather than in a derived so the previous blob is
	// always revoked. Building it in a derived leaked one object URL per
	// recompute and never released any of them.
	let objectUrl: string | null = $state(null);

	$effect(() => {
		if (!selectedFile) {
			objectUrl = null;
			return;
		}

		const url = URL.createObjectURL(selectedFile);
		objectUrl = url;

		return () => URL.revokeObjectURL(url);
	});

	let previewUrl = $derived(
		objectUrl ?? (pictureUrl.trim() || currentPicture),
	);

	async function handleUpload() {
		if (!selectedFile) return;

		uploading = true;

		const formData = new FormData();
		formData.append("image", selectedFile);

		const res = await apiClient.uploadUserImage(formData);
		if (!res.success) {
			handleApiError(res.error);
			uploading = false;
			return;
		}

		toast.success("Profile picture updated");
		selectedFile = undefined;
		uploading = false;
		invalidateAll();
	}

	async function handleSetUrl() {
		if (!pictureUrl.trim()) return;

		settingUrl = true;

		const res = await apiClient.updateMe({ pictureUrl: pictureUrl.trim() });
		if (!res.success) {
			handleApiError(res.error);
			settingUrl = false;
			return;
		}

		toast.success("Profile picture updated");
		pictureUrl = "";
		settingUrl = false;
		invalidateAll();
	}
</script>

<div class="flex items-center gap-4">
	<img
		class="h-16 w-16 shrink-0 rounded-full object-cover ring-1 ring-border"
		src={previewUrl}
		alt=""
	/>

	<div class="flex min-w-0 flex-col gap-2">
		<Input
			type="file"
			accept="image/png,image/jpeg"
			class="max-w-sm"
			onchange={(e) => {
				const file = (e.target as HTMLInputElement).files?.[0];
				selectedFile = file;
				pictureUrl = "";
			}}
		/>

		<div>
			<Button onclick={handleUpload} disabled={!selectedFile || uploading}>
				{uploading ? "Uploading..." : "Upload"}
				{#if uploading}
					<Spinner />
				{/if}
			</Button>
		</div>
	</div>
</div>

<Separator class="my-1" />

<div class="flex flex-col gap-2">
	<Label for="pictureUrl">Or set from URL</Label>
	<div class="flex flex-col gap-2 sm:flex-row">
		<Input
			id="pictureUrl"
			type="url"
			placeholder="https://example.com/image.png"
			class="max-w-sm"
			bind:value={pictureUrl}
			oninput={() => {
				selectedFile = undefined;
			}}
		/>

		<Button
			variant="outline"
			onclick={handleSetUrl}
			disabled={!pictureUrl.trim() || settingUrl}
		>
			Set from URL
			{#if settingUrl}
				<Spinner />
			{/if}
		</Button>
	</div>
</div>
