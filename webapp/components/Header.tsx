import React, { FC } from "react";
import { FiSearch } from "react-icons/fi";

type Props = {
  query: string;
  setQuery: (val: string) => any;
  fetchResults: () => any;
};

const Header: FC<Props> = ({ query, setQuery, fetchResults }) => {
  return (
    <header className="py-4">
      <form
        action="#"
        onSubmit={(e) => {
          e.preventDefault();
          fetchResults();
        }}
      >
        <div className="flex flex-col sm:flex-row gap-4">
          <div className="flex gap-4 items-center px-4 py-2 rounded bg-[#F6F6F6] border border-zinc-200 focus-within:border-zinc-400 flex-1">
            <FiSearch className="text-zinc-400" />
            <input
              placeholder="Enter search query"
              className="bg-transparent outline-none border-none w-0 flex-1"
              type="text"
              value={query}
              onChange={(e) => setQuery(e.target.value)}
            />
          </div>
          <button className="px-4 py-2 rounded bg-violet-500 text-white" type="submit">
            Search
          </button>
        </div>
      </form>
    </header>
  );
};

export default Header;
