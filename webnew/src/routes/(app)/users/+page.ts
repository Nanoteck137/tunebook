import { error } from "@sveltejs/kit";
import type { PageLoad } from "./$types";

const USERS_PAGE_SIZE = 30;

export const load: PageLoad = async ({ parent }) => {
	const data = await parent();

	const users = await data.apiClient.searchUsers({
		query: { perPage: USERS_PAGE_SIZE.toString() },
	});
	if (!users.success) {
		throw error(users.error.code, { message: users.error.message });
	}

	return {
		...data,
		page: users.data.page,
		users: users.data.users,
	};
};
