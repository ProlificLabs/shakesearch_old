import React from "react";
import "./App.css";

import { SearchBox } from "./components/SearchBox.component";
import { List } from "./components/List.component";

export default function App() {
  return (
    <div className="flex justify-center flex-col p-20">
      <SearchBox />
      <List />
    </div>
  );
}
