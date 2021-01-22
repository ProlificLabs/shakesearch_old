import './style.scss';

class Controller {
  static search(ev) {
    ev.preventDefault();
    const form = document.getElementById('form');
    const { query } = Object.fromEntries(new FormData(form));
    Controller.queryRX = new RegExp(`(${query})`, 'gi');

    const response = fetch(`/search?q=${query}`).then((response) => {
      response.json().then((results) => {
        Controller.updateTable(results, query);
      });
    });
  }

  static updateTable(results, query) {
    const table = document.getElementById('table-body');
    const rows = [];

    for (let result of results) {
      const highlightedResult = Controller.highlightQuery(result);
      const row = `
        <tr>
          <td class="result-item">
            <div class="result-body">${highlightedResult}</div>
          </td>
        </tr>
      `;

      rows.push(row);
    }

    table.innerHTML = rows.join('');
  }

  static highlightQuery(result) {
    return result.replace(
      Controller.queryRX,
      '<span class="highlight-query">$1</span>'
    );
  }

  static insertPagination() {
    const paginationEL = document.getElementById('pagination');
    paginationEL.innerHTML = this.paginator.html();
  }
}

const form = document.getElementById('form');
form.addEventListener('submit', Controller.search);
