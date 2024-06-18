import React, { useEffect, useState } from "react";
import { formatDistanceToNow } from "date-fns";

export interface DocumentProps {
  id: string;
  link: string;
  content: string;
  document_id: string;
  className?: string;
  date: string;
}
const DocumentCard: React.FC<DocumentProps> = ({ link, content, date }) => {
  const [timeElapsed, setTimeElapsed] = useState("");

  useEffect(() => {
    const parsedDate = new Date(date);

    const updateTimeElapsed = () => {
      setTimeElapsed(formatDistanceToNow(parsedDate, { addSuffix: true }));
    };

    updateTimeElapsed();

    const intervalId = setInterval(updateTimeElapsed, 60000);

    return () => clearInterval(intervalId);
  }, [date]);

  return (
    <div className="w-full max-w-sm sm:max-w-md bg-white shadow-lg rounded-lg overflow-hidden my-4">
      <div className="bg-gray-100 px-3 py-4">
        <div className="flex flex-col">
          <div className="flex flex-row items-center justify-between flex-wrap">
            <p className="font-bold text-sm truncate">
              {link.length < 16 ? link : `${link.slice(0, 16)}...`}
            </p>
            <p className="text-gray-600 rounded-md p-1 bg-white text-sm">
              Updated {timeElapsed}
            </p>
          </div>
        </div>
      </div>
      <div className="p-2 ml-2 w-full">
        <div className="px-3 py-4 text-muted-foreground break-all">
          AI Chat interactions, the client could indicate to NATS that multiple
          replies should be allowed. This flexibility allows for more dynamic
          and scalable communication between clients and servers
          {content}
        </div>
      </div>
    </div>
  );
};

export default DocumentCard;
