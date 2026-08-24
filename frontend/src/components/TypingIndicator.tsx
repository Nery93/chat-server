import "./TypingIndicator.css";

interface TypingIndicatorProps {
  users: string[];
}

export function TypingIndicator({ users }: TypingIndicatorProps) {
  if (users.length === 0) return <div className="typing-indicator typing-indicator--empty" />;

  const label =
    users.length === 1
      ? `${users[0]} está a escrever...`
      : users.length === 2
        ? `${users[0]} e ${users[1]} estão a escrever...`
        : `${users.length} pessoas estão a escrever...`;

  return (
    <div className="typing-indicator">
      <span className="typing-indicator__dots">
        <span />
        <span />
        <span />
      </span>
      {label}
    </div>
  );
}
