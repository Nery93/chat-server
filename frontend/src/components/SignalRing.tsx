import type { ConnectionStatus } from "../types";
import "./SignalRing.css";

interface SignalRingProps {
  status: ConnectionStatus;
  size?: "large" | "small";
}

export function SignalRing({ status, size = "small" }: SignalRingProps) {
  return (
    <span className={`signal-ring signal-ring--${size} signal-ring--${status}`} aria-hidden="true">
      <span className="signal-ring__pulse" />
      <span className="signal-ring__pulse signal-ring__pulse--delay" />
      <span className="signal-ring__core" />
    </span>
  );
}
