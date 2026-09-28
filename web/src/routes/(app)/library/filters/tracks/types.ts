import type { TrackFilter } from "$lib/api/types";

export const sortTypes = [
	{
		label: "Position",
		value: "position",
		reverse: "position-reverse",
		direction: "asc",
	},
	{
		label: "Name",
		value: "name-a-z",
		reverse: "name-z-a",
		direction: "asc",
	},
	{
		label: "Created",
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

export type SortType = (typeof sortValues)[number];

export const defaultSort: SortType = "position";

export function searchFilters(filters: TrackFilter[], query: string) {
	const q = query.trim().toLowerCase();
	if (q === "") {
		return [...filters];
	}

	return filters.filter(
		(filter) =>
			filter.name.toLowerCase().includes(q) ||
			filter.filter.toLowerCase().includes(q),
	);
}

export function sortFilters(filters: TrackFilter[], sort: SortType) {
	const list = [...filters];

	switch (sort) {
		case "position":
			return list;
		case "position-reverse":
			return list.reverse();
		case "name-a-z":
			list.sort((a, b) => a.name.localeCompare(b.name));
			return list;
		case "name-z-a":
			list.sort((a, b) => b.name.localeCompare(a.name));
			return list;
		case "created-new":
			list.sort((a, b) => Date.parse(b.created) - Date.parse(a.created));
			return list;
		case "created-old":
			list.sort((a, b) => Date.parse(a.created) - Date.parse(b.created));
			return list;
		case "updated-new":
			list.sort((a, b) => Date.parse(b.updated) - Date.parse(a.updated));
			return list;
		case "updated-old":
			list.sort((a, b) => Date.parse(a.updated) - Date.parse(b.updated));
			return list;
	}
}
