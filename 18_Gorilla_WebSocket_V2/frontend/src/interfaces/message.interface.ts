export type Request =
  | { type: "welcome"; payload: Welcome }
  | { type: "message"; payload: Message }
  | { type: "coords"; payload: Coords }
  | { type: "notify"; payload: Notification };

export type Response =
  | { type: "welcome"; payload: Welcome }
  | { type: "message"; payload: Message }
  | { type: "coords"; payload: Coords }
  | { type: "notify"; payload: Notification };

export interface Welcome {
  clientId: string;
  clientName: string;
  message: string;
}

export interface Message {
  id?: string;
  msg?: string;
  senderId?: string;
  recieverId?: string;
  timestamp?: string;
}

export interface Coords {
  clientId: string;
  clientName: string;
  CoordX: number;
  CoordY: number;
}

export interface Notification {
  category: "friend-request" | "update";
  heading: string;
  content: string;
  senderId?: string;
  recieverId?: string;
  timestamp?: string;
}
