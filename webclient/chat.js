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

function add_msg(message, type) {
  let hist = document.getElementById('block__chat');
  let msgBlock = document.createElement('div');
  msgBlock.classList.add('msg', type);
  msgBlock.appendChild(
    document.createTextNode(message)
  );
  hist.appendChild(msgBlock);
}

// Program flow starts here
var config = {
  apiKey: "AIzaSyAR1zmwaotOgqX3EKEkUDPzM26FaujxEKY",
  authDomain: "solace-171915.firebaseapp.com",
  databaseURL: "https://solace-171915.firebaseio.com",
  projectId: "solace-171915",
  storageBucket: "solace-171915.appspot.com",
  messagingSenderId: "930975983513"
};
firebase.initializeApp(config); // initialise the Firebase web SDK

var CONFIG;

var confRef = firebase.database().ref(`/config/${branch}`);
confRef.on('value', (snapshot) => {
  CONFIG = snapshot.val();

  firebase.auth().onAuthStateChanged(function(user) {
    console.log(user);
    if (user && user.emailVerified) {
      function held_accept(id) {
        done('block__held');
        sock.send(JSON.stringify({'type': 'approve', 'mid': id, 'uid': user.uid}));
      }
      function held(message, id) {
        update_var('heldmsg', message);
        var btnAccept = document.getElementById('heldmsg_accept');
        btnAccept.onclick = function() { held_accept(id); }
        show('block__held');
      }
      function report() {
        done('block__report');
        sock.send(JSON.stringify({'type': 'report', 'uid': user.uid}));
        sock.close();
      }
      function exit() {
        sock.close();
        done('block__report');
      }
      //show('block__chat');
      var sock = new WebSocket(`wss://${CONFIG.websocket}`);
      sock.onerror = function(e){
        err(`Failed to connect to server.`)
      }
      sock.onclose = function(e) {
        console.log("Socket closed.");
        done('block__chat');
        show('block__closed');
      }
      sock.onmessage = function(e) {
        let msg = JSON.parse(e.data);
        console.log('received: ', msg);
        switch(msg.type){
          case 'init':
            done('text__loading');
            show('block__chat');
            document.getElementById('button__report').classList.add('loggedin');
            document.getElementById('btnReport').onclick = function() { report() };
            document.getElementById('btnExit').onclick = function() { exit() };
            break;
          case 'msg':
            add_msg(msg.text, 'anon');
            break;
          case 'hold':
            held(msg.text, msg.id);
            break;
        }
      }
      sock.onopen = function(e) {
        console.log("Connected to ", e.currentTarget.url);
        sock.send(JSON.stringify({'type': 'id', 'uid': user.uid}));
        document.getElementById('text__loading').textContent = 'Finding someone chatty...';
        window.addEventListener('beforeunload', function(e){
          sock.close();
          return null;
        });
        let input = document.getElementById('cb');
        input.addEventListener('keydown', (event) => {
          if(event.key === "Enter") {
            sock.send(JSON.stringify({'type': 'msg', 'text': input.value, 'uid': user.uid}));
            add_msg(input.value, 'user');
            input.value = "";
          }
        });
      }
    } else if (user) {
      document.getElementById('btnResendVerification').onclick = function() { user.sendEmailVerification() }
      show('block__unverified');
    } else {
      done('text__loading');
      show('block__loggedout');
    }
  }, (error) => {
    err(error);
  });
});
