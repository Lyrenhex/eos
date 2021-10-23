// set the version number
const VERSION = "3.1.0";

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
    }).then(function (reg) {
      console.log('Service worker has been registered for scope:' + reg.scope);
    });
  }
}

class UserData {
  #name;
  #positives;
  #negatives;
  #neutrals;
  #moods;
  #theme;
  constructor() {
    this.storage = new StorageSpace((usage) => { update_var('usage', Math.round(usage / 1024)); });

    this.#name = this.storage.getItem('name');
    this.#positives = JSON.parse(this.storage.getItem('positives'));
    this.#negatives = JSON.parse(this.storage.getItem('negatives'));
    this.#neutrals = JSON.parse(this.storage.getItem('neutrals'));
    this.#moods = JSON.parse(this.storage.getItem('moods'));
    this.#theme = this.storage.getItem('theme');

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
    if (this.#theme == null) this.#theme = "";

    for (var year in this.#moods.years) {
      if (this.#moods.years[year].colour == undefined) {
        this.#moods.years[year].colour = generateNewColour();

        this.storage.setItem('moods', JSON.stringify(this.#moods));
      }
    }
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
  get theme() {
    return this.#theme;
  }

  setName(string) {
    if (string === "") this.#name = null;
    else this.#name = string;
    while (!this.storage.setItem('name', this.#name)) {
      if (!this.freeSpace()) {
        break;
      }
    }
    refresh_name();
  }

  addPositive(string) {
    if (string === "") return;
    this.#positives.push(string);
    while (!this.storage.setItem('positives', JSON.stringify(this.#positives))) {
      if (!this.freeSpace()) {
        break;
      }
    }
    refresh_comments();
  }
  addNegative(string) {
    if (string === "") return;
    this.#negatives.push(string);
    while (!this.storage.setItem('negatives', JSON.stringify(this.#negatives))) {
      if (!this.freeSpace()) {
        break;
      }
    }
    refresh_comments();
  }
  addNeutral(string) {
    if (string === "") return;
    this.#neutrals.push(string);
    while (!this.storage.setItem('neutrals', JSON.stringify(this.#neutrals))) {
      if (!this.freeSpace()) {
        break;
      }
    }
    refresh_comments();
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
        colour: generateNewColour(),
        month: [{ mood: 0, num: 0 }, { mood: 0, num: 0 }, { mood: 0, num: 0 },
        { mood: 0, num: 0 }, { mood: 0, num: 0 }, { mood: 0, num: 0 },
        { mood: 0, num: 0 }, { mood: 0, num: 0 }, { mood: 0, num: 0 },
        { mood: 0, num: 0 }, { mood: 0, num: 0 }, { mood: 0, num: 0 }]
      }
      newYear.month[month].mood += mood;
      newYear.month[month].num++;
      this.#moods.years.push(newYear);
    }

    while (!this.storage.setItem('moods', JSON.stringify(this.#moods))) {
      if (!this.freeSpace()) {
        break;
      }
    }
    refresh_graphs();
  }

  setTheme() {
    this.#theme = document.getElementById("theme_selector").value;
    while (!this.storage.setItem('theme', this.#theme)) {
      if (!this.freeSpace()) {
        break;
      }
    }
    refresh_theme();
  }

  freeSpace() {
    if (this.#negatives.length > 0) {
      this.#negatives.shift();
      this.storage.setItem('negatives', JSON.stringify(this.#negatives));
    } else if (this.#neutrals.length > 0) {
      this.#neutrals.shift();
      this.storage.setItem('neutrals', JSON.stringify(this.#neutrals));
    } else if (this.#positives.length > 0) {
      this.#positives.shift();
      this.storage.setItem('positives', JSON.stringify(this.#positives));
    } else if (this.#moods.years.length > 1) {
      this.#moods.years.shift();
    } else {
      console.error("No space could be made by deleting old comments or dropping data from previous years.");
      return false;
    }
    return true;
  }
}

var data = new UserData();

var graphs;

var MOOD;

window.onresize = function () {
  if (window.innerWidth >= 800) {
    show('menu');
  }
};

window.onload = () => {
  update_var('version_number', `${VERSION}`);
  update_var('usage', Math.round(data.storage.usage() / 1024));
  update_var('capacity', Math.round(data.storage.capacity() / 1024));
  window.onresize();
  createGraphs();
  refresh();
}

function createGraphs() {
  let dayGraph = createDayGraph();
  let monthGraph = createMonthGraph();
  let yearGraph = createYearGraph();
  graphs = {
    day: dayGraph,
    month: monthGraph,
    year: yearGraph
  };
}

function exportData() {
  let dataObj = {
    name: data.name,
    positives: data.positives,
    negatives: data.negatives,
    neutrals: data.neutrals,
    moods: data.moods,
    theme: data.theme,
  };

  var file = new Blob([JSON.stringify(dataObj)], { type: "application/json" });

  // credit: thank you to 0x000f at https://stackoverflow.com/questions/1066452/easiest-way-to-open-a-download-window-without-navigating-away-from-the-page
  var a = document.createElement('a');
  a.href = URL.createObjectURL(file);
  a.download = "eos_data.json";
  a.target = "_blank";
  document.body.appendChild(a);
  a.click();
  document.body.removeChild(a);
}

function refresh_name() {
  if (data.name !== null) {
    update_var('name', `, ${data.name}`);
  } else {
    update_var('name', '');
  }
}

function importData() {
  if (document.querySelector("#file__import").files.length == 0) {
    return;
  }

  // first file selected by user
  var file = document.querySelector("#file__import").files[0];

  // read the file
  var reader = new FileReader();

  // file reading finished successfully
  reader.addEventListener('load', function (e) {
    // contents of file in variable     
    var text = e.target.result;

    let dataObj = JSON.parse(text);

    var size = 0;
    for (var key in dataObj) {
      if (dataObj.hasOwnProperty(key)) {
        size += (key + dataObj[key]).length * 2;
      }
    }

    if (size > data.storage.capacity()) {
      console.error(`Imported data is too large for storage: ${size} B > ${data.storage.capacity()} B`);
      return;
    }

    data.storage.clear();
    if (dataObj.name != undefined) data.storage.setItem('name', dataObj.name);
    if (dataObj.positives != undefined) data.storage.setItem('positives', JSON.stringify(dataObj.positives));
    if (dataObj.negatives != undefined) data.storage.setItem('negatives', JSON.stringify(dataObj.negatives));
    if (dataObj.neutrals != undefined) data.storage.setItem('neutrals', JSON.stringify(dataObj.neutrals));
    if (dataObj.moods != undefined) data.storage.setItem('moods', JSON.stringify(dataObj.moods));
    if (dataObj.theme != undefined) data.storage.setItem('theme', dataObj.theme);

    data = new UserData();
    refresh();

    console.log("Imported .json: ", data);
  });

  // file reading failed
  reader.addEventListener('error', function (event) {
    console.error('Failed reading text file:', event);
  });

  // read as text file
  reader.readAsText(file);
}

function refresh_comments() {
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
    if (data !== "") {
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
    if (data !== "") {
      let li = document.createElement('li');
      li.classList.add('spectral');
      let li_text = document.createTextNode(data);
      li.appendChild(li_text);
      ul.appendChild(li);
    }
  });
}

function refresh_graphs() {
  graphs.day.update_day(data.moods);
  graphs.month.update_month(data.moods);
  graphs.year.update_year(data.moods);
}

function refresh_theme() {
  document.getElementsByName('theme').forEach((link, i) => {
    link.title = '';
    link.disabled = true;
  });

  if (data.theme !== "") {
    document.getElementById(`theme.${data.theme}`).disabled = false;
    document.getElementById(`theme.${data.theme}`).title = "main";
  }
}

function refresh() {
  refresh_name();
  refresh_comments();
  refresh_graphs();
  refresh_theme();
}

function mood(mood) {
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
  section('tracker');
  toggle(`mood__${MOOD}`);
  undone('block__mood');
}

function monthNext() {
  var active = document.getElementsByClassName('activeYear')[0];
  var index = YEARS.indexOf(active.id);
  if (index - 1 < 0) {
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
  if (index + 1 >= YEARS.length) {
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
  data.storage.clear();
  data = new UserData();
  refresh();
}
