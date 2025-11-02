import client from "@/client/client";
import { gql } from "@apollo/client";

type Task = {
  id: number;
  title: string;
  note: string;
  completed: number;
  created_at: string;
  updated_at: string;
};

type GetTasksQuery = {
  tasks: Task[];
};

const GET_TASKS = gql`
  query GetTasks {
    tasks {
      id
      title
      note
      completed
      created_at
      updated_at
    }
  }
`;

export async function getTasks(): Promise<Task[]> {
  const { data } = await client.query<GetTasksQuery>({ query: GET_TASKS });
  return data.tasks;
}
