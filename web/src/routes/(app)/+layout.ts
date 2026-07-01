import { getApiAddress, setApiClientAuth } from "$lib";
import { ApiClient } from "$lib/api/client";
import type { GetMe } from "$lib/api/types";
import { redirect } from "@sveltejs/kit";
import type { LayoutLoad } from "./$types";

export const prerender = false;
export const ssr = false;

export const load: LayoutLoad = async ({ url }) => {
  const apiClient = new ApiClient(getApiAddress(url));
  const token = localStorage.getItem("token") ?? undefined;
  setApiClientAuth(apiClient, token);

  let user: GetMe | null = null;
  if (token) {
    const res = await apiClient.getMe();
    if (!res.success) {
      console.error("Get Me API Error", res.error.message);
      user = null;
    } else {
      user = res.data;
    }
  }

  if (!user && !url.pathname.startsWith("/login")) {
    throw redirect(303, "/login");
  }

  return {
    apiClient,
    user,
  };
};
