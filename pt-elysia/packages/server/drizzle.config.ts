import type { Config } from "drizzle-kit";

export default {
  schema: "./internal/db/**/*.schema.ts",
  out: "./internal/db/drizzle/",
} satisfies Config;
