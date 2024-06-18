import { useEffect, useLayoutEffect, useRef, useState } from "react";
import axios from "axios";
import { ChatMessage } from "@/models/chat";
import { useParams } from "react-router-dom";
import { v4 as uuidv4 } from "uuid";
import { Persona } from "@/models/settings";
import { router } from "@/main";
import { toast } from "react-toastify";
import { RetrievedKnowledge } from "./components/retrieved-knowledge";
import { useMessages } from "@/context/ChatContext";
import { ChatInput } from "./chat-input";
import MessagesList from "./messages-list";
import { PersonaSelection } from "./persona-list-section";

const API_BASE_URL = import.meta.env.VITE_PLATFORM_API_BASE_URL;

export function ChatComponent() {
  const { messages, setMessages, addMessage } = useMessages();
  const [personas, setPersonas] = useState<Persona[]>([]);
  const [selectedPersona, setSelectedPersona] = useState<string>("");
  const [newMessage, setNewMessage] = useState<ChatMessage | null>(null);
  const [isDeactivateSendingButton, setIsDeactivateSendingButton] = useState(false);
  const textInputRef = useRef<HTMLInputElement>(null);

  const { chatId } = useParams<{ chatId?: string }>();

  useLayoutEffect(() => {
    if (chatId) {
      fetchMessages();
    } else {
      fetchPersonas();
      setMessages([]);
    }
  }, [chatId]);

  useEffect(() => {
    if (newMessage && newMessage.message) {
      const intervalId = setInterval(() => {
        updateMessageContent();
      }, 25);

      return () => clearInterval(intervalId);
    }
  }, [newMessage]);

  async function fetchMessages() {
    try {
      const response = await axios.get(`${API_BASE_URL}/chat/detail/${chatId}`);
      if (response.status === 200) {
        setMessages(response.data.data.messages);
      } else {
        setMessages([]);
      }
    } catch (error) {
      console.error("Error fetching messages:", error);
    }
  }

  async function fetchPersonas() {
    try {
      const response = await axios.get(`${API_BASE_URL}/llm/list`);
      if (response.status === 200) {
        setPersonas(response.data.data);
      } else {
        setPersonas([]);
      }
    } catch (error) {
      console.error("Error fetching personas:", error);
      setPersonas([]);
    }
  }

  async function createChat(text: string) {
    if (!selectedPersona) {
      if (!personas.length) {
        toast.error("You don't have any available Assistants yet.");
        return;
      }
    }

    try {
      const response = await axios.post(`${API_BASE_URL}/chat/create`, {
        description: text,
        one_shot: true,
        persona_id: selectedPersona || personas[0].id,
      });

      if (response.status === 201) {
        await createMessages(text, response.data.data.id);
      }
    } catch (error) {
      console.error("Error creating chat:", error);
    }
  }

  async function createMessages(text: string, passedChatId?: string) {
    const currentChatId = chatId || passedChatId;
    if (!currentChatId) {
      console.error("No chat ID available");
      return;
    }

    const userMessage: ChatMessage = {
      id: uuidv4(),
      message: text,
      chat_session_id: currentChatId,
      message_type: "user",
      time_sent: new Date().toISOString(),
    };

    try {
      const response = await fetch(`${API_BASE_URL}/chat/send-message`, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
          Authorization: `Bearer ${JSON.parse(localStorage.getItem("access_token") as string)}`,
        },
        body: JSON.stringify({ message: text, chat_session_id: currentChatId }),
      });

      if (!response.ok || !response.body) {
        throw new Error(response.statusText);
      }

      addMessage(userMessage);
      textInputRef.current!.value = "";

      const reader = response.body.getReader();
      const decoder = new TextDecoder("utf-8");
      let messageContent = "";

      while (true) {
        const { done, value } = await reader.read();
        if (done) break;
        const chunk = decoder.decode(value, { stream: true });
        messageContent += chunk;
        processStreamedChunk(chunk, currentChatId);
      }
    } catch (error) {
      console.error("Error creating messages:", error);
    }
  }

  function processStreamedChunk(chunk: string, chatSessionId: string) {
    const streams = chunk.split("\n");

    streams.forEach(stream => {
      if (stream) {
        const [event, data] = stream.split("data:");
        if (event === "event:document") {
          const doc = JSON.parse(data);
          updateDocumentMessage(doc.Document);
        } else if (event === "event:message") {
          const response = JSON.parse(data);
          addMessage({ ...response.Message, message: "" });
          setNewMessage(response.Message);
        } else if (event === "event:error") {
          const response = JSON.parse(data);
          router.navigate(`/chat/${chatSessionId}`);
          toast.error(response.Message.error);
        }
      }
    });
  }

  function updateDocumentMessage(document: any) {
    setMessages(prev => {
      const messageIndex = prev.findIndex(
        message => message.id === document.message_id
      );

      if (messageIndex !== -1) {
        const updatedMessages = [...prev];
        updatedMessages[messageIndex] = {
          ...updatedMessages[messageIndex],
          citations: (updatedMessages[messageIndex].citations || []).concat(document),
        };

        return updatedMessages;
      }

      return prev;
    });
  }

  function updateMessageContent() {
    setMessages(prevMessages =>
      prevMessages.map(prevMessage =>
        prevMessage.id === newMessage!.id
          ? { ...prevMessage, message: newMessage!.message.substring(0, newMessage!.message.length) }
          : prevMessage
      )
    );
  }

  async function handleSubmit() {
    if (textInputRef.current?.value) {
      setIsDeactivateSendingButton(true);
      chatId ? await createMessages(textInputRef.current.value) : await createChat(textInputRef.current.value);
      setIsDeactivateSendingButton(false);
    }
  }

  function chunkArray(array: any[], size: number) {
    const chunkedArray = [];
    for (let i = 0; i < array.length; i += size) {
      chunkedArray.push(array.slice(i, i + size));
    }
    return chunkedArray;
  }

  return (
    <div className="flex h-screen">
      <div className="flex flex-grow flex-col m-5 w-4/6">
        {messages.length === 0 ? (
          <PersonaSelection
            chunkArray={chunkArray}
            personas={personas}
            selectedPersona={selectedPersona}
            setSelectedPersona={setSelectedPersona}
          />
        ) : (
          <MessagesList messages={messages} newMessage={newMessage} />
        )}
        <ChatInput
          onSubmit={handleSubmit}
          textInputRef={textInputRef}
          isDeactivateSendingButton={isDeactivateSendingButton}
        />
      </div>
      <div className="hidden lg:my-8 lg:w-2/5 xl:w-3/12 lg:flex lg:flex-col lg:bg-white lg:rounded-md lg:rounded-l-none lg:overflow-x-hidden lg:no-scrollbar">
        <RetrievedKnowledge withHeader messages={messages} />
      </div>
    </div>
  );
}
