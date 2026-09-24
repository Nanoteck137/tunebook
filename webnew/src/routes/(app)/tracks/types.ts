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
