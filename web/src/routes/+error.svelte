<script>
  import { goto, invalidateAll } from "$app/navigation";
  import { page } from "$app/state";
  import { Button } from "@nanoteck137/nano-ui";
</script>

<div class="flex min-h-[80vh] flex-col items-center justify-center gap-6 p-4">
  <div
    class="flex h-28 w-28 items-center justify-center rounded-full bg-gradient-to-tr from-logo-1 via-logo-2 to-logo-3"
  >
    <span class="text-4xl font-bold text-black">{page.status}</span>
  </div>

  <div class="flex flex-col items-center gap-2">
    <h1 class="text-2xl font-bold">
      {#if page.status === 404}
        Page not found
      {:else if page.status === 400}
        Something went wrong
      {:else if page.status === 500}
        Server error
      {:else}
        Error
      {/if}
    </h1>
    <p class="max-w-md text-center text-muted-foreground">
      {page.error?.message ?? "An unexpected error occurred."}
    </p>
  </div>

  <div class="flex items-center gap-3">
    <Button href="/">Go Home</Button>
    {#if page.status === 400}
      <Button
        variant="outline"
        onclick={() => {
          localStorage.removeItem("token");
          invalidateAll();
          goto("/");
        }}
      >
        Logout
      </Button>
    {/if}
  </div>
</div>
