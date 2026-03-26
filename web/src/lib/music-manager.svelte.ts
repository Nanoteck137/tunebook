/* eslint-disable no-unused-vars */
import type { ApiClient } from "$lib/api/client";
import type { MediaItem } from "$lib/api/types";
import { getContext, setContext } from "svelte";

type SavedQueue = {
  items: string[];
  index: number;
};

class Queue {
  items = $state<MediaItem[]>([]);
  index = $state(0);

  setQueueItems(items: MediaItem[], index?: number) {
    this.items = items;
    this.index = index ?? 0;
  }

  getCurrentMediaItem() {
    if (this.items.length === 0) return null;
    return this.items[this.index];
  }

  isEndOfQueue() {
    return this.index >= this.items.length - 1;
  }

  isQueueEmpty() {
    return this.items.length === 0;
  }

  setQueueIndex(index: number) {
    if (index >= this.items.length) {
      this.index = 0;
      return;
    }

    if (index < 0) {
      return;
    }

    this.index = index;

    this.saveQueue();
  }

  clearQueue() {
    this.items = [];
    this.index = 0;

    this.saveQueue();
  }

  saveQueue() {
    const q: SavedQueue = {
      items: this.items.map((i) => i.trackId),
      index: this.index,
    };

    localStorage.setItem("queue", JSON.stringify(q));
  }

  async loadQueue(loadTracks: (ids: string[]) => Promise<MediaItem[]>) {
    const s = localStorage.getItem("queue");
    if (s) {
      const q: SavedQueue = JSON.parse(s);
      await loadTracks(q.items);
      this.index = q.index;
    }
  }
}

function getVolume(): number {
  const volume = localStorage.getItem("player-volume");
  if (volume) {
    return parseFloat(volume);
  }

  return 1.0;
}

function getMuted(): boolean {
  const muted = localStorage.getItem("player-muted");
  if (muted) {
    return muted === "true";
  }

  return false;
}

export class MusicManager {
  private audio: HTMLAudioElement;

  private apiClient: ApiClient;
  queue: Queue;

  loading = $state(false);
  playing = $state(false);

  currentTime = $state(0);
  duration = $state(0);

  volume = $state(getVolume());
  muted = $state(getMuted());

  showPlayer = $state(false);

  currentItem = $state<MediaItem | null>(null);

  private currentTrackId: string | null = null;
  private trackEventSent: boolean = false;

  constructor(apiClient: ApiClient) {
    this.audio = new Audio();

    this.apiClient = apiClient;
    this.queue = new Queue();

    this.setupAudio();

    this.volume = getVolume();
    this.muted = getMuted();

    if (this.muted) {
      this.audio.volume = 0.0;
    } else {
      this.audio.volume = this.volume;
    }

    /*
    this.queue.loadQueue(async (ids) => {
      const tracks = await apiClient.loadTrackFromIds({ ids });
      if (!tracks.success) {
        // TODO(patrik): Handle error
        return [];
      }

      return tracks.data.tracks.map(
        (track) =>
          ({
            trackId: track.id,
            name: track.name,
          }) as MediaItem,
      );
    });
    */

    this.queueUpdate();
  }

  setupAudio() {
    this.audio.addEventListener("canplay", () => {
      // console.log("canplay");
      this.loading = false;
    });

    this.audio.addEventListener("loadstart", () => {
      // console.log("loadstart");
      this.loading = true;
    });

    this.audio.addEventListener("loadedmetadata", () => {
      // console.log("loadedmetadata");
      this.currentTime = this.audio.currentTime;
      this.duration = this.audio.duration;
    });

    this.audio.addEventListener("progress", () => {
      // console.log("progress");
    });

    this.audio.addEventListener("timeupdate", () => {
      this.currentTime = this.audio.currentTime;
    });

    this.audio.addEventListener("loadeddata", () => {
      // console.log("loadeddata");
    });

    this.audio.addEventListener("playing", () => {
      // console.log("playing");
      this.playing = true;
    });

    this.audio.addEventListener("pause", () => {
      // console.log("pause");
      this.playing = false;
    });

    this.audio.addEventListener("load", () => {
      // console.log("load");
    });

    this.audio.addEventListener("ended", async () => {
      await this.nextTrack();
    });

    // musicManager.emitter.on("requestPlay", () => {
    //   audio.play();
    // });

    // musicManager.emitter.on("requestPause", () => {
    //   audio.pause();
    // });

    // musicManager.emitter.on("requestPlayPause", () => {
    //   if (playing) {
    //     audio.pause();
    //   } else {
    //     audio.play();
    //   }
    // });
  }

  private resetTrackEventTracking(trackId: string) {
    this.currentTrackId = trackId;
    this.trackEventSent = false;
    console.log("resetTrackEventTracking", trackId);
  }

  private checkTrackProgress() {
    if (!this.currentTrackId || this.trackEventSent) return;

    const progressPercent = (this.currentTime / this.duration) * 100;
    if (progressPercent >= 80) {
      this.sendTrackEvent();
    }
  }

  private async sendTrackEvent() {
    if (this.trackEventSent) return;
    if (!this.currentTrackId) return;

    const progressPercent = (this.currentTime / this.duration) * 100;
    if (progressPercent < 10) {
      return;
    }

    const res = await this.apiClient.addTrackEvent(this.currentTrackId, {
      position: this.currentTime,
      source: "web-player",
    });
    if (!res.success) {
      console.log("failed to add track event", res.error);
      return;
    }

    this.trackEventSent = true;
  }

  setPosition(position: number) {
    this.audio.currentTime = position;
  }

  async clearQueue() {
    await this.queue.clearQueue();
    this.queueUpdate();
  }

  /*
  async queueRequest(request: QueueRequest, options: QueueRequestOptions) {
    // onPlay={async (shuffle) => {
    //   await musicManager.clearQueue();
    //   await musicManager.addFromPlaylist(data.playlist.id, { shuffle });
    //   musicManager.requestPlay();
    // }}
    // onTrackPlay={async (trackId) => {
    //   await musicManager.clearQueue();
    //   await musicManager.addFromPlaylist(data.playlist.id);
    //   await musicManager.setQueueIndex(
    //     musicManager.queue.items.findIndex((t) => t.track.id === trackId),
    //   );
    //   musicManager.requestPlay();
    // }}

    if (!options.append) {
      await this.clearQueue();
    }

    switch (request.type) {
      case "addPlaylist": {
        await this.addFromPlaylist(request.playlistId, request.filterId, {
          shuffle: options.shuffle,
          front: options.append === "front",
        });
        break;
      }
      case "addAlbum": {
        await this.addFromAlbum(request.albumId, {
          shuffle: options.shuffle,
          front: options.append === "front",
        });
        break;
      }
      case "addArtist": {
        await this.addFromArtist(request.artistId, {
          shuffle: options.shuffle,
          front: options.append === "front",
        });
        break;
      }
      case "addTaglist": {
        await this.addFromTaglist(request.taglistId, {
          shuffle: options.shuffle,
          front: options.append === "front",
        });
        break;
      }
      case "addFilter": {
        await this.addFromFilter(request.filter, {
          shuffle: options.shuffle,
          front: options.append === "front",
        });
        break;
      }
    }

    if (options.queueIndex) {
      await this.setQueueIndex(options.queueIndex);
    }

    if (options.queueIndexToTrackId) {
      await this.setQueueIndex(
        this.queue.items.findIndex(
          (t) => t.trackId === options.queueIndexToTrackId,
        ),
      );
    }

    if (!options.skipPlay) {
      this.requestPlay();
    }
  }
    */

  /*
  async addFromPlaylist(
    playlistId: string,
    filterId?: string,
    settings?: AddToQueueSettings,
  ) {
    await this.queue.addFromPlaylist(playlistId, filterId, settings);
    this.emitter.emit("onQueueUpdated");
  }

  async addFromTaglist(taglistId: string, settings?: AddToQueueSettings) {
    await this.queue.addFromTaglist(taglistId, settings);
    this.emitter.emit("onQueueUpdated");
  }

  async addFromFilter(filter: string, settings?: AddToQueueSettings) {
    await this.queue.addFromFilter(filter, settings);
    this.emitter.emit("onQueueUpdated");
  }

  async addFromArtist(artistId: string, settings?: AddToQueueSettings) {
    await this.queue.addFromArtist(artistId, settings);
    this.emitter.emit("onQueueUpdated");
  }

  async addFromAlbum(albumId: string, settings?: AddToQueueSettings) {
    await this.queue.addFromAlbum(albumId, settings);
    this.emitter.emit("onQueueUpdated");
  }

  async addFromIds(trackIds: string[], settings?: AddToQueueSettings) {
    await this.queue.addFromIds(trackIds, settings);
    this.emitter.emit("onQueueUpdated");
  }
    */

  async setQueueIndex(index: number) {
    await this.markTrack();

    this.queue.setQueueIndex(index);
    this.queueUpdate();
  }

  queueUpdate() {
    console.log("update");
    this.showPlayer = !this.queue.isQueueEmpty();
    const mediaItem = this.queue.getCurrentMediaItem();

    if (mediaItem) {
      if (this.currentItem?.trackId === mediaItem.trackId) return;

      this.resetTrackEventTracking(mediaItem.trackId);

      const src = this.apiClient.url.streamTrack(mediaItem.trackId).toString();
      console.log(src);
      this.audio.src = src;

      this.currentItem = mediaItem;
    } else {
      this.audio.removeAttribute("src");
      this.audio.load();
    }
  }

  async markTrack() {
    // this.queue.markTrack(this.getPlayerPosition());
  }

  async nextTrack() {
    await this.sendTrackEvent();
    await this.setQueueIndex(this.queue.index + 1);
    this.play();
  }

  async previousTrack() {
    await this.sendTrackEvent();
    await this.setQueueIndex(this.queue.index - 1);
    this.play();
  }

  play() {
    this.audio.play();
  }

  pause() {
    this.audio.pause();
  }

  requestPlayPause() {
    if (this.playing) {
      this.pause();
    } else {
      this.play();
    }
  }

  async addTracks(params: { clear?: boolean; trackId?: string }) {
    const res = await this.apiClient.getMediaFromFilter({
      filter: "",
    });
    if (!res.success) {
      console.log("error getting media", res.error);
      return;
    }

    this.queue.items.push(...res.data.items);
    if (params.trackId) {
      this.queue.index = this.queue.items.findIndex(
        (item) => item.trackId == params.trackId,
      );
    } else {
      this.queue.index = 0;
    }
    this.queueUpdate();

    this.play();
  }
}

const MUSIC_MANAGER_KEY = Symbol("MUSIC_MANAGER");

export function setMusicManager(apiClient: ApiClient) {
  return setContext(MUSIC_MANAGER_KEY, new MusicManager(apiClient));
}

export function getMusicManager() {
  return getContext<ReturnType<typeof setMusicManager>>(MUSIC_MANAGER_KEY);
}
