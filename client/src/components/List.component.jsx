import React from "react";
import { ListItem } from "./ListItem.component";

export const List = React.memo(({ searchResult, searchQuery }) => {
  return (
    <div className="mt-5 container m-auto grid grid-cols-2 gap-4 md:grid-cols-3 lg:grid-cols-5">
      {searchResult.map((item) => (
        <ListItem key={item} text={item} searchQuery={searchQuery} />
      ))}
    </div>
  );
});
