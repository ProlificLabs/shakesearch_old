const Controller = {
  search: (ev) => {
    ev.preventDefault();
    const form = document.getElementById("form");
    const data = Object.fromEntries(new FormData(form));
    const response = fetch(`/search?q=${data.query}&f=${data.filter}&d=${data.dataset}&l=${data.limit}`).then((response) => {
      response.json().then((results) => {
        Controller.updateTable(results);
      });
    });
  },

  updateTable: (results) => {
       // Define variables
       var tabs = $('.tabs');
       var container = $('.container');
       var html_tabs = '';
       var html_content = '';
    console.log(results)
  var index = 0
    for (x in results) {
      html_tabs +='<li><a href="#tab'+index+'">'+x+'</a></li>';

      for (y in results[x]) {
        for (z in results[x][y])
        html_content +='<div id="tab'+index+'"><p>'+results[x][y][z]+'</p></div>';
      }

      index += index + 1;
    }
    tabs.html(html_tabs);
    container.html(html_content);

    // Set tabs and content html
    tabs.html(html_tabs);
    container.html(html_content);
    // Looping links
    $.each($('.tabs li a'),function(count, item) {
       // Set on click handler
       $(this).on('click',function() {
          // Hide all div content
          container.find('div').removeClass('active');
          var current = $(this).attr('href');
          // Remove active class on links
          $('.tabs li a').removeClass('active');
          // Set active class on current and ul parent
          $(this).addClass('active');
          $(this).parents('ul').find('li').removeClass('active');
          $(this).parent().addClass('active');
          // Show current container
          $(current).addClass('active');
       });  
    }).eq(1).click().addClass('active');    
  },
};

const form = document.getElementById("form");
form.addEventListener("submit", Controller.search);
