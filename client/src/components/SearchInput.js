import React from 'react';
export default function SearchInput(props) {
  const textRef = React.useRef(null);
  const [inputText, setInputText] = React.useState('');
  return (
    <div>
      <form onSubmit={(event) => {
        event.preventDefault();
        props.submitInput(inputText)
      }} >
        <label htmlFor="email" className="block text-sm font-medium leading-6 text-gray-900">
          Search Input
        </label>
        <div className="mt-2 flex flex-row space-x-4">
          <input
            type="text"
            name="text"
            id="text"
            className="block w-full rounded-md border-0 py-1.5 text-gray-900 shadow-sm ring-1 ring-inset ring-gray-300 placeholder:text-gray-400 focus:ring-2 focus:ring-inset focus:ring-indigo-600 sm:text-sm sm:leading-6"
            placeholder="Hamlet"
            aria-describedby="search-input-description"
            ref={textRef}
            onChange={
              (event) => {
                setInputText(event.target.value);
              }
            }
          />
          <button
            type='submit'
            className="rounded-md bg-indigo-600 py-1.5 px-2.5 text-sm font-semibold text-white shadow-sm hover:bg-indigo-500 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-indigo-600"
          >
            Search
          </button>
          <button
            type="reset"
            className="relative ml-3 inline-flex items-center rounded-md bg-white px-3 py-2 text-sm font-semibold text-gray-900 ring-1 ring-inset ring-gray-300 hover:bg-gray-50 focus-visible:outline-offset-0"
            onClick={() => {
              props.clearResults();
            }}
          >
            Clear
          </button>
        </div>
        <p className="mt-2 text-sm text-gray-500" id="search-input-description">
          Search for text (case insensitive)
        </p>
      </form>
    </div >
  )
}
