import { Dispatch, SetStateAction, useContext, useState } from "react";
import { Button } from "@/components/ui/button";
import { ChevronDown, ChevronUp, PlusCircle } from "lucide-react";
import CognixSmall from "@/assets/svgs/cognix-sm.svg?react";
import SideBarClosedIcon from "@/assets/svgs/sidebar-closed-icon.svg?react";

import React from "react";
import { Link, NavLink } from "react-router-dom";
import { AuthContext } from "@/context/AuthContext";
import { router } from "@/main";
import { settings } from "@/lib/utils";
import { Sidebar } from "./sidebar";

export interface SideBarProps {
  isSideBarOpen: boolean;
  setIsSideBarOpen: Dispatch<SetStateAction<boolean>>;
  isUserMenuOpen: boolean;
  onToggle: () => void;
}

const Navbar: React.FC<SideBarProps> = ({
  isSideBarOpen,
  setIsSideBarOpen,
  isUserMenuOpen,
  onToggle,
}) => {
  const [isHistoryOpen, setIsHistoryOpen] = useState(false);
  const [isSettingsOpen, setIsSettingsOpen] = useState(false);
  const { chats } = useContext(AuthContext);

  if (isSideBarOpen) {
    return (
      <Sidebar
        isSideBarOpen={isSideBarOpen}
        setIsSideBarOpen={setIsSideBarOpen}
        isUserMenuOpen={isUserMenuOpen}
        onToggle={onToggle}
      />
    );
  }

  return (
    <div className="flex flex-col ml-3 mr-2 space-y-5">
      <div className="space-y-6">
        <div className="flex pt-9 pl-3">
          <SideBarClosedIcon
            className="cursor-pointer"
            onClick={() => {
              setIsSideBarOpen(!isSideBarOpen);
            }}
          />
        </div>
        <div className="ml-1.5">
          <Link to={"/"}>
            <CognixSmall className="h-9 w-9" />
          </Link>
        </div>
        <div>
          <Button
            variant="outline"
            size="icon"
            className="ml-1.5 bg-primary h-9 w-9"
            type="button"
            onClick={() => {
              router.navigate("/");
            }}
          >
            <PlusCircle className="h-4 w-4" />
          </Button>
        </div>
        <div
          className="flex items-center cursor-pointer"
          onClick={() => {
            setIsHistoryOpen(!isHistoryOpen);
          }}
        >
          {isHistoryOpen ? (
            <ChevronUp className="ml-4 h-4 w-4" />
          ) : (
            <ChevronDown className="ml-4 h-4 w-4" />
          )}
        </div>
      </div>
      {isHistoryOpen && (
        <div className="flex flex-col ml-1 space-y-3 text-2sm font-thin text-muted">
          {chats.slice(0, 4).map((chat) => (
            <NavLink
              key={chat.id}
              to={`/chat/${chat.id}`}
              className="flex flex-row items-center"
            >
              <span className="text-ellipsis overflow-hidden">{chat.description.substring(0, 6)}</span>
            </NavLink>
          ))}
        </div>
      )}
      <div
        className="flex cursor-pointer pt-3"
        onClick={() => {
          setIsSettingsOpen(!isSettingsOpen);
        }}
      >
        {isSettingsOpen ? (
          <ChevronUp className="ml-4 h-4 w-4" />
        ) : (
          <ChevronDown className="ml-4 h-4 w-4" />
        )}
      </div>
      {isSettingsOpen && (
        <div className="flex flex-col ml-4 space-y-4 text-muted">
          {settings.map((setting) => (
            <NavLink
              key={setting.id}
              to={setting.link}
              className={({ isActive }) => (isActive ? "active" : "inactive")}
            >
              <div className="h-6 w-4 mr-6 flex-shrink-0 fill-foreground/70 group-[.is-active]:fill-accent/95">
                {setting.icon}
              </div>
            </NavLink>
          ))}
        </div>
      )}
    </div>
  );
};

export { Navbar };
