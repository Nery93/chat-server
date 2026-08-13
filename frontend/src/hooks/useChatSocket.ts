import { useCallback, useRef, useState } from "react";
import type { ConnectionStatus, LogEntry, ServerMessage } from "../types";

const WS_BASE = import.meta.env.VITE_WS_URL ?? "ws://localhost:8080";

function makeId() {
  return `${Date.now()}-${Math.random().toString(36).slice(2, 8)}`;
}

export function useChatSocket() {
  const [status, setStatus] = useState<ConnectionStatus>("idle");
  const [messages, setMessages] = useState<LogEntry[]>([]);
  const [username, setUsername] = useState("");
  const socketRef = useRef<WebSocket | null>(null);
  const leftIntentionallyRef = useRef(false);

  const connect = useCallback((room: string, user: string) => {
    if (socketRef.current) return;

    leftIntentionallyRef.current = false;
    setStatus("connecting");
    setMessages([]);
    setUsername(user);

    const url = `${WS_BASE}/ws/${encodeURIComponent(room)}?user=${encodeURIComponent(user)}`;
    const socket = new WebSocket(url);
    socketRef.current = socket;

    socket.onopen = () => setStatus("online");

    socket.onmessage = (event) => {
      let parsed: ServerMessage;
      try {
        parsed = JSON.parse(event.data);
      } catch {
        return;
      }
      setMessages((prev) => [
        ...prev,
        { ...parsed, id: makeId(), receivedAt: Date.now() },
      ]);
    };

    socket.onclose = () => {
      setStatus(leftIntentionallyRef.current ? "idle" : "offline");
      socketRef.current = null;
    };

    socket.onerror = () => {
      setStatus("offline");
    };
  }, []);

  const disconnect = useCallback(() => {
    leftIntentionallyRef.current = true;
    socketRef.current?.close();
  }, []);

  const send = useCallback((text: string) => {
    const socket = socketRef.current;
    if (!socket || socket.readyState !== WebSocket.OPEN || !text.trim()) return;
    socket.send(JSON.stringify({ type: "chat", user: username, text }));
  }, [username]);

  return { status, messages, username, connect, disconnect, send };
}
