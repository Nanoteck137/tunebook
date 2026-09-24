<script lang="ts">
	import { Input } from "$lib/components/ui";
	import { X } from "@lucide/svelte";
	import { cn } from "$lib/utils";

	let {
		value,
		setValue,
		search,
		placeholder = "Search...",
		debounce = 500,
		inputClass = "",
		iconSize = 14,
		autofocus = false,
		class: className = "",
	}: {
		value: string;
		setValue: (value: string) => void;
		search: (value: string) => void;
		placeholder?: string;
		debounce?: number;
		inputClass?: string;
		iconSize?: number;
		autofocus?: boolean;
		class?: string;
	} = $props();

	let timer: ReturnType<typeof setTimeout>;
	let inputEl = $state<HTMLInputElement | null>(null);

	function onInput(e: Event) {
		const target = e.target as HTMLInputElement;
		const current = target.value;
		setValue(current);

		clearTimeout(timer);
		timer = setTimeout(() => {
			search(current);
		}, debounce);
	}

	function onSubmit(e: SubmitEvent) {
		e.preventDefault();
		clearTimeout(timer);
		search(value);
	}

	function onClear() {
		if (inputEl) {
			inputEl.value = "";
			inputEl.focus();
		}
		setValue("");
		clearTimeout(timer);
		search("");
	}
</script>

<form
	action=""
	method="get"
	onsubmit={onSubmit}
	class={cn("relative", className)}
>
	<Input
		bind:ref={inputEl}
		class={cn("pr-8", inputClass)}
		{placeholder}
		{value}
		oninput={onInput}
		{autofocus}
	/>
	{#if value}
		<button
			type="button"
			class="absolute top-1/2 right-1.5 -translate-y-1/2 rounded-full p-0.5 text-muted-foreground hover:text-foreground"
			onclick={onClear}
			aria-label="Clear search"
		>
			<X size={iconSize} />
		</button>
	{/if}
</form>
