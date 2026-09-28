import { error } from "@sveltejs/kit";
import type { PageLoad } from "./$types";

export const load: PageLoad = async ({ parent }) => {
	const data = await parent();

	const jobsResult = await data.apiClient.getJobs();
	if (!jobsResult.success) {
		throw error(jobsResult.error.code, {
			message: jobsResult.error.message,
			type: jobsResult.error.type,
		});
	}

	return {
		jobs: jobsResult.data.jobs,
	};
};
