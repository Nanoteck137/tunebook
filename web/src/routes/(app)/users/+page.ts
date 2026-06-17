import { getPagedQueryOptions } from "$lib/utils";
import { error } from "@sveltejs/kit";
import type { PageLoad } from "./$types";

export const load: PageLoad = async ({ parent, url }) => {
  const data = await parent();

  const query = getPagedQueryOptions(url.searchParams);

  const users = await data.apiClient.searchUsers({ query });
  if (!users.success) {
    throw error(users.error.code, { message: users.error.message });
  }

  return {
    ...data,
    page: users.data.page,
    users: users.data.users,
  };
};
