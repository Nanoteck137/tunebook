import { z } from "zod";

export const sortTypes = [
	{
		label: "Position",
		value: "position",
		reverse: "position-reverse",
		field: "position",
		direction: "asc",
	},
	{
		label: "Name",
		value: "name-a-z",
		reverse: "name-z-a",
		field: "name",
		direction: "asc",
	},
	{
		label: "Tracks",
		value: "tracks-most",
		reverse: "tracks-least",
		field: "trackCount",
		direction: "desc",
	},
	{
		label: "Created",
		value: "created-new",
		reverse: "created-old",
		field: "created",
		direction: "desc",
	},
	{
		label: "Updated",
		value: "updated-new",
		reverse: "updated-old",
		field: "updated",
		direction: "desc",
	},
] as const;

const sortValues = [
	...sortTypes.map((t) => t.value),
	...sortTypes.map((t) => t.reverse),
] as const;

export const SortTypeEnum = z.enum(sortValues);
export type SortType = (typeof sortValues)[number];

export const defaultSort: SortType = "position";

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

	for (const type of sortTypes) {
		const reversed = filter.sort === type.reverse;
		if (filter.sort === type.value || reversed) {
			const dir = reversed
				? type.direction === "asc"
					? "desc"
					: "asc"
				: type.direction;
			query["sort"] = (dir === "asc" ? "+" : "-") + type.field;
			return;
		}
	}

	query["sort"] = "+position";
}
