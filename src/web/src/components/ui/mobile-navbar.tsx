import { Dispatch, memo, SetStateAction, useState } from "react";
import Menu from "../../assets/svgs/menu.svg?react";
import Document from "../../assets/svgs/document.svg?react";
import { RetrievedKnowledgeDialog } from "../dialogs/RetrievedKnowledgeDialog";
import { MobileSidebar } from "./mobile-sidebar";

export interface SideBarProps {
  isSideBarOpen: boolean;
  setIsSideBarOpen: Dispatch<SetStateAction<boolean>>;
  isUserMenuOpen: boolean;
  onToggle: () => void;
}

export const MobileNavBar = memo(
  ({ isSideBarOpen, setIsSideBarOpen,   isUserMenuOpen,
    onToggle, }: SideBarProps) => {
    const [openModal, setOpenModal] = useState(false);

    if (isSideBarOpen) {
      return (
        <MobileSidebar
          setIsSideBarOpen={setIsSideBarOpen}
          isSideBarOpen={isSideBarOpen}
          isUserMenuOpen={isUserMenuOpen}
          onToggle={onToggle}
        />
      );
    }

    if (openModal) {
      return <RetrievedKnowledgeDialog setOpenModal={setOpenModal} />;
    }

    return (
      <div
        className={`md:hidden w-full fixed bg-foreground text-white transition-all duration-300 ease-in-out`}
      >
        <div className="p-4 h-full flex justify-between bg-foreground z-60">
          <Menu
            onClick={() => setIsSideBarOpen(!isSideBarOpen)}
            height={16}
            width={32}
            color="white"
            className="cursor-pointer"
          />
          <Document
            onClick={() => setOpenModal(!openModal)}
            height={16}
            width={32}
            color="white"
          />
        </div>
      </div>
    );
  }
);
