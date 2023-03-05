import React from "react";

export const SearchBox = () => {
  return (
    <div className="flex justify-center">
      <div className="relative">
        <input
          type="text"
          className="h-14 w-96 pl-10 pr-20 rounded-lg z-0 focus:shadow focus:outline-none border-2 border-zinc-400"
          placeholder="Search anything..."
        />
        <div className="absolute top-2 right-2">
          <button className="h-10 w-20 text-white rounded-lg bg-red-500 hover:bg-red-600">
            Search
          </button>
        </div>
      </div>
    </div>
  );
};
