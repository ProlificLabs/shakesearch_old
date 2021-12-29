const Controller = {
  search: (ev) => {
    ev.preventDefault();
    const form = document.getElementById("form");
    const data = Object.fromEntries(new FormData(form));
    fetch(`/search?q=${data.query}`).then((response) => {
      response.json().then((results) => {
        Controller.updateTable(results);
      });
    });
  },

  updateTable: (results) => {
    const table = document.getElementById("table-body");
    for (let result of results) {
      const tr = document.createElement("tr");
      const td = document.createElement("td");
      td.innerText = result;
      tr.appendChild(td);
      table.appendChild(tr);
    }
  },
};

const form = document.getElementById("form");
form.addEventListener("submit", Controller.search);
