const Controller = {
  limit: 20,  // Number of results to display per fetch
  offset: 0,  // Starting offset

  search: (ev) => {
    ev.preventDefault();
    const form = document.getElementById("form");
    const data = Object.fromEntries(new FormData(form));
    const response = fetch(`/search?q=${data.query}&limit=${Controller.limit}&offset=${Controller.offset}`).then((response) => {
      response.json().then((results) => {
        Controller.updateTable(results);

        // Display the load-more button if results were found.
        if (results.length > 0) {
          document.getElementById("load-more").style.display = 'block';
        }
      });
    });
  },

  updateTable: (results) => {
    const table = document.getElementById("table-body");
    if(Controller.offset === 0) table.innerHTML = ""; // Clear the table when it's a new search
    for (let result of results) {
      const row = document.createElement('tr');
      const cell = document.createElement('td');
      cell.textContent = result;
      row.appendChild(cell);
      table.appendChild(row);
    }
  },
};


const form = document.getElementById("form");
form.addEventListener("submit", (ev) => {
  Controller.offset = 0; // Reset offset when a new search is performed
  Controller.search(ev);
});

const loadMoreButton = document.getElementById("load-more");
loadMoreButton.addEventListener("click", () => {
  Controller.offset += Controller.limit; // Increase offset for each "load more" action
  Controller.search(new Event('click'));
});