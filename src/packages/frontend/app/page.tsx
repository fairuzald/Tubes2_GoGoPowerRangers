"use client";
import InputQuery from "@/components/input-query";
import { useQueryContext } from "@/components/query-provider";
import type { PathInfo } from "@/types/result";
import React from 'react'; // Added import for React
import toast from "react-hot-toast";

// Cache to store URL info results
const infoCache = new Map<string, PathInfo>();

const fetchInfoUrl = async (url: string) => {
  if (!url) throw new Error('No url provided');

  // Check cache before making a request to the server
  if (infoCache.has(url)) {
    return infoCache.get(url)!;
  }

  try {
    const response = await fetch("/api/url-info", {
      headers: {
        'Content-Type': 'application/json',
      },
      method: 'POST',
      body: JSON.stringify({ url }),
    });

    const data = await response.json();
    if (!response.ok) {
      throw new Error(data.message || 'Something went wrong');
    }

    const info = data.data;
    infoCache.set(url, info); // Store info in cache

    return info;
  } catch (err) {
    console.error(err);
    const errMsg = err instanceof Error ? err.message : 'Something went wrong';
    toast.error(errMsg);
  }
};

export default function Home() {
  const { state, dispatch } = useQueryContext();

  const onSubmit = async (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    // Validation check
    if (!state.selectedSource || !state.selectedDestination) return toast.error('Please select source and destination from select input');

    // Binding loading toast to a variable
    const loading = toast.loading('Finding path...');

    // Fetching data from the server
    try {
      const url = state.method ==='ids' ? '/api/ids' : '/api/bfs';
      const response = await fetch(url, {
        cache: 'no-cache',
        headers: {
          'Content-Type': 'application/json',
        },
        method: 'POST',
        body: JSON.stringify({ source: state.selectedSource, destination: state.selectedDestination }),
      });

      const data = await response.json();
      if (!response.ok) {
        throw new Error(data.message || 'Something went wrong');
      }
      
      const result = data.paths as string[][];

      // Using Set to store unique URLs
      const uniquePaths = new Set<string>();
      result.forEach(path => path.forEach(url => uniquePaths.add(url)));

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
      const filteredResult = Object.values(resultWithInfo).filter((path) => path.length > 0);
      dispatch({ type: "SET_RESULT", payload: filteredResult });

    } catch (err) {
      console.error(err);
      const errMsg = err instanceof Error ? err.message : 'Something went wrong';
      toast.error(errMsg);
    }
    finally {
      toast.dismiss(loading);
    }

  };

  return (
    <main className="flex min-h-screen flex-col items-center p-24 w-full">
      <InputQuery /> 
      <form onSubmit={onSubmit} className="flex flex-col items-center">
        <button type="submit" className="btn btn-primary mt-4">Find Path</button>
      </form>
    </main>
  );
}
