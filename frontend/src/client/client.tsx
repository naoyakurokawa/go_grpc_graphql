import { ApolloClient, InMemoryCache } from "@apollo/client";

const client = new ApolloClient({
  uri: "http://bff:8080/query",
  cache: new InMemoryCache(),
});

export default client;