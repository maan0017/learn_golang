import { useCallback, useEffect, useRef, useState } from "react";
import type {
  Coords,
  Message,
  Request,
  Response,
} from "../interfaces/message.interface";

export const useWebSocket = (url: string) => {
  const [socket, setSocket] = useState<WebSocket | null>(null);
  const [connectionStatus, setConnectionStatus] = useState<
    "connecting" | "connected" | "disconnected"
  >("disconnected");
  const [messages, setMessages] = useState<Message[]>([]);
  const [clientId, setClientId] = useState<string>("");
  const [logs, setLogs] = useState<string[]>([]);
  const log = (log: string) => {
    const logMsg = `[${new Date().toLocaleTimeString()}] ${log}`;
    setLogs((prev) => [...prev, logMsg]);
  };

  const [clientsCords, setClientsCoords] = useState<Coords[]>([]);

  const reconnectTimeoutRef = useRef<number | null>(null);
  const reconnectAttempts = useRef<number>(0);

  const connect = useCallback(() => {
    if (socket?.readyState === WebSocket.OPEN) return;

    setConnectionStatus("connecting");
    log(`connecting to websocket ${url}`);
    const ws = new WebSocket(url);

    ws.onopen = () => {
      console.log("WebSocket connected");
      log(`✅ successfully connected to websocket`);
      setConnectionStatus("connected");
      setSocket(ws);
      reconnectAttempts.current = 0;
    };

    ws.onmessage = (event) => {
      try {
        log(`event: ${JSON.stringify(event)}`);
        log(`event data: ${JSON.stringify(event.data)}`);
        console.log(event.data);
        const res: Response = JSON.parse(event.data);

        switch (res.type) {
          case "welcome": {
            if (res.payload?.clientId) {
              setClientId(res.payload.clientId);

              const messageWithTimestamp: Message = {
                msg: res.payload.message,
                senderId: "server",
                recieverId: res.payload.clientId,
                timestamp: new Date().toLocaleTimeString(),
              };

              setMessages((prev) => [...prev, messageWithTimestamp]);
            }
            break;
          }

          case "message": {
            const messageWithTimestamp: Message = {
              ...res.payload,
              timestamp: new Date().toLocaleTimeString(),
            };

            setMessages((prev) => [...prev, messageWithTimestamp]);
            break;
          }

          // Handling coords
          case "coords": {
            setClientsCoords((prev) => {
              if (!Array.isArray(prev)) prev = [];

              // remove old entry of same client
              const filtered = prev.filter(
                (c) => c.clientId !== res.payload.clientId,
              );

              return [...filtered, res.payload]; // keep all active clients
            });
            break;
          }

          case "notify": {
            // Example: push notifications
            console.log("Notification:", res.payload);
            break;
          }

          default:
            console.warn("Unknown message type");
        }
      } catch (error) {
        console.error("Error parsing message:", error);
        log(`Error parsing message: ${error}`);
        // Handle plain text messages
        // setMessages((prev) => [
        //   ...prev,
        //   {
        //     type: "text",
        //     data: event.data,
        //     timestamp: new Date().toLocaleTimeString(),
        //   },
        // ]);
      }
    };

    ws.onclose = (event) => {
      console.log("WebSocket disconnected:", event.code, event.reason);
      setConnectionStatus("disconnected");
      log("disconnected from websocket");
      setSocket(null);

      // Auto-reconnect with exponential backoff
      if (reconnectAttempts.current < 5) {
        const timeout = Math.min(
          1000 * Math.pow(2, reconnectAttempts.current),
          10000,
        );
        reconnectTimeoutRef.current = setTimeout(() => {
          reconnectAttempts.current++;
          connect();
        }, timeout);
      }
    };

    ws.onerror = (error) => {
      console.error("WebSocket error:", error);
      log(`WebSocket error: ${error}`);
    };

    setSocket(ws);
  }, [socket]);

  const disconnect = useCallback(() => {
    if (reconnectTimeoutRef.current) {
      clearTimeout(reconnectTimeoutRef.current);
    }
    reconnectAttempts.current = 0;

    if (socket) {
      socket.close();
      setSocket(null);
    }
    setConnectionStatus("disconnected");
  }, [socket]);

  const sendMessage = useCallback(
    (msg: Request) => {
      if (socket?.readyState === WebSocket.OPEN) {
        socket.send(JSON.stringify(msg));
      }
    },
    [socket],
  );

  const clearMessages = () => {
    setMessages([]);
  };

  useEffect(() => {
    return () => {
      if (reconnectTimeoutRef.current) {
        clearTimeout(reconnectTimeoutRef.current);
      }
      if (socket) {
        socket.close();
      }
    };
  }, [socket]);

  return {
    messages,
    sendMessage,
    connectionStatus,
    clientId,
    connect,
    disconnect,
    clearMessages,
    logs,
    clientsCords,
  };
};
