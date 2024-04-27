"use client";
import { useQueryContext } from "@/components/query-provider";
import { Button } from "@/components/ui/button";
import { Switch } from "@/components/ui/switch";
import { makeApiRequest } from "@/libs/helper";
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

  let datas = null;

  try {
    await makeApiRequest({
      isToast: false,
      method: "POST",
      endpoint: "/url-info",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({
        url,
      }),
      loadingMessage: "Fetching data...",
      successMessage: "Data fetched successfully!",
      onSuccess(data: any) {
        const dataInfo = data.data as PathInfo;
        infoCache.set(url, dataInfo as PathInfo);
        datas = dataInfo;
      },
    });

    return datas;
  } catch (err) {
    console.error(err);
    const errMsg = err instanceof Error ? err.message : "Something went wrong";
    toast.error(errMsg);
    throw err; // Rethrow the error to propagate it further if needed
  }
};
interface ApiResponse {
  message: string;
  paths: string[][];
  runtime: number;
  articleCount: number;
}

const SwitchAPIReq = () => {
  const { state, dispatch } = useQueryContext();
  const [loading, setLoading] = useState(false);
  const [saving, setSaving] = useState(false);

  const handleSave = async () => {
    if (state.result.length <= 0) {
      return;
    }

    setSaving(true);
    try {
      await makeApiRequest({
        endpoint: "/save",
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        loadingMessage: "",
        successMessage: "",
        body: JSON.stringify({
          source: state.selectedSource,
          destination: state.selectedDestination,
          paths: state.result,
        }),
        onSuccess: () => {
          toast.success("Data saved successfully!");
        },
      });
    } catch (error) {
      toast.error("Failed to save data!");
    } finally {
      setSaving(false);
    }
  };

  const onSubmit = async () => {
    // Validation check
    if (!state.selectedSource || !state.selectedDestination) {
      toast.error("Please select source and destination from the select input");
      return;
    }
    setLoading(true);
    try {
      const url = !state.isBFS ? "/ids" : "/bfs";
      await makeApiRequest({
        method: "POST",
        endpoint: url,
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({
          source: state.selectedSource,
          destination: state.selectedDestination,
        }),
        loadingMessage:
          "Finding sortest path using " +
          (state.isBFS ? "BFS" : "IDS") +
          " algorithm...",
        successMessage: "Process completed successfully!",
        onSuccess: async (data: ApiResponse) => {
          const result = data.paths as string[][];

          // Using Set to store unique URLs
          const uniquePaths = new Set<string>();
          result.forEach((path) => path.forEach((url) => uniquePaths.add(url)));

          // Using object for the final result with URL as key
          const uniquePathsWithInfo: Record<string, PathInfo> = {};

          // add depth
          const uniquePathsWithDepth: Record<string, number> = {};
          const toastId = toast.loading("Fetching data info...");
          for (const url of uniquePaths) {
            const info = await fetchInfoUrl(url);
            if (info) {
              uniquePathsWithInfo[url] = info;
            }
          }
          toast.dismiss(toastId);

          for (let i = 0; i < result.length; i++) {
            for (let j = 0; j < result[i].length; j++) {
              const url = result[i][j];
              // Mapping the title
              if (!uniquePathsWithDepth[url]) {
                uniquePathsWithDepth[url] = j;
              }
            }
          }

          let dictionary: {
            [key: string]: { source: string; targets: Set<string> };
          } = {};

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
          dispatch({ type: "SET_NODES", payload: uniquePathsWithDepth });
          dispatch({ type: "SET_LINK_NODES", payload: dictionary });
          dispatch({ type: "SET_RUNTIME", payload: data.runtime as number });
          dispatch({
            type: "SET_ARTICLE_COUNT",
            payload: data.articleCount as number,
          });
        },
      });

      // Handle additional logic based on the response if needed
    } catch (err) {
      console.error(err);
      const errMsg =
        err instanceof Error ? err.message : "Something went wrong";
      toast.error(errMsg);
    } finally {
      setLoading(false);
    }
  };

  const handleOperations = async () => {
    await onSubmit();
    await handleSave();
  };

  return (
    <section className="flex flex-col gap-5">
      <div className="flex items-center justify-center gap-4">
        <span className="font-montserrat text-[21px] font-semibold">IDS</span>
        <Switch
          className="bg-white"
          checked={state.isBFS}
          onCheckedChange={(checked) =>
            dispatch({ type: "SET_ISBFS", payload: checked })
          }
        />
        <span className="font-montserrat text-[21px] font-semibold">BFS</span>
      </div>

      {process.env.NODE_ENV == "development" && (
        <Button
          size={"lg"}
          className="text-2xl sm:text-3xl bg-yellow-primary hover:bg-yellow-hover transition ease-in-out delay-150 hover:scale-102 hover:-translate-y-1 duration-300"
          onClick={handleOperations}
          disabled={loading}
        >
          {loading ? "Loading..." : "Go!"}
        </Button>
      )}
    </section>
  );
};

export default SwitchAPIReq;
