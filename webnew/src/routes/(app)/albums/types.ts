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

export const {
	sortTypes: decadeTypes,
	SortTypeEnum: DecadeTypeEnum,
	defaultSort: defaultDecade,
} = defineEnumTypes(
	[
		{ label: "None", value: "none" },
		{ label: "60s", value: "1960" },
		{ label: "70s", value: "1970" },
		{ label: "80s", value: "1980" },
		{ label: "90s", value: "1990" },
		{ label: "2000s", value: "2000" },
		{ label: "2010s", value: "2010" },
		{ label: "2020s", value: "2020" },
	] as const,
	"none",
);

export type DecadeType = (typeof decadeTypes)[number]["value"];

export const FullFilter = z.object({
	query: z.string(),
	sort: SortTypeEnum.default(defaultSort),
	filters: z.object({
		decade: DecadeTypeEnum.default(defaultDecade),
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

	if (filter.filters.decade !== "none") {
		switch (filter.filters.decade) {
			case "1960":
				filters.push("(year >= 1960 and year <= 1969)");
				break;
			case "1970":
				filters.push("(year >= 1970 and year <= 1979)");
				break;
			case "1980":
				filters.push("(year >= 1980 and year <= 1989)");
				break;
			case "1990":
				filters.push("(year >= 1990 and year <= 1999)");
				break;
			case "2000":
				filters.push("(year >= 2000 and year <= 2009)");
				break;
			case "2010":
				filters.push("(year >= 2010 and year <= 2019)");
				break;
			case "2020":
				filters.push("(year >= 2020 and year <= 2029)");
				break;
		}
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
