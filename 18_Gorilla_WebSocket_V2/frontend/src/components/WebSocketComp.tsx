import React, { useState, useEffect, useRef, useCallback } from "react";
import type { Message } from "../interfaces/message.interface";

const WebSocketClient: React.FC = () => {
  const [socket, setSocket] = useState<WebSocket | null>(null);
  const [connectionStatus, setConnectionStatus] = useState<
    "connecting" | "connected" | "disconnected"
  >("disconnected");
  const [messages, setMessages] = useState<Message[]>([]);
  const [inputMessage, setInputMessage] = useState<string>("");
  const [clientId, setClientId] = useState<string>("");

  const messagesEndRef = useRef<HTMLDivElement>(null);
  const reconnectTimeoutRef = useRef<NodeJS.Timeout>(null);
  const reconnectAttempts = useRef<number>(0);

  const scrollToBottom = () => {
    messagesEndRef.current?.scrollIntoView({ behavior: "smooth" });
  };

  useEffect(() => {
    scrollToBottom();
  }, [messages]);

  const connect = useCallback(() => {
    if (socket?.readyState === WebSocket.OPEN) {
      return;
    }

    setConnectionStatus("connecting");

    const ws = new WebSocket("ws://localhost:8080/ws");

    ws.onopen = () => {
      console.log("WebSocket connected");
      setConnectionStatus("connected");
      setSocket(ws);
      reconnectAttempts.current = 0;
    };

    ws.onmessage = (event) => {
      try {
        const message: Message = JSON.parse(event.data);

        // Add timestamp to message
        const messageWithTimestamp = {
          ...message,
          timestamp: new Date().toLocaleTimeString(),
        };

        setMessages((prev) => [...prev, messageWithTimestamp]);

        // Set client ID from welcome message
        if (message.type === "welcome" && message.clientId) {
          setClientId(message.clientId);
        }
      } catch (error) {
        console.error("Error parsing message:", error);
        // Handle plain text messages
        setMessages((prev) => [
          ...prev,
          {
            type: "text",
            data: event.data,
            timestamp: new Date().toLocaleTimeString(),
          },
        ]);
      }
    };

    ws.onclose = (event) => {
      console.log("WebSocket disconnected:", event.code, event.reason);
      setConnectionStatus("disconnected");
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

  const sendMessage = useCallback(() => {
    if (socket?.readyState === WebSocket.OPEN && inputMessage.trim()) {
      socket.send(inputMessage.trim());
      setInputMessage("");
    }
  }, [socket, inputMessage]);

  const handleKeyPress = (e: React.KeyboardEvent) => {
    if (e.key === "Enter" && !e.shiftKey) {
      e.preventDefault();
      sendMessage();
    }
  };

  const clearMessages = () => {
    setMessages([]);
  };

  // Cleanup on unmount
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

  return (
    <div className="max-w-4xl mx-auto p-6 bg-white shadow-lg rounded-lg">
      <div className="mb-6">
        <h1 className="text-3xl font-bold text-gray-800 mb-4">
          WebSocket Client
        </h1>

        {/* Connection Status */}
        <div className="flex items-center justify-between mb-4 p-4 bg-gray-50 rounded-lg">
          <div className="flex items-center space-x-3">
            <span className="text-lg">{getStatusIcon()}</span>
            <span className={`font-semibold capitalize ${getStatusColor()}`}>
              {connectionStatus}
            </span>
            {clientId && (
              <span className="text-sm text-gray-600">ID: {clientId}</span>
            )}
          </div>

          <div className="space-x-2">
            <button
              onClick={connect}
              disabled={connectionStatus === "connected"}
              className="px-4 py-2 bg-blue-500 text-white rounded hover:bg-blue-600 disabled:bg-gray-300 disabled:cursor-not-allowed"
            >
              Connect
            </button>
            <button
              onClick={disconnect}
              disabled={connectionStatus === "disconnected"}
              className="px-4 py-2 bg-red-500 text-white rounded hover:bg-red-600 disabled:bg-gray-300 disabled:cursor-not-allowed"
            >
              Disconnect
            </button>
            <button
              onClick={clearMessages}
              className="px-4 py-2 bg-gray-500 text-white rounded hover:bg-gray-600"
            >
              Clear Messages
            </button>
          </div>
        </div>
      </div>

      {/* Messages Display */}
      <div className="mb-4">
        <div className="bg-gray-100 border rounded-lg h-96 overflow-y-auto p-4 space-y-2">
          {messages.length === 0 ? (
            <div className="text-gray-500 text-center">
              No messages yet. Connect to start receiving messages!
            </div>
          ) : (
            messages.map((msg, index) => (
              <div key={index} className="bg-white p-3 rounded shadow-sm">
                <div className="flex justify-between items-start">
                  <div className="flex-1">
                    <div className="text-xs text-gray-500 mb-1">
                      {msg.type.toUpperCase()}
                      {msg.from && ` • From: ${msg.from}`}
                      {msg.timestamp && ` • ${msg.timestamp}`}
                    </div>
                    <div className="text-gray-800">
                      {msg.message || msg.data || "No content"}
                    </div>
                  </div>
                </div>
              </div>
            ))
          )}
          <div ref={messagesEndRef} />
        </div>
      </div>

      {/* Message Input */}
      <div className="flex space-x-2">
        <input
          type="text"
          value={inputMessage}
          onChange={(e) => setInputMessage(e.target.value)}
          onKeyPress={handleKeyPress}
          placeholder="Type your message here..."
          className="flex-1 px-4 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
          disabled={connectionStatus !== "connected"}
        />
        <button
          onClick={sendMessage}
          disabled={connectionStatus !== "connected" || !inputMessage.trim()}
          className="px-6 py-2 bg-green-500 text-white rounded-lg hover:bg-green-600 disabled:bg-gray-300 disabled:cursor-not-allowed"
        >
          Send
        </button>
      </div>

      {/* Connection Info */}
      <div className="mt-4 text-sm text-gray-600">
        <p>WebSocket URL: ws://localhost:8080/ws</p>
        <p>Server must be running on port 8080</p>
      </div>
    </div>
  );
};

export default WebSocketClient;
