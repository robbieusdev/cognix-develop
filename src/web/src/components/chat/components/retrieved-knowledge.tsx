import { memo } from "react";
import DocumentCard from "./document-card";
import FileIcon from "@/assets/svgs/file-icon.svg?react";
import { ChatMessage } from "@/models/chat";

interface Props {
  messages: ChatMessage[];
  withHeader?: boolean;
}

export const RetrievedKnowledge = memo(({ messages, withHeader }: Props) => {

  return (
    <div className="content-start space-x-2 pl-4">
      {withHeader && (
        <div className="flex content-start space-x-2 pt-5 pl-3">
          <FileIcon />
          <span className="font-bold">Retrieved Knowledge</span>
        </div>
      )}
      {messages && messages.length ? (
        <div>
          {messages?.map((message) =>
            message.citations?.map((citation, index) => (
              <DocumentCard
                key={index}
                id={citation.id}
                link={citation.link}
                content={citation.content}
                document_id={citation.document_id}
                date={citation.updated_date}
              />
            ))
          )}
        </div>
      ) : (
        <div>
          <div className="flex pt-5">
            <span className="font-thin text-sm text-muted">
              When you run ask a question, the
            </span>
          </div>
          <div className="flex pt-1">
            <span className="font-thin text-sm text-muted">
              retrieved knowledge will show up here
            </span>
          </div>
        </div>
      )}
    </div>
  );
});
