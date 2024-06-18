import React, { Dispatch, memo, SetStateAction } from "react";
import { Card } from "../ui/card";
import { Persona } from "@/models/settings";

interface Props {
  chunkArray: (array: any[], size: number) => any[];
  personas: Persona[];
  selectedPersona: string;
  setSelectedPersona: Dispatch<SetStateAction<string>>;
}

export const CardGrid = memo(
  ({ chunkArray, personas, selectedPersona, setSelectedPersona }: Props) => {
    return (
      <div className="max-w-3xl flex justify-center items-center container mx-auto p-4">
        <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-2 gap-4">
          {chunkArray(personas, 2).map((chunk, chunkIndex) => (
            <React.Fragment key={chunkIndex}>
              {chunk.map((persona: Persona) => (
                <Card
                  key={persona.id}
                  title={persona.name}
                  text={persona.description}
                  selected={selectedPersona === persona.id}
                  onClick={() => setSelectedPersona(persona.id)}
                />
              ))}
            </React.Fragment>
          ))}
        </div>
      </div>
    );
  }
);
