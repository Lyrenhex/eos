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

var USER;
var DB;
var MOOD;
var MOOD_TIME_DATA = {
nums: [],
dates: [],
months: []
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

// sign-in methods.
function signin_google() {
    var provider = new firebase.auth.GoogleAuthProvider();
    firebase.auth().signInWithRedirect(provider);
}
function signin_email() {
    var email = document.getElementById('login_email').value;
    var passw = document.getElementById('login_password').value;

    firebase.auth().signInWithEmailAndPassword(email, passw).catch(function(error){
        if (error.code === "auth/user-not-found"){
        firebase.auth().createUserWithEmailAndPassword(email, passw).catch(function(error){
            if (error.code === "auth/invalid-email") {
            alert("Invalid email address.");
            } else if (error.code === "auth/weak-password") {
            alert("Password not strong enough.");
            }
        });
        }
    });

    return false;
}

firebase.auth().onAuthStateChanged(function(user) {
    if (user) {
        update_var('id', `${user.uid}`);
        USER = user;
        DB = firebase.database();
        var positiveCommentsRef = DB.ref(`/users/${user.uid}/positives`);
        positiveCommentsRef.on('child_added', function(data) {
            let li = document.createElement('li');
            let li_text = document.createTextNode(data.val());
            li.appendChild(li_text);
            let ul = document.getElementById('comments_positive');
            ul.appendChild(li);
        });
        var neutralCommentsRef = DB.ref(`/users/${user.uid}/neutral`);
        neutralCommentsRef.on('child_added', function(data) {
            let li = document.createElement('li');
            let li_text = document.createTextNode(data.val());
            li.appendChild(li_text);
            let ul = document.getElementById('comments_neutral');
            ul.appendChild(li);
        });
        var negativeCommentsRef = DB.ref(`/users/${user.uid}/negatives`);
        negativeCommentsRef.on('child_added', function(data) {
            let li = document.createElement('li');
            let li_text = document.createTextNode(data.val());
            li.appendChild(li_text);
            let ul = document.getElementById('comments_negative');
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
                    year = years[year];
                    var tracker = document.getElementById('annual_moods_graphs');
                    tracker.appendChild(year);
                }
            }
        });
    } else {
        show('block__loggedout');
    }
}, (error) => {
    err(error);
});