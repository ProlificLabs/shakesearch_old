const Controller = {
  search: (ev) => {
    ev.preventDefault();
    const form = document.getElementById("form");
    const data = Object.fromEntries(new FormData(form));
    const response = fetch(`/search?q=${data.query}`).then((response) => {
      response.json().then((results) => {
        Controller.updateTable(results);
      });
    });
  },

  updateTable: (results) => {
    const table = document.getElementById("table-body");
    if (!results || results.length == 0) {
      table.innerHTML = "No results found";
    } else {
      table.innerHTML = "";
      let rows = [];
      for (let result of results) {
        let row = document.createElement("tr");
        
        let title = document.createElement("h3");
        row.appendChild(title);
        title.appendChild(document.createTextNode(result.Title));
        
        let matchList = document.createElement("ul");
        row.appendChild(matchList);
        for (let match of result.Matches) {
          let matchItem = document.createElement("li");
          matchList.appendChild(matchItem);
          let matchLink = document.createElement("a");
          matchItem.appendChild(matchLink);

          matchLink.appendChild(document.createTextNode(match.Text));
          matchLink.setAttribute("href", `/work/${result.Index}?idx=${match.Index}`);
        }
        
        rows.push(row);
      }
      for (let row of rows) {
        table.appendChild(row);
      }
    }
  },
};

const form = document.getElementById("form");
form.addEventListener("submit", Controller.search);
