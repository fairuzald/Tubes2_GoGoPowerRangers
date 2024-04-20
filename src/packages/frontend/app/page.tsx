"use client";
import { Button } from "@/components/ui/button";
import CardGrid from "@/components/individual-paths";
import InputQuery from "@/components/input-query";
import { useQueryContext } from "@/components/query-provider";
import type { PathInfo } from "@/types/result";
import Image from "next/image";
import { useState } from "react";
import toast from "react-hot-toast";

// Cache to store URL info results
const infoCache = new Map<string, PathInfo>();

const fetchInfoUrl = async (url: string) => {
  if (!url) throw new Error("No URL provided");

  // Check cache before making a request to the server
  if (infoCache.has(url)) {
    return infoCache.get(url)!;
  }

  try {
    const response = await fetch("/api/url-info", {
      headers: {
        "Content-Type": "application/json",
      },
      method: "POST",
      body: JSON.stringify({ url }),
    });

    const data = await response.json();
    if (!response.ok) {
      throw new Error(data.message || "Something went wrong");
    }

    const info = data.data;
    infoCache.set(url, info);

    return info;
  } catch (err) {
    console.error(err);
    const errMsg = err instanceof Error ? err.message : "Something went wrong";
    toast.error(errMsg);
    throw err; // Rethrow the error to propagate it further if needed
  }
};

export default function Home() {
  const { state, dispatch } = useQueryContext();
  const [loading, setLoading] = useState(false);

  const onSubmit = async () => {

    // Validation check
    if (!state.selectedSource || !state.selectedDestination) {
      toast.error("Please select source and destination from the select input");
      return;
    }
    setLoading(true);
    const loadingToast = toast.loading("Finding path...");
    try {
      const url = state.method === "ids" ? "/api/ids" : "/api/bfs";
      const response = await fetch(url, {
        cache: "no-cache",
        headers: {
          "Content-Type": "application/json",
        },
        method: "POST",
        body: JSON.stringify({
          source: state.selectedSource,
          destination: state.selectedDestination,
        }),
      });

      const data = await response.json();
      if (!response.ok) {
        throw new Error(data.message || "Something went wrong");
      }

      const result = data.paths as string[][];

      // Using Set to store unique URLs
      const uniquePaths = new Set<string>();
      result.forEach((path) => path.forEach((url) => uniquePaths.add(url)));

      // Using object for the final result with URL as key
      const resultWithInfo: Record<string, PathInfo[]> = {};

      // Mapping over the result to fetch info for each URL
      for (const path of result) {
        for (const url of path) {
          if (!resultWithInfo[url]) {
            const pathInfo = await fetchInfoUrl(url);
            if (pathInfo) {
              resultWithInfo[url] = [pathInfo];
            }
          }
        }
      }

      const filteredResult = Object.values(resultWithInfo).filter(
        (path) => path.length > 0
      );
      dispatch({ type: "SET_RESULT", payload: filteredResult });
    } catch (err) {
      console.error(err);
      const errMsg = err instanceof Error ? err.message : "Something went wrong";
      toast.error(errMsg);
    } finally {
      toast.dismiss(loadingToast);
      setLoading(false); // Reset loading state after fetching data
    }
  };

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
        <Button
          size={"lg"}
          className="text-2xl sm:text-3xl bg-yellow-primary hover:bg-yellow-hover transition ease-in-out delay-150 hover:scale-102 duration-300"
          onClick={onSubmit}
          disabled={loading}
        >
          {loading ? "Loading..." : "Go!"}
        </Button>
        <CardGrid />
      </div>
    </main>
  );
}
