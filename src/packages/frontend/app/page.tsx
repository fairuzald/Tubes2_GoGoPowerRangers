"use client";
import { Button } from "@/components/ui/button";
import CardGrid from "@/components/individual-paths";
import InputQuery from "@/components/input-query";
import { useQueryContext } from "@/components/query-provider";
import type { PathInfo } from "@/types/result";
import Image from "next/image";
import { useState } from "react";
import toast from "react-hot-toast";
import SwitchOption from "@/components/switch-option";
// import Graph from "@/components/graph";
import ForceGraph from "@/components/graph";


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
      <div className="mt-4 flex flex-col gap-10 items-center justify-center">
        {/* Main title section */}
        <section className="container mx-auto">
          <Image
            src="/wikirace.png"
            alt="logo wikirace"
            width={400}
            height={200}
            className="mx-auto"
          />
          <h1 className="text-center text-3xl font-bold mt-6">
            Find the shortest paths from
          </h1>
        </section>
        {/* Input Query */}
        <InputQuery />
        {/* Submit button */}
        <SwitchOption/>
        
        {/* <Button
          size={"lg"}
          className="text-2xl sm:text-3xl bg-yellow-primary hover:bg-yellow-hover transition ease-in-out delay-150 hover:scale-102 hover:-translate-y-1 duration-300"
          onClick={onSubmit}
          disabled={loading}
        >
          {loading ? "Loading..." : "Go!"}
        </Button> */}
        <CardGrid />
        <ForceGraph/>
        {/* <Graph/> */}

        {/* <Results/> */}
      </div>
    </main>
  );
}
