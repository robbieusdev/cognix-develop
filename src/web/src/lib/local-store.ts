import { useState } from "react";

type AnyObject = Record<string, any>;

interface UseLocalStorage {
  set: (key: string, value: AnyObject) => void;
  get: (key: string, initValue: AnyObject) => void;
  remove: (key: string) => void;
  clear: () => void;
}
export function useLocalStorage(): UseLocalStorage {
  const set = (key: string, value: AnyObject) => {
    localStorage.setItem(key, JSON.stringify(value));
  };
  const get = (key: string, initValue: AnyObject) => {
    const [value] = useState<AnyObject | null>(() => {
      const storedValue = localStorage.getItem(key);
      return storedValue ? JSON.parse(storedValue) : initValue;
    });
    return value;
  };
  const remove = (key: string) => {
    localStorage.removeItem(key);
  };
  const clear = () => {
    localStorage.clear();
  };
  return { get, set, remove, clear };
}
