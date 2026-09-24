import { defineEnumTypes } from "$lib/utils";
import { z } from "zod";

export const { sortTypes, SortTypeEnum, defaultSort } = defineEnumTypes(
	[
		{ label: "Name (A–Z)", value: "name-a-z" },
		{ label: "Name (Z–A)", value: "name-z-a" },
		{ label: "Added (New–Old)", value: "created-new" },
		{ label: "Added (Old–New)", value: "created-old" },
		{ label: "Updated (New–Old)", value: "updated-new" },
		{ label: "Updated (Old–New)", value: "updated-old" },
	] as const,
	"name-a-z",
);

export type SortType = (typeof sortTypes)[number]["value"];

export const FullFilter = z.object({
	query: z.string(),
	sort: SortTypeEnum.default(defaultSort),
	filters: z.object({
		tags: z.array(z.string()),
	}),
	excludes: z.object({
		tags: z.array(z.string()),
	}),
});
export type FullFilter = z.infer<typeof FullFilter>;

export function constructFilterSort(
	filter: FullFilter,
	query: Record<string, string>,
) {
	const filters: string[] = [];

	if (filter.query !== "") {
		filters.push(`name contains "${filter.query}"`);
	}

	filter.filters.tags.forEach((t) => {
		filters.push(`tags has "${t}"`);
	});

	filter.excludes.tags.forEach((t) => {
		filters.push(`not tags has "${t}"`);
	});

	query["filter"] = filters.join(" and ");

	switch (filter.sort) {
		case "name-a-z":
			query["sort"] = "+name";
			break;
		case "name-z-a":
			query["sort"] = "-name";
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
