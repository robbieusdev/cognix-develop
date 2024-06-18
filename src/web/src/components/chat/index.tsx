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

export function ChatComponent() {
  const { messages, setMessages, addMessage } = useMessages();
  const [personas, setPersonas] = useState<Persona[]>([]);
  const [selectedPersona, setSelectedPersona] = useState<string>("");
  const [newMessage, setNewMessage] = useState<ChatMessage | null>();
  const [isDeactivateSendingButton, setIsDeactivateSendingButton] =
    useState(false);
  const textInputRef = useRef<HTMLInputElement>(null);

  const { chatId } = useParams<{
    chatId?: string;
  }>();

  async function onHandleSubmit() {
    if (textInputRef.current?.value) {
      setIsDeactivateSendingButton(true);
      chatId
        ? await createMessages(textInputRef.current?.value ?? "")
        : await createChat(textInputRef.current?.value ?? "");
      setIsDeactivateSendingButton(false);
    }
  }

  async function getMessages(): Promise<void> {
    await axios
      .get(`${import.meta.env.VITE_PLATFORM_API_CHAT_DETAIL_URL}/${chatId}`)
      .then(function (response) {
        if (response.status == 200) {
          setMessages(response.data.data.messages);
        } else {
          setMessages([]);
        }
      })
      .catch(function (error) {
        console.error("Error fetching messages:", error);
      });
  }

  async function createChat(text: string) {
    if (!selectedPersona) {
      if (!personas || personas.length === 0) {
        toast.error("You don't have any available Assistants yet.");
        return;
      }
    }
    await axios
      .post(import.meta.env.VITE_PLATFORM_API_CHAT_CREATE_URL, {
        description: text,
        one_shot: true,
        persona_id: selectedPersona || personas[0].id,
      })
      .then(async function (response) {
        if (response.status == 201) {
          await createMessages(text, response.data.data.id);
        } else {
          setPersonas([]);
        }
      });
  }

  async function createMessages(
    text: string,
    passedChatId?: string
  ): Promise<void> {
    try {
      const currentChatId = chatId ?? passedChatId;
      if (!currentChatId) {
        throw new Error("Chat ID is not provided.");
      }
      const userMessage: ChatMessage = {
        id: uuidv4(),
        message: text,
        chat_session_id: currentChatId,
        message_type: "user",
        time_sent: new Date().toISOString(),
      };
      const response = await fetch(
        `${import.meta.env.VITE_PLATFORM_API_CHAT_SEND_MESSAGE_URL}`,
        {
          method: "POST",
          body: JSON.stringify({
            message: text,
            chat_session_id: currentChatId,
          }),
          headers: {
            "Content-Type": "application/json",
            Authorization: `Bearer ${JSON.parse(
              String(localStorage.getItem("access_token"))
            )}`,
          },
        }
      );
      if (!response.ok) {
        throw new Error(`Failed to send message: ${response.statusText}`);
      }
      addMessage(userMessage);
      if (textInputRef?.current) {
        textInputRef.current.value = "";
      }
      const reader = response.body
        ?.pipeThrough(new TextDecoderStream())
        .getReader();

      if (!reader) {
        throw new Error("Failed to get reader from response body.");
      }

      let buffer = "";
      while (true) {
        const { value, done } = await reader.read();
        if (done) break;
        buffer += value;
        const events = buffer.split("\n\n");
        buffer = events.pop() || "";
        for (const event of events) {
          const [eventLine, dataLine] = event.split("\n");
          if (!eventLine || !dataLine) continue;
          const eventType = eventLine.split(":")[1]?.trim();
          const data = dataLine.split("data:")[1]?.trim();
          if (!eventType || !data) {
            console.warn("Invalid event or data:", { eventLine, dataLine });
            continue;
          }
          try {
            const parsedData = JSON.parse(data);
            if (eventType === "document") {
              const doc = parsedData;
              setMessages((prev) => {
                const messageIndex = prev.findIndex(
                  (message) => message.id === doc.Document.message_id
                );
                if (messageIndex !== -1) {
                  const updatedMessages = [...prev];
                  updatedMessages[messageIndex] = {
                    ...updatedMessages[messageIndex],
                    citations: (
                      updatedMessages[messageIndex].citations || []
                    ).concat(doc.Document),
                  };
                  return updatedMessages;
                }
                return prev;
              });
            } else if (eventType === "message") {
              setMessages((prev) => [
                ...(prev ?? []),
                { ...parsedData.Message, message: "" },
              ]);
              setNewMessage(parsedData.Message);
            } else if (eventType === "error") {
              router.navigate(`/chat/${currentChatId}`);
              toast.error(parsedData.Message.error);
            }
          } catch (error) {
            console.error("Error parsing JSON:", error, "Data:", data);
          }
        }
      }
    } catch (error) {
      console.error("Error in createMessages:", error);
    }
  }

  async function getPersonas() {
    try {
      const response = await axios.get(
        import.meta.env.VITE_PLATFORM_API_LLM_LIST_URL
      );
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

  function chunkArray(array: any[], size: number) {
    const chunkedArray = [];
    for (let i = 0; i < array.length; i += size) {
      chunkedArray.push(array.slice(i, i + size));
    }
    return chunkedArray;
  }

  useLayoutEffect(() => {
    if (chatId) {
      getMessages();
    } else {
      getPersonas();
      setMessages([]);
    }
  }, [chatId]);

  useEffect(() => {
    let index = 0;
    const intervalId = setInterval(() => {
      if (newMessage && newMessage.message) {
        setMessages((prevMessages) =>
          prevMessages?.map((prevMessage) =>
            prevMessage.id === newMessage.id
              ? {
                  ...prevMessage,
                  message: newMessage.message.substr(0, index + 1),
                }
              : prevMessage
          )
        );
        index++;
        if (index >= newMessage.message.length) {
          clearInterval(intervalId);
        }
      }
    }, 25);
    if (!chatId && newMessage) {
      router.navigate(`/chat/${newMessage.chat_session_id}`);
    }
    return () => {
      clearInterval(intervalId);
    };
  }, [newMessage]);
  
  return (
    <div className="flex h-screen">
      <div className="flex flex-grow flex-col m-5 w-4/6">
        {messages.length == 0 ? (
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
          onSubmit={onHandleSubmit}
          textInputRef={textInputRef}
          isDeactivateSendingButton={isDeactivateSendingButton}
        />
      </div>
      <div className="hidden lg:my-8 lg:w-2/5 xl:w-3/12 px-2 lg:flex lg:flex-col lg:bg-white lg:rounded-md lg:rounded-l-none lg:overflow-x-hidden lg:no-scrollbar">
        <RetrievedKnowledge withHeader messages={messages} />
      </div>
    </div>
  );
}
