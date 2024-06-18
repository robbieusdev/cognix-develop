export type LLMSchema = {
  name: string;
  model_id: string;
  url?: string;
  api_key: string;
  endpoint?: string;
  system_prompt?: string;
  task_prompt?: string;
  description?: string;
};
