import { useState, type FormEvent } from "react";
import "./Composer.css";

interface ComposerProps {
  onSend: (text: string) => void;
  onTyping: () => void;
  disabled: boolean;
}

export function Composer({ onSend, onTyping, disabled }: ComposerProps) {
  const [text, setText] = useState("");

  const submit = (e: FormEvent) => {
    e.preventDefault();
    const trimmed = text.trim();
    if (!trimmed) return;
    onSend(trimmed);
    setText("");
  };

  return (
    <form className="composer" onSubmit={submit}>
      <input
        className="composer__input"
        type="text"
        value={text}
        onChange={(e) => {
          setText(e.target.value);
          if (e.target.value.trim()) onTyping();
        }}
        placeholder="escreve uma mensagem..."
        disabled={disabled}
        autoComplete="off"
      />
      <button className="composer__send" type="submit" disabled={disabled || !text.trim()}>
        enviar
      </button>
    </form>
  );
}
