import * as React from "react";
import { cn } from "@/lib/utils";

export interface InputProps
  extends React.InputHTMLAttributes<HTMLInputElement> {
  fieldState?: any;
}

const Input = React.forwardRef<HTMLInputElement, InputProps>(
  ({ className, fieldState, ...props }, ref) => {
    return (
      <input
        className={cn(
          "flex items-center min-h-[50px] w-full rounded-md px-3 py-2 text-sm shadow-lg placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-1 disabled:cursor-not-allowed disabled:opacity-50",
          className
        )}
        ref={ref}
        {...props}
      />
    );
  }
);

Input.displayName = "Textarea";

export { Input };
