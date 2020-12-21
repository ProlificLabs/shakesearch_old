import React, { useEffect, useState } from 'react';

const List = ({data}) => {
  const [list, setList] = useState([])

  useEffect(() => {
    if(data) {
      setList(data)
      return
    } 
    setList([])
  }, [data])

  const listItems = list.map((item, index) =>
    <li key={index} className="rounded-xl bg-gray-200 bg-opacity-30 p-5">
      <div className="space-y-4">
        <div className="space-y-2">
          <div className="text-sm">
            <p className="text-gray-500" dangerouslySetInnerHTML={{__html: item.text}}>
            </p>
          </div>
            <div className="text-xs leading-6 font-medium space-y-2 pt-5">
            <p className="space-x-5">
              <span className="text-gray-600">words: </span>
              <span className="text-gray-400">
              { item.matched_words && item.matched_words.map((w, i) => 
              <span key={i}>{w} </span>
              )}
            </span></p>
          </div>
        </div>
      </div>
    </li>
  );

  return <div className="flex-1 overflow-y-scroll h-600">
    <ul className="space-y-5">
      {listItems}
    </ul>
  </div>
}


export default List;
