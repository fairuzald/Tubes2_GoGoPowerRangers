"use client";
import { useQueryContext } from "@/components/query-provider";
import { Button } from "@/components/ui/button";
import { Switch } from "@/components/ui/switch";
import type { PathInfo } from "@/types/result";
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

const SwitchAPIReq = () => {
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
      const url = !state.isBFS ? "/api/ids" : "/api/bfs";
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
      const uniquePathsWithInfo: Record<string, PathInfo> = {};

      // Fetch info for each unique URL
      for (const url of uniquePaths) {
        const info = await fetchInfoUrl(url);
        if (info) {
          uniquePathsWithInfo[url] = info;
        }
      }


      let dictionary: { [key: string]: { source: string; targets: Set<string> } } = {};

      // Map the result with the info using linkNodes
      let resultsWithInfo = [];
      for (let i = 0; i < result.length; i++) {
        let arr = [];
        for (let j = 0; j < result[i].length; j++) {
          const url = result[i][j];
          const info = uniquePathsWithInfo[url];
          arr.push(info);

          // Update dictionary with index information
          if (!dictionary[url]) {
            dictionary[url] = { source: url, targets: new Set<string>() };
          }
          dictionary[url].targets.add(result[i][j + 1] || "");
        }
        resultsWithInfo.push(arr);
      }


      dispatch({ type: "SET_RESULT", payload: resultsWithInfo });
      dispatch({ type: "SET_NODES", payload: new Set(uniquePaths) });
      dispatch({ type: "SET_LINK_NODES", payload: dictionary });
    } catch (err) {
      console.error(err);
      const errMsg =
        err instanceof Error ? err.message : "Something went wrong";
      toast.error(errMsg);
    } finally {
      toast.dismiss(loadingToast);
      setLoading(false); // Reset loading state after fetching data
    }
  };

  return (
    <section className="flex flex-col gap-5">
      <div className="flex items-center justify-center gap-4">
        <span className="font-montserrat text-[21px] font-semibold">IDS</span>
        <Switch
          className="bg-white"
          checked={state.isBFS}
          onCheckedChange={(checked) => dispatch({ type: "SET_ISBFS", payload: checked })}
        />
        <span className="font-montserrat text-[21px] font-semibold">
          BFS
        </span>
      </div>

      {/* Submit button */}
      <Button
        size={"lg"}
        className="text-2xl sm:text-3xl bg-yellow-primary hover:bg-yellow-hover transition ease-in-out delay-150 hover:scale-102 hover:-translate-y-1 duration-300"
        onClick={onSubmit}
        disabled={loading}
      >
        {loading ? "Loading..." : "Go!"}
      </Button>
    </section>
  );
};

export default SwitchAPIReq;
