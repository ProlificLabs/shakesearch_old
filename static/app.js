const Controller = {
  search: (ev) => {
    // remove previous search results
    const searchList = document.getElementById("search-list");
    searchList.innerHTML = "";
    const statusInfo = document.getElementById("status-info");
    statusInfo.innerHTML = "";

    // add spinner
    const spinner = document.createElement("div");
    spinner.setAttribute("id", "spinner");
    spinner.classList.add("spinner");
    searchList.prepend(spinner);

    ev.preventDefault();
    const form = document.getElementById("form");
    const data = Object.fromEntries(new FormData(form));
    const response = fetch(`/search?q=${data.query}`).then((response) => {
      if (response.status === 400) {
        // display error message
        const statusInfo = document.getElementById("status-info");
        const errorMessage = document.createElement("h3");
        errorMessage.innerHTML = `<span style="color: #FD5F00;">Your search did not match anything. </span> <br>
          <br> Suggestions: 
           <br> Make sure all words are spelled correctly. 
           <br> Try different keywords.
           <br> Try more general keywords.
           <br> Try fewer keywords.`;
        statusInfo.prepend(errorMessage);
      }
      response.json().then((results) => {
        Controller.updateTable(results, data.query);
        try {
          spinner.parentElement.removeChild(spinner);
        } catch (err) {
          // do nothing because that means spinner was already removed
        }
      })
      .catch(error => {
        if (error instanceof SyntaxError) {
          // do nothing because this happens when the query has no matches
        } 
        else {
          const statusInfo = document.getElementById("status-info");
          statusInfo.innerHTML = "";
          const errorMessage = document.createElement("h3");
          errorMessage.innerHTML = `<span style="color: #FD5F00;">Something went wrong. Try again later... </span>`;
          statusInfo.prepend(errorMessage);
          console.log(error);
        }
        try {
          spinner.parentElement.removeChild(spinner);
        } catch (err) {
          // do nothing because that means spinner was already removed
        }
      });
    });
  }, 

  updateTable: (results, query) => {
    const searchList = document.getElementById("search-list");
    searchList.innerHTML = ''
    for (let result of results) {
      const paragraph = result.Paragraph;

      
      let text = formatParagraph(paragraph);
      const firstLine = paragraph[0];
      const lastLine = paragraph[paragraph.length - 1];

      searchList.insertAdjacentHTML( 'beforeend', `<li id=${firstLine.LineIndex + "-" + lastLine.LineIndex}><div class="card">
          <pre>${text}  <pre/>
        <div/> <li/>`);

      console.log("searchList: ");
      console.log(searchList);
      const preList = searchList.getElementsByTagName('pre')

      // delete pre element that got added by joining <li> elements
      const redundantPre = preList[preList.length - 1];
      redundantPre.parentNode.removeChild(redundantPre);
    }

    const cards = document.getElementsByClassName("card");
    for (let i = 0; i < cards.length; i++) {

      let button1 = document.createElement("button");
      let button2 = document.createElement("button");

      const card = cards[i]
      card.insertBefore(button1, card.firstChild)
      card.appendChild(button2)


      if (cards[i].parentElement.id.split("-")[0] !== "0") {
        button1.innerHTML = "...";
        button1.classList.add("add-lines");
        button1.addEventListener("click", function(event) {
          console.log("button1 is clicked");
          handleAddLinesClick(event, true);
        });
      } else {
        // delete button1
        button1.parentNode.removeChild(button1);
      }

      
      if (parseInt(cards[i].parentElement.id.split("-")[1]) < 169432) {
        button2.innerHTML = "...";
        button2.classList.add("add-lines");
        button2.addEventListener("click", function(event) {
          console.log("button2 is clicked");
          handleAddLinesClick(event, false);
        });
      } else {
        // delete button1
        button2.parentNode.removeChild(button2);
      }

      // add link to read entire work
      const lineNumbers = card.parentElement.id.split("-");
      const middleLine = Math.ceil((lineNumbers[0] + lineNumbers[1]) / 2)
      const openWork = `<div class="open-work-wrapper"><button class="open-work">Open Book</button></div>`;
      card.insertAdjacentHTML("beforeend", openWork)

      const openBookButtons = card.getElementsByClassName("open-work");
      for (let i = 0; i < openBookButtons.length; i++) {
        openBookButtons[i].addEventListener("click", (event) => {
          console.log("clicked Open book button")
          handleOpenWorkClick(event, query);
        })
      }
    }
  },
};

const handleOpenWorkClick = (event, query) => {
  const startAndEndId = event.target.parentElement.parentElement.parentElement.id.split("-");

  const middleLine = Math.ceil((parseInt(startAndEndId[0]) + parseInt(startAndEndId[1])) / 2);
  const url = `/read-work?q=${query}&line=${middleLine}`;

  fetch(url).then((response) => {
    response
    .json()
    .then((results) => {
      Controller.updateTable(results, query);
    });
  })
}

const formatParagraph = (paragraph) => {
  let text = ""
  for (let i = 0; i < paragraph.length; i++) {
    text += paragraph[i].TextResult + "\n"
  }
  return text
}

const form = document.getElementById("form");
form.addEventListener("submit", Controller.search);


const handleAddLinesClick = (event, addLinesUp) => {
  const wrapper = event.target.parentElement.parentElement
  wrapperId = wrapper.id;
  const lineIndexes = wrapperId.split("-");
  const url = addLinesUp ? `/add-lines/up?q=${lineIndexes[0]}` : `/add-lines/down?q=${lineIndexes[1]}`;

  fetch(url).then((response) => {
    response
    .json()
    .then((result) => {
      addExtraLines(wrapper, event, addLinesUp, result);
    });
  });
  
}

const addExtraLines = (wrapper, event, addLinesUp, result) => {
  const card = event.target.parentElement
  const pre = card.getElementsByTagName("pre")[0]

  if (addLinesUp) {
    result = result.reverse();
  } 
  const text = formatParagraph(result);
  console.log("text: " + text);
  if (addLinesUp) {
    // Change text to add new lines
    pre.innerHTML = text + pre.innerHTML;
    
    // Add indexes of lines of that card
    let newId = wrapper.id.split("-");
    newId[0] = result[0].LineIndex;
    wrapper.id = newId.join("-");
    if (parseInt(newId[0]) == 0) {
      // remove button
      const button = card.getElementsByTagName("button")[0];
      button.parentElement.removeChild(button);
    }
  } else {
    pre.innerHTML = pre.innerHTML + text;
    let newId = wrapper.id.split("-");
    newId[1] = result[result.length-1].LineIndex;
    wrapper.id = newId.join("-");

    if (parseInt(newId[1]) == 169432) {
      // remove button 
      const button = card.getElementsByTagName("button")[1];
      button.parentElement.removeChild(button);
    }
  }
}