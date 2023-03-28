import React, { useState } from "react";
import "./App.css";
import ResultsTable from "./components/ResultsTable/ResultsTable";

function App() {
  const [searchTerm, setSearchTerm] = useState("");
  const [searchResults, setSearchResults] = useState([]);
  const [caseSensitive, setCaseSensitive] = useState(false);
  const [wholeWord, setWholeWord] = useState(false);
  const [results, setresults] = useState(false);

  // handle input changes
  function handleInputChange(event) {
    setSearchTerm(event.target.value);
  }

  function handleCaseSensitiveChange(event) {
    setCaseSensitive(event.target.checked);
  }

  function handleWholeWordChange(event) {
    setWholeWord(event.target.checked);
  }

  // fetch data from backend and update state with the results
  async function fetchSearchResults() {
    const query = `q=${searchTerm}&caseSensitive=${caseSensitive}&wholeWord=${wholeWord}`;
    const response = await fetch(`http://127.0.0.1:3001/search?${query}`);
    const data = await response.json();
    setSearchResults(data);
    setresults(true);
  }

  // handle submission of data to backend
  async function handleSubmit(event) {
    event.preventDefault();
    await fetchSearchResults();
  }

  return (
    <div className="App">
      <div className="search-bar">
        {/* data submission form */}
        <form onSubmit={handleSubmit}>
          <input
            type="text"
            placeholder="try searching for hamlet"
            value={searchTerm}
            onChange={handleInputChange}
          />
          <button type="submit">
            <i className="fas fa-search"></i>
          </button>
          <label>
            Case Sensitive
            <input
              type="checkbox"
              checked={caseSensitive}
              onChange={handleCaseSensitiveChange}
            />
          </label>
          <label>
            Whole Words
            <input
              type="checkbox"
              checked={wholeWord}
              onChange={handleWholeWordChange}
            />
          </label>
        </form>
      </div>
      {/* render table of results */}
      {results ? (
        <ResultsTable data={searchResults} searchTerm={searchTerm} />
      ) : (
        <div>Welcome to te better Shakesearch!</div>
      )}
    </div>
  );
}

export default App;
