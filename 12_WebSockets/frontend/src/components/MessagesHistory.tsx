import { ChatMessage } from "@/types/chat.types";
import { useEffect, useRef } from "react";

type MessagesHistoryProps = {
  messages: ChatMessage[];
  connectionId: string;
};

export default function MessagesHistory({
  messages,
  connectionId,
}: MessagesHistoryProps) {
  const bottomRef = useRef<HTMLDivElement | null>(null);

  // Scroll to bottom whenever messages change
  useEffect(() => {
    bottomRef.current?.scrollIntoView({ behavior: "smooth" });
  }, [messages]);

  return (
    <div className="min-w-full h-[70vh] overflow-y-auto border p-2 mb-4 rounded flex flex-col gap-2">
      {messages.map((msg, i) => {
        const isMe = msg.senderId === connectionId;

        return (
          <div
            key={i}
            className={`flex ${isMe ? "justify-end" : "justify-start"}`}
          >
            <div
              className={`max-w-[70%] px-3 py-2 rounded-xl text-sm shadow
      ${
        isMe
          ? "bg-green-500 text-white rounded-br-none"
          : "bg-gray-200 text-gray-900 rounded-bl-none"
      }`}
            >
              {!isMe && (
                <p className="text-xs font-semibold text-gray-600 mb-1">
                  {msg.senderId === "server"
                    ? "server"
                    : msg.sender || "Unknown"}
                </p>
              )}
              <p>{msg.message}</p>
              <p className="text-[10px] text-gray-500 mt-1 text-right">
                {new Date(msg.timestamp).toLocaleTimeString([], {
                  hour: "2-digit",
                  minute: "2-digit",
                  hour12: true,
                })}
              </p>
            </div>
          </div>
        );
      })}
      {/* Invisible div used as scroll target */}
      <div ref={bottomRef} />
    </div>
  );
}
