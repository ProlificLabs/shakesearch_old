import { MAX_LINES_VISIBLE } from "@/config";
import { SearchResultApi, SearchResultGrouped } from "@/types/SearchResult";
import React, { FC, useEffect, useMemo, useState } from "react";
import ResultItem from "./ResultItem";

type Props = {
  result: SearchResultGrouped;
  query: string;
};

const transformHTML = (value: string, query: string) => {
  query = query.trim();
  const words = query
    .split(" ")
    .filter((value, index, array) => array.indexOf(value) === index);
  let transformedHtml = (value || "").toString();
  words.forEach((word) => {
    var re = new RegExp(word, "ig");
    transformedHtml = transformedHtml.replace(
      re,
      `<i class='highlight'>${word}</i>`
    );
  });

  return transformedHtml;
};

const ResultItemGroup: FC<Props> = ({ result, query }) => {
  const { play_name, group } = result;

  const [showViewMore, setShowViewMore] = useState(
    group.length > MAX_LINES_VISIBLE
  );
  const [collapsed, setCollapsed] = useState(true);
  const [groupItems, setGroupItems] = useState(
    group.length > MAX_LINES_VISIBLE ? group.slice(0, MAX_LINES_VISIBLE) : group
  );

  const playHtml = useMemo(() => {
    const transformedHtml = transformHTML(play_name, query);
    return transformedHtml;
  }, [query, play_name]);

  const toggleCollapse = () => {
    const newValue = !collapsed;
    setCollapsed(newValue);

    if (newValue) {
      setGroupItems(group.slice(0, MAX_LINES_VISIBLE));
    } else {
      setGroupItems(group);
    }
  };

  return (
    <div className="bg-white p-4 rounded-md border border-zinc-200">
      <h3
        className="text-xl font-bold mb-2"
        dangerouslySetInnerHTML={{ __html: playHtml }}
      />
      <ul className="flex flex-col gap-1">
        {groupItems.map((item, i) => (
          <li key={`${item.play_name}_${i}`}>
            <ResultItem result={item} query={query} />
          </li>
        ))}
        {showViewMore && (
          <li className="flex justify-end">
            <button
              className="text-xs px-2 py-0.5 border rounded border-violet-500 text-violet-500"
              onClick={toggleCollapse}
            >
              {collapsed ? "View More" : "View Less"}
            </button>
          </li>
        )}
      </ul>
    </div>
  );
};

export default ResultItemGroup;
