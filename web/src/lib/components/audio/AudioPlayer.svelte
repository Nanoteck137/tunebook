<script lang="ts">
  import LargePlayer from "$lib/components/audio/LargePlayer.svelte";
  import SmallPlayer from "$lib/components/audio/SmallPlayer.svelte";
  import { getMusicManager } from "$lib/music-manager.svelte";

  const musicManager = getMusicManager();

  $effect(() => {
    if (musicManager.showPlayer) {
      document.body.setAttribute("data-player", "true");
    } else {
      document.body.setAttribute("data-player", "false");
    }
  });
</script>

<!-- TODO(patrik): Fix this because input fields need space -->
<!-- <svelte:window
  onkeypress={(e) => {
    if (e.key === " ") {
      e.preventDefault();

      musicManager.requestPlayPause();
    }
  }}
/> -->

<!-- playing={musicManager.playing}
    loading={musicManager.loading}
    mediaItem={currentMediaItem}
    {currentTime}
    {duration}
    {volume}
    {queue}
    {currentQueueIndex}
    audioMuted={muted}
    onPlay={() => {
      audio.play();
    }}
    onPause={() => {
      audio.pause();
    }}
    onNextTrack={() => {
      musicManager.nextTrack();
    }}
    onPrevTrack={() => {
      musicManager.previousTrack();
    }}
    onSeek={(e) => {
      audio.currentTime = e;
    }}
    onVolumeChanged={(e) => {
      if (!muted) {
        audio.volume = e;
      }

      volume = e;
      localStorage.setItem("player-volume", e.toString());
    }}
    onToggleMuted={() => {
      muted = !muted;
      localStorage.setItem("player-muted", muted ? "true" : "false");

      if (muted) {
        audio.volume = 0;
      } else {
        audio.volume = volume;
      }
    }} -->

{#if musicManager.showPlayer}
  <LargePlayer />
  <SmallPlayer />
{/if}
