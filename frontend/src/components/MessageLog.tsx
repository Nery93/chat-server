import { useEffect, useRef } from "react";
import ReactMarkdown from "react-markdown";
import remarkGfm from "remark-gfm";
import type { LogEntry } from "../types";
import "./MessageLog.css";

interface MessageLogProps {
  entries: LogEntry[];
  ownUsername: string;
}

function formatTime(ts: number) {
  return new Date(ts).toLocaleTimeString("pt-PT", {
    hour: "2-digit",
    minute: "2-digit",
  });
}

export function MessageLog({ entries, ownUsername }: MessageLogProps) {
  const bottomRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    bottomRef.current?.scrollIntoView({ block: "end" });
  }, [entries.length]);

  if (entries.length === 0) {
    return (
      <div className="message-log">
        <p className="message-log__empty">
          Sem mensagens ainda — o registo desta sala começa aqui.
        </p>
      </div>
    );
  }

  return (
    <div className="message-log">
      {entries.map((entry) => {
        if (entry.type !== "chat") {
          const verb = entry.type === "join" ? "entrou na sala" : entry.type === "leave" ? "saiu da sala" : entry.text;
          return (
            <p key={entry.id} className="message-log__event">
              <span className="message-log__time">{formatTime(entry.receivedAt)}</span>
              {entry.user} {verb}
            </p>
          );
        }

        const mine = entry.user === ownUsername;
        return (
          <div
            key={entry.id}
            className={`message-log__row ${mine ? "message-log__row--mine" : "message-log__row--other"}`}
          >
            <span className="message-log__time">{formatTime(entry.receivedAt)}</span>
            <span className="message-log__user">{mine ? "você" : entry.user}</span>
            <div className="message-log__text">
              <ReactMarkdown remarkPlugins={[remarkGfm]}>{entry.text}</ReactMarkdown>
            </div>
          </div>
        );
      })}
      <div ref={bottomRef} />
    </div>
  );
}
