import React, { useState, useEffect } from "react";
import "./App.css";

import { SearchBox } from "./components/SearchBox.component";
import { List } from "./components/List.component";

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
    console.log("searchQuery changed", searchQuery);
    if (searchQuery) {
      fetch(`/search?q=${searchQuery}`)
        .then((response) => response.json())
        .then((data) => {
          setSearchResult(data);
          console.log(data);
        })
        .catch((error) => console.error(error));
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
