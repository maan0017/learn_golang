import { useState, type FormEvent } from "react";
import { useWebSocket } from "./hooks/useWebSocket";

export default function App() {
  const [inputMsg, setInputMsg] = useState<string>("");

  const { ready, messages, send, logs } = useWebSocket(
    "ws://localhost:8080/ws",
  );

  const handleSendMsg = (e: FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    send({
      type: "message",
      payload: { from: "asur", msg: inputMsg, timestamps: Date.now() },
    });
    setInputMsg("");
  };

  return (
    <div className="w-full min-h-screen font-serif p-10 bg-black">
      <h1 className="text-white">
        WebSocket {ready ? "🟢 Connected" : "🔴 Disconnected"}
      </h1>

      <form onSubmit={handleSendMsg}>
        <input
          name="msg"
          value={inputMsg}
          onChange={(e) => setInputMsg(e.target.value)}
          placeholder="type a message..."
          className="m-1 px-1 text-white border"
        />
        <button
          type="submit"
          disabled={!ready}
          className="w-20 px-5 py-1 m-2 bg-green-800 hover:bg-green-500 cursor-pointer duration-500"
        >
          Send
        </button>
      </form>

      <div className="w-full flex">
        {/* messages */}
        <div className="w-1/2">
          <h1 className="text-green-500">**Messages**</h1>
          <ul className="text-white">
            {[...messages].reverse().map((m, i) => (
              <li key={i}>
                <code>{JSON.stringify(m)}</code>
              </li>
            ))}
          </ul>
        </div>

        {/* logs */}
        <div className="w-1/2">
          <h1 className="text-green-500">**Logs**</h1>
          <ul className="text-white">
            {[...logs].reverse().map((m, i) => (
              <li key={i}>
                <code>{m}</code>
              </li>
            ))}
          </ul>
        </div>
      </div>
    </div>
  );
}
