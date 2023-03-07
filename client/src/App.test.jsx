import React from "react";
import { render, screen, fireEvent, act } from "@testing-library/react";
import user from "@testing-library/user-event";

import App from "./App";
import { MOCKED_SEARCH_RESULTS } from "./tests/apiMock";

describe("<App />", () => {
  beforeEach(() => {
    render(<App />);
  });

  it("page render search bar", async () => {
    const button = screen.getByRole("button", { name: "Search" });

    expect(button).toBeTruthy;
  });

  it("page should render search Result", async () => {
    jest.mock("./utils/api", () => ({
      getSearchResults: jest.fn(() => Promise.resolve(MOCKED_SEARCH_RESULTS)),
    }));

    const searchField = screen.queryByTestId("searchField");

    await act(async () => {
      fireEvent.change(searchField, { target: { value: "Hamlet" } });
    });

    const searchButton = screen.getByRole("button", { name: "Search" });

    await act(async () => {
      user.type(searchField, "Hamlet");
      user.click(searchButton);
    });

    const searchResults = screen.queryByTestId("searchResults");
    const listItems = screen.queryAllByTestId("listItem");

    expect(searchResults).toBeTruthy;
    expect(listItems.length > 0).toBeTruthy;
  });
});
