/* eslint-disable no-unused-vars */
import type { ApiClient } from "$lib/api/client";
import { getContext, setContext } from "svelte";

export type MediaRef = {
  id: string;
  name: string;
};

export type MediaItem = {
  trackId: string;
  name: string;
  album: MediaRef;
  artists: MediaRef[];
  coverArt: string;
};

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

  private trackEventSent = $state(false);

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

  private resetTrackEventTracking() {
    console.log("resetTrackEventTracking");
    this.trackEventSent = false;

    console.log("after resetTrackEventTracking", this.trackEventSent);
  }

  private async sendTrackEvent() {
    console.log("sendTrackEvent", this.currentItem?.name, this.trackEventSent);

    if (this.trackEventSent) {
      console.log("sendTrackEvent already sent", this.currentItem?.trackId);
      return;
    }

    if (!this.currentItem) {
      console.log("sendTrackEvent no current item");
      return;
    }

    console.log("SENDING EVENT");
    // const res = await this.apiClient.addTrackEvent(this.currentItem.trackId, {
    //   position: this.currentTime,
    //   source: "web-player",
    // });
    // if (!res.success) {
    //   console.log("failed to add track event", res.error);
    //   return;
    // }

    this.trackEventSent = true;
  }

  setPosition(position: number) {
    this.audio.currentTime = position;
  }

  async clearQueue(update = true) {
    await this.sendTrackEvent();

    this.queue.clearQueue();
    if (update) {
      this.queueUpdate();
    }
  }

  async setQueueIndex(index: number) {
    await this.sendTrackEvent();

    this.queue.setQueueIndex(index);
    this.queueUpdate();
  }

  queueUpdate() {
    this.showPlayer = !this.queue.isQueueEmpty();
    const mediaItem = this.queue.getCurrentMediaItem();

    if (mediaItem) {
      if (mediaItem?.trackId !== this.currentItem?.trackId) {
        this.resetTrackEventTracking();
      }

      const src = this.apiClient.url.streamTrack(mediaItem.trackId).toString();
      this.audio.src = src;

      this.currentItem = mediaItem;
    } else {
      this.audio.removeAttribute("src");
      this.audio.load();
    }
  }

  async nextTrack() {
    await this.setQueueIndex(this.queue.index + 1);
    this.play();
  }

  async previousTrack() {
    await this.setQueueIndex(this.queue.index - 1);
    this.play();
  }

  play() {
    this.audio.play();
  }

  pause() {
    this.audio.pause();
  }

  async addTracks(params: { clear?: boolean; trackId?: string }) {
    /*
    if (params.clear) {
      await this.clearQueue(false);
    }

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
    */
  }

  async addAlbumTracks(params: {
    albumId: string;
    clear?: boolean;
    trackId?: string;
  }) {
    if (params.clear) {
      await this.clearQueue(false);
    }

    const res = await this.apiClient.getAlbumTracks(params.albumId);
    if (!res.success) {
      console.log("error getting media", res.error);
      return;
    }

    const items: MediaItem[] = res.data.tracks.map(
      (t) =>
        ({
          trackId: t.id,
          name: t.name,
          album: { id: t.albumId, name: t.albumName },
          artists: t.artists.map((a) => ({ id: a.id, name: a.name })),
          coverArt: t.coverArt.small,
        }) as MediaItem,
    );

    this.queue.items.push(...items);
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
