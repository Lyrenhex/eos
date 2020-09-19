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

class Graph {
  constructor(trackChart, reportChart) {
    this.tracker = trackChart;
    this.report = reportChart;
  }

  #update(moods) {
    this.tracker.data.datasets = moods;
    this.report.data.datasets = moods;
    this.tracker.update();
    this.report.update();
  }

  update_day(data) {
    var moods = [0, 0, 0, 0, 0, 0, 0];
    console.log(data);
    for (var day in data.day) {
      moods[day] = (data.day[day].mood / data.day[day].num);
    }
    this.#update([{
      label: "Average mood",
      data: moods
    }]);
  }

  update_month(data) {
    var moods = [0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0];
    for (var month in data.month) {
      moods[month] = (data.month[month].mood / data.month[month].num);
    }
    this.#update([{
      label: "Average mood",
      data: moods
    }]);
  }

  update_year(data) {
    let datasets = [];

    let year;
    for (year in data.years) {
      if (data.years[year].year !== 0) {
        let dataset = {
          label: `${data.years[year].year}`,
          data: [0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0]
        };
        data.years[year].month.forEach((month, i) => {
          dataset.data[i] = month.mood / month.num;
        });
        datasets.push(dataset);
      }
    }

    this.#update(datasets);
  }
}

function createDayGraph() {
  var tracker = document.getElementById('moodchart_days');
  var report = document.getElementById('pr_moodchart_days');
  var chart_tracker = new Chart(tracker, {
    type: 'line',
    data: {
      labels: ['Sun', 'Mon', 'Tues', 'Wed', 'Thurs', 'Fri', 'Sat'],
      datasets: []
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
  var chart_report = new Chart(report, {
    type: 'line',
    data: {
      labels: ['Sun', 'Mon', 'Tues', 'Wed', 'Thurs', 'Fri', 'Sat'],
      datasets: []
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
  return new Graph(chart_tracker, chart_report);
}

function createMonthGraph() {
  var tracker = document.getElementById('moodchart_months');
  var report = document.getElementById('pr_moodchart_months');
  var chart_tracker = new Chart(tracker, {
    type: 'line',
    data: {
      labels: ['Jan', 'Feb', 'Mar', 'Apr', 'May', 'Jun', 'Jul', 'Aug', 'Sep', 'Oct', 'Nov', 'Dec'],
      datasets: []
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
  var chart_report = new Chart(report, {
    type: 'line',
    data: {
      labels: ['Jan', 'Feb', 'Mar', 'Apr', 'May', 'Jun', 'Jul', 'Aug', 'Sep', 'Oct', 'Nov', 'Dec'],
      datasets: []
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
  return new Graph(chart_tracker, chart_report);
}

function createYearGraph() {
  var tracker = document.getElementById('moodchart_years');
  var report = document.getElementById('pr_moodchart_years');
  var chart_tracker = new Chart(tracker, {
    type: 'line',
    data: {
      labels: ['Jan', 'Feb', 'Mar', 'Apr', 'May', 'Jun', 'Jul', 'Aug', 'Sep', 'Oct', 'Nov', 'Dec'],
      datasets: []
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
  var chart_report = new Chart(report, {
    type: 'line',
    data: {
      labels: ['Jan', 'Feb', 'Mar', 'Apr', 'May', 'Jun', 'Jul', 'Aug', 'Sep', 'Oct', 'Nov', 'Dec'],
      datasets: []
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
  return new Graph(chart_tracker, chart_report);
}
