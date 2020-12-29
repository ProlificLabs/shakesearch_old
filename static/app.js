const Controller = {
  search: (ev) => {
    ev.preventDefault();
    const form = document.getElementById("form");
    const data = Object.fromEntries(new FormData(form));
    const response = fetch(`/search?q=${data.query}`).then((response) => {
      response.json().then((results) => {
        Controller.updateTable(results, data.query);
      });
    });
  },

  updateTable: (results, query) => {
    const table = document.getElementById("table-body");
    let rows = '';
    for (let result of results) {
      let regExpQuery = new RegExp('('+query+')', 'gi');
      let formattedString = result.replace(regExpQuery, `<span style='background: yellow;'>$&</span>`);
      rows += `<tr>...${formattedString}...<tr/>`;
    }
    table.innerHTML = rows;
  },
};

const form = document.getElementById("form");
form.addEventListener("submit", Controller.search);
