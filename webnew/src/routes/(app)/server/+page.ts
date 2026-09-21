import { isRoleAdmin } from "$lib/utils";
import { error, redirect } from "@sveltejs/kit";
import type { PageLoad } from "./$types";

export const load: PageLoad = async ({ parent }) => {
  const data = await parent();

  if (!isRoleAdmin(data.user?.role ?? "")) {
    redirect(301, "/");
  }

  const [mediaSettings, systemInfo, jobsResult] = await Promise.all([
    data.apiClient.getMediaSettings(),
    data.apiClient.getSystemInfo(),
    data.apiClient.getJobs(),
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

  if (!jobsResult.success) {
    throw error(jobsResult.error.code, {
      message: jobsResult.error.message,
      type: jobsResult.error.type,
    });
  }

  return {
    ...data,
    mediaSettings: mediaSettings.data,
    systemInfo: systemInfo.data,
    jobs: jobsResult.data.jobs,
  };
};
