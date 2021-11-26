# Eos

Eos is a progressive web application for the purpose of mood tracking (in particular for people with depression).
Eos was designed based on a personal irritation with existing mood tracking solutions, many of which required identification
of specific emotions (eg, happy, sad, frustrated, etc), or provided scales which were too precise (at least for some
neurodivergent people) -- this includes the standardised 10-point scale used often in medical settings. Instead, Eos uses a
5-point scale centered on 0 ('neutral').

Eos is licensed under the **MIT License**.

This project is a self-contained, local PWA. The 'live' version of Eos runs on a GitHub Pages instance hosting the code within
this repository, but the application can run offline. Data is stored using the `LocalStorage` API found in all modern web browsers,
though there is infrastructure to support moving over to a different database technology should the need arise (such as `IndexedDB`).

To use Eos, either visit [the Live instance](https://eos.lyrenhex.com), or open `app/index.html`.
