import { z } from "zod";

export const sortTypes = [
	{
		label: "Name",
		value: "name-a-z",
		reverse: "name-z-a",
		direction: "asc",
	},
	{
		label: "Year",
		value: "year-new",
		reverse: "year-old",
		direction: "desc",
	},
	{
		label: "Added",
		value: "created-new",
		reverse: "created-old",
		direction: "desc",
	},
	{
		label: "Updated",
		value: "updated-new",
		reverse: "updated-old",
		direction: "desc",
	},
] as const;

const sortValues = [
	...sortTypes.map((t) => t.value),
	...sortTypes.map((t) => t.reverse),
] as const;

export const SortTypeEnum = z.enum(sortValues);
export type SortType = (typeof sortValues)[number];

export const defaultSort: SortType = "name-a-z";

export const AlbumFilter = z.object({
	sort: SortTypeEnum.default(defaultSort),
});
export type AlbumFilter = z.infer<typeof AlbumFilter>;

export function applySort(filter: AlbumFilter, query: Record<string, string>) {
	switch (filter.sort) {
		case "name-a-z":
			query["sort"] = "+name";
			break;
		case "name-z-a":
			query["sort"] = "-name";
			break;
		case "year-new":
			query["sort"] = "-year";
			break;
		case "year-old":
			query["sort"] = "+year";
			break;
		case "created-new":
			query["sort"] = "-created";
			break;
		case "created-old":
			query["sort"] = "+created";
			break;
		case "updated-new":
			query["sort"] = "-updated";
			break;
		case "updated-old":
			query["sort"] = "+updated";
			break;
	}
}
