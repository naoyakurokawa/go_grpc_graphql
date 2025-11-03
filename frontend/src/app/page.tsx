"use client";

import { FormEvent, useState } from "react";
import { useMutation, useQuery } from "@apollo/client";
import {
  CREATE_TASK,
  DELETE_TASK,
  GET_TASKS,
  UPDATE_TASK,
  Task,
} from "../lib/getTodos";

type GetTasksResponse = {
  tasks: Task[];
};

export default function Home() {
  const { data, loading, error } = useQuery<GetTasksResponse>(GET_TASKS, {
    fetchPolicy: "cache-and-network",
  });

  const [createTaskMutation, { loading: creating }] = useMutation(
    CREATE_TASK,
    {
      refetchQueries: [{ query: GET_TASKS }],
      awaitRefetchQueries: true,
    },
  );

  const [updateTaskMutation] = useMutation(UPDATE_TASK, {
    refetchQueries: [{ query: GET_TASKS }],
    awaitRefetchQueries: true,
  });

  const [deleteTaskMutation] = useMutation(DELETE_TASK, {
    refetchQueries: [{ query: GET_TASKS }],
    awaitRefetchQueries: true,
  });

  const [updatingId, setUpdatingId] = useState<number | null>(null);
  const [deletingId, setDeletingId] = useState<number | null>(null);
  const [togglingId, setTogglingId] = useState<number | null>(null);

  const handleCreate = async (event: FormEvent<HTMLFormElement>) => {
    event.preventDefault();
    const form = event.currentTarget;
    const formData = new FormData(form);
    const title = (formData.get("title") as string | null)?.trim();
    const note = (formData.get("note") as string | null)?.trim();

    if (!title || !note) {
      return;
    }

    await createTaskMutation({
      variables: { input: { title, note } },
    });

    form.reset();
  };

  const handleUpdate = async (event: FormEvent<HTMLFormElement>) => {
    event.preventDefault();
    const form = event.currentTarget;
    const formData = new FormData(form);
    const idRaw = formData.get("id");

    if (typeof idRaw !== "string" || !idRaw.trim()) {
      return;
    }

    const id = Number(idRaw);
    if (!Number.isFinite(id)) {
      return;
    }

    const title = (formData.get("title") as string | null)?.trim();
    const note = (formData.get("note") as string | null)?.trim();

    const input: {
      id: number;
      title?: string;
      note?: string;
    } = { id };

    if (title) {
      input.title = title;
    }

    if (note) {
      input.note = note;
    }

    if (input.title === undefined && input.note === undefined) {
      return;
    }

    setUpdatingId(id);
    try {
      await updateTaskMutation({
        variables: { input },
      });
    } finally {
      setUpdatingId(null);
      const details = form.closest("details");
      if (details instanceof HTMLDetailsElement) {
        details.open = false;
      }
    }
  };

  const handleDelete = async (id: number) => {
    setDeletingId(id);
    try {
      await deleteTaskMutation({
        variables: { id },
      });
    } finally {
      setDeletingId(null);
    }
  };

  const handleToggleCompletion = async (task: Task) => {
    setTogglingId(task.id);
    try {
      await updateTaskMutation({
        variables: {
          input: {
            id: task.id,
            completed: task.completed ? 0 : 1,
          },
        },
      });
    } finally {
      setTogglingId(null);
    }
  };

  const tasks = data?.tasks ?? [];

  return (
    <main className="p-6">
      <h1 className="text-2xl font-bold mb-4">Todo List</h1>

      <section className="mb-6">
        <form onSubmit={handleCreate} className="flex flex-wrap gap-2">
          <input
            name="title"
            type="text"
            placeholder="Title"
            className="border rounded px-3 py-2 flex-1 min-w-[200px]"
            required
          />
          <input
            name="note"
            type="text"
            placeholder="Note"
            className="border rounded px-3 py-2 flex-1 min-w-[200px]"
            required
          />
          <button
            type="submit"
            className="bg-blue-600 text-white px-4 py-2 rounded disabled:opacity-60 disabled:cursor-not-allowed"
            disabled={creating}
          >
            {creating ? "追加中..." : "追加"}
          </button>
        </form>
      </section>

      {loading && (
        <p className="text-sm text-gray-500">読み込み中...</p>
      )}

      {error && (
        <p className="text-sm text-red-600">
          データの取得に失敗しました: {error.message}
        </p>
      )}

      {!loading && tasks.length === 0 && (
        <p className="text-sm text-gray-500">タスクがありません。</p>
      )}

      <ul className="space-y-4 mt-4">
        {tasks.map((task) => {
          const isCompleted = Boolean(task.completed);

          return (
            <li
              key={task.id}
              className={`border rounded p-4 shadow-sm flex flex-col gap-3 ${isCompleted ? "bg-gray-200" : ""
                }`}
            >
              <div>
                <p
                  className={`font-semibold ${isCompleted ? "line-through text-gray-600" : ""
                    }`}
                >
                  {task.title}
                </p>
                <p
                  className={`text-sm text-gray-600 ${isCompleted ? "line-through" : ""
                    }`}
                >
                  {task.note}
                </p>
              </div>
              <div className="flex items-start gap-3">
                <details className="relative">
                  <summary className="list-none">
                    <span className="inline-block bg-yellow-500 text-white px-3 py-1 rounded cursor-pointer">
                      編集
                    </span>
                  </summary>
                  <div className="mt-2 border rounded p-3 bg-white shadow-lg">
                    <form
                      onSubmit={handleUpdate}
                      className="flex flex-col gap-2"
                    >
                      <input type="hidden" name="id" value={task.id} />
                      <label className="text-sm font-medium">
                        タイトル
                        <input
                          name="title"
                          defaultValue={task.title}
                          className="mt-1 border rounded px-2 py-1 w-full"
                        />
                      </label>
                      <label className="text-sm font-medium">
                        メモ
                        <input
                          name="note"
                          defaultValue={task.note}
                          className="mt-1 border rounded px-2 py-1 w-full"
                        />
                      </label>
                      <button
                        type="submit"
                        className="bg-yellow-500 text-white px-3 py-1 rounded self-end disabled:opacity-60 disabled:cursor-not-allowed"
                        disabled={updatingId === task.id}
                      >
                        {updatingId === task.id ? "保存中..." : "保存"}
                      </button>
                    </form>
                  </div>
                </details>
                <button
                  type="button"
                  onClick={() => handleToggleCompletion(task)}
                  className="bg-green-600 text-white px-3 py-1 rounded disabled:opacity-60 disabled:cursor-not-allowed"
                  disabled={togglingId === task.id}
                >
                  {togglingId === task.id
                    ? "更新中..."
                    : isCompleted
                      ? "再オープン"
                      : "完了"}
                </button>
                <button
                  type="button"
                  onClick={() => handleDelete(task.id)}
                  className="bg-red-600 text-white px-3 py-1 rounded disabled:opacity-60 disabled:cursor-not-allowed"
                  disabled={deletingId === task.id}
                >
                  {deletingId === task.id ? "削除中..." : "削除"}
                </button>
              </div>
            </li>
          );
        })}
      </ul>
    </main>
  );
}
