import { defineEnumTypes } from "$lib/utils";
import { z } from "zod";

export type TrackKind = "playlist" | "album";

const orderLabels: Record<TrackKind, { asc: string; desc: string }> = {
	playlist: { asc: "Playlist order", desc: "Reverse order" },
	album: { asc: "Track number", desc: "Reverse track number" },
};

function makeSortTypes(kind: TrackKind) {
	const order = orderLabels[kind];

	return [
		{ label: order.asc, value: "order-asc" },
		{ label: order.desc, value: "order-desc" },
		{ label: "Title (A–Z)", value: "title-asc" },
		{ label: "Title (Z–A)", value: "title-desc" },
		{ label: "Added (Newest first)", value: "added-desc" },
		{ label: "Added (Oldest first)", value: "added-asc" },
	] as const;
}

export const trackSortTypes = {
	playlist: makeSortTypes("playlist"),
	album: makeSortTypes("album"),
};

export const { SortTypeEnum, defaultSort } = defineEnumTypes(
	makeSortTypes("playlist"),
	"order-asc",
);

export type SortType = (typeof trackSortTypes.playlist)[number]["value"];

export function trackSortQuery(kind: TrackKind, sort: SortType): string {
	switch (sort) {
		case "order-asc":
			return kind === "album" ? "+number" : "+position";
		case "order-desc":
			return kind === "album" ? "-number" : "-position";
		case "title-asc":
			return "+name";
		case "title-desc":
			return "-name";
		case "added-desc":
			return "-created";
		case "added-asc":
			return "+created";
	}
}

export type Column = {
	label: string;
	asc: SortType;
	desc: SortType;
	className: string;
	title?: string;
};

const columnBase =
	"flex items-center gap-1 rounded-sm transition-colors hover:text-foreground px-1 py-0.5";

export function trackColumns(kind: TrackKind): Column[] {
	return [
		{
			label: "#",
			asc: "order-asc",
			desc: "order-desc",
			title: `Sort by ${orderLabels[kind].asc.toLowerCase()}`,
			className:
				"flex h-6 w-12 shrink-0 items-center justify-center gap-0.5 rounded-sm transition-colors hover:text-foreground",
		},
		{
			label: "Title",
			asc: "title-asc",
			desc: "title-desc",
			className: `-ml-1 min-w-0 flex-1 ${columnBase}`,
		},
		{
			label: "Added",
			asc: "added-desc",
			desc: "added-asc",
			className: `-ml-1 hidden shrink-0 md:flex ${columnBase}`,
		},
	];
}

export const TrackFilter = z.object({
	query: z.string(),
	sort: SortTypeEnum.default(defaultSort),
});
export type TrackFilter = z.infer<typeof TrackFilter>;

export function buildTrackQuery(
	filter: TrackFilter,
	kind: TrackKind,
	query: Record<string, string>,
) {
	const filters: string[] = [];

	if (filter.query !== "") {
		filters.push(`name contains "${filter.query}"`);
	}

	query["filter"] = filters.join(" and ");
	query["sort"] = trackSortQuery(kind, filter.sort);
}
