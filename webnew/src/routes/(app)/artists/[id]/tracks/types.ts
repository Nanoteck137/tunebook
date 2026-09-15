import { defineEnumTypes } from "$lib/utils";
import { z } from "zod";

export const { sortTypes, SortTypeEnum, defaultSort } = defineEnumTypes(
	[
		{ label: "Name (A-Z)", value: "name-a-z" },
		{ label: "Name (Z-A)", value: "name-z-a" },
		{ label: "Artist", value: "artist" },
		{ label: "Album", value: "album" },
		{ label: "Duration", value: "duration" },
		{ label: "Year", value: "year" },
		{ label: "Added (New–Old)", value: "created-new" },
		{ label: "Added (Old-New)", value: "created-old" },
	] as const,
	"name-a-z",
);

export type SortType = (typeof sortTypes)[number]["value"];

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
			query["sort"] = "+artist";
			break;
		case "album":
			query["sort"] = "+album";
			break;
		case "duration":
			query["sort"] = "+duration";
			break;
		case "year":
			query["sort"] = "-year";
			break;
		case "created-new":
			query["sort"] = "-created";
			break;
		case "created-old":
			query["sort"] = "+created";
			break;
	}
}
