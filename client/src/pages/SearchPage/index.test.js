/* eslint-disable testing-library/prefer-screen-queries */
import React from "react";
import { render, fireEvent, waitFor } from "@testing-library/react";
import "@testing-library/jest-dom"
import axios from "axios";
import SearchPage from "./index.js";

afterEach(() => {
  axios.get.mockClear();
});

function mockCall() {
  axios.get.mockResolvedValueOnce({
    data: [
      { Snippet: "Chapter IV. The King is alive! This is the end", Matches: 2 },
      { Snippet: "Section III. The King is dead! Here we go", Matches: 1 },
      { Snippet: "Following we stand! Hail the King! Long and", Matches: 1 }
    ]
  });
}

test("render results when submit a search request", async () => {
  mockCall();

  const { getByTestId, findAllByTestId } = render(<SearchPage />);
  const searchInput = getByTestId('search-input');
  const submitButton = getByTestId('submit-button');

  fireEvent.change(searchInput, { target: { value: 'King alive' } });
  fireEvent.click(submitButton);

  await waitFor(() => expect(axios.get).toHaveBeenCalledTimes(1));

  expect(axios.get).toHaveBeenCalledWith('/search?q=king&q=alive', {
    headers: {
      'Content-Type': 'application/json',
    },
  });

  const resultRows = await findAllByTestId('resultRow');

  expect(resultRows[0]).toBeInTheDocument();
  expect(resultRows[1]).toBeInTheDocument();
  expect(resultRows[2]).toBeInTheDocument();
});
