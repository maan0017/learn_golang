"use client";

import MessageInput from "@/components/MessageInput";
import MessagesHistory from "@/components/MessagesHistory";
import { useWebSocket } from "@/hooks/useWebsocket";
import { useEffect, useState } from "react";

export default function Home() {
  const { messages, sendMessage, isConnected, connectionId } = useWebSocket(
    "ws://localhost:8080/ws",
  );
  const [input, setInput] = useState<string>("");

  const handleSendMessage = () => {
    if (input.trim()) {
      sendMessage(input);
      setInput("");
    }
  };

  useEffect(() => {
    const handleMouseMove = (e: MouseEvent) => {
      sendMessage(
        JSON.stringify({
          type: "coords",
          data: {
            CoordX: e.clientX,
            CoordY: e.clientY,
          },
        }),
      );
    };

    // attach listener
    window.addEventListener("mousemove", handleMouseMove);

    // cleanup
    return () => {
      window.removeEventListener("mousemove", handleMouseMove);
    };
  }, []);

  return (
    <main className="w-full h-screen flex flex-col justify-end items-center p-10">
      <h1 className="text-xl font-bold mb-2">Next.js + Go Chat</h1>
      <span>My Id :{connectionId}</span>
      <p className="text-sm mb-4">
        Status: {isConnected ? "🟢 Connected" : "🔴 Disconnected"}
      </p>

      <MessagesHistory messages={messages} connectionId={connectionId} />

      <MessageInput
        input={input}
        setInput={setInput}
        handleSendMessage={handleSendMessage}
      />
    </main>
  );
}
