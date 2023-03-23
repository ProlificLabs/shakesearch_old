import React from 'react';
import axios from 'axios';
import Header from '../../components/Header';
import ResultTable from '../../components/ResultTable';
import SearchInput from '../../components/SearchInput';

export default function SearchPage() {
  const [results, setResults] = React.useState(null);

  const submitInput = (input) => {
    axios.get(`https://7533-187-18-138-44.ngrok.io/search?q=${input}`,
      {
        headers: {
          "Content-Type": "application/json",
        }
      })
      .then((response) => {
        console.log(response);
        setResults(response.data);
      })
  };

  const clearResults = () => {
    setResults(null);
  };

  return (
    <div className="max-h-screen">
      <Header title="Shake Search" subtitle="Search for a text string in the complete works of Shakespeare" />
      <SearchInput submitInput={submitInput} clearResults={clearResults} />
      {results && <ResultTable results={results} />}
    </div>
  );
}
