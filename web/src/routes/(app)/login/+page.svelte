<script lang="ts">
  import { invalidateAll } from "$app/navigation";
  import { getApiClient, handleApiError, setApiClientAuth } from "$lib";
  import type { AuthQuickConnectInitiate } from "$lib/api/types.js";
  import { getMusicManager } from "$lib/music-manager.svelte";
  import { Button } from "@nanoteck137/nano-ui";
  import { LogIn, QrCode, RefreshCw } from "lucide-svelte";
  import toast from "svelte-5-french-toast";

  const { data } = $props();
  const apiClient = getApiClient();
  const musicManager = getMusicManager();

  type LoginSuccess = {
    isSuccess: true;
    token: string;
  };
  type LoginError = {
    isSuccess: false;
    message: string;
  };
  type LoginResult = LoginSuccess | LoginError;

  async function loginWithPolling(providerId: string): Promise<LoginResult> {
    const res = await apiClient.authProviderInitiate({ providerId });
    if (!res.success) {
      handleApiError(res.error);
      return {
        isSuccess: false,
        message: `failed to initiate auth: ${res.error.message}`,
      };
    }

    const { requestId, challenge, expiresAt, authUrl } = res.data;

    const win = window.open(authUrl, "auth_window", "width=500,height=600");

    return new Promise<LoginResult>((resolve) => {
      const expiresAtDate = new Date(expiresAt);

      const pollInterval = setInterval(async () => {
        try {
          if (new Date() > expiresAtDate) {
            clearInterval(pollInterval);
            win?.close();
            resolve({ isSuccess: false, message: "authentication timeout" });
            return;
          }

          const res = await apiClient.authGetProviderStatus({
            requestId,
            challenge,
          });
          if (!res.success) {
            clearInterval(pollInterval);
            resolve({
              isSuccess: false,
              message: `authentication polling failed: ${res.error.message}`,
            });
            return;
          }

          if (res.data.status === "completed") {
            clearInterval(pollInterval);

            const res = await apiClient.authFinishProvider({
              requestId,
              challenge,
            });
            if (!res.success) {
              resolve({
                isSuccess: false,
                message: `authentication failed to get code: ${res.error.message}`,
              });
              return;
            }

            win?.close();

            resolve({
              isSuccess: true,
              token: res.data.token,
            });
          } else if (res.data.status === "failed") {
            clearInterval(pollInterval);
            resolve({
              isSuccess: false,
              message: "authentication failed for unknown reason",
            });
          } else if (res.data.status === "expired") {
            clearInterval(pollInterval);
            resolve({
              isSuccess: false,
              message: "authentication session expired",
            });
          } else if (res.data.status === "pending") {
            return;
          } else {
            clearInterval(pollInterval);
            resolve({
              isSuccess: false,
              message: "authentication failed for unknown reason",
            });
          }
        } catch (error) {
          clearInterval(pollInterval);
          console.error("auth catch error", error);
          resolve({
            isSuccess: false,
            message: "authentication failed for unknown reason",
          });
        }
      }, 2000);
    });
  }

  let quickAuth = $state<AuthQuickConnectInitiate | null>(null);

  async function initiateQuickCode() {
    const res = await apiClient.authQuickConnectInitiate();
    if (!res.success) {
      handleApiError(res.error);
      return;
    }

    quickAuth = res.data;
  }

  function resetQuickCode() {
    quickAuth = null;
  }

  $effect(() => {
    if (!quickAuth) return;

    const expiresAtDate = new Date(quickAuth.expiresAt);

    const pollInterval = setInterval(async () => {
      try {
        if (new Date() > expiresAtDate) {
          clearInterval(pollInterval);
          quickAuth = null;
          return;
        }

        if (!quickAuth) return;

        const res = await apiClient.authGetQuickConnectStatus({
          code: quickAuth.code,
          challenge: quickAuth.challenge,
        });
        if (!res.success) {
          clearInterval(pollInterval);
          return handleApiError(res.error);
        }

        console.log(res.data);

        if (res.data.status === "completed") {
          const res = await apiClient.authFinishQuickConnect({
            code: quickAuth.code,
            challenge: quickAuth.challenge,
          });
          if (!res.success) {
            clearInterval(pollInterval);
            quickAuth = null;
            return handleApiError(res.error);
          }

          clearInterval(pollInterval);

          localStorage.setItem("token", res.data.token);
          setApiClientAuth(apiClient, res.data.token);

          musicManager.initQueue();
          invalidateAll();
        } else if (res.data.status === "pending") {
          return;
        } else if (res.data.status === "expired") {
          clearInterval(pollInterval);
          quickAuth = null;
        } else {
          clearInterval(pollInterval);
          quickAuth = null;
        }
      } catch (error) {
        clearInterval(pollInterval);
        console.error("quick auth catch error", error);
      }
    }, 2000);

    return () => {
      clearInterval(pollInterval);
    };
  });
</script>

<svelte:head>
  <title>Login - Tunebook</title>
</svelte:head>

<div class="mt-16 flex flex-col items-center justify-center gap-8 p-4">
  <div
    class="flex h-32 w-32 items-center justify-center rounded-full bg-gradient-to-tr from-logo-1 via-logo-2 to-logo-3"
  >
    <svg
      xmlns="http://www.w3.org/2000/svg"
      class="h-16 w-16 text-black"
      viewBox="0 0 24 24"
      fill="none"
      stroke="currentColor"
      stroke-width="2"
      stroke-linecap="round"
      stroke-linejoin="round"
    >
      <circle cx="12" cy="12" r="10" /><circle cx="12" cy="12" r="3" />
    </svg>
  </div>
  <h1
    class="bg-gradient-to-tr from-logo-1 via-logo-2 to-logo-3 bg-clip-text text-5xl font-bold text-transparent"
  >
    Tunebook
  </h1>
  <p class="text-muted-foreground">Your personal music streaming server</p>

  <div class="flex flex-col gap-3">
    {#each data.providers as provider (provider.id)}
      <Button
        class="h-12 w-full justify-start gap-3 px-4 text-base"
        variant="outline"
        onclick={async () => {
          const res = await loginWithPolling(provider.id);
          if (!res.isSuccess) {
            toast.error(`login failed: ${res.message}`);
            return;
          }

          localStorage.setItem("token", res.token);
          setApiClientAuth(apiClient, res.token);
          musicManager.initQueue();
          invalidateAll();
        }}
      >
        <LogIn size={18} />
        Login with {provider.displayName}
      </Button>
    {/each}
  </div>

  <div class="flex items-center gap-3">
    <div class="h-px w-16 bg-border"></div>
    <span class="text-xs text-muted-foreground">or</span>
    <div class="h-px w-16 bg-border"></div>
  </div>

  {#if quickAuth}
    <div class="flex flex-col items-center gap-4">
      <div
        class="flex flex-col items-center gap-2 rounded-lg border bg-card px-8 py-6"
      >
        <QrCode size={20} class="text-muted-foreground" />
        <p class="text-xs text-muted-foreground">
          Enter this code on another device
        </p>
        <p class="text-4xl font-bold tracking-widest">{quickAuth.code}</p>
      </div>
      <div class="flex items-center gap-2">
        <RefreshCw size={14} class="animate-spin text-muted-foreground" />
        <p class="text-xs text-muted-foreground">Waiting for approval...</p>
      </div>
      <Button variant="ghost" size="sm" onclick={resetQuickCode}>
        Cancel
      </Button>
    </div>
  {:else}
    <Button variant="ghost" onclick={initiateQuickCode}>
      <QrCode size={16} />
      Login with Quick Code
    </Button>
  {/if}
</div>
