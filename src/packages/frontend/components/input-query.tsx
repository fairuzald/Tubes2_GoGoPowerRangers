"use client";
import React, { useState } from "react";
import AutoCompleteInput from "@/components/autocomplete-input";
import useAutoCompleteData from "@/hooks/useAutoComplete";
import { MoveHorizontal } from "lucide-react";
import { useQueryContext } from "./query-provider";

export default function InputQuery() {
  // Use the custom hook to access state and dispatch function
  const { state, dispatch } = useQueryContext();
  const { source, destination } = state;

  // Fetch autocomplete data for source and destination
  const [sourceData, setSourceData] = useAutoCompleteData(source);
  const [destinationData, setDestinationData] =
    useAutoCompleteData(destination);

  // Swap source and destination
  const [isSwapped, setIsSwapped] = useState(false);

  const onSwap = () => {
    dispatch({ type: "SWAP" });
    const temp = sourceData;
    setSourceData(destinationData);
    setDestinationData(temp);
    setIsSwapped(!isSwapped);
  };
  console.log(destination, source);

  return (
    <section className="flex flex-col lg:flex-row gap-7 w-full max-w-[1500px] items-center justify-center">
      {/* AutoCompleteInput for source */}
      <div className="space-y-3">
        {!isSwapped ? (
          <h2 className="text-xl font-semibold text-white text-center">From</h2>
        ) : (
          <h2 className="text-xl font-semibold text-white text-center">To</h2>
        )}
        <AutoCompleteInput
          data={sourceData}
          type="text"
          placeholder="Type a Wikipedia article..."
          value={source}
          onChange={(e) =>
            dispatch({ type: "SET_SOURCE", payload: e.target.value })
          }
          onSelectValue={(url) =>
            dispatch({ type: "SET_SELECTED_SOURCE", payload: url })
          }
        />
      </div>
      {/* Button to trigger swap */}
      <button
        onClick={onSwap}
        className="w-fit h-fit rotate-90 lg:rotate-0 lg:mt-12"
      >
        <MoveHorizontal className="text-white" size={40} />
      </button>
      {/* AutoCompleteInput for destination */}
      <div className="space-y-3">
        {!isSwapped ? (
          <h2 className="text-xl font-semibold text-white text-center">To</h2>
        ) : (
          <h2 className="text-xl font-semibold text-white text-center">From</h2>
        )}
        <AutoCompleteInput
          data={destinationData}
          type="text"
          placeholder="Type a Wikipedia article..."
          value={destination}
          onChange={(e) =>
            dispatch({ type: "SET_DESTINATION", payload: e.target.value })
          }
          onSelectValue={(url) =>
            dispatch({ type: "SET_SELECTED_DESTINATION", payload: url })
          }
        />
      </div>
    </section>
  );
}
