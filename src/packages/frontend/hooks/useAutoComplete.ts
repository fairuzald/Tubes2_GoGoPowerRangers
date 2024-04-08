import { AutoCompleteData } from '@/types/autocomplete';
import debounce from 'lodash.debounce';
import { useEffect, useState } from 'react';
import toast from 'react-hot-toast';

const useAutoComplete = (searchQuery: string): [AutoCompleteData[], (data: AutoCompleteData[]) => void] => {
  const [data, setData] = useState<AutoCompleteData[]>([]);

  const fetchData = async (search: string) => {
    try {
      const response = await fetch(`/api/autocomplete?search=${search}`, {
        cache: 'no-cache',
        headers: {
          'Content-Type': 'application/json',
        },
        method: 'GET',
      });
      if (!response.ok) {
        throw new Error('Network response was not ok');
      }

      const fetchedData: AutoCompleteData[] = await response.json();
      setData(fetchedData);
    } catch (error) {
      console.error('Error fetching data:', error);
      const errMsg = error instanceof Error ? error.message : 'Unknown error occurred';
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
