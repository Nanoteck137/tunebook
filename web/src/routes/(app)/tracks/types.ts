import { z } from "zod";

export const sortTypes = [
	{
		label: "Name",
		value: "name-a-z",
		reverse: "name-z-a",
		direction: "asc",
	},
	{
		label: "Artist",
		value: "artist",
		reverse: "artist-desc",
		direction: "asc",
	},
	{
		label: "Album",
		value: "album",
		reverse: "album-desc",
		direction: "asc",
	},
	{
		label: "Duration",
		value: "duration",
		reverse: "duration-desc",
		direction: "asc",
	},
	{ label: "Year", value: "year", reverse: "year-old", direction: "desc" },
	{
		label: "Added",
		value: "created-new",
		reverse: "created-old",
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

export const FullFilter = z.object({
	query: z.string(),
	sort: SortTypeEnum.default(defaultSort),
});
export type FullFilter = z.infer<typeof FullFilter>;

export function trackSortQuery(sort: SortType): string {
	switch (sort) {
		case "name-a-z":
			return "+name";
		case "name-z-a":
			return "-name";
		case "artist":
			return "+artist";
		case "artist-desc":
			return "-artist";
		case "album":
			return "+album";
		case "album-desc":
			return "-album";
		case "duration":
			return "+duration";
		case "duration-desc":
			return "-duration";
		case "year":
			return "-year";
		case "year-old":
			return "+year";
		case "created-new":
			return "-created";
		case "created-old":
			return "+created";
	}
}

export function constructFilterSort(
	filter: FullFilter,
	query: Record<string, string>,
) {
	const filters: string[] = [];

	if (filter.query !== "") {
		filters.push(`name contains "${filter.query}"`);
	}

	query["filter"] = filters.join(" and ");
	query["sort"] = trackSortQuery(filter.sort);
}
