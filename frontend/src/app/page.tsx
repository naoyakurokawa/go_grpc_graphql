"use client";

import { FormEvent, useMemo, useState } from "react";
import { useMutation, useQuery } from "@apollo/client";
import {
  CREATE_TASK,
  DELETE_TASK,
  GET_CATEGORIES,
  GET_TASKS,
  UPDATE_TASK,
  Category,
  Task,
} from "../lib/getTodos";

type GetTasksResponse = {
  tasks: Task[];
};

type GetCategoriesResponse = {
  categories: Category[];
};

export default function Home() {
  const [selectedCategoryId, setSelectedCategoryId] = useState<string>("");

  const categoryFilterValue = selectedCategoryId
    ? Number(selectedCategoryId)
    : null;

  const taskVariables = useMemo(
    () => ({ categoryId: categoryFilterValue }),
    [categoryFilterValue],
  );

  const { data, loading, error } = useQuery<GetTasksResponse>(GET_TASKS, {
    variables: taskVariables,
    fetchPolicy: "cache-and-network",
  });

  const {
    data: categoriesData,
    loading: categoriesLoading,
    error: categoriesError,
  } = useQuery<GetCategoriesResponse>(GET_CATEGORIES);

  const categories = categoriesData?.categories ?? [];
  const categoryNameMap = useMemo(() => {
    const map = new Map<number, string>();
    categories.forEach((category) => {
      map.set(category.id, category.name);
    });
    return map;
  }, [categories]);

  const taskRefetchQueries = useMemo(
    () => [{ query: GET_TASKS, variables: taskVariables }],
    [taskVariables],
  );

  const [createTaskMutation, { loading: creating }] = useMutation(
    CREATE_TASK,
    {
      refetchQueries: taskRefetchQueries,
      awaitRefetchQueries: true,
    },
  );

  const [updateTaskMutation] = useMutation(UPDATE_TASK, {
    refetchQueries: taskRefetchQueries,
    awaitRefetchQueries: true,
  });

  const [deleteTaskMutation] = useMutation(DELETE_TASK, {
    refetchQueries: taskRefetchQueries,
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
    const categoryIdRaw = (formData.get("category_id") as string | null)?.trim();

    if (!title || !note || !categoryIdRaw) {
      return;
    }

    const categoryId = Number(categoryIdRaw);
    if (!Number.isFinite(categoryId)) {
      return;
    }

    await createTaskMutation({
      variables: { input: { title, note, category_id: categoryId } },
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
      category_id?: number;
    } = { id };

    if (title) {
      input.title = title;
    }

    if (note) {
      input.note = note;
    }

    const categoryIdRaw = (formData.get("category_id") as string | null)?.trim();
    if (categoryIdRaw) {
      const categoryId = Number(categoryIdRaw);
      if (Number.isFinite(categoryId)) {
        input.category_id = categoryId;
      }
    }

    if (
      input.title === undefined &&
      input.note === undefined &&
      input.category_id === undefined
    ) {
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

      <section className="mb-6 flex flex-col gap-2">
        <label className="text-sm font-medium">
          カテゴリで絞り込み
          <select
            value={selectedCategoryId}
            onChange={(event) => setSelectedCategoryId(event.target.value)}
            className="mt-1 border rounded px-3 py-2 w-full max-w-xs"
            disabled={categoriesLoading}
          >
            <option value="">すべて</option>
            {categories.map((category) => (
              <option key={category.id} value={category.id}>
                {category.name}
              </option>
            ))}
          </select>
        </label>
        {categoriesError && (
          <p className="text-xs text-red-600">
            カテゴリの取得に失敗しました: {categoriesError.message}
          </p>
        )}
      </section>

      <section className="mb-6">
        <form onSubmit={handleCreate} className="flex flex-wrap gap-2">
          <input
            name="title"
            type="text"
            placeholder="Title"
            className="border rounded px-3 py-2 flex-1 min-w-[200px]"
            required
          />
          <textarea
            name="note"
            placeholder="Note"
            rows={3}
            className="border rounded px-3 py-2 flex-1 min-w-[200px] resize-y"
            required
          />
          <select
            name="category_id"
            defaultValue=""
            required
            className="border rounded px-3 py-2 flex-1 min-w-[200px]"
            disabled={categoriesLoading || categories.length === 0}
          >
            <option value="" disabled>
              カテゴリを選択
            </option>
            {categories.map((category) => (
              <option key={category.id} value={category.id}>
                {category.name}
              </option>
            ))}
          </select>
          <button
            type="submit"
            className="bg-blue-600 text-white px-4 py-2 rounded disabled:opacity-60 disabled:cursor-not-allowed"
            disabled={creating || categories.length === 0}
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
          const categoryName =
            task.category_id != null
              ? categoryNameMap.get(task.category_id)
              : undefined;

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
                  className={`text-sm text-gray-600 whitespace-pre-line ${
                    isCompleted ? "line-through" : ""
                  }`}
                >
                  {task.note}
                </p>
                <p className="text-xs text-gray-500">
                  カテゴリ: {categoryName ?? "未選択"}
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
                        <textarea
                          name="note"
                          defaultValue={task.note}
                          rows={3}
                          className="mt-1 border rounded px-2 py-1 w-full resize-y"
                        />
                      </label>
                      <label className="text-sm font-medium">
                        カテゴリ
                        <select
                          name="category_id"
                          defaultValue={
                            task.category_id != null
                              ? String(task.category_id)
                              : ""
                          }
                          className="mt-1 border rounded px-2 py-1 w-full"
                          disabled={categories.length === 0}
                        >
                          <option value="">
                            {task.category_id == null ? "未選択" : "変更しない"}
                          </option>
                          {categories.map((category) => (
                            <option key={category.id} value={category.id}>
                              {category.name}
                            </option>
                          ))}
                        </select>
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
