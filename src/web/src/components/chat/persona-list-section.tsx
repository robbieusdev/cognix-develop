import React, { memo, Dispatch, SetStateAction } from "react";
import { Persona } from "@/models/settings";
import { CardGrid } from "./card-grid";

interface PersonaSelectionProps {
  personas: Persona[];
  selectedPersona: string;
  setSelectedPersona: Dispatch<SetStateAction<string>>;
  chunkArray: (array: any[], size: number) => any[];
}

export const PersonaSelection: React.FC<PersonaSelectionProps> = memo(
  ({ personas, selectedPersona, setSelectedPersona, chunkArray }) => {
    return (
      <div className="flex flex-col flex-grow overflow-x-hidden no-scrollbar">
        <div className="flex items-center justify-center pt-8">
          <span className="text-4xl font-bold">
            Which assistant do you want
          </span>
        </div>
        <div className="flex items-center justify-center pt-1">
          <span className="text-4xl font-bold">to chat with today?</span>
        </div>
        <div className="flex items-center justify-center pt-8">
          <span className="font-thin text-base text-muted">
            Or ask a question immediately to use the CogniX assistant
          </span>
        </div>
        <div className="pt-10 pb-2">
          {personas && (
            <CardGrid
              chunkArray={chunkArray}
              personas={personas}
              selectedPersona={selectedPersona}
              setSelectedPersona={setSelectedPersona}
            />
          )}
        </div>
      </div>
    );
  }
);
