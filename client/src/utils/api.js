export const getSearchResults = async (searchQuery) => {
  const res = await fetch(
    `${process.env.REACT_APP_API_URL}/search?q=${searchQuery}`
  );

  return res.json();
};
