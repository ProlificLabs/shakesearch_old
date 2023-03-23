import { useEffect, useState } from "react";

function Pagination(props) {
  return (
    <div
      className="flex items-center justify-between border-t border-gray-200 bg-white px-2 py-3 sm:px-3"
      aria-label="Pagination"
    >
      <div className="hidden sm:block">
        <p className="text-sm text-gray-700">
          Showing <span className="font-medium">{props.total > 0 ? (((props.page - 1) * 10) + 1) : 0}</span> to <span className="font-medium">{((props.page - 1) * 10) + 10 > props.total ? props.total : ((props.page - 1) * 10) + 10}</span> of{' '}
          <span className="font-medium">{props.total}</span> results
        </p>
      </div>
      <div className="flex flex-1 justify-between sm:justify-end">
        <a
          href="#previous"
          className="relative inline-flex items-center rounded-md bg-white px-3 py-2 text-sm font-semibold text-gray-900 ring-1 ring-inset ring-gray-300 hover:bg-gray-50 focus-visible:outline-offset-0"
          onClick={() => {
            if ((props.page - 1) >= 1)
              props.setPage(props.page - 1)
          }}
        >
          Previous
        </a>
        <a
          href="#next"
          className="relative ml-3 inline-flex items-center rounded-md bg-white px-3 py-2 text-sm font-semibold text-gray-900 ring-1 ring-inset ring-gray-300 hover:bg-gray-50 focus-visible:outline-offset-0"
          onClick={() => {
            if ((props.page + 1) <= Math.ceil(props.total / 10))
              props.setPage(props.page + 1)
          }}
        >
          Next
        </a>
      </div>
    </div>
  )
}

export default function ResultTable(props) {
  const [page, setPage] = useState(1);
  const [total, setTotal] = useState(0);

  useEffect(() => {
    setTotal(props.results.length);
  }, [props.results]);
  return (
    <div className="">
      <div className="mt-8 flow-root">
        <div className="-my-2 -mx-4 overflow-x-auto sm:-mx-6 lg:-mx-8">
          <div className="inline-block min-w-full py-2 align-middle sm:px-6 lg:px-8">
            <table className="min-w-full divide-y divide-gray-300">
              <thead>
                <tr>
                  <th scope="col" className="px-3 py-3.5 text-left text-sm font-semibold text-gray-900">
                    Results: {total} occurrance{total > 1 ? 's' : ''} found
                  </th>
                </tr>
              </thead>
              <tbody className="bg-white">
                {props.results.length > 0 && props.results.slice(((page - 1) * 10), page * 10).map((result, resultIdx) => (
                  <tr key={result.Idx} className={resultIdx % 2 === 0 ? undefined : 'bg-gray-50'}>
                    <td className="flex-wrap py-4 pl-4 pr-3 text-sm font-normal text-gray-900 sm:pl-3">
                      "...{result}..."
                    </td>
                  </tr>
                ))}
              </tbody>

            </table>
            {props.results.length > 0 && < Pagination page={page} setPage={setPage} total={total} />}
          </div>
        </div>
      </div>
    </div>
  )
}
