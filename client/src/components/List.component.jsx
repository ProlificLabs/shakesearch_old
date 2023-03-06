import React, { useEffect, useState } from "react";
import { ListItem } from "./ListItem.component";
import { Pagination } from "./Pagination.component";
import {
  DEFAULT_CURRENT_PAGE,
  DEFAULT_RESULTS_PER_PAGE,
} from "../consts/common.const";

export const List = React.memo(({ searchResult, searchQuery }) => {
  const resultsPerPage = DEFAULT_RESULTS_PER_PAGE;

  const [currentPage, setCurrentPage] = useState(DEFAULT_CURRENT_PAGE);
  const [perPageResults, setPerPageResults] = useState([]);
  const [startIndex, setStartIndex] = useState(0);
  const [endIndex, setEndIndex] = useState(0);

  useEffect(() => {
    const startIndex = (currentPage - 1) * resultsPerPage;
    const endIndex = startIndex + resultsPerPage;
    setStartIndex(startIndex);
    setEndIndex(endIndex);

    setPerPageResults(searchResult.slice(startIndex, endIndex));
  }, [currentPage, searchResult, resultsPerPage]);

  useEffect(() => {
    setCurrentPage(DEFAULT_CURRENT_PAGE);
  }, [searchQuery, searchResult]);

  return (
    <>
      <div className="mt-5 container m-auto grid grid-cols-2 gap-4 md:grid-cols-3 lg:grid-cols-5">
        {perPageResults.map((item) => (
          <ListItem key={item} text={item} searchQuery={searchQuery} />
        ))}
      </div>
      <Pagination
        searchResult={searchResult}
        currentPage={currentPage}
        setCurrentPage={setCurrentPage}
        startIndex={startIndex}
        endIndex={endIndex}
        resultsPerPage={resultsPerPage}
      />
    </>
  );
});
