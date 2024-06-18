import { Persona } from "src/models/settings";

export interface MessageFeedback {
  id: string;
  chat_message_id: string;
  user_id: string;
  up_votes?: boolean;
}

export interface Document {
  id: string;
  link: string;
  document_id: string;
  content: string;
  updated_date: string;
}

export interface ChatMessage {
  chat_session_id: string;
  citations?: Document[];
  error?: string;
  id: string;
  latest_child_message?: number;
  message: string;
  message_type: string;
  parent_message?: number;
  rephrased_query?: string;
  time_sent: string;
  token_count?: number;
  feedback?: MessageFeedback;
}

export interface ChatSession {
  created_date: string;
  deleted_date?: string | null;
  description: string;
  id: string;
  messages: ChatMessage[];
  one_shot: boolean;
  persona: Persona;
  persona_id: number;
  user_id: string;
}
