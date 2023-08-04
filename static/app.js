let offset = 0;

const Controller = {
  search: async (ev) => {
    ev.preventDefault();
    const form = document.getElementById("form");
    const data = Object.fromEntries(new FormData(form));
    const req = await fetch(`/search?q=${data.query}&offset=${offset}`);
    const res = await req.json();
    Controller.updateTable(res);
  },

  loadMore: () => {
    offset += 20; // Increment offset by 20
    Controller.search(new Event("submit")); // Trigger a new search
  },

  updateTable: (results) => {
    const table = document.getElementById("table-body");
    const rows = [];
    for (let result of results) {
      rows.push(`<tr><td>${result}</td></tr>`);
    }
    table.innerHTML += rows.join(""); // Append new rows to the table
  },
};

const form = document.getElementById("form");
form.addEventListener("submit", Controller.search);

const loadMoreButton = document.getElementById("load-more");
loadMoreButton.addEventListener("click", Controller.loadMore);
