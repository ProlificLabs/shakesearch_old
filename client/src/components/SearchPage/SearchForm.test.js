/* eslint-disable testing-library/prefer-screen-queries */
import React from "react";
import { render, fireEvent, waitFor } from "@testing-library/react";
import "@testing-library/jest-dom"
import SearchForm from "./SearchForm";


test("search input is empty and enabled", async () => {

  const { getByTestId } = render(<SearchForm />);
  const searchInput = getByTestId('search-input');

  expect(searchInput.value).toBe('');

  fireEvent.change(searchInput, { target: { value: 'King alive' } });
  expect(searchInput.value).toBe('King alive');
});
