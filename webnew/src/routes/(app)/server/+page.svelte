<script lang="ts">
	import { PUBLIC_COMMIT, PUBLIC_VERSION } from "$env/static/public";
	import {
		Database,
		Disc3,
		DiscAlbum,
		Heart,
		Info,
		LayoutGrid,
		ListMusic,
		ListFilter,
		Music,
		Play,
		Server,
		Users,
	} from "@lucide/svelte";
	import { Card, Separator } from "$lib/components/ui";
	import { formatPlayTime } from "$lib/utils";

	const { data } = $props();

	const formatsCount = $derived(data.mediaSettings.formats.length);
	const deviceSpecsCount = $derived(data.mediaSettings.deviceSpecs.length);

	const statTiles = $derived([
		{ icon: Users, label: "Users", value: data.serverStats.users },
		{ icon: Disc3, label: "Artists", value: data.serverStats.artists },
		{ icon: DiscAlbum, label: "Albums", value: data.serverStats.albums },
		{ icon: Music, label: "Tracks", value: data.serverStats.tracks },
		{ icon: ListMusic, label: "Playlists", value: data.serverStats.playlists },
		{ icon: Heart, label: "Favorites", value: data.serverStats.favorites },
		{
			icon: ListFilter,
			label: "Track filters",
			value: data.serverStats.trackFilters,
		},
		{ icon: LayoutGrid, label: "Queues", value: data.serverStats.queues },
	]);

	function formatUptime(ms: number) {
		const s = Math.floor(ms / 1000);
		const d = Math.floor(s / 86400);
		const h = Math.floor((s % 86400) / 3600);
		const m = Math.floor((s % 3600) / 60);

		if (d > 0) return `${d}d ${h}h`;
		if (h > 0) return `${h}h ${m}m`;
		return `${m}m`;
	}

	function formatBytes(bytes: number) {
		if (bytes <= 0) return "0 B";

		const units = ["B", "kB", "MB", "GB", "TB"];
		const i = Math.min(
			Math.floor(Math.log(bytes) / Math.log(1024)),
			units.length - 1,
		);
		return `${(bytes / 1024 ** i).toFixed(i === 0 ? 0 : 1)} ${units[i]}`;
	}

	const uptime = $derived(
		formatUptime(Date.now() - data.systemInfo.startedAt),
	);
	const startedAt = $derived(new Date(data.systemInfo.startedAt));

	const overview = $derived([
		{ label: "Server version", value: data.systemInfo.version, mono: true },
		{ label: "Commit", value: data.systemInfo.commit, mono: true },
		{ label: "Uptime", value: uptime },
		{ label: "Started", value: startedAt.toLocaleString() },
		{ label: "Client version", value: `v${PUBLIC_VERSION}`, mono: true },
		{ label: "Client commit", value: PUBLIC_COMMIT, mono: true },
		{ label: "Media formats", value: String(formatsCount) },
		{ label: "Device specs", value: String(deviceSpecsCount) },
		{ label: "Data directory", value: data.serverStats.dataDir, mono: true },
		{
			label: "Database file",
			value: data.serverStats.databaseFile,
			mono: true,
		},
		{
			label: "Database size",
			value: formatBytes(data.serverStats.databaseSize),
		},
	]);
</script>

<div class="flex flex-col gap-6">
	<div class="flex items-center gap-3">
		<Server size={24} />
		<div class="flex flex-col">
			<h1 class="text-xl font-bold">Server</h1>
			<p class="text-xs text-muted-foreground">
				v{data.systemInfo.version}
				{#if data.systemInfo.commit}({data.systemInfo.commit}){/if}
			</p>
		</div>
	</div>

	<Card.Root>
		<Card.Content>
			<div class="flex items-center gap-2">
				<Info size={18} />
				<h2 class="text-lg font-semibold">Overview</h2>
			</div>

			<Separator class="my-4" />

			<div class="grid grid-cols-1 gap-x-8 gap-y-4 sm:grid-cols-2">
				{#each overview as item (item.label)}
					<div class="flex items-center justify-between gap-4">
						<span class="text-sm text-muted-foreground">{item.label}</span>
						<span class="text-sm font-medium {item.mono ? 'font-mono' : ''}">
							{item.value}
						</span>
					</div>
				{/each}
			</div>
		</Card.Content>
	</Card.Root>

	<Card.Root>
		<Card.Content>
			<div class="flex items-center gap-2">
				<Database size={18} />
				<h2 class="text-lg font-semibold">Stats</h2>
			</div>

			<Separator class="my-4" />

			<div class="grid grid-cols-2 gap-4 sm:grid-cols-3 md:grid-cols-4">
				{#each statTiles as tile (tile.label)}
					<div class="flex flex-col gap-2 rounded border p-4">
						<div class="flex items-center gap-2 text-sm text-muted-foreground">
							<tile.icon size={16} />
							<span>{tile.label}</span>
						</div>
						<span class="text-2xl font-bold">{tile.value}</span>
					</div>
				{/each}
			</div>
		</Card.Content>
	</Card.Root>

	<Card.Root>
		<Card.Content>
			<div class="flex items-center gap-2">
				<Play size={18} />
				<h2 class="text-lg font-semibold">Usage</h2>
			</div>

			<Separator class="my-4" />

			<div class="grid grid-cols-1 gap-x-8 gap-y-4 sm:grid-cols-2">
				<div class="flex items-center justify-between gap-4">
					<span class="text-sm text-muted-foreground">Total plays</span>
					<span class="text-sm font-medium">
						{data.serverStats.totalPlays}
					</span>
				</div>
				<div class="flex items-center justify-between gap-4">
					<span class="text-sm text-muted-foreground"
						>Total listening time</span
					>
					<span class="text-sm font-medium">
						{formatPlayTime(data.serverStats.totalListeningTime)}
					</span>
				</div>
			</div>
		</Card.Content>
	</Card.Root>

	<Card.Root>
		<Card.Content>
			<div class="flex items-center gap-2">
				<Disc3 size={18} />
				<h2 class="text-lg font-semibold">Media Configuration</h2>
			</div>

			<Separator class="my-4" />

			<div class="flex flex-col gap-6">
				<div>
					<h3 class="mb-2 text-sm font-medium text-muted-foreground">
						Formats
					</h3>
					<div class="overflow-x-auto">
						<table class="w-full text-left text-sm">
							<thead>
								<tr class="text-muted-foreground">
									<th class="px-3 py-2 font-medium">Name</th>
									<th class="px-3 py-2 font-medium">Format</th>
									<th class="px-3 py-2 font-medium">Ext</th>
									<th class="px-3 py-2 font-medium">High</th>
									<th class="px-3 py-2 font-medium">Medium</th>
									<th class="px-3 py-2 font-medium">Low</th>
								</tr>
							</thead>
							<tbody>
								{#each data.mediaSettings.formats as format (format.name)}
									<tr class="border-t">
										<td class="px-3 py-2">{format.name}</td>
										<td class="px-3 py-2 font-mono text-xs">{format.format}</td
										>
										<td class="px-3 py-2 font-mono text-xs">{format.ext}</td>
										<td class="px-3 py-2">{format.qualityHighBitrate}kbps</td>
										<td class="px-3 py-2">{format.qualityMediumBitrate}kbps</td
										>
										<td class="px-3 py-2">{format.qualityLowBitrate}kbps</td>
									</tr>
								{/each}
							</tbody>
						</table>
					</div>
				</div>

				<Separator />

				<div>
					<h3 class="mb-2 text-sm font-medium text-muted-foreground">
						Device Specs
					</h3>
					<div class="overflow-x-auto">
						<table class="w-full text-left text-sm">
							<thead>
								<tr class="text-muted-foreground">
									<th class="px-3 py-2 font-medium">Name</th>
									<th class="px-3 py-2 font-medium">Preferred Format</th>
									<th class="px-3 py-2 font-medium">Allowed Formats</th>
								</tr>
							</thead>
							<tbody>
								{#each data.mediaSettings.deviceSpecs as spec (spec.name)}
									<tr class="border-t">
										<td class="px-3 py-2">{spec.name}</td>
										<td class="px-3 py-2 font-mono text-xs"
											>{spec.preferedFormat}</td
										>
										<td class="px-3 py-2 font-mono text-xs"
											>{spec.allowedFormats.join(", ")}</td
										>
									</tr>
								{/each}
							</tbody>
						</table>
					</div>
				</div>
			</div>
		</Card.Content>
	</Card.Root>
</div>
