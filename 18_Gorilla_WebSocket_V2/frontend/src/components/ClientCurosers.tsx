import { motion, useSpring, useTransform } from "framer-motion";
import { useEffect, useMemo } from "react";
import type { Coords } from "../interfaces/message.interface";
import { FaLocationArrow } from "react-icons/fa";

type ClientCursorsProps = {
  clientId: string; // current user's ID
  clientsCursors: Coords[] | null | undefined;
};

const COLORS = [
  "text-red-500",
  "text-blue-500",
  "text-green-500",
  "text-purple-500",
  "text-orange-500",
  "text-teal-500",
  "text-pink-500",
  "text-indigo-500",
  "text-yellow-500",
  "text-emerald-500",
];

// Individual cursor component for better performance
function AnimatedCursor({
  client,
  colorClass,
}: {
  client: Coords;
  colorClass: string;
}) {
  // Use springs for ultra-smooth animation
  const springX = useSpring(0, {
    stiffness: 400,
    damping: 40,
    mass: 0.8,
  });
  const springY = useSpring(0, {
    stiffness: 400,
    damping: 40,
    mass: 0.8,
  });

  // Update spring values when coordinates change
  useEffect(() => {
    springX.set(client.CoordX);
    springY.set(client.CoordY);
  }, [client.CoordX, client.CoordY, springX, springY]);

  // Transform springs to CSS transform values
  const x = useTransform(springX, (value) => `${value}px`);
  const y = useTransform(springY, (value) => `${value}px`);

  return (
    <motion.div
      className="absolute pointer-events-none will-change-transform"
      style={{
        left: x,
        top: y,
        // Precise positioning - cursor tip points exactly where it should
        transform: "translate(-2px, -2px)",
      }}
      initial={{ opacity: 0, scale: 0.8 }}
      animate={{ opacity: 1, scale: 1 }}
      exit={{ opacity: 0, scale: 0.8 }}
      transition={{
        opacity: { duration: 0.15 },
        scale: { duration: 0.15 },
      }}
    >
      {/* Cursor icon with precise positioning */}
      <div className="relative">
        <FaLocationArrow
          className={`w-4 h-4 ${colorClass} drop-shadow-sm`}
          style={{
            // Fine-tune the cursor tip position
            transform: "rotate(-85deg)",
            filter: "drop-shadow(0 1px 2px rgba(0,0,0,0.3))",
          }}
        />

        {/* Name label with smooth follow animation */}
        <motion.div
          className="absolute top-5 left-2 whitespace-nowrap"
          initial={{ opacity: 0, y: -5 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ delay: 0.1, duration: 0.2 }}
        >
          <span
            className={`
            text-xs font-medium px-2 py-1 rounded-md
            bg-white/90 backdrop-blur-sm border border-gray-200
            shadow-sm text-gray-700
            ${colorClass.replace("text-", "border-l-4 border-l-")}
          `}
          >
            {client.clientName ?? client.clientId}
          </span>
        </motion.div>
      </div>
    </motion.div>
  );
}

export default function ClientCursors({
  clientId,
  clientsCursors,
}: ClientCursorsProps) {
  // Memoize filtered clients to prevent unnecessary re-renders
  const filteredClients = useMemo(() => {
    if (!Array.isArray(clientsCursors)) return [];
    return clientsCursors.filter((client) => client.clientId !== clientId);
  }, [clientsCursors, clientId]);

  // Create a stable color mapping based on clientId
  const getColorForClient = useMemo(() => {
    const colorMap = new Map<string, string>();
    return (clientId: string) => {
      if (!colorMap.has(clientId)) {
        // Generate consistent color based on client ID hash
        const hash = clientId
          .split("")
          .reduce((acc, char) => char.charCodeAt(0) + ((acc << 5) - acc), 0);
        const colorIndex = Math.abs(hash) % COLORS.length;
        colorMap.set(clientId, COLORS[colorIndex]);
      }
      return colorMap.get(clientId)!;
    };
  }, []);

  if (filteredClients.length === 0) return null;

  return (
    <div className="absolute inset-0 pointer-events-none z-[9999] overflow-hidden">
      {filteredClients.map((client) => (
        <AnimatedCursor
          key={client.clientId}
          client={client}
          colorClass={getColorForClient(client.clientId)}
        />
      ))}
    </div>
  );
}
