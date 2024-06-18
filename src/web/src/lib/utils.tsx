import { Persona } from "@/models/settings";
import { type ClassValue, clsx } from "clsx";
import { twMerge } from "tailwind-merge";

import ConnectorsIcon from "@/assets/svgs/connectors.svg?react";
import FeedbackIcon from "@/assets/svgs/feedback.svg?react";
import LLMIcon from "@/assets/svgs/llm.svg?react";
import EmbeddingIcon from "@/assets/svgs/embedding.svg?react";
import UsersIcon from "@/assets/svgs/users.svg?react";
import ConfigIcon from "@/assets/svgs/config.svg?react";

export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs));
}

export function capitalize(str: string) {
  return str.charAt(0).toUpperCase() + str.slice(1);
}

export const virtuosoClassName = cn(
  "[&::-webkit-scrollbar]:bg-transparent [&::-webkit-scrollbar]:w-2",
  "[&::-webkit-scrollbar-thumb]:dark:bg-muted [&::-webkit-scrollbar-thumb]:bg-muted-foreground/40 [&::-webkit-scrollbar-thumb]:rounded"
);

export function dataConverter(dateString: string) {
  return new Date(dateString).toLocaleDateString("en-US", {
    month: "long",
    day: "numeric",
    year: "numeric",
  });
}

export function reassembleLLMData(data: Persona[]) {
  for (const record in data) {
    data[record].model_id = data[record].llm.model_id;
    data[record].endpoint = data[record].llm.endpoint;
  }
  return data;
}

interface IRoute {
  id: number;
  text: string;
  icon: React.ReactNode;
  link: string;
}

export const settings: IRoute[] = [
  {
    id: 1,
    text: "Connectors",
    icon: <ConnectorsIcon className="h-4 w-4" />,
    link: "/settings/connectors/existing-connectors",
  },
  {
    id: 2,
    text: "Feedback",
    icon: <FeedbackIcon className="h-4 w-4" />,
    link: "/settings/feedback",
  },
  {
    id: 3,
    text: "LLMs",
    icon: <LLMIcon className="h-4 w-4" />,
    link: "/settings/llms",
  },
  {
    id: 4,
    text: "Embeddings",
    icon: <EmbeddingIcon className="h-4 w-4" />,
    link: "/settings/embeddings",
  },
  {
    id: 5,
    text: "Users",
    icon: <UsersIcon className="h-4 w-4" />,
    link: "/settings/users",
  },
  {
    id: 6,
    text: "Config Map",
    icon: <ConfigIcon className="h-4 w-4" />,
    link: "/settings/config",
  },
];

// export function reassembleLLMInstance (data: Persona) {
//   data.model_id = data.llm.model_id
//   data.endpoint = data.llm.endpoint
//   data.task_prompt = data.prompt?.task_prompt
//   data.system_prompt = data.prompt?.system_prompt
//   data.url = data.prompt.url
//   return data
// }