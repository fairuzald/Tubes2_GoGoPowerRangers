"use client";
import InputQuery from "@/components/input-query";
import { useQueryContext } from "@/components/query-provider";
import Hero from "@/components/hero";
import SearchButton from "@/components/search-button";
import CardGrid from "@/components/individual-paths";

export default function Home() {
  const { state } = useQueryContext();
  console.log(state);

  return (
    <main className="flex min-h-screen flex-col items-center px-24 py-12 w-full bg-[#14213d] gap-10">
      <Hero />
      <InputQuery />
      <SearchButton />
      <CardGrid/>
    </main>
  );
}
