import Cognix from "@/assets/svgs/cognix.svg?react";
import { AuthContext } from "@/context/AuthContext";
import { settings } from "@/lib/utils";
import { router } from "@/main";
import { ChevronDown, ChevronUp, PlusCircle, X } from "lucide-react";
import { Dispatch, memo, SetStateAction, useContext, useState } from "react";
import { Button } from "./button";
import { Link, NavLink } from "react-router-dom";
import { UserMenu } from "@/components/ui/user-menu";

interface Props {
  isSideBarOpen: boolean;
  setIsSideBarOpen: Dispatch<SetStateAction<boolean>>;
  isUserMenuOpen: boolean;
  onToggle: () => void;
}

export const MobileSidebar = memo(
  ({ isSideBarOpen, setIsSideBarOpen, isUserMenuOpen, onToggle }: Props) => {
    const [isHistoryOpen, setIsHistoryOpen] = useState(false);
    const [isSettingsOpen, setIsSettingsOpen] = useState(false);
    const { firstName, lastName, chats } = useContext(AuthContext);

    return (
      <div className="flex flex-col justify-between ml-2 mr-2 h-full z-100">
        <div>
          <div className="space-y-9">
            <div className="flex items-center mt-8 space-x-3">
              <X
                width={20}
                height={20}
                fill="#111"
                className="cursor-pointer"
                onClick={() => {
                  setIsSideBarOpen(!isSideBarOpen);
                }}
              />
              <Link to={"/"}>
                <Cognix width={140} height={42} />
              </Link>
            </div>
            <div className="mb-4 space-y-5">
              <Button
                variant="outline"
                size="lg"
                className="shadow-none bg-primary w-5/6"
                type="button"
                onClick={() => {
                  router.navigate("/");
                }}
              >
                <PlusCircle className="mr-2 h-9" width={160} height={160} />
                New chat
              </Button>
            </div>
            <div
              className="flex items-center cursor-pointer"
              onClick={() => {
                setIsHistoryOpen(!isHistoryOpen);
              }}
            >
              {isHistoryOpen ? (
                <ChevronUp className="h-4 w-4" />
              ) : (
                <ChevronDown className="h-4 w-4" />
              )}
              <span className="ml-2">Chat history</span>
            </div>
          </div>
          {isHistoryOpen && (
            <div className="flex flex-col ml-3 pt-4 space-y-4 text-2sm font-thin text-muted">
              {chats.slice(0, 4).map((chat) => (
                <NavLink
                  key={chat.id}
                  to={`/chat/${chat.id}`}
                  className="flex flex-row items-center"
                >
              <span className="text-ellipsis overflow-hidden">{chat.description.substring(0, 20)}</span>
                </NavLink>
              ))}
            </div>
          )}
          <div
            className="flex items-center cursor-pointer pt-3"
            onClick={() => {
              setIsSettingsOpen(!isSettingsOpen);
            }}
          >
            {isSettingsOpen ? (
              <ChevronUp className="h-4 w-4" />
            ) : (
              <ChevronDown className="h-4 w-4" />
            )}
            <span className="ml-2">Settings</span>
          </div>
          {isSettingsOpen && (
            <div className="flex flex-col ml-3 pt-4 space-y-4 text-muted text-2sm">
              {settings.map((setting) => (
                <NavLink
                  key={setting.id}
                  to={setting.link}
                  className="flex flex-row items-center"
                >
                  <div className="mr-2">{setting.icon}</div>
                  <span className="truncate">{setting.text}</span>
                </NavLink>
              ))}
            </div>
          )}
        </div>
        {(firstName || lastName) && (
          <UserMenu
            firstName={firstName}
            lastName={lastName}
            isUserMenuOpen={isUserMenuOpen}
            onToggle={onToggle}
          />
        )}
      </div>
    );
  }
);
