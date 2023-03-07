import React from "react";
import { PageinationButton } from "./Buttons/PageinationButton.component";
import { BUTTON_DIRECTION } from "../consts/component.const";

export const Pagination = ({
  currentPage,
  searchResult,
  setCurrentPage,
  startIndex,
  endIndex,
  resultsPerPage,
}) => {
  const totalResult = searchResult.length;
  const totalPages = Math.ceil(totalResult / resultsPerPage);

  const handlePrevButton = () => {
    if (currentPage > 1) {
      setCurrentPage(currentPage - 1);
    }
  };

  const handleNextButton = () => {
    if (currentPage < totalPages) {
      setCurrentPage(currentPage + 1);
    }
  };

  return (
    <>
      <div className="flex flex-col items-center mt-2">
        <span className="text-sm text-gray-700 ">
          Showing
          <span className="font-semibold text-gray-900 ml-1 mr-1">
            {startIndex + 1}
          </span>
          to
          <span className="font-semibold text-gray-900 ml-1 mr-1">
            {endIndex > totalResult ? totalResult : endIndex}
          </span>
          of
          <span className="font-semibold text-gray-900 ml-1 mr-1">
            {totalResult}
          </span>
          Entries
        </span>
        <div className="inline-flex mt-2 xs:mt-0">
          {currentPage > 1 && (
            <PageinationButton
              direction={BUTTON_DIRECTION.prev}
              handleButtonClick={handlePrevButton}
            />
          )}

          {!!(currentPage < totalPages) && (
            <PageinationButton
              direction={BUTTON_DIRECTION.next}
              handleButtonClick={handleNextButton}
            />
          )}
        </div>
      </div>
    </>
  );
};
