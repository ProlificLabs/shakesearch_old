import React, { useState, useEffect } from "react";
import Modal from "../Modal/Modal.js";
import "./ResultsTable.css";

export default function ResultsTable(props) {
  const [data, setdata] = useState(props.data);
  const [searchterm, setsearchterm] = useState(props.searchTerm);
  const [totalwords, settotalwords] = useState(0);
  const [order, setorder] = useState("ASCENDING");

  // initialize
  useEffect(() => {
    setdata(props.data);
    setsearchterm(props.searchTerm);
    let wordcounter = 0;
    props.data.forEach((e) => (wordcounter += Number(e.wordCount)));
    settotalwords(wordcounter);
  }, [props.data]);

  // sort the results table in alphabetical order
  const alphabetical_sorting = (column_name) => {
    if (order === "ASCENDING") {
      const sorted = [...data].sort((a, b) =>
        a[column_name].toLowerCase() > b[column_name].toLowerCase() ? 1 : -1
      );
      setdata(sorted);
      setorder("DESCENDING");
    }
    if (order === "DESCENDING") {
      const sorted = [...data].sort((a, b) =>
        a[column_name].toLowerCase() < b[column_name].toLowerCase() ? 1 : -1
      );
      setdata(sorted);
      setorder("ASCENDING");
    }
  };

  // sort the results table in numeric order
  const numeric_sorting = (column_name) => {
    if (order === "ASCENDING") {
      const sorted = [...data].sort((a, b) =>
        Number(a[column_name]) > Number(b[column_name]) ? 1 : -1
      );
      setdata(sorted);
      setorder("DESCENDING");
    }
    if (order === "DESCENDING") {
      const sorted = [...data].sort((a, b) =>
        Number(a[column_name]) < Number(b[column_name]) ? 1 : -1
      );
      setdata(sorted);
      setorder("ASCENDING");
    }
  };

  return (
    <div className="search-results">
      "{searchterm}" appears a total of {totalwords} times
      <table className="results-table">
        <thead>
          <tr>
            {/* sort buttons and table headers */}
            <th onClick={() => alphabetical_sorting("title")}>
              Title <i className="fas fa-sort"></i>
            </th>
            <th onClick={() => numeric_sorting("wordCount")}>
              Word Count <i className="fas fa-sort"></i>
            </th>
            <th></th>
          </tr>
        </thead>
        <tbody>
          {/* results table data */}
          {data.map((d) => (
            <tr key={d.title}>
              <td>{d.title}</td>
              <td>{d.wordCount}</td>
              <td>
                <Modal
                  title={d.title}
                  wordCount={d.wordCount}
                  artPiece={d.artPiece}
                  searchTerm={searchterm}
                />
              </td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
}
