/* eslint-disable @typescript-eslint/no-explicit-any */
/* eslint-disable @typescript-eslint/no-unused-expressions */
const WS_URL = "ws://localhost:3000/ws";

// const WS_URL =
//   import.meta.env.VITE_WS_URL ??
//   (location.protocol === "https:" ? "wss://" : "ws://") +
//     (import.meta.env.VITE_WS_HOST ?? location.hostname + ":3000") +
//     "/ws";

export type Listener = (data: any) => void;

export class WSClient {
  private ws?: WebSocket;
  private listeners = new Set<Listener>();
  private heartbeatTimer?: number;
  private reconnectTimer?: number;
  private backoff = 500; // ms
  private readonly maxBackoff = 8000;
  private isManualClose = false;

  connect() {
    this.isManualClose = false;
    this.ws = new WebSocket(WS_URL);

    this.ws.onopen = () => {
      this.backoff = 500;
      this.startHeartbeat();
    };

    this.ws.onmessage = (ev) => {
      try {
        const msg = JSON.parse(ev.data);
        if (msg?.type === "ping") {
          this.send({ type: "pong" });
          return;
        }
        this.listeners.forEach((fn) => fn(msg));
      } catch {
        // ignore bad JSON
      }
    };

    this.ws.onclose = () => {
      this.stopHeartbeat();
      if (!this.isManualClose) this.scheduleReconnect();
    };

    this.ws.onerror = () => {
      // will trigger onclose after
    };
  }

  private scheduleReconnect() {
    if (this.reconnectTimer) return;
    const delay = this.backoff;
    this.reconnectTimer = window.setTimeout(() => {
      this.reconnectTimer && window.clearTimeout(this.reconnectTimer);
      this.reconnectTimer = undefined;
      this.backoff = Math.min(this.backoff * 2, this.maxBackoff);
      this.connect();
    }, delay) as unknown as number;
  }

  private startHeartbeat() {
    this.heartbeatTimer = window.setInterval(() => {
      this.send({ type: "ping" }); // optional; server also pings
    }, 10000) as unknown as number;
  }
  private stopHeartbeat() {
    if (this.heartbeatTimer) window.clearInterval(this.heartbeatTimer);
    this.heartbeatTimer = undefined;
  }

  send(obj: any) {
    const s = this.ws;
    if (!s || s.readyState !== WebSocket.OPEN) return false;
    s.send(JSON.stringify(obj));
    return true;
  }

  subscribe(fn: Listener) {
    this.listeners.add(fn);
    return () => this.listeners.delete(fn);
  }

  close() {
    this.isManualClose = true;
    this.stopHeartbeat();
    this.ws?.close(1000, "client closing");
  }

  readyState() {
    return this.ws?.readyState ?? WebSocket.CLOSED;
  }
}
