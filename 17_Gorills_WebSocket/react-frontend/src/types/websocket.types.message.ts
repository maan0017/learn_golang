// shared/types.ts (copy into your React app; mirror it in Go)
export type WSMessage =
  | { type: "hello"; payload: { id: string } }
  | { type: "chat"; payload: { from: string; text: string; at: number } }
  | { type: "pong"; payload?: null }
  | { type: "error"; payload: { code: string; message: string } };
