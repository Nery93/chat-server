export type MessageType = "chat" | "join" | "leave" | "system";

export interface ServerMessage {
  type: MessageType;
  user: string;
  text: string;
}

export interface LogEntry extends ServerMessage {
  id: string;
  receivedAt: number;
}

export type ConnectionStatus = "idle" | "connecting" | "online" | "offline";
