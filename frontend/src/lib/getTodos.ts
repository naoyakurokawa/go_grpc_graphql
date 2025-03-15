import client from "@/client/client";
import { gql } from "@apollo/client";

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

export async function getTasks() {
  const { data } = await client.query({ query: GET_TASKS });
  return data.tasks;
}