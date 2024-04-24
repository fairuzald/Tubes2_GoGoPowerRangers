import ForceGraph from "@/components/graph";
import CardGridResult from "@/components/individual-paths";
import InputQuery from "@/components/input-query";
import SwitchAPIReq from "@/components/switch-api-req";
import { Metadata } from "next";
import Image from "next/image";

export default function Home() {
  return (
    <main className="flex min-h-screen flex-col items-center px-24 py-14 w-full bg-[#141414] relative z-10">
      {/* Video background */}
      <video
        className="absolute inset-0 w-full h-full object-cover -z-10"
        autoPlay
        loop
        muted
      >
        <source src={"bg1.mp4"} type="video/mp4" />
      </video>
      <div className="mt-4 flex flex-col gap-10 items-center justify-center w-full">
        {/* Main title section */}
        <section className="container mx-auto">
          <Image
            src="/wikirace.png"
            alt="logo wikirace"
            width={400}
            height={200}
            className="mx-auto"
            priority
          />
          <h1 className="text-center text-3xl font-bold mt-6">
            Find the shortest paths from
          </h1>
        </section>

        {/* Input Query */}
        <InputQuery />

        {/* Submit button */}
        <SwitchAPIReq />
        
        {/* Grid Result Card */}
        <CardGridResult />

        {/* Bonus: Graph Visualization */}
        <ForceGraph/>
      </div>
    </main>
  );
}

export const metadata: Metadata = {
  title: "Wikirace | Go Go Power Rangers!",
  description: "Find the shortest paths from source to destination on Wikipedia using BFS and IDS algorithms.",
  generator: "Next.js",
  keywords: ["wikirace", "wikipedia", "bfs", "ids", "shortest path"],
  applicationName: "Wikirace",
  category: "Game",
};