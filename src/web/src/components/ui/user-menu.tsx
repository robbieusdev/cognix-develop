import { ChevronUp, ChevronDown } from "lucide-react";
import { memo } from "react";
import { UserAccordion } from "./user-accordion";

interface Props {
  isUserMenuOpen: boolean;
  onToggle: () => void;
  firstName?: string;
  lastName?: string;
}

export const UserMenu = memo(
  ({ isUserMenuOpen, onToggle, firstName, lastName }: Props) => {
    const initials = `${firstName?.charAt(0) ?? ""}${
      lastName?.charAt(0) ?? ""
    }`;
    return (
      <div
        className="flex items-center justify-between p-4 text-white cursor-pointer shadow-lg bg-dark-gray"
        onClick={onToggle}
      >
        <div className="flex items-center space-x-3">
          <div className="flex items-center justify-center w-8 h-8 bg-accent rounded-md">
            <span className="text-sm text-ellipsis overflow-hidden">{initials}</span>
          </div>
          <span className="text-md font-medium">
            {firstName} {lastName}
          </span>
        </div>
        <div className="flex items-center  lg:ml-1">
          {isUserMenuOpen ? (
            <ChevronUp className="w-5 h-5" />
          ) : (
            <ChevronDown className="w-5 h-5" />
          )}
        </div>
        {isUserMenuOpen && (
          <div className="absolute bottom-9 w-2/5 lg:w-4/5">
            <UserAccordion />
          </div>
        )}
      </div>
    );
  }
);
