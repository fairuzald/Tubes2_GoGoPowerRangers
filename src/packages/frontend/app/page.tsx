"use client";
import InputQuery from "@/components/input-query";
import { useQueryContext } from "@/components/query-provider";

export default function Home() {
  const { state } = useQueryContext();
  console.log(state);

  return (
    <main className="flex min-h-screen flex-col items-center p-24 w-full">
      <InputQuery />
    </main>
  );
}
