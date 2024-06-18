import { ChatSession } from "@/models/chat";
import axios from "axios";
import React, { createContext, useEffect, useState } from "react";

type IAuth = ReturnType<typeof AuthProvider>["value"];

export const AuthContext = createContext<IAuth>({} as IAuth);

export default function AuthProvider({
  children,
}: {
  children: React.ReactNode;
}) {
  const [id, setId] = useState<string>();
  const [userName, setUserName] = useState<string>();
  const [firstName, setFirstName] = useState<string>();
  const [lastName, setLastName] = useState<string>();
  const [roles, setRoles] = useState<string[]>();
  const [chats, setChats] = useState<ChatSession[]>([]);

  const fetchMeToState = async () => {
    try {
      const [userRes, chatsRes] = await Promise.all([
        axios.get(import.meta.env.VITE_PLATFORM_API_USER_INFO_URL),
        axios.get(import.meta.env.VITE_PLATFORM_API_USER_CHATS_URL),
      ]);
      const userData = userRes.data.data;
      const chatsData = chatsRes.data.data;
      setId(userData.id);
      setUserName(userData.user_name);
      setFirstName(userData.first_name);
      setLastName(userData.last_name);
      setRoles(userData.roles);
      setChats(chatsData);
    } catch (error) {
      console.error("Error fetching data:", error);
    }
  };

  useEffect(() => {
    setTimeout(() => fetchMeToState(), 1000);
  }, []);

  const value = {
    id,
    setId,
    userName,
    setUserName,
    firstName,
    setFirstName,
    lastName,
    setLastName,
    roles,
    setRoles,
    chats,
    setChats,
    fetchMeToState,
  } as const;
  
  return {
    ...(<AuthContext.Provider value={value}>{children}</AuthContext.Provider>),
    value,
  };
}
