import { isRoleAdmin } from "$lib/utils";
import { error, redirect } from "@sveltejs/kit";
import type { PageLoad } from "./$types";

export const load: PageLoad = async ({ parent }) => {
  const data = await parent();

  if (!isRoleAdmin(data.user?.role ?? "")) {
    redirect(301, "/");
  }

  const mediaSettings = await data.apiClient.getMediaSettings();
  if (!mediaSettings.success) {
    throw error(mediaSettings.error.code, {
      message: mediaSettings.error.message,
      type: mediaSettings.error.type,
    });
  }

  return {
    ...data,
    mediaSettings: mediaSettings.data,
  };
};
