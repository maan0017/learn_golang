import ClientCurosers from "./components/ClientCurosers";
import { Logger } from "./components/Logger";
import Messenger from "./components/Messenger";
import { useLiveMouseTracker } from "./hooks/useLiveMouseTracker";
import { useWebSocket } from "./hooks/useWebSocket";

function App() {
  const {
    messages,
    sendMessage,
    connectionStatus,
    clientId,
    connect,
    disconnect,
    clearMessages,
    logs,
    clientsCords,
  } = useWebSocket("ws://localhost:8080/ws");

  useLiveMouseTracker(clientId, sendMessage);

  return (
    <div className="w-full min-h-screen bg-black/90 text-white flex flex-col md:flex-row relative">
      {/* Left/Main section - Messenger */}
      <div className="flex-1 p-4 overflow-y-auto">
        <Messenger
          messages={messages}
          sendMessage={sendMessage}
          connectionStatus={connectionStatus}
          clientId={clientId}
          connect={connect}
          disconnect={disconnect}
          clearMessages={clearMessages}
        />
      </div>

      {/* Right/Sidebar - Logger (hidden on small screens) */}
      <Logger logs={logs} />

      {/* Overlay - Live client cursors */}
      <ClientCurosers clientId={clientId} clientsCursors={clientsCords} />
    </div>
  );
}

export default App;
