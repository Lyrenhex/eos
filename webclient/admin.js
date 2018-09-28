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
var serverAddr = `ws://${window.location.hostname}:9874`;

var reportID;

var sock = new WebSocket(serverAddr);
function signin_email() {
  var email = document.getElementById('login_email').value;
  var passw = document.getElementById('login_password').value;

  sock.send(JSON.stringify({'type':'login', 'emailAddress':email, 'password': passw}))

  return false;
}
sock.onerror = function(e) {
  alert(e);
}
sock.onclose = function(e) {
  alert("Server connection lost.");
}
sock.onmessage = function(e) {
  let msg = JSON.parse(e.data);
  console.log(msg);
  switch(msg.type){
    case 'login':
      if (msg.flag) {
        if (msg.user.Admin) {
          section('report');
          let reportUri = new URL(window.location.href);
          reportID = reportUri.searchParams.get('id');
          if (reportID != "") {
            sock.send(JSON.stringify({
              type: "admin:access",
              cid: reportID
            }));
            update_var('reportID', reportID);
          } else
            alert("No reportID.");
        } else 
          section('permerr');
      }
      break;
    case 'admin:chatlog':
      let chatlog = document.getElementById('report_log');
      msg.chatlog.forEach((message, key) => {
        let newMsg = document.createElement('p');
        newMsg.innerHTML = `[${message.aiDecision ? "Allowed" : "Rejected"}] <span class="userid">${message.sender}</span>: <span class="message">${escapeHtml(message.message)}</span>`;
        chatlog.appendChild(newMsg);
      });
      break;
    case 'admin:success':
      section('success');
  }
}

function submit_flag() {
  sock.send(JSON.stringify({
    type: "admin:flag",
    cid: reportID
  }));
}

function submit_flag() {
  let reportUser = document.getElementById('report_malign').value;
  sock.send(JSON.stringify({
    type: "admin:decision",
    cid: reportID,
    data: reportUser
  }));
}