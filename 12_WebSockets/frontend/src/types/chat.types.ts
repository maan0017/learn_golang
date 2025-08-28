export interface ChatMessage {
  id: string;
  senderId: string; // UUID from backend
  sender?: string; // optional username
  message: string;
  timestamp: string; // store as ISO string
  system?: boolean; // for "welcome"/"joined" messages
}
