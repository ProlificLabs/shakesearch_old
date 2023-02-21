import { SearchResultGrouped } from "@/types/SearchResult";
import React, { FC } from "react";
import ResultItemGroup from "./ResultItemGroup";

type Props = {
  loading: boolean;
  results: SearchResultGrouped[];
  error: string;
  query: string;
};

const Body: FC<Props> = ({ results, loading, error, query }) => {
  return (
    <main className="pb-4">
      {
        <ul className="flex flex-col gap-4">
          {results.map((result, i) => (
            <li key={`${result.play_name}_group_${i}`}>
              <ResultItemGroup result={result} query={query} />
            </li>
          ))}
        </ul>
      }
    </main>
  );
};

export default Body;
