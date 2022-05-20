import { useState } from "react";
import "./App.css";
import {
  ErrorMessage,
  ResultList,
  ResultSummary,
  ScrollToTop,
  SearchBox,
  VisuallyHidden,
} from "./components";

function App() {
  const [searchResult, setSearchResult] = useState();
  const [isLoading, setIsLoading] = useState(false);
  const [searchValue, setSearchValue] = useState("");
  const [scrollToTop, setScrollToTop] = useState(-1);
  const [errorMessage, setErrorMessage] = useState("");

  const fetchResults = async (searchTerm) => {
    if (!searchTerm) return;

    setErrorMessage("");
    setSearchValue(searchTerm);
    setSearchResult();
    setIsLoading(true);
    try {
      const response = await fetch(
        `/search?q=${searchTerm.trim().toLowerCase()}`
      );
      const result = await response.json();
      setSearchResult(result);
    } catch (error) {
      setErrorMessage("Could not fetch search results. Please try again.");
    }
    setIsLoading(false);
  };

  return (
    <>
      <header>
        <VisuallyHidden>
          <h1>ShakeSearch</h1>
        </VisuallyHidden>
      </header>
      <main className="container mx-auto px-4 md:px-10 pt-10 max-w-3xl h-screen">
        <SearchBox fetchResults={fetchResults} />
        <ResultSummary response={searchResult} isLoading={isLoading} />
        <ErrorMessage message={errorMessage} />
        <ResultList
          isLoading={isLoading}
          response={searchResult}
          searchValue={searchValue}
          scrollToTop={scrollToTop}
        />
        <ScrollToTop scroll={setScrollToTop} response={searchResult} />
      </main>
    </>
  );
}

export default App;
