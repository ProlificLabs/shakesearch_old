export const ResultSummary = ({ response, isLoading }) => {
  if (!response || response.length === 0) {
    return null;
  }

  const resultCount = response.length;

  return (
    <div className="text-sm font-light mt-2 ml-auto mr-auto text-slate-500 text-center">
      {isLoading ? (
        <div className="h-3 mt-4">
          <div className="animate-pulse h-2 bg-slate-200 rounded w-48">
            &nbsp;
          </div>
        </div>
      ) : (
        <div>
          Showing {resultCount} {`result${resultCount > 1 ? "s" : ""}`}
        </div>
      )}
    </div>
  );
};
