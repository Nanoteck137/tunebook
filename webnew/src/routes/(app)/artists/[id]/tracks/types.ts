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

export const TrackFilter = z.object({
	sort: SortTypeEnum.default(defaultSort),
});
export type TrackFilter = z.infer<typeof TrackFilter>;

export function applySort(filter: TrackFilter, query: Record<string, string>) {
	switch (filter.sort) {
		case "name-a-z":
			query["sort"] = "+name";
			break;
		case "name-z-a":
			query["sort"] = "-name";
			break;
		case "artist":
			query["sort"] = "+artistName";
			break;
		case "artist-desc":
			query["sort"] = "-artistName";
			break;
		case "album":
			query["sort"] = "+albumName";
			break;
		case "album-desc":
			query["sort"] = "-albumName";
			break;
		case "duration":
			query["sort"] = "+duration";
			break;
		case "duration-desc":
			query["sort"] = "-duration";
			break;
		case "year":
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
	}
}
