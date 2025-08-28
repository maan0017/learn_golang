import { useState, useEffect, useCallback, useRef } from "react";
import type { Coords, Request } from "../interfaces/message.interface";

const COORD_INTERVAL = 50; // send every 50ms (~20fps)
const MIN_MOVEMENT_THRESHOLD = 2; // pixels - only send if moved more than this

export const useLiveMouseTracker = (
  clientId: string,
  sendMessage: (msg: Request) => void,
) => {
  const [mouseCoord, setMouseCoord] = useState({ x: 0, y: 0 });
  const lastSentRef = useRef<number>(0);
  const lastSentCoordsRef = useRef({ x: 0, y: 0 });
  const currentCoordsRef = useRef({ x: 0, y: 0 });
  const rafIdRef = useRef<number | null>(null);
  const intervalIdRef = useRef<number | null>(null);
  const hasMountedRef = useRef(false);

  // Optimized mouse tracker with RAF for smooth updates
  const trackMouse = useCallback((e: MouseEvent) => {
    // Use pageX/pageY to get coordinates relative to the document, not just viewport
    const newCoords = {
      x: e.pageX, // Document-relative X coordinate
      y: e.pageY, // Document-relative Y coordinate
    };
    currentCoordsRef.current = newCoords;

    // Cancel previous RAF to avoid stacking
    if (rafIdRef.current) {
      cancelAnimationFrame(rafIdRef.current);
    }

    // Update state on next frame for smooth UI updates
    rafIdRef.current = requestAnimationFrame(() => {
      setMouseCoord(newCoords);
    });
  }, []);

  // Optimized coordinate sender with movement threshold
  const sendCoordinates = useCallback(() => {
    const now = Date.now();
    const currentCoords = currentCoordsRef.current;
    const lastSentCoords = lastSentCoordsRef.current;

    // Check if enough time has passed
    if (now - lastSentRef.current < COORD_INTERVAL) {
      return;
    }

    // Check if mouse moved enough to warrant sending
    const deltaX = Math.abs(currentCoords.x - lastSentCoords.x);
    const deltaY = Math.abs(currentCoords.y - lastSentCoords.y);
    const moved =
      deltaX > MIN_MOVEMENT_THRESHOLD || deltaY > MIN_MOVEMENT_THRESHOLD;

    if (!moved && hasMountedRef.current) {
      return;
    }

    // Prepare and send coordinates
    const coords: Coords = {
      clientId,
      clientName: "user_" + clientId,
      CoordX: currentCoords.x,
      CoordY: currentCoords.y,
    };

    const req: Request = {
      type: "coords",
      payload: coords,
    };

    try {
      sendMessage(req);
      lastSentRef.current = now;
      lastSentCoordsRef.current = { ...currentCoords };
    } catch (error) {
      console.warn("Failed to send coordinates:", error);
    }
  }, [clientId, sendMessage]);

  // Attach/detach optimized mouse listener
  useEffect(() => {
    // Use passive listener for better performance
    const options = { passive: true };
    window.addEventListener("mousemove", trackMouse, options);

    return () => {
      window.removeEventListener("mousemove", trackMouse);
      if (rafIdRef.current) {
        cancelAnimationFrame(rafIdRef.current);
      }
    };
  }, [trackMouse]);

  // Optimized interval-based coordinate sending
  useEffect(() => {
    hasMountedRef.current = true;

    // Use a more precise interval approach
    intervalIdRef.current = setInterval(sendCoordinates, COORD_INTERVAL);

    return () => {
      if (intervalIdRef.current) {
        clearInterval(intervalIdRef.current);
      }
    };
  }, [sendCoordinates]);

  // Cleanup on unmount
  useEffect(() => {
    return () => {
      hasMountedRef.current = false;
      if (rafIdRef.current) {
        cancelAnimationFrame(rafIdRef.current);
      }
      if (intervalIdRef.current) {
        clearInterval(intervalIdRef.current);
      }
    };
  }, []);

  return mouseCoord;
};
