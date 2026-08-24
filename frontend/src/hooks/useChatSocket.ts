import { useCallback, useRef, useState } from "react";
import type { ConnectionStatus, LogEntry, ServerMessage } from "../types";

const WS_BASE = import.meta.env.VITE_WS_URL ?? "ws://localhost:8080";
const TYPING_EXPIRY_MS = 2500;
const TYPING_SEND_THROTTLE_MS = 2000;

function makeId() {
  return `${Date.now()}-${Math.random().toString(36).slice(2, 8)}`;
}

export function useChatSocket() {
  const [status, setStatus] = useState<ConnectionStatus>("idle");
  const [messages, setMessages] = useState<LogEntry[]>([]);
  const [onlineUsers, setOnlineUsers] = useState<string[]>([]);
  const [typingUsers, setTypingUsers] = useState<string[]>([]);
  const [username, setUsername] = useState("");
  const socketRef = useRef<WebSocket | null>(null);
  const leftIntentionallyRef = useRef(false);
  const typingTimersRef = useRef(new Map<string, ReturnType<typeof setTimeout>>());
  const lastTypingSentRef = useRef(0);

  const connect = useCallback((room: string, user: string) => {
    if (socketRef.current) return;

    leftIntentionallyRef.current = false;
    setStatus("connecting");
    setMessages([]);
    setOnlineUsers([]);
    setTypingUsers([]);
    typingTimersRef.current.forEach((timer) => clearTimeout(timer));
    typingTimersRef.current.clear();
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

      if (parsed.type === "userlist") {
        setOnlineUsers(parsed.users ?? []);
        return;
      }

      if (parsed.type === "typing") {
        if (parsed.user === user) return;

        const timers = typingTimersRef.current;
        const existing = timers.get(parsed.user);
        if (existing) clearTimeout(existing);

        setTypingUsers((prev) => (prev.includes(parsed.user) ? prev : [...prev, parsed.user]));

        timers.set(
          parsed.user,
          setTimeout(() => {
            timers.delete(parsed.user);
            setTypingUsers((prev) => prev.filter((u) => u !== parsed.user));
          }, TYPING_EXPIRY_MS),
        );
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
      setOnlineUsers([]);
      typingTimersRef.current.forEach((timer) => clearTimeout(timer));
      typingTimersRef.current.clear();
      setTypingUsers([]);
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

  const sendTyping = useCallback(() => {
    const socket = socketRef.current;
    if (!socket || socket.readyState !== WebSocket.OPEN) return;

    const now = Date.now();
    if (now - lastTypingSentRef.current < TYPING_SEND_THROTTLE_MS) return;
    lastTypingSentRef.current = now;

    socket.send(JSON.stringify({ type: "typing", user: username, text: "" }));
  }, [username]);

  return { status, messages, onlineUsers, typingUsers, username, connect, disconnect, send, sendTyping };
}
