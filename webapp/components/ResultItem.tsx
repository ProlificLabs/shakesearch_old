import { SearchResultApi } from "@/types/SearchResult";
import React, { FC, useMemo } from "react";

type Props = {
  result: SearchResultApi;
  query: string;
};

const transformHTML = (value: string, query: string) => {
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

const ResultItem: FC<Props> = ({ result, query }) => {
  const text = result.line;
  const { play_name, character, scene_number, act_number } = result;

  const textHtml = useMemo(() => {
    const transformedHtml = transformHTML(text, query);
    return transformedHtml;
  }, [query, text]);

  const playHtml = useMemo(() => {
    const transformedHtml = transformHTML(play_name, query);
    return transformedHtml;
  }, [query, play_name]);

  const sceneHtml = useMemo(() => {
    if (!scene_number) return "";
    const transformedHtml = transformHTML(scene_number, query);
    return transformedHtml;
  }, [query, scene_number]);

  const actHtml = useMemo(() => {
    if (!act_number) return "";
    const transformedHtml = transformHTML(act_number, query);
    return transformedHtml;
  }, [query, act_number]);

  const characterHtml = useMemo(() => {
    if (!character) return "";
    const transformedHtml = transformHTML(character, query);
    return transformedHtml;
  }, [query, character]);

  return (
    <div>
      <div
        className="text-gray-500"
        dangerouslySetInnerHTML={{ __html: textHtml }}
      />
      <div className="text-gray-500 text-xs flex flex-wrap gap-2 mt-2">
        {characterHtml.length > 0 && (
          <div
            className="border border-gray-400 px-2 py-0.5 rounded-full"
            dangerouslySetInnerHTML={{ __html: characterHtml }}
          />
        )}
        {actHtml.length > 0 && (
          <div
            className="border border-gray-400 px-2 py-0.5 rounded-full"
            dangerouslySetInnerHTML={{ __html: "act: " + actHtml }}
          />
        )}
        {sceneHtml.length > 0 && (
          <div
            className="border border-gray-400 px-2 py-0.5 rounded-full"
            dangerouslySetInnerHTML={{ __html: "scene: " + sceneHtml }}
          />
        )}
      </div>
    </div>
  );
};

export default ResultItem;
