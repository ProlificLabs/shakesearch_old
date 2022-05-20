import { AutoSizer, List } from "react-virtualized";
import { ResultItem } from "../ResultItem";
import { LoadingSkeleton } from "../LoadingSkeleton";

export const ResultList = ({
  response,
  searchValue,
  isLoading,
  scrollToTop,
}) => {
  if (isLoading) {
    return <LoadingSkeleton />;
  }

  if (!response) {
    return null;
  }

  const rowRenderer = ({ key, index, style }) => {
    return (
      <div style={style} key={key}>
        <ResultItem
          itemId={key}
          result={response[index]}
          searchValue={searchValue}
        />
      </div>
    );
  };

  const getRowHeight = ({ index }) => {
    const rowCount =
      response[index].replace(/[\r\n]{3,}/g, "\r\n\r\n").split("\r\n").length ||
      0;
    return rowCount * 20 + 60;
  };

  return (
    <div className="mt-8 sm:mx-4 w-full h-full">
      {searchValue && !response?.length > 0 ? (
        <ResultItem
          itemId="notfound"
          result="No results found"
          searchValue={searchValue}
        />
      ) : (
        <AutoSizer>
          {({ height, width }) => (
            <List
              width={width}
              height={height}
              rowCount={response?.length}
              rowHeight={getRowHeight}
              rowRenderer={rowRenderer}
              scrollToIndex={scrollToTop}
            />
          )}
        </AutoSizer>
      )}
    </div>
  );
};
