import { ChatComponent } from "@/components/chat";
import { useDocumentTitle } from "@/lib/hooks/use-document-title";

export function Chat() {
  useDocumentTitle("Chat");
  return (
    <div>
      <ChatComponent />
    </div>
  );
}

export { Chat as Component };
