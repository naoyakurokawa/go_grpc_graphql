import { ApolloClient, InMemoryCache, createHttpLink } from "@apollo/client";

const getGraphQLEndpoint = () => {
  if (typeof window === "undefined") {
    return (
      process.env.BFF_GRAPHQL_ENDPOINT ??
      process.env.NEXT_PUBLIC_GRAPHQL_ENDPOINT ??
      "http://bff:8080/query"
    );
  }

  return (
    process.env.NEXT_PUBLIC_GRAPHQL_ENDPOINT ?? "http://localhost:8080/query"
  );
};

const httpLink = createHttpLink({
  uri: getGraphQLEndpoint(),
  credentials: "include",
});

const client = new ApolloClient({
  link: httpLink,
  cache: new InMemoryCache(),
});

export default client;
