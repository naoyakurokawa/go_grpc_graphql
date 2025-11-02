import type { JSX } from "react";

import { getTasks } from "../lib/getTodos";

export default async function Home(): Promise<JSX.Element> {
  const tasks = await getTasks();

  return (
    <main className="p-6">
      <h1 className="text-2xl font-bold mb-4">Todo List</h1>
      <ul className="list-disc pl-5">
        {tasks.map((task: { id: string; title: string, note: string }) => (
          <li key={task.id} className="mb-2">
            {task.title}: {task.note}
          </li>
        ))}
      </ul>
    </main>
  );
}
