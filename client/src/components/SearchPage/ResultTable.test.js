import React from "react";
import { render, screen } from "@testing-library/react";
import "@testing-library/jest-dom";
import ResultTable from "./ResultTable";


test("render empty state", () => {
  // render the component
  render(<ResultTable results={[]} sumOfMatches={0} timeToComplete={0} />);

  const numOfResults = screen.getByTestId("numOfResults");
  const numOfMatches = screen.getByTestId("numOfMatches");
  const resultsBody = screen.getByTestId("resultsBody");

  //assert the expected result
  expect(numOfResults).toHaveTextContent("0");
  expect(numOfMatches).toHaveTextContent("0");
  expect(resultsBody).toBeEmptyDOMElement();
});

test("render non-empty state", () => {
  // render the component
  render(
    <ResultTable
      results=
      {[
        { Snippet: "The King is alive!", Matches: 2 },
        { Snippet: "The King is dead!", Matches: 1 },
        { Snippet: "Hail the King!", Matches: 1 }
      ]}
      sumOfMatches={4}
      timeToComplete={0}
    />);

  const numOfResults = screen.getByTestId("numOfResults");
  const numOfMatches = screen.getByTestId("numOfMatches");
  const resultsBody = screen.getByTestId("resultsBody");

  //assert the expected results
  expect(numOfResults).toHaveTextContent("3");
  expect(numOfMatches).toHaveTextContent("4");
  expect(resultsBody).not.toBeEmptyDOMElement();
});
