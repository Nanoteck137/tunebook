<script lang="ts">
	import { page } from "$app/state";
	import { Button, InputGroup, Separator } from "$lib/components/ui";
	import {
		Disc,
		ListMusic,
		Music,
		Mic,
		SearchIcon,
		Users,
		XIcon,
	} from "@lucide/svelte";

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
		{ label: "All", href: "/search", icon: SearchIcon },
		{ label: "Tracks", href: "/search/tracks", icon: Music },
		{ label: "Artists", href: "/search/artists", icon: Mic },
		{ label: "Albums", href: "/search/albums", icon: Disc },
		{ label: "Playlists", href: "/search/playlists", icon: ListMusic },
		{ label: "Users", href: "/search/users", icon: Users },
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
	{#each tabs as { label, href, icon: Icon }}
		<Button
			class="transition-colors hover:bg-accent hover:text-accent-foreground dark:hover:bg-accent {page
				.url.pathname === href
				? 'bg-accent text-accent-foreground'
				: 'text-muted-foreground'}"
			variant="ghost"
			href="{href}?query={value}"
		>
			<Icon />
			{label}
		</Button>
	{/each}
</nav>
