/*
COPYRIGHT (C) Damian Heaton 2017 - 2018

The storage of this software on a computer via means of browser 'caching', and the execution of such code by user browsers is permitted. This work cannot be duplicated, copied, distributed, or modified (neither privately nor publicly) without the express, written consent of Damian Heaton, whom can be contacted (at time of notice) at damian@lyrenhex.me. This software cannot be used for commercial purposes.

THIS SOFTWARE IS DISTRIBUTED "AS IS", WITHOUT WARRANTY OF ANY FORM (EITHER EXPRESS OR IMPLIED), INCLUDING (BUT NOT LIMITED TO) ANY IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE.

TL;DR:
You may:
- Store this software on a computer via a browser's cache, though not by any other means.
- Store this software for purposes of search indexing.
- Execute this code for the purpose of using the software.

You may NOT:
- Distribute, copy, replicate, or duplicate this copyrighted software without express, written permission from Damian Heaton.
- Use this software, in any way, shape or form, for any commercial purpose(s).
- Claim this software as your own, or attempt to imply affiliation with the software in any way that could be detrimental or unlawful, or to suggest that the software, or Damian Heaton, are represented by, or represent, yourself.
*/

Date.prototype.getWeek = function() { // get the current week of the year (start Sun)
  var onejan = new Date(this.getFullYear(),0,1);
  var millisecsInDay = 86400000;
  return Math.ceil((((this - onejan) /millisecsInDay) + onejan.getDay()+1)/7);
};

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
function err(error) {
  let edump = document.getElementById('error_dump');
  edump.textContent = JSON.stringify(error, null, 4);
  show('block__error');
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
function signout() {
  location.reload();
}
function section(id) {
  let activeSection = document.getElementsByClassName("section active")[0];
  let newActiveSection = document.getElementById(`section__${id}`);
  activeSection.classList.remove('active');
  newActiveSection.classList.add('active');
}

function dayGraph(ctx, data) {
  var moods = [0, 0, 0, 0, 0, 0, 0];
  for(day in data.Day){
    moods[day] = (data.Day[day].Mood / data.Day[day].Num);
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
function fortGraph(ctx, data) {
  var date = new Date();
  var moods = [[0, 0, 0, 0, 0, 0, 0], [0, 0, 0, 0, 0, 0, 0]];
  for(day in data.Day){
    if(data.Day[day].weeks !== undefined){
      if(data.Day[day].weeks[date.getWeek()] !== undefined){
        moods[0][day] = (data.Day[day].weeks[date.getWeek()].Mood / data.Day[day].weeks[date.getWeek()].Num);
      }
      if(data.Day[day].weeks[date.getWeek() - 1] !== undefined){
        moods[1][day] = (data.Day[day].weeks[date.getWeek() - 1].Mood / data.Day[day].weeks[date.getWeek() - 1].Num);
      }
    }
  }
  chart = new Chart(ctx, {
    type: 'line',
    data: {
      labels: ['Sun', 'Mon', 'Tues', 'Wed', 'Thurs', 'Fri', 'Sat'],
      datasets: [
        {
          label: "This week",
          data: moods[0],
          borderColor: 'rgba(0, 192, 255, 0.8)',
          backgroundColor: 'rgba(0, 192, 255, 0.4)'
        },
        {
          label: "Last week",
          data: moods[1]
        }
      ]
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
  for(month in data.Month){
    moods[month] = (data.Month[month].Mood / data.Month[month].Num);
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
  for(year in data.Years) {
    if (year.Year !== 0) {
      if (moods[data.Years[year].Year] === undefined) {
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
