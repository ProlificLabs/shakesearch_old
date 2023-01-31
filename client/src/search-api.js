import axios from "axios";

const SearchAPI = {
  fetchResults: ({ queryKey }) => {
    const [, { query }] = queryKey;

    return axios
      .get(`${process.env.REACT_APP_API_URL}/search?q=${query}`)
      .then(({ data }) => data);
  },
};

export default SearchAPI;
