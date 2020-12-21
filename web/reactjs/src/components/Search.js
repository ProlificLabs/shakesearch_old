import React, { useState } from 'react';

const Search = ({disabled, onChange}) => {
  const [search, setSearch] = useState("")

  const handleClick = () => {
    onChange(search)
  }

  const handleChange = (e) => {
    setSearch(e.target.value)
  }

  const handleKeyDown = (e) => {
    if (e.key === 'Enter') {
      onChange(e.target.value)
    }
  }

  return <div className="flex-1 px-2 flex items-center justify-center my-10">
    <div className="w-full">
      <label htmlFor="search" className="sr-only">Search</label>
      <div className="relative">
        <div className="pointer-events-none absolute inset-y-0 left-0 pl-3 flex items-center">
          <svg className="flex-shrink-0 h-5 w-5 text-gray-400" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20" fill="currentColor" aria-hidden="true">
          <path fillRule="evenodd" d="M8 4a4 4 0 100 8 4 4 0 000-8zM2 8a6 6 0 1110.89 3.476l4.817 4.817a1 1 0 01-1.414 1.414l-4.816-4.816A6 6 0 012 8z" clipRule="evenodd" />
  </svg>
  </div>
  <input name="search" id="search" value={search} onChange={handleChange} onKeyDown={handleKeyDown} className="block w-full bg-white border border-gray-300 rounded-md py-2 pl-10 pr-3 text-sm placeholder-gray-500 focus:outline-none focus:text-gray-900 focus:placeholder-gray-400 focus:ring-1 focus:ring-gray-900 focus:border-gray-900 sm:text-sm disabled:opacity-50" placeholder="Search" type="search" disabled={disabled} />
  </div>
  </div>
  <div className="px-10">
    <button onClick={handleClick} type="button" className="bg-white rounded-md font-medium text-purple-600 hover:text-purple-400 
focus:outline-none disabled:opacity-50:w
" disabled={disabled} >
            Search
          </button>
          </div>

      </div>
}

export default Search;
