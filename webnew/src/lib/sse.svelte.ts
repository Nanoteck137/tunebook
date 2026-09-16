import { handleApiError } from "$lib";
import type { ApiClient } from "$lib/api/client";
import { getContext, setContext } from "svelte";

type SseHandler = (data: unknown) => void;

class SseConnection {
	apiClient: ApiClient;

	private eventSource: EventSource | null = null;
	private handlers = new Map<string, Set<SseHandler>>();

	private stopped = false;
	private reconnectTimer: ReturnType<typeof setTimeout> | null = null;

	constructor(apiClient: ApiClient) {
		this.apiClient = apiClient;
	}

	on(eventType: string, handler: SseHandler) {
		let set = this.handlers.get(eventType);
		if (!set) {
			set = new Set();
			this.handlers.set(eventType, set);
		}

		set.add(handler);
	}

	async start() {
		this.stopped = false;
		await this.connect();
	}

	stop() {
		this.stopped = true;

		if (this.reconnectTimer) {
			clearTimeout(this.reconnectTimer);
			this.reconnectTimer = null;
		}

		this.eventSource?.close();
		this.eventSource = null;
	}

	private async connect() {
		if (this.stopped) return;

		const res = await this.apiClient.createSseToken();
		if (!res.success) {
			handleApiError(res.error);
			this.scheduleReconnect();
			return;
		}

		const url = this.apiClient.url.sseHandler();
		url.searchParams.set("token", res.data.token);

		const eventSource = new EventSource(url);
		this.eventSource = eventSource;

		for (const eventType of this.handlers.keys()) {
			eventSource.addEventListener(eventType, (e) => {
				const handlers = this.handlers.get(eventType);
				if (!handlers) return;

				let data: unknown = null;
				try {
					data = JSON.parse((e as MessageEvent).data);
				} catch {
					data = (e as MessageEvent).data;
				}

				for (const handler of handlers) {
					handler(data);
				}
			});
		}

		eventSource.onerror = () => {
			if (eventSource.readyState === EventSource.CLOSED) {
				eventSource.close();
				this.eventSource = null;
				this.scheduleReconnect();
			}
		};
	}

	private scheduleReconnect() {
		if (this.stopped) return;

		this.reconnectTimer = setTimeout(() => {
			this.reconnectTimer = null;
			this.connect();
		}, 2000);
	}
}

const SSE_KEY = Symbol("SSE");

export function setSseConnection(apiClient: ApiClient) {
	return setContext(SSE_KEY, new SseConnection(apiClient));
}

export function getSseConnection() {
	return getContext<ReturnType<typeof setSseConnection>>(SSE_KEY);
}