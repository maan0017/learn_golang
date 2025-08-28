interface LoggerProps {
  logs: string[];
}

export const Logger = ({ logs }: LoggerProps) => {
  return (
    <div className="w-1/2 p-2 h-screen overflow-y-auto">
      <h1 className="text-green-500">**Logs**</h1>
      <ul className="text-white">
        {[...logs].reverse().map((m, i) => (
          <li key={i}>
            <code>{m}</code>
          </li>
        ))}
      </ul>
    </div>
  );
};
