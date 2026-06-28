import { defineEnumTypes } from "$lib/utils";

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
