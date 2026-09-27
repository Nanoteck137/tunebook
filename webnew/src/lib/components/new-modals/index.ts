import type { Snippet } from "svelte";

export interface Modal<T> {
	class?: string;
	children?: Snippet<[]>;

	// eslint-disable-next-line no-unused-vars
	onResult: (res: T) => void;
}

export { default as PlaylistSelectorModal } from "./PlaylistSelectorModal.svelte";
export { default as QuickCodeModal } from "./QuickCodeModal.svelte";
