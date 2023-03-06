export const getSearchResults = async (searchQuery) => {
  const res = await fetch(`/search?q=${searchQuery}`);

  return res.json();
};
