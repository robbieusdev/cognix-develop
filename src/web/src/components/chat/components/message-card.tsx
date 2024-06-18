import React, { useContext, useState } from "react";
import BotIcon from "@/assets/svgs/cognix-sm.svg?react";
import CopyIcon from "@/assets/svgs/copy-icon.svg?react";
import ThumbUpIcon from "@/assets/svgs/thumb-up.svg?react";
import ThumbDownIcon from "@/assets/svgs/thumn-down.svg?react";
import FileWhiteIcon from "@/assets/svgs/file-white-icon.svg?react";
import axios from "axios";
import { AuthContext } from "@/context/AuthContext";
import { Document, MessageFeedback } from "@/models/chat";

export interface MessageProps {
  id: string;
  sender?: string;
  message: string;
  timestamp: string;
  className?: string;
  citations?: Document[];
  isResponse?: boolean;
  feedback?: MessageFeedback;
}

const MessageCard: React.FC<MessageProps> = ({
  id,
  sender,
  message,
  citations,
  feedback,
  isResponse,
  className,
}) => {
  const { firstName, lastName } = useContext(AuthContext);

  const [feedbackValue, setFeedbackValue] = useState(feedback?.up_votes);

  async function feedbackMessage(vote: "upvote" | "downvote"): Promise<void> {
    await axios
      .post(import.meta.env.VITE_PLATFORM_API_CHAT_FEEDBACK_MESSAGE_URL, {
        id: id,
        vote: vote,
      })
      .then((response) => {
        if (response.status == 200) {
          setFeedbackValue(vote === "upvote");
        }
      });
  }

  return (
    <div className="flex flex-wrap w-full">
      <div className={`flex flex-col py-4 lg:w-10/12 w-12/12 ${className}`}>
        <div className="flex items-start mb-2">
          {isResponse ? (
            <BotIcon className="w-10 h-10" />
          ) : (
            <div className="w-10 h-10 rounded-full bg-accent flex items-center justify-center">
              <span className="text-md text-white">
                {firstName && `${firstName.charAt(0)}`}
              </span>
              <span className="text-md text-white">
                {lastName && `${lastName.charAt(0)}`}
              </span>
            </div>
          )}
          <div className="flex flex-grow items-start ml-2">
            <div className="text-sm font-bold">{sender}</div>
          </div>
        </div>
        <div className="ml-12">
          <div className="-mt-6 text-muted-foreground break-all">{message}</div>
          <div>
            {sender !== "You" && <div className="pt-2 font-bold">Sources:</div>}
            {citations?.map((citation) => (
              <div
                key={citation.id}
                className="inline-flex cursor-pointer items-center m-1 px-2 py-1 space-x-2 bg-main rounded-lg shadow-md"
              >
                <FileWhiteIcon className="w-4 h-4" />
                <span>{citation.link}</span>
              </div>
            ))}
          </div>
          {isResponse && (
            <div className="flex items-center mt-5 space-x-3 text-muted">
              <div
                onClick={() => {
                  navigator.clipboard.writeText(message);
                }}
              >
                <CopyIcon className="w-5 h-5 cursor-pointer" />
              </div>
              <div
                onClick={() => {
                  feedbackMessage("upvote");
                }}
              >
                <ThumbUpIcon
                  className="w-5 h-5 cursor-pointer"
                  fill={feedbackValue ? "#047C57" : ""}
                />
              </div>
              <div
                className="w-5 h-5"
                onClick={() => {
                  feedbackMessage("downvote");
                }}
              >
                <ThumbDownIcon
                  className="w-5 h-5 cursor-pointer"
                  fill={
                    feedbackValue != undefined && !feedbackValue
                      ? "#DB3A34"
                      : ""
                  }
                />
              </div>
            </div>
          )}
        </div>
      </div>
    </div>
  );
};

export default MessageCard;
