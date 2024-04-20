"use client";
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
  const [destinationData, setDestinationData] = useAutoCompleteData(destination);

  // Swap source and destination
  const onSwap = () => {
    dispatch({ type: "SWAP" });
    const temp = sourceData;
    setSourceData(destinationData);
    setDestinationData(temp);
  };

  return (
    <div className="flex gap-7 w-full max-w-[1300px]">
      {/* AutoCompleteInput for source */}
      <div className="space-y-3">
        <h2 className="text-xl font-semibold text-white text-center">From</h2>
        <AutoCompleteInput
          data={sourceData}
          type="text"
          placeholder="Source...."
          value={source}
          onChange={(e) => dispatch({ type: "SET_SOURCE", payload: e.target.value })}
          onSelectValue={(url) => dispatch({ type: "SET_SELECTED_SOURCE", payload: url })}
        />
      </div>
      {/* Button to trigger swap */}
      <button onClick={onSwap} className="w-fit h-fit mt-12">
        <MoveHorizontal className="text-white" size={40} />
      </button>
      {/* AutoCompleteInput for destination */}
      <div className="space-y-3">
        <h2 className="text-xl font-semibold text-white text-center">To</h2>
        <AutoCompleteInput
          data={destinationData}
          type="text"
          placeholder="Destination...."
          value={destination}
          onChange={(e) => dispatch({ type: "SET_DESTINATION", payload: e.target.value })}
          onSelectValue={(url) => dispatch({ type: "SET_SELECTED_DESTINATION", payload: url })}
        />
      </div>
    </div>
  );
}
