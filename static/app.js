const WORK_ID_ATTR_NAME = "work-id"
const MATCH_ID_ATTR_NAME = "match-id"
const QUERY_ATTR_NAME = "query-text"

const Controller = {
  search: (ev) => {
    ev.preventDefault();
    const form = document.getElementById("form");
    const data = Object.fromEntries(new FormData(form));
    const response = fetch(`/search?q=${data.query}`).then((response) => {
      response.json().then((results) => {
        Controller.displaySearchResults(results, data.query);
      });
    });
  },

  getWork: (ev) => {
    ev.preventDefault();
    let workId = ev.target.getAttribute(WORK_ID_ATTR_NAME);
    let matchId = ev.target.getAttribute(MATCH_ID_ATTR_NAME);
    if (workId && matchId) {
      const response = fetch(`/work?work=${workId}`).then((response) => {
        response.json().then((work) => {
          let query = document.getElementById("table-body").getAttribute(QUERY_ATTR_NAME);
          Controller.displayWork(work, query, matchId);
        });
      });
    } else {
      console.error("no work id or match id found");
    }
  },

  displaySearchResults: (results, query) => {
    const table = document.getElementById("table-body");
    table.setAttribute(QUERY_ATTR_NAME, query);
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
        result.Matches.forEach((match, matchIndex) => {
          let matchItem = document.createElement("li");
          matchItem.setAttribute(WORK_ID_ATTR_NAME, result.Index);
          matchItem.setAttribute(MATCH_ID_ATTR_NAME, matchIndex);
          matchItem.appendChild(document.createTextNode(match));
          matchItem.onclick = Controller.getWork;
          matchItem.style.color = 'blue';
          matchItem.style.textDecoration = 'underline';
          matchItem.style.cursor = 'pointer';
          matchList.appendChild(matchItem);
        });
        
        rows.push(row);
      }
      for (let row of rows) {
        table.appendChild(row);
      }
    }
  },

  displayWork: (work, query, matchId) => {
    const table = document.getElementById("table-body");
    if (!work) {
      table.innerHTML = "Work not found";
    } else {
      table.innerHTML = "";
      let row = document.createElement("tr");
      table.appendChild(row);

      let title = document.createElement("h2");
      row.appendChild(title);
      title.appendChild(document.createTextNode(work.Title));
    
      let text = document.createElement("p");
      row.appendChild(text);
      text.innerHTML = work.Text.replace(/\r/g, '').replace(/\n/g, '<br>');

      let matchNumber = 0;
      let found = false;
      let index = 0;
      let lastElement;
      let regex = new RegExp(query, 'i');
      while (!found) {
        let node = text.childNodes[index];
        if (node instanceof Element) {
          lastElement = node;
        }
        let content = node.textContent || node.innerText;
        if (content && regex.test(content)) {
          if (matchNumber == matchId) {
            found = true;
            console.log("found! " + content);
            lastElement.scrollIntoView();
          } else {
            matchNumber++;
          }
        }
        index++;
      }
    }
  }
};

const form = document.getElementById("form");
form.addEventListener("submit", Controller.search);
