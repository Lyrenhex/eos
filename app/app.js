/*
COPYRIGHT (C) Damian Heaton 2017 - 2020

The storage of this file on a computer via means of browser 'caching', and the execution of such code by user browsers is permitted. This work cannot be duplicated, copied, distributed, or modified (neither privately nor publicly) without the express, written consent of Damian Heaton, whom can be contacted (at time of notice) at damian@lyrenhex.me. This software cannot be used for commercial purposes.

THIS SOFTWARE IS DISTRIBUTED "AS IS", WITHOUT WARRANTY OF ANY FORM (EITHER EXPRESS OR IMPLIED), INCLUDING (BUT NOT LIMITED TO) ANY IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE.

TL;DR:
You may:
- Store this file on a computer via a browser's cache, though not by any other means.
- Store this file for purposes of search indexing.
- Execute this code for the purpose of using the software.

You may NOT:
- Distribute, copy, replicate, or duplicate this copyrighted software without express, written permission from Damian Heaton.
- Use this software, in any way, shape or form, for any commercial purpose(s).
- Claim this software as your own, or attempt to imply affiliation with the software in any way that could be detrimental or unlawful, or to suggest that the software, or Damian Heaton, are represented by, or represent, yourself.
*/

// set the version number
const VERSION = "3.0.0";

// make sure the client *allows* service workers...
// apparently ChromeOS doesn't?
if (navigator.serviceWorker != undefined) {
  // register the service worker if not active
  if (navigator.serviceWorker.controller) {
    console.log('active service worker found, no need to register')
  } else {
    //Register the ServiceWorker
    navigator.serviceWorker.register('sw.js', {
      scope: './'
    }).then(function(reg) {
      console.log('Service worker has been registered for scope:'+ reg.scope);
    });
  }
}

class UserData {
  #name;
  #positives;
  #negatives;
  #neutrals;
  #moods;
  constructor (storage) {
    this.storage = storage;

    this.#name = this.storage.getItem('name');
    this.#positives = JSON.parse(this.storage.getItem('positives'));
    this.#negatives = JSON.parse(this.storage.getItem('negatives'));
    this.#neutrals = JSON.parse(this.storage.getItem('neutrals'));
    this.#moods = JSON.parse(this.storage.getItem('moods'));

    if (this.#positives == null) this.#positives = [];
    if (this.#negatives == null) this.#negatives = [];
    if (this.#neutrals == null) this.#neutrals = [];
    if (this.#moods == null) this.#moods = {
      day: [{ mood: 0, num: 0 }, { mood: 0, num: 0 }, { mood: 0, num: 0 },
        { mood: 0, num: 0 }, { mood: 0, num: 0 }, { mood: 0, num: 0 },
        { mood: 0, num: 0 }],
      month: [{ mood: 0, num: 0 }, { mood: 0, num: 0 }, { mood: 0, num: 0 },
        { mood: 0, num: 0 }, { mood: 0, num: 0 }, { mood: 0, num: 0 },
        { mood: 0, num: 0 }, { mood: 0, num: 0 }, { mood: 0, num: 0 },
        { mood: 0, num: 0 }, { mood: 0, num: 0 }, { mood: 0, num: 0 }],
      years: []
    };
  }
  get name() {
    return this.#name;
  }
  get positives() {
    return this.#positives;
  }
  get negatives() {
    return this.#negatives;
  }
  get neutrals() {
    return this.#neutrals;
  }
  get moods() {
    return this.#moods;
  }

  setName(string) {
    if (string === "") this.#name = null;
    else this.#name = string;
    this.storage.setItem('name', this.#name);
    refresh();
  }

  addPositive(string) {
    if (string === "") return;
    this.#positives.push(string);
    this.storage.setItem('positives', JSON.stringify(this.#positives));
    refresh();
  }
  addNegative(string) {
    if (string === "") return;
    this.#negatives.push(string);
    this.storage.setItem('negatives', JSON.stringify(this.#negatives));
    refresh();
  }
  addNeutral(string) {
    if (string === "") return;
    this.#neutrals.push(string);
    this.storage.setItem('neutrals', JSON.stringify(this.#neutrals));
    refresh();
  }
  addMood(day, month, year, mood) {
    this.#moods.day[day].mood += mood;
    this.#moods.day[day].num++;
    this.#moods.month[month].mood += mood;
    this.#moods.month[month].num++;

    let yearRecorded = false;
    this.#moods.years.forEach((y, i) => {
      if (y.year == year) {
        this.#moods.years[i].month[month].mood += mood;
        this.#moods.years[i].month[month].num++;
        yearRecorded = true;
      }
    });

    if (!yearRecorded) {
      let newYear = {
        year: year,
        month: [{ mood: 0, num: 0 }, { mood: 0, num: 0 }, { mood: 0, num: 0 },
          { mood: 0, num: 0 }, { mood: 0, num: 0 }, { mood: 0, num: 0 },
          { mood: 0, num: 0 }, { mood: 0, num: 0 }, { mood: 0, num: 0 },
          { mood: 0, num: 0 }, { mood: 0, num: 0 }, { mood: 0, num: 0 }]
      }
      newYear.month[month].mood += mood;
      newYear.month[month].num++;
      this.#moods.years.push(newYear);
    }

    this.storage.setItem('moods', JSON.stringify(this.#moods));
    refresh();
  }
}

var storage = window.localStorage;
var data = new UserData(storage);

var YEARS = [];
var USER;
var DB;
var MOOD;
var MOOD_TIME_DATA = {
  nums: [],
  dates: [],
  months: []
}

window.onresize = function() {
  if (window.innerWidth >= 800) {
    show('menu');
  }
};

document.onreadystatechange = function() {
  window.onresize();
  update_var('version_number', `${VERSION}`);

  refresh();
}

function exportData() {
  let dataObj = {
    name: data.name,
    positives: data.positives,
    negatives: data.negatives,
    neutrals: data.neutrals,
    moods: data.moods
  };

  var file = new Blob([JSON.stringify(dataObj)], {type: "application/json"});

  // credit: thank you to 0x000f at https://stackoverflow.com/questions/1066452/easiest-way-to-open-a-download-window-without-navigating-away-from-the-page
  var a = document.createElement('a');
  a.href = URL.createObjectURL(file);
  a.download = "eos_data.json";
  a.target = "_new";
  document.body.appendChild(a);
  a.click();
  document.body.removeChild(a);
}

function refresh() {
  if(data.name !== null){
    update_var('name', `, ${data.name}`);
  }

  let ul = document.getElementById('comments_positive');
  while (ul.firstChild) {
    ul.removeChild(ul.lastChild);
  }
  let ul2 = document.getElementById('comments');
  while (ul2.firstChild) {
    ul2.removeChild(ul2.lastChild);
  }
  data.positives.forEach((data, key) => {
    let li = document.createElement('li');
    li.classList.add('spectral');
    let li_text = document.createTextNode(data);
    li.appendChild(li_text);
    ul.appendChild(li);

    let li2 = document.createElement('li');
    li2.classList.add('spectral');
    let li_text2 = document.createTextNode(data);
    li2.appendChild(li_text2);
    ul2.appendChild(li2);
  });
  ul = document.getElementById('comments_neutral');
  while (ul.firstChild) {
    ul.removeChild(ul.lastChild);
  }
  data.neutrals.forEach((data, key) => {
    if(data !== ""){
      let li = document.createElement('li');
      li.classList.add('spectral');
      let li_text = document.createTextNode(data);
      li.appendChild(li_text);
      ul.appendChild(li);
    }
  });
  ul = document.getElementById('comments_negative');
  while (ul.firstChild) {
    ul.removeChild(ul.lastChild);
  }
  data.negatives.forEach((data, key) => {
    if(data !== ""){
      let li = document.createElement('li');
      li.classList.add('spectral');
      let li_text = document.createTextNode(data);
      li.appendChild(li_text);
      ul.appendChild(li);
    }
  });

  var ctx = document.getElementById('moodchart_months');
  monthGraph(ctx, data.moods);
  var ctx = document.getElementById('moodchart_days');
  dayGraph(ctx, data.moods);

  var ctx = document.getElementById('pr_moodchart_months');
  monthGraph(ctx, data.moods);
  var ctx = document.getElementById('pr_moodchart_days');
  dayGraph(ctx, data.moods);

  //var ctx = document.getElementById('moodchart_days_fortnight');
  //fortGraph(ctx, msg.user.Moods);

  /* var years = yearGraph(data.moods);
  for(year in years){
    YEARS.push(year);
    year = years[year];
    var tracker = document.getElementById('annual_moods_graphs');
    tracker.appendChild(year);
  }
  for(year in years){
    YEARS.push(year);
    year = years[year];
    var tracker = document.getElementById('pr_annual_moods_graphs');
    tracker.appendChild(year);
  }
  document.getElementById(`graph.${YEARS[YEARS.length-1]}`).classList.add('activeYear'); */
}

function mood(mood){
  let date = new Date();
  let month = date.getMonth();
  let day = date.getDay();
  let year = date.getUTCFullYear();
  data.addMood(day, month, year, mood);
  done('block__mood');
  show(`mood__${mood}`);
  MOOD = mood;
}

function ecstatic_submit() {
  let comment = document.getElementById('ecstatic_comment').value;
  data.addPositive(comment);
  mood_continue();
}
function happy_submit() {
  let comment = document.getElementById('happy_comment').value;
  data.addPositive(comment);
  mood_continue();
}
function neutral_submit() {
  let comment = document.getElementById('neutral_comment').value;
  data.addNeutral(comment);
  mood_continue();
}
function negative_submit() {
  let comment = document.getElementById('negative_comment').value;
  data.addNegative(comment);
  mood_continue();
}
function danger_submit() {
  let comment = document.getElementById('danger_comment').value;
  data.addNegative(comment);
  mood_continue();
}
function mood_continue() {
  done(`mood__${MOOD}`);
  section('tracker');
}

function monthNext() {
  var active = document.getElementsByClassName('activeYear')[0];
  var index = YEARS.indexOf(active.id);
  if (index-1 < 0){
    index = YEARS.length - 1;
  } else {
    index--;
  }
  active.classList.remove('activeYear');
  active = document.getElementById(`graph.${YEARS[index]}`);
  active.classList.add('activeYear');
}
function monthPrev() {
  var active = document.getElementsByClassName('activeYear')[0];
  var index = YEARS.indexOf(active.id);
  if (index+1 >= YEARS.length){
    index = 0;
  } else {
    index++;
  }
  active.classList.remove('activeYear');
  active = document.getElementById(`graph.${YEARS[index]}`);
  active.classList.add('activeYear');
}

function change_name() {
  let name = document.getElementById('account_name').value;
  data.setName(name);
  return false;
}
function deleteData() {
  storage.clear();
  data = new UserData(storage);
  refresh();
}