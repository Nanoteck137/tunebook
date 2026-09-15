import { defineEnumTypes } from "$lib/utils";
import { z } from "zod";

export const { sortTypes, SortTypeEnum, defaultSort } = defineEnumTypes(
  [
    { label: "Custom", value: "custom" },
    { label: "Name (A-Z)", value: "name-a-z" },
    { label: "Name (Z-A)", value: "name-z-a" },
    { label: "Tracks (Most)", value: "tracks-most" },
    { label: "Tracks (Least)", value: "tracks-least" },
    { label: "Created (New–Old)", value: "created-new" },
    { label: "Created (Old-New)", value: "created-old" },
    { label: "Updated (New–Old)", value: "updated-new" },
    { label: "Updated (Old-New)", value: "updated-old" },
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
