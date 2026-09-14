<script lang="ts">
	import { page } from "$app/state";
	import { Button, InputGroup, Separator } from "$lib/components/ui";
	import { SearchIcon, XIcon } from "@lucide/svelte";

	type Props = {
		searchBarPlaceholder: string;
		value: string;
		setValue: (value: string) => void;
		search: (value: string) => void;
		searchWithValue: () => void;
		clearSearch: () => void;
	};

	const {
		searchBarPlaceholder,
		value,
		setValue,
		search,
		searchWithValue,
		clearSearch,
	}: Props = $props();

	const tabs = [
		{ label: "All", href: "/search" },
		{ label: "Tracks", href: "/search/tracks" },
		{ label: "Artists", href: "/search/artists" },
		{ label: "Albums", href: "/search/albums" },
		{ label: "Playlists", href: "/search/playlists" },
		{ label: "Users", href: "/search/users" },
	];

	let timer: ReturnType<typeof setTimeout>;
	function onInput(e: Event) {
		const target = e.target as HTMLInputElement;
		const current = target.value;
		setValue(current);

		clearTimeout(timer);
		timer = setTimeout(() => {
			search(current);
		}, 500);
	}
</script>

<form
	action=""
	method="get"
	onsubmit={(e) => {
		e.preventDefault();
		clearTimeout(timer);
		searchWithValue();
	}}
>
	<InputGroup.Root>
		<InputGroup.Input
			id="query"
			name="query"
			placeholder={searchBarPlaceholder}
			autocomplete="off"
			{value}
			oninput={onInput}
			autofocus
		/>
		<InputGroup.Addon>
			<SearchIcon />
		</InputGroup.Addon>
		<InputGroup.Addon align="inline-end">
			<InputGroup.Button type="submit">
				<SearchIcon />
			</InputGroup.Button>

			<Separator class="min-h-4" orientation="vertical" />

			<InputGroup.Button
				onclick={() => {
					const e = document.getElementById("query") as HTMLInputElement;
					e.value = "";
					e.focus();
					clearSearch();
				}}
			>
				<XIcon />
			</InputGroup.Button>
		</InputGroup.Addon>
	</InputGroup.Root>
</form>

<nav class="flex flex-wrap gap-1">
	{#each tabs as { label, href }}
		<Button
			class="transition-colors"
			variant={page.url.pathname === href ? "default" : "outline"}
			href="{href}?query={value}"
		>
			{label}
		</Button>
	{/each}
</nav>
