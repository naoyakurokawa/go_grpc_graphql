import { getTodos } from "../lib/getTodos";

export default async function Home() {
  const todos = await getTodos();

  return (
    <main className="p-6">
      <h1 className="text-2xl font-bold mb-4">Todo List</h1>
      <ul className="list-disc pl-5">
        {todos.map((todo: { id: string; text: string }) => (
          <li key={todo.id} className="mb-2">
            {todo.text}
          </li>
        ))}
      </ul>
    </main>
  );
}