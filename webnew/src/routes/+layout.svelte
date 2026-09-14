<script lang="ts">
	import { ModeWatcher } from "mode-watcher";
	import "../app.css";
	import { Toaster, Tooltip } from "$lib/components/ui";

	const { children } = $props();
</script>

<ModeWatcher />
<Toaster />

<Tooltip.Provider>
	{@render children()}
</Tooltip.Provider>
