import type { Page, UserData } from "$lib/api/types";
import type { PageLoad } from "./$types";

export const load: PageLoad = async ({ parent, url }) => {
  const data = await parent();

  const query = url.searchParams.get("query") ?? "";
  const page = url.searchParams.get("page") ?? "0";

  let users = [] as UserData[];
  let userPage: Page | null = null;

  if (query) {
    const res = await data.apiClient.searchUsers({
      query: { query, page },
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
