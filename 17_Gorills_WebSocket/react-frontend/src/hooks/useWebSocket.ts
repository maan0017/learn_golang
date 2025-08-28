/* eslint-disable @typescript-eslint/no-unused-expressions */
import { useEffect, useRef, useState } from "react";

export function useWebSocket(url: string) {
  const wsRef = useRef<WebSocket | null>(null);
  const [ready, setReady] = useState(false);
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  const [messages, setMessages] = useState<any[]>([]);
  const [logs, setLogs] = useState<string[]>([]);

  // --- config
  const heartbeatInterval = 10_000; // 10s
  const reconnectBaseDelay = 500; // ms
  const reconnectMaxDelay = 8000; // ms

  // --- refs for timers and backoff
  const hbTimer = useRef<number | null>(null);
  const reconnectTimer = useRef<number | null>(null);
  const backoff = useRef(reconnectBaseDelay);

  // helper to push log entries
  const log = (msg: string) => {
    const entry = `[${new Date().toLocaleTimeString()}] ${msg}`;
    setLogs((prev) => [...prev, entry]);
  };

  const connect = () => {
    if (wsRef.current?.readyState === WebSocket.OPEN) return;

    log(`Connecting to ${url} ...`);
    const socket = new WebSocket(url);
    wsRef.current = socket;

    socket.onopen = () => {
      setReady(true);
      backoff.current = reconnectBaseDelay;
      startHeartbeat();
      log("✅ Connection established");
    };

    socket.onmessage = (ev) => {
      console.log("event:", ev);
      try {
        const msg = JSON.parse(ev.data);
        console.log("msg:", msg);
        if (
          msg?.type === "ping" ||
          (msg?.type === "message" && msg?.data?.type === "ping")
        ) {
          send({ type: "pong" });
          log("↔️ Received ping, sent pong");
          return;
        }
        setMessages((prev) => [...prev, msg]);
        log(`📩 Message received: ${JSON.stringify(msg)}`);
      } catch {
        setMessages((prev) => [...prev, ev.data]);
        log(`📩 Non-JSON message received: ${ev.data}`);
      }
    };

    socket.onclose = (ev) => {
      setReady(false);
      stopHeartbeat();
      log(
        `⚠️ Connection closed (code=${ev.code}, reason=${ev.reason || "N/A"})`,
      );
      scheduleReconnect();
    };

    socket.onerror = (err) => {
      log(`❌ WebSocket error: ${JSON.stringify(err)}`);
      socket.close(); // triggers onclose
    };
  };

  const scheduleReconnect = () => {
    if (reconnectTimer.current) return;
    const delay = backoff.current;
    log(`🔄 Scheduling reconnect in ${delay}ms`);
    reconnectTimer.current = window.setTimeout(() => {
      reconnectTimer.current && window.clearTimeout(reconnectTimer.current);
      reconnectTimer.current = null;
      backoff.current = Math.min(backoff.current * 2, reconnectMaxDelay);
      connect();
    }, delay) as unknown as number;
  };

  const startHeartbeat = () => {
    hbTimer.current = window.setInterval(() => {
      send({ type: "ping" });
      log("💓 Sent ping");
    }, heartbeatInterval) as unknown as number;
  };

  const stopHeartbeat = () => {
    if (hbTimer.current) window.clearInterval(hbTimer.current);
    hbTimer.current = null;
  };

  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  const send = (obj: any) => {
    const socket = wsRef.current;
    if (socket && socket.readyState === WebSocket.OPEN) {
      socket.send(JSON.stringify(obj));
      log(`📤 Sent: ${JSON.stringify(obj)}`);
      return true;
    }
    log("⚠️ Tried to send but socket not open");
    return false;
  };

  useEffect(() => {
    connect();
    return () => {
      stopHeartbeat();
      reconnectTimer.current && window.clearTimeout(reconnectTimer.current);
      wsRef.current?.close(1000, "client closing");
      log("🛑 Hook cleanup: connection closed");
    };
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [url]);

  return { ready, messages, send, logs };
}
