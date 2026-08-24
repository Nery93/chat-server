import { useState } from "react";
import { useChatSocket } from "./hooks/useChatSocket";
import { SignalRing } from "./components/SignalRing";
import { TuningPanel } from "./components/TuningPanel";
import { MessageLog } from "./components/MessageLog";
import { Composer } from "./components/Composer";
import { OnlineUsers } from "./components/OnlineUsers";
import { TypingIndicator } from "./components/TypingIndicator";
import "./App.css";

function App() {
  const { status, messages, onlineUsers, typingUsers, username, connect, disconnect, send, sendTyping } =
    useChatSocket();
  const [room, setRoom] = useState("");

  const handleConnect = (roomName: string, name: string) => {
    setRoom(roomName);
    connect(roomName, name);
  };

  const isConnected = status === "online" || status === "connecting";

  if (!isConnected) {
    return (
      <div className="app app--tuning">
        <TuningPanel status={status} onConnect={handleConnect} />
      </div>
    );
  }

  return (
    <div className="app">
      <div className="console">
        <header className="console__header">
          <div className="console__room">
            <SignalRing status={status} size="small" />
            <span>{room}</span>
          </div>
          <div className="console__identity">
            <span className="console__username">~{username}</span>
            <button className="console__leave" onClick={disconnect}>
              sair
            </button>
          </div>
        </header>

        <OnlineUsers users={onlineUsers} ownUsername={username} />

        <MessageLog entries={messages} ownUsername={username} />

        <TypingIndicator users={typingUsers} />

        <Composer onSend={send} onTyping={sendTyping} disabled={status !== "online"} />
      </div>
    </div>
  );
}

export default App;
