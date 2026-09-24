import { defineEnumTypes } from "$lib/utils";

export const { sortTypes, SortTypeEnum, defaultSort } = defineEnumTypes(
	[
		{ label: "Name (A–Z)", value: "name-a-z" },
		{ label: "Name (Z–A)", value: "name-z-a" },
		{ label: "Artist (A–Z)", value: "artist" },
		{ label: "Album (A–Z)", value: "album" },
		{ label: "Duration (Short–Long)", value: "duration" },
		{ label: "Year (New–Old)", value: "year" },
		{ label: "Added (New–Old)", value: "created-new" },
		{ label: "Added (Old–New)", value: "created-old" },
	] as const,
	"name-a-z",
);

export type SortType = (typeof sortTypes)[number]["value"];
