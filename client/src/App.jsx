import React, { useState, useEffect } from "react";
import "./App.css";

import { SearchBox } from "./components/SearchBox.component";
import { List } from "./components/List.component";

import { getSearchResults } from "./utils/api";

const AppComponent = ({ searchResult, setSearchQuery, searchQuery }) => {
  return (
    <div className="flex justify-center flex-col p-20">
      <SearchBox setSearchQuery={setSearchQuery} />
      <List searchResult={searchResult} searchQuery={searchQuery} />
    </div>
  );
};

function App() {
  const [searchQuery, setSearchQuery] = useState("");
  const [searchResult, setSearchResult] = useState([]);

  useEffect(() => {
    const fetchData = async () => {
      const searchResponse = await getSearchResults(searchQuery);
      setSearchResult(searchResponse);
    };

    if (searchQuery) {
      fetchData().catch(console.error);
    }
  }, [searchQuery]);

  return (
    <AppComponent
      searchResult={searchResult}
      setSearchQuery={setSearchQuery}
      searchQuery={searchQuery}
    />
  );
}

export default App;
