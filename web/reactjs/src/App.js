import React, { useState } from 'react';

import './App.css';

import Header from './components/Header.js';
import List from './components/List.js';
import Search from './components/Search.js';
import Loading from './components/Loading.js';
import Error from './components/Error.js';
import Info from './components/Info.js';


function App() {
  const [data, setData] = useState([])
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState(false)
  const [empty, setEmpty] = useState(false)

  const handleSearch = async (value) => {
    setData([])
    setError(false)
    setEmpty(false)

    if(!value) return setEmpty(true)

    setLoading(true)

    try {
      const response = await fetch(`/search?q=${value}`)
      if(response.status !== 200) throw new Error(`Error ${response.status}`)
      const results = await response.json()
      if(!results.highlights && results.highlights.length === 0) setEmpty(true)
      if(results.highlights) {
        setData(results.highlights)
      }
    } catch(e) {
      console.log(e)
      setError(true)
    }
    setLoading(false)
  }

  return (
    <div className="App">
      <Header />
      <div className="bg-white">
        <div className="mx-auto py-12 px-10 max-w-4xl">
          <Info />
          <Search onChange={handleSearch} disabled={loading} />
          { error && <Error /> }
          { empty && <p className="text-center">No result for your search</p> }
          { loading ? <Loading /> :
          <div className="space-y-12 flex flex-row justify-between gap-10">
            { data && data.length > 0 && <List data={data} /> }
          </div>
          }
        </div>
      </div>
    </div>
  );
}

export default App;
