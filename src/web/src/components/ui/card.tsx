import React, { memo, MouseEventHandler } from "react";

interface CardProps {
  title: string;
  text: string;
  selected: boolean;
  onClick: MouseEventHandler<HTMLDivElement>;
}

export const Card: React.FC<CardProps> = memo(
  ({ title, text, selected, onClick }) => {
    return (
      <div
        className={`w-80 p-4 bg-white rounded-lg shadow-md cursor-pointer ${
          selected ? "border-2 border-primary" : ""
        }`}
        onClick={onClick}
      >
        <h2 className="text-lg font-semibold mb-2">{title}</h2>
        <p className="text-gray-600">{text}</p>
      </div>
    );
  }
);
