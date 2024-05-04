import { boolean } from "drizzle-orm/mysql-core";
import { integer, sqliteTable, text } from "drizzle-orm/sqlite-core";
import { createInsertSchema, createSelectSchema } from "drizzle-typebox";

export const userTable = sqliteTable("users", {
  id: integer("id").notNull().primaryKey({ autoIncrement: true }),
  name: text("name").notNull(),
  likesPizza: integer("likes_pizza", { mode: "boolean" }).notNull(),
});

export const insertUserSchema = createInsertSchema(userTable);
export const selectUserSchema = createSelectSchema(userTable);
