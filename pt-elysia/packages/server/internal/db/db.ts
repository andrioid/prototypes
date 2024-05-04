import { Elysia } from "elysia";
import Database from "bun:sqlite";
import { migrate } from "drizzle-orm/bun-sqlite/migrator";
import { drizzle } from "drizzle-orm/bun-sqlite";
import dzConfig from "../../drizzle.config";

const sqlite = new Database("app.db");
export const db = drizzle(sqlite);
migrate(db, { migrationsFolder: dzConfig.out });

// Not a great idea. Not available to sub routes
//export const dbPlugin = new Elysia().decorate("db", db);
