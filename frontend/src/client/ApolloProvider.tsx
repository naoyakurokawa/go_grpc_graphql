"use client";

import { ApolloProvider } from "@apollo/client";
import client from "@/client/client";

type Props = {
  children: React.ReactNode;
};

export function ApolloProviderWrapper({ children }: Props) {
  return <ApolloProvider client={client}>{children}</ApolloProvider>;
}
