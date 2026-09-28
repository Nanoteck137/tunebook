import { z } from "zod";

export const sortTypes = [
	{ label: "Newest", value: "newest", reverse: "oldest", direction: "desc" },
	{ label: "Name", value: "name-a-z", reverse: "name-z-a", direction: "asc" },
	{
		label: "Played",
		value: "percent-desc",
		reverse: "percent-asc",
		direction: "desc",
	},
] as const;

const sortValues = [
	...sortTypes.map((t) => t.value),
	...sortTypes.map((t) => t.reverse),
] as const;

export const SortTypeEnum = z.enum(sortValues);
export type SortType = (typeof sortValues)[number];

export const defaultSort: SortType = "newest";

// PushTrackHistory only ever records "completed" (>= 80% played) or
// "skipped"; anything under 10% is discarded before insert.
export const statusTypes = [
	{ label: "All", value: "all" },
	{ label: "Completed", value: "completed" },
	{ label: "Skipped", value: "skipped" },
] as const;

export const StatusTypeEnum = z.enum(statusTypes.map((t) => t.value));
export type StatusType = z.infer<typeof StatusTypeEnum>;

export const defaultStatus: StatusType = "all";

export const FullFilter = z.object({
	query: z.string(),
	sort: SortTypeEnum.default(defaultSort),
	status: StatusTypeEnum.default(defaultStatus),
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

	if (filter.status !== defaultStatus) {
		filters.push(`status = "${filter.status}"`);
	}

	query["filter"] = filters.join(" and ");

	switch (filter.sort) {
		case "newest":
			query["sort"] = "-listenedAt";
			break;
		case "oldest":
			query["sort"] = "+listenedAt";
			break;
		case "name-a-z":
			query["sort"] = "+name";
			break;
		case "name-z-a":
			query["sort"] = "-name";
			break;
		case "percent-desc":
			query["sort"] = "-percentPlayed";
			break;
		case "percent-asc":
			query["sort"] = "+percentPlayed";
			break;
	}
}
