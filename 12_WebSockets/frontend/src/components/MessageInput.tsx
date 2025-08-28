"use client";

import { SendHorizontal } from "lucide-react";
import { useRef, useEffect } from "react";

type MessageInputProps = {
  input: string;
  setInput: React.Dispatch<React.SetStateAction<string>>;
  handleSendMessage: () => void;
};

export default function MessageInput({
  input,
  setInput,
  handleSendMessage,
}: MessageInputProps) {
  const textareaRef = useRef<HTMLTextAreaElement>(null);

  // Auto resize textarea height
  useEffect(() => {
    const textarea = textareaRef.current;
    if (textarea) {
      textarea.style.height = "auto"; // reset first
      textarea.style.height = `${textarea.scrollHeight}px`;
    }
  }, [input]);

  const handleSubmit = (e?: React.FormEvent) => {
    e?.preventDefault();
    if (input.trim() !== "") {
      handleSendMessage();
    }
  };

  const handleKeyDown = (e: React.KeyboardEvent<HTMLTextAreaElement>) => {
    if (e.key === "Enter" && !e.shiftKey) {
      e.preventDefault(); // prevent new line
      handleSubmit(); // send message
    }
  };

  return (
    <form onSubmit={handleSubmit} className="max-w-3xl w-full mx-auto p-2">
      <div className="flex items-end bg-white/10 border border-green-500 rounded-xl overflow-hidden px-2 py-1">
        <textarea
          ref={textareaRef}
          value={input}
          onChange={(e) => setInput(e.target.value)}
          onKeyDown={handleKeyDown}
          placeholder="Enter something..."
          className="my-auto flex-grow resize-none px-2 py-1 outline-none bg-transparent text-white placeholder:text-white/50 text-sm"
          rows={1}
        />
        <button
          title="Send message"
          type="submit"
          className="p-2 text-green-500 hover:text-green-600 cursor-pointer"
        >
          <SendHorizontal size={24} />
        </button>
      </div>
    </form>
  );
}
