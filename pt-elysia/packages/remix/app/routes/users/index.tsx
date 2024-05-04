import { treaty } from "@elysiajs/eden";

import type { App } from "@eventpuffin/server";

const app = treaty<App>("http://localhost:3000");

export default function UsersPage() {
  return <></>;
}
