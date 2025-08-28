"use client";

import { ChatMessage } from "@/types/chat.types";
import { useEffect, useRef, useState } from "react";
import ReconnectingWebSocket from "reconnecting-websocket";

export function useWebSocket(url: string) {
  const socketRef = useRef<ReconnectingWebSocket | null>(null);
  const [messages, setMessages] = useState<ChatMessage[]>([]);
  const [isConnected, setIsConnected] = useState(false);
  const [connectionId, setConnectionId] = useState<string>("");

  useEffect(() => {
    // Configure reconnecting-websocket
    const socket = new ReconnectingWebSocket(url, [], {
      //   reconnectInterval: 2000, // retry every 2s
      maxRetries: 10, // stop after 10 retries (optional)
    });

    socketRef.current = socket;

    socket.addEventListener("open", () => {
      console.log("✅ Connected to WebSocket");
      setIsConnected(true);
    });

    socket.addEventListener("message", (event) => {
      console.log("event :> ", event);
      try {
        const parsed = JSON.parse(event.data);

        // const msg: ChatMessage = {
        //   id: parsed.id,
        //   senderId: parsed.senderId,
        //   sender: parsed.sender ?? undefined,
        //   message: parsed.message,
        //   timestamp: new Date(parsed.timestamp).toISOString(),
        //   system: parsed.system ?? false,
        // };

        // if the server sends a "welcome" event with your ID
        // if (
        //   msg.system &&
        //   msg.senderId === "server" &&
        //   msg.message.startsWith("WELCOME")
        // ) {
        //   setConnectionId(msg.id);
        // }

        console.log("server res: ", parsed);

        // setMessages((prev) => [...prev, msg]);
      } catch (error) {
        console.error("❌ Failed to parse message:", event.data, error);
      }
    });

    socket.addEventListener("close", () => {
      console.log("❌ Disconnected, will try to reconnect...");
      setIsConnected(false);
    });

    socket.addEventListener("error", (err) => {
      console.error("WebSocket error:", err);
    });

    return () => {
      socket.close();
    };
  }, [url]);

  const sendMessage = (text: string) => {
    if (socketRef.current && socketRef.current.readyState === WebSocket.OPEN) {
      // const msg: ChatMessage = {
      //   id: crypto.randomUUID(),
      //   senderId: connectionId,
      //   message: text,
      //   timestamp: new Date().toISOString(),
      // };
      console.log("sending msg: ", JSON.stringify(text));
      socketRef.current.send(text);
    }
  };

  return { messages, sendMessage, isConnected, connectionId };
}
