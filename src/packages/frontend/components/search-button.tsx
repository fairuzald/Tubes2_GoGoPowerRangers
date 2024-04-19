import React from "react";
import { Button } from "./button";

interface SearchButtonProps {
  isFetchingResults: boolean;
  fetchShortestPaths: () => void;
}

const SearchButton: React.FC<SearchButtonProps> = ({
  isFetchingResults,
  fetchShortestPaths,
}) => {
  if (isFetchingResults) {
    return null;
  }

  return (
    <Button
      className="w-60 h-18 mx-auto mb-10 text-4xl rounded-lg 
                 sm:w-50 sm:h-15 sm:text-3xl bg-[#fca311] hover:bg-[#ffd966] transition ease-in-out delay-150 hover:-translate-y-1 hover:scale-110 duration-300  "
      onClick={fetchShortestPaths}
    >
      Go!
    </Button>
  );
};

export default SearchButton;
