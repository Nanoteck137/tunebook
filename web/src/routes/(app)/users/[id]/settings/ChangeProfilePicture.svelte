<script lang="ts">
  import { invalidateAll } from "$app/navigation";
  import { getApiClient, handleApiError } from "$lib";
  import Errors from "$lib/components/Errors.svelte";
  import FormItem from "$lib/components/FormItem.svelte";
  import { Button, Input, Label, Separator } from "@nanoteck137/nano-ui";
  import toast from "svelte-5-french-toast";
  import Spinner from "$lib/components/Spinner.svelte";

  let { currentPicture }: { currentPicture: string } = $props();

  const apiClient = getApiClient();

  let selectedFile: File | undefined = $state();
  let uploading = $state(false);

  let pictureUrl = $state("");
  let settingUrl = $state(false);
  let urlError = $state("");

  let previewUrl = $derived(
    selectedFile
      ? URL.createObjectURL(selectedFile)
      : pictureUrl
        ? pictureUrl
        : currentPicture,
  );

  async function handleUpload() {
    if (!selectedFile) return;

    uploading = true;
    urlError = "";

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
    if (!pictureUrl) return;

    settingUrl = true;
    urlError = "";

    const res = await apiClient.updateMe({ pictureUrl });
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
          pictureUrl = "";
        }
      }}
    />

    <div>
      <Button
        onclick={handleUpload}
        disabled={!selectedFile || uploading}
      >
        {uploading ? "Uploading..." : "Upload"}
        {#if uploading}
          <Spinner />
        {/if}
      </Button>
    </div>
  </div>
</div>

<Separator class="my-4" />

<div class="flex flex-col gap-4">
  <FormItem>
    <Label for="pictureUrl">Image URL</Label>
    <Input
      id="pictureUrl"
      type="url"
      placeholder="https://example.com/image.png"
      class="max-w-sm"
      bind:value={pictureUrl}
      oninput={() => {
        selectedFile = undefined;
        urlError = "";
      }}
    />
    <Errors errors={urlError ? [urlError] : undefined} />
  </FormItem>

  <div>
    <Button
      onclick={handleSetUrl}
      disabled={!pictureUrl || settingUrl}
    >
      Set from URL
      {#if settingUrl}
        <Spinner />
      {/if}
    </Button>
  </div>
</div>
