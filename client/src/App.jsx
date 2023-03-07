import React, { useState, useEffect, useMemo } from "react";
import "./App.css";

import { SearchBox } from "./components/SearchBox.component";
import { List } from "./components/List.component";
import { Message } from "./components/Message.component";
import { getSearchResults } from "./utils/api";

const AppComponent = ({
  searchResult,
  setSearchQuery,
  searchQuery,
  errorMessage,
  isLoading,
}) => {
  return (
    <div className="flex justify-center flex-col p-20">
      <SearchBox setSearchQuery={setSearchQuery} />

      {errorMessage ? (
        <Message text={errorMessage} warning />
      ) : isLoading ? (
        <Message text={"Loading ..."} />
      ) : (
        <List searchResult={searchResult} searchQuery={searchQuery} />
      )}
    </div>
  );
};

function App() {
  const [searchQuery, setSearchQuery] = useState("");
  const [searchResult, setSearchResult] = useState([]);
  const [errorMessage, setErrorMessage] = useState("");
  const [isLoading, setIsLoading] = useState(false);

  const cachedResults = useMemo(() => {
    const cached = sessionStorage.getItem(searchQuery);
    return cached ? JSON.parse(cached) : null;
  }, [searchQuery]);

  useEffect(() => {
    const fetchData = async () => {
      setIsLoading(true);

      if (cachedResults) {
        setSearchResult(cachedResults);
        setIsLoading(false);
      } else {
        try {
          const searchResponse = await getSearchResults(searchQuery);
          setSearchResult(searchResponse);
          sessionStorage.setItem(searchQuery, JSON.stringify(searchResponse));
          setIsLoading(false);
        } catch (error) {
          console.log("fetch data error", error);
          setIsLoading(false);
          setErrorMessage(
            "Fetch data error from server, please try again later"
          );
        }
      }
    };

    if (searchQuery) {
      fetchData();
    }
  }, [searchQuery]);

  return (
    <AppComponent
      searchResult={searchResult}
      setSearchQuery={setSearchQuery}
      searchQuery={searchQuery}
      isLoading={isLoading}
      errorMessage={errorMessage}
    />
  );
}

export default App;
