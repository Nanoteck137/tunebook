import { error } from "@sveltejs/kit";
import type { PageLoad } from "./$types";

export const load: PageLoad = async ({ parent }) => {
	const data = await parent();

	const [mediaSettings, systemInfo, serverStats] = await Promise.all([
		data.apiClient.getMediaSettings(),
		data.apiClient.getSystemInfo(),
		data.apiClient.getSystemStats(),
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

	if (!serverStats.success) {
		throw error(serverStats.error.code, {
			message: serverStats.error.message,
			type: serverStats.error.type,
		});
	}

	return {
		mediaSettings: mediaSettings.data,
		systemInfo: systemInfo.data,
		serverStats: serverStats.data,
	};
};
