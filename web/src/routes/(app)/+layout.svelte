<script lang="ts">
  import {
    DiscAlbum,
    FileMusic,
    Home,
    ListMusic,
    ListVideo,
    LogIn,
    LogOut,
    Menu,
    Search,
    Server,
    User,
    Users,
  } from "lucide-svelte";
  import AudioPlayer from "$lib/components/audio/AudioPlayer.svelte";
  import Link from "$lib/components/Link.svelte";
  import { browser } from "$app/environment";
  import { fade, fly } from "svelte/transition";
  import { Button, buttonVariants, DropdownMenu } from "@nanoteck137/nano-ui";
  import toast, { Toaster } from "svelte-5-french-toast";
  import { getApiAddress, handleApiError, setApiClient } from "$lib";
  import { setMusicManager } from "$lib/music-manager.svelte";
  import { goto, invalidateAll } from "$app/navigation";
  import { setQuickPlaylist } from "$lib/quick-playlist.svelte";
  import { setFavorites } from "$lib/favorites.svelte";
  import {
    initPlaylistModalManager,
    showPlaylistModal,
  } from "$lib/playlist-modal.svelte";
  import PlaylistSelectorModal from "$lib/components/new-modals/PlaylistSelectorModal.svelte";
  import { isRoleAdmin } from "$lib/utils.js";
  import { page } from "$app/state";

  let { children, data } = $props();

  let apiClient = setApiClient(
    getApiAddress(page.url),
    localStorage.getItem("token") ?? undefined,
  );

  setMusicManager(apiClient);

  setFavorites(apiClient);

  initPlaylistModalManager(apiClient);

  let quickPlaylist = setQuickPlaylist(apiClient);

  $effect(() => {
    quickPlaylist.setPlaylistId(data.user?.quickPlaylist ?? null);
  });

  let showSideMenu = $state(false);

  function close() {
    showSideMenu = false;
  }

  $effect(() => {
    if (showSideMenu) {
      if (browser) document.body.style.overflow = "hidden";
    } else {
      if (browser) document.body.style.overflow = "";
    }
  });
</script>

<svelte:head>
  <title>Tunebook</title>
</svelte:head>

<Toaster position="bottom-right" />

<PlaylistSelectorModal />

<header
  class="sticky top-0 z-50 w-full border-b border-border/40 bg-background/95 backdrop-blur supports-[backdrop-filter]:bg-background/60"
>
  <div
    class="container flex h-14 max-w-screen-2xl items-center gap-4 px-4 sm:px-8"
  >
    <button
      onclick={() => {
        showSideMenu = true;
      }}
    >
      <Menu size="20" />
    </button>

    <a
      class="bg-gradient-to-tr from-logo-1 via-logo-2 to-logo-3 bg-clip-text text-2xl font-medium text-transparent"
      href="/">Tunebook</a
    >

    <div class="hidden items-center gap-1 md:flex">
      <a
        href="/albums"
        class="rounded-md px-3 py-1.5 text-sm font-medium transition-colors hover:bg-accent hover:text-accent-foreground {page.url.pathname.startsWith(
          '/albums',
        )
          ? 'bg-accent text-accent-foreground'
          : 'text-muted-foreground'}"
      >
        Albums
      </a>
      <a
        href="/artists"
        class="rounded-md px-3 py-1.5 text-sm font-medium transition-colors hover:bg-accent hover:text-accent-foreground {page.url.pathname.startsWith(
          '/artists',
        )
          ? 'bg-accent text-accent-foreground'
          : 'text-muted-foreground'}"
      >
        Artists
      </a>
      <a
        href="/tracks"
        class="rounded-md px-3 py-1.5 text-sm font-medium transition-colors hover:bg-accent hover:text-accent-foreground {page.url.pathname.startsWith(
          '/tracks',
        )
          ? 'bg-accent text-accent-foreground'
          : 'text-muted-foreground'}"
      >
        Tracks
      </a>
      {#if data.user}
        <a
          href="/playlists"
          class="rounded-md px-3 py-1.5 text-sm font-medium transition-colors hover:bg-accent hover:text-accent-foreground {page.url.pathname.startsWith(
            '/playlists',
          )
            ? 'bg-accent text-accent-foreground'
            : 'text-muted-foreground'}"
        >
          Playlists
        </a>
      {/if}
    </div>

    <div class="flex-grow"></div>

    <div class="flex items-center gap-2">
      {#if data.user}
        <button
          class={buttonVariants({ variant: "ghost", size: "icon" })}
          onclick={async () => {
            const id = await showPlaylistModal({
              selectedId: data.user?.quickPlaylist ?? undefined,
            });
            if (!id) return;

            const res = await apiClient.setQuickPlaylist({
              playlistId: id,
            });
            if (!res.success) {
              handleApiError(res.error);
              invalidateAll();
              return;
            }

            toast.success("Successfully set quick playlist");
            invalidateAll();
          }}
        >
          <ListVideo />
        </button>
      {/if}

      <Button href="/search" size="icon" variant="ghost">
        <Search />
      </Button>

      {#if data.user}
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
                  goto("/", { invalidateAll: true });
                }}
              >
                <LogOut />
                Logout
              </DropdownMenu.Item>
            </DropdownMenu.Group>
          </DropdownMenu.Content>
        </DropdownMenu.Root>
      {/if}
    </div>
  </div>
</header>

<main class="container px-4 py-4 sm:px-8">
  {@render children()}
</main>

<footer class="fixed bottom-0 z-40 w-full">
  <AudioPlayer />

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
    {#if data.user}
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
    {/if}
  </nav>
</footer>

{#if showSideMenu}
  <!-- svelte-ignore a11y_consider_explicit_label -->
  <button
    class="fixed inset-0 z-50 bg-black/80"
    onclick={() => {
      showSideMenu = false;
    }}
    transition:fade={{ duration: 200 }}
  ></button>

  <aside
    class={`fixed bottom-0 top-0 z-50 flex w-72 flex-col bg-sidebar text-sidebar-foreground`}
    transition:fly={{ x: -400 }}
  >
    <div class="flex h-14 items-center gap-4 border-b px-8">
      <button
        onclick={() => {
          showSideMenu = false;
        }}
      >
        <Menu size="20" />
      </button>
      <a
        class="text-2xl font-medium"
        href="/"
        onclick={() => {
          showSideMenu = false;
        }}
      >
        Tunebook
      </a>
    </div>

    <div class="flex flex-col gap-2 px-4 py-4">
      <Link title="Home" href="/" icon={Home} onClick={close} />
      <Link title="Artists" href="/artists" icon={Users} onClick={close} />
      <Link title="Albums" href="/albums" icon={DiscAlbum} onClick={close} />
      <Link title="Tracks" href="/tracks" icon={FileMusic} onClick={close} />

      {#if data.user}
        <Link
          title="Playlists"
          href="/playlists"
          icon={ListMusic}
          onClick={close}
        />
      {/if}
    </div>
    <div class="flex-grow"></div>
    <div class="flex flex-col gap-2 px-4 py-2">
      {#if data.user}
        <!-- TODO(patrik): Temp -->
        <img class="w-16" src={data.user.picture.small} alt="" />

        <Link
          title={data.user.displayName}
          href="/users/{data.user.id}"
          icon={User}
          onClick={close}
        />

        {#if data.user.role === "super_user"}
          <Link title="Server" href="/server" icon={Server} onClick={close} />
        {/if}

        <Link
          title="Logout"
          icon={LogOut}
          onClick={() => {
            localStorage.removeItem("token");
            invalidateAll();
            goto("/");

            close();
          }}
        />
      {:else}
        <Link title="Login" href="/login" icon={LogIn} onClick={close} />
      {/if}
    </div>
    <div class="h-4"></div>
  </aside>
{/if}
