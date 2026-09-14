import { isRoleAdmin } from "$lib/utils";
import { error, redirect } from "@sveltejs/kit";
import type { PageLoad } from "./$types";

export const load: PageLoad = async ({ parent }) => {
  const data = await parent();

  if (!isRoleAdmin(data.user?.role ?? "")) {
    redirect(301, "/");
  }

  const [mediaSettings, systemInfo] = await Promise.all([
    data.apiClient.getMediaSettings(),
    data.apiClient.getSystemInfo(),
  ]);

  if (!mediaSettings.success) {
    throw error(mediaSettings.error.code, {
      message: mediaSettings.error.message,
      type: mediaSettings.error.type,
    });
  }

  if (!systemInfo.success) {
    throw error(systemInfo.error.code, {
      message: systemInfo.error.message,
      type: systemInfo.error.type,
    });
  }

  return {
    ...data,
    mediaSettings: mediaSettings.data,
    systemInfo: systemInfo.data,
  };
};
