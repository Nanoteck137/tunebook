<script lang="ts">
  import { invalidateAll } from "$app/navigation";
  import { getApiClient, handleApiError } from "$lib";
  import { Button } from "@nanoteck137/nano-ui";
  import toast from "svelte-5-french-toast";

  let { currentPicture }: { currentPicture: string } = $props();

  const apiClient = getApiClient();

  let selectedFile: File | undefined = $state();
  let uploading = $state(false);

  let previewUrl = $derived(
    selectedFile ? URL.createObjectURL(selectedFile) : currentPicture,
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
</script>

<div class="flex items-center gap-4">
  <img
    class="h-16 w-16 rounded-full object-cover"
    src={previewUrl}
    alt=""
  />

  <div class="flex flex-col gap-2">
    <input
      type="file"
      accept="image/png,image/jpeg"
      onchange={(e) => {
        const file = (e.target as HTMLInputElement).files?.[0];
        if (file) {
          selectedFile = file;
        }
      }}
    />

    <div>
      <Button
        onclick={handleUpload}
        disabled={!selectedFile || uploading}
      >
        {uploading ? "Uploading..." : "Upload"}
      </Button>
    </div>
  </div>
</div>
