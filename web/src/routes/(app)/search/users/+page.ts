import type { Page, UserData } from "$lib/api/types";
import type { PageLoad } from "./$types";

const SEARCH_USERS_PAGE_SIZE = 30;

export const load: PageLoad = async ({ parent, url }) => {
	const data = await parent();

	const query = url.searchParams.get("query") ?? "";

	let users = [] as UserData[];
	let userPage: Page | null = null;

	if (query) {
		const res = await data.apiClient.searchUsers({
			query: { query, perPage: SEARCH_USERS_PAGE_SIZE.toString() },
		});

		if (res.success) {
			users = res.data.users;
			userPage = res.data.page;
		}
	}

	return {
		...data,
		query,
		users,
		page: userPage,
	};
};
