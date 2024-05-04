import { Elysia, t } from "elysia";
import { swagger } from "@elysiajs/swagger";

import { userAPI } from "./domain/user";

const app = new Elysia()
  .group("/api", (app) => app.use(userAPI))
  .get("/", () => "Hello Elysia");

export type App = typeof app;

app.use(swagger());
app.listen(3000);

console.log(
  `🦊 Elysia is running at http://${app.server?.hostname}:${app.server?.port}`
);
