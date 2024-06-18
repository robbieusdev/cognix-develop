import { z } from "zod";

export const formSchema = z.object({
  name: z.string(),
  model_id: z.string(),
  url: z.string().optional(),
  api_key: z.string(),
  endpoint: z.string().optional(),
  system_prompt: z.string().optional(),
  task_prompt: z.string().optional(),
  description: z.string().optional(),
});
