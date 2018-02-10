/*
COPYRIGHT (C) Damian Heaton 2017

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

var serverAddr = `ws://${window.location.hostname}:9874`;


var YEARS = [];

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

//update_var('version_number', `version ${version}`);

var USER;
var DB;
var MOOD;
var MOOD_TIME_DATA = {
  nums: [],
  dates: [],
  months: []
}

var sock = new WebSocket(serverAddr);
function signin_email() {
  var email = document.getElementById('login_email').value;
  var passw = document.getElementById('login_password').value;

  sock.send(JSON.stringify({'type':'login', 'emailAddress':email, 'password': passw}))

  return false;
}
sock.onerror = function(e) {
  console.error(e);
  // TODO: onerror
}
sock.onclose = function(e) {
  // TODO: onclose
}
sock.onmessage = function(e) {
  let msg = JSON.parse(e.data);
  switch(msg.type){
    case 'version':
      done('text__loading');
      show('block__login');
      update_var('version_number', `version ${msg.data}`)
      break;
    case 'login':
      if (msg.flag) {
        done('block__login');
        show('block__mood');
        document.getElementById('btn__menu').classList.add('loggedin');
        console.log(msg)
        if(msg.user.Name !== ""){
          update_var('name', `, ${msg.user.Name}`);
        }
        var ctx = document.getElementById('moodchart_months');
        monthGraph(ctx, msg.user.Moods);
        var ctx = document.getElementById('moodchart_days');
        dayGraph(ctx, msg.user.Moods);
        //var ctx = document.getElementById('moodchart_days_fortnight');
        //fortGraph(ctx, msg.user.Moods);

        /* var years = yearGraph(msg.user.Moods);
        for(year in years){
          YEARS.push(year);
          year = years[year];
          var tracker = document.getElementById('annual_moods_graphs');
          tracker.appendChild(year);
        }
        document.getElementById(`graph.${YEARS[YEARS.length-1]}`).classList.add('active'); */
      }
  }
}

firebase.auth().onAuthStateChanged(function(user) {
  done('text__loading');
  if (user) {
    document.title = `Eos: How are you?`;
    if (!user.emailVerified) {
      user.sendEmailVerification();
    }
    if(user.displayName !== null){
      update_var('name', `, ${user.displayName}`);
      document.title = `Eos: How are you, ${user.displayName}?`;
    }
    done('block__login');
    done('block__login_email');
    show('block__mood');
    USER = user;
    DB = firebase.database();
    document.getElementById('buttons_system').classList.add('loggedin');
    var commentsRef = DB.ref(`/users/${user.uid}/positives`);
    commentsRef.on('child_added', function(data) {
      let li = document.createElement('li');
      li.classList.add('spectral');
      let li_text = document.createTextNode(data.val());
      li.appendChild(li_text);
      let ul = document.getElementById('comments');
      ul.appendChild(li);
    });
    var moodsRef = DB.ref(`/users/${user.uid}/moods`);
    moodsRef.once('value', function(data) {
      data = data.val();
      if(data !== null){
        var ctx = document.getElementById('moodchart_months');
        monthGraph(ctx, data);
        var ctx = document.getElementById('moodchart_days');
        dayGraph(ctx, data);
        var ctx = document.getElementById('moodchart_days_fortnight');
        fortGraph(ctx, data);

        var years = yearGraph(data);
        for(year in years){
          YEARS.push(year);
          year = years[year];
          var tracker = document.getElementById('annual_moods_graphs');
          tracker.appendChild(year);
        }
        document.getElementById(`graph.${YEARS[YEARS.length-1]}`).classList.add('active');
      }
    });
  } else {
    document.title = 'Eos Login';
    show('block__login');
  }
}, (error) => {
  err(error);
});

function mood(mood){
  let date = new Date();
  let month = date.getMonth();
  let day = date.getDay();
  let json = {
    type: 'mood',
    day: day,
    month: month,
    mood: mood
  }
  console.log(JSON.stringify(json))
  sock.send(JSON.stringify(json))
  done('block__mood');
  show(`mood__${mood}`);
  MOOD = mood;
}

function ecstatic_submit() {
  let comment = document.getElementById('ecstatic_comment').value;
  if (comment !== "") {
    let ref = DB.ref(`/users/${USER.uid}/positives`);
    let newRef = ref.push();
    newRef.set(comment);
  }
  mood_continue();
}
function happy_submit() {
  let comment = document.getElementById('happy_comment').value;
  if (comment !== "") {
    let ref = DB.ref(`/users/${USER.uid}/positives`);
    let newRef = ref.push();
    newRef.set(comment);
  }
  mood_continue();
}
function neutral_submit() {
  let comment = document.getElementById('neutral_comment').value;
  if (comment !== "") {
    let ref = DB.ref(`/users/${USER.uid}/neutral`);
    let newRef = ref.push();
    newRef.set(comment);
  }
  mood_continue();
}
function negative_submit() {
  let comment = document.getElementById('negative_comment').value;
  if (comment !== "") {
    let ref = DB.ref(`/users/${USER.uid}/negatives`);
    let newRef = ref.push();
    newRef.set(comment);
  }
  mood_continue();
}
function danger_submit() {
  let comment = document.getElementById('danger_comment').value;
  if (comment !== "") {
    let ref = DB.ref(`/users/${USER.uid}/negatives`);
    let newRef = ref.push();
    newRef.set(comment);
  }
  mood_continue();
}
function mood_continue() {
  done(`mood__${MOOD}`);
  show('block__chat');
}

function monthPrev() {
  var active = document.getElementsByClassName('active')[0];
  var index = YEARS.indexOf(active.id);
  if (index-1 < 0){
    index = YEARS.length - 1;
  } else {
    index--;
  }
  active.classList.remove('active');
  active = document.getElementById(`graph.${YEARS[index]}`);
  active.classList.add('active');
}
function monthNext() {
  var active = document.getElementsByClassName('active')[0];
  var index = YEARS.indexOf(active.id);
  if (index+1 >= YEARS.length){
    index = 0;
  } else {
    index++;
  }
  active.classList.remove('active');
  active = document.getElementById(`graph.${YEARS[index]}`);
  active.classList.add('active');
}
