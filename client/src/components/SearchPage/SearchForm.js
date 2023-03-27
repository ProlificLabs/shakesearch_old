import React from 'react';
export default function SearchInput(props) {
  const textRef = React.useRef(null);
  const [inputText, setInputText] = React.useState('');
  const searchReminder = 'Please enter a search term';
  return (
    <div>
      <form onSubmit={(event) => {
        event.preventDefault();
        if (inputText.length > 0) {
          props.submitInput(inputText);
        } else {
          alert(searchReminder);
        }
      }} >
        <label htmlFor="email" className="block text-sm font-medium leading-6 text-gray-900">
          Enter a search term
        </label>
        <div className="mt-2 flex flex-row space-x-4">
          <input
            type="text"
            name="text"
            data-testid="search-input"
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
            onKeyDown={
              (event) => {
                if (event.key === 'Enter') {
                  event.preventDefault();
                  if (inputText.length > 0) {
                    props.submitInput(inputText);
                  } else {
                    alert(searchReminder);
                  }
                }
              }
            }
          />
          < button
            type='submit'
            className="rounded-md bg-indigo-600 py-1.5 px-2.5 text-sm font-semibold text-white shadow-sm hover:bg-indigo-500 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-indigo-600"
            data-testid="submit-button"
          >
            Search
          </button>
          <button
            type="reset"
            className="relative ml-3 inline-flex items-center rounded-md bg-white px-3 py-2 text-sm font-semibold text-gray-900 ring-1 ring-inset ring-gray-300 hover:bg-gray-50 focus-visible:outline-offset-0"
            data-testid="clear-button"
            onClick={() => {
              props.clearResults();
            }}
          >
            Clear
          </button>
        </div>
        <p className="mt-2 text-sm text-gray-500" id="search-input-description">
          Search with free keywords or wrap with "" for exact matches (case insensitive)
        </p>
      </form >
    </div >
  )
}
