import * as React from "react";
import { cn } from "@/lib/utils";

export interface TextAreaProps
  extends React.TextareaHTMLAttributes<HTMLTextAreaElement> { // Change input attributes to textarea attributes
  fieldState?: any;
}

const TextArea = React.forwardRef<HTMLTextAreaElement, TextAreaProps>(
  ({ className, fieldState, ...props }, ref) => {
    return (
      <textarea
        className={cn(
          "resize-none min-h-[90px] w-full rounded-md px-3 py-2 text-md shadow-lg placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-1 disabled:cursor-not-allowed disabled:opacity-50",
          className
        )}
        ref={ref}
        {...props}
      />
    );
  }
);

TextArea.displayName = "Textarea";

export { TextArea };
