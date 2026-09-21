import { getApiClient, handleApiError } from "$lib";
import { z } from "zod";

export const MissingItem = z.object({
	id: z.string(),
	name: z.string(),
});
export type MissingItemTy = z.infer<typeof MissingItem>;

export const LibrarySyncStateEvent = z.object({
	errors: z.array(z.string()),

	numArtists: z.number(),
	numAlbums: z.number(),
	numTracks: z.number(),

	missingArtists: z.array(MissingItem),
	missingAlbums: z.array(MissingItem),
	missingTracks: z.array(MissingItem),

	artistsSyncDurationMs: z.number(),
	albumsSyncDurationMs: z.number(),
	tracksSyncDurationMs: z.number(),
	totalSyncDurationMs: z.number(),
});
export type LibrarySyncStateEventTy = z.infer<typeof LibrarySyncStateEvent>;

export const TaskSyncStateEventTask = z.object({
	name: z.string(),
	displayName: z.string(),
	isRunning: z.boolean(),
});
export type TaskSyncStateEventTaskTy = z.infer<typeof TaskSyncStateEventTask>;

export const TaskSyncStateEvent = z.object({
	tasks: z.array(TaskSyncStateEventTask),
});
export type TaskSyncStateEventTy = z.infer<typeof TaskSyncStateEvent>;

export const Job = z.object({
	id: z.string(),
	name: z.string(),
	displayName: z.string(),
	status: z.string(),
	error: z.string(),
	attempts: z.number(),
	maxAttempts: z.number(),
	created: z.number(),
	updated: z.number(),
});
export type JobTy = z.infer<typeof Job>;

export const JobSyncStateEvent = z.object({
	jobs: z.array(Job),
});
export type JobSyncStateEventTy = z.infer<typeof JobSyncStateEvent>;

export function jobStatusVariant(status: string) {
	switch (status) {
		case "completed":
			return "default";
		case "running":
			return "secondary";
		case "failed":
			return "destructive";
		case "pending":
			return "outline";
		default:
			return "outline";
	}
}

export function formatJobDate(ms: number) {
	return new Date(ms).toLocaleString();
}

// connectServerSse opens an EventSource for the server pages and dispatches
// each received event to its registered handler.
export async function connectServerSse(
	events: Record<string, (data: unknown) => void>,
): Promise<EventSource | null> {
	const apiClient = getApiClient();

	const res = await apiClient.createSseToken();
	if (!res.success) {
		handleApiError(res.error);
		return null;
	}

	const url = apiClient.url.sseHandler();
	url.searchParams.set("token", res.data.token);

	const eventSource = new EventSource(url);

	for (const [type, handler] of Object.entries(events)) {
		eventSource.addEventListener(type, (e) => {
			handler(JSON.parse((e as MessageEvent).data));
		});
	}

	return eventSource;
}
