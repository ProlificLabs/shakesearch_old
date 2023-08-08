const DEFAULT_OFFSET = 20
const DEFAULT_LIMIT = 20

const Pagination = {
  query: "",
  offset: DEFAULT_OFFSET,
  limit: DEFAULT_LIMIT,
  resetPagination(query) {
    Pagination.query = query;
    Pagination.offset = DEFAULT_OFFSET;
    Pagination.limit = DEFAULT_LIMIT;
  },
  increasePagination() {
    Pagination.offset += Pagination.limit;
  }
}

const Controller = {
  search: (ev) => {
    ev.preventDefault();
    const searchButton = document.getElementById("search");
    searchButton.setAttribute("disabled", true)
    const form = document.getElementById("form");
    const data = Object.fromEntries(new FormData(form));
    const query = data.query;
    if (!query || query.length === 0) {
      alert("Please enter a search query")
      return
    }
    const response = fetch(`/search?q=${query}`).then((response) => {
      if (response.status !== 200) {
        response.text().then((results) => {
          searchButton.removeAttribute("disabled")
          alert(results)
        })
        return;
      }
      response.json().then((results) => {
        Pagination.resetPagination(query);
        Controller.updateTable(query, results);
        searchButton.removeAttribute("disabled")
      });
    });
  },

  loadMore: (ev) => {
    ev.preventDefault();
    const loadMoreButton = document.getElementById("load-more");
    loadMoreButton.setAttribute("disabled", true)
    const response = fetch(`/search?q=${Pagination.query}&l=${Pagination.limit}&o=${Pagination.offset}`).then((response) => {
      if (response.status !== 200) {
        response.text().then((results) => {
          loadMore.removeAttribute("disabled")
          alert(results)
        })
        return;
      }
      response.json().then((results) => {
        Pagination.increasePagination()
        Controller.updateTable(Pagination.query, results, true);
        loadMore.removeAttribute("disabled")
      });
    });
  },

  updateTable: (query, results, appendResults = false) => {
    const table = document.getElementById("table-body");
    const rows = [];
    for (let result of results) {
      const markedText = result.replace(new RegExp(query, 'gi'), `<mark>${query}</mark>`);
      rows.push(`<tr><td>${markedText}</td></tr>`);
    }
    if (appendResults) {
      table.innerHTML +=  rows;
    } else {
      table.innerHTML = rows;
    }
  },
};

const form = document.getElementById("form");
form.addEventListener("submit", Controller.search);
const loadMore = document.getElementById("load-more");
loadMore.addEventListener("click", Controller.loadMore);
