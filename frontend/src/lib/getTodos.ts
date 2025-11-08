import client from "@/client/client";
import { gql } from "@apollo/client";

export type Task = {
  id: number;
  title: string;
  note: string;
  category_id?: number | null;
  completed: number;
  created_at: string;
  updated_at: string;
};

export type Category = {
  id: number;
  name: string;
};

export type GetTasksQuery = {
  tasks: Task[];
};

export const GET_TASKS = gql`
  query GetTasks($categoryId: Uint64) {
    tasks(category_id: $categoryId) {
      id
      title
      note
      category_id
      completed
      created_at
      updated_at
    }
  }
`;

export async function getTasks(categoryId?: number): Promise<Task[]> {
  const { data } = await client.query<GetTasksQuery>({
    query: GET_TASKS,
    variables: { categoryId },
  });
  return data.tasks;
}

export type CreateTaskMutation = {
  createTask: Task;
};

export type CreateTaskInput = {
  title: string;
  note: string;
  category_id: number;
};

export const CREATE_TASK = gql`
  mutation CreateTask($input: NewTask!) {
    createTask(input: $input) {
      id
      title
      note
      category_id
      completed
      created_at
      updated_at
    }
  }
`;

export async function createTask(input: CreateTaskInput): Promise<Task> {
  const { data } = await client.mutate<CreateTaskMutation>({
    mutation: CREATE_TASK,
    variables: { input },
  });

  if (!data) {
    throw new Error("Failed to create task");
  }

  return data.createTask;
}

export type UpdateTaskMutation = {
  updateTask: Task;
};

export type UpdateTaskInput = {
  id: number;
  title?: string;
  note?: string;
  category_id?: number;
  completed?: number;
};

export const UPDATE_TASK = gql`
  mutation UpdateTask($input: UpdateTask!) {
    updateTask(input: $input) {
      id
      title
      note
      category_id
      completed
      created_at
      updated_at
    }
  }
`;

export async function updateTask(input: UpdateTaskInput): Promise<Task> {
  const { data } = await client.mutate<UpdateTaskMutation>({
    mutation: UPDATE_TASK,
    variables: { input },
  });

  if (!data) {
    throw new Error("Failed to update task");
  }

  return data.updateTask;
}

export type DeleteTaskMutation = {
  deleteTask: boolean;
};

export const DELETE_TASK = gql`
  mutation DeleteTask($id: Uint64!) {
    deleteTask(id: $id)
  }
`;

export type GetCategoriesQuery = {
  categories: Category[];
};

export const GET_CATEGORIES = gql`
  query GetCategories {
    categories {
      id
      name
    }
  }
`;

export async function deleteTask(id: number): Promise<boolean> {
  const { data } = await client.mutate<DeleteTaskMutation>({
    mutation: DELETE_TASK,
    variables: { id },
  });

  if (!data) {
    throw new Error("Failed to delete task");
  }

  return data.deleteTask;
}
