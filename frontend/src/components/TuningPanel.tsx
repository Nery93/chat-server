import { useState, type FormEvent } from "react";
import { SignalRing } from "./SignalRing";
import type { ConnectionStatus } from "../types";
import "./TuningPanel.css";

interface TuningPanelProps {
  status: ConnectionStatus;
  onConnect: (room: string, user: string) => void;
}

export function TuningPanel({ status, onConnect }: TuningPanelProps) {
  const [name, setName] = useState("");
  const [room, setRoom] = useState("geral");

  const submit = (e: FormEvent) => {
    e.preventDefault();
    if (!name.trim() || !room.trim()) return;
    onConnect(room.trim(), name.trim());
  };

  const connecting = status === "connecting";

  return (
    <div className="tuning-panel">
      <SignalRing status={status} size="large" />
      <h1 className="tuning-panel__title">chat-server</h1>
      <p className="tuning-panel__subtitle">test console</p>

      <form className="tuning-panel__form" onSubmit={submit}>
        <label className="tuning-panel__field">
          <span>nome</span>
          <input
            value={name}
            onChange={(e) => setName(e.target.value)}
            placeholder="como te chamas?"
            disabled={connecting}
            autoFocus
          />
        </label>
        <label className="tuning-panel__field">
          <span>sala</span>
          <input
            value={room}
            onChange={(e) => setRoom(e.target.value)}
            placeholder="geral"
            disabled={connecting}
          />
        </label>
        <button type="submit" disabled={connecting || !name.trim() || !room.trim()}>
          {connecting ? "a sintonizar..." : "entrar na sala"}
        </button>
      </form>

      {status === "offline" && (
        <p className="tuning-panel__hint">
          Sinal perdido. Confirma que o servidor está a correr e tenta outra vez.
        </p>
      )}
    </div>
  );
}
