import "./ResultItem.css";

export const ResultItem = ({ searchValue, result, itemId }) => {
  return (
    <div className=" text-slate-800 p-4">
      <pre className="itemContainer font-quote font-extralight text-sm overflow-x-scroll border-l border-slate-300 pl-6">
        {result.split("\r\n").map((line, index) => {
          return (
            <ResultLine
              key={`${itemId}-${index}`}
              line={line}
              searchValue={searchValue}
            />
          );
        })}
      </pre>
    </div>
  );
};

const ResultLine = ({ line, searchValue }) => {
  const searchValueIndex = line
    .toLowerCase()
    .indexOf(searchValue.toLowerCase());

  if (searchValueIndex === -1) {
    const [speakerNameIndex, speakerName] = nameIndex(line);
    return (
      <div className={`${line ? "" : "mt-4"}`}>
        <span className="font-normal">{speakerName}</span>
        <span>
          {speakerNameIndex > -1 ? line.substring(speakerNameIndex) : line}
        </span>
      </div>
    );
  }

  const preSearchValue = line.substring(0, searchValueIndex);
  const postSearchValue = line.substring(
    searchValueIndex + searchValue.length,
    line.length
  );
  const exactSearchItem = line.substring(
    searchValueIndex,
    searchValueIndex + searchValue.length
  );
  return (
    <div className="bg-emerald-50 rounded-xl">
      <span>{preSearchValue}</span>
      <span className="text-white bg-emerald-800 rounded-lg px-1 font-extrabold">
        {exactSearchItem}
      </span>
      <span>{postSearchValue}</span>
    </div>
  );
};

const nameIndex = (line) => {
  const dotIndex = line.indexOf(".");
  if (dotIndex === -1) {
    return [-1];
  }

  const nameString = line.substring(0, dotIndex + 1);
  if (nameString.toUpperCase() === nameString) {
    return [dotIndex + 1, nameString];
  }
  return [-1];
};
