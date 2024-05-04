import Elysia, { t } from "elysia";
import { insertUserSchema, userTable } from "./schema";
import { db } from "../../internal/db/db";
import { eq } from "drizzle-orm";

export const userAPI = new Elysia({ prefix: "/user" })
  .get("/", async ({ params }) => {
    const result = await db.select().from(userTable);
    return result;
  })
  .get(
    "/:id",
    async ({ params: { id } }) => {
      const result = await db
        .select()
        .from(userTable)
        .where(eq(userTable.id, id));
      return result;
    },
    {
      params: t.Object({
        id: t.Numeric(),
      }),
    }
  )
  .post(
    "/",
    async ({ body }) => {
      const res = await db
        .insert(userTable)
        .values(body)
        .returning({ id: userTable.id });
      if (res[0]) {
        return res[0].id;
      }
    },
    {
      body: insertUserSchema,
    }
  );
