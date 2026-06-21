<script lang="ts">
  import { invalidateAll } from "$app/navigation";
  import { getApiClient, handleApiError, setApiClientAuth } from "$lib";
  import type { AuthQuickConnectInitiate } from "$lib/api/types.js";

  // const {} = $props();
  const apiClient = getApiClient();

  let auth = $state<AuthQuickConnectInitiate | null>(null);

  async function initiate() {
    const res = await apiClient.authQuickConnectInitiate();
    if (!res.success) {
      return handleApiError(res.error);
    }

    auth = res.data;
    console.log(auth);
  }

  $effect(() => {
    if (!auth) {
      initiate();
    }
  });

  $effect(() => {
    if (!auth) {
      return;
    }

    const expiresAtDate = new Date(auth.expiresAt);

    const pollInterval = setInterval(async () => {
      try {
        if (new Date() > expiresAtDate) {
          clearInterval(pollInterval);

          auth = null;

          return;
        }

        if (!auth) return;

        const res = await apiClient.authGetQuickConnectStatus({
          code: auth.code,
          challenge: auth.challenge,
        });
        if (!res.success) {
          clearInterval(pollInterval);
          // resolve({
          //   isSuccess: false,
          //   message: `authentication polling failed: ${res.error.message}`,
          // });
          return handleApiError(res.error);
        }

        if (res.data.status === "completed") {
          const res = await apiClient.authFinishQuickConnect({
            code: auth.code,
            challenge: auth.challenge,
          });
          if (!res.success) {
            clearInterval(pollInterval);
            auth = null;

            return handleApiError(res.error);
          }

          clearInterval(pollInterval);

          localStorage.setItem("token", res.data.token);
          setApiClientAuth(apiClient, res.data.token);

          invalidateAll();
          // } else if (res.data.status === "pending") {
        } else if (res.data.status === "expired") {
          clearInterval(pollInterval);
          auth = null;
        } else {
          clearInterval(pollInterval);
          auth = null;
        }
      } catch (error) {
        clearInterval(pollInterval);
        console.error("auth catch error", error);
        // resolve({
        //   isSuccess: false,
        //   message: `authentication failed for unknown reason`,
        // });
      }
    }, 2000);

    return () => {
      clearInterval(pollInterval);
    };
  });
</script>

<p>Code: {auth?.code}</p>
