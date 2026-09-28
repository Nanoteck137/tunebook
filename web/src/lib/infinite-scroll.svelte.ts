import { SvelteSet } from "svelte/reactivity";

export type InfiniteScrollInitialState<T> = {
	items: T[];
	hasMore: boolean;
	page: number;
};

export type InfiniteScrollLoadResult<T> = {
	items: T[];
	hasMore: boolean;
} | null;

export type InfiniteScrollOptions<T> = {
	initialLoad: () => InfiniteScrollInitialState<T>;
	load: (page: number) => Promise<InfiniteScrollLoadResult<T>>;
	itemKey?: (item: T) => string;
};

export class InfiniteScrollController<T> {
	items = $state<T[]>([]);
	hasMore = $state(false);
	loading = $state(false);
	error = $state(false);
	page = $state(0);

	#initialLoad: () => InfiniteScrollInitialState<T>;
	#load: (page: number) => Promise<InfiniteScrollLoadResult<T>>;
	#itemKey: (item: T) => string;
	#seen = new SvelteSet<string>();
	#generation = 0;

	constructor(options: InfiniteScrollOptions<T>) {
		const { initialLoad, load, itemKey = (item: T) => String(item) } = options;

		this.#initialLoad = initialLoad;
		this.#load = load;
		this.#itemKey = itemKey;

		this.reset();
		$effect(() => {
			this.reset();
		});
	}

	setInitial(state: InfiniteScrollInitialState<T>) {
		this.#generation++;
		this.items = [...state.items];
		this.hasMore = state.hasMore;
		this.page = state.page;
		this.#seen = new SvelteSet(state.items.map(this.#itemKey));
		this.loading = false;
		this.error = false;
	}

	reset() {
		this.setInitial(this.#initialLoad());
		queueMicrotask(() => {
			if (this.hasMore && this.items.length === 0) {
				this.loadMore();
			}
		});
	}

	async loadMore() {
		if (this.loading || !this.hasMore) return;

		this.loading = true;
		this.error = false;

		const generation = this.#generation;
		const nextPage = this.page + 1;
		const res = await this.#load(nextPage);
		if (generation !== this.#generation) return;

		if (res) {
			const newItems = res.items.filter((item) => {
				const key = this.#itemKey(item);
				if (this.#seen.has(key)) return false;
				this.#seen.add(key);
				return true;
			});
			if (newItems.length > 0) {
				this.items = [...this.items, ...newItems];
			}
			this.hasMore = res.hasMore;
			this.page = nextPage;
		} else {
			this.error = true;
		}

		this.loading = false;
	}
}
