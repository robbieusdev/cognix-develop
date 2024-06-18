import React, { memo } from "react";
import { Button } from "../ui/button";
import { Input } from "../ui/input";
import SendIcon from "@/assets/svgs/send-icon.svg?react";

interface ChatInputProps {
  onSubmit: () => Promise<void>;
  isDeactivateSendingButton: boolean;
  textInputRef: any;
}

export const ChatInput: React.FC<ChatInputProps> = memo(
  ({ onSubmit, isDeactivateSendingButton, textInputRef }: ChatInputProps) => {
    return (
      <div className="w-full flex flex-col justify-start lg:justify-center">
        <div className="flex flex-row items-start justify-center py-4 gap-4">
          <Input
            placeholder="Ask me anything..."
            className="lg:w-2/3 md:w-2/3"
            ref={textInputRef}
            onKeyDown={(event) => {
              if (event.key === "Enter" && !isDeactivateSendingButton) {
                onSubmit();
              }
            }}
          />
          <Button
            size="icon"
            variant="outline"
            className="w-12 h-12 bg-primary hover:bg-foreground"
            type="button"
            onClick={onSubmit}
            disabled={isDeactivateSendingButton}
          >
            <SendIcon className="size-5" />
          </Button>
        </div>
        <div className="flex items-center justify-center pb-5">
          <span className="text-xs font-thin text-muted text-center">
            CogniX can make mistakes. Consider checking critical information.
          </span>
        </div>
      </div>
    );
  }
);
