import React from "react";
import { hightText } from "../utils/text";

export const ListItem = ({ text, searchQuery }) => {
  return (
    <div className="block max-w-sm p-6 bg-white border border-gray-200 rounded-lg shadow hover:bg-gray-100 ">
      <h5 className="mb-2 text-2xl font-bold tracking-tight text-gray-900 ">
        Noteworthy technology acquisitions 2021
      </h5>
      <div className="font-normal text-gray-700">
        {hightText(text, searchQuery)}
      </div>
    </div>
  );
};
