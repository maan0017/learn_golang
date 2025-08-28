import { useEffect, useRef, useState } from "react";
import type { Message, Request } from "../interfaces/message.interface";

interface MessengerProps {
  messages: Message[];
  sendMessage: (req: Request) => void;
  connectionStatus: "connecting" | "connected" | "disconnected";
  clientId: string;
  connect: () => void;
  disconnect: () => void;
  clearMessages: () => void;
}

export default function Messenger({
  messages,
  sendMessage,
  connectionStatus,
  clientId,
  connect,
  disconnect,
  clearMessages,
}: MessengerProps) {
  const messagesEndRef = useRef<HTMLDivElement>(null);

  const [inputMessage, setInputMessage] = useState<string>("");

  const scrollToBottom = () => {
    if (!messagesEndRef.current) return;
    messagesEndRef.current?.scrollIntoView({ behavior: "smooth" });
  };

  const handleKeyPress = (e: React.KeyboardEvent) => {
    if (e.key === "Enter" && !e.shiftKey) {
      e.preventDefault();
      handleSendMessage();
    }
  };

  const handleSendMessage = () => {
    const msg: Message = {
      senderId: clientId,
      recieverId: "all",
      msg: inputMessage,
      timestamp: new Date().toLocaleTimeString(),
    };
    const req: Request = {
      type: "message",
      payload: msg,
    };

    sendMessage(req);
    setInputMessage("");
  };

  const getStatusColor = () => {
    switch (connectionStatus) {
      case "connected":
        return "text-green-600";
      case "connecting":
        return "text-yellow-600";
      case "disconnected":
        return "text-red-600";
    }
  };

  const getStatusIcon = () => {
    switch (connectionStatus) {
      case "connected":
        return "🟢";
      case "connecting":
        return "🟡";
      case "disconnected":
        return "🔴";
    }
  };

  useEffect(() => {
    scrollToBottom();
  }, [messages]);

  return (
    <div className="w-full max-w-2xl mx-auto p-4 sm:p-6 bg-white shadow-lg rounded-lg">
      {/* Header */}
      <div className="mb-6">
        <h1 className="text-2xl sm:text-3xl font-bold text-gray-800 mb-4">
          WebSocket Client
        </h1>

        {/* Connection Status */}
        <div className="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-3 mb-4 p-4 bg-gray-50 rounded-lg">
          <div className="flex items-center flex-wrap gap-2">
            <span className="text-lg">{getStatusIcon()}</span>
            <span className={`font-semibold capitalize ${getStatusColor()}`}>
              {connectionStatus}
            </span>
            {clientId && (
              <span className="text-sm text-gray-600 truncate">
                ID: {clientId}
              </span>
            )}
          </div>

          <div className="flex flex-wrap gap-2">
            <button
              type="button"
              onClick={connect}
              disabled={connectionStatus === "connected"}
              className="px-4 py-2 bg-blue-500 text-white rounded-lg hover:bg-blue-600 disabled:bg-gray-300 disabled:cursor-not-allowed"
            >
              Connect
            </button>
            <button
              type="button"
              onClick={disconnect}
              disabled={connectionStatus === "disconnected"}
              className="px-4 py-2 bg-red-500 text-white rounded-lg hover:bg-red-600 disabled:bg-gray-300 disabled:cursor-not-allowed"
            >
              Disconnect
            </button>
            <button
              type="button"
              onClick={clearMessages}
              className="px-4 py-2 bg-gray-500 text-white rounded-lg hover:bg-gray-600"
            >
              Clear
            </button>
          </div>
        </div>
      </div>

      {/* Messages Display */}
      <div className="mb-4">
        <div className="bg-gray-100 border rounded-lg h-80 sm:h-96 overflow-y-auto p-3 sm:p-4 space-y-3">
          {messages.length === 0 ? (
            <div className="text-gray-500 text-center text-sm sm:text-base">
              No messages yet. Connect to start receiving messages!
            </div>
          ) : (
            messages.map((msg, index) => {
              const isMine = msg.senderId === clientId || msg.senderId === "me";
              return (
                <div
                  key={index}
                  className={`flex ${isMine ? "justify-end" : "justify-start"}`}
                >
                  <div
                    className={`max-w-[80%] md:max-w-md p-3 rounded-lg shadow-sm break-words ${
                      isMine
                        ? "bg-blue-500 text-white rounded-br-none"
                        : "bg-white text-gray-800 rounded-bl-none"
                    }`}
                  >
                    {/* Meta Info */}
                    <div className="text-[10px] sm:text-xs opacity-70 mb-1">
                      {msg.senderId && !isMine && `From: ${msg.senderId}`}{" "}
                      {msg.timestamp && `• ${msg.timestamp}`}
                    </div>
                    {/* Message Body */}
                    <div className="text-sm sm:text-base">
                      {msg.msg || "No content"}
                    </div>
                  </div>
                </div>
              );
            })
          )}
          <div ref={messagesEndRef} />
        </div>
      </div>

      {/* Message Input */}
      <div className="flex flex-col sm:flex-row gap-2">
        <input
          type="text"
          value={inputMessage}
          onChange={(e) => setInputMessage(e.target.value)}
          onKeyPress={handleKeyPress}
          placeholder="Type your message..."
          className="flex-1 px-3 sm:px-4 py-2 border border-gray-300 text-black rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500 text-sm sm:text-base"
          disabled={connectionStatus !== "connected"}
        />
        <button
          type="submit"
          onClick={handleSendMessage}
          disabled={connectionStatus !== "connected" || !inputMessage.trim()}
          className="px-5 sm:px-6 py-2 bg-green-500 text-white rounded-lg hover:bg-green-600 disabled:bg-gray-300 disabled:cursor-not-allowed text-sm sm:text-base"
        >
          Send
        </button>
      </div>

      {/* Connection Info */}
      <div className="mt-4 text-xs sm:text-sm text-gray-600">
        <p>WebSocket URL: ws://localhost:8080/ws</p>
        <p>Server must be running on port 8080</p>
      </div>
    </div>
  );
}
