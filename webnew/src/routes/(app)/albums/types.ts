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
