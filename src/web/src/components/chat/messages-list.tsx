import React, { useEffect } from "react";
import { ChatMessage } from "@/models/chat";
import { dataConverter } from "@/lib/utils";
import MessageCard from "./components/message-card";

interface MessagesListProps {
  messages: ChatMessage[];
  newMessage: ChatMessage | null | undefined;
}

const MessagesList: React.FC<MessagesListProps> = ({
  messages,
  newMessage,
}) => {
  useEffect(() => {
    let index = 0;
    const intervalId = setInterval(() => {
      if (newMessage && newMessage.message) {
        messages = messages?.map((message) =>
          message.id === newMessage.id
            ? { ...message, message: newMessage.message.substr(0, index + 1) }
            : message
        );
        index++;
        if (index >= newMessage.message.length) {
          return () => clearInterval(intervalId);
        }
      }
    }, 25);
  }, [newMessage, messages]);

  return (
    <div className="flex flex-col flex-grow lg:mx-10 md:mx-10 overflow-y-scroll no-scrollbar">
      <div className="flex flex-grow items-start lg:my-4 my-10">
        <hr className="my-2 mr-4 flex-grow border-t border-gray-300" />
        <div className="text-muted-foreground text-sm font-thin">
          {dataConverter(messages[0]?.time_sent)}
        </div>
        <hr className="my-2 ml-4 flex-grow border-t border-gray-300" />
      </div>
      {messages.map((message, index) => (
        <MessageCard
          key={index}
          id={message.id}
          sender={message.message_type === "user" ? "You" : "AI Chat"}
          isResponse={message.message_type !== "user"}
          message={message.message ?? message.error}
          timestamp={message.time_sent}
          citations={message.citations}
          feedback={message.feedback}
        />
      ))}
    </div>
  );
};

export default MessagesList;
