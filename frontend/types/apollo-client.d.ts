declare module "@apollo/client" {
  import type { DocumentNode } from "graphql";

  export interface OperationVariables {
    [key: string]: unknown;
  }

  export interface ApolloQueryResult<TData> {
    data: TData;
  }

  export class InMemoryCache {
    constructor(options?: Record<string, unknown>);
  }

  export class ApolloClient<TCacheShape = unknown> {
    constructor(options: { uri: string; cache: InMemoryCache });
    query<TData = unknown, TVariables extends OperationVariables = OperationVariables>(
      options: { query: DocumentNode; variables?: TVariables }
    ): Promise<ApolloQueryResult<TData>>;
  }

  export function gql(
    literals: TemplateStringsArray,
    ...placeholders: unknown[]
  ): DocumentNode;
}
