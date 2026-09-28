<script lang="ts">
	import { ModeWatcher } from "mode-watcher";
	import "../app.css";
	import { Toaster, Tooltip } from "$lib/components/ui";
	import { onMount } from "svelte";

	const { children } = $props();

	onMount(() => {
		document.documentElement.style.removeProperty("background-color");
	});
</script>

<ModeWatcher />
<Toaster />

<Tooltip.Provider>
	{@render children()}
</Tooltip.Provider>
