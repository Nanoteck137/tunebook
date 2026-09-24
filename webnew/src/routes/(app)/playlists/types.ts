import { defineEnumTypes } from "$lib/utils";
import { z } from "zod";

export const { sortTypes, SortTypeEnum, defaultSort } = defineEnumTypes(
	[
		{ label: "Custom", value: "custom" },
		{ label: "Name (A–Z)", value: "name-a-z" },
		{ label: "Name (Z–A)", value: "name-z-a" },
		{ label: "Tracks (Most)", value: "tracks-most" },
		{ label: "Tracks (Least)", value: "tracks-least" },
		{ label: "Added (New–Old)", value: "created-new" },
		{ label: "Added (Old–New)", value: "created-old" },
		{ label: "Updated (New–Old)", value: "updated-new" },
		{ label: "Updated (Old–New)", value: "updated-old" },
	] as const,
	"custom",
);

export type SortType = (typeof sortTypes)[number]["value"];

export const FullFilter = z.object({
	query: z.string(),
	sort: SortTypeEnum.default(defaultSort),
	filters: z.object({}),
	excludes: z.object({}),
});
export type FullFilter = z.infer<typeof FullFilter>;

export function constructFilterSort(
	filter: FullFilter,
	query: Record<string, string>,
	currentUserId?: string,
) {
	const filters: string[] = [];

	if (filter.query !== "") {
		filters.push(`name contains "${filter.query}"`);
	}

	if (currentUserId) {
		filters.push(`ownerId = "${currentUserId}"`);
	}

	query["filter"] = filters.join(" and ");

	switch (filter.sort) {
		case "custom":
			query["sort"] = "+position";
			break;
		case "name-a-z":
			query["sort"] = "+name";
			break;
		case "name-z-a":
			query["sort"] = "-name";
			break;
		case "tracks-most":
			query["sort"] = "-trackCount";
			break;
		case "tracks-least":
			query["sort"] = "+trackCount";
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
