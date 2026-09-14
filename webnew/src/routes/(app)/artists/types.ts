import { defineEnumTypes } from "$lib/utils";
import { z } from "zod";

export const { sortTypes, SortTypeEnum, defaultSort } = defineEnumTypes(
  [
    { label: "Name (A-Z)", value: "name-a-z" },
    { label: "Name (Z-A)", value: "name-z-a" },
    { label: "Created (New–Old)", value: "created-new" },
    { label: "Created (Old-New)", value: "created-old" },
    { label: "Updated (New–Old)", value: "updated-new" },
    { label: "Updated (Old-New)", value: "updated-old" },
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
