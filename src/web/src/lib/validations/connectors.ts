import { z } from "zod";

export const formSchema = z.object({
  name: z.string(),
  source: z.string(),
  connector_specific_config: z.string(),
  refresh_freq: z.string(),
  credential_id: z.string(),
});
