import { useState } from "react";

export const SearchBox = ({ fetchResults }) => {
  const [searchTerm, setSearchTerm] = useState("");

  const handleSubmit = async (event) => {
    event.preventDefault();
    fetchResults(searchTerm.trim());
  };

  return (
    <div className="flex flex-col">
      <form className="group relative flex-1" onSubmit={handleSubmit}>
        <svg
          width="20"
          height="20"
          fill="currentColor"
          className="absolute left-3 top-1/2 -mt-2.5 text-slate-400 pointer-events-none group-focus-within:text-emerald-600"
          aria-hidden="true"
        >
          <path
            fillRule="evenodd"
            clipRule="evenodd"
            d="M8 4a4 4 0 100 8 4 4 0 000-8zM2 8a6 6 0 1110.89 3.476l4.817 4.817a1 1 0 01-1.414 1.414l-4.816-4.816A6 6 0 012 8z"
          />
        </svg>
        <svg
          onClick={() => setSearchTerm("")}
          width="20"
          height="20"
          viewBox="0 0 72 72"
          fill="currentColor"
          className="absolute right-4 top-1/2 -mt-2.5 text-slate-500 cursor-pointer group-focus-within:text-emerald-600"
          aria-hidden="true"
        >
          <g id="line">
            <line
              x1="17.5"
              x2="54.5"
              y1="17.5"
              y2="54.5"
              fill="none"
              stroke="currentColor"
              strokeLinecap="round"
              strokeLinejoin="round"
              strokeMiterlimit="10"
              strokeWidth="2"
            />
            <line
              x1="54.5"
              x2="17.5"
              y1="17.5"
              y2="54.5"
              fill="none"
              stroke="currentColor"
              strokeLinecap="round"
              strokeLinejoin="round"
              strokeMiterlimit="10"
              strokeWidth="2"
            />
          </g>
        </svg>
        <input
          className="disabled:bg-slate-100 disabled:text-slate-400 focus:ring-1 focus:ring-emerald-600 focus:outline-none appearance-none w-full text-md font-light leading-6 text-slate-900 placeholder-slate-400 rounded-full py-4 pl-10 ring-1 ring-slate-300 shadow-sm"
          type="text"
          aria-label="Filter projects"
          placeholder="Try searching for 'Castle'"
          autoFocus
          value={searchTerm}
          onChange={(event) => setSearchTerm(event.target.value)}
        />
      </form>
      <button
        onClick={handleSubmit}
        aria-label="Search button"
        className="bg-emerald-50 p-2 rounded-full w-32 font-light mt-8 mr-auto ml-auto flex-1 hover:shadow-sm text-slate-700"
      >
        Search
      </button>
    </div>
  );
};
