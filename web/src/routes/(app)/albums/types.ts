import { z } from "zod";

export const sortTypes = [
  { label: "Name (A-Z)", value: "name-a-z" },
  { label: "Name (Z-A)", value: "name-z-a" },
  { label: "Created (New–Old)", value: "created-new" },
  { label: "Created (Old-New)", value: "created-old" },
  { label: "Updated (New–Old)", value: "updated-new" },
  { label: "Updated (Old-New)", value: "updated-old" },
] as const;
export type SortType = (typeof sortTypes)[number]["value"];
export const SortTypeEnum = z.enum(
  sortTypes.map((f) => f.value) as [SortType, ...SortType[]],
);

export const defaultSort: SortType = "name-a-z";

export const decadeTypes = [
  { label: "None", value: "none" },
  { label: "60s", value: "1960" },
  { label: "70s", value: "1970" },
  { label: "80s", value: "1980" },
  { label: "90s", value: "1990" },
  { label: "2000s", value: "2000" },
  { label: "2010s", value: "2010" },
  { label: "2020s", value: "2020" },
] as const;
export type DecadeType = (typeof decadeTypes)[number]["value"];
export const DecadeTypeEnum = z.enum(
  decadeTypes.map((f) => f.value) as [DecadeType, ...DecadeType[]],
);
export const defaultDecade: DecadeType = "none";

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
