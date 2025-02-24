import client from "@/client/client";
import { gql } from "@apollo/client";

const GET_TODOS = gql`
  query GetTodos {
    todos {
      id
      text
    }
  }
`;

export async function getTodos() {
  const { data } = await client.query({ query: GET_TODOS });
  return data.todos;
}