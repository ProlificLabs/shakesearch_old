const Controller = {
  initialPage: 1,

  search: (ev) => {
    ev.preventDefault();
    const query = Controller.getQuery();
    const response = fetch(`/search?q=${query}&pageSize=20&page=1`).then((response) => {
      response.json().then((results) => {
        Controller.updateTable(results);
      });
    });
  },

  updateTable: (results) => {
    if (results.length === 0) return
    const table = document.getElementById("table-body");
    for (let result of results) {
      const newTaleRow = table.insertRow(-1);
      const newCell = newTaleRow.insertCell(0);
      newCell.innerText = result;
    }
  },

  loadMore: (currentPage) => {
    Controller.initialPage++;
    const query = Controller.getQuery();
    const response = fetch(`/search?q=${query}&pageSize=20&page=${Controller.initialPage}`).then((response) => {
      response.json().then((results) => {
        Controller.updateTable(results);
      });
    });
  },

  getQuery: () => {
    const form = document.getElementById("form");
    const data = Object.fromEntries(new FormData(form));
    return data.query
  }
};

const form = document.getElementById("form");
const loadMore = document.getElementById("load-more");
form.addEventListener("submit", Controller.search);
loadMore.addEventListener("click", Controller.loadMore);
