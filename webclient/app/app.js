/*
COPYRIGHT (C) Damian Heaton 2017 - 2018

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

// EDIT THIS LINE AS NECESSARY. Usually, the server should operate on the same hostname as the web app, but the port may need changing depending on which port your server is configured to use.
var serverAddr = `${location.protocol == "https:" ? "wss:" : "ws:"}//${window.location.hostname}:9874`;


var CHATID = "";


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
function update_acc() {
  let newEmail = document.getElementById('account_email').value;
  let newPass = document.getElementById('account_password').value;
  let newName = document.getElementById('account_name').value;

  sock.send(JSON.stringify({'type':'details', 'emailAddress': newEmail, 'password': newPass, 'data': newName}));
}
sock.onerror = function(e) {
  show('block__error');
}
sock.onclose = function(e) {
  // TODO: onclose
}
sock.onmessage = function(e) {
  let msg = JSON.parse(e.data);
  console.log(msg);
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
        setState();
        window.onresize = setState;

        var file = new Blob([JSON.stringify(msg.user)], {type: "application/json"});
        var downloadLink = document.getElementById('downloadLink');
        downloadLink.href = URL.createObjectURL(file);

        document.getElementById('account_email').placeholder = msg.user.EmailAddr;
        document.getElementById('account_name').placeholder = msg.user.Name;
        

        if(msg.user.Name !== ""){
          update_var('name', `, ${msg.user.Name}`);
        }
        msg.user.Positives.forEach((data, key) => {
          if(data !== ""){
            let li = document.createElement('li');
            li.classList.add('spectral');
            let li_text = document.createTextNode(data);
            li.appendChild(li_text);
            let ul = document.getElementById('comments_positive');
            ul.appendChild(li);

            let li2 = document.createElement('li');
            li2.classList.add('spectral');
            let li_text2 = document.createTextNode(data);
            li2.appendChild(li_text2);
            let ul2 = document.getElementById('comments');
            ul2.appendChild(li2);
          }
        })
        msg.user.Neutrals.forEach((data, key) => {
          if(data !== ""){
            let li = document.createElement('li');
            li.classList.add('spectral');
            let li_text = document.createTextNode(data);
            li.appendChild(li_text);
            let ul = document.getElementById('comments_neutral');
            ul.appendChild(li);
          }
        })
        msg.user.Negatives.forEach((data, key) => {
          if(data !== ""){
            let li = document.createElement('li');
            li.classList.add('spectral');
            let li_text = document.createTextNode(data);
            li.appendChild(li_text);
            let ul = document.getElementById('comments_negative');
            ul.appendChild(li);
          }
        })

        var ctx = document.getElementById('moodchart_months');
        monthGraph(ctx, msg.user.Moods);
        var ctx = document.getElementById('moodchart_days');
        dayGraph(ctx, msg.user.Moods);

        var ctx = document.getElementById('pr_moodchart_months');
        monthGraph(ctx, msg.user.Moods);
        var ctx = document.getElementById('pr_moodchart_days');
        dayGraph(ctx, msg.user.Moods);

        //var ctx = document.getElementById('moodchart_days_fortnight');
        //fortGraph(ctx, msg.user.Moods);

        var years = yearGraph(msg.user.Moods);
        for(year in years){
          YEARS.push(year);
          year = years[year];
          var tracker = document.getElementById('annual_moods_graphs');
          tracker.appendChild(year);
        }
        var years = yearGraph(msg.user.Moods);
        for(year in years){
          YEARS.push(year);
          year = years[year];
          var tracker = document.getElementById('pr_annual_moods_graphs');
          tracker.appendChild(year);
        }
        document.getElementById(`graph.${YEARS[YEARS.length-1]}`).classList.add('activeYear');
      }
      break;
    case "chat:ready":
      if (msg.flag) { // chat connection with partner established
        CHATID = msg.cid;
        document.getElementById("chatbox")
            .addEventListener("keyup", function(event) {
            event.preventDefault();
            if (event.keyCode === 13) {
                document.getElementById("chatbox__send").click();
            }
        });
        if(document.getElementById('chat_flow_1')
          .classList.contains("shown"))
          done('chat_flow_1');
        done('text__loading');
        show('chat_flow_2');
      } else { // waiting on another user to start chat
        done('chat_flow_1');
        undone('text__loading');
        document.getElementById('text__loading').innerText = "Finding you someone to talk to";
      }
      break;
    case "chat:message":
      let chatlog = document.getElementById('chatlog');
      let newMessage = document.createElement('p');
      newMessage.classList.add(msg.flag ? "otherUser" : "user");
      newMessage.innerHTML = `${msg.flag ? "Peer: " : ""}${msg.data}`;
      chatlog.appendChild(newMessage);
      break;
    case "chat:rejected":
      if (confirm("Woah there! Are you sure that you're saying something nice? Remember, the other person is likely in a difficult place, much like you might be!")) {
        sock.send(JSON.stringify({
          type: "chat:verify",
          cid: CHATID,
          mid: msg.mid
        }));
      }
      break;
    case "chat:banned":
      if(document.getElementById('chat_flow_1')
      .classList.contains("shown"))
        done('chat_flow_1');
      done('text__loading');
      show('chat_flow_banned');
      break;
    case "chat:closed":
      toggle('chat_flow_2');
      show('chat_flow_end');
  }
}

function mood(mood){
  let date = new Date();
  let month = date.getMonth();
  let day = date.getDay();
  let year = date.getUTCFullYear();
  let json = {
    type: 'mood',
    day: day,
    month: month,
    year: year,
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
  let json = {
    type: 'comment',
    mood: 1,
    data: comment
  }
  sock.send(JSON.stringify(json));
  mood_continue();
}
function happy_submit() {
  let comment = document.getElementById('happy_comment').value;
  let json = {
    type: 'comment',
    mood: 1,
    data: comment
  }
  sock.send(JSON.stringify(json));
  mood_continue();
}
function neutral_submit() {
  let comment = document.getElementById('neutral_comment').value;
  let json = {
    type: 'comment',
    mood: 0,
    data: comment
  }
  sock.send(JSON.stringify(json));
  mood_continue();
}
function negative_submit() {
  let comment = document.getElementById('negative_comment').value;
  let json = {
    type: 'comment',
    mood: -1,
    data: comment
  }
  sock.send(JSON.stringify(json));
  mood_continue();
}
function danger_submit() {
  let comment = document.getElementById('danger_comment').value;
  let json = {
    type: 'comment',
    mood: -1,
    data: comment
  }
  sock.send(JSON.stringify(json));
  mood_continue();
}
function mood_continue() {
  done(`mood__${MOOD}`);
  show('block__chat');
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

function setState() {
  if (window.innerWidth >= 800) {
    show('menu');
  }
}

function deleteData() {
  let json = {
    type: 'delete'
  }
  sock.send(JSON.stringify(json));
}

function startChat() {
  if (document.getElementById('chat_flow_end').classList.contains('shown'))
  document.getElementById('chat_flow_end').classList.remove('shown')
  let json = {
    type: "chat:start"
  }
  sock.send(JSON.stringify(json));
}
function sendChatMsg() {
  let chatbox = document.getElementById('chatbox');
  let textToSend = chatbox.value;
  let json = {
    type: "chat:send",
    cid: CHATID,
    data: textToSend
  }
  sock.send(JSON.stringify(json));
  chatbox.value = "";
}
function sendChatReport() {
  let json = {
    type: "chat:report",
    cid: CHATID
  }
  sock.send(JSON.stringify(json));
}
function endChat() {
  let json = {
    type: "chat:close"
  }
  sock.send(JSON.stringify(json));
}