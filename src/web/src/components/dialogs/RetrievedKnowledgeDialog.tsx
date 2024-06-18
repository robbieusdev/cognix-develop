import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";
import { Dispatch, memo, SetStateAction } from "react";
import { RetrievedKnowledge } from "../chat/components/retrieved-knowledge";
import { X } from "lucide-react";
import Document from "@/assets/svgs/file-icon.svg?react";
import { useMessages } from "@/context/ChatContext";

interface Props {
  setOpenModal: Dispatch<SetStateAction<boolean>>;
}

export const RetrievedKnowledgeDialog = memo(({ setOpenModal }: Props) => {
  const { messages } = useMessages();

  return (
    <Dialog modal open={true} onOpenChange={() => {}}>
      <DialogContent className="sm:max-w-[765px] h-full block">
        <DialogHeader className="flex-row justify-between items-center h-5">
          <DialogTitle className="flex flex-row items-center justify-start gap-1">
            <Document width={16} height={16} />
            Retrieved Knowledge
          </DialogTitle>
          <X onClick={() => setOpenModal(false)}>Close</X>
        </DialogHeader>
        <div className="w-full h-full flex flex-col overflow-hidden px-0.5 bg-white">
          <RetrievedKnowledge messages={messages} />
        </div>
      </DialogContent>
    </Dialog>
  );
});
