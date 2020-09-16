Date.prototype.getWeek = function() { // get the current week of the year (start Sun)
  var onejan = new Date(this.getFullYear(),0,1);
  var millisecsInDay = 86400000;
  return Math.ceil((((this - onejan) /millisecsInDay) + onejan.getDay()+1)/7);
};

function escapeHtml(str) {
  var div = document.createElement('div');
  div.appendChild(document.createTextNode(str));
  return div.innerHTML;
}

function done(blockId) {
  let block = document.getElementById(blockId);
  block.classList.add('done');
}
function undone(blockId) {
  let block = document.getElementById(blockId);
  block.classList.remove('done');
}
function show(blockId) {
  let block = document.getElementById(blockId);
  block.classList.add('shown');
}
function toggle(id) {
  let block = document.getElementById(id);
  block.classList.toggle('shown');
}
function update_var(varName, text){
  let spans = document.getElementsByName(varName);
  for (element in spans) {
    element = spans[element];
    element.textContent = text;
  }
}
function section(id) {
  let activeSection = document.getElementsByClassName("section active")[0];
  let newActiveSection = document.getElementById(`section__${id}`);
  activeSection.classList.remove('active');
  newActiveSection.classList.add('active');
}

function dayGraph(ctx, data) {
  var moods = [0, 0, 0, 0, 0, 0, 0];
  for(day in data.day){
    moods[day] = (data.day[day].mood / data.day[day].num);
  }
  chart = new Chart(ctx, {
    type: 'line',
    data: {
      labels: ['Sun', 'Mon', 'Tues', 'Wed', 'Thurs', 'Fri', 'Sat'],
      datasets: [{
        label: "Average mood",
        data: moods
      }]
    },
    options: {
      responsive: true,
      tooltips: {
        mode: 'index',
        intersect: false
      },
      hover: {
        mode: 'nearest',
        intersect: false
      },
      scales: {
        xAxes: [{
          display: true,
          scaleLabel: {
            display: true,
            labelString: 'Day of week'
          }
        }],
        yAxes: [{
          display: true,
          scaleLabel: {
            display: true,
            labelString: 'Average mood'
          }
        }]
      }
    }
  });
}
function monthGraph(ctx, data){
  var moods = [0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0];
  for(month in data.month){
    moods[month] = (data.month[month].mood / data.month[month].num);
  }
  let chart = new Chart(ctx, {
    type: 'line',
    data: {
      labels: ['Jan', 'Feb', 'Mar', 'Apr', 'May', 'Jun', 'Jul', 'Aug', 'Sep', 'Oct', 'Nov', 'Dec'],
      datasets: [{
        label: "Average mood",
        data: moods
      }]
    },
    options: {
      responsive: true,
      tooltips: {
        mode: 'index',
        intersect: false
      },
      hover: {
        mode: 'nearest',
        intersect: false
      },
      scales: {
        xAxes: [{
          display: true,
          scaleLabel: {
            display: true,
            labelString: 'Month'
          }
        }],
        yAxes: [{
          display: true,
          scaleLabel: {
            display: true,
            labelString: 'Average mood'
          }
        }]
      }
    }
  });
}
function yearGraph(data){
  var moods = {};
  for(year in data.years) {
    if (year.year !== 0) {
      if (moods[data.years[year].year] === undefined) {
        moods[data.Years[year].Year] = [0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0];
      }
      data.Years[year].Month.forEach((month, i) => {
        moods[data.Years[year].Year][i] = month.Mood / month.Num;
      })
    }
  }

  var graphs = {};
  for (year in moods) {
    let ctx = document.createElement('canvas');
    ctx.id = `graph.${year}`;
    let chart = new Chart(ctx, {
      type: 'line',
      data: {
        labels: ['Jan', 'Feb', 'Mar', 'Apr', 'May', 'Jun', 'Jul', 'Aug', 'Sep', 'Oct', 'Nov', 'Dec'],
        datasets: [{
          label: `${year}`,
          data: moods[year]
        }]
      },
      options: {
        responsive: true,
        tooltips: {
          mode: 'index',
          intersect: false
        },
        hover: {
          mode: 'nearest',
          intersect: false
        },
        scales: {
          xAxes: [{
            display: true,
            scaleLabel: {
              display: true,
              labelString: 'Month'
            }
          }],
          yAxes: [{
            display: true,
            scaleLabel: {
              display: true,
              labelString: 'Average mood'
            }
          }]
        }
      }
    });
    graphs[year] = ctx;
  }
  return graphs;
}
