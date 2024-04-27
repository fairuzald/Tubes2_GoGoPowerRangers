import { makeApiRequest } from "@/libs/helper";
import { AutoCompleteData } from "@/types/autocomplete";
import debounce from "lodash.debounce";
import { useEffect, useState } from "react";
import toast from "react-hot-toast";

const useAutoComplete = (
  searchQuery: string
): [AutoCompleteData[], (data: AutoCompleteData[]) => void] => {
  const [data, setData] = useState<AutoCompleteData[]>([]);

  const fetchData = async (search: string) => {
    try {
      await makeApiRequest({
        method: "GET",
        endpoint: "/autocomplete?search=" + search,
        headers: {
          "Content-Type": "application/json",
        },
        loadingMessage: "Fetching data...",
        successMessage: "Data fetched successfully!",
        onSuccess: async (data: any) => {
          setData(data.data as AutoCompleteData[]);
        },
      });
    } catch (error) {
      console.error("Error fetching data:", error);
      const errMsg =
        error instanceof Error ? error.message : "Unknown error occurred";
      toast.error(errMsg);
    }
  };

  useEffect(() => {
    const debouncedFetchData = debounce(fetchData, 1000);

    if (searchQuery) {
      debouncedFetchData(searchQuery);
    } else {
      setData([]);
    }

    return () => {
      debouncedFetchData.cancel();
    };
  }, [searchQuery]);

  return [data, setData];
};

export default useAutoComplete;
