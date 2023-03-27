import React from 'react';
import axios from 'axios';
import NProgress from 'nprogress';
import 'nprogress/nprogress.css';
import Header from '../../components/Header';
import ResultTable from '../../components/SearchPage/ResultTable';
import SearchInput from '../../components/SearchPage/SearchForm';

export default function SearchPage() {
  const [results, setResults] = React.useState(null);
  const [timeToComplete, setTimeToComplete] = React.useState(0);
  const [sumOfMatches, setSumOfMatches] = React.useState(0);

  const sumMatches = (resultsObj) => {
    let sum = 0;
    for (let i = 0; i < resultsObj.length; i++) {
      sum = sum + resultsObj[i].Matches;
    }
    return sum;
  }

  const submitInput = async (input) => {
    NProgress.start();
    const tokens = input.trim().match(/(?:"[^"]*"|\S)+/g);
    let strArray = [];
    for (let i = 0; i < tokens.length; i++) {
      strArray.push(`q=${encodeURIComponent(tokens[i].toLowerCase())}`);
    }
    const beginTime = new Date().getTime();
    axios.get(`/search?${strArray.join('&')}`,
      {
        headers: {
          "Content-Type": "application/json",
        }
      })
      .then(async (response) => {
        setTimeToComplete(new Date().getTime() - beginTime);
        try {
          // validate response data
          if (typeof response.data !== 'object') {
            NProgress.done();
            sendAlert();
            throw new Error('Invalid response');
          }
          setResults(response.data);
          setSumOfMatches(sumMatches(response.data));
        } catch (err) {
          console.error(err);
          NProgress.done();
        }
        NProgress.done();
      })
      .catch((error) => {
        console.error(error);
        setTimeToComplete(new Date().getTime() - beginTime);
        NProgress.done();
        sendAlert();
      })
  };

  const clearResults = () => {
    setResults(null);
    setTimeToComplete(0);
    setSumOfMatches(0);
  };

  function sendAlert() {
    alert("Sorry, we have got an error. Please, try again.");
  }

  return (
    <div className="max-h-screen">
      <Header
        title="Shake Search"
        subtitle="Search for a text string in the complete works of Shakespeare"
      />
      <SearchInput submitInput={submitInput} clearResults={clearResults} />
      {results &&
        <ResultTable
          results={results}
          sumOfMatches={sumOfMatches}
          timeToComplete={timeToComplete} />}
    </div>
  );


}
