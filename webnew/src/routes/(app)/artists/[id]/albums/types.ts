import { defineEnumTypes } from "$lib/utils";
import { z } from "zod";

export const { sortTypes, SortTypeEnum, defaultSort } = defineEnumTypes(
	[
		{ label: "Name (A–Z)", value: "name-a-z" },
		{ label: "Name (Z–A)", value: "name-z-a" },
		{ label: "Year (New–Old)", value: "year-new" },
		{ label: "Year (Old–New)", value: "year-old" },
		{ label: "Added (New–Old)", value: "created-new" },
		{ label: "Added (Old–New)", value: "created-old" },
		{ label: "Updated (New–Old)", value: "updated-new" },
		{ label: "Updated (Old–New)", value: "updated-old" },
	] as const,
	"name-a-z",
);

export type SortType = (typeof sortTypes)[number]["value"];

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
