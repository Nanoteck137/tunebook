<script lang="ts">
	import {
		DiscAlbum,
		FileMusic,
		Heart,
		ListMusic,
		ListVideo,
		LogOut,
		MonitorIcon,
		MoonIcon,
		Search,
		Server,
		SunIcon,
		User,
		Users,
	} from "@lucide/svelte";
	import AudioPlayer from "$lib/components/audio/AudioPlayer.svelte";
	import MobilePlayer from "$lib/components/audio/MobilePlayer.svelte";
	import { getMusicManager } from "$lib/music-manager.svelte";
	import { onMount } from "svelte";
	import { getApiAddress, handleApiError, setApiClient } from "$lib";
	import { setMusicManager } from "$lib/music-manager.svelte";
	import { goto, invalidateAll } from "$app/navigation";
	import { setQuickPlaylist } from "$lib/quick-playlist.svelte";
	import { setFavorites } from "$lib/favorites.svelte";
	import {
		initPlaylistModalManager,
		showPlaylistModal,
	} from "$lib/playlist-modal.svelte";
	import { isRoleAdmin } from "$lib/utils";
	import { page } from "$app/state";
	import { Button, buttonVariants, DropdownMenu } from "$lib/components/ui";
	import { toast } from "svelte-sonner";
	import PlaylistSelectorModal from "$lib/components/new-modals/PlaylistSelectorModal.svelte";
	import { resetMode, setMode } from "mode-watcher";

	let { children, data } = $props();

	let apiClient = setApiClient(
		getApiAddress(page.url),
		localStorage.getItem("token") ?? undefined,
	);

	const musicManager = setMusicManager(apiClient);

	onMount(() => {
		if (data.user) {
			musicManager.initQueue();
		}
	});

	setFavorites(apiClient);

	initPlaylistModalManager(apiClient);

	let quickPlaylist = setQuickPlaylist(apiClient);

	$effect(() => {
		quickPlaylist.setPlaylistId(data.user?.quickPlaylist ?? null);
	});
</script>

<svelte:head>
	<title>Tunebook</title>
</svelte:head>

<PlaylistSelectorModal />

{#if data.user}
	<header
		class="sticky top-0 z-50 w-full border-b border-border/40 bg-background/95 backdrop-blur supports-backdrop-filter:bg-background/60" >
		<div
			class="old-container flex h-14 max-w-screen-2xl items-center gap-4 px-4 sm:px-8"
		>
			<a
				class="bg-linear-to-tr from-logo-1 via-logo-2 to-logo-3 bg-clip-text text-2xl font-medium text-transparent"
				href="/"
			>
				Tunebook
			</a>

			<div class="hidden items-center gap-1 md:flex">
				<Button
					href="/albums"
					class="text-muted-foreground hover:bg-accent hover:text-accent-foreground dark:hover:bg-accent {page.url.pathname.startsWith(
						'/albums',
					)
						? 'bg-accent text-accent-foreground'
						: 'text-muted-foreground'}"
					variant="ghost"
				>
					Albums
				</Button>

				<Button
					href="/artists"
					class="text-muted-foreground hover:bg-accent hover:text-accent-foreground dark:hover:bg-accent {page.url.pathname.startsWith(
						'/artists',
					)
						? 'bg-accent text-accent-foreground'
						: 'text-muted-foreground'}"
					variant="ghost"
				>
					Artists
				</Button>

				<Button
					href="/tracks"
					class="text-muted-foreground hover:bg-accent hover:text-accent-foreground dark:hover:bg-accent {page.url.pathname.startsWith(
						'/tracks',
					)
						? 'bg-accent text-accent-foreground'
						: 'text-muted-foreground'}"
					variant="ghost"
				>
					Tracks
				</Button>

				<Button
					href="/playlists"
					class="text-muted-foreground hover:bg-accent hover:text-accent-foreground dark:hover:bg-accent {page.url.pathname.startsWith(
						'/playlists',
					)
						? 'bg-accent text-accent-foreground'
						: 'text-muted-foreground'}"
					variant="ghost"
				>
					Playlists
				</Button>

				<Button
					href="/favorites"
					class="text-muted-foreground hover:bg-accent hover:text-accent-foreground dark:hover:bg-accent {page.url.pathname.startsWith(
						'/favorites',
					)
						? 'bg-accent text-accent-foreground'
						: 'text-muted-foreground'}"
					variant="ghost"
				>
					Favorites
				</Button>
			</div>

			<div class="grow"></div>

			<div class="flex items-center gap-2">
				<button
					class={buttonVariants({ variant: "ghost", size: "icon" })}
					onclick={async () => {
						const playlist = await showPlaylistModal({
							selectedId: data.user?.quickPlaylist ?? undefined,
							title: "Set Quick Playlist",
							description: "Choose which playlist to use for quick adds",
						});
						if (!playlist) return;

						const res = await apiClient.setQuickPlaylist({
							playlistId: playlist.id,
						});
						if (!res.success) {
							handleApiError(res.error);
							invalidateAll();
							return;
						}

						toast.success("Settings Quick Playlist: " + playlist.name);
						invalidateAll();
					}}
				>
					<ListVideo />
				</button>

				<Button
					href="/search"
					size="icon"
					variant="ghost"
					class="transition-colors hover:bg-accent hover:text-accent-foreground dark:hover:bg-accent  {page.url.pathname.startsWith(
						'/search',
					)
						? 'bg-accent text-accent-foreground'
						: 'text-foreground'}"
				>
					<Search />
				</Button>

				<DropdownMenu.Root>
					<DropdownMenu.Trigger
						class={buttonVariants({ variant: "ghost", size: "icon" })}
					>
						<SunIcon
							class="scale-100 rotate-0 transition-all! dark:scale-0 dark:-rotate-90"
						/>
						<MoonIcon
							class="absolute scale-0 rotate-90 transition-all! dark:scale-100 dark:rotate-0"
						/>
					</DropdownMenu.Trigger>
					<DropdownMenu.Content class="w-56" align="end">
						<DropdownMenu.Item onclick={() => setMode("light")}>
							<SunIcon />
							Light
						</DropdownMenu.Item>
						<DropdownMenu.Item onclick={() => setMode("dark")}>
							<MoonIcon />
							Dark
						</DropdownMenu.Item>
						<DropdownMenu.Item onclick={() => resetMode()}>
							<MonitorIcon />
							System
						</DropdownMenu.Item>
					</DropdownMenu.Content>
				</DropdownMenu.Root>

				<DropdownMenu.Root>
					<DropdownMenu.Trigger>
						<img
							class="w-8 min-w-8 rounded-full"
							src={data.user.picture.small}
							alt=""
						/>
					</DropdownMenu.Trigger>
					<DropdownMenu.Content class="w-56" align="end">
						<DropdownMenu.Group>
							<DropdownMenu.GroupHeading>
								{data.user.displayName}
							</DropdownMenu.GroupHeading>

							<DropdownMenu.Separator />

							<DropdownMenu.Item
								onSelect={() => {
									if (!data.user) return;

									goto(`/users/${data.user.id}`);
								}}
							>
								<User />
								Account
							</DropdownMenu.Item>

							{#if isRoleAdmin(data.user.role)}
								<DropdownMenu.Item
									onSelect={() => {
										goto(`/server`);
									}}
								>
									<Server />
									Server
								</DropdownMenu.Item>
							{/if}

							<DropdownMenu.Separator />

							<DropdownMenu.Item
								onSelect={() => {
									localStorage.removeItem("token");
									musicManager.reset();
									goto("/", { invalidateAll: true });
								}}
							>
								<LogOut />
								Logout
							</DropdownMenu.Item>
						</DropdownMenu.Group>
					</DropdownMenu.Content>
				</DropdownMenu.Root>
			</div>
		</div>
	</header>
{/if}

<main class="old-container px-4 py-4 sm:px-8">
	{@render children()}
</main>

{#if getMusicManager().showPlayer}
	<MobilePlayer />
{/if}

<footer class="fixed bottom-0 z-40 w-full">
	<AudioPlayer />

	{#if data.user}
		<nav
			class="flex items-center justify-around border-t bg-background py-1 md:hidden"
		>
			<a
				href="/albums"
				class="flex flex-col items-center gap-0.5 px-3 py-1 text-xs font-medium transition-colors {page.url.pathname.startsWith(
					'/albums',
				)
					? 'text-primary'
					: 'text-muted-foreground'}"
			>
				<DiscAlbum size={18} />
				Albums
			</a>
			<a
				href="/artists"
				class="flex flex-col items-center gap-0.5 px-3 py-1 text-xs font-medium transition-colors {page.url.pathname.startsWith(
					'/artists',
				)
					? 'text-primary'
					: 'text-muted-foreground'}"
			>
				<Users size={18} />
				Artists
			</a>
			<a
				href="/tracks"
				class="flex flex-col items-center gap-0.5 px-3 py-1 text-xs font-medium transition-colors {page.url.pathname.startsWith(
					'/tracks',
				)
					? 'text-primary'
					: 'text-muted-foreground'}"
			>
				<FileMusic size={18} />
				Tracks
			</a>
			<a
				href="/playlists"
				class="flex flex-col items-center gap-0.5 px-3 py-1 text-xs font-medium transition-colors {page.url.pathname.startsWith(
					'/playlists',
				)
					? 'text-primary'
					: 'text-muted-foreground'}"
			>
				<ListMusic size={18} />
				Playlists
			</a>
			<a
				href="/favorites"
				class="flex flex-col items-center gap-0.5 px-3 py-1 text-xs font-medium transition-colors {page.url.pathname.startsWith(
					'/favorites',
				)
					? 'text-primary'
					: 'text-muted-foreground'}"
			>
				<Heart size={18} />
				Favorites
			</a>
		</nav>
	{/if}
</footer>
