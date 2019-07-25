# Security Policy

## Supported Versions

| Version | Supported          |
| ------- | ------------------ |
| 3.0.x   | TBD                |
| 2.1.x   | TBD                |
| 2.0.x   | :white_check_mark: |
| < 2.0   | :x:                |

## Supported Browsers

| Browser                  | Supported          |
| ------------------------ | ------------------ |
| Google Chrome >= 49      | :white_check_mark: |
| Google Chrome <= 48      | :x:                |
| Mozilla Firefox >= 31    | :white_check_mark: |
| Mozilla Firefox <= 31    | :x:                |
| Edge (EdgeHTML)          | :x:                |
| Internet Explorer        | :x:                |

If any browsers not included on here are based on a supported browser, please test the vulnerability with the supported browser.
The child browser is not, itself, officially supported (though it is assumed to work based on its ancestor).

## Reporting a Vulnerability

To report a security vulnerability, please email dh64784@gmail.com with details about the vulnerability.
This must include steps to reproduce and the impact of the vulnerability. If you have traced the origin of the vulnerability in the code,
this is also appreciated. If you have fixed a vulnerability, please *also* provide a diff to me over email, along with a list of changes 
and their explanations; I will then authorise whether this is an acceptable patch, and advise you on whether to create a pull request.
**Do not create a pull request unless asked to**.

However, I may publish some low-risk vulnerabilites in the Issue tracker, if I have assessed them to not be relatively safe to do so.
If this is the case, then a pull request may be created without contacting me first, or without describing the vulnerability (merely 
reference the Issue Number).

Note that any forks made to develop a security patch should be made private, **unless** there is a public bug report about the 
security vulnerability.

If asked to submit a pull request, please include:

- An explanation of the security vulnerability;
- A list of changes that have been made and explanations.

This information is expected within the email to me, so should be a simple copy-paste job.

Vulnerability reports will be assessed and responded to within a few days, though may take up to a week (dependant on how busy I am).
If I approve a vulnerability report, the vulnerability will be placed higher priority than bugfixes or feature updates, though will be
addressed in order of relative risk.

The following will not be considered security vulnerabilities that may be addressed, though this list is not exhaustive:

- Social engineering, including vulnerabilites arising from **user misuse** of the Developer Console.
- Vulnerabilites arising from other libraries, modules, or technologies used by the application; such issues should be reported to the
relevant project, not Eos. (This policy includes Browser bugs which are not compliant with the W3C specification.)
