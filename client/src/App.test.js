import React from 'react';
import { render, screen } from '@testing-library/react';
import '@testing-library/jest-dom'

import App from './App';


test('renders header', () => {
  render(<App />);
  const headerElement = screen.getByTestId("header");
  expect(headerElement).toBeInTheDocument();
});
